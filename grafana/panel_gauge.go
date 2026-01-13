package grafana

// GaugePanel represents a Grafana gauge panel.
type GaugePanel struct {
	BasePanel
	Options GaugeOptions `json:"options,omitempty"`
}

// GaugeOptions contains gauge panel options.
type GaugeOptions struct {
	ShowThresholdLabels  bool          `json:"showThresholdLabels,omitempty"`
	ShowThresholdMarkers bool          `json:"showThresholdMarkers,omitempty"`
	ReduceOptions        ReduceOptions `json:"reduceOptions,omitempty"`
	Orientation          string        `json:"orientation,omitempty"`
	TextMode             string        `json:"textMode,omitempty"`
}

// Gauge creates a new GaugePanel.
func Gauge(title string) *GaugePanel {
	return &GaugePanel{
		BasePanel: BasePanel{
			Type:  "gauge",
			Title: title,
			GridPos: GridPos{
				W: 8,
				H: 6,
			},
		},
		Options: GaugeOptions{
			ShowThresholdMarkers: true,
			ReduceOptions: ReduceOptions{
				Calcs: []string{ReduceLastNotNA},
			},
			Orientation: OrientationAuto,
		},
	}
}

// WithDescription sets the panel description.
func (p *GaugePanel) WithDescription(desc string) *GaugePanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *GaugePanel) WithDatasource(ds string) *GaugePanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *GaugePanel) WithSize(w, h int) *GaugePanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *GaugePanel) WithPosition(x, y int) *GaugePanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *GaugePanel) WithTargets(targets ...any) *GaugePanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *GaugePanel) AddTarget(target any) *GaugePanel {
	p.Targets = append(p.Targets, target)
	return p
}

// WithMin sets the minimum value.
func (p *GaugePanel) WithMin(min float64) *GaugePanel {
	p.FieldConfig.Defaults.Min = &min
	return p
}

// WithMax sets the maximum value.
func (p *GaugePanel) WithMax(max float64) *GaugePanel {
	p.FieldConfig.Defaults.Max = &max
	return p
}

// WithUnit sets the display unit.
func (p *GaugePanel) WithUnit(unit string) *GaugePanel {
	p.FieldConfig.Defaults.Unit = unit
	return p
}

// WithDecimals sets the number of decimal places.
func (p *GaugePanel) WithDecimals(decimals int) *GaugePanel {
	p.FieldConfig.Defaults.Decimals = &decimals
	return p
}

// ShowThresholdLabels shows threshold labels on the gauge.
func (p *GaugePanel) ShowThresholdLabels() *GaugePanel {
	p.Options.ShowThresholdLabels = true
	return p
}

// HideThresholdLabels hides threshold labels.
func (p *GaugePanel) HideThresholdLabels() *GaugePanel {
	p.Options.ShowThresholdLabels = false
	return p
}

// ShowThresholdMarkers shows threshold markers on the gauge.
func (p *GaugePanel) ShowThresholdMarkers() *GaugePanel {
	p.Options.ShowThresholdMarkers = true
	return p
}

// HideThresholdMarkers hides threshold markers.
func (p *GaugePanel) HideThresholdMarkers() *GaugePanel {
	p.Options.ShowThresholdMarkers = false
	return p
}

// WithReduceCalc sets the reduce calculation.
func (p *GaugePanel) WithReduceCalc(calc string) *GaugePanel {
	p.Options.ReduceOptions.Calcs = []string{calc}
	return p
}

// Horizontal sets horizontal orientation.
func (p *GaugePanel) Horizontal() *GaugePanel {
	p.Options.Orientation = OrientationHorizontal
	return p
}

// Vertical sets vertical orientation.
func (p *GaugePanel) Vertical() *GaugePanel {
	p.Options.Orientation = OrientationVertical
	return p
}
