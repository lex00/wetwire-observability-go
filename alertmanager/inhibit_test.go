package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewInhibitRule(t *testing.T) {
	rule := NewInhibitRule()
	if rule == nil {
		t.Error("NewInhibitRule() returned nil")
	}
}

func TestInhibitRule_FluentAPI(t *testing.T) {
	rule := NewInhibitRule().
		WithSourceMatchers(
			Eq("severity", "critical"),
		).
		WithTargetMatchers(
			Eq("severity", "warning"),
		).
		WithEqual("alertname", "cluster")

	if len(rule.SourceMatchers) != 1 {
		t.Errorf("len(SourceMatchers) = %d, want 1", len(rule.SourceMatchers))
	}
	if len(rule.TargetMatchers) != 1 {
		t.Errorf("len(TargetMatchers) = %d, want 1", len(rule.TargetMatchers))
	}
	if len(rule.Equal) != 2 {
		t.Errorf("len(Equal) = %d, want 2", len(rule.Equal))
	}
}

func TestInhibitRule_WithSourceMatch(t *testing.T) {
	rule := NewInhibitRule().
		WithSourceMatch(map[string]string{
			"severity": "critical",
		})

	if rule.SourceMatch["severity"] != "critical" {
		t.Errorf("SourceMatch[severity] = %v", rule.SourceMatch["severity"])
	}
}

func TestInhibitRule_WithTargetMatch(t *testing.T) {
	rule := NewInhibitRule().
		WithTargetMatch(map[string]string{
			"severity": "warning",
		})

	if rule.TargetMatch["severity"] != "warning" {
		t.Errorf("TargetMatch[severity] = %v", rule.TargetMatch["severity"])
	}
}

func TestInhibitRule_Serialize(t *testing.T) {
	rule := NewInhibitRule().
		WithSourceMatchers(
			Eq("severity", "critical"),
			Eq("alertname", "HighCPU"),
		).
		WithTargetMatchers(
			Eq("severity", "warning"),
		).
		WithEqual("alertname", "cluster", "namespace")

	data, err := yaml.Marshal(rule)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"source_matchers:",
		"severity=\"critical\"",
		"alertname=\"HighCPU\"",
		"target_matchers:",
		"severity=\"warning\"",
		"equal:",
		"- alertname",
		"- cluster",
		"- namespace",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestInhibitRule_SerializeWithMatch(t *testing.T) {
	rule := NewInhibitRule().
		WithSourceMatch(map[string]string{
			"severity": "critical",
		}).
		WithTargetMatch(map[string]string{
			"severity": "warning",
		}).
		WithEqual("alertname")

	data, err := yaml.Marshal(rule)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"source_match:",
		"severity: critical",
		"target_match:",
		"severity: warning",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestInhibitRule_Unmarshal(t *testing.T) {
	input := `
source_matchers:
  - severity="critical"
target_matchers:
  - severity="warning"
equal:
  - alertname
  - cluster
`
	var rule InhibitRule
	if err := yaml.Unmarshal([]byte(input), &rule); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(rule.SourceMatchers) != 1 {
		t.Errorf("len(SourceMatchers) = %d, want 1", len(rule.SourceMatchers))
	}
	if len(rule.TargetMatchers) != 1 {
		t.Errorf("len(TargetMatchers) = %d, want 1", len(rule.TargetMatchers))
	}
	if len(rule.Equal) != 2 {
		t.Errorf("len(Equal) = %d, want 2", len(rule.Equal))
	}
}

func TestCriticalInhibitsWarning(t *testing.T) {
	rule := CriticalInhibitsWarning()

	if len(rule.SourceMatchers) != 1 {
		t.Errorf("len(SourceMatchers) = %d, want 1", len(rule.SourceMatchers))
	}
	if len(rule.TargetMatchers) != 1 {
		t.Errorf("len(TargetMatchers) = %d, want 1", len(rule.TargetMatchers))
	}
	if len(rule.Equal) != 1 || rule.Equal[0] != "alertname" {
		t.Errorf("Equal = %v, want [alertname]", rule.Equal)
	}
}

func TestErrorInhibitsWarning(t *testing.T) {
	rule := ErrorInhibitsWarning()

	if len(rule.SourceMatchers) != 1 {
		t.Errorf("len(SourceMatchers) = %d, want 1", len(rule.SourceMatchers))
	}
	if len(rule.TargetMatchers) != 1 {
		t.Errorf("len(TargetMatchers) = %d, want 1", len(rule.TargetMatchers))
	}
}

func TestInhibitRuleInConfig(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithInhibitRules(
			CriticalInhibitsWarning(),
			ErrorInhibitsWarning(),
		)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "inhibit_rules:") {
		t.Errorf("yaml.Marshal() missing inhibit_rules:\nGot:\n%s", yamlStr)
	}
}
