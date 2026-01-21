// Package importer provides functionality to import existing configuration files
// and generate equivalent Go code.
package importer

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-observability-go/alertmanager"
)

// ParseAlertmanagerConfig parses an alertmanager.yml file and returns an AlertmanagerConfig.
func ParseAlertmanagerConfig(path string) (*alertmanager.AlertmanagerConfig, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseAlertmanagerConfigFromBytes(data)
}

// ParseAlertmanagerConfigFromBytes parses alertmanager.yml content from bytes.
func ParseAlertmanagerConfigFromBytes(data []byte) (*alertmanager.AlertmanagerConfig, error) {
	var config alertmanager.AlertmanagerConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &config, nil
}

// ValidateAlertmanagerConfig validates a parsed Alertmanager configuration.
func ValidateAlertmanagerConfig(config *alertmanager.AlertmanagerConfig) []string {
	var warnings []string

	// Check for required route
	if config.Route == nil {
		warnings = append(warnings, "missing route configuration")
	} else if config.Route.Receiver == "" {
		warnings = append(warnings, "root route missing receiver")
	}

	// Check for at least one receiver
	if len(config.Receivers) == 0 {
		warnings = append(warnings, "no receivers defined")
	}

	// Check that route receiver exists in receivers list
	if config.Route != nil && config.Route.Receiver != "" {
		found := false
		for _, r := range config.Receivers {
			if r.Name == config.Route.Receiver {
				found = true
				break
			}
		}
		if !found {
			warnings = append(warnings, fmt.Sprintf("route receiver %q not found in receivers list", config.Route.Receiver))
		}
	}

	return warnings
}
