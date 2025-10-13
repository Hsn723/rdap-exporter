package collector

import (
	"context"
	"log/slog"
	"net/url"
	"strings"
	"time"

	"github.com/hsn723/rdap-exporter/config"
	"github.com/openrdap/rdap"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "rdap"

	// Error labels
	ErrServerURLParse  = "rdap_server_url_parse_error"
	ErrNoInfo          = "rdap_no_info"
	ErrNotDomain       = "rdap_response_not_domain"
	ErrWrongDateFormat = "rdap_wrong_date_format"
)

// RdapExporter implements a Prometheus Collector.
type RdapExporter struct {
	config         config.Config
	domainStatuses *prometheus.GaugeVec
	domainEvents   *prometheus.GaugeVec
	domainErrors   *prometheus.CounterVec
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
		domainErrors: prometheus.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Name:      "domain_error",
			Help:      "Count of errors encountered per domain and error type.",
		}, []string{"domain", "error"}),
		logger: logger,
	}
}

// Describe the Prometheus metrics being exported.
func (e *RdapExporter) Describe(ch chan<- *prometheus.Desc) {
	e.domainStatuses.Describe(ch)
	e.domainEvents.Describe(ch)
	e.domainErrors.Describe(ch)
}

// Collect metrics.
func (e *RdapExporter) Collect(ch chan<- prometheus.Metric) {
	e.domainStatuses.Collect(ch)
	e.domainEvents.Collect(ch)
	e.domainErrors.Collect(ch)
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

func collectRdapInfo(ctx context.Context, e *RdapExporter, domain config.Domain) {
	req := &rdap.Request{
		Type:  rdap.DomainRequest,
		Query: domain.Name,
	}
	if domain.RdapServerUrl != "" {
		RdapServerUrl, err := url.Parse(domain.RdapServerUrl)
		if err != nil {
			e.domainErrors.WithLabelValues(domain.Name, ErrServerURLParse).Inc()
			e.logger.Error("could not parse RdapServerUrl", "error", err, "domain", domain.Name, "rdap_server_url", domain.RdapServerUrl)
			return
		}
		req.Server = RdapServerUrl
	}

	req = req.WithContext(ctx)
	req.Timeout = time.Duration(e.config.Timeout) * time.Second
	client := &rdap.Client{}
	resp, err := client.Do(req)
	if err != nil {
		e.domainErrors.WithLabelValues(domain.Name, ErrNoInfo).Inc()
		e.logger.Error("could not get RDAP info", "error", err, "domain", domain.Name)
		return
	}
	data, ok := resp.Object.(*rdap.Domain)
	if !ok {
		e.domainErrors.WithLabelValues(domain.Name, ErrNotDomain).Inc()
		e.logger.Error("RDAP response is not a domain object", "domain", domain.Name)
		return
	}
	for _, rawStatus := range data.Status {
		status := normalizeLabel(rawStatus)
		e.domainStatuses.WithLabelValues(domain.Name, status).Set(1)
	}
	for _, event := range data.Events {
		date, err := time.Parse(time.RFC3339, event.Date)
		if err != nil {
			e.domainErrors.WithLabelValues(domain.Name, ErrWrongDateFormat).Inc()
			e.logger.Error("wrong date format", "error", err, "domain", domain.Name, "event", event.Action)
		}
		action := normalizeLabel(event.Action)
		e.domainEvents.WithLabelValues(domain.Name, action).Set(float64(date.Unix()))
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
