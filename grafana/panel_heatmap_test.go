package grafana

import "testing"

func TestHeatmap(t *testing.T) {
	p := Heatmap("Latency Distribution")
	if p.Title != "Latency Distribution" {
		t.Errorf("Title = %q, want Latency Distribution", p.Title)
	}
	if p.Type != "heatmap" {
		t.Errorf("Type = %q, want heatmap", p.Type)
	}
}

func TestHeatmap_WithDescription(t *testing.T) {
	p := Heatmap("Test").WithDescription("A test heatmap")
	if p.Description != "A test heatmap" {
		t.Errorf("Description = %q, want A test heatmap", p.Description)
	}
}

func TestHeatmap_WithSize(t *testing.T) {
	p := Heatmap("Test").WithSize(24, 10)
	if p.GridPos.W != 24 || p.GridPos.H != 10 {
		t.Errorf("GridPos = %+v, want W=24 H=10", p.GridPos)
	}
}

func TestHeatmap_ColorScheme(t *testing.T) {
	p := Heatmap("Test").WithColorScheme("Spectral")
	if p.Options.Color.Scheme != "Spectral" {
		t.Errorf("Color.Scheme = %q, want Spectral", p.Options.Color.Scheme)
	}
}

func TestHeatmap_PresetColorSchemes(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*HeatmapPanel) *HeatmapPanel
		expected string
	}{
		{"Blues", func(p *HeatmapPanel) *HeatmapPanel { return p.Blues() }, "Blues"},
		{"Reds", func(p *HeatmapPanel) *HeatmapPanel { return p.Reds() }, "Reds"},
		{"Greens", func(p *HeatmapPanel) *HeatmapPanel { return p.Greens() }, "Greens"},
		{"Spectral", func(p *HeatmapPanel) *HeatmapPanel { return p.Spectral() }, "Spectral"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setup(Heatmap("Test"))
			if p.Options.Color.Scheme != tt.expected {
				t.Errorf("Color.Scheme = %q, want %q", p.Options.Color.Scheme, tt.expected)
			}
		})
	}
}

func TestHeatmap_ShowLegend(t *testing.T) {
	p := Heatmap("Test").ShowLegend()
	if !p.Options.Legend.Show {
		t.Error("Legend.Show should be true")
	}
}

func TestHeatmap_HideLegend(t *testing.T) {
	p := Heatmap("Test").ShowLegend().HideLegend()
	if p.Options.Legend.Show {
		t.Error("Legend.Show should be false")
	}
}

func TestHeatmap_ShowTooltip(t *testing.T) {
	p := Heatmap("Test").ShowTooltip()
	if !p.Options.Tooltip.Show {
		t.Error("Tooltip.Show should be true")
	}
}

func TestHeatmap_ShowYHistogram(t *testing.T) {
	p := Heatmap("Test").ShowYHistogram()
	if !p.Options.Tooltip.YHistogram {
		t.Error("YHistogram should be true")
	}
}

func TestHeatmap_Calculate(t *testing.T) {
	p := Heatmap("Test").Calculate()
	if !p.Options.Calculate {
		t.Error("Calculate should be true")
	}
}

func TestHeatmap_WithBucketSize(t *testing.T) {
	p := Heatmap("Test").WithXBucketSize(10).WithYBucketSize(5)
	if p.Options.Calculation.XBuckets.Value != "10" {
		t.Errorf("XBuckets.Value = %q, want 10", p.Options.Calculation.XBuckets.Value)
	}
	if p.Options.Calculation.YBuckets.Value != "5" {
		t.Errorf("YBuckets.Value = %q, want 5", p.Options.Calculation.YBuckets.Value)
	}
}

func TestHeatmap_FluentAPI(t *testing.T) {
	p := Heatmap("Test").
		WithDescription("desc").
		WithSize(24, 12).
		Spectral().
		ShowLegend().
		ShowTooltip().
		Calculate()

	if p.Title != "Test" {
		t.Error("Fluent API should preserve title")
	}
	if p.Options.Color.Scheme != "Spectral" {
		t.Error("Fluent API should set color scheme")
	}
	if !p.Options.Legend.Show {
		t.Error("Fluent API should show legend")
	}
}
