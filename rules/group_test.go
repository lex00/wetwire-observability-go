package rules

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewRuleGroup(t *testing.T) {
	group := NewRuleGroup("test-alerts")
	if group == nil {
		t.Error("NewRuleGroup() returned nil")
	}
	if group.Name != "test-alerts" {
		t.Errorf("Name = %v, want test-alerts", group.Name)
	}
}

func TestRuleGroup_WithInterval(t *testing.T) {
	group := NewRuleGroup("test").WithInterval(30 * Second)
	if group.Interval != 30*Second {
		t.Errorf("Interval = %v, want 30s", group.Interval)
	}
}

func TestRuleGroup_WithLimit(t *testing.T) {
	group := NewRuleGroup("test").WithLimit(100)
	if group.Limit != 100 {
		t.Errorf("Limit = %v, want 100", group.Limit)
	}
}

func TestRuleGroup_WithRules(t *testing.T) {
	group := NewRuleGroup("test").WithRules(
		NewAlertingRule("HighCPU").WithExpr("cpu > 90"),
		NewRecordingRule("cpu:avg").WithExpr("avg(cpu)"),
	)

	if len(group.Rules) != 2 {
		t.Errorf("len(Rules) = %d, want 2", len(group.Rules))
	}
}

func TestRuleGroup_AddRule(t *testing.T) {
	group := NewRuleGroup("test")
	group.AddRule(NewAlertingRule("Alert1").WithExpr("expr1"))
	group.AddRule(NewAlertingRule("Alert2").WithExpr("expr2"))

	if len(group.Rules) != 2 {
		t.Errorf("len(Rules) = %d, want 2", len(group.Rules))
	}
}

func TestRuleGroup_Serialize(t *testing.T) {
	group := NewRuleGroup("cpu-alerts").
		WithInterval(1 * Minute).
		WithRules(
			NewAlertingRule("HighCPU").
				WithExpr("cpu_usage > 90").
				WithFor(5 * Minute).
				WithLabels(map[string]string{"severity": "critical"}).
				WithAnnotations(map[string]string{"summary": "High CPU usage"}),
		)

	data, err := yaml.Marshal(group)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"name: cpu-alerts",
		"interval: 1m",
		"rules:",
		"alert: HighCPU",
		"expr: cpu_usage > 90",
		"for: 5m",
		"labels:",
		"severity: critical",
		"annotations:",
		"summary: High CPU usage",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRuleGroup_SerializeWithRecording(t *testing.T) {
	group := NewRuleGroup("recording").
		WithRules(
			NewRecordingRule("job:cpu:avg").
				WithExpr("avg by (job) (cpu_usage)").
				WithLabels(map[string]string{"env": "production"}),
		)

	data, err := yaml.Marshal(group)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"name: recording",
		"rules:",
		"record: job:cpu:avg",
		"expr: avg by (job) (cpu_usage)",
		"labels:",
		"env: production",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRuleGroup_Unmarshal(t *testing.T) {
	input := `
name: test-group
interval: 30s
rules:
  - alert: HighCPU
    expr: cpu > 90
    for: 5m
    labels:
      severity: critical
  - record: cpu:avg
    expr: avg(cpu)
`
	var group RuleGroup
	if err := yaml.Unmarshal([]byte(input), &group); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if group.Name != "test-group" {
		t.Errorf("Name = %v", group.Name)
	}
	if group.Interval != 30*Second {
		t.Errorf("Interval = %v", group.Interval)
	}
	if len(group.Rules) != 2 {
		t.Errorf("len(Rules) = %d, want 2", len(group.Rules))
	}

	// Rules are unmarshaled as map[string]any
	alertRule := group.Rules[0].(map[string]any)
	if alertRule["alert"] != "HighCPU" {
		t.Errorf("alert rule name = %v", alertRule["alert"])
	}

	recordRule := group.Rules[1].(map[string]any)
	if recordRule["record"] != "cpu:avg" {
		t.Errorf("record rule name = %v", recordRule["record"])
	}
}

func TestNewRulesFile(t *testing.T) {
	file := NewRulesFile()
	if file == nil {
		t.Error("NewRulesFile() returned nil")
	}
}

func TestRulesFile_WithGroups(t *testing.T) {
	file := NewRulesFile().WithGroups(
		NewRuleGroup("alerts"),
		NewRuleGroup("recording"),
	)

	if len(file.Groups) != 2 {
		t.Errorf("len(Groups) = %d, want 2", len(file.Groups))
	}
}

func TestRulesFile_Serialize(t *testing.T) {
	file := NewRulesFile().WithGroups(
		NewRuleGroup("alerts").WithRules(
			NewAlertingRule("Test").WithExpr("test > 0"),
		),
	)

	data, err := yaml.Marshal(file)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "groups:") {
		t.Errorf("yaml.Marshal() missing groups:\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "name: alerts") {
		t.Errorf("yaml.Marshal() missing name:\nGot:\n%s", yamlStr)
	}
}
