package grafana

import (
	"testing"
)

func TestNewRow(t *testing.T) {
	row := NewRow("Test Row")
	if row == nil {
		t.Fatal("NewRow() returned nil")
	}
	if row.Title != "Test Row" {
		t.Errorf("Title = %v, want Test Row", row.Title)
	}
}

func TestRow_WithPanels(t *testing.T) {
	// Use any for panels since we haven't defined panel types yet
	row := NewRow("Test").WithPanels("panel1", "panel2")
	if len(row.Panels) != 2 {
		t.Errorf("len(Panels) = %d, want 2", len(row.Panels))
	}
}

func TestRow_AddPanel(t *testing.T) {
	row := NewRow("Test")
	row.AddPanel("panel1")
	row.AddPanel("panel2")
	if len(row.Panels) != 2 {
		t.Errorf("len(Panels) = %d, want 2", len(row.Panels))
	}
}

func TestRow_Collapsed(t *testing.T) {
	row := NewRow("Test").Collapsed()
	if !row.IsCollapsed {
		t.Error("IsCollapsed should be true")
	}
}

func TestRow_Expanded(t *testing.T) {
	row := NewRow("Test").Collapsed().Expanded()
	if row.IsCollapsed {
		t.Error("IsCollapsed should be false")
	}
}

func TestRow_WithHeight(t *testing.T) {
	row := NewRow("Test").WithHeight(300)
	if row.Height != 300 {
		t.Errorf("Height = %d, want 300", row.Height)
	}
}

func TestRow_FluentAPI(t *testing.T) {
	row := NewRow("Metrics Overview").
		WithHeight(250).
		Collapsed().
		WithPanels("panel1", "panel2", "panel3")

	if row.Title != "Metrics Overview" {
		t.Errorf("Title = %v", row.Title)
	}
	if row.Height != 250 {
		t.Errorf("Height = %d", row.Height)
	}
	if !row.IsCollapsed {
		t.Error("IsCollapsed should be true")
	}
	if len(row.Panels) != 3 {
		t.Errorf("len(Panels) = %d", len(row.Panels))
	}
}
