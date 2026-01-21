package importer

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
)

// GenerateGrafanaGoCode generates Go source code from a GrafanaDashboard.
func GenerateGrafanaGoCode(dashboard *GrafanaDashboard, packageName string) ([]byte, error) {
	gen := &grafanaCodeGenerator{
		packageName: packageName,
		dashboard:   dashboard,
		varCounter:  0,
	}
	return gen.generate()
}

type grafanaCodeGenerator struct {
	packageName string
	dashboard   *GrafanaDashboard
	varCounter  int
}

func (g *grafanaCodeGenerator) generate() ([]byte, error) {
	var buf bytes.Buffer

	// Write package and imports
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"github.com/lex00/wetwire-observability-go/grafana\"\n")
	buf.WriteString(")\n\n")

	// Generate variables if present
	if g.dashboard.Templating != nil && len(g.dashboard.Templating.List) > 0 {
		for _, v := range g.dashboard.Templating.List {
			if err := g.writeVariable(&buf, v); err != nil {
				return nil, err
			}
		}
	}

	// Generate panels by row
	panelVars, rowVars := g.collectPanelAndRowVars()

	// Write panel variables
	for _, pv := range panelVars {
		if err := g.writePanel(&buf, pv.panel, pv.varName); err != nil {
			return nil, err
		}
	}

	// Write row variables
	for _, rv := range rowVars {
		if err := g.writeRow(&buf, rv.row, rv.varName, rv.panelVarNames); err != nil {
			return nil, err
		}
	}

	// Generate the main dashboard
	if err := g.writeDashboard(&buf, rowVars); err != nil {
		return nil, err
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Return unformatted code with error for debugging
		return buf.Bytes(), fmt.Errorf("failed to format generated code: %w", err)
	}

	return formatted, nil
}

type panelVar struct {
	panel   GrafanaPanel
	varName string
}

type rowVar struct {
	row           GrafanaPanel
	varName       string
	panelVarNames []string
}

func (g *grafanaCodeGenerator) collectPanelAndRowVars() ([]panelVar, []rowVar) {
	var panelVars []panelVar
	var rowVars []rowVar

	panelIndex := 0
	rowIndex := 0

	var currentRowPanels []string

	for _, panel := range g.dashboard.Panels {
		if panel.Type == "row" {
			// If we have accumulated panels from before this row, create a default row
			if len(currentRowPanels) > 0 && rowIndex == 0 {
				rv := rowVar{
					row:           GrafanaPanel{Type: "row", Title: "Dashboard"},
					varName:       "DefaultRow",
					panelVarNames: currentRowPanels,
				}
				rowVars = append(rowVars, rv)
				rowIndex++
			}

			// Process panels in collapsed row
			var rowPanelVars []string
			if panel.Collapsed {
				for _, nested := range panel.Panels {
					varName := g.getPanelVarName(nested, panelIndex)
					panelVars = append(panelVars, panelVar{panel: nested, varName: varName})
					rowPanelVars = append(rowPanelVars, varName)
					panelIndex++
				}
			}

			rowVarName := g.getRowVarName(panel, rowIndex)
			rowVars = append(rowVars, rowVar{
				row:           panel,
				varName:       rowVarName,
				panelVarNames: rowPanelVars,
			})
			rowIndex++
			currentRowPanels = nil
		} else {
			// Regular panel
			varName := g.getPanelVarName(panel, panelIndex)
			panelVars = append(panelVars, panelVar{panel: panel, varName: varName})
			currentRowPanels = append(currentRowPanels, varName)
			panelIndex++
		}
	}

	// Handle remaining panels not in a row
	if len(currentRowPanels) > 0 {
		if len(rowVars) == 0 {
			// No rows at all - create a single default row
			rv := rowVar{
				row:           GrafanaPanel{Type: "row", Title: "Dashboard"},
				varName:       "DefaultRow",
				panelVarNames: currentRowPanels,
			}
			rowVars = append(rowVars, rv)
		} else {
			// Add to the last row
			rowVars[len(rowVars)-1].panelVarNames = append(rowVars[len(rowVars)-1].panelVarNames, currentRowPanels...)
		}
	}

	return panelVars, rowVars
}

func (g *grafanaCodeGenerator) getPanelVarName(panel GrafanaPanel, index int) string {
	if panel.Title != "" {
		return g.sanitizeVarName(panel.Title) + "Panel"
	}
	return fmt.Sprintf("Panel%d", index)
}

func (g *grafanaCodeGenerator) getRowVarName(row GrafanaPanel, index int) string {
	if row.Title != "" {
		return g.sanitizeVarName(row.Title) + "Row"
	}
	return fmt.Sprintf("Row%d", index)
}

