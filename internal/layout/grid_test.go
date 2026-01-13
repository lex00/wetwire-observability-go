package layout

import (
	"testing"

	"github.com/lex00/wetwire-observability-go/grafana"
)

func TestCalculateGridPositions_SinglePanel(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Row 1").WithPanels(
				grafana.TimeSeries("Panel 1"),
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}

	panel := dashboard.Rows[0].Panels[0].(*grafana.TimeSeriesPanel)
	if panel.GridPos.X != 0 {
		t.Errorf("Panel X = %d, want 0", panel.GridPos.X)
	}
	if panel.GridPos.Y != 1 { // After row header
		t.Errorf("Panel Y = %d, want 1", panel.GridPos.Y)
	}
}

func TestCalculateGridPositions_TwoPanelsSideBySide(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Row 1").WithPanels(
				grafana.TimeSeries("Panel 1").WithSize(12, 8),
				grafana.TimeSeries("Panel 2").WithSize(12, 8),
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}

	panel1 := dashboard.Rows[0].Panels[0].(*grafana.TimeSeriesPanel)
	panel2 := dashboard.Rows[0].Panels[1].(*grafana.TimeSeriesPanel)

	// Panel 1 should be at x=0
	if panel1.GridPos.X != 0 {
		t.Errorf("Panel 1 X = %d, want 0", panel1.GridPos.X)
	}

	// Panel 2 should be at x=12 (after panel 1)
	if panel2.GridPos.X != 12 {
		t.Errorf("Panel 2 X = %d, want 12", panel2.GridPos.X)
	}

	// Both should be on same Y
	if panel1.GridPos.Y != panel2.GridPos.Y {
		t.Errorf("Panels should be on same Y: %d vs %d", panel1.GridPos.Y, panel2.GridPos.Y)
	}
}

func TestCalculateGridPositions_Wrapping(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Row 1").WithPanels(
				grafana.TimeSeries("Panel 1").WithSize(12, 8),
				grafana.TimeSeries("Panel 2").WithSize(12, 8),
				grafana.TimeSeries("Panel 3").WithSize(12, 8), // Should wrap
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}

	panel3 := dashboard.Rows[0].Panels[2].(*grafana.TimeSeriesPanel)

	// Panel 3 should wrap to x=0
	if panel3.GridPos.X != 0 {
		t.Errorf("Panel 3 X = %d, want 0", panel3.GridPos.X)
	}
}

func TestCalculateGridPositions_MultipleRows(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Row 1").WithPanels(
				grafana.TimeSeries("Panel 1").WithSize(24, 8),
			),
			grafana.NewRow("Row 2").WithPanels(
				grafana.TimeSeries("Panel 2").WithSize(24, 8),
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}

	panel1 := dashboard.Rows[0].Panels[0].(*grafana.TimeSeriesPanel)
	panel2 := dashboard.Rows[1].Panels[0].(*grafana.TimeSeriesPanel)

	// Row 2 should be below Row 1
	if panel2.GridPos.Y <= panel1.GridPos.Y+panel1.GridPos.H {
		t.Errorf("Panel 2 Y (%d) should be below Panel 1 (Y=%d, H=%d)",
			panel2.GridPos.Y, panel1.GridPos.Y, panel1.GridPos.H)
	}
}

func TestCalculateGridPositions_MixedPanelSizes(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Row 1").WithPanels(
				grafana.TimeSeries("Full Width").WithSize(24, 8),
			),
			grafana.NewRow("Row 2").WithPanels(
				grafana.Stat("Stat 1").WithSize(6, 4),
				grafana.Stat("Stat 2").WithSize(6, 4),
				grafana.Stat("Stat 3").WithSize(6, 4),
				grafana.Stat("Stat 4").WithSize(6, 4),
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}

	stat1 := dashboard.Rows[1].Panels[0].(*grafana.StatPanel)
	stat2 := dashboard.Rows[1].Panels[1].(*grafana.StatPanel)
	stat3 := dashboard.Rows[1].Panels[2].(*grafana.StatPanel)
	stat4 := dashboard.Rows[1].Panels[3].(*grafana.StatPanel)

	// All stats should fit on same row
	if stat1.GridPos.X != 0 {
		t.Errorf("Stat 1 X = %d, want 0", stat1.GridPos.X)
	}
	if stat2.GridPos.X != 6 {
		t.Errorf("Stat 2 X = %d, want 6", stat2.GridPos.X)
	}
	if stat3.GridPos.X != 12 {
		t.Errorf("Stat 3 X = %d, want 12", stat3.GridPos.X)
	}
	if stat4.GridPos.X != 18 {
		t.Errorf("Stat 4 X = %d, want 18", stat4.GridPos.X)
	}
}

func TestCalculateGridPositions_CollapsedRow(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Row 1").Collapsed().WithPanels(
				grafana.TimeSeries("Panel 1").WithSize(24, 8),
			),
			grafana.NewRow("Row 2").WithPanels(
				grafana.TimeSeries("Panel 2").WithSize(24, 8),
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}

	// Row 2 should account for collapsed Row 1
	panel2 := dashboard.Rows[1].Panels[0].(*grafana.TimeSeriesPanel)
	// Collapsed row takes 1 row height
	if panel2.GridPos.Y < 2 {
		t.Errorf("Panel 2 Y = %d should be >= 2 (after collapsed row)", panel2.GridPos.Y)
	}
}

func TestCalculateGridPositions_EmptyDashboard(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test")

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}
}

func TestCalculateGridPositions_EmptyRow(t *testing.T) {
	dashboard := grafana.NewDashboard("test", "Test").
		WithRows(
			grafana.NewRow("Empty Row"),
			grafana.NewRow("Row 2").WithPanels(
				grafana.TimeSeries("Panel 1"),
			),
		)

	err := CalculateGridPositions(dashboard)
	if err != nil {
		t.Fatalf("CalculateGridPositions() error = %v", err)
	}
}

func TestFullWidth(t *testing.T) {
	if FullWidth() != 24 {
		t.Errorf("FullWidth() = %d, want 24", FullWidth())
	}
}

func TestHalfWidth(t *testing.T) {
	if HalfWidth() != 12 {
		t.Errorf("HalfWidth() = %d, want 12", HalfWidth())
	}
}

func TestThirdWidth(t *testing.T) {
	if ThirdWidth() != 8 {
		t.Errorf("ThirdWidth() = %d, want 8", ThirdWidth())
	}
}

func TestQuarterWidth(t *testing.T) {
	if QuarterWidth() != 6 {
		t.Errorf("QuarterWidth() = %d, want 6", QuarterWidth())
	}
}
