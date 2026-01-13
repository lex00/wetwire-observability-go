package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewRoute(t *testing.T) {
	r := NewRoute("default-receiver")
	if r.Receiver != "default-receiver" {
		t.Errorf("Receiver = %v, want default-receiver", r.Receiver)
	}
}

func TestRoute_FluentAPI(t *testing.T) {
	r := NewRoute("pagerduty-critical").
		WithGroupBy("alertname", "cluster").
		WithGroupWait(30 * Second).
		WithGroupInterval(5 * Minute).
		WithRepeatInterval(4 * Hour).
		WithMatchers(SeverityCritical()).
		WithContinue(false)

	if r.Receiver != "pagerduty-critical" {
		t.Errorf("Receiver = %v, want pagerduty-critical", r.Receiver)
	}
	if len(r.GroupBy) != 2 || r.GroupBy[0] != "alertname" {
		t.Errorf("GroupBy = %v, want [alertname cluster]", r.GroupBy)
	}
	if r.GroupWait != 30*Second {
		t.Errorf("GroupWait = %v, want 30s", r.GroupWait)
	}
	if r.GroupInterval != 5*Minute {
		t.Errorf("GroupInterval = %v, want 5m", r.GroupInterval)
	}
	if r.RepeatInterval != 4*Hour {
		t.Errorf("RepeatInterval = %v, want 4h", r.RepeatInterval)
	}
	if len(r.Matchers) != 1 {
		t.Errorf("len(Matchers) = %d, want 1", len(r.Matchers))
	}
}

func TestRoute_NestedRoutes(t *testing.T) {
	root := NewRoute("default").
		WithGroupBy("alertname").
		WithRoutes(
			NewRoute("pagerduty-critical").
				WithMatchers(SeverityCritical()),
			NewRoute("slack-warning").
				WithMatchers(SeverityWarning()),
			NewRoute("email-info").
				WithMatchers(SeverityInfo()),
		)

	if len(root.Routes) != 3 {
		t.Errorf("len(Routes) = %d, want 3", len(root.Routes))
	}
	if root.Routes[0].Receiver != "pagerduty-critical" {
		t.Errorf("Routes[0].Receiver = %v, want pagerduty-critical", root.Routes[0].Receiver)
	}
}

func TestRoute_AddRoute(t *testing.T) {
	root := NewRoute("default")
	root.AddRoute(NewRoute("critical").WithMatchers(SeverityCritical()))
	root.AddRoute(NewRoute("warning").WithMatchers(SeverityWarning()))

	if len(root.Routes) != 2 {
		t.Errorf("len(Routes) = %d, want 2", len(root.Routes))
	}
}

func TestRoute_ConvenienceMethods(t *testing.T) {
	r := NewRoute("team-backend").
		Severity("critical").
		Team("backend").
		Service("api").
		Environment("production")

	if len(r.Matchers) != 4 {
		t.Errorf("len(Matchers) = %d, want 4", len(r.Matchers))
	}

	expected := []struct {
		label string
		value string
	}{
		{"severity", "critical"},
		{"team", "backend"},
		{"service", "api"},
		{"env", "production"},
	}

	for i, exp := range expected {
		if r.Matchers[i].Label != exp.label || r.Matchers[i].Value != exp.value {
			t.Errorf("Matcher[%d] = %v=%v, want %v=%v", i, r.Matchers[i].Label, r.Matchers[i].Value, exp.label, exp.value)
		}
	}
}

