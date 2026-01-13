package grafana

import (
	"testing"
)

func TestTimeSeries(t *testing.T) {
	p := TimeSeries("CPU Usage")
	if p == nil {
		t.Fatal("TimeSeries() returned nil")
	}
	if p.Title != "CPU Usage" {
		t.Errorf("Title = %v, want CPU Usage", p.Title)
	}
	if p.Type != "timeseries" {
		t.Errorf("Type = %v, want timeseries", p.Type)
	}
}

func TestTimeSeries_WithDescription(t *testing.T) {
	p := TimeSeries("CPU Usage").WithDescription("CPU usage over time")
	if p.Description != "CPU usage over time" {
		t.Errorf("Description = %v", p.Description)
	}
}

func TestTimeSeries_WithTargets(t *testing.T) {
	p := TimeSeries("CPU Usage").WithTargets("target1", "target2")
	if len(p.Targets) != 2 {
		t.Errorf("len(Targets) = %d, want 2", len(p.Targets))
	}
}

func TestTimeSeries_AddTarget(t *testing.T) {
	p := TimeSeries("CPU Usage")
	p.AddTarget("target1")
	p.AddTarget("target2")
	if len(p.Targets) != 2 {
		t.Errorf("len(Targets) = %d, want 2", len(p.Targets))
	}
}

func TestTimeSeries_WithDatasource(t *testing.T) {
	p := TimeSeries("CPU Usage").WithDatasource("$datasource")
	if p.Datasource != "$datasource" {
		t.Errorf("Datasource = %v, want $datasource", p.Datasource)
	}
}

func TestTimeSeries_WithSize(t *testing.T) {
	p := TimeSeries("CPU Usage").WithSize(12, 8)
	if p.GridPos.W != 12 {
		t.Errorf("Width = %d, want 12", p.GridPos.W)
	}
	if p.GridPos.H != 8 {
		t.Errorf("Height = %d, want 8", p.GridPos.H)
	}
}

func TestTimeSeries_WithPosition(t *testing.T) {
	p := TimeSeries("CPU Usage").WithPosition(6, 0)
	if p.GridPos.X != 6 {
		t.Errorf("X = %d, want 6", p.GridPos.X)
	}
	if p.GridPos.Y != 0 {
		t.Errorf("Y = %d, want 0", p.GridPos.Y)
	}
}

func TestTimeSeries_LegendPosition(t *testing.T) {
	p := TimeSeries("CPU Usage").WithLegendPosition(LegendBottom)
	if p.Options.Legend.DisplayMode != "list" {
		t.Errorf("Legend.DisplayMode = %v", p.Options.Legend.DisplayMode)
	}
	if p.Options.Legend.Placement != LegendBottom {
		t.Errorf("Legend.Placement = %v, want %v", p.Options.Legend.Placement, LegendBottom)
	}
}

func TestTimeSeries_HideLegend(t *testing.T) {
	p := TimeSeries("CPU Usage").HideLegend()
	if p.Options.Legend.DisplayMode != "hidden" {
		t.Errorf("Legend.DisplayMode = %v, want hidden", p.Options.Legend.DisplayMode)
	}
}

func TestTimeSeries_WithTooltip(t *testing.T) {
	p := TimeSeries("CPU Usage").WithTooltip(TooltipAll)
	if p.Options.Tooltip.Mode != TooltipAll {
		t.Errorf("Tooltip.Mode = %v, want %v", p.Options.Tooltip.Mode, TooltipAll)
	}
}

func TestTimeSeries_LineWidth(t *testing.T) {
	p := TimeSeries("CPU Usage").WithLineWidth(2)
	if p.FieldConfig.Defaults.Custom.LineWidth != 2 {
		t.Errorf("LineWidth = %d, want 2", p.FieldConfig.Defaults.Custom.LineWidth)
	}
}

func TestTimeSeries_FillOpacity(t *testing.T) {
	p := TimeSeries("CPU Usage").WithFillOpacity(20)
	if p.FieldConfig.Defaults.Custom.FillOpacity != 20 {
		t.Errorf("FillOpacity = %d, want 20", p.FieldConfig.Defaults.Custom.FillOpacity)
	}
}

func TestTimeSeries_DrawBars(t *testing.T) {
	p := TimeSeries("CPU Usage").DrawBars()
	if p.FieldConfig.Defaults.Custom.DrawStyle != DrawStyleBars {
		t.Errorf("DrawStyle = %v, want %v", p.FieldConfig.Defaults.Custom.DrawStyle, DrawStyleBars)
	}
}

