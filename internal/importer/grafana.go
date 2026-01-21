// Package importer provides functionality to import existing configuration files
// and generate equivalent Go code.
package importer

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/lex00/wetwire-observability-go/grafana"
)

// GrafanaDashboard represents a Grafana dashboard JSON structure for importing.
type GrafanaDashboard struct {
	UID           string         `json:"uid"`
	Title         string         `json:"title"`
	Description   string         `json:"description,omitempty"`
	Tags          []string       `json:"tags,omitempty"`
	Time          *TimeRange     `json:"time,omitempty"`
	Refresh       string         `json:"refresh,omitempty"`
	Timezone      string         `json:"timezone,omitempty"`
	Editable      bool           `json:"editable"`
	Panels        []GrafanaPanel `json:"panels"`
	Templating    *Templating    `json:"templating,omitempty"`
	SchemaVersion int            `json:"schemaVersion"`
	Version       int            `json:"version,omitempty"`
}

// TimeRange represents the dashboard time range.
type TimeRange struct {
	From string `json:"from"`
	To   string `json:"to"`
}

// Templating represents dashboard templating/variables.
type Templating struct {
	List []GrafanaVariable `json:"list"`
}

// GrafanaVariable represents a template variable.
type GrafanaVariable struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Label       string `json:"label,omitempty"`
	Description string `json:"description,omitempty"`
	Query       any    `json:"query,omitempty"`
	Datasource  any    `json:"datasource,omitempty"`
	Regex       string `json:"regex,omitempty"`
	Sort        int    `json:"sort,omitempty"`
	Refresh     int    `json:"refresh,omitempty"`
	Multi       bool   `json:"multi,omitempty"`
	IncludeAll  bool   `json:"includeAll,omitempty"`
	AllValue    string `json:"allValue,omitempty"`
	Current     any    `json:"current,omitempty"`
	Hide        int    `json:"hide,omitempty"`
}

// GrafanaPanel represents a panel in the dashboard JSON.
type GrafanaPanel struct {
	ID          int                    `json:"id"`
	Type        string                 `json:"type"`
	Title       string                 `json:"title"`
	Description string                 `json:"description,omitempty"`
	Datasource  any                    `json:"datasource,omitempty"`
	GridPos     GridPos                `json:"gridPos"`
	Targets     []GrafanaTarget        `json:"targets,omitempty"`
	Options     map[string]any         `json:"options,omitempty"`
	FieldConfig *GrafanaFieldConfig    `json:"fieldConfig,omitempty"`
	Collapsed   bool                   `json:"collapsed,omitempty"`
	Panels      []GrafanaPanel         `json:"panels,omitempty"` // For collapsed rows
	Transparent bool                   `json:"transparent,omitempty"`
	Extra       map[string]any         `json:"-"` // Capture unknown fields
}

// GridPos represents panel position.
type GridPos struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

// GrafanaTarget represents a query target.
type GrafanaTarget struct {
	RefID        string `json:"refId"`
	Expr         string `json:"expr,omitempty"`
	LegendFormat string `json:"legendFormat,omitempty"`
	Interval     string `json:"interval,omitempty"`
	Instant      bool   `json:"instant,omitempty"`
	Range        bool   `json:"range,omitempty"`
	Hide         bool   `json:"hide,omitempty"`
	Datasource   any    `json:"datasource,omitempty"`
}

// GrafanaFieldConfig represents field configuration.
type GrafanaFieldConfig struct {
	Defaults  *FieldDefaults  `json:"defaults,omitempty"`
	Overrides []FieldOverride `json:"overrides,omitempty"`
}

// FieldDefaults represents default field settings.
type FieldDefaults struct {
	Unit       string          `json:"unit,omitempty"`
	Min        *float64        `json:"min,omitempty"`
	Max        *float64        `json:"max,omitempty"`
	Decimals   *int            `json:"decimals,omitempty"`
	Color      map[string]any  `json:"color,omitempty"`
	Thresholds map[string]any  `json:"thresholds,omitempty"`
	Custom     map[string]any  `json:"custom,omitempty"`
}

// FieldOverride represents a field override.
type FieldOverride struct {
	Matcher    map[string]any `json:"matcher"`
	Properties []any          `json:"properties"`
}

// ParseGrafanaDashboard parses a Grafana dashboard JSON file.
func ParseGrafanaDashboard(path string) (*GrafanaDashboard, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return ParseGrafanaDashboardFromBytes(data)
}

