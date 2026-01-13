package rules

// RecordingRule represents a Prometheus recording rule.
type RecordingRule struct {
	// Record is the name of the metric to record.
	Record string `yaml:"record"`

	// Expr is the PromQL expression to evaluate.
	Expr string `yaml:"expr"`

	// Labels are additional labels to attach to the recorded metric.
	Labels map[string]string `yaml:"labels,omitempty"`
}

// isRule implements the Rule interface.
func (*RecordingRule) isRule() {}

// NewRecordingRule creates a new RecordingRule with the given name.
func NewRecordingRule(name string) *RecordingRule {
	return &RecordingRule{Record: name}
}

// WithExpr sets the PromQL expression.
func (r *RecordingRule) WithExpr(expr string) *RecordingRule {
	r.Expr = expr
	return r
}

// WithLabels sets additional labels.
func (r *RecordingRule) WithLabels(labels map[string]string) *RecordingRule {
	r.Labels = labels
	return r
}
