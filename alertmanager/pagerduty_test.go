package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewPagerDutyConfig(t *testing.T) {
	config := NewPagerDutyConfig()
	if config == nil {
		t.Error("NewPagerDutyConfig() returned nil")
	}
}

func TestPagerDutyConfig_FluentAPI(t *testing.T) {
	config := NewPagerDutyConfig().
		WithRoutingKey(NewSecret("abc123")).
		WithClient("Prometheus").
		WithClientURL("https://prometheus.example.com").
		WithDescription("{{ .Annotations.summary }}").
		WithSeverity(PagerDutySeverityCritical).
		WithClass("monitoring").
		WithComponent("api-server").
		WithGroup("production").
		WithSendResolved(true)

	if config.RoutingKey != "abc123" {
		t.Errorf("RoutingKey = %v", config.RoutingKey)
	}
	if config.Client != "Prometheus" {
		t.Errorf("Client = %v, want Prometheus", config.Client)
	}
	if config.Severity != PagerDutySeverityCritical {
		t.Errorf("Severity = %v, want critical", config.Severity)
	}
	if config.Class != "monitoring" {
		t.Errorf("Class = %v, want monitoring", config.Class)
	}
	if config.SendResolved == nil || !*config.SendResolved {
		t.Error("SendResolved should be true")
	}
}

func TestPagerDutyConfig_WithDetails(t *testing.T) {
	config := NewPagerDutyConfig().
		WithDetails(map[string]string{
			"firing":       "{{ .Alerts.Firing | len }}",
			"resolved":     "{{ .Alerts.Resolved | len }}",
			"alertname":    "{{ .CommonLabels.alertname }}",
		})

	if len(config.Details) != 3 {
		t.Errorf("len(Details) = %d, want 3", len(config.Details))
	}
}

func TestPagerDutyConfig_WithImages(t *testing.T) {
	config := NewPagerDutyConfig().
		WithImages(
			NewPagerDutyImage("https://grafana.example.com/chart.png").
				WithHref("https://grafana.example.com").
				WithAlt("CPU Usage"),
		)

	if len(config.Images) != 1 {
		t.Errorf("len(Images) = %d, want 1", len(config.Images))
	}
	if config.Images[0].Src != "https://grafana.example.com/chart.png" {
		t.Errorf("Images[0].Src = %v", config.Images[0].Src)
	}
	if config.Images[0].Href != "https://grafana.example.com" {
		t.Errorf("Images[0].Href = %v", config.Images[0].Href)
	}
}

func TestPagerDutyConfig_WithLinks(t *testing.T) {
	config := NewPagerDutyConfig().
		WithLinks(
			NewPagerDutyLink("https://runbook.example.com").WithText("Runbook"),
			NewPagerDutyLink("https://grafana.example.com").WithText("Dashboard"),
		)

	if len(config.Links) != 2 {
		t.Errorf("len(Links) = %d, want 2", len(config.Links))
	}
	if config.Links[0].Text != "Runbook" {
		t.Errorf("Links[0].Text = %v, want Runbook", config.Links[0].Text)
	}
}

func TestPagerDutyConfig_Serialize(t *testing.T) {
	config := NewPagerDutyConfig().
		WithRoutingKey(NewSecret("test-key")).
		WithClient("Alertmanager").
		WithClientURL("https://alertmanager.example.com").
		WithSeverity(PagerDutySeverityWarning).
		WithSendResolved(true)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"send_resolved: true",
		"routing_key: test-key",
		"client: Alertmanager",
		"client_url: https://alertmanager.example.com",
		"severity: warning",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPagerDutyConfig_SerializeWithDetails(t *testing.T) {
	config := NewPagerDutyConfig().
		WithRoutingKey(NewSecret("key")).
		WithDetails(map[string]string{
			"env":     "production",
			"service": "api",
		})

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "details:") {
		t.Errorf("yaml.Marshal() missing details:\nGot:\n%s", yamlStr)
	}
}

func TestPagerDutyConfig_SerializeWithImages(t *testing.T) {
	config := NewPagerDutyConfig().
		WithRoutingKey(NewSecret("key")).
		WithImages(
			NewPagerDutyImage("https://example.com/chart.png").
				WithHref("https://example.com").
				WithAlt("Chart"),
		)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"images:",
		"src: https://example.com/chart.png",
		"href: https://example.com",
		"alt: Chart",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPagerDutyConfig_Unmarshal(t *testing.T) {
	input := `
send_resolved: true
routing_key: test-routing-key
client: Prometheus
client_url: https://prometheus.example.com
severity: critical
details:
  env: production
links:
  - href: https://runbook.example.com
    text: Runbook
`
	var config PagerDutyConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if string(config.RoutingKey) != "test-routing-key" {
		t.Errorf("RoutingKey = %v", config.RoutingKey)
	}
	if config.Severity != "critical" {
		t.Errorf("Severity = %v, want critical", config.Severity)
	}
	if len(config.Details) != 1 {
		t.Errorf("len(Details) = %d, want 1", len(config.Details))
	}
	if len(config.Links) != 1 {
		t.Errorf("len(Links) = %d, want 1", len(config.Links))
	}
}

func TestPagerDutySeverityConstants(t *testing.T) {
	tests := []struct {
		constant string
		want     string
	}{
		{PagerDutySeverityCritical, "critical"},
		{PagerDutySeverityError, "error"},
		{PagerDutySeverityWarning, "warning"},
		{PagerDutySeverityInfo, "info"},
	}

	for _, tt := range tests {
		if tt.constant != tt.want {
			t.Errorf("constant = %v, want %v", tt.constant, tt.want)
		}
	}
}

func TestPagerDutyConfig_RoutingKeyFile(t *testing.T) {
	config := NewPagerDutyConfig().
		WithRoutingKeyFile("/etc/alertmanager/pagerduty-key")

	if config.RoutingKeyFile != "/etc/alertmanager/pagerduty-key" {
		t.Errorf("RoutingKeyFile = %v", config.RoutingKeyFile)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	if !strings.Contains(string(data), "routing_key_file:") {
		t.Errorf("yaml.Marshal() missing routing_key_file\nGot:\n%s", data)
	}
}
