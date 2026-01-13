package grafana

import "testing"

func TestPieChart(t *testing.T) {
	p := PieChart("Request Distribution")
	if p.Title != "Request Distribution" {
		t.Errorf("Title = %q, want Request Distribution", p.Title)
	}
	if p.Type != "piechart" {
		t.Errorf("Type = %q, want piechart", p.Type)
	}
}

func TestPieChart_WithDescription(t *testing.T) {
	p := PieChart("Test").WithDescription("A test pie chart")
	if p.Description != "A test pie chart" {
		t.Errorf("Description = %q, want A test pie chart", p.Description)
	}
}

func TestPieChart_WithSize(t *testing.T) {
	p := PieChart("Test").WithSize(10, 10)
	if p.GridPos.W != 10 || p.GridPos.H != 10 {
		t.Errorf("GridPos = %+v, want W=10 H=10", p.GridPos)
	}
}

func TestPieChart_PieType(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*PieChartPanel) *PieChartPanel
		expected string
	}{
		{"Pie", func(p *PieChartPanel) *PieChartPanel { return p.Pie() }, "pie"},
		{"Donut", func(p *PieChartPanel) *PieChartPanel { return p.Donut() }, "donut"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setup(PieChart("Test"))
			if p.Options.PieType != tt.expected {
				t.Errorf("PieType = %q, want %q", p.Options.PieType, tt.expected)
			}
		})
	}
}

func TestPieChart_LegendPlacement(t *testing.T) {
	p := PieChart("Test").LegendRight()
	if p.Options.Legend.Placement != LegendRight {
		t.Errorf("Legend.Placement = %q, want right", p.Options.Legend.Placement)
	}

	p = PieChart("Test").LegendBottom()
	if p.Options.Legend.Placement != LegendBottom {
		t.Errorf("Legend.Placement = %q, want bottom", p.Options.Legend.Placement)
	}
}

func TestPieChart_HideLegend(t *testing.T) {
	p := PieChart("Test").HideLegend()
	if p.Options.Legend.DisplayMode != "hidden" {
		t.Errorf("Legend.DisplayMode = %q, want hidden", p.Options.Legend.DisplayMode)
	}
}

func TestPieChart_ShowLabels(t *testing.T) {
	p := PieChart("Test").ShowLabels()
	if !p.Options.Labels {
		t.Error("Labels should be true")
	}
}

func TestPieChart_ShowTooltip(t *testing.T) {
	p := PieChart("Test").ShowTooltip()
	if p.Options.Tooltip.Mode != "single" {
		t.Errorf("Tooltip.Mode = %q, want single", p.Options.Tooltip.Mode)
	}
}

func TestPieChart_FluentAPI(t *testing.T) {
	p := PieChart("Test").
		WithDescription("desc").
		WithSize(12, 10).
		Donut().
		LegendRight().
		ShowLabels()

	if p.Title != "Test" {
		t.Error("Fluent API should preserve title")
	}
	if p.Options.PieType != "donut" {
		t.Error("Fluent API should set pie type")
	}
}
