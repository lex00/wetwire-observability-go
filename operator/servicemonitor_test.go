package operator

import (
	"strings"
	"testing"
)

func TestServiceMon(t *testing.T) {
	sm := ServiceMon("api-server", "monitoring")
	if sm.Name != "api-server" {
		t.Errorf("Name = %q, want api-server", sm.Name)
	}
	if sm.Namespace != "monitoring" {
		t.Errorf("Namespace = %q, want monitoring", sm.Namespace)
	}
	if sm.Kind != "ServiceMonitor" {
		t.Errorf("Kind = %q, want ServiceMonitor", sm.Kind)
	}
	if sm.APIVersion != "monitoring.coreos.com/v1" {
		t.Errorf("APIVersion = %q", sm.APIVersion)
	}
}

func TestServiceMonitor_WithLabels(t *testing.T) {
	sm := ServiceMon("api", "default").
		WithLabels(map[string]string{"team": "backend"})
	if sm.Labels["team"] != "backend" {
		t.Errorf("Labels[team] = %q, want backend", sm.Labels["team"])
	}
}

func TestServiceMonitor_SelectService(t *testing.T) {
	sm := ServiceMon("api", "default").
		SelectService("app", "api-server")
	if sm.Spec.Selector.MatchLabels["app"] != "api-server" {
		t.Error("Selector.MatchLabels should contain app=api-server")
	}
}

func TestServiceMonitor_SelectServiceByLabels(t *testing.T) {
	sm := ServiceMon("api", "default").
		SelectServiceByLabels(map[string]string{"app": "api", "version": "v1"})
	if len(sm.Spec.Selector.MatchLabels) != 2 {
		t.Errorf("len(MatchLabels) = %d, want 2", len(sm.Spec.Selector.MatchLabels))
	}
}

func TestServiceMonitor_InNamespace(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		InNamespace("production")
	if len(sm.Spec.NamespaceSelector.MatchNames) != 1 {
		t.Errorf("len(MatchNames) = %d, want 1", len(sm.Spec.NamespaceSelector.MatchNames))
	}
	if sm.Spec.NamespaceSelector.MatchNames[0] != "production" {
		t.Error("NamespaceSelector should contain production")
	}
}

func TestServiceMonitor_InNamespaces(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		InNamespaces("prod", "staging")
	if len(sm.Spec.NamespaceSelector.MatchNames) != 2 {
		t.Errorf("len(MatchNames) = %d, want 2", len(sm.Spec.NamespaceSelector.MatchNames))
	}
}

func TestServiceMonitor_InAllNamespaces(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		InAllNamespaces()
	if !sm.Spec.NamespaceSelector.Any {
		t.Error("NamespaceSelector.Any should be true")
	}
}

func TestServiceMonitor_WithEndpoint(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		WithEndpoint(NewEndpoint("metrics").WithInterval("30s"))
	if len(sm.Spec.Endpoints) != 1 {
		t.Errorf("len(Endpoints) = %d, want 1", len(sm.Spec.Endpoints))
	}
}

func TestServiceMonitor_AddEndpoint(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		AddEndpoint(NewEndpoint("metrics")).
		AddEndpoint(NewEndpoint("health"))
	if len(sm.Spec.Endpoints) != 2 {
		t.Errorf("len(Endpoints) = %d, want 2", len(sm.Spec.Endpoints))
	}
}

func TestServiceMonitor_WithJobLabel(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		WithJobLabel("app.kubernetes.io/name")
	if sm.Spec.JobLabel != "app.kubernetes.io/name" {
		t.Errorf("JobLabel = %q", sm.Spec.JobLabel)
	}
}

func TestServiceMonitor_Serialize(t *testing.T) {
	sm := ServiceMon("api-server", "monitoring").
		SelectService("app", "api").
		InNamespace("production").
		WithEndpoint(NewEndpoint("http").WithPath("/metrics"))

	data, err := sm.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "apiVersion: monitoring.coreos.com/v1") {
		t.Error("Expected apiVersion")
	}
	if !strings.Contains(yamlStr, "kind: ServiceMonitor") {
		t.Error("Expected kind: ServiceMonitor")
	}
	if !strings.Contains(yamlStr, "name: api-server") {
		t.Error("Expected name: api-server")
	}
	if !strings.Contains(yamlStr, "app: api") {
		t.Error("Expected selector with app: api")
	}
}

func TestServiceMonitor_FluentAPI(t *testing.T) {
	sm := ServiceMon("api", "monitoring").
		WithLabels(map[string]string{"prometheus": "main"}).
		SelectService("app", "api").
		InNamespaces("prod", "staging").
		WithEndpoint(NewEndpoint("metrics").WithInterval("30s")).
		WithJobLabel("app")

	if sm.Name != "api" {
		t.Error("Fluent API should preserve name")
	}
	if len(sm.Spec.Endpoints) != 1 {
		t.Error("Fluent API should add endpoint")
	}
}
