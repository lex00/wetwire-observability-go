package grafana

// Row represents a row in a Grafana dashboard for logical panel grouping.
type Row struct {
	// Title is the display title of the row.
	Title string `json:"title"`

	// Panels contains the panels in this row.
	// Use any to allow different panel types.
	Panels []any `json:"panels,omitempty"`

	// IsCollapsed indicates whether the row is collapsed by default.
	IsCollapsed bool `json:"collapsed,omitempty"`

	// Height is the row height in pixels (optional).
	Height int `json:"height,omitempty"`
}

// NewRow creates a new Row with the given title.
func NewRow(title string) *Row {
	return &Row{Title: title}
}

// WithPanels sets the panels in this row.
func (r *Row) WithPanels(panels ...any) *Row {
	r.Panels = panels
	return r
}

// AddPanel adds a panel to this row.
func (r *Row) AddPanel(panel any) *Row {
	r.Panels = append(r.Panels, panel)
	return r
}

// Collapsed marks this row as collapsed by default.
func (r *Row) Collapsed() *Row {
	r.IsCollapsed = true
	return r
}

// Expanded marks this row as expanded by default.
func (r *Row) Expanded() *Row {
	r.IsCollapsed = false
	return r
}

// WithHeight sets the row height in pixels.
func (r *Row) WithHeight(height int) *Row {
	r.Height = height
	return r
}