// ParseGrafanaDashboardFromBytes parses a Grafana dashboard from JSON bytes.
func ParseGrafanaDashboardFromBytes(data []byte) (*GrafanaDashboard, error) {
	var dashboard GrafanaDashboard
	if err := json.Unmarshal(data, &dashboard); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &dashboard, nil
}

// ValidateGrafanaDashboard validates a parsed Grafana dashboard.
func ValidateGrafanaDashboard(dashboard *GrafanaDashboard) []string {
	var warnings []string

	if dashboard.Title == "" {
		warnings = append(warnings, "dashboard has no title")
	}

	if dashboard.UID == "" {
		warnings = append(warnings, "dashboard has no UID (one will be generated)")
	}

	panelCount := 0
	rowCount := 0
	for _, panel := range dashboard.Panels {
		if panel.Type == "row" {
			rowCount++
			panelCount += len(panel.Panels) // Collapsed row panels
		} else {
			panelCount++
		}
	}

	if panelCount == 0 {
		warnings = append(warnings, "dashboard has no panels")
	}

	if rowCount == 0 && panelCount > 0 {
		warnings = append(warnings, "dashboard has no rows; panels will be grouped into a default row")
	}

	// Check for unsupported panel types
	supportedTypes := map[string]bool{
		"row":        true,
		"timeseries": true,
		"stat":       true,
		"table":      true,
		"gauge":      true,
		"bargauge":   true,
		"piechart":   true,
		"text":       true,
		"logs":       true,
		"heatmap":    true,
		"graph":      true, // Legacy, will be converted to timeseries
	}

	for _, panel := range dashboard.Panels {
		if panel.Type == "row" {
			for _, nested := range panel.Panels {
				if !supportedTypes[nested.Type] {
					warnings = append(warnings, fmt.Sprintf("unsupported panel type %q (panel %q)", nested.Type, nested.Title))
				}
			}
		} else if !supportedTypes[panel.Type] {
			warnings = append(warnings, fmt.Sprintf("unsupported panel type %q (panel %q)", panel.Type, panel.Title))
		}
	}

	return warnings
}

// ConvertToWetwire converts a GrafanaDashboard to a wetwire grafana.Dashboard.
func ConvertToWetwire(gd *GrafanaDashboard) *grafana.Dashboard {
	d := grafana.NewDashboard(gd.UID, gd.Title)

	if gd.Description != "" {
		d.WithDescription(gd.Description)
	}
	if len(gd.Tags) > 0 {
		d.WithTags(gd.Tags...)
	}
	if gd.Time != nil {
		d.WithTime(gd.Time.From, gd.Time.To)
	}
	if gd.Refresh != "" {
		d.WithRefresh(gd.Refresh)
	}
	if gd.Timezone != "" {
		d.WithTimezone(gd.Timezone)
	}
	if gd.Editable {
		d.Editable()
	}
	d.SchemaVersion = gd.SchemaVersion
	d.Version = gd.Version

	// Convert panels to rows
	d.Rows = convertPanelsToRows(gd.Panels)

	// Convert variables
	if gd.Templating != nil {
		for _, v := range gd.Templating.List {
			d.Variables = append(d.Variables, convertVariable(v))
		}
	}

	return d
}

// convertPanelsToRows converts flat panel array to row-based structure.
func convertPanelsToRows(panels []GrafanaPanel) []*grafana.Row {
	var rows []*grafana.Row
	var currentRow *grafana.Row

	for _, panel := range panels {
		if panel.Type == "row" {
			// Start a new row
			currentRow = grafana.NewRow(panel.Title)
			if panel.Collapsed {
				currentRow.Collapsed()
			}
			rows = append(rows, currentRow)

			// Add collapsed row panels
			if panel.Collapsed && len(panel.Panels) > 0 {
				for _, nested := range panel.Panels {
					currentRow.AddPanel(convertPanel(nested))
				}
			}
		} else {
			// Regular panel
			if currentRow == nil {
				// Create default row if no row exists
				currentRow = grafana.NewRow("Dashboard")
				rows = append(rows, currentRow)
			}
			currentRow.AddPanel(convertPanel(panel))
		}
	}

	// If no panels were added and no rows, return empty
	if len(rows) == 0 && len(panels) > 0 {
		// All panels without row markers - create a single row
		row := grafana.NewRow("Dashboard")
		for _, panel := range panels {
			row.AddPanel(convertPanel(panel))
		}
		rows = append(rows, row)
	}

	return rows
}

