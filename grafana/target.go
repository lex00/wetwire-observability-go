package grafana

import "github.com/lex00/wetwire-observability-go/promql"

// PrometheusTarget represents a Prometheus query target.
type PrometheusTarget struct {
	// RefID is the unique reference ID for this target.
	RefID string `json:"refId"`

	// Expr is the PromQL expression.
	Expr string `json:"expr"`

	// LegendFormat is the legend template (e.g., "{{instance}}").
	LegendFormat string `json:"legendFormat,omitempty"`

	// Interval is the query interval (e.g., "$__rate_interval").
	Interval string `json:"interval,omitempty"`

	// IsInstant indicates whether this is an instant query.
	IsInstant bool `json:"instant,omitempty"`

	// Hidden hides this target from the legend.
	Hidden bool `json:"hide,omitempty"`

	// Datasource is the data source reference.
	Datasource string `json:"datasource,omitempty"`

	// EditorMode is the query editor mode.
	EditorMode string `json:"editorMode,omitempty"`

	// Range indicates a range query (for instant toggle).
	IsRange bool `json:"range,omitempty"`
}

// PromTarget creates a PrometheusTarget from a string expression.
func PromTarget(expr string) *PrometheusTarget {
	return &PrometheusTarget{
		RefID: "A",
		Expr:  expr,
	}
}

// PromTargetExpr creates a PrometheusTarget from a typed PromQL expression.
func PromTargetExpr(expr promql.Expr) *PrometheusTarget {
	return &PrometheusTarget{
		RefID: "A",
		Expr:  expr.String(),
	}
}

// WithRefID sets the reference ID.
func (t *PrometheusTarget) WithRefID(refID string) *PrometheusTarget {
	t.RefID = refID
	return t
}

// WithLegendFormat sets the legend format template.
func (t *PrometheusTarget) WithLegendFormat(format string) *PrometheusTarget {
	t.LegendFormat = format
	return t
}

// WithInterval sets the query interval.
func (t *PrometheusTarget) WithInterval(interval string) *PrometheusTarget {
	t.Interval = interval
	return t
}

// Instant sets the query to instant mode.
func (t *PrometheusTarget) Instant() *PrometheusTarget {
	t.IsInstant = true
	t.IsRange = false
	return t
}

// Range sets the query to range mode.
func (t *PrometheusTarget) Range() *PrometheusTarget {
	t.IsInstant = false
	t.IsRange = true
	return t
}

// Hide hides this target from the legend.
func (t *PrometheusTarget) Hide() *PrometheusTarget {
	t.Hidden = true
	return t
}

// WithDatasource sets the data source reference.
func (t *PrometheusTarget) WithDatasource(ds string) *PrometheusTarget {
	t.Datasource = ds
	return t
}

// LokiQueryTarget represents a Loki query target.
type LokiQueryTarget struct {
	// RefID is the unique reference ID for this target.
	RefID string `json:"refId"`

	// Expr is the LogQL expression.
	Expr string `json:"expr"`

	// LegendFormat is the legend template.
	LegendFormat string `json:"legendFormat,omitempty"`

	// MaxLines limits the number of log lines returned.
	MaxLines int `json:"maxLines,omitempty"`

	// Hidden hides this target from the legend.
	Hidden bool `json:"hide,omitempty"`

	// Datasource is the data source reference.
	Datasource string `json:"datasource,omitempty"`

	// QueryType is the query type (range, instant).
	QueryType string `json:"queryType,omitempty"`
}

// LokiTarget creates a new Loki query target.
func LokiTarget(expr string) *LokiQueryTarget {
	return &LokiQueryTarget{
		RefID: "A",
		Expr:  expr,
	}
}

// WithRefID sets the reference ID.
func (t *LokiQueryTarget) WithRefID(refID string) *LokiQueryTarget {
	t.RefID = refID
	return t
}

// WithLegendFormat sets the legend format template.
func (t *LokiQueryTarget) WithLegendFormat(format string) *LokiQueryTarget {
	t.LegendFormat = format
	return t
}

// WithMaxLines sets the maximum number of log lines.
func (t *LokiQueryTarget) WithMaxLines(max int) *LokiQueryTarget {
	t.MaxLines = max
	return t
}

// Hide hides this target from the legend.
func (t *LokiQueryTarget) Hide() *LokiQueryTarget {
	t.Hidden = true
	return t
}

// WithDatasource sets the data source reference.
func (t *LokiQueryTarget) WithDatasource(ds string) *LokiQueryTarget {
	t.Datasource = ds
	return t
}
