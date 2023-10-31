package cmd

import (
	"context"
	"log/slog"
	"net"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hsn723/rdap-exporter/collector"
	"github.com/hsn723/rdap-exporter/config"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "rdap-exporter",
		Short: "rdap-exporter provides metrics for domain RDAP info.",
		RunE:  runRoot,
	}

	configFile string

	version string
	commit  string
	date    string
	builtBy string
)

func init() {
	rootCmd.Flags().StringVarP(&configFile, "config", "c", config.DefaultConfigFile, "path to configuration file")
}

func runRoot(cmd *cobra.Command, _ []string) error {
	slog.Info("rdap-exporter", "version", version, "commit", commit, "date", date, "built_by", builtBy)
	conf, err := config.Load(configFile)
	if err != nil {
		return err
	}
	slog.Info("loaded configuration", "config", configFile)
	rdapExporter := collector.NewRdapExporter(*conf)
	prometheus.MustRegister(rdapExporter)

	ctx, cancelFunc := context.WithTimeout(cmd.Context(), time.Duration(conf.Timeout)*time.Second)
	defer cancelFunc()
	go rdapExporter.StartMetricsCollection(ctx)

	http.Handle("/metrics", promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte(`
			<html>
			<head><title>RDAP Exporter</title></head>
			<body>
			<h1>RDAP Exporter</h1>
			<p><a href='metrics'>Metrics</a></p>
			</body>
			</html>
		`))
		if err != nil {
			slog.Error("HTTP write failed", "error", err)
		}
	})

	listenAddress := net.JoinHostPort("0.0.0.0", strconv.FormatUint(conf.ListenPort, 10))
	slog.Info("start listening server", "listen_address", listenAddress)
	return http.ListenAndServe(listenAddress, nil)
}

// Execute runs the root command.
func Execute() {
	ctx := context.Background()
	if err := rootCmd.ExecuteContext(ctx); err != nil {
		slog.Error("failed to run rdap-exporter", "error", err)
		os.Exit(1)
	}
}
