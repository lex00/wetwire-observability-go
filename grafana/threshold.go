package grafana

// Threshold mode constants.
const (
	ThresholdModeAbsolute   = "absolute"
	ThresholdModePercentage = "percentage"
)

// Common Grafana colors.
const (
	ColorGreen        = "green"
	ColorYellow       = "yellow"
	ColorOrange       = "orange"
	ColorRed          = "red"
	ColorBlue         = "blue"
	ColorPurple       = "purple"
	ColorSuperLightGreen = "super-light-green"
	ColorLightGreen   = "light-green"
	ColorSemiDarkGreen = "semi-dark-green"
	ColorDarkGreen    = "dark-green"
	ColorSuperLightYellow = "super-light-yellow"
	ColorLightYellow  = "light-yellow"
	ColorSemiDarkYellow = "semi-dark-yellow"
	ColorDarkYellow   = "dark-yellow"
	ColorSuperLightOrange = "super-light-orange"
	ColorLightOrange  = "light-orange"
	ColorSemiDarkOrange = "semi-dark-orange"
	ColorDarkOrange   = "dark-orange"
	ColorSuperLightRed = "super-light-red"
	ColorLightRed     = "light-red"
	ColorSemiDarkRed  = "semi-dark-red"
	ColorDarkRed      = "dark-red"
)

// BaseStep creates a base threshold step (no value, just color).
// This is typically used as the first step in a threshold configuration.
func BaseStep(color string) ThresholdStep {
	return ThresholdStep{
		Value: nil,
		Color: color,
	}
}

// Step creates a threshold step with a specific value and color.
func Step(value float64, color string) ThresholdStep {
	return ThresholdStep{
		Value: &value,
		Color: color,
	}
}

// AbsoluteThresholds creates a threshold configuration with absolute mode.
// Values are compared directly against the data values.
func AbsoluteThresholds(steps ...ThresholdStep) *ThresholdStyle {
	return &ThresholdStyle{
		Mode:  ThresholdModeAbsolute,
		Steps: steps,
	}
}

// PercentageThresholds creates a threshold configuration with percentage mode.
// Values are interpreted as percentages of the range (min to max).
func PercentageThresholds(steps ...ThresholdStep) *ThresholdStyle {
	return &ThresholdStyle{
		Mode:  ThresholdModePercentage,
		Steps: steps,
	}
}

// GreenYellowRed creates a standard 3-step threshold (green -> yellow -> red).
// This is useful for metrics where higher values are worse (e.g., error rate, latency).
func GreenYellowRed(yellowAt, redAt float64) *ThresholdStyle {
	return AbsoluteThresholds(
		BaseStep(ColorGreen),
		Step(yellowAt, ColorYellow),
		Step(redAt, ColorRed),
	)
}

// RedYellowGreen creates a reverse 3-step threshold (red -> yellow -> green).
// This is useful for metrics where higher values are better (e.g., availability, success rate).
func RedYellowGreen(yellowAt, greenAt float64) *ThresholdStyle {
	return AbsoluteThresholds(
		BaseStep(ColorRed),
		Step(yellowAt, ColorYellow),
		Step(greenAt, ColorGreen),
	)
}

// GreenRed creates a simple 2-step threshold (green -> red).
func GreenRed(redAt float64) *ThresholdStyle {
	return AbsoluteThresholds(
		BaseStep(ColorGreen),
		Step(redAt, ColorRed),
	)
}

// RedGreen creates a simple 2-step threshold (red -> green).
func RedGreen(greenAt float64) *ThresholdStyle {
	return AbsoluteThresholds(
		BaseStep(ColorRed),
		Step(greenAt, ColorGreen),
	)
}

// SingleColor creates a threshold with a single color (no color changes).
func SingleColor(color string) *ThresholdStyle {
	return AbsoluteThresholds(
		BaseStep(color),
	)
}

// SLOThresholds creates thresholds suitable for SLO displays.
// Green when meeting objective, yellow when at risk, red when breaching.
func SLOThresholds(targetPercent float64) *ThresholdStyle {
	// When SLO target is e.g. 99.9%, values below 99% are red, 99-99.9% yellow, above green
	warningAt := targetPercent - (100 - targetPercent) // e.g., 99.9 - 0.1 = 99.8
	if warningAt < 0 {
		warningAt = 0
	}
	return AbsoluteThresholds(
		BaseStep(ColorRed),
		Step(warningAt, ColorYellow),
		Step(targetPercent, ColorGreen),
	)
}
