package grafana

import "testing"

func TestBarGauge(t *testing.T) {
	p := BarGauge("Request Rate")
	if p.Title != "Request Rate" {
		t.Errorf("Title = %q, want Request Rate", p.Title)
	}
	if p.Type != "bargauge" {
		t.Errorf("Type = %q, want bargauge", p.Type)
	}
}

func TestBarGauge_WithDescription(t *testing.T) {
	p := BarGauge("Test").WithDescription("A test bar gauge")
	if p.Description != "A test bar gauge" {
		t.Errorf("Description = %q, want A test bar gauge", p.Description)
	}
}

func TestBarGauge_WithSize(t *testing.T) {
	p := BarGauge("Test").WithSize(8, 6)
	if p.GridPos.W != 8 || p.GridPos.H != 6 {
		t.Errorf("GridPos = %+v, want W=8 H=6", p.GridPos)
	}
}

func TestBarGauge_Horizontal(t *testing.T) {
	p := BarGauge("Test").Horizontal()
	if p.Options.Orientation != OrientationHorizontal {
		t.Errorf("Orientation = %q, want horizontal", p.Options.Orientation)
	}
}

func TestBarGauge_Vertical(t *testing.T) {
	p := BarGauge("Test").Vertical()
	if p.Options.Orientation != OrientationVertical {
		t.Errorf("Orientation = %q, want vertical", p.Options.Orientation)
	}
}

func TestBarGauge_DisplayMode(t *testing.T) {
	tests := []struct {
		name     string
		setup    func(*BarGaugePanel) *BarGaugePanel
		expected string
	}{
		{"Basic", func(p *BarGaugePanel) *BarGaugePanel { return p.Basic() }, "basic"},
		{"Gradient", func(p *BarGaugePanel) *BarGaugePanel { return p.Gradient() }, "gradient"},
		{"LCD", func(p *BarGaugePanel) *BarGaugePanel { return p.LCD() }, "lcd"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := tt.setup(BarGauge("Test"))
			if p.Options.DisplayMode != tt.expected {
				t.Errorf("DisplayMode = %q, want %q", p.Options.DisplayMode, tt.expected)
			}
		})
	}
}

func TestBarGauge_WithMinMax(t *testing.T) {
	p := BarGauge("Test").WithMin(0).WithMax(100)
	if *p.FieldConfig.Defaults.Min != 0 {
		t.Errorf("Min = %v, want 0", *p.FieldConfig.Defaults.Min)
	}
	if *p.FieldConfig.Defaults.Max != 100 {
		t.Errorf("Max = %v, want 100", *p.FieldConfig.Defaults.Max)
	}
}

func TestBarGauge_FluentAPI(t *testing.T) {
	p := BarGauge("Test").
		WithDescription("desc").
		WithSize(12, 8).
		Horizontal().
		Gradient().
		WithMin(0).
		WithMax(100)

	if p.Title != "Test" {
		t.Error("Fluent API should preserve title")
	}
	if p.Options.DisplayMode != "gradient" {
		t.Error("Fluent API should set display mode")
	}
}
