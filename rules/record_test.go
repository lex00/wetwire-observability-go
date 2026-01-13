package rules

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewRecordingRule(t *testing.T) {
	rule := NewRecordingRule("cpu:avg")
	if rule == nil {
		t.Error("NewRecordingRule() returned nil")
	}
	if rule.Record != "cpu:avg" {
		t.Errorf("Record = %v, want cpu:avg", rule.Record)
	}
}

func TestRecordingRule_WithExpr(t *testing.T) {
	rule := NewRecordingRule("test").WithExpr("avg(cpu)")
	if rule.Expr != "avg(cpu)" {
		t.Errorf("Expr = %v", rule.Expr)
	}
}

func TestRecordingRule_WithLabels(t *testing.T) {
	rule := NewRecordingRule("test").WithLabels(map[string]string{
		"env":  "production",
		"tier": "backend",
	})
	if rule.Labels["env"] != "production" {
		t.Errorf("Labels[env] = %v", rule.Labels["env"])
	}
}

func TestRecordingRule_FluentAPI(t *testing.T) {
	rule := NewRecordingRule("job:http_requests:rate5m").
		WithExpr("sum(rate(http_requests_total[5m])) by (job)").
		WithLabels(map[string]string{
			"source": "prometheus",
		})

	if rule.Record != "job:http_requests:rate5m" {
		t.Errorf("Record = %v", rule.Record)
	}
	if rule.Expr != "sum(rate(http_requests_total[5m])) by (job)" {
		t.Errorf("Expr = %v", rule.Expr)
	}
	if rule.Labels["source"] != "prometheus" {
		t.Errorf("Labels[source] = %v", rule.Labels["source"])
	}
}

func TestRecordingRule_Serialize(t *testing.T) {
	rule := NewRecordingRule("instance:memory:usage").
		WithExpr("1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)").
		WithLabels(map[string]string{"type": "memory"})

	data, err := yaml.Marshal(rule)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"record: instance:memory:usage",
		"expr: 1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes)",
		"labels:",
		"type: memory",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRecordingRule_Unmarshal(t *testing.T) {
	input := `
record: test:metric
expr: sum(rate(requests[5m]))
labels:
  source: prometheus
`
	var rule RecordingRule
	if err := yaml.Unmarshal([]byte(input), &rule); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if rule.Record != "test:metric" {
		t.Errorf("Record = %v", rule.Record)
	}
	if rule.Expr != "sum(rate(requests[5m]))" {
		t.Errorf("Expr = %v", rule.Expr)
	}
	if rule.Labels["source"] != "prometheus" {
		t.Errorf("Labels[source] = %v", rule.Labels["source"])
	}
}

func TestRecordingRuleNamingConvention(t *testing.T) {
	// Recording rule names should follow the convention:
	// level:metric:operations
	// e.g., job:http_requests:rate5m

	validNames := []string{
		"job:http_requests:rate5m",
		"instance:memory:usage",
		"cluster:cpu:avg",
	}

	for _, name := range validNames {
		rule := NewRecordingRule(name)
		if rule.Record != name {
			t.Errorf("Record = %v, want %v", rule.Record, name)
		}
	}
}
