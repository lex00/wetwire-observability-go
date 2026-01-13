package rules

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRulesFile_SerializeMethod(t *testing.T) {
	file := NewRulesFile().WithGroups(
		NewRuleGroup("test-alerts").
			WithInterval(30*Second).
			WithRules(
				NewAlertingRule("HighCPU").
					WithExpr("cpu_usage > 90").
					WithFor(5*Minute).
					Critical().
					WithSummary("CPU usage is high"),
			),
	)

	data, err := file.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yaml := string(data)
	if !strings.Contains(yaml, "groups:") {
		t.Error("Expected 'groups:' in output")
	}
	if !strings.Contains(yaml, "name: test-alerts") {
		t.Error("Expected 'name: test-alerts' in output")
	}
	if !strings.Contains(yaml, "interval: 30s") {
		t.Error("Expected 'interval: 30s' in output")
	}
	if !strings.Contains(yaml, "alert: HighCPU") {
		t.Error("Expected 'alert: HighCPU' in output")
	}
	if !strings.Contains(yaml, "expr: cpu_usage > 90") {
		t.Error("Expected 'expr: cpu_usage > 90' in output")
	}
	if !strings.Contains(yaml, "for: 5m") {
		t.Error("Expected 'for: 5m' in output")
	}
	if !strings.Contains(yaml, "severity: critical") {
		t.Error("Expected 'severity: critical' in output")
	}
}

func TestRulesFile_SerializeRecording(t *testing.T) {
	file := NewRulesFile().WithGroups(
		NewRuleGroup("recording-rules").
			WithRules(
				NewRecordingRule("job:http_requests:rate5m").
					WithExpr("sum(rate(http_requests_total[5m])) by (job)"),
			),
	)

	data, err := file.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yaml := string(data)
	if !strings.Contains(yaml, "record: job:http_requests:rate5m") {
		t.Error("Expected 'record: job:http_requests:rate5m' in output")
	}
	if !strings.Contains(yaml, "expr: sum(rate(http_requests_total[5m])) by (job)") {
		t.Error("Expected expr in output")
	}
}

func TestRulesFile_SerializeMixed(t *testing.T) {
	file := NewRulesFile().WithGroups(
		NewRuleGroup("api-rules").
			WithRules(
				NewRecordingRule("api:http_requests:rate5m").
					WithExpr("sum(rate(http_requests_total[5m])) by (service)"),
				NewAlertingRule("HighErrorRate").
					WithExpr("api:http_errors:rate5m / api:http_requests:rate5m > 0.05").
					WithFor(5*Minute).
					Warning(),
			),
	)

	data, err := file.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yaml := string(data)
	if !strings.Contains(yaml, "record:") {
		t.Error("Expected recording rule in output")
	}
	if !strings.Contains(yaml, "alert:") {
		t.Error("Expected alerting rule in output")
	}
}

func TestRulesFile_SerializeToFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "rules.yml")

	file := NewRulesFile().WithGroups(
		NewRuleGroup("test").
			WithRules(
				NewAlertingRule("TestAlert").
					WithExpr("up == 0").
					WithFor(1*Minute),
			),
	)

	err := file.SerializeToFile(path)
	if err != nil {
		t.Fatalf("SerializeToFile() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	yaml := string(data)
	if !strings.Contains(yaml, "alert: TestAlert") {
		t.Error("Expected 'alert: TestAlert' in file")
	}
}

func TestRulesFile_MustSerialize(t *testing.T) {
	file := NewRulesFile().WithGroups(
		NewRuleGroup("test").
			WithRules(
				NewAlertingRule("TestAlert").
					WithExpr("up == 0"),
			),
	)

	// Should not panic
	data := file.MustSerialize()
	if len(data) == 0 {
		t.Error("MustSerialize() returned empty data")
	}
}

func TestRuleGroup_SerializeStandalone(t *testing.T) {
	group := NewRuleGroup("standalone-group").
		WithInterval(1*Minute).
		WithRules(
			NewAlertingRule("Alert1").WithExpr("expr1"),
			NewAlertingRule("Alert2").WithExpr("expr2"),
		)

	data, err := group.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yaml := string(data)
	// When serializing a single group, it should wrap in groups array
	if !strings.Contains(yaml, "groups:") {
		t.Error("Expected 'groups:' wrapper in output")
	}
	if !strings.Contains(yaml, "name: standalone-group") {
		t.Error("Expected 'name: standalone-group' in output")
	}
}