func TestRoute_Serialize(t *testing.T) {
	r := NewRoute("default").
		WithGroupBy("alertname", "cluster").
		WithGroupWait(30 * Second).
		WithGroupInterval(5 * Minute).
		WithRepeatInterval(4 * Hour)

	data, err := yaml.Marshal(r)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"receiver: default",
		"group_by:",
		"alertname",
		"cluster",
		"group_wait: 30s",
		"group_interval: 5m",
		"repeat_interval: 4h",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRoute_SerializeWithMatchers(t *testing.T) {
	r := NewRoute("pagerduty").
		WithMatchers(
			Eq("severity", "critical"),
			NotEq("team", "testing"),
		)

	data, err := yaml.Marshal(r)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"receiver: pagerduty",
		"matchers:",
		`severity="critical"`,
		`team!="testing"`,
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRoute_SerializeNestedRoutes(t *testing.T) {
	root := NewRoute("default").
		WithGroupBy("alertname").
		WithRoutes(
			NewRoute("critical").
				WithMatchers(SeverityCritical()),
			NewRoute("warning").
				WithMatchers(SeverityWarning()),
		)

	data, err := yaml.Marshal(root)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"receiver: default",
		"routes:",
		"receiver: critical",
		"receiver: warning",
		`severity="critical"`,
		`severity="warning"`,
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRoute_Unmarshal(t *testing.T) {
	input := `
receiver: default
group_by:
  - alertname
  - cluster
group_wait: 30s
group_interval: 5m
repeat_interval: 4h
routes:
  - receiver: critical
    matchers:
      - severity="critical"
`
	var r Route
	if err := yaml.Unmarshal([]byte(input), &r); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if r.Receiver != "default" {
		t.Errorf("Receiver = %v, want default", r.Receiver)
	}
	if len(r.GroupBy) != 2 {
		t.Errorf("len(GroupBy) = %d, want 2", len(r.GroupBy))
	}
	if r.GroupWait != 30*Second {
		t.Errorf("GroupWait = %v, want 30s", r.GroupWait)
	}
	if len(r.Routes) != 1 {
		t.Errorf("len(Routes) = %d, want 1", len(r.Routes))
	}
	if r.Routes[0].Receiver != "critical" {
		t.Errorf("Routes[0].Receiver = %v, want critical", r.Routes[0].Receiver)
	}
}

func TestRoute_MuteTimeIntervals(t *testing.T) {
	r := NewRoute("default").
		WithMuteTimeIntervals("maintenance-window", "weekends")

	if len(r.MuteTimeIntervals) != 2 {
		t.Errorf("len(MuteTimeIntervals) = %d, want 2", len(r.MuteTimeIntervals))
	}

	data, err := yaml.Marshal(r)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "mute_time_intervals:") {
		t.Errorf("yaml.Marshal() missing mute_time_intervals\nGot:\n%s", yamlStr)
	}
}

func TestRoute_Continue(t *testing.T) {
	r := NewRoute("first").
		WithContinue(true).
		WithMatchers(SeverityCritical())

	if !r.Continue {
		t.Error("Continue should be true")
	}

	data, err := yaml.Marshal(r)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	if !strings.Contains(string(data), "continue: true") {
		t.Errorf("yaml.Marshal() missing continue: true\nGot:\n%s", data)
	}
}

func TestComplexRoutingTree(t *testing.T) {
	// This test verifies a real-world routing tree configuration
	root := NewRoute("default").
		WithGroupBy("alertname", "cluster").
		WithGroupWait(30 * Second).
		WithGroupInterval(5 * Minute).
		WithRepeatInterval(4 * Hour).
		WithRoutes(
			// Critical alerts go to PagerDuty
			NewRoute("pagerduty-critical").
				WithMatchers(SeverityCritical()).
				WithGroupWait(10 * Second),

			// Team-specific routing
			NewRoute("slack-backend").
				WithMatchers(Team("backend")).
				WithContinue(true).
				WithRoutes(
					NewRoute("pagerduty-backend").
						WithMatchers(SeverityCritical()),
				),

			// Warnings go to Slack
			NewRoute("slack-warning").
				WithMatchers(SeverityWarning()),

			// Everything else to email
			NewRoute("email-default"),
		)

	data, err := yaml.Marshal(root)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Verify round-trip
	var restored Route
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(restored.Routes) != 4 {
		t.Errorf("len(Routes) = %d, want 4", len(restored.Routes))
	}

	// Check nested route
	if len(restored.Routes[1].Routes) != 1 {
		t.Errorf("len(Routes[1].Routes) = %d, want 1", len(restored.Routes[1].Routes))
	}

	t.Logf("Generated routing tree:\n%s", string(data))
}
