package monitoring

import "github.com/lex00/wetwire-observability-go/alertmanager"

// Global Configuration

// GlobalConfig defines global Alertmanager settings.
var GlobalConfig = &alertmanager.GlobalConfig{
	SMTPSmarthost:    "smtp.example.com:587",
	SMTPFrom:         "alertmanager@example.com",
	SMTPAuthUsername: "alertmanager",
	SMTPAuthPassword: "${SMTP_PASSWORD}",
	ResolveTimeout:   5 * alertmanager.Minute,
}

// Receivers

// DefaultReceiver catches unrouted alerts.
var DefaultReceiver = alertmanager.NewReceiver("default").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-general").
			WithTitle("[{{ .Status | toUpper }}] {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.summary }}").
			WithSendResolved(true),
	)

// PlatformSlackReceiver sends platform team alerts to Slack.
var PlatformSlackReceiver = alertmanager.NewReceiver("platform-slack").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-platform").
			WithTitle("[{{ .CommonLabels.severity | toUpper }}] {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.description }}").
			WithColor("{{ if eq .CommonLabels.severity \"critical\" }}danger{{ else if eq .CommonLabels.severity \"warning\" }}warning{{ else }}good{{ end }}").
			WithSendResolved(true).
			WithFields(
				alertmanager.NewSlackField("Service", "{{ .CommonLabels.service }}").WithShort(true),
				alertmanager.NewSlackField("Severity", "{{ .CommonLabels.severity }}").WithShort(true),
			).
			WithActions(
				alertmanager.NewSlackAction("Dashboard", "{{ .CommonAnnotations.dashboard }}"),
				alertmanager.NewSlackAction("Runbook", "{{ .CommonAnnotations.runbook_url }}"),
			),
	)

// PlatformPagerDutyReceiver pages platform team for critical alerts.
var PlatformPagerDutyReceiver = alertmanager.NewReceiver("platform-pagerduty").
	WithPagerDutyConfigs(
		alertmanager.NewPagerDutyConfig().
			WithRoutingKeyFile("/etc/alertmanager/secrets/pagerduty-platform").
			WithDescription("{{ .CommonAnnotations.summary }}").
			WithSeverity("{{ .CommonLabels.severity }}").
			WithComponent("{{ .CommonLabels.service }}").
			WithGroup("platform").
			WithSendResolved(true),
	)

// BackendSlackReceiver sends backend team alerts to Slack.
var BackendSlackReceiver = alertmanager.NewReceiver("backend-slack").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-backend").
			WithTitle("[{{ .CommonLabels.severity | toUpper }}] {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.description }}").
			WithColor("{{ if eq .CommonLabels.severity \"critical\" }}danger{{ else if eq .CommonLabels.severity \"warning\" }}warning{{ else }}good{{ end }}").
			WithSendResolved(true).
			WithFields(
				alertmanager.NewSlackField("Service", "{{ .CommonLabels.service }}").WithShort(true),
				alertmanager.NewSlackField("Severity", "{{ .CommonLabels.severity }}").WithShort(true),
			),
	)

// BackendPagerDutyReceiver pages backend team for critical alerts.
var BackendPagerDutyReceiver = alertmanager.NewReceiver("backend-pagerduty").
	WithPagerDutyConfigs(
		alertmanager.NewPagerDutyConfig().
			WithRoutingKeyFile("/etc/alertmanager/secrets/pagerduty-backend").
			WithDescription("{{ .CommonAnnotations.summary }}").
			WithSeverity("{{ .CommonLabels.severity }}").
			WithComponent("{{ .CommonLabels.service }}").
			WithGroup("backend").
			WithSendResolved(true),
	)

// SRESlackReceiver sends SRE alerts to Slack.
var SRESlackReceiver = alertmanager.NewReceiver("sre-slack").
	WithSlackConfigs(
		alertmanager.NewSlackConfig().
			WithChannel("#alerts-sre").
			WithTitle(":rotating_light: {{ .CommonLabels.alertname }}").
			WithText("{{ .CommonAnnotations.description }}").
			WithColor("danger").
			WithSendResolved(true),
	)

// NullReceiver silently drops alerts.
var NullReceiver = alertmanager.NewReceiver("null")

