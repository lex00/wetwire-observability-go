// Package testrunner provides test execution with AI personas for evaluating observability configurations.
package testrunner

import "strings"

// Persona represents an evaluation perspective for observability configurations.
type Persona struct {
	// ID is the unique identifier for this persona.
	ID string

	// Name is the display name.
	Name string

	// Description explains the persona's perspective.
	Description string

	// Criteria are the evaluation criteria for this persona.
	Criteria []Criterion
}

// Criterion defines a specific evaluation criterion.
type Criterion struct {
	// ID is the unique identifier for this criterion.
	ID string

	// Name is the display name.
	Name string

	// Description explains what this criterion evaluates.
	Description string

	// Weight is the scoring weight (1-10).
	Weight int

	// Category groups related criteria.
	Category string
}

// Built-in personas
var personas = map[string]*Persona{
	"sre": {
		ID:          "sre",
		Name:        "SRE",
		Description: "Site Reliability Engineer focused on service reliability, alerting, and incident response",
		Criteria: []Criterion{
			{ID: "slo-coverage", Name: "SLO Coverage", Description: "Are Service Level Objectives defined and monitored?", Weight: 9, Category: "reliability"},
			{ID: "alerting", Name: "Alerting Rules", Description: "Are there meaningful alerts with proper severity levels?", Weight: 9, Category: "alerting"},
			{ID: "burn-rate", Name: "Burn Rate Alerts", Description: "Are multi-window burn rate alerts configured?", Weight: 7, Category: "alerting"},
			{ID: "runbooks", Name: "Runbook Links", Description: "Do alerts link to runbooks or documentation?", Weight: 6, Category: "alerting"},
			{ID: "paging", Name: "Paging Policy", Description: "Is there a sensible paging escalation policy?", Weight: 8, Category: "alerting"},
			{ID: "dashboards", Name: "Dashboard Coverage", Description: "Are there dashboards for key services?", Weight: 7, Category: "visibility"},
			{ID: "latency", Name: "Latency Monitoring", Description: "Is latency tracked with percentiles (p50, p95, p99)?", Weight: 8, Category: "metrics"},
			{ID: "saturation", Name: "Saturation Metrics", Description: "Are resource saturation metrics monitored?", Weight: 7, Category: "metrics"},
			{ID: "recording-rules", Name: "Recording Rules", Description: "Are recording rules used for expensive queries?", Weight: 5, Category: "performance"},
		},
	},
	"developer": {
		ID:          "developer",
		Name:        "Developer",
		Description: "Application developer focused on debugging, tracing, and application metrics",
		Criteria: []Criterion{
			{ID: "app-metrics", Name: "Application Metrics", Description: "Are custom application metrics exposed?", Weight: 8, Category: "metrics"},
			{ID: "error-tracking", Name: "Error Tracking", Description: "Are errors tracked with context?", Weight: 9, Category: "debugging"},
			{ID: "request-tracing", Name: "Request Tracing", Description: "Can requests be traced through the system?", Weight: 7, Category: "debugging"},
			{ID: "log-correlation", Name: "Log Correlation", Description: "Are logs correlated with metrics/traces?", Weight: 6, Category: "debugging"},
			{ID: "api-metrics", Name: "API Metrics", Description: "Are API endpoints monitored (rate, errors)?", Weight: 8, Category: "metrics"},
			{ID: "db-metrics", Name: "Database Metrics", Description: "Are database queries and connection pools monitored?", Weight: 7, Category: "metrics"},
			{ID: "cache-metrics", Name: "Cache Metrics", Description: "Are cache hit rates and latency tracked?", Weight: 6, Category: "metrics"},
			{ID: "queue-metrics", Name: "Queue Metrics", Description: "Are message queue depths and processing times tracked?", Weight: 6, Category: "metrics"},
		},
	},
	"security": {
		ID:          "security",
		Name:        "Security Analyst",
		Description: "Security analyst focused on security monitoring, compliance, and threat detection",
		Criteria: []Criterion{
			{ID: "auth-metrics", Name: "Auth Metrics", Description: "Are authentication attempts monitored?", Weight: 9, Category: "security"},
			{ID: "auth-failures", Name: "Auth Failure Alerts", Description: "Are there alerts for authentication failures?", Weight: 9, Category: "alerting"},
			{ID: "rate-limiting", Name: "Rate Limiting", Description: "Is rate limiting monitored and alerted?", Weight: 8, Category: "security"},
			{ID: "secrets", Name: "Secret Exposure", Description: "Are there checks for exposed secrets?", Weight: 10, Category: "security"},
			{ID: "audit-logs", Name: "Audit Logging", Description: "Are security-relevant actions logged?", Weight: 8, Category: "compliance"},
			{ID: "tls-monitoring", Name: "TLS Monitoring", Description: "Are TLS certificate expirations monitored?", Weight: 7, Category: "security"},
			{ID: "anomaly-detection", Name: "Anomaly Detection", Description: "Is there anomaly detection for traffic patterns?", Weight: 6, Category: "security"},
			{ID: "access-patterns", Name: "Access Patterns", Description: "Are unusual access patterns detected?", Weight: 7, Category: "security"},
		},
	},
	"beginner": {
		ID:          "beginner",
		Name:        "Beginner",
		Description: "New to observability, focused on basic monitoring setup and best practices",
		Criteria: []Criterion{
			{ID: "basic-metrics", Name: "Basic Metrics", Description: "Are basic request/error/duration metrics present?", Weight: 10, Category: "basics"},
			{ID: "up-check", Name: "Up Check", Description: "Is there an 'up' metric for health checking?", Weight: 9, Category: "basics"},
			{ID: "scrape-targets", Name: "Scrape Targets", Description: "Are scrape targets properly configured?", Weight: 9, Category: "basics"},
			{ID: "labels", Name: "Consistent Labels", Description: "Are labels consistent across metrics?", Weight: 7, Category: "basics"},
			{ID: "documentation", Name: "Documentation", Description: "Are metrics and alerts documented?", Weight: 6, Category: "basics"},
			{ID: "simple-alerts", Name: "Simple Alerts", Description: "Are there basic health alerts?", Weight: 8, Category: "alerting"},
		},
	},
}

// GetPersona returns a persona by ID.
func GetPersona(id string) *Persona {
	return personas[strings.ToLower(id)]
}

// GetAllPersonas returns all available personas.
func GetAllPersonas() []*Persona {
	result := make([]*Persona, 0, len(personas))
	for _, p := range personas {
		result = append(result, p)
	}
	return result
}

// ListPersonaNames returns the IDs of all available personas.
func ListPersonaNames() []string {
	names := make([]string, 0, len(personas))
	for id := range personas {
		names = append(names, id)
	}
	return names
}
