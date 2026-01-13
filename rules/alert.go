package rules

// AlertingRule represents a Prometheus alerting rule.
type AlertingRule struct {
	// Alert is the name of the alert.
	Alert string `yaml:"alert"`

	// Expr is the PromQL expression to evaluate.
	Expr string `yaml:"expr"`

	// For is the duration the condition must be true before firing.
	For Duration `yaml:"for,omitempty"`

	// KeepFiringFor keeps the alert firing for this duration after resolving.
	KeepFiringFor Duration `yaml:"keep_firing_for,omitempty"`

	// Labels are additional labels to attach to the alert.
	Labels map[string]string `yaml:"labels,omitempty"`

	// Annotations provide additional information about the alert.
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// isRule implements the Rule interface.
func (*AlertingRule) isRule() {}

// NewAlertingRule creates a new AlertingRule with the given name.
func NewAlertingRule(name string) *AlertingRule {
	return &AlertingRule{Alert: name}
}

// WithExpr sets the PromQL expression.
func (a *AlertingRule) WithExpr(expr string) *AlertingRule {
	a.Expr = expr
	return a
}

// WithFor sets the duration the condition must be true.
func (a *AlertingRule) WithFor(d Duration) *AlertingRule {
	a.For = d
	return a
}

// WithKeepFiringFor sets the keep firing duration.
func (a *AlertingRule) WithKeepFiringFor(d Duration) *AlertingRule {
	a.KeepFiringFor = d
	return a
}

// WithLabels sets additional labels.
func (a *AlertingRule) WithLabels(labels map[string]string) *AlertingRule {
	a.Labels = labels
	return a
}

// WithAnnotations sets annotations.
func (a *AlertingRule) WithAnnotations(annotations map[string]string) *AlertingRule {
	a.Annotations = annotations
	return a
}

// Critical sets severity to critical.
func (a *AlertingRule) Critical() *AlertingRule {
	return a.withLabel("severity", "critical")
}

// Warning sets severity to warning.
func (a *AlertingRule) Warning() *AlertingRule {
	return a.withLabel("severity", "warning")
}

// Info sets severity to info.
func (a *AlertingRule) Info() *AlertingRule {
	return a.withLabel("severity", "info")
}

// WithSummary sets the summary annotation.
func (a *AlertingRule) WithSummary(summary string) *AlertingRule {
	return a.withAnnotation("summary", summary)
}

// WithDescription sets the description annotation.
func (a *AlertingRule) WithDescription(description string) *AlertingRule {
	return a.withAnnotation("description", description)
}

// WithRunbook sets the runbook_url annotation.
func (a *AlertingRule) WithRunbook(url string) *AlertingRule {
	return a.withAnnotation("runbook_url", url)
}

// withLabel is a helper to add a single label.
func (a *AlertingRule) withLabel(key, value string) *AlertingRule {
	if a.Labels == nil {
		a.Labels = make(map[string]string)
	}
	a.Labels[key] = value
	return a
}

// withAnnotation is a helper to add a single annotation.
func (a *AlertingRule) withAnnotation(key, value string) *AlertingRule {
	if a.Annotations == nil {
		a.Annotations = make(map[string]string)
	}
	a.Annotations[key] = value
	return a
}
