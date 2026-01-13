package grafana

import "strings"

// Sort options for query variables.
const (
	SortDisabled          = 0
	SortAlphabetical      = 1
	SortNumerical         = 3
	SortAlphabeticalDesc  = 2
	SortNumericalDesc     = 4
	SortAlphabeticalCI    = 5
	SortAlphabeticalCIAsc = 6
)

// Refresh options for query variables.
const (
	RefreshNever             = 0
	RefreshOnDashboardLoad   = 1
	RefreshOnTimeRangeChange = 2
)

// Hide options for variables.
const (
	HideNone      = 0
	HideLabelOnly = 1
	HideVariable  = 2
)

// Variable represents a Grafana template variable.
type Variable struct {
	// Name is the variable name (used in queries as $name).
	Name string `json:"name"`

	// Type is the variable type: query, custom, interval, datasource, textbox, constant.
	Type string `json:"type"`

	// Label is the display label (optional, defaults to Name).
	Label string `json:"label,omitempty"`

	// Description is an optional description of the variable.
	Description string `json:"description,omitempty"`

	// Query is the query expression or value(s) depending on type.
	Query string `json:"query,omitempty"`

	// Datasource is the data source for query variables.
	Datasource string `json:"datasource,omitempty"`

	// Regex is a regex filter for the results.
	Regex string `json:"regex,omitempty"`

	// Sort is the sort order for results.
	Sort int `json:"sort,omitempty"`

	// Refresh determines when to refresh the variable.
	Refresh int `json:"refresh,omitempty"`

	// Multi allows selecting multiple values.
	Multi bool `json:"multi,omitempty"`

	// IncludeAllOption adds an "All" option.
	IncludeAllOption bool `json:"includeAll,omitempty"`

	// AllValue is the custom value to use when "All" is selected.
	AllValue string `json:"allValue,omitempty"`

	// Current is the currently selected value.
	Current string `json:"current,omitempty"`

	// HideOption controls visibility: 0=show, 1=hide label, 2=hide variable.
	HideOption int `json:"hide,omitempty"`

	// Auto enables auto interval calculation (interval variables only).
	Auto bool `json:"auto,omitempty"`

	// AutoCount is the step count for auto interval calculation.
	AutoCount int `json:"auto_count,omitempty"`

	// AutoMin is the minimum interval for auto calculation.
	AutoMin string `json:"auto_min,omitempty"`

	// Options contains the available options (custom variables).
	Options []VariableOption `json:"options,omitempty"`
}

// VariableOption represents a single option in a variable.
type VariableOption struct {
	Text     string `json:"text"`
	Value    string `json:"value"`
	Selected bool   `json:"selected,omitempty"`
}

// QueryVar creates a query variable that populates from a data source query.
func QueryVar(name, query string) *Variable {
	return &Variable{
		Name:    name,
		Type:    "query",
		Query:   query,
		Refresh: RefreshOnDashboardLoad,
	}
}

// CustomVar creates a custom variable with predefined values.
func CustomVar(name string, values ...string) *Variable {
	query := strings.Join(values, ",")
	options := make([]VariableOption, len(values))
	for i, v := range values {
		options[i] = VariableOption{Text: v, Value: v}
	}
	return &Variable{
		Name:    name,
		Type:    "custom",
		Query:   query,
		Options: options,
	}
}

// IntervalVar creates an interval variable for time range selection.
func IntervalVar(name string, intervals ...string) *Variable {
	query := strings.Join(intervals, ",")
	options := make([]VariableOption, len(intervals))
	for i, v := range intervals {
		options[i] = VariableOption{Text: v, Value: v}
	}
	return &Variable{
		Name:    name,
		Type:    "interval",
		Query:   query,
		Options: options,
	}
}

// DatasourceVar creates a variable that lists data sources by type.
func DatasourceVar(name, dsType string) *Variable {
	return &Variable{
		Name:  name,
		Type:  "datasource",
		Query: dsType,
	}
}

// TextboxVar creates a textbox variable with an optional default value.
func TextboxVar(name, defaultValue string) *Variable {
	return &Variable{
		Name:  name,
		Type:  "textbox",
		Query: defaultValue,
	}
}

// ConstantVar creates a constant (hidden) variable.
func ConstantVar(name, value string) *Variable {
	return &Variable{
		Name:  name,
		Type:  "constant",
		Query: value,
	}
}

// WithDatasource sets the data source for query variables.
func (v *Variable) WithDatasource(ds string) *Variable {
	v.Datasource = ds
	return v
}

// WithRegex sets a regex filter for results.
func (v *Variable) WithRegex(regex string) *Variable {
	v.Regex = regex
	return v
}

// WithSort sets the sort order.
func (v *Variable) WithSort(sort int) *Variable {
	v.Sort = sort
	return v
}

// WithRefresh sets the refresh behavior.
func (v *Variable) WithRefresh(refresh int) *Variable {
	v.Refresh = refresh
	return v
}

// WithLabel sets the display label.
func (v *Variable) WithLabel(label string) *Variable {
	v.Label = label
	return v
}

// WithDescription sets the description.
func (v *Variable) WithDescription(description string) *Variable {
	v.Description = description
	return v
}

// WithDefault sets the default/current value.
func (v *Variable) WithDefault(value string) *Variable {
	v.Current = value
	return v
}

// MultiSelect enables selecting multiple values.
func (v *Variable) MultiSelect() *Variable {
	v.Multi = true
	return v
}

// IncludeAll adds an "All" option.
func (v *Variable) IncludeAll() *Variable {
	v.IncludeAllOption = true
	return v
}

// WithAllValue sets a custom value for the "All" option.
func (v *Variable) WithAllValue(value string) *Variable {
	v.AllValue = value
	return v
}

// Hide hides the variable completely.
func (v *Variable) Hide() *Variable {
	v.HideOption = HideVariable
	return v
}

// HideLabel hides only the label.
func (v *Variable) HideLabel() *Variable {
	v.HideOption = HideLabelOnly
	return v
}

// AutoOption enables auto interval calculation (interval variables only).
func (v *Variable) AutoOption(stepCount int, minInterval string) *Variable {
	v.Auto = true
	v.AutoCount = stepCount
	v.AutoMin = minInterval
	return v
}
