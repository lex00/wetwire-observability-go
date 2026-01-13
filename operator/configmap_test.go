package operator

import (
	"strings"
	"testing"

	"github.com/lex00/wetwire-observability-go/grafana"
)

func TestConfigMap(t *testing.T) {
	cm := ConfigMap("my-dashboard", "monitoring")
	if cm.Name != "my-dashboard" {
		t.Errorf("Name = %q, want my-dashboard", cm.Name)
	}
	if cm.Namespace != "monitoring" {
		t.Errorf("Namespace = %q, want monitoring", cm.Namespace)
	}
	if cm.Kind != "ConfigMap" {
		t.Errorf("Kind = %q, want ConfigMap", cm.Kind)
	}
	if cm.APIVersion != "v1" {
		t.Errorf("APIVersion = %q, want v1", cm.APIVersion)
	}
}

func TestConfigMap_WithLabels(t *testing.T) {
	cm := ConfigMap("dashboard", "default").
		WithLabels(map[string]string{"app": "grafana"})
	if cm.Labels["app"] != "grafana" {
		t.Errorf("Labels[app] = %q, want grafana", cm.Labels["app"])
	}
}

func TestConfigMap_AddLabel(t *testing.T) {
	cm := ConfigMap("dashboard", "default").
		AddLabel("team", "platform").
		AddLabel("env", "prod")
	if len(cm.Labels) != 2 {
		t.Errorf("len(Labels) = %d, want 2", len(cm.Labels))
	}
}

func TestConfigMap_WithData(t *testing.T) {
	cm := ConfigMap("config", "default").
		WithData(map[string]string{
			"key1": "value1",
			"key2": "value2",
		})
	if len(cm.Data) != 2 {
		t.Errorf("len(Data) = %d, want 2", len(cm.Data))
	}
}

func TestConfigMap_AddData(t *testing.T) {
	cm := ConfigMap("config", "default").
		AddData("prometheus.yml", "global: {}")
	if cm.Data["prometheus.yml"] != "global: {}" {
		t.Errorf("Data[prometheus.yml] = %q", cm.Data["prometheus.yml"])
	}
}

func TestDashboardConfigMap(t *testing.T) {
	dashboard := grafana.NewDashboard("test-uid", "Test Dashboard").
		WithTags("test", "example")

	cm := DashboardConfigMap("test-dashboard", "monitoring", dashboard)

	if cm.Name != "test-dashboard" {
		t.Errorf("Name = %q, want test-dashboard", cm.Name)
	}
	if cm.Labels["grafana_dashboard"] != "1" {
		t.Errorf("Labels[grafana_dashboard] = %q, want 1", cm.Labels["grafana_dashboard"])
	}
	if _, ok := cm.Data["test-dashboard.json"]; !ok {
		t.Error("Expected dashboard JSON in Data")
	}
}

func TestDashboardConfigMap_WithFolder(t *testing.T) {
	dashboard := grafana.NewDashboard("uid", "Dashboard")

	cm := DashboardConfigMap("dashboard", "monitoring", dashboard).
		WithFolder("Platform")

	if cm.Metadata.Annotations["grafana_folder"] != "Platform" {
		t.Errorf("Annotations[grafana_folder] = %q, want Platform", cm.Metadata.Annotations["grafana_folder"])
	}
}

func TestDashboardConfigMap_WithCustomLabel(t *testing.T) {
	dashboard := grafana.NewDashboard("uid", "Dashboard")

	cm := DashboardConfigMap("dashboard", "monitoring", dashboard).
		WithGrafanaLabel("grafana_dashboard", "true")

	if cm.Labels["grafana_dashboard"] != "true" {
		t.Errorf("Labels[grafana_dashboard] = %q, want true", cm.Labels["grafana_dashboard"])
	}
}

func TestDashboardConfigMap_Serialize(t *testing.T) {
	dashboard := grafana.NewDashboard("api-overview", "API Overview").
		WithTags("api", "overview")

	cm := DashboardConfigMap("api-dashboard", "monitoring", dashboard).
		WithFolder("Platform")

	data, err := cm.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "apiVersion: v1") {
		t.Error("Expected apiVersion: v1")
	}
	if !strings.Contains(yamlStr, "kind: ConfigMap") {
		t.Error("Expected kind: ConfigMap")
	}
	if !strings.Contains(yamlStr, "name: api-dashboard") {
		t.Error("Expected name: api-dashboard")
	}
	if !strings.Contains(yamlStr, "grafana_dashboard") {
		t.Error("Expected grafana_dashboard label")
	}
	if !strings.Contains(yamlStr, "grafana_folder") {
		t.Error("Expected grafana_folder annotation")
	}
}

func TestConfigMap_FluentAPI(t *testing.T) {
	cm := ConfigMap("config", "monitoring").
		WithLabels(map[string]string{"app": "test"}).
		WithAnnotations(map[string]string{"note": "test"}).
		AddData("key", "value")

	if cm.Name != "config" {
		t.Error("Fluent API should preserve name")
	}
	if len(cm.Data) != 1 {
		t.Error("Fluent API should add data")
	}
}

func TestConfigMap_ForGrafanaSidecar(t *testing.T) {
	cm := ConfigMap("dashboard", "monitoring").
		ForGrafanaSidecar()

	if cm.Labels["grafana_dashboard"] != "1" {
		t.Errorf("Labels[grafana_dashboard] = %q, want 1", cm.Labels["grafana_dashboard"])
	}
}