func TestTimeSeries_DrawPoints(t *testing.T) {
	p := TimeSeries("CPU Usage").DrawPoints()
	if p.FieldConfig.Defaults.Custom.DrawStyle != DrawStylePoints {
		t.Errorf("DrawStyle = %v, want %v", p.FieldConfig.Defaults.Custom.DrawStyle, DrawStylePoints)
	}
}

func TestTimeSeries_WithUnit(t *testing.T) {
	p := TimeSeries("CPU Usage").WithUnit(UnitPercent)
	if p.FieldConfig.Defaults.Unit != UnitPercent {
		t.Errorf("Unit = %v, want %v", p.FieldConfig.Defaults.Unit, UnitPercent)
	}
}

func TestTimeSeries_FluentAPI(t *testing.T) {
	p := TimeSeries("CPU Usage").
		WithDescription("CPU usage percentage").
		WithDatasource("Prometheus").
		WithSize(12, 8).
		WithPosition(0, 0).
		WithLegendPosition(LegendBottom).
		WithTooltip(TooltipAll).
		WithLineWidth(2).
		WithFillOpacity(10).
		WithUnit(UnitPercent)

	if p.Title != "CPU Usage" {
		t.Errorf("Title = %v", p.Title)
	}
	if p.GridPos.W != 12 || p.GridPos.H != 8 {
		t.Errorf("GridPos = %+v", p.GridPos)
	}
}

func TestStat(t *testing.T) {
	p := Stat("Total Requests")
	if p == nil {
		t.Fatal("Stat() returned nil")
	}
	if p.Title != "Total Requests" {
		t.Errorf("Title = %v, want Total Requests", p.Title)
	}
	if p.Type != "stat" {
		t.Errorf("Type = %v, want stat", p.Type)
	}
}

func TestStat_WithDescription(t *testing.T) {
	p := Stat("Total Requests").WithDescription("Total HTTP requests")
	if p.Description != "Total HTTP requests" {
		t.Errorf("Description = %v", p.Description)
	}
}

func TestStat_ColorByValue(t *testing.T) {
	p := Stat("Total Requests").ColorByValue()
	if p.Options.ColorMode != ColorModeValue {
		t.Errorf("ColorMode = %v, want %v", p.Options.ColorMode, ColorModeValue)
	}
}

func TestStat_ColorByBackground(t *testing.T) {
	p := Stat("Total Requests").ColorByBackground()
	if p.Options.ColorMode != ColorModeBackground {
		t.Errorf("ColorMode = %v, want %v", p.Options.ColorMode, ColorModeBackground)
	}
}

func TestStat_WithReduceCalc(t *testing.T) {
	p := Stat("Total Requests").WithReduceCalc(ReduceMax)
	if p.Options.ReduceOptions.Calcs[0] != ReduceMax {
		t.Errorf("Calcs[0] = %v, want %v", p.Options.ReduceOptions.Calcs[0], ReduceMax)
	}
}

func TestStat_GraphMode(t *testing.T) {
	p := Stat("Total Requests").WithGraphMode(GraphModeArea)
	if p.Options.GraphMode != GraphModeArea {
		t.Errorf("GraphMode = %v, want %v", p.Options.GraphMode, GraphModeArea)
	}
}

func TestStat_NoGraph(t *testing.T) {
	p := Stat("Total Requests").NoGraph()
	if p.Options.GraphMode != GraphModeNone {
		t.Errorf("GraphMode = %v, want %v", p.Options.GraphMode, GraphModeNone)
	}
}

func TestStat_Orientation(t *testing.T) {
	p := Stat("Total Requests").WithOrientation(OrientationHorizontal)
	if p.Options.Orientation != OrientationHorizontal {
		t.Errorf("Orientation = %v, want %v", p.Options.Orientation, OrientationHorizontal)
	}
}

func TestStat_FluentAPI(t *testing.T) {
	p := Stat("Error Rate").
		WithDescription("HTTP error rate").
		WithDatasource("Prometheus").
		WithSize(6, 4).
		ColorByBackground().
		WithReduceCalc(ReduceLast).
		WithGraphMode(GraphModeArea).
		WithUnit(UnitPercent)

	if p.Title != "Error Rate" {
		t.Errorf("Title = %v", p.Title)
	}
	if p.Options.ColorMode != ColorModeBackground {
		t.Errorf("ColorMode = %v", p.Options.ColorMode)
	}
}
