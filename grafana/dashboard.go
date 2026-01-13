package grafana

// Dashboard represents a Grafana dashboard.
type Dashboard struct {
	// UID is the unique identifier for the dashboard.
	UID string `json:"uid"`

	// Title is the display title of the dashboard.
	Title string `json:"title"`

	// Description is an optional description of the dashboard.
	Description string `json:"description,omitempty"`

	// Tags are labels for organizing and filtering dashboards.
	Tags []string `json:"tags,omitempty"`

	// Time is the default time range for the dashboard.
	Time *TimeRange `json:"time,omitempty"`

	// Refresh is the auto-refresh interval (e.g., "30s", "1m", "5m").
	Refresh string `json:"refresh,omitempty"`

	// Timezone is the dashboard timezone (e.g., "browser", "utc", "America/New_York").
	Timezone string `json:"timezone,omitempty"`

	// IsEditable indicates whether the dashboard can be edited.
	IsEditable bool `json:"editable,omitempty"`

	// Rows contains the dashboard rows for panel organization.
	Rows []*Row `json:"rows,omitempty"`

	// Variables contains template variables for the dashboard.
	Variables []any `json:"templating,omitempty"`

	// Annotations contains annotation settings.
	Annotations []any `json:"annotations,omitempty"`

	// Links contains dashboard links.
	Links []any `json:"links,omitempty"`

	// Version is the dashboard version (incremented on save).
	Version int `json:"version,omitempty"`

	// SchemaVersion is the dashboard schema version.
	SchemaVersion int `json:"schemaVersion,omitempty"`
}

// NewDashboard creates a new Dashboard with the given uid and title.
func NewDashboard(uid, title string) *Dashboard {
	return &Dashboard{
		UID:           uid,
		Title:         title,
		SchemaVersion: 39, // Current Grafana schema version
	}
}

// WithDescription sets the dashboard description.
func (d *Dashboard) WithDescription(description string) *Dashboard {
	d.Description = description
	return d
}

// WithTags sets the dashboard tags.
func (d *Dashboard) WithTags(tags ...string) *Dashboard {
	d.Tags = tags
	return d
}

// WithTime sets the default time range.
func (d *Dashboard) WithTime(from, to string) *Dashboard {
	d.Time = NewTimeRange(from, to)
	return d
}

// WithTimeRange sets the default time range from a TimeRange.
func (d *Dashboard) WithTimeRange(tr *TimeRange) *Dashboard {
	d.Time = tr
	return d
}

// WithRefresh sets the auto-refresh interval.
func (d *Dashboard) WithRefresh(refresh string) *Dashboard {
	d.Refresh = refresh
	return d
}

// WithTimezone sets the dashboard timezone.
func (d *Dashboard) WithTimezone(timezone string) *Dashboard {
	d.Timezone = timezone
	return d
}

// Editable marks the dashboard as editable.
func (d *Dashboard) Editable() *Dashboard {
	d.IsEditable = true
	return d
}

// ReadOnly marks the dashboard as read-only.
func (d *Dashboard) ReadOnly() *Dashboard {
	d.IsEditable = false
	return d
}

// WithRows sets the dashboard rows.
func (d *Dashboard) WithRows(rows ...*Row) *Dashboard {
	d.Rows = rows
	return d
}

// AddRow adds a row to the dashboard.
func (d *Dashboard) AddRow(row *Row) *Dashboard {
	d.Rows = append(d.Rows, row)
	return d
}

// WithVariables sets the template variables.
func (d *Dashboard) WithVariables(variables ...any) *Dashboard {
	d.Variables = variables
	return d
}

// AddVariable adds a template variable.
func (d *Dashboard) AddVariable(variable any) *Dashboard {
	d.Variables = append(d.Variables, variable)
	return d
}
