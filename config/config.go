package config

import (
	"github.com/spf13/viper"
)

const (
	// DefaultConfigFile is the default path used to load the configuration.
	DefaultConfigFile    = "/etc/rdap-exporter/config.toml"
	defaultCheckInterval = 60
	defaultTimeout       = 30
	defaultPort          = 9099
)

// Config contains the configuration for rdap-exporter.
type Config struct {
	// Domains is a list of domains to scan.
	Domains []string `mapstructure:"domains"`
	// CheckInterval is the interval between scans.
	CheckInterval uint64 `mapstructure:"check_interval"`
	// Timeout is the timeout for rdap queries.
	Timeout uint64 `mapstructure:"timeout"`
	// ListenPort is the port on which rdap-exporter listens.
	ListenPort uint64 `mapstructure:"listen_port"`
}

// Load the configuration from file.
func Load(confFile string) (conf *Config, err error) {
	viper.SetConfigFile(confFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}
	conf = &Config{
		CheckInterval: defaultCheckInterval,
		Timeout:       defaultTimeout,
		ListenPort:    defaultPort,
	}
	if err := viper.Unmarshal(&conf); err != nil {
		return nil, err
	}
	return conf, nil
}
