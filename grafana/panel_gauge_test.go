package grafana

import "testing"

func TestGauge(t *testing.T) {
	p := Gauge("CPU Usage")
	if p.Title != "CPU Usage" {
		t.Errorf("Title = %q, want CPU Usage", p.Title)
	}
	if p.Type != "gauge" {
		t.Errorf("Type = %q, want gauge", p.Type)
	}
}

func TestGauge_WithDescription(t *testing.T) {
	p := Gauge("Test").WithDescription("A test gauge")
	if p.Description != "A test gauge" {
		t.Errorf("Description = %q, want A test gauge", p.Description)
	}
}

func TestGauge_WithSize(t *testing.T) {
	p := Gauge("Test").WithSize(8, 6)
	if p.GridPos.W != 8 || p.GridPos.H != 6 {
		t.Errorf("GridPos = %+v, want W=8 H=6", p.GridPos)
	}
}

func TestGauge_WithMinMax(t *testing.T) {
	p := Gauge("Test").WithMin(0).WithMax(100)
	if *p.FieldConfig.Defaults.Min != 0 {
		t.Errorf("Min = %v, want 0", *p.FieldConfig.Defaults.Min)
	}
	if *p.FieldConfig.Defaults.Max != 100 {
		t.Errorf("Max = %v, want 100", *p.FieldConfig.Defaults.Max)
	}
}

func TestGauge_ShowThresholdLabels(t *testing.T) {
	p := Gauge("Test").ShowThresholdLabels()
	if !p.Options.ShowThresholdLabels {
		t.Error("ShowThresholdLabels should be true")
	}
}

func TestGauge_ShowThresholdMarkers(t *testing.T) {
	p := Gauge("Test").ShowThresholdMarkers()
	if !p.Options.ShowThresholdMarkers {
		t.Error("ShowThresholdMarkers should be true")
	}
}

func TestGauge_FluentAPI(t *testing.T) {
	p := Gauge("Test").
		WithDescription("desc").
		WithSize(12, 8).
		WithMin(0).
		WithMax(100).
		ShowThresholdLabels().
		ShowThresholdMarkers()

	if p.Title != "Test" {
		t.Error("Fluent API should preserve title")
	}
	if p.Options.ShowThresholdLabels != true {
		t.Error("Fluent API should set threshold labels")
	}
}
