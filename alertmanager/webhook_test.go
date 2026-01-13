package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewWebhookConfig(t *testing.T) {
	config := NewWebhookConfig()
	if config == nil {
		t.Error("NewWebhookConfig() returned nil")
	}
}

func TestWebhookConfig_FluentAPI(t *testing.T) {
	config := NewWebhookConfig().
		WithURL("https://webhook.example.com/alerts").
		WithMaxAlerts(10).
		WithSendResolved(true)

	if config.URL != "https://webhook.example.com/alerts" {
		t.Errorf("URL = %v", config.URL)
	}
	if config.MaxAlerts == nil || *config.MaxAlerts != 10 {
		t.Errorf("MaxAlerts = %v, want 10", config.MaxAlerts)
	}
	if config.SendResolved == nil || !*config.SendResolved {
		t.Error("SendResolved should be true")
	}
}

func TestWebhookConfig_WithURLFile(t *testing.T) {
	config := NewWebhookConfig().
		WithURLFile("/etc/alertmanager/secrets/webhook-url")

	if config.URLFile != "/etc/alertmanager/secrets/webhook-url" {
		t.Errorf("URLFile = %v", config.URLFile)
	}
}

func TestWebhookConfig_WithHTTPConfig(t *testing.T) {
	httpConfig := &HTTPConfig{
		BearerToken: "my-token",
	}
	config := NewWebhookConfig().
		WithURL("https://webhook.example.com").
		WithHTTPConfig(httpConfig)

	if config.HTTPConfig == nil {
		t.Error("HTTPConfig should not be nil")
	}
	if config.HTTPConfig.BearerToken != "my-token" {
		t.Errorf("BearerToken = %v", config.HTTPConfig.BearerToken)
	}
}

func TestWebhookConfig_Serialize(t *testing.T) {
	config := NewWebhookConfig().
		WithURL("https://webhook.example.com/alerts").
		WithMaxAlerts(5).
		WithSendResolved(true)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"send_resolved: true",
		"url: https://webhook.example.com/alerts",
		"max_alerts: 5",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestWebhookConfig_SerializeWithHTTPConfig(t *testing.T) {
	config := NewWebhookConfig().
		WithURL("https://webhook.example.com").
		WithHTTPConfig(&HTTPConfig{
			BearerToken: "secret-token",
		})

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "http_config:") {
		t.Errorf("yaml.Marshal() missing http_config\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "bearer_token: secret-token") {
		t.Errorf("yaml.Marshal() missing bearer_token\nGot:\n%s", yamlStr)
	}
}

func TestWebhookConfig_Unmarshal(t *testing.T) {
	input := `
send_resolved: true
url: https://webhook.example.com/alerts
max_alerts: 10
http_config:
  bearer_token: test-token
`
	var config WebhookConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if config.URL != "https://webhook.example.com/alerts" {
		t.Errorf("URL = %v", config.URL)
	}
	if config.MaxAlerts == nil || *config.MaxAlerts != 10 {
		t.Errorf("MaxAlerts = %v, want 10", config.MaxAlerts)
	}
	if config.HTTPConfig == nil {
		t.Error("HTTPConfig should not be nil")
	}
	if config.HTTPConfig.BearerToken != "test-token" {
		t.Errorf("BearerToken = %v", config.HTTPConfig.BearerToken)
	}
}

func TestWebhookReceiver(t *testing.T) {
	receiver := WebhookReceiver("webhook-alerts", "https://webhook.example.com")
	if receiver.Name != "webhook-alerts" {
		t.Errorf("Name = %v", receiver.Name)
	}
	if len(receiver.WebhookConfigs) != 1 {
		t.Errorf("len(WebhookConfigs) = %d, want 1", len(receiver.WebhookConfigs))
	}
	if receiver.WebhookConfigs[0].URL != "https://webhook.example.com" {
		t.Errorf("URL = %v", receiver.WebhookConfigs[0].URL)
	}
}
