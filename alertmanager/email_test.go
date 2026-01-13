package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewEmailConfig(t *testing.T) {
	config := NewEmailConfig()
	if config == nil {
		t.Error("NewEmailConfig() returned nil")
	}
}

func TestEmailConfig_FluentAPI(t *testing.T) {
	config := NewEmailConfig().
		WithTo("team@example.com").
		WithFrom("alertmanager@example.com").
		WithSmarthost("smtp.example.com:587").
		WithAuthUsername("user").
		WithAuthPassword(NewSecret("password")).
		WithRequireTLS(true).
		WithSendResolved(true)

	if config.To != "team@example.com" {
		t.Errorf("To = %v, want team@example.com", config.To)
	}
	if config.From != "alertmanager@example.com" {
		t.Errorf("From = %v", config.From)
	}
	if config.Smarthost != "smtp.example.com:587" {
		t.Errorf("Smarthost = %v", config.Smarthost)
	}
	if config.AuthUsername != "user" {
		t.Errorf("AuthUsername = %v", config.AuthUsername)
	}
	if config.RequireTLS == nil || !*config.RequireTLS {
		t.Error("RequireTLS should be true")
	}
	if config.SendResolved == nil || !*config.SendResolved {
		t.Error("SendResolved should be true")
	}
}

func TestEmailConfig_WithHeaders(t *testing.T) {
	config := NewEmailConfig().
		WithTo("team@example.com").
		WithHeaders(map[string]string{
			"Subject":  "[{{ .Status }}] {{ .CommonLabels.alertname }}",
			"Reply-To": "no-reply@example.com",
		})

	if len(config.Headers) != 2 {
		t.Errorf("len(Headers) = %d, want 2", len(config.Headers))
	}
}

func TestEmailConfig_WithHTML(t *testing.T) {
	config := NewEmailConfig().
		WithTo("team@example.com").
		WithHTML("<h1>{{ .Status }}</h1>").
		WithText("Status: {{ .Status }}")

	if config.HTML != "<h1>{{ .Status }}</h1>" {
		t.Errorf("HTML = %v", config.HTML)
	}
	if config.Text != "Status: {{ .Status }}" {
		t.Errorf("Text = %v", config.Text)
	}
}

func TestEmailConfig_Serialize(t *testing.T) {
	config := NewEmailConfig().
		WithTo("team@example.com").
		WithFrom("alertmanager@example.com").
		WithSmarthost("smtp.example.com:587").
		WithSendResolved(true)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"send_resolved: true",
		"to: team@example.com",
		"from: alertmanager@example.com",
		"smarthost: smtp.example.com:587",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestEmailConfig_SerializeWithAuth(t *testing.T) {
	config := NewEmailConfig().
		WithTo("team@example.com").
		WithAuthUsername("user").
		WithAuthPassword(NewSecret("secret-pass")).
		WithAuthIdentity("identity")

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"auth_username: user",
		"auth_password: secret-pass",
		"auth_identity: identity",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestEmailConfig_Unmarshal(t *testing.T) {
	input := `
send_resolved: true
to: team@example.com
from: alertmanager@example.com
smarthost: smtp.example.com:587
auth_username: user
auth_password: secret
require_tls: true
headers:
  Subject: "[ALERT] {{ .CommonLabels.alertname }}"
`
	var config EmailConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if config.To != "team@example.com" {
		t.Errorf("To = %v", config.To)
	}
	if config.Smarthost != "smtp.example.com:587" {
		t.Errorf("Smarthost = %v", config.Smarthost)
	}
	if config.AuthUsername != "user" {
		t.Errorf("AuthUsername = %v", config.AuthUsername)
	}
	if config.RequireTLS == nil || !*config.RequireTLS {
		t.Error("RequireTLS should be true")
	}
	if len(config.Headers) != 1 {
		t.Errorf("len(Headers) = %d, want 1", len(config.Headers))
	}
}

func TestEmailReceiver(t *testing.T) {
	receiver := EmailReceiver("email-alerts", "team@example.com")
	if receiver.Name != "email-alerts" {
		t.Errorf("Name = %v", receiver.Name)
	}
	if len(receiver.EmailConfigs) != 1 {
		t.Errorf("len(EmailConfigs) = %d, want 1", len(receiver.EmailConfigs))
	}
	if receiver.EmailConfigs[0].To != "team@example.com" {
		t.Errorf("To = %v", receiver.EmailConfigs[0].To)
	}
}