func (g *grafanaCodeGenerator) writeVariable(buf *bytes.Buffer, v GrafanaVariable) error {
	varName := g.sanitizeVarName(v.Name) + "Var"

	// Get query as string
	queryStr := ""
	switch q := v.Query.(type) {
	case string:
		queryStr = q
	case map[string]any:
		if query, ok := q["query"].(string); ok {
			queryStr = query
		}
	}

	buf.WriteString(fmt.Sprintf("// %s is a %s variable.\n", varName, v.Type))

	switch v.Type {
	case "query":
		buf.WriteString(fmt.Sprintf("var %s = grafana.QueryVar(%q, %q)", varName, v.Name, queryStr))
	case "custom":
		buf.WriteString(fmt.Sprintf("var %s = grafana.CustomVar(%q, %q)", varName, v.Name, queryStr))
	case "interval":
		buf.WriteString(fmt.Sprintf("var %s = grafana.IntervalVar(%q, %q)", varName, v.Name, queryStr))
	case "datasource":
		buf.WriteString(fmt.Sprintf("var %s = grafana.DatasourceVar(%q, %q)", varName, v.Name, queryStr))
	case "textbox":
		buf.WriteString(fmt.Sprintf("var %s = grafana.TextboxVar(%q, %q)", varName, v.Name, queryStr))
	case "constant":
		buf.WriteString(fmt.Sprintf("var %s = grafana.ConstantVar(%q, %q)", varName, v.Name, queryStr))
	default:
		buf.WriteString(fmt.Sprintf("var %s = &grafana.Variable{Name: %q, Type: %q, Query: %q}", varName, v.Name, v.Type, queryStr))
	}

	// Add method chains
	if v.Label != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithLabel(%q)", v.Label))
	}
	if v.Description != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithDescription(%q)", v.Description))
	}
	if v.Regex != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithRegex(%q)", v.Regex))
	}
	if v.Sort != 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithSort(%d)", v.Sort))
	}
	if v.Refresh != 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithRefresh(%d)", v.Refresh))
	}
	if v.Multi {
		buf.WriteString(".\n\tMultiSelect()")
	}
	if v.IncludeAll {
		buf.WriteString(".\n\tIncludeAll()")
	}
	if v.AllValue != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithAllValue(%q)", v.AllValue))
	}
	if v.Hide == 2 {
		buf.WriteString(".\n\tHide()")
	} else if v.Hide == 1 {
		buf.WriteString(".\n\tHideLabel()")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *grafanaCodeGenerator) writePanel(buf *bytes.Buffer, panel GrafanaPanel, varName string) error {
	buf.WriteString(fmt.Sprintf("// %s is a %s panel.\n", varName, panel.Type))

	// Write panel constructor
	switch panel.Type {
	case "timeseries", "graph":
		buf.WriteString(fmt.Sprintf("var %s = grafana.TimeSeries(%q)", varName, panel.Title))
	case "stat":
		buf.WriteString(fmt.Sprintf("var %s = grafana.Stat(%q)", varName, panel.Title))
	case "table":
		buf.WriteString(fmt.Sprintf("var %s = grafana.Table(%q)", varName, panel.Title))
	case "gauge":
		buf.WriteString(fmt.Sprintf("var %s = grafana.Gauge(%q)", varName, panel.Title))
	case "bargauge":
		buf.WriteString(fmt.Sprintf("var %s = grafana.BarGauge(%q)", varName, panel.Title))
	case "piechart":
		buf.WriteString(fmt.Sprintf("var %s = grafana.PieChart(%q)", varName, panel.Title))
	case "text":
		buf.WriteString(fmt.Sprintf("var %s = grafana.Text(%q)", varName, panel.Title))
		// Add content via method chain if present
		if panel.Options != nil {
			if content, ok := panel.Options["content"].(string); ok && content != "" {
				buf.WriteString(fmt.Sprintf(".\n\tWithContent(%q)", content))
			}
		}
	case "logs":
		buf.WriteString(fmt.Sprintf("var %s = grafana.Logs(%q)", varName, panel.Title))
	case "heatmap":
		buf.WriteString(fmt.Sprintf("var %s = grafana.Heatmap(%q)", varName, panel.Title))
	default:
		buf.WriteString(fmt.Sprintf("var %s = grafana.TimeSeries(%q)", varName, panel.Title))
	}

	// Add method chains
	if panel.Description != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithDescription(%q)", panel.Description))
	}

	// Datasource
	if ds, ok := panel.Datasource.(string); ok && ds != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithDatasource(%q)", ds))
	} else if dsMap, ok := panel.Datasource.(map[string]any); ok {
		if uid, ok := dsMap["uid"].(string); ok {
			buf.WriteString(fmt.Sprintf(".\n\tWithDatasource(%q)", uid))
		}
	}

	// Size
	if panel.GridPos.W > 0 && panel.GridPos.H > 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithSize(%d, %d)", panel.GridPos.W, panel.GridPos.H))
	}

	// Transparent
	if panel.Transparent {
		buf.WriteString(".\n\tTransparent()")
	}

	// Field config - unit
	if panel.FieldConfig != nil && panel.FieldConfig.Defaults != nil {
		if panel.FieldConfig.Defaults.Unit != "" {
			buf.WriteString(fmt.Sprintf(".\n\tWithUnit(%q)", panel.FieldConfig.Defaults.Unit))
		}
		if panel.FieldConfig.Defaults.Min != nil {
			buf.WriteString(fmt.Sprintf(".\n\tWithMin(%v)", *panel.FieldConfig.Defaults.Min))
		}
		if panel.FieldConfig.Defaults.Max != nil {
			buf.WriteString(fmt.Sprintf(".\n\tWithMax(%v)", *panel.FieldConfig.Defaults.Max))
		}
		if panel.FieldConfig.Defaults.Decimals != nil {
			buf.WriteString(fmt.Sprintf(".\n\tWithDecimals(%d)", *panel.FieldConfig.Defaults.Decimals))
		}
	}

	buf.WriteString("\n\n")

	// Write targets separately if present
	if len(panel.Targets) > 0 {
		for i, t := range panel.Targets {
			targetVarName := fmt.Sprintf("%sTarget%d", varName, i)
			g.writeTarget(buf, t, targetVarName)
		}

		// Add targets to panel using init
		buf.WriteString("func init() {\n")
		for i := range panel.Targets {
			targetVarName := fmt.Sprintf("%sTarget%d", varName, i)
			buf.WriteString(fmt.Sprintf("\t%s.AddTarget(%s)\n", varName, targetVarName))
		}
		buf.WriteString("}\n\n")
	}

	return nil
}

