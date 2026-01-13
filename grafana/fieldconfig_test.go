package grafana

import "testing"

func TestNewFieldConfig(t *testing.T) {
	fc := NewFieldConfig()
	if fc == nil {
		t.Fatal("NewFieldConfig() returned nil")
	}
}

func TestFieldConfig_WithUnit(t *testing.T) {
	fc := NewFieldConfig().WithUnit(UnitPercent)
	if fc.Defaults.Unit != UnitPercent {
		t.Errorf("Unit = %q, want percent", fc.Defaults.Unit)
	}
}

func TestFieldConfig_WithDecimals(t *testing.T) {
	fc := NewFieldConfig().WithDecimals(2)
	if fc.Defaults.Decimals == nil || *fc.Defaults.Decimals != 2 {
		t.Errorf("Decimals = %v, want 2", fc.Defaults.Decimals)
	}
}

func TestFieldConfig_WithMin(t *testing.T) {
	fc := NewFieldConfig().WithMin(0)
	if fc.Defaults.Min == nil || *fc.Defaults.Min != 0 {
		t.Errorf("Min = %v, want 0", fc.Defaults.Min)
	}
}

func TestFieldConfig_WithMax(t *testing.T) {
	fc := NewFieldConfig().WithMax(100)
	if fc.Defaults.Max == nil || *fc.Defaults.Max != 100 {
		t.Errorf("Max = %v, want 100", fc.Defaults.Max)
	}
}

func TestFieldConfig_WithNoValue(t *testing.T) {
	fc := NewFieldConfig().WithNoValue("N/A")
	if fc.Defaults.NoValue != "N/A" {
		t.Errorf("NoValue = %q, want N/A", fc.Defaults.NoValue)
	}
}

func TestFieldConfig_WithThresholds(t *testing.T) {
	thresholds := GreenYellowRed(50, 80)
	fc := NewFieldConfig().WithThresholds(thresholds)
	if fc.Defaults.Thresholds == nil {
		t.Fatal("Thresholds is nil")
	}
	if len(fc.Defaults.Thresholds.Steps) != 3 {
		t.Errorf("len(Thresholds.Steps) = %d, want 3", len(fc.Defaults.Thresholds.Steps))
	}
}

func TestFieldConfig_WithColor(t *testing.T) {
	fc := NewFieldConfig().WithColor(ColorByValue())
	if fc.Defaults.Color == nil {
		t.Fatal("Color is nil")
	}
	if fc.Defaults.Color.Mode != "thresholds" {
		t.Errorf("Color.Mode = %q, want thresholds", fc.Defaults.Color.Mode)
	}
}

func TestFieldConfig_WithFixedColor(t *testing.T) {
	fc := NewFieldConfig().WithFixedColor("blue")
	if fc.Defaults.Color == nil {
		t.Fatal("Color is nil")
	}
	if fc.Defaults.Color.Mode != "fixed" {
		t.Errorf("Color.Mode = %q, want fixed", fc.Defaults.Color.Mode)
	}
	if fc.Defaults.Color.FixedColor != "blue" {
		t.Errorf("Color.FixedColor = %q, want blue", fc.Defaults.Color.FixedColor)
	}
}

func TestFieldConfig_FluentAPI(t *testing.T) {
	fc := NewFieldConfig().
		WithUnit(UnitPercent).
		WithDecimals(1).
		WithMin(0).
		WithMax(100).
		WithThresholds(GreenYellowRed(50, 80))

	if fc.Defaults.Unit != UnitPercent {
		t.Error("Fluent API should set unit")
	}
	if *fc.Defaults.Decimals != 1 {
		t.Error("Fluent API should set decimals")
	}
	if *fc.Defaults.Min != 0 {
		t.Error("Fluent API should set min")
	}
	if *fc.Defaults.Max != 100 {
		t.Error("Fluent API should set max")
	}
}

func TestColorByValue(t *testing.T) {
	c := ColorByValue()
	if c.Mode != "thresholds" {
		t.Errorf("Mode = %q, want thresholds", c.Mode)
	}
}

func TestColorByPalette(t *testing.T) {
	c := ColorByPalette()
	if c.Mode != "palette-classic" {
		t.Errorf("Mode = %q, want palette-classic", c.Mode)
	}
}

func TestFixedColor(t *testing.T) {
	c := FixedColor("red")
	if c.Mode != "fixed" {
		t.Errorf("Mode = %q, want fixed", c.Mode)
	}
	if c.FixedColor != "red" {
		t.Errorf("FixedColor = %q, want red", c.FixedColor)
	}
}

func TestFieldOverrideByName(t *testing.T) {
	override := OverrideByName("requests").
		SetUnit(UnitShort).
		SetDecimals(0).
		Build()

	if override.Matcher.ID != "byName" {
		t.Errorf("Matcher.ID = %q, want byName", override.Matcher.ID)
	}
	if override.Matcher.Options != "requests" {
		t.Errorf("Matcher.Options = %q, want requests", override.Matcher.Options)
	}
	if len(override.Properties) != 2 {
		t.Errorf("len(Properties) = %d, want 2", len(override.Properties))
	}
}

func TestFieldOverrideByRegex(t *testing.T) {
	override := OverrideByRegex(".*_total").Build()

	if override.Matcher.ID != "byRegexp" {
		t.Errorf("Matcher.ID = %q, want byRegexp", override.Matcher.ID)
	}
	if override.Matcher.Options != ".*_total" {
		t.Errorf("Matcher.Options = %q, want .*_total", override.Matcher.Options)
	}
}

func TestFieldConfig_AddOverride(t *testing.T) {
	fc := NewFieldConfig().
		AddOverride(OverrideByName("value").SetUnit(UnitPercent))

	if len(fc.Overrides) != 1 {
		t.Errorf("len(Overrides) = %d, want 1", len(fc.Overrides))
	}
}
