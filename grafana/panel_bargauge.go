package grafana

// Bar gauge display mode options.
const (
	DisplayModeBasic    = "basic"
	DisplayModeGradient = "gradient"
	DisplayModeLCD      = "lcd"
)

// BarGaugePanel represents a Grafana bar gauge panel.
type BarGaugePanel struct {
	BasePanel
	Options BarGaugeOptions `json:"options,omitempty"`
}

// BarGaugeOptions contains bar gauge panel options.
type BarGaugeOptions struct {
	DisplayMode   string        `json:"displayMode,omitempty"`
	Orientation   string        `json:"orientation,omitempty"`
	ReduceOptions ReduceOptions `json:"reduceOptions,omitempty"`
	ShowUnfilled  bool          `json:"showUnfilled,omitempty"`
	MinVizWidth   int           `json:"minVizWidth,omitempty"`
	MinVizHeight  int           `json:"minVizHeight,omitempty"`
}

// BarGauge creates a new BarGaugePanel.
func BarGauge(title string) *BarGaugePanel {
	return &BarGaugePanel{
		BasePanel: BasePanel{
			Type:  "bargauge",
			Title: title,
			GridPos: GridPos{
				W: 8,
				H: 6,
			},
		},
		Options: BarGaugeOptions{
			DisplayMode:  DisplayModeGradient,
			Orientation:  OrientationHorizontal,
			ShowUnfilled: true,
			ReduceOptions: ReduceOptions{
				Calcs: []string{ReduceLastNotNA},
			},
		},
	}
}

// WithDescription sets the panel description.
func (p *BarGaugePanel) WithDescription(desc string) *BarGaugePanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *BarGaugePanel) WithDatasource(ds string) *BarGaugePanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *BarGaugePanel) WithSize(w, h int) *BarGaugePanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *BarGaugePanel) WithPosition(x, y int) *BarGaugePanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *BarGaugePanel) WithTargets(targets ...any) *BarGaugePanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *BarGaugePanel) AddTarget(target any) *BarGaugePanel {
	p.Targets = append(p.Targets, target)
	return p
}

// WithMin sets the minimum value.
func (p *BarGaugePanel) WithMin(min float64) *BarGaugePanel {
	p.FieldConfig.Defaults.Min = &min
	return p
}

// WithMax sets the maximum value.
func (p *BarGaugePanel) WithMax(max float64) *BarGaugePanel {
	p.FieldConfig.Defaults.Max = &max
	return p
}

// WithUnit sets the display unit.
func (p *BarGaugePanel) WithUnit(unit string) *BarGaugePanel {
	p.FieldConfig.Defaults.Unit = unit
	return p
}

// WithDecimals sets the number of decimal places.
func (p *BarGaugePanel) WithDecimals(decimals int) *BarGaugePanel {
	p.FieldConfig.Defaults.Decimals = &decimals
	return p
}

// Basic sets the display mode to basic.
func (p *BarGaugePanel) Basic() *BarGaugePanel {
	p.Options.DisplayMode = DisplayModeBasic
	return p
}

// Gradient sets the display mode to gradient.
func (p *BarGaugePanel) Gradient() *BarGaugePanel {
	p.Options.DisplayMode = DisplayModeGradient
	return p
}

// LCD sets the display mode to LCD.
func (p *BarGaugePanel) LCD() *BarGaugePanel {
	p.Options.DisplayMode = DisplayModeLCD
	return p
}

// Horizontal sets horizontal orientation.
func (p *BarGaugePanel) Horizontal() *BarGaugePanel {
	p.Options.Orientation = OrientationHorizontal
	return p
}

// Vertical sets vertical orientation.
func (p *BarGaugePanel) Vertical() *BarGaugePanel {
	p.Options.Orientation = OrientationVertical
	return p
}

// ShowUnfilled shows the unfilled portion of the bar.
func (p *BarGaugePanel) ShowUnfilled() *BarGaugePanel {
	p.Options.ShowUnfilled = true
	return p
}

// HideUnfilled hides the unfilled portion of the bar.
func (p *BarGaugePanel) HideUnfilled() *BarGaugePanel {
	p.Options.ShowUnfilled = false
	return p
}

// WithReduceCalc sets the reduce calculation.
func (p *BarGaugePanel) WithReduceCalc(calc string) *BarGaugePanel {
	p.Options.ReduceOptions.Calcs = []string{calc}
	return p
}
