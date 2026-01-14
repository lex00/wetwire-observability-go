package design

import (
	"strings"
	"testing"
)

func TestNewPromptBuilder(t *testing.T) {
	pb := NewPromptBuilder()
	if pb == nil {
		t.Fatal("NewPromptBuilder returned nil")
	}
}

func TestPromptBuilder_SystemPrompt(t *testing.T) {
	pb := NewPromptBuilder()
	system := pb.SystemPrompt()

	// Should include key observability concepts
	if !strings.Contains(system, "Prometheus") {
		t.Error("system prompt should mention Prometheus")
	}
	if !strings.Contains(system, "Alertmanager") {
		t.Error("system prompt should mention Alertmanager")
	}
	if !strings.Contains(system, "Grafana") {
		t.Error("system prompt should mention Grafana")
	}
	if !strings.Contains(system, "wetwire") {
		t.Error("system prompt should mention wetwire")
	}
}

func TestPromptBuilder_BuildUserPrompt(t *testing.T) {
	pb := NewPromptBuilder()
	prompt := pb.BuildUserPrompt("Add monitoring for an API server")

	if !strings.Contains(prompt, "API server") {
		t.Error("user prompt should include the request")
	}
}

func TestPromptBuilder_WithContext(t *testing.T) {
	pb := NewPromptBuilder().
		WithContext("existing scrape config", "var APIScrape = prometheus.ScrapeConfig{}")

	prompt := pb.BuildUserPrompt("Add more scrapes")

	if !strings.Contains(prompt, "existing scrape config") {
		t.Error("prompt should include context")
	}
}

func TestPromptBuilder_ForPrometheus(t *testing.T) {
	pb := NewPromptBuilder().ForPrometheus()
	system := pb.SystemPrompt()

	if !strings.Contains(system, "scrape") {
		t.Error("Prometheus prompt should mention scrape configs")
	}
}

func TestPromptBuilder_ForAlertmanager(t *testing.T) {
	pb := NewPromptBuilder().ForAlertmanager()
	system := pb.SystemPrompt()

	lower := strings.ToLower(system)
	if !strings.Contains(lower, "route") || !strings.Contains(lower, "receiver") {
		t.Error("Alertmanager prompt should mention routes and receivers")
	}
}

func TestPromptBuilder_ForGrafana(t *testing.T) {
	pb := NewPromptBuilder().ForGrafana()
	system := pb.SystemPrompt()

	if !strings.Contains(system, "dashboard") || !strings.Contains(system, "panel") {
		t.Error("Grafana prompt should mention dashboards and panels")
	}
}

func TestPromptBuilder_ForRules(t *testing.T) {
	pb := NewPromptBuilder().ForRules()
	system := pb.SystemPrompt()

	if !strings.Contains(system, "alert") || !strings.Contains(system, "recording") {
		t.Error("Rules prompt should mention alerting and recording rules")
	}
}

func TestPromptBuilder_FullRequest(t *testing.T) {
	pb := NewPromptBuilder().
		ForPrometheus().
		WithContext("package", "monitoring").
		WithContext("imports", "prometheus")

	system := pb.SystemPrompt()
	user := pb.BuildUserPrompt("Add kubernetes service discovery")

	// System should have base + prometheus focus
	if !strings.Contains(system, "Prometheus") {
		t.Error("system should mention Prometheus")
	}

	// User should have context + request
	if !strings.Contains(user, "kubernetes") {
		t.Error("user should include request")
	}
}
