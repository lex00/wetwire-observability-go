package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewAlertmanagerConfig(t *testing.T) {
	config := NewAlertmanagerConfig()
	if config == nil {
		t.Error("NewAlertmanagerConfig() returned nil")
	}
}

func TestAlertmanagerConfig_FluentAPI(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithGlobal(NewGlobalConfig().
			WithResolveTimeout(5 * Minute)).
		WithRoute(NewRoute("default")).
		WithReceivers(NewReceiver("default")).
		WithTemplates("/etc/alertmanager/templates/*.tmpl")

	if config.Global == nil {
		t.Error("Global is nil")
	}
	if config.Route == nil {
		t.Error("Route is nil")
	}
	if len(config.Receivers) != 1 {
		t.Errorf("len(Receivers) = %d, want 1", len(config.Receivers))
	}
	if len(config.Templates) != 1 {
		t.Errorf("len(Templates) = %d, want 1", len(config.Templates))
	}
}

func TestNewGlobalConfig(t *testing.T) {
	g := NewGlobalConfig().
		WithSMTP("smtp.example.com:587", "alerts@example.com").
		WithSMTPAuth("user", "password").
		WithSlackAPIURL("https://hooks.slack.com/services/xxx").
		WithResolveTimeout(5 * Minute)

	if g.SMTPSmarthost != "smtp.example.com:587" {
		t.Errorf("SMTPSmarthost = %v, want smtp.example.com:587", g.SMTPSmarthost)
	}
	if g.SMTPFrom != "alerts@example.com" {
		t.Errorf("SMTPFrom = %v, want alerts@example.com", g.SMTPFrom)
	}
	if g.SMTPAuthUsername != "user" {
		t.Errorf("SMTPAuthUsername = %v, want user", g.SMTPAuthUsername)
	}
	if g.SlackAPIURL != "https://hooks.slack.com/services/xxx" {
		t.Errorf("SlackAPIURL = %v", g.SlackAPIURL)
	}
	if g.ResolveTimeout != 5*Minute {
		t.Errorf("ResolveTimeout = %v, want 5m", g.ResolveTimeout)
	}
}

func TestNewReceiver(t *testing.T) {
	r := NewReceiver("pagerduty-critical")
	if r.Name != "pagerduty-critical" {
		t.Errorf("Name = %v, want pagerduty-critical", r.Name)
	}
}

func TestAlertmanagerConfig_Serialize(t *testing.T) {
	config := &AlertmanagerConfig{
		Global: &GlobalConfig{
			ResolveTimeout: 5 * Minute,
		},
		Route: NewRoute("default").
			WithGroupBy("alertname", "cluster").
			WithGroupWait(30 * Second).
			WithRoutes(
				NewRoute("critical").
					WithMatchers(SeverityCritical()),
			),
		Receivers: []*Receiver{
			{Name: "default"},
			{Name: "critical"},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"global:",
		"resolve_timeout: 5m",
		"route:",
		"receiver: default",
		"group_by:",
		"alertname",
		"routes:",
		"receivers:",
		"name: default",
		"name: critical",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestAlertmanagerConfig_Unmarshal(t *testing.T) {
	input := `
global:
  resolve_timeout: 5m
route:
  receiver: default
  group_by:
    - alertname
  routes:
    - receiver: critical
      matchers:
        - severity="critical"
receivers:
  - name: default
  - name: critical
`
	var config AlertmanagerConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if config.Global == nil {
		t.Error("Global is nil")
	}
	if config.Global.ResolveTimeout != 5*Minute {
		t.Errorf("ResolveTimeout = %v, want 5m", config.Global.ResolveTimeout)
	}
	if config.Route == nil {
		t.Error("Route is nil")
	}
	if config.Route.Receiver != "default" {
		t.Errorf("Route.Receiver = %v, want default", config.Route.Receiver)
	}
	if len(config.Receivers) != 2 {
		t.Errorf("len(Receivers) = %d, want 2", len(config.Receivers))
	}
}

func TestCompleteAlertmanagerConfig(t *testing.T) {
	// This test verifies a complete alertmanager.yml configuration
	config := &AlertmanagerConfig{
		Global: NewGlobalConfig().
			WithResolveTimeout(5 * Minute).
			WithSlackAPIURL("https://hooks.slack.com/services/xxx"),

		Route: NewRoute("default-receiver").
			WithGroupBy("alertname", "cluster", "service").
			WithGroupWait(30 * Second).
			WithGroupInterval(5 * Minute).
			WithRepeatInterval(4 * Hour).
			WithRoutes(
				// Critical alerts to PagerDuty
				NewRoute("pagerduty-critical").
					WithMatchers(SeverityCritical()).
					WithGroupWait(10 * Second),

				// Backend team alerts
				NewRoute("slack-backend").
					WithMatchers(Team("backend")).
					WithContinue(true).
					WithRoutes(
						NewRoute("pagerduty-backend-critical").
							WithMatchers(SeverityCritical()),
					),

				// Frontend team alerts
				NewRoute("slack-frontend").
					WithMatchers(Team("frontend")),

				// Warning alerts to Slack
				NewRoute("slack-warnings").
					WithMatchers(SeverityWarning()),

				// Info alerts are muted during maintenance
				NewRoute("slack-info").
					WithMatchers(SeverityInfo()).
					WithMuteTimeIntervals("maintenance"),
			),

		Receivers: []*Receiver{
			{Name: "default-receiver"},
			{Name: "pagerduty-critical"},
			{Name: "pagerduty-backend-critical"},
			{Name: "slack-backend"},
			{Name: "slack-frontend"},
			{Name: "slack-warnings"},
			{Name: "slack-info"},
		},

		Templates: []string{
			"/etc/alertmanager/templates/*.tmpl",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Verify round-trip
	var restored AlertmanagerConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if restored.Route == nil {
		t.Fatal("restored.Route is nil")
	}
	if len(restored.Route.Routes) != 5 {
		t.Errorf("len(Routes) = %d, want 5", len(restored.Route.Routes))
	}
	if len(restored.Receivers) != 7 {
		t.Errorf("len(Receivers) = %d, want 7", len(restored.Receivers))
	}

	// Check nested route
	backendRoute := restored.Route.Routes[1]
	if len(backendRoute.Routes) != 1 {
		t.Errorf("len(backend.Routes) = %d, want 1", len(backendRoute.Routes))
	}

	t.Logf("Generated config:\n%s", string(data))
}
