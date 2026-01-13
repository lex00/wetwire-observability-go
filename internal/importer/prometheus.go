// Package importer provides functionality to import existing configuration files
// and generate equivalent Go code.
package importer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-observability-go/prometheus"
)

// ParsePrometheusConfig parses a prometheus.yml file and returns a PrometheusConfig.
func ParsePrometheusConfig(path string) (*prometheus.PrometheusConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParsePrometheusConfigFromBytes(data)
}

// ParsePrometheusConfigFromBytes parses prometheus.yml content from bytes.
func ParsePrometheusConfigFromBytes(data []byte) (*prometheus.PrometheusConfig, error) {
	var config prometheus.PrometheusConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// ValidatePrometheusConfig validates a parsed Prometheus configuration.
func ValidatePrometheusConfig(config *prometheus.PrometheusConfig) []string {
	var warnings []string

	// Check for features that may need manual adjustment
	for _, sc := range config.ScrapeConfigs {
		if sc.JobName == "" {
			warnings = append(warnings, "scrape_config missing job_name")
		}

		// Check for service discovery types
		if len(sc.KubernetesSDConfigs) > 0 {
			for i, k := range sc.KubernetesSDConfigs {
				if k.Role == "" {
					warnings = append(warnings, fmt.Sprintf("scrape_config[%s].kubernetes_sd_configs[%d] missing role", sc.JobName, i))
				}
			}
		}
	}

	return warnings
}
