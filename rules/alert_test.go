package rules

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewAlertingRule(t *testing.T) {
	rule := NewAlertingRule("HighCPU")
	if rule == nil {
		t.Error("NewAlertingRule() returned nil")
	}
	if rule.Alert != "HighCPU" {
		t.Errorf("Alert = %v, want HighCPU", rule.Alert)
	}
}

func TestAlertingRule_WithExpr(t *testing.T) {
	rule := NewAlertingRule("Test").WithExpr("cpu > 90")
	if rule.Expr != "cpu > 90" {
		t.Errorf("Expr = %v", rule.Expr)
	}
}

func TestAlertingRule_WithFor(t *testing.T) {
	rule := NewAlertingRule("Test").WithFor(5 * Minute)
	if rule.For != 5*Minute {
		t.Errorf("For = %v, want 5m", rule.For)
	}
}

func TestAlertingRule_WithLabels(t *testing.T) {
	rule := NewAlertingRule("Test").WithLabels(map[string]string{
		"severity": "critical",
		"team":     "platform",
	})
	if rule.Labels["severity"] != "critical" {
		t.Errorf("Labels[severity] = %v", rule.Labels["severity"])
	}
}

func TestAlertingRule_WithAnnotations(t *testing.T) {
	rule := NewAlertingRule("Test").WithAnnotations(map[string]string{
		"summary":     "Test alert",
		"description": "Detailed description",
	})
	if rule.Annotations["summary"] != "Test alert" {
		t.Errorf("Annotations[summary] = %v", rule.Annotations["summary"])
	}
}

func TestAlertingRule_SeverityHelpers(t *testing.T) {
	critical := NewAlertingRule("Critical").Critical()
	if critical.Labels["severity"] != "critical" {
		t.Errorf("Critical() severity = %v", critical.Labels["severity"])
	}

	warning := NewAlertingRule("Warning").Warning()
	if warning.Labels["severity"] != "warning" {
		t.Errorf("Warning() severity = %v", warning.Labels["severity"])
	}

	info := NewAlertingRule("Info").Info()
	if info.Labels["severity"] != "info" {
		t.Errorf("Info() severity = %v", info.Labels["severity"])
	}
}

func TestAlertingRule_AnnotationHelpers(t *testing.T) {
	rule := NewAlertingRule("Test").
		WithSummary("High CPU usage on {{ $labels.instance }}").
		WithDescription("CPU usage is {{ $value }}%").
		WithRunbook("https://runbook.example.com/high-cpu")

	if rule.Annotations["summary"] != "High CPU usage on {{ $labels.instance }}" {
		t.Errorf("summary = %v", rule.Annotations["summary"])
	}
	if rule.Annotations["description"] != "CPU usage is {{ $value }}%" {
		t.Errorf("description = %v", rule.Annotations["description"])
	}
	if rule.Annotations["runbook_url"] != "https://runbook.example.com/high-cpu" {
		t.Errorf("runbook_url = %v", rule.Annotations["runbook_url"])
	}
}

func TestAlertingRule_FluentAPI(t *testing.T) {
	rule := NewAlertingRule("HighCPU").
		WithExpr("avg(cpu_usage) > 90").
		WithFor(5 * Minute).
		Critical().
		WithSummary("High CPU usage").
		WithDescription("CPU is at {{ $value }}%").
		WithRunbook("https://example.com/runbook")

	if rule.Alert != "HighCPU" {
		t.Errorf("Alert = %v", rule.Alert)
	}
	if rule.Expr != "avg(cpu_usage) > 90" {
		t.Errorf("Expr = %v", rule.Expr)
	}
	if rule.For != 5*Minute {
		t.Errorf("For = %v", rule.For)
	}
	if rule.Labels["severity"] != "critical" {
		t.Errorf("severity = %v", rule.Labels["severity"])
	}
	if rule.Annotations["summary"] != "High CPU usage" {
		t.Errorf("summary = %v", rule.Annotations["summary"])
	}
}

func TestAlertingRule_Serialize(t *testing.T) {
	rule := NewAlertingRule("HighMemory").
		WithExpr("memory_usage > 80").
		WithFor(10 * Minute).
		Warning().
		WithSummary("High memory usage")

	data, err := yaml.Marshal(rule)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"alert: HighMemory",
		"expr: memory_usage > 80",
		"for: 10m",
		"labels:",
		"severity: warning",
		"annotations:",
		"summary: High memory usage",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestAlertingRule_Unmarshal(t *testing.T) {
	input := `
alert: TestAlert
expr: test > 0
for: 5m
labels:
  severity: critical
annotations:
  summary: Test summary
`
	var rule AlertingRule
	if err := yaml.Unmarshal([]byte(input), &rule); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if rule.Alert != "TestAlert" {
		t.Errorf("Alert = %v", rule.Alert)
	}
	if rule.Expr != "test > 0" {
		t.Errorf("Expr = %v", rule.Expr)
	}
	if rule.For != 5*Minute {
		t.Errorf("For = %v", rule.For)
	}
	if rule.Labels["severity"] != "critical" {
		t.Errorf("severity = %v", rule.Labels["severity"])
	}
}

func TestAlertingRule_KeepFiringFor(t *testing.T) {
	rule := NewAlertingRule("Test").
		WithExpr("test > 0").
		WithKeepFiringFor(15 * Minute)

	if rule.KeepFiringFor != 15*Minute {
		t.Errorf("KeepFiringFor = %v, want 15m", rule.KeepFiringFor)
	}
}