// convertPanel converts a GrafanaPanel to a wetwire panel type.
func convertPanel(panel GrafanaPanel) any {
	switch panel.Type {
	case "timeseries", "graph":
		p := grafana.TimeSeries(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "stat":
		p := grafana.Stat(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "table":
		p := grafana.Table(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "gauge":
		p := grafana.Gauge(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "bargauge":
		p := grafana.BarGauge(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "piechart":
		p := grafana.PieChart(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "text":
		p := grafana.Text(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "logs":
		p := grafana.Logs(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	case "heatmap":
		p := grafana.Heatmap(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	default:
		// Unknown panel type - use timeseries as fallback
		p := grafana.TimeSeries(panel.Title)
		applyBasePanel(&p.BasePanel, panel)
		return p
	}
}

// applyBasePanel applies common panel fields.
func applyBasePanel(base *grafana.BasePanel, panel GrafanaPanel) {
	base.ID = panel.ID
	base.Description = panel.Description
	base.GridPos = grafana.GridPos{
		X: panel.GridPos.X,
		Y: panel.GridPos.Y,
		W: panel.GridPos.W,
		H: panel.GridPos.H,
	}
	base.Transparent = panel.Transparent

	// Convert datasource
	if ds, ok := panel.Datasource.(string); ok {
		base.Datasource = ds
	} else if dsMap, ok := panel.Datasource.(map[string]any); ok {
		if uid, ok := dsMap["uid"].(string); ok {
			base.Datasource = uid
		}
	}

	// Convert targets
	for _, t := range panel.Targets {
		target := &grafana.PrometheusTarget{
			RefID:        t.RefID,
			Expr:         t.Expr,
			LegendFormat: t.LegendFormat,
			Interval:     t.Interval,
			IsInstant:    t.Instant,
			IsRange:      t.Range,
			Hidden:       t.Hide,
		}
		base.Targets = append(base.Targets, target)
	}

	// Convert field config
	if panel.FieldConfig != nil && panel.FieldConfig.Defaults != nil {
		base.FieldConfig.Defaults.Unit = panel.FieldConfig.Defaults.Unit
		base.FieldConfig.Defaults.Min = panel.FieldConfig.Defaults.Min
		base.FieldConfig.Defaults.Max = panel.FieldConfig.Defaults.Max
		base.FieldConfig.Defaults.Decimals = panel.FieldConfig.Defaults.Decimals
	}
}

// convertVariable converts a GrafanaVariable to a wetwire variable.
func convertVariable(v GrafanaVariable) *grafana.Variable {
	var variable *grafana.Variable

	// Extract query as string
	queryStr := ""
	switch q := v.Query.(type) {
	case string:
		queryStr = q
	case map[string]any:
		if query, ok := q["query"].(string); ok {
			queryStr = query
		}
	}

	switch v.Type {
	case "query":
		variable = grafana.QueryVar(v.Name, queryStr)
	case "custom":
		variable = grafana.CustomVar(v.Name, queryStr)
	case "interval":
		variable = grafana.IntervalVar(v.Name, queryStr)
	case "datasource":
		variable = grafana.DatasourceVar(v.Name, queryStr)
	case "textbox":
		variable = grafana.TextboxVar(v.Name, queryStr)
	case "constant":
		variable = grafana.ConstantVar(v.Name, queryStr)
	default:
		variable = &grafana.Variable{Name: v.Name, Type: v.Type, Query: queryStr}
	}

	if v.Label != "" {
		variable.WithLabel(v.Label)
	}
	if v.Description != "" {
		variable.WithDescription(v.Description)
	}
	if v.Regex != "" {
		variable.WithRegex(v.Regex)
	}
	if v.Sort != 0 {
		variable.WithSort(v.Sort)
	}
	if v.Refresh != 0 {
		variable.WithRefresh(v.Refresh)
	}
	if v.Multi {
		variable.MultiSelect()
	}
	if v.IncludeAll {
		variable.IncludeAll()
	}
	if v.AllValue != "" {
		variable.WithAllValue(v.AllValue)
	}
	if v.Hide == grafana.HideVariable {
		variable.Hide()
	} else if v.Hide == grafana.HideLabelOnly {
		variable.HideLabel()
	}

	return variable
}
