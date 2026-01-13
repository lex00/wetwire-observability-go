package operator

import (
	"strings"
	"testing"
)

func TestAMConfig(t *testing.T) {
	am := AMConfig("alerts", "monitoring")
	if am.Name != "alerts" {
		t.Errorf("Name = %q, want alerts", am.Name)
	}
	if am.Namespace != "monitoring" {
		t.Errorf("Namespace = %q, want monitoring", am.Namespace)
	}
	if am.Kind != "AlertmanagerConfig" {
		t.Errorf("Kind = %q, want AlertmanagerConfig", am.Kind)
	}
	if am.APIVersion != "monitoring.coreos.com/v1alpha1" {
		t.Errorf("APIVersion = %q", am.APIVersion)
	}
}

func TestAlertmanagerConfig_WithLabels(t *testing.T) {
	am := AMConfig("alerts", "default").
		WithLabels(map[string]string{"alertmanager": "main"})
	if am.Labels["alertmanager"] != "main" {
		t.Errorf("Labels[alertmanager] = %q, want main", am.Labels["alertmanager"])
	}
}

func TestAlertmanagerConfig_AddLabel(t *testing.T) {
	am := AMConfig("alerts", "default").
		AddLabel("team", "platform").
		AddLabel("env", "prod")
	if len(am.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(am.Labels))
	}
}

func TestAlertmanagerConfig_WithRoute(t *testing.T) {
	route := NewAMRoute("slack-critical").
		WithGroupBy("alertname", "namespace").
		WithGroupWait("30s").
		WithGroupInterval("5m").
		WithRepeatInterval("4h")

	am := AMConfig("alerts", "monitoring").
		WithRoute(route)

	if am.Spec.Route == nil {
		t.Fatal("Route should not be nil")
	}
	if am.Spec.Route.Receiver != "slack-critical" {
		t.Errorf("Route.Receiver = %q", am.Spec.Route.Receiver)
	}
	if len(am.Spec.Route.GroupBy) != 2 {
		t.Errorf("len(GroupBy) = %d, want 2", len(am.Spec.Route.GroupBy))
	}
}

func TestAlertmanagerConfig_RouteMatchers(t *testing.T) {
	route := NewAMRoute("pagerduty-critical").
		AddMatcher(AMMatcherEq("severity", "critical")).
		AddMatcher(AMMatcherNeq("team", "test"))

	if len(route.Matchers) != 2 {
		t.Errorf("len(Matchers) = %d, want 2", len(route.Matchers))
	}
}

func TestAlertmanagerConfig_ChildRoutes(t *testing.T) {
	criticalRoute := NewAMRoute("pagerduty").
		AddMatcher(AMMatcherEq("severity", "critical"))

	warningRoute := NewAMRoute("slack-warnings").
		AddMatcher(AMMatcherEq("severity", "warning"))

	root := NewAMRoute("default").
		WithGroupBy("alertname").
		AddRoute(criticalRoute).
		AddRoute(warningRoute)

	if len(root.Routes) != 2 {
		t.Errorf("len(Routes) = %d, want 2", len(root.Routes))
	}
}

func TestAlertmanagerConfig_SlackReceiver(t *testing.T) {
	receiver := NewAMReceiver("slack-critical").
		WithSlackConfig(NewAMSlackConfig().
			WithChannel("#alerts").
			WithAPIURLSecret("alertmanager-slack", "webhook-url").
			WithSendResolved(true))

	if receiver.Name != "slack-critical" {
		t.Errorf("Name = %q", receiver.Name)
	}
	if len(receiver.SlackConfigs) != 1 {
		t.Errorf("len(SlackConfigs) = %d, want 1", len(receiver.SlackConfigs))
	}
	if receiver.SlackConfigs[0].Channel != "#alerts" {
		t.Errorf("Channel = %q", receiver.SlackConfigs[0].Channel)
	}
}

func TestAlertmanagerConfig_PagerDutyReceiver(t *testing.T) {
	receiver := NewAMReceiver("pagerduty-critical").
		WithPagerDutyConfig(NewAMPagerDutyConfig().
			WithRoutingKeySecret("alertmanager-pagerduty", "routing-key").
			WithSeverity("critical"))

	if len(receiver.PagerDutyConfigs) != 1 {
		t.Errorf("len(PagerDutyConfigs) = %d, want 1", len(receiver.PagerDutyConfigs))
	}
}

