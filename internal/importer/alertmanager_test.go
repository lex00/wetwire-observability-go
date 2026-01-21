package importer

import (
	"strings"
	"testing"
)

func TestParseAlertmanagerConfigFromBytes_Basic(t *testing.T) {
	yaml := `
global:
  smtp_smarthost: 'localhost:25'
  smtp_from: 'alertmanager@example.org'
  resolve_timeout: 5m

route:
  receiver: 'team-notifications'
  group_by: ['alertname', 'severity']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 4h

receivers:
  - name: 'team-notifications'
    slack_configs:
      - channel: '#alerts'
        api_url: 'https://hooks.slack.com/services/xxx'
`
	config, err := ParseAlertmanagerConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("ParseAlertmanagerConfigFromBytes failed: %v", err)
	}

	// Check global config
	if config.Global == nil {
		t.Fatal("expected Global config")
	}
	if config.Global.SMTPSmarthost != "localhost:25" {
		t.Errorf("expected SMTPSmarthost 'localhost:25', got %q", config.Global.SMTPSmarthost)
	}

	// Check route
	if config.Route == nil {
		t.Fatal("expected Route config")
	}
	if config.Route.Receiver != "team-notifications" {
		t.Errorf("expected receiver 'team-notifications', got %q", config.Route.Receiver)
	}

	// Check receivers
	if len(config.Receivers) != 1 {
		t.Fatalf("expected 1 receiver, got %d", len(config.Receivers))
	}
	if config.Receivers[0].Name != "team-notifications" {
		t.Errorf("expected receiver name 'team-notifications', got %q", config.Receivers[0].Name)
	}
}

func TestValidateAlertmanagerConfig(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		wantWarn int
	}{
		{
			name: "valid config",
			yaml: `
route:
  receiver: 'default'
receivers:
  - name: 'default'
`,
			wantWarn: 0,
		},
		{
			name:     "missing route",
			yaml:     `receivers: []`,
			wantWarn: 2, // missing route + no receivers
		},
		{
			name: "missing receiver",
			yaml: `
route:
  receiver: 'nonexistent'
receivers:
  - name: 'default'
`,
			wantWarn: 1, // receiver not found
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config, err := ParseAlertmanagerConfigFromBytes([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("parse failed: %v", err)
			}

			warnings := ValidateAlertmanagerConfig(config)
			if len(warnings) != tt.wantWarn {
				t.Errorf("expected %d warnings, got %d: %v", tt.wantWarn, len(warnings), warnings)
			}
		})
	}
}

func TestGenerateAlertmanagerGoCode(t *testing.T) {
	yaml := `
route:
  receiver: 'slack'
  group_by: ['alertname']
receivers:
  - name: 'slack'
    slack_configs:
      - channel: '#alerts'
`
	config, err := ParseAlertmanagerConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("parse failed: %v", err)
	}

	code, err := GenerateAlertmanagerGoCode(config, "monitoring")
	if err != nil {
		t.Fatalf("generate failed: %v", err)
	}

	codeStr := string(code)

	// Check package declaration
	if !strings.Contains(codeStr, "package monitoring") {
		t.Error("expected package declaration")
	}

	// Check import
	if !strings.Contains(codeStr, `"github.com/lex00/wetwire-observability-go/alertmanager"`) {
		t.Error("expected alertmanager import")
	}

	// Check receiver generation
	if !strings.Contains(codeStr, "SlackReceiver") {
		t.Error("expected SlackReceiver variable")
	}

	// Check route generation
	if !strings.Contains(codeStr, "RootRoute") {
		t.Error("expected RootRoute variable")
	}

	// Check main config
	if !strings.Contains(codeStr, "var Config = &alertmanager.AlertmanagerConfig") {
		t.Error("expected Config variable")
	}
}
