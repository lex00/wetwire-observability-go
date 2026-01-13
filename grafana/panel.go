package grafana

// Common units.
const (
	UnitNone         = ""
	UnitShort        = "short"
	UnitPercent      = "percent"
	UnitPercentUnit  = "percentunit"
	UnitBytes        = "bytes"
	UnitDecBytes     = "decbytes"
	UnitBitsPerSec   = "bps"
	UnitBytesPerSec  = "Bps"
	UnitSeconds      = "s"
	UnitMilliseconds = "ms"
	UnitMicroseconds = "µs"
	UnitDateTimeISO  = "dateTimeAsIso"
	UnitDateTimeUS   = "dateTimeAsUS"
)

// Legend placement options.
const (
	LegendBottom = "bottom"
	LegendRight  = "right"
)

// Tooltip mode options.
const (
	TooltipSingle = "single"
	TooltipAll    = "all"
	TooltipNone   = "none"
)

// Draw style options.
const (
	DrawStyleLine   = "line"
	DrawStyleBars   = "bars"
	DrawStylePoints = "points"
)

// Color mode options for stat panels.
const (
	ColorModeValue      = "value"
	ColorModeBackground = "background"
	ColorModeNone       = "none"
)

// Graph mode options for stat panels.
const (
	GraphModeNone = "none"
	GraphModeArea = "area"
)

// Orientation options for stat panels.
const (
	OrientationAuto       = "auto"
	OrientationHorizontal = "horizontal"
	OrientationVertical   = "vertical"
)

// Reduce calculation options.
const (
	ReduceLast      = "last"
	ReduceLastNotNA = "lastNotNull"
	ReduceFirst     = "first"
	ReduceMax       = "max"
	ReduceMin       = "min"
	ReduceMean      = "mean"
	ReduceSum       = "sum"
	ReduceCount     = "count"
	ReduceRange     = "range"
	ReduceDiff      = "diff"
	ReduceDelta     = "delta"
	ReduceStep      = "step"
)

// GridPos represents panel position and size in the dashboard grid.
type GridPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// Panel is the interface implemented by all panel types.
type Panel interface {
	GetTitle() string
	GetType() string
	GetGridPos() GridPos
}

// BasePanel contains fields common to all panel types.
type BasePanel struct {
	// ID is the unique panel ID (auto-assigned).
	ID int `json:"id,omitempty"`

	// Type is the panel type (timeseries, stat, table, etc.).
	Type string `json:"type"`

	// Title is the panel title.
	Title string `json:"title"`

	// Description is an optional panel description.
	Description string `json:"description,omitempty"`

	// Datasource is the data source reference.
	Datasource string `json:"datasource,omitempty"`

	// GridPos is the panel position and size.
	GridPos GridPos `json:"gridPos"`

	// Targets are the query targets.
	Targets []any `json:"targets,omitempty"`

	// FieldConfig contains field configuration.
	FieldConfig FieldConfig `json:"fieldConfig,omitempty"`

	// Links are panel links.
	Links []any `json:"links,omitempty"`

	// Transparent makes the panel background transparent.
	Transparent bool `json:"transparent,omitempty"`
}

// GetTitle returns the panel title.
func (b *BasePanel) GetTitle() string {
	return b.Title
}

// GetType returns the panel type.
func (b *BasePanel) GetType() string {
	return b.Type
}

// GetGridPos returns the panel position.
func (b *BasePanel) GetGridPos() GridPos {
	return b.GridPos
}

// FieldConfig contains field configuration.
type FieldConfig struct {
	Defaults  FieldDefaults   `json:"defaults,omitempty"`
	Overrides []FieldOverride `json:"overrides,omitempty"`
}

// FieldDefaults contains default field settings.
type FieldDefaults struct {
	Unit       string          `json:"unit,omitempty"`
	Min        *float64        `json:"min,omitempty"`
	Max        *float64        `json:"max,omitempty"`
	Decimals   *int            `json:"decimals,omitempty"`
	NoValue    string          `json:"noValue,omitempty"`
	Color      *ColorConfig    `json:"color,omitempty"`
	Thresholds *ThresholdStyle `json:"thresholds,omitempty"`
	Custom     CustomFieldConfig `json:"custom,omitempty"`
}

// ColorConfig represents color configuration.
type ColorConfig struct {
	Mode       string `json:"mode,omitempty"`
	FixedColor string `json:"fixedColor,omitempty"`
}

// ThresholdStyle represents threshold configuration.
type ThresholdStyle struct {
	Mode  string           `json:"mode,omitempty"`
	Steps []ThresholdStep `json:"steps,omitempty"`
}

// ThresholdStep represents a single threshold step.
type ThresholdStep struct {
	Value *float64 `json:"value,omitempty"`
	Color string   `json:"color"`
}

// CustomFieldConfig contains panel-specific custom field configuration.
type CustomFieldConfig struct {
	DrawStyle       string `json:"drawStyle,omitempty"`
	LineInterpolation string `json:"lineInterpolation,omitempty"`
	LineWidth       int    `json:"lineWidth,omitempty"`
	FillOpacity     int    `json:"fillOpacity,omitempty"`
	GradientMode    string `json:"gradientMode,omitempty"`
	SpanNulls       bool   `json:"spanNulls,omitempty"`
	ShowPoints      string `json:"showPoints,omitempty"`
	PointSize       int    `json:"pointSize,omitempty"`
	Stacking        *StackingConfig `json:"stacking,omitempty"`
	AxisPlacement   string `json:"axisPlacement,omitempty"`
	BarAlignment    int    `json:"barAlignment,omitempty"`
}

// StackingConfig represents stacking configuration.
type StackingConfig struct {
	Mode  string `json:"mode,omitempty"`
	Group string `json:"group,omitempty"`
}

// FieldOverride represents a field override.
type FieldOverride struct {
	Matcher    FieldMatcher `json:"matcher"`
	Properties []any        `json:"properties"`
}

// FieldMatcher identifies which fields to override.
type FieldMatcher struct {
	ID      string `json:"id"`
	Options string `json:"options,omitempty"`
}
