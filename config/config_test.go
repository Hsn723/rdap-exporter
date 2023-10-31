package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoad(t *testing.T) {
	cases := []struct {
		title    string
		file     string
		expected *Config
		isErr    bool
	}{
		{
			title: "Defaults",
			file:  "../testdata/config.toml",
			expected: &Config{
				Domains: []string{
					"example.com",
					"example.net",
				},
				CheckInterval: defaultCheckInterval,
				Timeout:       defaultTimeout,
				ListenPort:    defaultPort,
			},
		},
		{
			title: "Full",
			file:  "../testdata/config_full.toml",
			expected: &Config{
				Domains: []string{
					"example.com",
					"example.net",
				},
				CheckInterval: 100,
				Timeout:       100,
				ListenPort:    9999,
			},
		},
	}
	for _, tc := range cases {
		t.Run(tc.title, func(t *testing.T) {
			t.Helper()
			actual, err := Load(tc.file)
			if tc.isErr {
				assert.Error(t, err)
				assert.Nil(t, actual)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expected, actual)
			}
		})
	}
}
