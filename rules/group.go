// Package rules provides types for Prometheus rule configuration synthesis.
package rules

import "github.com/lex00/wetwire-observability-go/prometheus"

// Duration is an alias for prometheus.Duration for consistency.
type Duration = prometheus.Duration

// Convenience duration constants.
const (
	Second = prometheus.Second
	Minute = prometheus.Minute
	Hour   = prometheus.Hour
)

// RulesFile represents a complete rules file containing multiple groups.
type RulesFile struct {
	// Groups contains the rule groups.
	Groups []*RuleGroup `yaml:"groups"`
}

// RuleGroup represents a group of rules evaluated together.
type RuleGroup struct {
	// Name is the unique name of the group.
	Name string `yaml:"name"`

	// Interval is the evaluation interval for rules in this group.
	Interval Duration `yaml:"interval,omitempty"`

	// Limit is the maximum number of alerts to produce.
	Limit int `yaml:"limit,omitempty"`

	// Rules contains the alerting and recording rules.
	// Use any to allow both AlertingRule and RecordingRule.
	Rules []any `yaml:"rules,omitempty"`
}

// Rule is an interface implemented by AlertingRule and RecordingRule.
type Rule interface {
	isRule()
}

// NewRulesFile creates a new RulesFile.
func NewRulesFile() *RulesFile {
	return &RulesFile{}
}

// WithGroups sets the rule groups.
func (f *RulesFile) WithGroups(groups ...*RuleGroup) *RulesFile {
	f.Groups = groups
	return f
}

// AddGroup adds a rule group.
func (f *RulesFile) AddGroup(group *RuleGroup) *RulesFile {
	f.Groups = append(f.Groups, group)
	return f
}

// NewRuleGroup creates a new RuleGroup with the given name.
func NewRuleGroup(name string) *RuleGroup {
	return &RuleGroup{Name: name}
}

// WithInterval sets the evaluation interval.
func (g *RuleGroup) WithInterval(d Duration) *RuleGroup {
	g.Interval = d
	return g
}

// WithLimit sets the alert limit.
func (g *RuleGroup) WithLimit(limit int) *RuleGroup {
	g.Limit = limit
	return g
}

// WithRules sets the rules in this group.
func (g *RuleGroup) WithRules(rules ...any) *RuleGroup {
	g.Rules = rules
	return g
}

// AddRule adds a rule to this group.
func (g *RuleGroup) AddRule(rule any) *RuleGroup {
	g.Rules = append(g.Rules, rule)
	return g
}
