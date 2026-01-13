package grafana

import "strconv"

// Color scheme presets.
const (
	ColorSchemeBlues    = "Blues"
	ColorSchemeReds     = "Reds"
	ColorSchemeGreens   = "Greens"
	ColorSchemeSpectral = "Spectral"
	ColorSchemeYlOrRd   = "YlOrRd"
	ColorSchemeYlGnBu   = "YlGnBu"
)

// HeatmapPanel represents a Grafana heatmap panel.
type HeatmapPanel struct {
	BasePanel
	Options HeatmapOptions `json:"options,omitempty"`
}

// HeatmapOptions contains heatmap panel options.
type HeatmapOptions struct {
	Calculate   bool                  `json:"calculate,omitempty"`
	Calculation HeatmapCalculation    `json:"calculation,omitempty"`
	Color       HeatmapColor          `json:"color,omitempty"`
	CellGap     int                   `json:"cellGap,omitempty"`
	Legend      HeatmapLegend         `json:"legend,omitempty"`
	Tooltip     HeatmapTooltip        `json:"tooltip,omitempty"`
	YAxis       HeatmapYAxis          `json:"yAxis,omitempty"`
	ShowValue   string                `json:"showValue,omitempty"`
}

// HeatmapCalculation contains calculation options.
type HeatmapCalculation struct {
	XBuckets HeatmapBucketConfig `json:"xBuckets,omitempty"`
	YBuckets HeatmapBucketConfig `json:"yBuckets,omitempty"`
}

// HeatmapBucketConfig contains bucket configuration.
type HeatmapBucketConfig struct {
	Mode  string `json:"mode,omitempty"`
	Value string `json:"value,omitempty"`
	Scale HeatmapScale `json:"scale,omitempty"`
}

// HeatmapScale contains scale configuration.
type HeatmapScale struct {
	Type string `json:"type,omitempty"`
}

// HeatmapColor contains color configuration.
type HeatmapColor struct {
	Mode     string  `json:"mode,omitempty"`
	Scheme   string  `json:"scheme,omitempty"`
	Fill     string  `json:"fill,omitempty"`
	Min      float64 `json:"min,omitempty"`
	Max      float64 `json:"max,omitempty"`
	Exponent float64 `json:"exponent,omitempty"`
}

// HeatmapLegend contains legend configuration.
type HeatmapLegend struct {
	Show bool `json:"show,omitempty"`
}

// HeatmapTooltip contains tooltip configuration.
type HeatmapTooltip struct {
	Show       bool `json:"show,omitempty"`
	YHistogram bool `json:"yHistogram,omitempty"`
}

// HeatmapYAxis contains Y axis configuration.
type HeatmapYAxis struct {
	AxisDisplay string `json:"axisDisplay,omitempty"`
	Unit        string `json:"unit,omitempty"`
	Reverse     bool   `json:"reverse,omitempty"`
}

// Heatmap creates a new HeatmapPanel.
func Heatmap(title string) *HeatmapPanel {
	return &HeatmapPanel{
		BasePanel: BasePanel{
			Type:  "heatmap",
			Title: title,
			GridPos: GridPos{
				W: 24,
				H: 10,
			},
		},
		Options: HeatmapOptions{
			Color: HeatmapColor{
				Mode:   "scheme",
				Scheme: ColorSchemeSpectral,
			},
			CellGap: 1,
			Legend: HeatmapLegend{
				Show: true,
			},
			Tooltip: HeatmapTooltip{
				Show: true,
			},
		},
	}
}

// WithDescription sets the panel description.
func (p *HeatmapPanel) WithDescription(desc string) *HeatmapPanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *HeatmapPanel) WithDatasource(ds string) *HeatmapPanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *HeatmapPanel) WithSize(w, h int) *HeatmapPanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *HeatmapPanel) WithPosition(x, y int) *HeatmapPanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *HeatmapPanel) WithTargets(targets ...any) *HeatmapPanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *HeatmapPanel) AddTarget(target any) *HeatmapPanel {
	p.Targets = append(p.Targets, target)
	return p
}

// WithColorScheme sets the color scheme.
func (p *HeatmapPanel) WithColorScheme(scheme string) *HeatmapPanel {
	p.Options.Color.Scheme = scheme
	return p
}

// Blues sets the color scheme to Blues.
func (p *HeatmapPanel) Blues() *HeatmapPanel {
	p.Options.Color.Scheme = ColorSchemeBlues
	return p
}

// Reds sets the color scheme to Reds.
func (p *HeatmapPanel) Reds() *HeatmapPanel {
	p.Options.Color.Scheme = ColorSchemeReds
	return p
}

// Greens sets the color scheme to Greens.
func (p *HeatmapPanel) Greens() *HeatmapPanel {
	p.Options.Color.Scheme = ColorSchemeGreens
	return p
}

// Spectral sets the color scheme to Spectral.
func (p *HeatmapPanel) Spectral() *HeatmapPanel {
	p.Options.Color.Scheme = ColorSchemeSpectral
	return p
}

// ShowLegend shows the legend.
func (p *HeatmapPanel) ShowLegend() *HeatmapPanel {
	p.Options.Legend.Show = true
	return p
}

// HideLegend hides the legend.
func (p *HeatmapPanel) HideLegend() *HeatmapPanel {
	p.Options.Legend.Show = false
	return p
}

// ShowTooltip shows the tooltip.
func (p *HeatmapPanel) ShowTooltip() *HeatmapPanel {
	p.Options.Tooltip.Show = true
	return p
}

// HideTooltip hides the tooltip.
func (p *HeatmapPanel) HideTooltip() *HeatmapPanel {
	p.Options.Tooltip.Show = false
	return p
}

// ShowYHistogram shows Y histogram in the tooltip.
func (p *HeatmapPanel) ShowYHistogram() *HeatmapPanel {
	p.Options.Tooltip.YHistogram = true
	return p
}

// HideYHistogram hides Y histogram in the tooltip.
func (p *HeatmapPanel) HideYHistogram() *HeatmapPanel {
	p.Options.Tooltip.YHistogram = false
	return p
}

// Calculate enables heatmap calculation from raw data.
func (p *HeatmapPanel) Calculate() *HeatmapPanel {
	p.Options.Calculate = true
	return p
}

// NoCalculate disables heatmap calculation (use pre-bucketed data).
func (p *HeatmapPanel) NoCalculate() *HeatmapPanel {
	p.Options.Calculate = false
	return p
}

// WithXBucketSize sets the X bucket size.
func (p *HeatmapPanel) WithXBucketSize(size int) *HeatmapPanel {
	p.Options.Calculation.XBuckets.Value = strconv.Itoa(size)
	return p
}

// WithYBucketSize sets the Y bucket size.
func (p *HeatmapPanel) WithYBucketSize(size int) *HeatmapPanel {
	p.Options.Calculation.YBuckets.Value = strconv.Itoa(size)
	return p
}

// WithCellGap sets the gap between cells.
func (p *HeatmapPanel) WithCellGap(gap int) *HeatmapPanel {
	p.Options.CellGap = gap
	return p
}

// WithUnit sets the Y axis unit.
func (p *HeatmapPanel) WithUnit(unit string) *HeatmapPanel {
	p.Options.YAxis.Unit = unit
	return p
}

// ReverseYAxis reverses the Y axis direction.
func (p *HeatmapPanel) ReverseYAxis() *HeatmapPanel {
	p.Options.YAxis.Reverse = true
	return p
}
