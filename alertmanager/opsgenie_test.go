package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewOpsGenieConfig(t *testing.T) {
	config := NewOpsGenieConfig()
	if config == nil {
		t.Error("NewOpsGenieConfig() returned nil")
	}
}

func TestOpsGenieConfig_FluentAPI(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKey(NewSecret("api-key-123")).
		WithMessage("{{ .CommonLabels.alertname }}").
		WithDescription("{{ .CommonAnnotations.description }}").
		WithPriority("P1").
		WithSource("Alertmanager").
		WithSendResolved(true)

	if string(config.APIKey) != "api-key-123" {
		t.Errorf("APIKey = %v", config.APIKey)
	}
	if config.Message != "{{ .CommonLabels.alertname }}" {
		t.Errorf("Message = %v", config.Message)
	}
	if config.Priority != "P1" {
		t.Errorf("Priority = %v", config.Priority)
	}
	if config.SendResolved == nil || !*config.SendResolved {
		t.Error("SendResolved should be true")
	}
}

func TestOpsGenieConfig_WithAPIKeyFile(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKeyFile("/etc/alertmanager/secrets/opsgenie-key")

	if config.APIKeyFile != "/etc/alertmanager/secrets/opsgenie-key" {
		t.Errorf("APIKeyFile = %v", config.APIKeyFile)
	}
}

func TestOpsGenieConfig_WithTags(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKey(NewSecret("key")).
		WithTags("production", "critical", "{{ .CommonLabels.team }}")

	if len(config.Tags) != 3 {
		t.Errorf("len(Tags) = %d, want 3", len(config.Tags))
	}
	if config.Tags[0] != "production" {
		t.Errorf("Tags[0] = %v", config.Tags[0])
	}
}

func TestOpsGenieConfig_WithDetails(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKey(NewSecret("key")).
		WithDetails(map[string]string{
			"firing":    "{{ .Alerts.Firing | len }}",
			"alertname": "{{ .CommonLabels.alertname }}",
		})

	if len(config.Details) != 2 {
		t.Errorf("len(Details) = %d, want 2", len(config.Details))
	}
}

func TestOpsGenieConfig_WithResponders(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKey(NewSecret("key")).
		WithResponders(
			NewOpsGenieResponder("team", "ops-team"),
			NewOpsGenieResponder("user", "on-call@example.com"),
		)

	if len(config.Responders) != 2 {
		t.Errorf("len(Responders) = %d, want 2", len(config.Responders))
	}
	if config.Responders[0].Type != "team" {
		t.Errorf("Responders[0].Type = %v", config.Responders[0].Type)
	}
	if config.Responders[0].Name != "ops-team" {
		t.Errorf("Responders[0].Name = %v", config.Responders[0].Name)
	}
}

func TestOpsGenieConfig_Serialize(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKey(NewSecret("test-api-key")).
		WithMessage("Alert: {{ .CommonLabels.alertname }}").
		WithPriority("P2").
		WithSendResolved(true)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"send_resolved: true",
		"api_key: test-api-key",
		"message: 'Alert: {{ .CommonLabels.alertname }}'",
		"priority: P2",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestOpsGenieConfig_SerializeWithResponders(t *testing.T) {
	config := NewOpsGenieConfig().
		WithAPIKey(NewSecret("key")).
		WithResponders(
			NewOpsGenieResponder("team", "platform").WithID("team-123"),
		)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"responders:",
		"type: team",
		"name: platform",
		"id: team-123",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestOpsGenieConfig_Unmarshal(t *testing.T) {
	input := `
send_resolved: true
api_key: test-key
message: "{{ .CommonLabels.alertname }}"
priority: P1
tags:
  - production
  - critical
responders:
  - type: team
    name: ops
`
	var config OpsGenieConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if string(config.APIKey) != "test-key" {
		t.Errorf("APIKey = %v", config.APIKey)
	}
	if config.Priority != "P1" {
		t.Errorf("Priority = %v", config.Priority)
	}
	if len(config.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(config.Tags))
	}
	if len(config.Responders) != 1 {
		t.Errorf("len(Responders) = %d, want 1", len(config.Responders))
	}
}

func TestOpsGeniePriorityConstants(t *testing.T) {
	tests := []struct {
		constant string
		want     string
	}{
		{OpsGeniePriorityP1, "P1"},
		{OpsGeniePriorityP2, "P2"},
		{OpsGeniePriorityP3, "P3"},
		{OpsGeniePriorityP4, "P4"},
		{OpsGeniePriorityP5, "P5"},
	}

	for _, tt := range tests {
		if tt.constant != tt.want {
			t.Errorf("constant = %v, want %v", tt.constant, tt.want)
		}
	}
}

func TestNewOpsGenieResponder(t *testing.T) {
	responder := NewOpsGenieResponder("team", "platform-team")
	if responder.Type != "team" {
		t.Errorf("Type = %v", responder.Type)
	}
	if responder.Name != "platform-team" {
		t.Errorf("Name = %v", responder.Name)
	}
}

func TestOpsGenieResponder_WithID(t *testing.T) {
	responder := NewOpsGenieResponder("team", "ops").WithID("team-abc123")
	if responder.ID != "team-abc123" {
		t.Errorf("ID = %v", responder.ID)
	}
}

func TestOpsGenieReceiver(t *testing.T) {
	receiver := OpsGenieReceiver("opsgenie-alerts", NewSecret("api-key"))
	if receiver.Name != "opsgenie-alerts" {
		t.Errorf("Name = %v", receiver.Name)
	}
	if len(receiver.OpsGenieConfigs) != 1 {
		t.Errorf("len(OpsGenieConfigs) = %d, want 1", len(receiver.OpsGenieConfigs))
	}
	if string(receiver.OpsGenieConfigs[0].APIKey) != "api-key" {
		t.Errorf("APIKey = %v", receiver.OpsGenieConfigs[0].APIKey)
	}
}
