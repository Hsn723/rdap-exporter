package collector

import (
	"context"
	"log/slog"
	"strings"
	"time"

	"github.com/hsn723/rdap-exporter/config"
	"github.com/openrdap/rdap"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "rdap"
)

// RdapExporter implements a Prometheus Collector.
type RdapExporter struct {
	config         config.Config
	domainStatuses *prometheus.GaugeVec
	domainEvents   *prometheus.GaugeVec
	logger         *slog.Logger
}

// NewRdapExporter creates a new RdapExporter instance.
func NewRdapExporter(config config.Config, logger *slog.Logger) *RdapExporter {
	return &RdapExporter{
		config: config,
		domainStatuses: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "domain_status",
			Help:      "Domain status codes.",
		}, []string{"domain", "status"}),
		domainEvents: prometheus.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "domain_event",
			Help:      "Dates pertaining to the domain as a unix timestamp.",
		}, []string{"domain", "event"}),
		logger: logger,
	}
}

// Describe the Prometheus metrics being exported.
func (e *RdapExporter) Describe(ch chan<- *prometheus.Desc) {
	e.domainStatuses.Describe(ch)
	e.domainEvents.Describe(ch)
}

// Collect metrics.
func (e *RdapExporter) Collect(ch chan<- prometheus.Metric) {
	e.domainStatuses.Collect(ch)
	e.domainEvents.Collect(ch)
}

// StartMetricsCollection starts the metrics collection.
func (e *RdapExporter) StartMetricsCollection(ctx context.Context) {
	collectDomainRdapInfo(ctx, e)
	ticker := time.NewTicker(time.Duration(e.config.CheckInterval) * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		collectDomainRdapInfo(ctx, e)
	}
}

func collectRdapInfo(ctx context.Context, e *RdapExporter, domain string) {
	req := &rdap.Request{
		Type:  rdap.DomainRequest,
		Query: domain,
	}
	req = req.WithContext(ctx)
	req.Timeout = time.Duration(e.config.Timeout) * time.Second
	client := &rdap.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e.logger.Error("could not get RDAP info", "error", err, "domain", domain)
		return
	}
	data, ok := resp.Object.(*rdap.Domain)
	if !ok {
		e.logger.Error("RDAP response is not a domain object", "domain", domain)
		return
	}
	for _, rawStatus := range data.Status {
		status := normalizeLabel(rawStatus)
		e.domainStatuses.WithLabelValues(domain, status).Set(1)
	}
	for _, event := range data.Events {
		date, err := time.Parse(time.RFC3339, event.Date)
		if err != nil {
			e.logger.Error("wrong date format", "error", err, "domain", domain, "event", event.Action)
		}
		action := normalizeLabel(event.Action)
		e.domainEvents.WithLabelValues(domain, action).Set(float64(date.Unix()))
	}
}

func collectDomainRdapInfo(ctx context.Context, e *RdapExporter) {
	for _, domain := range e.config.Domains {
		go collectRdapInfo(ctx, e, domain)
	}
}

func normalizeLabel(s string) string {
	return strings.ToLower(strings.ReplaceAll(s, " ", "_"))
}
