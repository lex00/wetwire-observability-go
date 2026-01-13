package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewSlackConfig(t *testing.T) {
	config := NewSlackConfig()
	if config == nil {
		t.Error("NewSlackConfig() returned nil")
	}
}

func TestSlackConfig_FluentAPI(t *testing.T) {
	config := NewSlackConfig().
		WithAPIURL(NewSecret("https://hooks.slack.com/services/xxx")).
		WithChannel("#alerts").
		WithUsername("alertbot").
		WithIconEmoji(":warning:").
		WithTitle("Alert: {{ .Status }}").
		WithText("{{ .Annotations.description }}").
		WithColor(SlackColorDanger).
		WithSendResolved(true)

	if config.APIURL != "https://hooks.slack.com/services/xxx" {
		t.Errorf("APIURL = %v", config.APIURL)
	}
	if config.Channel != "#alerts" {
		t.Errorf("Channel = %v, want #alerts", config.Channel)
	}
	if config.Username != "alertbot" {
		t.Errorf("Username = %v, want alertbot", config.Username)
	}
	if config.IconEmoji != ":warning:" {
		t.Errorf("IconEmoji = %v", config.IconEmoji)
	}
	if config.Color != SlackColorDanger {
		t.Errorf("Color = %v, want danger", config.Color)
	}
	if config.SendResolved == nil || !*config.SendResolved {
		t.Error("SendResolved should be true")
	}
}

func TestSlackConfig_WithActions(t *testing.T) {
	config := NewSlackConfig().
		WithActions(
			NewSlackAction("View Dashboard", "https://grafana.example.com"),
			NewSlackAction("Silence", "https://alertmanager.example.com/silences").
				WithStyle("danger"),
		)

	if len(config.Actions) != 2 {
		t.Errorf("len(Actions) = %d, want 2", len(config.Actions))
	}
	if config.Actions[0].Text != "View Dashboard" {
		t.Errorf("Actions[0].Text = %v", config.Actions[0].Text)
	}
	if config.Actions[1].Style != "danger" {
		t.Errorf("Actions[1].Style = %v, want danger", config.Actions[1].Style)
	}
}

func TestSlackConfig_WithFields(t *testing.T) {
	config := NewSlackConfig().
		WithFields(
			NewSlackField("Severity", "{{ .Labels.severity }}").WithShort(true),
			NewSlackField("Service", "{{ .Labels.service }}").WithShort(true),
			NewSlackField("Description", "{{ .Annotations.description }}"),
		)

	if len(config.Fields) != 3 {
		t.Errorf("len(Fields) = %d, want 3", len(config.Fields))
	}
	if config.Fields[0].Short == nil || !*config.Fields[0].Short {
		t.Error("Fields[0].Short should be true")
	}
}

func TestSlackConfig_Serialize(t *testing.T) {
	config := NewSlackConfig().
		WithChannel("#alerts").
		WithUsername("alertbot").
		WithIconEmoji(":bell:").
		WithColor(SlackColorWarning).
		WithSendResolved(true)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"send_resolved: true",
		"channel: '#alerts'",
		"username: alertbot",
		"icon_emoji: ':bell:'",
		"color: warning",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestSlackConfig_SerializeWithActions(t *testing.T) {
	config := NewSlackConfig().
		WithChannel("#alerts").
		WithActions(
			NewSlackAction("View", "https://example.com").WithStyle("primary"),
		)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"actions:",
		"type: button",
		"text: View",
		"url: https://example.com",
		"style: primary",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestSlackConfig_SerializeWithFields(t *testing.T) {
	config := NewSlackConfig().
		WithChannel("#alerts").
		WithFields(
			NewSlackField("Severity", "critical").WithShort(true),
			NewSlackField("Service", "api"),
		)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"fields:",
		"title: Severity",
		"value: critical",
		"short: true",
		"title: Service",
		"value: api",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestSlackConfig_Unmarshal(t *testing.T) {
	input := `
send_resolved: true
channel: '#alerts'
username: alertbot
icon_emoji: ':warning:'
color: danger
actions:
  - type: button
    text: View
    url: https://example.com
`
	var config SlackConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if config.Channel != "#alerts" {
		t.Errorf("Channel = %v, want #alerts", config.Channel)
	}
	if config.SendResolved == nil || !*config.SendResolved {
		t.Error("SendResolved should be true")
	}
	if len(config.Actions) != 1 {
		t.Errorf("len(Actions) = %d, want 1", len(config.Actions))
	}
}

func TestSlackColorConstants(t *testing.T) {
	if SlackColorGood != "good" {
		t.Errorf("SlackColorGood = %v, want good", SlackColorGood)
	}
	if SlackColorWarning != "warning" {
		t.Errorf("SlackColorWarning = %v, want warning", SlackColorWarning)
	}
	if SlackColorDanger != "danger" {
		t.Errorf("SlackColorDanger = %v, want danger", SlackColorDanger)
	}
}

func TestNewSlackAction(t *testing.T) {
	action := NewSlackAction("Click Me", "https://example.com")
	if action.Type != "button" {
		t.Errorf("Type = %v, want button", action.Type)
	}
	if action.Text != "Click Me" {
		t.Errorf("Text = %v, want Click Me", action.Text)
	}
	if action.URL != "https://example.com" {
		t.Errorf("URL = %v", action.URL)
	}
}

func TestNewSlackField(t *testing.T) {
	field := NewSlackField("Title", "Value")
	if field.Title != "Title" {
		t.Errorf("Title = %v, want Title", field.Title)
	}
	if field.Value != "Value" {
		t.Errorf("Value = %v, want Value", field.Value)
	}
}
