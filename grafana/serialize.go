package grafana

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// dashboardOutput represents the Grafana dashboard JSON structure.
type dashboardOutput struct {
	UID           string            `json:"uid"`
	Title         string            `json:"title"`
	Description   string            `json:"description,omitempty"`
	Tags          []string          `json:"tags,omitempty"`
	Time          *TimeRange        `json:"time,omitempty"`
	Refresh       string            `json:"refresh,omitempty"`
	Timezone      string            `json:"timezone,omitempty"`
	Editable      bool              `json:"editable"`
	Panels        []any             `json:"panels"`
	Templating    *templatingOutput `json:"templating,omitempty"`
	Annotations   *annotationsOutput `json:"annotations,omitempty"`
	SchemaVersion int               `json:"schemaVersion"`
	Version       int               `json:"version,omitempty"`
}

type templatingOutput struct {
	List []any `json:"list"`
}

type annotationsOutput struct {
	List []any `json:"list"`
}

// Serialize converts the Dashboard to JSON bytes.
func (d *Dashboard) Serialize() ([]byte, error) {
	// Assign panel IDs
	AssignPanelIDs(d)

	// Calculate grid positions (inline to avoid import cycle)
	calculateGridPositions(d)

	// Build the output structure
	output := dashboardOutput{
		UID:           d.UID,
		Title:         d.Title,
		Description:   d.Description,
		Tags:          d.Tags,
		Time:          d.Time,
		Refresh:       d.Refresh,
		Timezone:      d.Timezone,
		Editable:      d.IsEditable,
		Panels:        flattenPanels(d),
		SchemaVersion: d.SchemaVersion,
		Version:       d.Version,
	}

	// Add templating if there are variables
	if len(d.Variables) > 0 {
		output.Templating = &templatingOutput{
			List: d.Variables,
		}
	}

	// Add annotations
	if len(d.Annotations) > 0 {
		output.Annotations = &annotationsOutput{
			List: d.Annotations,
		}
	}

	return json.MarshalIndent(output, "", "  ")
}

// SerializeToFile writes the Dashboard to a JSON file.
func (d *Dashboard) SerializeToFile(path string) error {
	data, err := d.Serialize()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// MustSerialize converts the Dashboard to JSON bytes, panicking on error.
func (d *Dashboard) MustSerialize() []byte {
	data, err := d.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}

// flattenPanels extracts all panels from rows into a flat array.
// This is the format Grafana expects in the JSON.
func flattenPanels(d *Dashboard) []any {
	var panels []any
	panelIndex := 0

	for rowIndex, row := range d.Rows {
		if row == nil {
			continue
		}

		// Add row panel (Grafana uses row type panels for collapsible rows)
		rowPanel := map[string]any{
			"id":        1000 + rowIndex, // Row IDs start at 1000
			"type":      "row",
			"title":     row.Title,
			"collapsed": row.IsCollapsed,
			"gridPos": map[string]int{
				"x": 0,
				"y": getRowY(d, rowIndex),
				"w": 24,
				"h": 1,
			},
		}

		if row.IsCollapsed {
			// For collapsed rows, panels go inside the row
			rowPanels := make([]any, 0, len(row.Panels))
			rowPanels = append(rowPanels, row.Panels...)
			rowPanel["panels"] = rowPanels
		}

		panels = append(panels, rowPanel)

		// If not collapsed, add panels after the row
		if !row.IsCollapsed {
			for _, p := range row.Panels {
				panels = append(panels, p)
				panelIndex++
			}
		}
	}

	return panels
}

// getRowY calculates the Y position for a row.
func getRowY(d *Dashboard, rowIndex int) int {
	y := 0
	for i := 0; i < rowIndex; i++ {
		row := d.Rows[i]
		if row == nil {
			continue
		}
		y += 1 // Row header height
		if !row.IsCollapsed {
			y += getRowContentHeight(row)
		}
	}
	return y
}

// getRowContentHeight calculates the height of a row's content.
func getRowContentHeight(row *Row) int {
	maxHeight := 0
	x := 0
	rowHeight := 0

	for _, p := range row.Panels {
		gridPos := getPanelGridPos(p)
		if gridPos == nil {
			continue
		}

		w := gridPos.W
		if w == 0 {
			w = 12
		}
		h := gridPos.H
		if h == 0 {
			h = 8
		}

		if x+w > 24 {
			rowHeight += maxHeight
			maxHeight = h
			x = w
		} else {
			if h > maxHeight {
				maxHeight = h
			}
			x += w
		}
	}

	return rowHeight + maxHeight
}

// getPanelGridPos returns the GridPos for any panel type.
func getPanelGridPos(p any) *GridPos {
	switch panel := p.(type) {
	case *TimeSeriesPanel:
		return &panel.GridPos
	case *StatPanel:
		return &panel.GridPos
	case *TablePanel:
		return &panel.GridPos
	default:
		return nil
	}
}

// AssignPanelIDs assigns unique IDs to all panels in the dashboard.
// IDs are assigned deterministically based on panel order.
func AssignPanelIDs(d *Dashboard) {
	if d == nil {
		return
	}

	id := 1
	for _, row := range d.Rows {
		if row == nil {
			continue
		}
		for _, p := range row.Panels {
			setPanelID(p, id)
			id++
		}
	}
}

// setPanelID sets the ID for any panel type.
func setPanelID(p any, id int) {
	switch panel := p.(type) {
	case *TimeSeriesPanel:
		panel.ID = id
	case *StatPanel:
		panel.ID = id
	case *TablePanel:
		panel.ID = id
	}
}

const (
	gridColumns        = 24
	defaultPanelWidth  = 12
	defaultPanelHeight = 8
	rowHeaderHeight    = 1
)

// calculateGridPositions calculates grid positions for all panels.
func calculateGridPositions(d *Dashboard) {
	if d == nil {
		return
	}

	currentY := 0

	for _, row := range d.Rows {
		if row == nil {
			continue
		}

		// Row header takes one row
		currentY += rowHeaderHeight

		if row.IsCollapsed {
			continue
		}

		// Layout panels in this row
		currentX := 0
		rowMaxHeight := 0

		for _, panelAny := range row.Panels {
			if panelAny == nil {
				continue
			}

			gridPos := getPanelGridPos(panelAny)
			if gridPos == nil {
				continue
			}

			width := gridPos.W
			height := gridPos.H

			if width == 0 {
				width = defaultPanelWidth
			}
			if height == 0 {
				height = defaultPanelHeight
			}

			// Wrap to next line if needed
			if currentX+width > gridColumns {
				currentX = 0
				currentY += rowMaxHeight
				rowMaxHeight = 0
			}

			// Set position
			setGridPos(panelAny, currentX, currentY, width, height)

			currentX += width
			if height > rowMaxHeight {
				rowMaxHeight = height
			}
		}

		if len(row.Panels) > 0 && rowMaxHeight > 0 {
			currentY += rowMaxHeight
		}
	}
}

// setGridPos sets the grid position for a panel.
func setGridPos(p any, x, y, w, h int) {
	switch panel := p.(type) {
	case *TimeSeriesPanel:
		panel.GridPos = GridPos{X: x, Y: y, W: w, H: h}
	case *StatPanel:
		panel.GridPos = GridPos{X: x, Y: y, W: w, H: h}
	case *TablePanel:
		panel.GridPos = GridPos{X: x, Y: y, W: w, H: h}
	}
}
