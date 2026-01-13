package grafana

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDashboard_Serialize(t *testing.T) {
	d := NewDashboard("test-uid", "Test Dashboard").
		WithDescription("A test dashboard").
		WithTags("test", "monitoring").
		WithTime("now-1h", "now").
		WithRefresh("30s")

	data, err := d.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Verify it's valid JSON
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Verify key fields
	if result["uid"] != "test-uid" {
		t.Errorf("uid = %v, want test-uid", result["uid"])
	}
	if result["title"] != "Test Dashboard" {
		t.Errorf("title = %v, want Test Dashboard", result["title"])
	}
}

func TestDashboard_Serialize_WithPanels(t *testing.T) {
	d := NewDashboard("api-overview", "API Overview").
		WithRows(
			NewRow("Overview").WithPanels(
				TimeSeries("Request Rate").WithSize(12, 8),
				TimeSeries("Error Rate").WithSize(12, 8),
			),
			NewRow("Details").WithPanels(
				Table("Recent Requests").WithSize(24, 10),
			),
		)

	data, err := d.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Verify it's valid JSON
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Check that panels are included
	panels, ok := result["panels"].([]any)
	if !ok {
		t.Fatal("Expected panels array in output")
	}
	if len(panels) == 0 {
		t.Error("Expected panels in output")
	}
}

func TestDashboard_Serialize_PanelIDs(t *testing.T) {
	d := NewDashboard("test", "Test").
		WithRows(
			NewRow("Row 1").WithPanels(
				TimeSeries("Panel 1"),
				TimeSeries("Panel 2"),
			),
		)

	data, err := d.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	panels := result["panels"].([]any)
	ids := make(map[float64]bool)

	for _, p := range panels {
		panel := p.(map[string]any)
		id, ok := panel["id"].(float64)
		if !ok {
			t.Error("Panel missing id field")
			continue
		}
		if ids[id] {
			t.Errorf("Duplicate panel ID: %v", id)
		}
		ids[id] = true
	}
}

func TestDashboard_SerializeToFile(t *testing.T) {
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "dashboard.json")

	d := NewDashboard("test", "Test Dashboard")

	err := d.SerializeToFile(path)
	if err != nil {
		t.Fatalf("SerializeToFile() error = %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read file: %v", err)
	}

	// Verify it's valid JSON
	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}
}

func TestDashboard_MustSerialize(t *testing.T) {
	d := NewDashboard("test", "Test")

	// Should not panic
	data := d.MustSerialize()
	if len(data) == 0 {
		t.Error("MustSerialize() returned empty data")
	}
}

func TestDashboard_Serialize_Formatted(t *testing.T) {
	d := NewDashboard("test", "Test")

	data, err := d.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Verify it's pretty-printed (has newlines)
	if !strings.Contains(string(data), "\n") {
		t.Error("Expected formatted JSON with newlines")
	}
}

func TestDashboard_Serialize_WithVariables(t *testing.T) {
	d := NewDashboard("test", "Test").
		WithVariables(
			QueryVar("namespace", "label_values(kube_pod_info, namespace)"),
			CustomVar("env", "dev", "staging", "prod"),
		)

	data, err := d.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("Output is not valid JSON: %v", err)
	}

	// Check templating section
	templating, ok := result["templating"].(map[string]any)
	if !ok {
		t.Fatal("Expected templating object in output")
	}

	list, ok := templating["list"].([]any)
	if !ok {
		t.Fatal("Expected templating.list array in output")
	}

	if len(list) != 2 {
		t.Errorf("Expected 2 variables, got %d", len(list))
	}
}

func TestAssignPanelIDs(t *testing.T) {
	d := NewDashboard("test", "Test").
		WithRows(
			NewRow("Row 1").WithPanels(
				TimeSeries("Panel 1"),
				TimeSeries("Panel 2"),
			),
			NewRow("Row 2").WithPanels(
				Stat("Stat 1"),
			),
		)

	AssignPanelIDs(d)

	panel1 := d.Rows[0].Panels[0].(*TimeSeriesPanel)
	panel2 := d.Rows[0].Panels[1].(*TimeSeriesPanel)
	panel3 := d.Rows[1].Panels[0].(*StatPanel)

	// IDs should be sequential and unique
	if panel1.ID == 0 {
		t.Error("Panel 1 ID should not be 0")
	}
	if panel2.ID == 0 {
		t.Error("Panel 2 ID should not be 0")
	}
	if panel3.ID == 0 {
		t.Error("Panel 3 ID should not be 0")
	}
	if panel1.ID == panel2.ID || panel2.ID == panel3.ID || panel1.ID == panel3.ID {
		t.Error("Panel IDs should be unique")
	}
}

func TestAssignPanelIDs_Deterministic(t *testing.T) {
	d1 := NewDashboard("test", "Test").
		WithRows(
			NewRow("Row").WithPanels(
				TimeSeries("A"),
				TimeSeries("B"),
			),
		)
	d2 := NewDashboard("test", "Test").
		WithRows(
			NewRow("Row").WithPanels(
				TimeSeries("A"),
				TimeSeries("B"),
			),
		)

	AssignPanelIDs(d1)
	AssignPanelIDs(d2)

	p1a := d1.Rows[0].Panels[0].(*TimeSeriesPanel)
	p2a := d2.Rows[0].Panels[0].(*TimeSeriesPanel)

	if p1a.ID != p2a.ID {
		t.Errorf("Panel IDs should be deterministic: %d != %d", p1a.ID, p2a.ID)
	}
}
