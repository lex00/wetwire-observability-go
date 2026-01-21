package monitoring

import "github.com/lex00/wetwire-observability-go/alertmanager"

// RootRoute is the top-level routing configuration.
// Alerts are routed based on team, severity, and service labels.
var RootRoute = alertmanager.NewRoute("default").
	WithGroupBy("alertname", "cluster", "service").
	WithGroupWait(30 * alertmanager.Second).
	WithGroupInterval(5 * alertmanager.Minute).
	WithRepeatInterval(4 * alertmanager.Hour).
	WithRoutes(
		// Critical alerts always go to PagerDuty in addition to Slack
		CriticalRoute,
		// Team-specific routing
		PlatformTeamRoute,
		DatabaseTeamRoute,
		SecurityTeamRoute,
		// Test/development alerts go to null receiver
		TestRoute,
	)

// CriticalRoute routes all critical alerts to PagerDuty.
var CriticalRoute = alertmanager.NewRoute("platform-pagerduty").
	Severity("critical").
	WithGroupWait(10 * alertmanager.Second).
	WithRepeatInterval(1 * alertmanager.Hour).
	WithContinue(true) // Continue to allow team-specific routing

// PlatformTeamRoute routes platform team alerts.
var PlatformTeamRoute = alertmanager.NewRoute("platform-slack").
	Team("platform").
	WithRoutes(
		// Critical platform alerts also page
		alertmanager.NewRoute("platform-pagerduty").
			Severity("critical"),
	)

// DatabaseTeamRoute routes database team alerts.
var DatabaseTeamRoute = alertmanager.NewRoute("database-slack").
	Team("database").
	WithRoutes(
		// Critical database alerts also page
		alertmanager.NewRoute("database-pagerduty").
			Severity("critical").
			WithRepeatInterval(30 * alertmanager.Minute),
	)

// SecurityTeamRoute routes security-related alerts.
var SecurityTeamRoute = alertmanager.NewRoute("security-slack").
	Team("security").
	WithContinue(true).
	WithRoutes(
		// All security alerts also go to email
		alertmanager.NewRoute("security-email"),
	)

// TestRoute drops alerts from test/development environments.
var TestRoute = alertmanager.NewRoute("null").
	Environment("test").
	WithMatchers(
		alertmanager.Regex("alertname", "Test.*"),
	)

// BusinessHoursRoute demonstrates time-based routing.
var BusinessHoursRoute = alertmanager.NewRoute("platform-slack").
	Team("platform").
	Severity("warning").
	WithMuteTimeIntervals("outside-business-hours").
	WithRoutes(
		// During business hours, warning alerts go to Slack
		alertmanager.NewRoute("platform-slack").
			WithActiveTimeIntervals("business-hours"),
	)

// ServiceSpecificRoutes demonstrates routing by service.
var ServiceSpecificRoutes = alertmanager.NewRoute("default").
	WithRoutes(
		// Payment service alerts have shorter group wait
		alertmanager.NewRoute("platform-pagerduty").
			Service("payment-service").
			WithGroupWait(10 * alertmanager.Second),
		// Auth service alerts
		alertmanager.NewRoute("security-slack").
			Service("auth-service"),
	)
