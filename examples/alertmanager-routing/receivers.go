package monitoring

import "github.com/lex00/wetwire-observability-go/alertmanager"

// Default Receiver - catches all unrouted alerts

// DefaultReceiver sends alerts to the general ops channel.
var DefaultReceiver = alertmanager.NewReceiver("default").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-general").
			WithTitle("{{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.summary }}").
			WithColor("{{ if eq .Status \"firing\" }}danger{{ else }}good{{ end }}").
			WithSendResolved(true),
	)

// Platform Team Receivers

// PlatformSlackReceiver sends platform alerts to the platform Slack channel.
var PlatformSlackReceiver = alertmanager.NewReceiver("platform-slack").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-platform").
			WithUsername("AlertBot").
			WithIconEmoji(":robot_face:").
			WithTitle("[{{ .Status | toUpper }}] {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.description }}").
			WithColor("{{ if eq .Status \"firing\" }}{{ if eq .CommonLabels.severity \"critical\" }}danger{{ else }}warning{{ end }}{{ else }}good{{ end }}").
			WithSendResolved(true).
			WithFields(
				alertmanager.NewSlackField("Severity", "{{ .CommonLabels.severity }}").WithShort(true),
				alertmanager.NewSlackField("Service", "{{ .CommonLabels.service }}").WithShort(true),
			).
			WithActions(
				alertmanager.NewSlackAction("View Runbook", "{{ .CommonAnnotations.runbook_url }}"),
				alertmanager.NewSlackAction("Silence", "{{ .ExternalURL }}/#/silences/new?filter=%7Balertname%3D%22{{ .CommonLabels.alertname }}%22%7D"),
			),
	)

// PlatformPagerDutyReceiver sends critical platform alerts to PagerDuty.
var PlatformPagerDutyReceiver = alertmanager.NewReceiver("platform-pagerduty").
	WithPagerDutyConfigs(
		alertmanager.NewPagerDutyConfig().
			WithRoutingKeyFile("/etc/alertmanager/secrets/pagerduty-platform-key").
			WithClient("Alertmanager").
			WithClientURL("{{ .ExternalURL }}").
			WithDescription("{{ .CommonAnnotations.summary }}").
			WithSeverity("{{ .CommonLabels.severity }}").
			WithComponent("{{ .CommonLabels.service }}").
			WithGroup("platform").
			WithSendResolved(true),
	)

// Database Team Receivers

// DatabaseSlackReceiver sends database alerts to the DBA Slack channel.
var DatabaseSlackReceiver = alertmanager.NewReceiver("database-slack").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-database").
			WithTitle("[{{ .CommonLabels.severity | toUpper }}] {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.description }}").
			WithColor("{{ if eq .CommonLabels.severity \"critical\" }}danger{{ else if eq .CommonLabels.severity \"warning\" }}warning{{ else }}good{{ end }}").
			WithSendResolved(true).
			WithFields(
				alertmanager.NewSlackField("Database", "{{ .CommonLabels.datname }}").WithShort(true),
				alertmanager.NewSlackField("Instance", "{{ .CommonLabels.instance }}").WithShort(true),
			),
	)

// DatabasePagerDutyReceiver sends critical database alerts to PagerDuty.
var DatabasePagerDutyReceiver = alertmanager.NewReceiver("database-pagerduty").
	WithPagerDutyConfigs(
		alertmanager.NewPagerDutyConfig().
			WithRoutingKeyFile("/etc/alertmanager/secrets/pagerduty-database-key").
			WithDescription("{{ .CommonAnnotations.summary }}").
			WithSeverity("critical").
			WithGroup("database").
			WithSendResolved(true),
	)

// Security Team Receivers

// SecuritySlackReceiver sends security alerts to the security channel.
var SecuritySlackReceiver = alertmanager.NewReceiver("security-slack").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-security").
			WithTitle(":lock: {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.description }}").
			WithColor("danger").
			WithSendResolved(true),
	)

// SecurityEmailReceiver sends security alerts via email.
var SecurityEmailReceiver = alertmanager.NewReceiver("security-email").
	WithEmailConfigs(
		alertmanager.NewEmailConfig().
			WithTo("security-oncall@example.com").
			WithSendResolved(true),
	)

// Null Receiver - silently drops alerts

// NullReceiver drops alerts (for testing or intentionally suppressed alerts).
var NullReceiver = alertmanager.NewReceiver("null")

// Receivers is the list of all receivers.
var Receivers = []*alertmanager.Receiver{
	DefaultReceiver,
	PlatformSlackReceiver,
	PlatformPagerDutyReceiver,
	DatabaseSlackReceiver,
	DatabasePagerDutyReceiver,
	SecuritySlackReceiver,
	SecurityEmailReceiver,
	NullReceiver,
}
