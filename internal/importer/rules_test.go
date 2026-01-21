package importer

import (
	"strings"
	"testing"
)

func TestParseRulesFileFromBytes(t *testing.T) {
	input := `
groups:
  - name: example
    interval: 1m
    rules:
      - alert: HighErrorRate
        expr: sum(rate(http_errors_total[5m])) / sum(rate(http_requests_total[5m])) > 0.1
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High error rate detected"
          description: "Error rate is above 10%"
      - record: job:http_requests:rate5m
        expr: sum(rate(http_requests_total[5m])) by (job)
        labels:
          job: api
`

	rf, err := ParseRulesFileFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse rules file: %v", err)
	}

	if len(rf.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(rf.Groups))
	}

	group := rf.Groups[0]
	if group.Name != "example" {
		t.Errorf("expected group name 'example', got %q", group.Name)
	}

	if group.Interval != "1m" {
		t.Errorf("expected interval '1m', got %q", group.Interval)
	}

	if len(group.Rules) != 2 {
		t.Errorf("expected 2 rules, got %d", len(group.Rules))
	}

	// Check alerting rule
	alert := group.Rules[0]
	if !alert.IsAlertingRule() {
		t.Error("expected first rule to be alerting rule")
	}
	if alert.Alert != "HighErrorRate" {
		t.Errorf("expected alert name 'HighErrorRate', got %q", alert.Alert)
	}
	if alert.For != "5m" {
		t.Errorf("expected for '5m', got %q", alert.For)
	}
	if alert.Labels["severity"] != "critical" {
		t.Errorf("expected severity 'critical', got %q", alert.Labels["severity"])
	}

	// Check recording rule
	record := group.Rules[1]
	if !record.IsRecordingRule() {
		t.Error("expected second rule to be recording rule")
	}
	if record.Record != "job:http_requests:rate5m" {
		t.Errorf("expected record name 'job:http_requests:rate5m', got %q", record.Record)
	}
}

func TestValidateRulesFile(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantWarn []string
	}{
		{
			name: "valid rules",
			input: `
groups:
  - name: test
    rules:
      - alert: TestAlert
        expr: up == 0
        for: 5m
        labels:
          severity: warning
`,
			wantWarn: nil,
		},
		{
			name: "no groups",
			input: `
groups: []
`,
			wantWarn: []string{"rules file has no groups"},
		},
		{
			name: "no rules in group",
			input: `
groups:
  - name: empty
    rules: []
`,
			wantWarn: []string{"has no rules"},
		},
		{
			name: "missing for duration",
			input: `
groups:
  - name: test
    rules:
      - alert: NoForAlert
        expr: up == 0
        labels:
          severity: warning
`,
			wantWarn: []string{"has no 'for' duration"},
		},
		{
			name: "missing severity",
			input: `
groups:
  - name: test
    rules:
      - alert: NoSeverityAlert
        expr: up == 0
        for: 5m
`,
			wantWarn: []string{"has no severity label"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rf, err := ParseRulesFileFromBytes([]byte(tt.input))
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			warnings := ValidateRulesFile(rf)

			if len(tt.wantWarn) == 0 && len(warnings) > 0 {
				t.Errorf("expected no warnings, got %v", warnings)
			}

			for _, want := range tt.wantWarn {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected warning containing %q, got %v", want, warnings)
				}
			}
		})
	}
}

func TestGenerateRulesGoCode(t *testing.T) {
	input := `
groups:
  - name: api-alerts
    interval: 30s
    rules:
      - alert: HighLatency
        expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m])) > 1
        for: 5m
        labels:
          severity: critical
          team: api
        annotations:
          summary: "High latency on API"
          description: "P99 latency is above 1 second"
          runbook_url: "https://runbooks.example.com/high-latency"
      - record: job:http_latency:p99
        expr: histogram_quantile(0.99, rate(http_request_duration_seconds_bucket[5m]))
  - name: infra-rules
    rules:
      - alert: NodeDown
        expr: up == 0
        for: 1m
        labels:
          severity: warning
`

	rf, err := ParseRulesFileFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse rules file: %v", err)
	}

	code, err := GenerateRulesGoCode(rf, "monitoring")
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	codeStr := string(code)
	t.Logf("Generated code:\n%s", codeStr)

	// Check package declaration
	if !strings.Contains(codeStr, "package monitoring") {
		t.Error("expected 'package monitoring' in generated code")
	}

	// Check import
	if !strings.Contains(codeStr, "github.com/lex00/wetwire-observability-go/rules") {
		t.Error("expected rules import in generated code")
	}

	// Check alerting rule
	if !strings.Contains(codeStr, "rules.NewAlertingRule") {
		t.Error("expected NewAlertingRule in generated code")
	}

	if !strings.Contains(codeStr, `"HighLatency"`) {
		t.Error("expected HighLatency alert in generated code")
	}

	// Check recording rule
	if !strings.Contains(codeStr, "rules.NewRecordingRule") {
		t.Error("expected NewRecordingRule in generated code")
	}

	if !strings.Contains(codeStr, "job:http_latency:p99") {
		t.Error("expected recording rule name in generated code")
	}

	// Check rule group
	if !strings.Contains(codeStr, "rules.NewRuleGroup") {
		t.Error("expected NewRuleGroup in generated code")
	}

	// Check severity convenience method
	if !strings.Contains(codeStr, "Critical()") {
		t.Error("expected Critical() method call in generated code")
	}

	// Check annotation convenience methods
	if !strings.Contains(codeStr, "WithSummary") {
		t.Error("expected WithSummary in generated code")
	}

	if !strings.Contains(codeStr, "WithDescription") {
		t.Error("expected WithDescription in generated code")
	}

	if !strings.Contains(codeStr, "WithRunbook") {
		t.Error("expected WithRunbook in generated code")
	}

	// Check rules file
	if !strings.Contains(codeStr, "rules.NewRulesFile") {
		t.Error("expected NewRulesFile in generated code")
	}

	// Check interval (go/format removes spaces around *)
	if !strings.Contains(codeStr, "30*rules.Second") && !strings.Contains(codeStr, "30 * rules.Second") {
		t.Error("expected formatted interval in generated code")
	}
}

func TestConvertToWetwireRules(t *testing.T) {
	input := `
groups:
  - name: test
    interval: 1m
    rules:
      - alert: TestAlert
        expr: up == 0
        for: 5m
        labels:
          severity: critical
`

	rf, err := ParseRulesFileFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	result := ConvertToWetwireRules(rf)

	if len(result.Groups) != 1 {
		t.Errorf("expected 1 group, got %d", len(result.Groups))
	}

	if result.Groups[0].Name != "test" {
		t.Errorf("expected group name 'test', got %q", result.Groups[0].Name)
	}

	if len(result.Groups[0].Rules) != 1 {
		t.Errorf("expected 1 rule, got %d", len(result.Groups[0].Rules))
	}
}
