package grafana

import (
	"strings"
	"testing"
)

func TestDataSourceProvisioning_Serialize(t *testing.T) {
	provisioning := NewDataSourceProvisioning("default").
		AddDataSource(PrometheusDataSource("prometheus", "http://prometheus:9090").AsDefault()).
		AddDataSource(LokiDataSource("loki", "http://loki:3100"))

	data, err := provisioning.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	// Check required fields
	if !strings.Contains(yamlStr, "apiVersion: 1") {
		t.Error("Expected apiVersion: 1")
	}
	if !strings.Contains(yamlStr, "name: prometheus") {
		t.Error("Expected name: prometheus")
	}
	if !strings.Contains(yamlStr, "type: prometheus") {
		t.Error("Expected type: prometheus")
	}
	if !strings.Contains(yamlStr, "name: loki") {
		t.Error("Expected name: loki")
	}
}

func TestDataSourceProvisioning_DeleteAllExisting(t *testing.T) {
	provisioning := NewDataSourceProvisioning("default").
		DeleteAllExisting().
		AddDataSource(PrometheusDataSource("prometheus", "http://prometheus:9090"))

	data, err := provisioning.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	if !strings.Contains(string(data), "deleteDatasources:") {
		t.Error("Expected deleteDatasources section")
	}
}

func TestDashboardProvisioning_Serialize(t *testing.T) {
	provisioning := NewDashboardProvisioning("default").
		WithFolder("Monitoring").
		WithPath("/var/lib/grafana/dashboards")

	data, err := provisioning.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "apiVersion: 1") {
		t.Error("Expected apiVersion: 1")
	}
	if !strings.Contains(yamlStr, "folder: Monitoring") {
		t.Error("Expected folder: Monitoring")
	}
	if !strings.Contains(yamlStr, "path: /var/lib/grafana/dashboards") {
		t.Error("Expected path")
	}
}

func TestDashboardProvisioning_WithOptions(t *testing.T) {
	provisioning := NewDashboardProvisioning("default").
		WithFolder("Services").
		WithPath("/dashboards").
		Editable().
		WithUpdateInterval(30)

	data, err := provisioning.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	if !strings.Contains(yamlStr, "editable: true") {
		t.Error("Expected editable: true")
	}
	if !strings.Contains(yamlStr, "updateIntervalSeconds: 30") {
		t.Error("Expected updateIntervalSeconds: 30")
	}
}

func TestDashboardProvisioning_DisableDelete(t *testing.T) {
	provisioning := NewDashboardProvisioning("default").
		WithPath("/dashboards").
		DisableDelete()

	data, err := provisioning.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	if !strings.Contains(string(data), "disableDeletion: true") {
		t.Error("Expected disableDeletion: true")
	}
}
