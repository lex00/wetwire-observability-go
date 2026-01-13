package grafana

// TimeSeriesPanel represents a Grafana time series panel.
type TimeSeriesPanel struct {
	BasePanel
	Options TimeSeriesOptions `json:"options,omitempty"`
}

// TimeSeriesOptions contains time series panel options.
type TimeSeriesOptions struct {
	Legend  LegendOptions  `json:"legend,omitempty"`
	Tooltip TooltipOptions `json:"tooltip,omitempty"`
}

// LegendOptions contains legend configuration.
type LegendOptions struct {
	DisplayMode string   `json:"displayMode,omitempty"`
	Placement   string   `json:"placement,omitempty"`
	ShowLegend  bool     `json:"showLegend,omitempty"`
	Calcs       []string `json:"calcs,omitempty"`
}

// TooltipOptions contains tooltip configuration.
type TooltipOptions struct {
	Mode string `json:"mode,omitempty"`
	Sort string `json:"sort,omitempty"`
}

// TimeSeries creates a new TimeSeriesPanel.
func TimeSeries(title string) *TimeSeriesPanel {
	return &TimeSeriesPanel{
		BasePanel: BasePanel{
			Type:  "timeseries",
			Title: title,
			GridPos: GridPos{
				W: 12,
				H: 8,
			},
		},
		Options: TimeSeriesOptions{
			Legend: LegendOptions{
				DisplayMode: "list",
				Placement:   LegendBottom,
				ShowLegend:  true,
			},
			Tooltip: TooltipOptions{
				Mode: TooltipSingle,
			},
		},
	}
}

// WithDescription sets the panel description.
func (p *TimeSeriesPanel) WithDescription(desc string) *TimeSeriesPanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *TimeSeriesPanel) WithDatasource(ds string) *TimeSeriesPanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *TimeSeriesPanel) WithSize(w, h int) *TimeSeriesPanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *TimeSeriesPanel) WithPosition(x, y int) *TimeSeriesPanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *TimeSeriesPanel) WithTargets(targets ...any) *TimeSeriesPanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *TimeSeriesPanel) AddTarget(target any) *TimeSeriesPanel {
	p.Targets = append(p.Targets, target)
	return p
}

// WithLegendPosition sets the legend position.
func (p *TimeSeriesPanel) WithLegendPosition(placement string) *TimeSeriesPanel {
	p.Options.Legend.DisplayMode = "list"
	p.Options.Legend.Placement = placement
	p.Options.Legend.ShowLegend = true
	return p
}

// HideLegend hides the legend.
func (p *TimeSeriesPanel) HideLegend() *TimeSeriesPanel {
	p.Options.Legend.DisplayMode = "hidden"
	p.Options.Legend.ShowLegend = false
	return p
}

// WithTooltip sets the tooltip mode.
func (p *TimeSeriesPanel) WithTooltip(mode string) *TimeSeriesPanel {
	p.Options.Tooltip.Mode = mode
	return p
}

// WithLineWidth sets the line width.
func (p *TimeSeriesPanel) WithLineWidth(width int) *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.LineWidth = width
	return p
}

// WithFillOpacity sets the fill opacity (0-100).
func (p *TimeSeriesPanel) WithFillOpacity(opacity int) *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.FillOpacity = opacity
	return p
}

// DrawBars sets the draw style to bars.
func (p *TimeSeriesPanel) DrawBars() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.DrawStyle = DrawStyleBars
	return p
}

// DrawPoints sets the draw style to points.
func (p *TimeSeriesPanel) DrawPoints() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.DrawStyle = DrawStylePoints
	return p
}

// DrawLine sets the draw style to line (default).
func (p *TimeSeriesPanel) DrawLine() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.DrawStyle = DrawStyleLine
	return p
}

// WithUnit sets the display unit.
func (p *TimeSeriesPanel) WithUnit(unit string) *TimeSeriesPanel {
	p.FieldConfig.Defaults.Unit = unit
	return p
}

// Stacked enables stacking.
func (p *TimeSeriesPanel) Stacked() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.Stacking = &StackingConfig{
		Mode:  "normal",
		Group: "A",
	}
	return p
}

// StackedPercent enables percentage stacking.
func (p *TimeSeriesPanel) StackedPercent() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.Stacking = &StackingConfig{
		Mode:  "percent",
		Group: "A",
	}
	return p
}

// SpanNulls enables spanning nulls with lines.
func (p *TimeSeriesPanel) SpanNulls() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.SpanNulls = true
	return p
}

// GradientFill enables gradient fill.
func (p *TimeSeriesPanel) GradientFill() *TimeSeriesPanel {
	p.FieldConfig.Defaults.Custom.GradientMode = "opacity"
	return p
}

// WithMin sets the Y-axis minimum.
func (p *TimeSeriesPanel) WithMin(min float64) *TimeSeriesPanel {
	p.FieldConfig.Defaults.Min = &min
	return p
}

// WithMax sets the Y-axis maximum.
func (p *TimeSeriesPanel) WithMax(max float64) *TimeSeriesPanel {
	p.FieldConfig.Defaults.Max = &max
	return p
}

// WithDecimals sets the number of decimal places.
func (p *TimeSeriesPanel) WithDecimals(decimals int) *TimeSeriesPanel {
	p.FieldConfig.Defaults.Decimals = &decimals
	return p
}
