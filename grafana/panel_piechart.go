package grafana

// Pie type options.
const (
	PieTypePie   = "pie"
	PieTypeDonut = "donut"
)

// PieChartPanel represents a Grafana pie chart panel.
type PieChartPanel struct {
	BasePanel
	Options PieChartOptions `json:"options,omitempty"`
}

// PieChartOptions contains pie chart panel options.
type PieChartOptions struct {
	PieType       string             `json:"pieType,omitempty"`
	ReduceOptions ReduceOptions      `json:"reduceOptions,omitempty"`
	Legend        PieChartLegend     `json:"legend,omitempty"`
	Tooltip       PieChartTooltip    `json:"tooltip,omitempty"`
	Labels        bool               `json:"displayLabels,omitempty"`
	Orientation   string             `json:"orientation,omitempty"`
}

// PieChartLegend contains legend options for pie charts.
type PieChartLegend struct {
	DisplayMode string   `json:"displayMode,omitempty"`
	Placement   string   `json:"placement,omitempty"`
	ShowLegend  bool     `json:"showLegend,omitempty"`
	Values      []string `json:"values,omitempty"`
}

// PieChartTooltip contains tooltip options for pie charts.
type PieChartTooltip struct {
	Mode string `json:"mode,omitempty"`
	Sort string `json:"sort,omitempty"`
}

// PieChart creates a new PieChartPanel.
func PieChart(title string) *PieChartPanel {
	return &PieChartPanel{
		BasePanel: BasePanel{
			Type:  "piechart",
			Title: title,
			GridPos: GridPos{
				W: 10,
				H: 10,
			},
		},
		Options: PieChartOptions{
			PieType: PieTypePie,
			ReduceOptions: ReduceOptions{
				Calcs: []string{ReduceLastNotNA},
			},
			Legend: PieChartLegend{
				DisplayMode: "list",
				Placement:   LegendRight,
				ShowLegend:  true,
			},
			Tooltip: PieChartTooltip{
				Mode: TooltipSingle,
			},
		},
	}
}

// WithDescription sets the panel description.
func (p *PieChartPanel) WithDescription(desc string) *PieChartPanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *PieChartPanel) WithDatasource(ds string) *PieChartPanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *PieChartPanel) WithSize(w, h int) *PieChartPanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *PieChartPanel) WithPosition(x, y int) *PieChartPanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *PieChartPanel) WithTargets(targets ...any) *PieChartPanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *PieChartPanel) AddTarget(target any) *PieChartPanel {
	p.Targets = append(p.Targets, target)
	return p
}

// Pie sets the chart type to pie.
func (p *PieChartPanel) Pie() *PieChartPanel {
	p.Options.PieType = PieTypePie
	return p
}

// Donut sets the chart type to donut.
func (p *PieChartPanel) Donut() *PieChartPanel {
	p.Options.PieType = PieTypeDonut
	return p
}

// LegendRight places the legend on the right.
func (p *PieChartPanel) LegendRight() *PieChartPanel {
	p.Options.Legend.Placement = LegendRight
	p.Options.Legend.ShowLegend = true
	return p
}

// LegendBottom places the legend at the bottom.
func (p *PieChartPanel) LegendBottom() *PieChartPanel {
	p.Options.Legend.Placement = LegendBottom
	p.Options.Legend.ShowLegend = true
	return p
}

// HideLegend hides the legend.
func (p *PieChartPanel) HideLegend() *PieChartPanel {
	p.Options.Legend.DisplayMode = "hidden"
	p.Options.Legend.ShowLegend = false
	return p
}

// ShowLabels shows labels on the pie slices.
func (p *PieChartPanel) ShowLabels() *PieChartPanel {
	p.Options.Labels = true
	return p
}

// HideLabels hides labels on the pie slices.
func (p *PieChartPanel) HideLabels() *PieChartPanel {
	p.Options.Labels = false
	return p
}

// ShowTooltip enables the tooltip.
func (p *PieChartPanel) ShowTooltip() *PieChartPanel {
	p.Options.Tooltip.Mode = TooltipSingle
	return p
}

// WithUnit sets the display unit.
func (p *PieChartPanel) WithUnit(unit string) *PieChartPanel {
	p.FieldConfig.Defaults.Unit = unit
	return p
}

// WithDecimals sets the number of decimal places.
func (p *PieChartPanel) WithDecimals(decimals int) *PieChartPanel {
	p.FieldConfig.Defaults.Decimals = &decimals
	return p
}

// WithReduceCalc sets the reduce calculation.
func (p *PieChartPanel) WithReduceCalc(calc string) *PieChartPanel {
	p.Options.ReduceOptions.Calcs = []string{calc}
	return p
}
