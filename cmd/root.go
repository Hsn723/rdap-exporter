package cmd

import (
	"context"
	"net"
	"net/http"
	"os"
	"strconv"

	"github.com/hsn723/rdap-exporter/collector"
	"github.com/hsn723/rdap-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promslog"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "rdap-exporter",
		Short: "rdap-exporter provides metrics for domain RDAP info.",
		RunE:  runRoot,
	}

	logger = promslog.New(&promslog.Config{})

	configFile    string
	webConfigFile string

	version string
	commit  string
	date    string
	builtBy string
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", config.DefaultConfigFile, "path to configuration file")
	rootCmd.Flags().StringVar(&webConfigFile, "web.config.file", "", "Path to configuration file that can enable TLS or authentication. See: https://github.com/prometheus/exporter-toolkit/blob/master/docs/web-configuration.md")
}

func runRoot(cmd *cobra.Command, _ []string) error {
	logger.Info("rdap-exporter", "version", version, "commit", commit, "date", date, "built_by", builtBy)
	conf, err := config.Load(configFile)
	if err != nil {
		return err
	}
	logger.Info("loaded configuration", "config", configFile)
	rdapExporter := collector.NewRdapExporter(*conf, logger)
	prometheus.MustRegister(rdapExporter)

	ctx, cancelFunc := context.WithCancel(cmd.Context())
	defer cancelFunc()
	go rdapExporter.StartMetricsCollection(ctx)

	http.Handle("/metrics", promhttp.Handler())
	lc := web.LandingConfig{
		Name:        "RDAP Exporter",
		Description: "Prometheus exporter for domain RDAP information",
		Version:     version,
		Links: []web.LandingLinks{
			{
				Address: "/metrics",
				Text:    "Metrics",
			},
		},
	}
	lp, err := web.NewLandingPage(lc)
	if err != nil {
		return err
	}
	http.Handle("/", lp)

	listenAddress := net.JoinHostPort("0.0.0.0", strconv.FormatUint(conf.ListenPort, 10))
	useSystemdSocket := false
	flags := &web.FlagConfig{
		WebListenAddresses: &[]string{listenAddress},
		WebSystemdSocket:   &useSystemdSocket,
		WebConfigFile:      &webConfigFile,
	}
	server := &http.Server{}
	return web.ListenAndServe(server, flags, logger)
}

// Execute runs the root command.
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logger.Error("failed to run rdap-exporter", "error", err)
		os.Exit(1)
	}
}
