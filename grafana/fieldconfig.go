package grafana

// Color mode constants.
const (
	ColorModeThresholds    = "thresholds"
	ColorModePaletteClassic = "palette-classic"
	ColorModeFixed         = "fixed"
	ColorModeContinuous    = "continuous-GrYlRd"
)

// NewFieldConfig creates a new FieldConfig.
func NewFieldConfig() *FieldConfig {
	return &FieldConfig{}
}

// WithUnit sets the display unit.
func (fc *FieldConfig) WithUnit(unit string) *FieldConfig {
	fc.Defaults.Unit = unit
	return fc
}

// WithDecimals sets the number of decimal places.
func (fc *FieldConfig) WithDecimals(decimals int) *FieldConfig {
	fc.Defaults.Decimals = &decimals
	return fc
}

// WithMin sets the minimum value.
func (fc *FieldConfig) WithMin(min float64) *FieldConfig {
	fc.Defaults.Min = &min
	return fc
}

// WithMax sets the maximum value.
func (fc *FieldConfig) WithMax(max float64) *FieldConfig {
	fc.Defaults.Max = &max
	return fc
}

// WithNoValue sets the text to display when there's no value.
func (fc *FieldConfig) WithNoValue(text string) *FieldConfig {
	fc.Defaults.NoValue = text
	return fc
}

// WithThresholds sets the threshold configuration.
func (fc *FieldConfig) WithThresholds(thresholds *ThresholdStyle) *FieldConfig {
	fc.Defaults.Thresholds = thresholds
	return fc
}

// WithColor sets the color configuration.
func (fc *FieldConfig) WithColor(color *ColorConfig) *FieldConfig {
	fc.Defaults.Color = color
	return fc
}

// WithFixedColor sets a fixed color.
func (fc *FieldConfig) WithFixedColor(color string) *FieldConfig {
	fc.Defaults.Color = &ColorConfig{
		Mode:       ColorModeFixed,
		FixedColor: color,
	}
	return fc
}

// AddOverride adds a field override.
func (fc *FieldConfig) AddOverride(override *FieldOverrideBuilder) *FieldConfig {
	fc.Overrides = append(fc.Overrides, override.Build())
	return fc
}

// ColorByValue returns a color config that colors by threshold value.
func ColorByValue() *ColorConfig {
	return &ColorConfig{
		Mode: ColorModeThresholds,
	}
}

// ColorByPalette returns a color config that colors by palette.
func ColorByPalette() *ColorConfig {
	return &ColorConfig{
		Mode: ColorModePaletteClassic,
	}
}

// FixedColor returns a color config with a fixed color.
func FixedColor(color string) *ColorConfig {
	return &ColorConfig{
		Mode:       ColorModeFixed,
		FixedColor: color,
	}
}

// ContinuousColor returns a color config with continuous coloring.
func ContinuousColor() *ColorConfig {
	return &ColorConfig{
		Mode: ColorModeContinuous,
	}
}

// FieldOverrideBuilder helps build field overrides.
type FieldOverrideBuilder struct {
	matcher    FieldMatcher
	properties []OverrideProperty
}

// OverrideProperty represents a single override property.
type OverrideProperty struct {
	ID    string `json:"id"`
	Value any    `json:"value"`
}

// OverrideByName creates an override that matches fields by exact name.
func OverrideByName(name string) *FieldOverrideBuilder {
	return &FieldOverrideBuilder{
		matcher: FieldMatcher{
			ID:      "byName",
			Options: name,
		},
	}
}

// OverrideByRegex creates an override that matches fields by regex pattern.
func OverrideByRegex(pattern string) *FieldOverrideBuilder {
	return &FieldOverrideBuilder{
		matcher: FieldMatcher{
			ID:      "byRegexp",
			Options: pattern,
		},
	}
}

// OverrideByType creates an override that matches fields by type.
func OverrideByType(fieldType string) *FieldOverrideBuilder {
	return &FieldOverrideBuilder{
		matcher: FieldMatcher{
			ID:      "byType",
			Options: fieldType,
		},
	}
}

// SetUnit sets the unit override.
func (b *FieldOverrideBuilder) SetUnit(unit string) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "unit",
		Value: unit,
	})
	return b
}

// SetDecimals sets the decimals override.
func (b *FieldOverrideBuilder) SetDecimals(decimals int) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "decimals",
		Value: decimals,
	})
	return b
}

// SetMin sets the minimum value override.
func (b *FieldOverrideBuilder) SetMin(min float64) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "min",
		Value: min,
	})
	return b
}

// SetMax sets the maximum value override.
func (b *FieldOverrideBuilder) SetMax(max float64) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "max",
		Value: max,
	})
	return b
}

// SetColor sets the color override.
func (b *FieldOverrideBuilder) SetColor(color *ColorConfig) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "color",
		Value: color,
	})
	return b
}

// SetThresholds sets the thresholds override.
func (b *FieldOverrideBuilder) SetThresholds(thresholds *ThresholdStyle) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "thresholds",
		Value: thresholds,
	})
	return b
}

// SetDisplayName sets the display name override.
func (b *FieldOverrideBuilder) SetDisplayName(name string) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "displayName",
		Value: name,
	})
	return b
}

// SetNoValue sets the no value text override.
func (b *FieldOverrideBuilder) SetNoValue(text string) *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID:    "noValue",
		Value: text,
	})
	return b
}

// Hide hides the field.
func (b *FieldOverrideBuilder) Hide() *FieldOverrideBuilder {
	b.properties = append(b.properties, OverrideProperty{
		ID: "custom.hidden",
		Value: true,
	})
	return b
}

// Build creates the FieldOverride.
func (b *FieldOverrideBuilder) Build() FieldOverride {
	props := make([]any, len(b.properties))
	for i, p := range b.properties {
		props[i] = map[string]any{
			"id":    p.ID,
			"value": p.Value,
		}
	}
	return FieldOverride{
		Matcher:    b.matcher,
		Properties: props,
	}
}
