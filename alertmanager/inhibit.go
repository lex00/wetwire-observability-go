package alertmanager

// NewInhibitRule creates a new InhibitRule.
func NewInhibitRule() *InhibitRule {
	return &InhibitRule{}
}

// WithSourceMatch sets the source match labels (deprecated, use WithSourceMatchers).
func (i *InhibitRule) WithSourceMatch(match map[string]string) *InhibitRule {
	i.SourceMatch = match
	return i
}

// WithSourceMatchers sets the source matchers.
func (i *InhibitRule) WithSourceMatchers(matchers ...*Matcher) *InhibitRule {
	i.SourceMatchers = matchers
	return i
}

// WithTargetMatch sets the target match labels (deprecated, use WithTargetMatchers).
func (i *InhibitRule) WithTargetMatch(match map[string]string) *InhibitRule {
	i.TargetMatch = match
	return i
}

// WithTargetMatchers sets the target matchers.
func (i *InhibitRule) WithTargetMatchers(matchers ...*Matcher) *InhibitRule {
	i.TargetMatchers = matchers
	return i
}

// WithEqual sets the labels that must be equal between source and target.
func (i *InhibitRule) WithEqual(labels ...string) *InhibitRule {
	i.Equal = labels
	return i
}

// CriticalInhibitsWarning creates a rule where critical alerts inhibit warning alerts
// for the same alertname.
func CriticalInhibitsWarning() *InhibitRule {
	return NewInhibitRule().
		WithSourceMatchers(Eq("severity", "critical")).
		WithTargetMatchers(Eq("severity", "warning")).
		WithEqual("alertname")
}

// ErrorInhibitsWarning creates a rule where error alerts inhibit warning alerts
// for the same alertname.
func ErrorInhibitsWarning() *InhibitRule {
	return NewInhibitRule().
		WithSourceMatchers(Eq("severity", "error")).
		WithTargetMatchers(Eq("severity", "warning")).
		WithEqual("alertname")
}

// CriticalInhibitsInfo creates a rule where critical alerts inhibit info alerts
// for the same alertname.
func CriticalInhibitsInfo() *InhibitRule {
	return NewInhibitRule().
		WithSourceMatchers(Eq("severity", "critical")).
		WithTargetMatchers(Eq("severity", "info")).
		WithEqual("alertname")
}

// ErrorInhibitsInfo creates a rule where error alerts inhibit info alerts
// for the same alertname.
func ErrorInhibitsInfo() *InhibitRule {
	return NewInhibitRule().
		WithSourceMatchers(Eq("severity", "error")).
		WithTargetMatchers(Eq("severity", "info")).
		WithEqual("alertname")
}

// WarningInhibitsInfo creates a rule where warning alerts inhibit info alerts
// for the same alertname.
func WarningInhibitsInfo() *InhibitRule {
	return NewInhibitRule().
		WithSourceMatchers(Eq("severity", "warning")).
		WithTargetMatchers(Eq("severity", "info")).
		WithEqual("alertname")
}
