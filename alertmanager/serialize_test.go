package alertmanager

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSerialize_AlertmanagerConfig(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithGlobal(
			NewGlobalConfig().
				WithSlackAPIURL("https://hooks.slack.com/services/xxx").
				WithResolveTimeout(5 * Minute),
		).
		WithRoute(
			NewRoute("default").
				WithGroupBy("alertname", "cluster").
				WithGroupWait(30 * Second),
		).
		WithReceivers(
			NewReceiver("default"),
		)

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"global:",
		"slack_api_url:",
		"resolve_timeout: 5m",
		"route:",
		"receiver: default",
		"group_by:",
		"group_wait: 30s",
		"receivers:",
		"name: default",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Serialize() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestSerialize_AlertmanagerConfigComplete(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithGlobal(
			NewGlobalConfig().
				WithSMTP("smtp.example.com:587", "alertmanager@example.com").
				WithSlackAPIURL("https://hooks.slack.com/services/xxx"),
		).
		WithRoute(
			NewRoute("default").
				WithGroupBy("alertname", "cluster").
				WithRoutes(
					NewRoute("pagerduty-critical").
						WithMatchers(Eq("severity", "critical")),
				),
		).
		WithReceivers(
			NewReceiver("default").WithSlackConfigs(
				NewSlackConfig().WithChannel("#alerts"),
			),
			NewReceiver("pagerduty-critical").WithPagerDutyConfigs(
				NewPagerDutyConfig().WithRoutingKey(NewSecret("key")),
			),
		).
		WithInhibitRules(
			CriticalInhibitsWarning(),
		).
		WithMuteTimeIntervals(
			WeekendsMuteInterval(),
		).
		WithTemplates("/etc/alertmanager/templates/*.tmpl")

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"global:",
		"smtp_smarthost:",
		"smtp_from:",
		"route:",
		"routes:",
		"receivers:",
		"slack_configs:",
		"channel:",
		"pagerduty_configs:",
		"routing_key:",
		"inhibit_rules:",
		"source_matchers:",
		"target_matchers:",
		"equal:",
		"mute_time_intervals:",
		"name: weekends",
		"templates:",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Serialize() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestSerialize_AlertmanagerConfigToFile(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithGlobal(
			NewGlobalConfig().WithResolveTimeout(5 * Minute),
		).
		WithRoute(NewRoute("default")).
		WithReceivers(NewReceiver("default"))

	tmpDir := t.TempDir()
	outputPath := filepath.Join(tmpDir, "alertmanager.yml")

	err := config.SerializeToFile(outputPath)
	if err != nil {
		t.Fatalf("SerializeToFile() error = %v", err)
	}

	// Verify file exists and has content
	data, err := os.ReadFile(outputPath)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if len(data) == 0 {
		t.Error("SerializeToFile() wrote empty file")
	}

	if !strings.Contains(string(data), "resolve_timeout: 5m") {
		t.Errorf("SerializeToFile() missing expected content\nGot:\n%s", data)
	}
}

func TestAlertmanagerConfig_MustSerialize(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithRoute(NewRoute("default")).
		WithReceivers(NewReceiver("default"))

	// Should not panic
	data := config.MustSerialize()
	if len(data) == 0 {
		t.Error("MustSerialize() returned empty data")
	}
}

func TestGlobalConfig_Serialize(t *testing.T) {
	global := NewGlobalConfig().
		WithSMTP("smtp.example.com:587", "alertmanager@example.com").
		WithSMTPAuth("user", "pass")

	data, err := global.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "smtp_smarthost:") {
		t.Errorf("Serialize() missing smtp_smarthost\nGot:\n%s", yamlStr)
	}
}

func TestSerialize_Route(t *testing.T) {
	route := NewRoute("default").
		WithGroupBy("alertname").
		WithGroupWait(30 * Second).
		WithRoutes(
			NewRoute("critical").WithMatchers(Eq("severity", "critical")),
		)

	data, err := route.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"receiver: default",
		"group_by:",
		"group_wait: 30s",
		"routes:",
		"receiver: critical",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Serialize() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestReceiver_Serialize(t *testing.T) {
	receiver := NewReceiver("multi-channel").
		WithSlackConfigs(
			NewSlackConfig().WithChannel("#alerts"),
		).
		WithPagerDutyConfigs(
			NewPagerDutyConfig().WithRoutingKey(NewSecret("key")),
		).
		WithEmailConfigs(
			NewEmailConfig().WithTo("team@example.com"),
		)

	data, err := receiver.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"name: multi-channel",
		"slack_configs:",
		"pagerduty_configs:",
		"email_configs:",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Serialize() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestSerialize_InhibitRule(t *testing.T) {
	rule := CriticalInhibitsWarning()

	data, err := rule.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "source_matchers:") {
		t.Errorf("Serialize() missing source_matchers\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "target_matchers:") {
		t.Errorf("Serialize() missing target_matchers\nGot:\n%s", yamlStr)
	}
}

func TestSerialize_MuteTimeInterval(t *testing.T) {
	mti := WeekendsMuteInterval()

	data, err := mti.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "name: weekends") {
		t.Errorf("Serialize() missing name\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "time_intervals:") {
		t.Errorf("Serialize() missing time_intervals\nGot:\n%s", yamlStr)
	}
}
