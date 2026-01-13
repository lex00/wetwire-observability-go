package operator

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-observability-go/rules"
)

func TestPromRule(t *testing.T) {
	pr := PromRule("api-alerts", "monitoring")
	if pr.Name != "api-alerts" {
		t.Errorf("Name = %q, want api-alerts", pr.Name)
	}
	if pr.Namespace != "monitoring" {
		t.Errorf("Namespace = %q, want monitoring", pr.Namespace)
	}
	if pr.Kind != "PrometheusRule" {
		t.Errorf("Kind = %q, want PrometheusRule", pr.Kind)
	}
	if pr.APIVersion != "monitoring.coreos.com/v1" {
		t.Errorf("APIVersion = %q", pr.APIVersion)
	}
}

func TestPrometheusRule_WithLabels(t *testing.T) {
	pr := PromRule("alerts", "default").
		WithLabels(map[string]string{"prometheus": "main"})
	if pr.Labels["prometheus"] != "main" {
		t.Errorf("Labels[prometheus] = %q, want main", pr.Labels["prometheus"])
	}
}

func TestPrometheusRule_AddLabel(t *testing.T) {
	pr := PromRule("alerts", "default").
		AddLabel("team", "platform").
		AddLabel("severity", "critical")
	if len(pr.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(pr.Labels))
	}
}

func TestPrometheusRule_WithAnnotations(t *testing.T) {
	pr := PromRule("alerts", "default").
		WithAnnotations(map[string]string{"description": "Main alerts"})
	if pr.Metadata.Annotations["description"] != "Main alerts" {
		t.Error("Annotations should contain description")
	}
}

func TestPrometheusRule_AddRuleGroup(t *testing.T) {
	group := rules.NewRuleGroup("api.rules").
		WithInterval(rules.Minute).
		AddRule(rules.NewAlertingRule("HighErrorRate").
			WithExpr("rate(http_errors_total[5m]) > 0.1").
			WithFor(5 * rules.Minute))

	pr := PromRule("alerts", "monitoring").
		AddRuleGroup(group)

	if len(pr.Spec.Groups) != 1 {
		t.Errorf("len(Groups) = %d, want 1", len(pr.Spec.Groups))
	}
}

func TestPrometheusRule_WithRuleGroups(t *testing.T) {
	group1 := rules.NewRuleGroup("api.rules")
	group2 := rules.NewRuleGroup("db.rules")

	pr := PromRule("alerts", "monitoring").
		WithRuleGroups(group1, group2)

	if len(pr.Spec.Groups) != 2 {
		t.Errorf("len(Groups) = %d, want 2", len(pr.Spec.Groups))
	}
}

func TestPrometheusRule_Serialize(t *testing.T) {
	group := rules.NewRuleGroup("api.rules").
		AddRule(rules.NewAlertingRule("HighErrorRate").
			WithExpr("rate(http_errors_total[5m]) > 0.1"))

	pr := PromRule("api-alerts", "monitoring").
		WithLabels(map[string]string{"prometheus": "main"}).
		AddRuleGroup(group)

	data, err := pr.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "apiVersion: monitoring.coreos.com/v1") {
		t.Error("Expected apiVersion")
	}
	if !strings.Contains(yamlStr, "kind: PrometheusRule") {
		t.Error("Expected kind: PrometheusRule")
	}
	if !strings.Contains(yamlStr, "name: api-alerts") {
		t.Error("Expected name: api-alerts")
	}
	if !strings.Contains(yamlStr, "prometheus: main") {
		t.Error("Expected label prometheus: main")
	}
	if !strings.Contains(yamlStr, "api.rules") {
		t.Error("Expected group name api.rules")
	}
	if !strings.Contains(yamlStr, "HighErrorRate") {
		t.Error("Expected alert name HighErrorRate")
	}
}

func TestPrometheusRule_FromRulesFile(t *testing.T) {
	rf := &rules.RulesFile{
		Groups: []*rules.RuleGroup{
			rules.NewRuleGroup("api.rules"),
			rules.NewRuleGroup("db.rules"),
		},
	}

	pr := PromRuleFromRulesFile("alerts", "monitoring", rf)

	if len(pr.Spec.Groups) != 2 {
		t.Errorf("len(Groups) = %d, want 2", len(pr.Spec.Groups))
	}
}

func TestPrometheusRule_FluentAPI(t *testing.T) {
	group := rules.NewRuleGroup("api.rules").
		WithInterval(30 * rules.Second)

	pr := PromRule("alerts", "monitoring").
		WithLabels(map[string]string{"prometheus": "main"}).
		WithAnnotations(map[string]string{"description": "Alert rules"}).
		AddRuleGroup(group)

	if pr.Name != "alerts" {
		t.Error("Fluent API should preserve name")
	}
	if len(pr.Spec.Groups) != 1 {
		t.Error("Fluent API should add group")
	}
}

func TestPrometheusRule_WithRecordingRules(t *testing.T) {
	group := rules.NewRuleGroup("recording.rules").
		AddRule(rules.NewRecordingRule("job:http_requests:rate5m").
			WithExpr("sum(rate(http_requests_total[5m])) by (job)"))

	pr := PromRule("recording", "monitoring").
		AddRuleGroup(group)

	data, err := pr.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	if !strings.Contains(string(data), "job:http_requests:rate5m") {
		t.Error("Expected recording rule name")
	}
}
