package grafana

// TablePanel represents a Grafana table panel.
type TablePanel struct {
	BasePanel
	Options TableOptions `json:"options,omitempty"`
}

// TableOptions contains table panel options.
type TableOptions struct {
	ShowHeader       bool          `json:"showHeader,omitempty"`
	Footer           TableFooter   `json:"footer,omitempty"`
	SortBy           []TableSortBy `json:"sortBy,omitempty"`
	EnableFiltering  bool          `json:"enableFiltering,omitempty"`
	EnablePagination bool          `json:"enablePagination,omitempty"`
	FrameIndex       int           `json:"frameIndex,omitempty"`
}

// TableFooter contains table footer configuration.
type TableFooter struct {
	Show             bool     `json:"show,omitempty"`
	Reducer          []string `json:"reducer,omitempty"`
	Fields           []string `json:"fields,omitempty"`
	EnablePagination bool     `json:"enablePagination,omitempty"`
}

// TableSortBy represents a sort configuration.
type TableSortBy struct {
	DisplayName string `json:"displayName"`
	Desc        bool   `json:"desc,omitempty"`
}

// Table creates a new TablePanel.
func Table(title string) *TablePanel {
	return &TablePanel{
		BasePanel: BasePanel{
			Type:  "table",
			Title: title,
			GridPos: GridPos{
				W: 12,
				H: 8,
			},
		},
		Options: TableOptions{
			ShowHeader: true,
		},
	}
}

// WithDescription sets the panel description.
func (p *TablePanel) WithDescription(desc string) *TablePanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *TablePanel) WithDatasource(ds string) *TablePanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *TablePanel) WithSize(w, h int) *TablePanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *TablePanel) WithPosition(x, y int) *TablePanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *TablePanel) WithTargets(targets ...any) *TablePanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *TablePanel) AddTarget(target any) *TablePanel {
	p.Targets = append(p.Targets, target)
	return p
}

// ShowHeader shows the table header.
func (p *TablePanel) ShowHeader() *TablePanel {
	p.Options.ShowHeader = true
	return p
}

// HideHeader hides the table header.
func (p *TablePanel) HideHeader() *TablePanel {
	p.Options.ShowHeader = false
	return p
}

// WithFooter enables or disables the table footer.
func (p *TablePanel) WithFooter(show bool) *TablePanel {
	p.Options.Footer.Show = show
	return p
}

// WithFooterCalcs sets the footer calculations.
func (p *TablePanel) WithFooterCalcs(calcs ...string) *TablePanel {
	p.Options.Footer.Reducer = calcs
	return p
}

// SortByColumn sets the sort column.
func (p *TablePanel) SortByColumn(column string, descending bool) *TablePanel {
	p.Options.SortBy = []TableSortBy{
		{DisplayName: column, Desc: descending},
	}
	return p
}

// EnableFiltering enables column filtering.
func (p *TablePanel) EnableFiltering() *TablePanel {
	p.Options.EnableFiltering = true
	return p
}

// EnablePagination enables table pagination.
func (p *TablePanel) EnablePagination() *TablePanel {
	p.Options.EnablePagination = true
	return p
}

// WithColumnWidth sets a specific column width.
func (p *TablePanel) WithColumnWidth(column string, width int) *TablePanel {
	p.FieldConfig.Overrides = append(p.FieldConfig.Overrides, FieldOverride{
		Matcher: FieldMatcher{
			ID:      "byName",
			Options: column,
		},
		Properties: []any{
			map[string]any{
				"id":    "custom.width",
				"value": width,
			},
		},
	})
	return p
}

// HideColumn hides a specific column.
func (p *TablePanel) HideColumn(column string) *TablePanel {
	p.FieldConfig.Overrides = append(p.FieldConfig.Overrides, FieldOverride{
		Matcher: FieldMatcher{
			ID:      "byName",
			Options: column,
		},
		Properties: []any{
			map[string]any{
				"id":    "custom.hidden",
				"value": true,
			},
		},
	})
	return p
}

// WithUnit sets the display unit for all columns.
func (p *TablePanel) WithUnit(unit string) *TablePanel {
	p.FieldConfig.Defaults.Unit = unit
	return p
}

// WithDecimals sets the number of decimal places.
func (p *TablePanel) WithDecimals(decimals int) *TablePanel {
	p.FieldConfig.Defaults.Decimals = &decimals
	return p
}
