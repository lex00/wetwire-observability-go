// Package layout provides grid position calculation for Grafana dashboards.
package layout

import (
	"github.com/lex00/wetwire-observability-go/grafana"
)

const (
	// GridColumns is the number of columns in the Grafana grid.
	GridColumns = 24

	// DefaultPanelWidth is the default panel width.
	DefaultPanelWidth = 12

	// DefaultPanelHeight is the default panel height.
	DefaultPanelHeight = 8

	// RowHeaderHeight is the height of a row header.
	RowHeaderHeight = 1
)

// Width helper functions for common panel sizes.

// FullWidth returns the full grid width (24).
func FullWidth() int {
	return GridColumns
}

// HalfWidth returns half the grid width (12).
func HalfWidth() int {
	return GridColumns / 2
}

// ThirdWidth returns a third of the grid width (8).
func ThirdWidth() int {
	return GridColumns / 3
}

// QuarterWidth returns a quarter of the grid width (6).
func QuarterWidth() int {
	return GridColumns / 4
}

// CalculateGridPositions calculates grid positions for all panels in a dashboard.
// Panels are laid out from top to bottom, left to right, with automatic wrapping.
func CalculateGridPositions(dashboard *grafana.Dashboard) error {
	if dashboard == nil {
		return nil
	}

	currentY := 0

	for _, row := range dashboard.Rows {
		if row == nil {
			continue
		}

		// Row header takes one row
		rowStartY := currentY
		currentY += RowHeaderHeight

		if row.IsCollapsed {
			// Collapsed rows just take the header height
			continue
		}

		// Layout panels in this row
		currentX := 0
		rowMaxHeight := 0

		for _, panelAny := range row.Panels {
			if panelAny == nil {
				continue
			}

			// Get the panel's grid position (works with interface)
			gridPos := getPanelGridPos(panelAny)
			if gridPos == nil {
				continue
			}

			width := gridPos.W
			height := gridPos.H

			// Set defaults if not specified
			if width == 0 {
				width = DefaultPanelWidth
			}
			if height == 0 {
				height = DefaultPanelHeight
			}

			// Check if we need to wrap to next line
			if currentX+width > GridColumns {
				currentX = 0
				currentY += rowMaxHeight
				rowMaxHeight = 0
			}

			// Set position
			setPanelPosition(panelAny, currentX, currentY, width, height)

			// Advance X position
			currentX += width

			// Track max height for this visual row
			if height > rowMaxHeight {
				rowMaxHeight = height
			}
		}

		// Move Y to after the tallest panel in this row
		if len(row.Panels) > 0 && rowMaxHeight > 0 {
			currentY += rowMaxHeight
		}

		_ = rowStartY // Unused for now, but could be used for row positioning
	}

	return nil
}

// getPanelGridPos returns a pointer to the panel's GridPos field.
func getPanelGridPos(panel any) *grafana.GridPos {
	switch p := panel.(type) {
	case *grafana.TimeSeriesPanel:
		return &p.GridPos
	case *grafana.StatPanel:
		return &p.GridPos
	case *grafana.TablePanel:
		return &p.GridPos
	default:
		return nil
	}
}

// setPanelPosition sets the panel's grid position.
func setPanelPosition(panel any, x, y, w, h int) {
	switch p := panel.(type) {
	case *grafana.TimeSeriesPanel:
		p.GridPos = grafana.GridPos{X: x, Y: y, W: w, H: h}
	case *grafana.StatPanel:
		p.GridPos = grafana.GridPos{X: x, Y: y, W: w, H: h}
	case *grafana.TablePanel:
		p.GridPos = grafana.GridPos{X: x, Y: y, W: w, H: h}
	}
}
