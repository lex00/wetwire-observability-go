// Package monitoring demonstrates Alertmanager routing patterns.
//
// This example shows multi-team alert routing with Slack, PagerDuty, and email
// receivers, including routing trees, inhibition rules, and mute time intervals.
package monitoring

import "github.com/lex00/wetwire-observability-go/alertmanager"

// GlobalSettings configures global Alertmanager settings.
var GlobalSettings = alertmanager.GlobalConfig{
	SMTPSmarthost:    "smtp.example.com:587",
	SMTPFrom:         "alertmanager@example.com",
	SMTPAuthUsername: "alertmanager",
	SMTPAuthPassword: "${SMTP_PASSWORD}",
	ResolveTimeout:   5 * alertmanager.Minute,
}

// Config is the complete Alertmanager configuration.
var Config = alertmanager.AlertmanagerConfig{
	Global:            &GlobalSettings,
	Route:             RootRoute,
	Receivers:         Receivers,
	InhibitRules:      InhibitionRules,
	MuteTimeIntervals: MuteIntervals,
}
