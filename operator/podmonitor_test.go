package operator

import (
	"strings"
	"testing"
)

func TestPodMon(t *testing.T) {
	pm := PodMon("api-pods", "monitoring")
	if pm.Name != "api-pods" {
		t.Errorf("Name = %q, want api-pods", pm.Name)
	}
	if pm.Namespace != "monitoring" {
		t.Errorf("Namespace = %q, want monitoring", pm.Namespace)
	}
	if pm.Kind != "PodMonitor" {
		t.Errorf("Kind = %q, want PodMonitor", pm.Kind)
	}
	if pm.APIVersion != "monitoring.coreos.com/v1" {
		t.Errorf("APIVersion = %q", pm.APIVersion)
	}
}

func TestPodMonitor_WithLabels(t *testing.T) {
	pm := PodMon("api", "default").
		WithLabels(map[string]string{"team": "backend"})
	if pm.Labels["team"] != "backend" {
		t.Errorf("Labels[team] = %q, want backend", pm.Labels["team"])
	}
}

func TestPodMonitor_SelectPods(t *testing.T) {
	pm := PodMon("api", "default").
		SelectPods("app", "api-server")
	if pm.Spec.Selector.MatchLabels["app"] != "api-server" {
		t.Error("Selector.MatchLabels should contain app=api-server")
	}
}

func TestPodMonitor_SelectPodsByLabels(t *testing.T) {
	pm := PodMon("api", "default").
		SelectPodsByLabels(map[string]string{"app": "api", "version": "v1"})
	if len(pm.Spec.Selector.MatchLabels) != 2 {
		t.Errorf("len(MatchLabels) = %d, want 2", len(pm.Spec.Selector.MatchLabels))
	}
}

func TestPodMonitor_InNamespace(t *testing.T) {
	pm := PodMon("api", "monitoring").
		InNamespace("production")
	if len(pm.Spec.NamespaceSelector.MatchNames) != 1 {
		t.Errorf("len(MatchNames) = %d, want 1", len(pm.Spec.NamespaceSelector.MatchNames))
	}
}

func TestPodMonitor_InNamespaces(t *testing.T) {
	pm := PodMon("api", "monitoring").
		InNamespaces("prod", "staging")
	if len(pm.Spec.NamespaceSelector.MatchNames) != 2 {
		t.Errorf("len(MatchNames) = %d, want 2", len(pm.Spec.NamespaceSelector.MatchNames))
	}
}

func TestPodMonitor_InAllNamespaces(t *testing.T) {
	pm := PodMon("api", "monitoring").
		InAllNamespaces()
	if !pm.Spec.NamespaceSelector.Any {
		t.Error("NamespaceSelector.Any should be true")
	}
}

func TestPodMonitor_WithPodMetricsEndpoint(t *testing.T) {
	pm := PodMon("api", "monitoring").
		WithPodMetricsEndpoint(NewPodMetricsEndpoint("metrics").WithInterval("30s"))
	if len(pm.Spec.PodMetricsEndpoints) != 1 {
		t.Errorf("len(PodMetricsEndpoints) = %d, want 1", len(pm.Spec.PodMetricsEndpoints))
	}
}

func TestPodMonitor_AddPodMetricsEndpoint(t *testing.T) {
	pm := PodMon("api", "monitoring").
		AddPodMetricsEndpoint(NewPodMetricsEndpoint("metrics")).
		AddPodMetricsEndpoint(NewPodMetricsEndpoint("health"))
	if len(pm.Spec.PodMetricsEndpoints) != 2 {
		t.Errorf("len(PodMetricsEndpoints) = %d, want 2", len(pm.Spec.PodMetricsEndpoints))
	}
}

func TestPodMonitor_WithJobLabel(t *testing.T) {
	pm := PodMon("api", "monitoring").
		WithJobLabel("app.kubernetes.io/name")
	if pm.Spec.JobLabel != "app.kubernetes.io/name" {
		t.Errorf("JobLabel = %q", pm.Spec.JobLabel)
	}
}

func TestPodMonitor_Serialize(t *testing.T) {
	pm := PodMon("api-pods", "monitoring").
		SelectPods("app", "api").
		InNamespace("production").
		WithPodMetricsEndpoint(NewPodMetricsEndpoint("http").WithPath("/metrics"))

	data, err := pm.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "apiVersion: monitoring.coreos.com/v1") {
		t.Error("Expected apiVersion")
	}
	if !strings.Contains(yamlStr, "kind: PodMonitor") {
		t.Error("Expected kind: PodMonitor")
	}
	if !strings.Contains(yamlStr, "name: api-pods") {
		t.Error("Expected name: api-pods")
	}
}

func TestPodMonitor_FluentAPI(t *testing.T) {
	pm := PodMon("api", "monitoring").
		WithLabels(map[string]string{"prometheus": "main"}).
		SelectPods("app", "api").
		InNamespaces("prod", "staging").
		WithPodMetricsEndpoint(NewPodMetricsEndpoint("metrics").WithInterval("30s")).
		WithJobLabel("app")

	if pm.Name != "api" {
		t.Error("Fluent API should preserve name")
	}
	if len(pm.Spec.PodMetricsEndpoints) != 1 {
		t.Error("Fluent API should add endpoint")
	}
}

func TestNewPodMetricsEndpoint(t *testing.T) {
	e := NewPodMetricsEndpoint("metrics")
	if e.Port != "metrics" {
		t.Errorf("Port = %q, want metrics", e.Port)
	}
}

func TestPodMetricsEndpoint_FluentAPI(t *testing.T) {
	e := NewPodMetricsEndpoint("metrics").
		WithPath("/metrics").
		WithInterval("30s").
		WithScrapeTimeout("10s")

	if e.Port != "metrics" {
		t.Error("Fluent API should preserve port")
	}
	if e.Path != "/metrics" {
		t.Error("Fluent API should set path")
	}
}
