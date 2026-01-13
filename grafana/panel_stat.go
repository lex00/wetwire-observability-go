package grafana

// StatPanel represents a Grafana stat panel.
type StatPanel struct {
	BasePanel
	Options StatOptions `json:"options,omitempty"`
}

// StatOptions contains stat panel options.
type StatOptions struct {
	ColorMode     string        `json:"colorMode,omitempty"`
	GraphMode     string        `json:"graphMode,omitempty"`
	JustifyMode   string        `json:"justifyMode,omitempty"`
	Orientation   string        `json:"orientation,omitempty"`
	ReduceOptions ReduceOptions `json:"reduceOptions,omitempty"`
	TextMode      string        `json:"textMode,omitempty"`
}

// ReduceOptions contains reduce calculation options.
type ReduceOptions struct {
	Values bool     `json:"values,omitempty"`
	Calcs  []string `json:"calcs,omitempty"`
	Fields string   `json:"fields,omitempty"`
	Limit  *int     `json:"limit,omitempty"`
}

// Stat creates a new StatPanel.
func Stat(title string) *StatPanel {
	return &StatPanel{
		BasePanel: BasePanel{
			Type:  "stat",
			Title: title,
			GridPos: GridPos{
				W: 6,
				H: 4,
			},
		},
		Options: StatOptions{
			ColorMode:   ColorModeValue,
			GraphMode:   GraphModeNone,
			Orientation: OrientationAuto,
			ReduceOptions: ReduceOptions{
				Calcs: []string{ReduceLastNotNA},
			},
		},
	}
}

// WithDescription sets the panel description.
func (p *StatPanel) WithDescription(desc string) *StatPanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *StatPanel) WithDatasource(ds string) *StatPanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *StatPanel) WithSize(w, h int) *StatPanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *StatPanel) WithPosition(x, y int) *StatPanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *StatPanel) WithTargets(targets ...any) *StatPanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *StatPanel) AddTarget(target any) *StatPanel {
	p.Targets = append(p.Targets, target)
	return p
}

// ColorByValue colors the text by value (default).
func (p *StatPanel) ColorByValue() *StatPanel {
	p.Options.ColorMode = ColorModeValue
	return p
}

// ColorByBackground colors the background.
func (p *StatPanel) ColorByBackground() *StatPanel {
	p.Options.ColorMode = ColorModeBackground
	return p
}

// NoColor disables coloring.
func (p *StatPanel) NoColor() *StatPanel {
	p.Options.ColorMode = ColorModeNone
	return p
}

// WithReduceCalc sets the reduce calculation.
func (p *StatPanel) WithReduceCalc(calc string) *StatPanel {
	p.Options.ReduceOptions.Calcs = []string{calc}
	return p
}

// WithReduceCalcs sets multiple reduce calculations.
func (p *StatPanel) WithReduceCalcs(calcs ...string) *StatPanel {
	p.Options.ReduceOptions.Calcs = calcs
	return p
}

// WithGraphMode sets the graph mode.
func (p *StatPanel) WithGraphMode(mode string) *StatPanel {
	p.Options.GraphMode = mode
	return p
}

// NoGraph disables the sparkline graph.
func (p *StatPanel) NoGraph() *StatPanel {
	p.Options.GraphMode = GraphModeNone
	return p
}

// WithSparkline enables the area sparkline.
func (p *StatPanel) WithSparkline() *StatPanel {
	p.Options.GraphMode = GraphModeArea
	return p
}

// WithOrientation sets the panel orientation.
func (p *StatPanel) WithOrientation(orientation string) *StatPanel {
	p.Options.Orientation = orientation
	return p
}

// Horizontal sets horizontal orientation.
func (p *StatPanel) Horizontal() *StatPanel {
	p.Options.Orientation = OrientationHorizontal
	return p
}

// Vertical sets vertical orientation.
func (p *StatPanel) Vertical() *StatPanel {
	p.Options.Orientation = OrientationVertical
	return p
}

// WithUnit sets the display unit.
func (p *StatPanel) WithUnit(unit string) *StatPanel {
	p.FieldConfig.Defaults.Unit = unit
	return p
}

// WithMin sets the minimum value.
func (p *StatPanel) WithMin(min float64) *StatPanel {
	p.FieldConfig.Defaults.Min = &min
	return p
}

// WithMax sets the maximum value.
func (p *StatPanel) WithMax(max float64) *StatPanel {
	p.FieldConfig.Defaults.Max = &max
	return p
}

// WithDecimals sets the number of decimal places.
func (p *StatPanel) WithDecimals(decimals int) *StatPanel {
	p.FieldConfig.Defaults.Decimals = &decimals
	return p
}

// WithNoValue sets the text to display when there's no value.
func (p *StatPanel) WithNoValue(text string) *StatPanel {
	p.FieldConfig.Defaults.NoValue = text
	return p
}