func TestAlertmanagerConfig_EmailReceiver(t *testing.T) {
	receiver := NewAMReceiver("email-team").
		WithEmailConfig(NewAMEmailConfig().
			WithTo("team@example.com").
			WithFrom("alerts@example.com"))

	if len(receiver.EmailConfigs) != 1 {
		t.Errorf("len(EmailConfigs) = %d, want 1", len(receiver.EmailConfigs))
	}
}

func TestAlertmanagerConfig_WebhookReceiver(t *testing.T) {
	receiver := NewAMReceiver("webhook-custom").
		WithWebhookConfig(NewAMWebhookConfig().
			WithURL("https://example.com/webhook"))

	if len(receiver.WebhookConfigs) != 1 {
		t.Errorf("len(WebhookConfigs) = %d, want 1", len(receiver.WebhookConfigs))
	}
}

func TestAlertmanagerConfig_AddReceiver(t *testing.T) {
	am := AMConfig("alerts", "monitoring").
		AddReceiver(NewAMReceiver("slack")).
		AddReceiver(NewAMReceiver("pagerduty"))

	if len(am.Spec.Receivers) != 2 {
		t.Errorf("len(Receivers) = %d, want 2", len(am.Spec.Receivers))
	}
}

func TestAlertmanagerConfig_Serialize(t *testing.T) {
	route := NewAMRoute("slack-critical").
		WithGroupBy("alertname").
		AddMatcher(AMMatcherEq("severity", "critical"))

	receiver := NewAMReceiver("slack-critical").
		WithSlackConfig(NewAMSlackConfig().
			WithChannel("#alerts"))

	am := AMConfig("team-alerts", "monitoring").
		WithLabels(map[string]string{"alertmanager": "main"}).
		WithRoute(route).
		AddReceiver(receiver)

	data, err := am.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "apiVersion: monitoring.coreos.com/v1alpha1") {
		t.Error("Expected apiVersion")
	}
	if !strings.Contains(yamlStr, "kind: AlertmanagerConfig") {
		t.Error("Expected kind: AlertmanagerConfig")
	}
	if !strings.Contains(yamlStr, "name: team-alerts") {
		t.Error("Expected name: team-alerts")
	}
	if !strings.Contains(yamlStr, "slack-critical") {
		t.Error("Expected receiver name slack-critical")
	}
}

func TestAlertmanagerConfig_FluentAPI(t *testing.T) {
	am := AMConfig("alerts", "monitoring").
		WithLabels(map[string]string{"alertmanager": "main"}).
		WithRoute(NewAMRoute("default").WithGroupBy("alertname")).
		AddReceiver(NewAMReceiver("default"))

	if am.Name != "alerts" {
		t.Error("Fluent API should preserve name")
	}
	if am.Spec.Route == nil {
		t.Error("Fluent API should set route")
	}
	if len(am.Spec.Receivers) != 1 {
		t.Error("Fluent API should add receiver")
	}
}

func TestAlertmanagerConfig_InhibitRule(t *testing.T) {
	rule := NewAMInhibitRule().
		WithSourceMatcher(AMMatcherEq("severity", "critical")).
		WithTargetMatcher(AMMatcherEq("severity", "warning")).
		WithEqual("alertname", "namespace")

	am := AMConfig("alerts", "monitoring").
		AddInhibitRule(rule)

	if len(am.Spec.InhibitRules) != 1 {
		t.Errorf("len(InhibitRules) = %d, want 1", len(am.Spec.InhibitRules))
	}
}

func TestAlertmanagerConfig_MuteTimeInterval(t *testing.T) {
	mute := NewAMMuteTimeInterval("weekends").
		AddTimeInterval(AMTimeInterval{
			Weekdays: []string{"saturday", "sunday"},
		})

	am := AMConfig("alerts", "monitoring").
		AddMuteTimeInterval(mute)

	if len(am.Spec.MuteTimeIntervals) != 1 {
		t.Errorf("len(MuteTimeIntervals) = %d, want 1", len(am.Spec.MuteTimeIntervals))
	}
	if am.Spec.MuteTimeIntervals[0].Name != "weekends" {
		t.Errorf("MuteTimeInterval name = %q", am.Spec.MuteTimeIntervals[0].Name)
	}
}
