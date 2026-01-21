// Package importer provides functionality to import existing configuration files
// and generate equivalent Go code.
package importer

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-observability-go/prometheus"
	"github.com/lex00/wetwire-observability-go/rules"
)

// RulesFile represents a Prometheus rules file for importing.
type RulesFile struct {
	Groups []RuleGroupInput `yaml:"groups"`
}

// RuleGroupInput represents a rule group in the input YAML.
type RuleGroupInput struct {
	Name     string         `yaml:"name"`
	Interval string         `yaml:"interval,omitempty"`
	Limit    int            `yaml:"limit,omitempty"`
	Rules    []RuleInput    `yaml:"rules"`
}

// RuleInput represents either an alerting or recording rule.
type RuleInput struct {
	// Alerting rule fields
	Alert         string            `yaml:"alert,omitempty"`
	For           string            `yaml:"for,omitempty"`
	KeepFiringFor string            `yaml:"keep_firing_for,omitempty"`

	// Recording rule fields
	Record string `yaml:"record,omitempty"`

	// Common fields
	Expr        string            `yaml:"expr"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// IsAlertingRule returns true if this is an alerting rule.
func (r *RuleInput) IsAlertingRule() bool {
	return r.Alert != ""
}

// IsRecordingRule returns true if this is a recording rule.
func (r *RuleInput) IsRecordingRule() bool {
	return r.Record != ""
}

// ParseRulesFile parses a Prometheus rules file.
func ParseRulesFile(path string) (*RulesFile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseRulesFileFromBytes(data)
}

// ParseRulesFileFromBytes parses Prometheus rules from YAML bytes.
func ParseRulesFileFromBytes(data []byte) (*RulesFile, error) {
	var rulesFile RulesFile
	if err := yaml.Unmarshal(data, &rulesFile); err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}

	return &rulesFile, nil
}

// ValidateRulesFile validates a parsed rules file.
func ValidateRulesFile(rf *RulesFile) []string {
	var warnings []string

	if len(rf.Groups) == 0 {
		warnings = append(warnings, "rules file has no groups")
	}

	for _, group := range rf.Groups {
		if group.Name == "" {
			warnings = append(warnings, "rule group has no name")
		}

		if len(group.Rules) == 0 {
			warnings = append(warnings, fmt.Sprintf("rule group %q has no rules", group.Name))
		}

		for i, rule := range group.Rules {
			if !rule.IsAlertingRule() && !rule.IsRecordingRule() {
				warnings = append(warnings, fmt.Sprintf("rule %d in group %q has neither 'alert' nor 'record' field", i, group.Name))
			}

			if rule.Expr == "" {
				warnings = append(warnings, fmt.Sprintf("rule %d in group %q has no expression", i, group.Name))
			}

			// Check alerting rule specific issues
			if rule.IsAlertingRule() {
				if rule.For == "" {
					warnings = append(warnings, fmt.Sprintf("alerting rule %q in group %q has no 'for' duration", rule.Alert, group.Name))
				}

				if _, ok := rule.Labels["severity"]; !ok {
					warnings = append(warnings, fmt.Sprintf("alerting rule %q in group %q has no severity label", rule.Alert, group.Name))
				}
			}
		}
	}

	return warnings
}

// ConvertToWetwireRules converts parsed rules to wetwire rules types.
func ConvertToWetwireRules(rf *RulesFile) *rules.RulesFile {
	result := rules.NewRulesFile()

	for _, group := range rf.Groups {
		rg := rules.NewRuleGroup(group.Name)

		if group.Interval != "" {
			if dur := parseDuration(group.Interval); dur != 0 {
				rg.WithInterval(prometheus.Duration(dur))
			}
		}

		if group.Limit > 0 {
			rg.WithLimit(group.Limit)
		}

		for _, rule := range group.Rules {
			if rule.IsAlertingRule() {
				ar := rules.NewAlertingRule(rule.Alert).
					WithExpr(rule.Expr)

				if rule.For != "" {
					if dur := parseDuration(rule.For); dur != 0 {
						ar.WithFor(prometheus.Duration(dur))
					}
				}

				if rule.KeepFiringFor != "" {
					if dur := parseDuration(rule.KeepFiringFor); dur != 0 {
						ar.WithKeepFiringFor(prometheus.Duration(dur))
					}
				}

				if len(rule.Labels) > 0 {
					ar.WithLabels(rule.Labels)
				}

				if len(rule.Annotations) > 0 {
					ar.WithAnnotations(rule.Annotations)
				}

				rg.AddRule(ar)
			} else if rule.IsRecordingRule() {
				rr := rules.NewRecordingRule(rule.Record).
					WithExpr(rule.Expr)

				if len(rule.Labels) > 0 {
					rr.WithLabels(rule.Labels)
				}

				rg.AddRule(rr)
			}
		}

		result.AddGroup(rg)
	}

	return result
}

// parseDuration parses a Prometheus-style duration string.
func parseDuration(s string) time.Duration {
	if s == "" {
		return 0
	}

	// Try standard Go duration first
	d, err := time.ParseDuration(s)
	if err == nil {
		return d
	}

	// Handle Prometheus-style durations (e.g., "5m", "1h", "30s")
	// This is a simplified parser; Go's ParseDuration handles most cases
	return 0
}
