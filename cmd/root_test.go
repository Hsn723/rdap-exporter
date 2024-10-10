package cmd

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/onsi/gomega"
	"github.com/stretchr/testify/assert"
)

func TestExecute(t *testing.T) {
	g := gomega.NewWithT(t)
	g.SetDefaultEventuallyTimeout(10 * time.Second)

	g.Expect(rootCmd.Flags().Set("config", "../testdata/config.toml")).To(gomega.Succeed())

	go Execute()

	g.Eventually(func() error {
		resp, err := http.Get("http://127.0.0.1:9099/")
		if err != nil {
			return err
		}
		defer resp.Body.Close()
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if !strings.Contains(string(data), "</html>") {
			return assert.AnError
		}
		return nil
	}).Should(gomega.Succeed())

	expectedMetrics := []string{
		"rdap_domain_event{domain=\"example.com\",event=\"expiration\"}",
		"rdap_domain_event{domain=\"example.net\",event=\"last_update_of_rdap_database\"}",
		"rdap_domain_status{domain=\"example.com\",status=\"client_delete_prohibited\"} 1",
		"rdap_domain_status{domain=\"example.net\",status=\"client_transfer_prohibited\"} 1",
	}
	g.Eventually(func() error {
		metrics, err := http.Get("http://127.0.0.1:9099/metrics")
		if err != nil {
			return err
		}
		defer metrics.Body.Close()
		rawData, err := io.ReadAll(metrics.Body)
		if err != nil {
			return err
		}
		data := string(rawData)
		for _, expect := range expectedMetrics {
			if !strings.Contains(data, expect) {
				return fmt.Errorf("missing output %s", expect)
			}
		}
		return nil
	}).Should(gomega.Succeed())
}
