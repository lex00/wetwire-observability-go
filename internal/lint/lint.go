// Package lint provides linting functionality for wetwire-observability-go.
package lint

import (
	"fmt"
	"strings"

	"github.com/lex00/wetwire-observability-go/internal/discover"
)

// LintOptions contains options for the linter.
type LintOptions struct {
	// DisabledRules is a list of rule IDs to skip (e.g., ["WOB001", "WOB002"]).
	DisabledRules []string

	// Fix indicates whether to automatically fix fixable issues.
	Fix bool
}

// LintIssue represents a single lint issue found during linting.
type LintIssue struct {
	// RuleID is the lint rule identifier (e.g., "WOB001").
	RuleID string

	// Severity is the severity level ("error" or "warning").
	Severity string

	// Message describes the issue.
	Message string

	// File is the path to the file where the issue was found.
	File string

	// Line is the line number where the issue was found.
	Line int

	// Fixable indicates whether this issue can be auto-fixed.
	Fixable bool
}

// LintResult contains the results of linting.
type LintResult struct {
	// Issues is the list of lint issues found.
	Issues []LintIssue

	// ResourceCount is the number of resources that were linted.
	ResourceCount int

	// FixRequested indicates whether fix mode was requested.
	FixRequested bool

	// DisabledRules is the list of rules that were disabled.
	DisabledRules []string
}

// Linter performs lint checks on discovered resources.
type Linter struct {
	options LintOptions
}

// NewLinter creates a new Linter with default options.
func NewLinter() *Linter {
	return &Linter{}
}

// NewLinterWithOptions creates a new Linter with the specified options.
func NewLinterWithOptions(opts LintOptions) *Linter {
	return &Linter{options: opts}
}

// LintAll lints all resources discovered at the given path.
func (l *Linter) LintAll(path string) (*LintResult, error) {
	resources, err := discover.Discover(path)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	return l.lintResources(resources)
}

// LintAllWithOptions lints all resources with the specified options.
func LintAllWithOptions(path string, opts LintOptions) (*LintResult, error) {
	linter := NewLinterWithOptions(opts)
	return linter.LintAll(path)
}

// lintResources performs the actual linting on discovered resources.
func (l *Linter) lintResources(resources *discover.DiscoveryResult) (*LintResult, error) {
	result := &LintResult{
		ResourceCount: resources.TotalCount(),
		FixRequested:  l.options.Fix,
		DisabledRules: l.options.DisabledRules,
		Issues:        []LintIssue{},
	}

	// Build a set of disabled rules for efficient lookup
	disabledSet := make(map[string]bool)
	for _, rule := range l.options.DisabledRules {
		disabledSet[strings.ToUpper(rule)] = true
	}

	// Run lint rules on each resource type

	// Lint Grafana dashboards (WOB120-149)
	if len(resources.Dashboards) > 0 {
		dashboardIssues := l.lintDashboards(resources.Dashboards)
		result.Issues = append(result.Issues, dashboardIssues...)
	}

	// Filter out issues from disabled rules
	filteredIssues := []LintIssue{}
	for _, issue := range result.Issues {
		if !disabledSet[strings.ToUpper(issue.RuleID)] {
			filteredIssues = append(filteredIssues, issue)
		}
	}
	result.Issues = filteredIssues

	return result, nil
}

// IsRuleDisabled checks if a rule is disabled.
func (l *Linter) IsRuleDisabled(ruleID string) bool {
	for _, disabled := range l.options.DisabledRules {
		if strings.EqualFold(disabled, ruleID) {
			return true
		}
	}
	return false
}