func (g *grafanaCodeGenerator) writeTarget(buf *bytes.Buffer, target GrafanaTarget, varName string) {
	buf.WriteString(fmt.Sprintf("var %s = grafana.PromTarget(%q)", varName, target.Expr))

	if target.RefID != "" && target.RefID != "A" {
		buf.WriteString(fmt.Sprintf(".\n\tWithRefID(%q)", target.RefID))
	}
	if target.LegendFormat != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithLegendFormat(%q)", target.LegendFormat))
	}
	if target.Interval != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithInterval(%q)", target.Interval))
	}
	if target.Instant {
		buf.WriteString(".\n\tInstant()")
	}
	if target.Range {
		buf.WriteString(".\n\tRange()")
	}
	if target.Hide {
		buf.WriteString(".\n\tHide()")
	}

	buf.WriteString("\n\n")
}

func (g *grafanaCodeGenerator) writeRow(buf *bytes.Buffer, row GrafanaPanel, varName string, panelVarNames []string) error {
	buf.WriteString(fmt.Sprintf("// %s groups related panels.\n", varName))
	buf.WriteString(fmt.Sprintf("var %s = grafana.NewRow(%q)", varName, row.Title))

	if row.Collapsed {
		buf.WriteString(".\n\tCollapsed()")
	}

	if len(panelVarNames) > 0 {
		buf.WriteString(".\n\tWithPanels(\n")
		for _, pv := range panelVarNames {
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", pv))
		}
		buf.WriteString("\t)")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *grafanaCodeGenerator) writeDashboard(buf *bytes.Buffer, rowVars []rowVar) error {
	buf.WriteString("// Dashboard is the main Grafana dashboard.\n")
	buf.WriteString(fmt.Sprintf("var Dashboard = grafana.NewDashboard(%q, %q)", g.dashboard.UID, g.dashboard.Title))

	if g.dashboard.Description != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithDescription(%q)", g.dashboard.Description))
	}

	if len(g.dashboard.Tags) > 0 {
		tags := make([]string, len(g.dashboard.Tags))
		for i, t := range g.dashboard.Tags {
			tags[i] = fmt.Sprintf("%q", t)
		}
		buf.WriteString(fmt.Sprintf(".\n\tWithTags(%s)", strings.Join(tags, ", ")))
	}

	if g.dashboard.Time != nil {
		buf.WriteString(fmt.Sprintf(".\n\tWithTime(%q, %q)", g.dashboard.Time.From, g.dashboard.Time.To))
	}

	if g.dashboard.Refresh != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithRefresh(%q)", g.dashboard.Refresh))
	}

	if g.dashboard.Timezone != "" && g.dashboard.Timezone != "browser" {
		buf.WriteString(fmt.Sprintf(".\n\tWithTimezone(%q)", g.dashboard.Timezone))
	}

	if g.dashboard.Editable {
		buf.WriteString(".\n\tEditable()")
	} else {
		buf.WriteString(".\n\tReadOnly()")
	}

	// Add rows
	if len(rowVars) > 0 {
		buf.WriteString(".\n\tWithRows(\n")
		for _, rv := range rowVars {
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", rv.varName))
		}
		buf.WriteString("\t)")
	}

	// Add variables
	if g.dashboard.Templating != nil && len(g.dashboard.Templating.List) > 0 {
		buf.WriteString(".\n\tWithVariables(\n")
		for _, v := range g.dashboard.Templating.List {
			varName := g.sanitizeVarName(v.Name) + "Var"
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t)")
	}

	buf.WriteString("\n")
	return nil
}

func (g *grafanaCodeGenerator) sanitizeVarName(name string) string {
	result := strings.Builder{}
	capitalize := true

	for _, c := range name {
		if c == '-' || c == '_' || c == '.' || c == '/' || c == ' ' || c == ':' {
			capitalize = true
			continue
		}
		if capitalize {
			result.WriteRune(toUpperRune(c))
			capitalize = false
		} else {
			result.WriteRune(c)
		}
	}

	s := result.String()
	if s == "" {
		return "Dashboard"
	}
	return s
}