// Routing Configuration

// RootRoute is the top-level routing tree.
var RootRoute = alertmanager.NewRoute("default").
	WithGroupBy("alertname", "service", "cluster").
	WithGroupWait(30 * alertmanager.Second).
	WithGroupInterval(5 * alertmanager.Minute).
	WithRepeatInterval(4 * alertmanager.Hour).
	WithRoutes(
		// SLO alerts go to SRE
		SLORoute,
		// Critical alerts get paged
		CriticalRoute,
		// Team-specific routing
		PlatformRoute,
		BackendRoute,
		// Info alerts are dropped
		InfoRoute,
	)

// SLORoute handles SLO-based alerts.
var SLORoute = alertmanager.NewRoute("sre-slack").
	WithMatchers(alertmanager.Eq("slo", "availability")).
	WithGroupWait(10 * alertmanager.Second).
	WithRepeatInterval(30 * alertmanager.Minute)

// CriticalRoute pages on-call for critical alerts.
var CriticalRoute = alertmanager.NewRoute("platform-pagerduty").
	Severity("critical").
	WithGroupWait(10 * alertmanager.Second).
	WithRepeatInterval(1 * alertmanager.Hour).
	WithContinue(true) // Continue to team routes

// PlatformRoute handles platform team alerts.
var PlatformRoute = alertmanager.NewRoute("platform-slack").
	Team("platform").
	WithRoutes(
		alertmanager.NewRoute("platform-pagerduty").
			Severity("critical"),
	)

// BackendRoute handles backend team alerts.
var BackendRoute = alertmanager.NewRoute("backend-slack").
	Team("backend").
	WithRoutes(
		alertmanager.NewRoute("backend-pagerduty").
			Severity("critical"),
	)

// InfoRoute drops info-level alerts.
var InfoRoute = alertmanager.NewRoute("null").
	Severity("info")

// Inhibition Rules

// InhibitRules define alert suppression.
var InhibitRules = []*alertmanager.InhibitRule{
	// Critical inhibits warning for same alert
	alertmanager.NewInhibitRule().
		WithSourceMatchers(alertmanager.SeverityCritical()).
		WithTargetMatchers(alertmanager.SeverityWarning()).
		WithEqual("alertname", "service"),
	// Critical inhibits info
	alertmanager.NewInhibitRule().
		WithSourceMatchers(alertmanager.SeverityCritical()).
		WithTargetMatchers(alertmanager.SeverityInfo()).
		WithEqual("alertname", "service"),
	// Warning inhibits info
	alertmanager.NewInhibitRule().
		WithSourceMatchers(alertmanager.SeverityWarning()).
		WithTargetMatchers(alertmanager.SeverityInfo()).
		WithEqual("alertname", "service"),
	// ServiceDown inhibits all other alerts for that service
	alertmanager.NewInhibitRule().
		WithSourceMatchers(alertmanager.Alertname("ServiceDown")).
		WithTargetMatchers(alertmanager.NotEq("alertname", "ServiceDown")).
		WithEqual("service"),
}

// Mute Time Intervals

// MaintenanceWindow mutes alerts during maintenance.
var MaintenanceWindow = &alertmanager.MuteTimeInterval{
	Name: "maintenance",
	TimeIntervals: []alertmanager.TimeInterval{
		{
			Weekdays: []alertmanager.WeekdayRange{"sunday"},
			Times: []alertmanager.TimeRange{
				{StartTime: "02:00", EndTime: "06:00"},
			},
		},
	},
}

// AlertmanagerConfig is the complete configuration.
var AlertmanagerConfig = alertmanager.AlertmanagerConfig{
	Global: GlobalConfig,
	Route:  RootRoute,
	Receivers: []*alertmanager.Receiver{
		DefaultReceiver,
		PlatformSlackReceiver,
		PlatformPagerDutyReceiver,
		BackendSlackReceiver,
		BackendPagerDutyReceiver,
		SRESlackReceiver,
		NullReceiver,
	},
	InhibitRules:      InhibitRules,
	MuteTimeIntervals: []*alertmanager.MuteTimeInterval{MaintenanceWindow},
}
