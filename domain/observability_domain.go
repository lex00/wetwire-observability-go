// Package domain provides the ObservabilityDomain implementation for wetwire-core-go.
package domain

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-observability-go/internal/discover"
	"github.com/lex00/wetwire-observability-go/internal/lint"
	"github.com/spf13/cobra"
)

// Version is set at build time
var Version = "dev"

// Re-export core types for convenience
type (
	Context      = coredomain.Context
	BuildOpts    = coredomain.BuildOpts
	LintOpts     = coredomain.LintOpts
	InitOpts     = coredomain.InitOpts
	ValidateOpts = coredomain.ValidateOpts
	ListOpts     = coredomain.ListOpts
	GraphOpts    = coredomain.GraphOpts
	Result       = coredomain.Result
	Error        = coredomain.Error
)

var (
	NewResult              = coredomain.NewResult
	NewResultWithData      = coredomain.NewResultWithData
	NewErrorResult         = coredomain.NewErrorResult
	NewErrorResultMultiple = coredomain.NewErrorResultMultiple
)

// ObservabilityDomain implements the Domain interface for observability (Prometheus, Alertmanager, Grafana).
type ObservabilityDomain struct{}

// Compile-time checks
var (
	_ coredomain.Domain       = (*ObservabilityDomain)(nil)
	_ coredomain.ListerDomain = (*ObservabilityDomain)(nil)
)

// Name returns "observability"
func (d *ObservabilityDomain) Name() string {
	return "observability"
}

// Version returns the current version
func (d *ObservabilityDomain) Version() string {
	return Version
}

// Builder returns the observability builder implementation
func (d *ObservabilityDomain) Builder() coredomain.Builder {
	return &observabilityBuilder{}
}

// Linter returns the observability linter implementation
func (d *ObservabilityDomain) Linter() coredomain.Linter {
	return &observabilityLinter{}
}

// Initializer returns the observability initializer implementation
func (d *ObservabilityDomain) Initializer() coredomain.Initializer {
	return &observabilityInitializer{}
}

// Validator returns the observability validator implementation
func (d *ObservabilityDomain) Validator() coredomain.Validator {
	return &observabilityValidator{}
}

// Lister returns the observability lister implementation
func (d *ObservabilityDomain) Lister() coredomain.Lister {
	return &observabilityLister{}
}

// CreateRootCommand creates the root command using the domain interface.
func CreateRootCommand(d coredomain.Domain) *cobra.Command {
	return coredomain.Run(d)
}

// observabilityBuilder implements domain.Builder
type observabilityBuilder struct{}

func (b *observabilityBuilder) Build(ctx *Context, path string, opts BuildOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discover.Discover(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	if resources.TotalCount() == 0 {
		return NewErrorResult("no resources found", Error{
			Path:    absPath,
			Message: "no observability resources found",
		}), nil
	}

	// Build output structure
	outputData := make(map[string]interface{})

	// Group resources by type
	if len(resources.PrometheusConfigs) > 0 {
		outputData["prometheus_configs"] = resourceRefsToMap(resources.PrometheusConfigs)
	}
	if len(resources.ScrapeConfigs) > 0 {
		outputData["scrape_configs"] = resourceRefsToMap(resources.ScrapeConfigs)
	}
	if len(resources.GlobalConfigs) > 0 {
		outputData["global_configs"] = resourceRefsToMap(resources.GlobalConfigs)
	}
	if len(resources.StaticConfigs) > 0 {
		outputData["static_configs"] = resourceRefsToMap(resources.StaticConfigs)
	}
	if len(resources.AlertmanagerConfigs) > 0 {
		outputData["alertmanager_configs"] = resourceRefsToMap(resources.AlertmanagerConfigs)
	}
	if len(resources.RulesFiles) > 0 {
		outputData["rules_files"] = resourceRefsToMap(resources.RulesFiles)
	}
	if len(resources.RuleGroups) > 0 {
		outputData["rule_groups"] = resourceRefsToMap(resources.RuleGroups)
	}
	if len(resources.AlertingRules) > 0 {
		outputData["alerting_rules"] = resourceRefsToMap(resources.AlertingRules)
	}
	if len(resources.RecordingRules) > 0 {
		outputData["recording_rules"] = resourceRefsToMap(resources.RecordingRules)
	}

	// Format output
	var jsonData []byte
	if opts.Format == "pretty" {
		jsonData, err = json.MarshalIndent(outputData, "", "  ")
	} else {
		jsonData, err = json.Marshal(outputData)
	}
	if err != nil {
		return nil, fmt.Errorf("serialization failed: %w", err)
	}

	// Handle output file
	if !opts.DryRun && opts.Output != "" {
		if err := os.WriteFile(opts.Output, jsonData, 0644); err != nil {
			return nil, fmt.Errorf("write output: %w", err)
		}
		return NewResult(fmt.Sprintf("Wrote %s", opts.Output)), nil
	}

	return NewResultWithData("Build completed", string(jsonData)), nil
}

// observabilityLinter implements domain.Linter
type observabilityLinter struct{}

func (l *observabilityLinter) Lint(ctx *Context, path string, opts LintOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Create lint options from LintOpts
	lintOpts := lint.LintOptions{
		DisabledRules: opts.Disable,
		Fix:           opts.Fix,
	}

	// Lint all resources with options
	lintResult, err := lint.LintAllWithOptions(absPath, lintOpts)
	if err != nil {
		return nil, fmt.Errorf("lint failed: %w", err)
	}

	// Handle Fix mode
	if opts.Fix && len(lintResult.Issues) > 0 {
		// Count fixable issues
		fixableCount := 0
		for _, issue := range lintResult.Issues {
			if issue.Fixable {
				fixableCount++
			}
		}
		if fixableCount > 0 {
			return NewResult(fmt.Sprintf("Linted %d resources: %d issues found, auto-fix not yet implemented for %d fixable issues",
				lintResult.ResourceCount, len(lintResult.Issues), fixableCount)), nil
		}
	}

	// Build message based on results
	if len(lintResult.Issues) > 0 {
		return NewResult(fmt.Sprintf("Linted %d resources: %d issues found",
			lintResult.ResourceCount, len(lintResult.Issues))), nil
	}

	// Build success message
	msg := fmt.Sprintf("Linted %d resources: no issues found", lintResult.ResourceCount)
	if len(opts.Disable) > 0 {
		msg += fmt.Sprintf(" (disabled rules: %v)", opts.Disable)
	}
	return NewResult(msg), nil
}

// observabilityInitializer implements domain.Initializer
type observabilityInitializer struct{}

func (i *observabilityInitializer) Init(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Use opts.Path if provided, otherwise fall back to path argument
	targetPath := opts.Path
	if targetPath == "" || targetPath == "." {
		targetPath = path
	}

	// Handle scenario initialization
	if opts.Scenario {
		return i.initScenario(ctx, targetPath, opts)
	}

	// Basic project initialization
	return i.initProject(ctx, targetPath, opts)
}

// initScenario creates a full scenario structure with prompts and expected outputs
func (i *observabilityInitializer) initScenario(ctx *Context, path string, opts InitOpts) (*Result, error) {
	name := opts.Name
	if name == "" {
		name = filepath.Base(path)
	}

	description := opts.Description
	if description == "" {
		description = "Observability scenario"
	}

	// Use core's scenario scaffolding
	scenario := coredomain.ScaffoldScenario(name, description, "observability")
	created, err := coredomain.WriteScenario(path, scenario)
	if err != nil {
		return nil, fmt.Errorf("write scenario: %w", err)
	}

	// Create observability-specific expected directories
	expectedDirs := []string{
		filepath.Join(path, "expected", "prometheus"),
		filepath.Join(path, "expected", "alertmanager"),
		filepath.Join(path, "expected", "rules"),
	}
	for _, dir := range expectedDirs {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("create directory %s: %w", dir, err)
		}
	}

	// Create example prometheus config in expected/prometheus/
	exampleConfig := `package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

// Production is the Prometheus configuration for production
var Production = prometheus.PrometheusConfig{
	Global: &prometheus.GlobalConfig{
		ScrapeInterval:     prometheus.Duration("15s"),
		EvaluationInterval: prometheus.Duration("15s"),
		ExternalLabels: map[string]string{
			"environment": "production",
		},
	},
	ScrapeConfigs: []*prometheus.ScrapeConfig{
		{
			JobName: "prometheus",
			StaticConfigs: []*prometheus.StaticConfig{
				{
					Targets: []string{"localhost:9090"},
				},
			},
		},
	},
}
`
	configPath := filepath.Join(path, "expected", "prometheus", "prometheus.go")
	if err := os.WriteFile(configPath, []byte(exampleConfig), 0644); err != nil {
		return nil, fmt.Errorf("write example config: %w", err)
	}
	created = append(created, "expected/prometheus/prometheus.go")

	return NewResultWithData(
		fmt.Sprintf("Created scenario %s with %d files", name, len(created)),
		created,
	), nil
}

// initProject creates a basic project with example prometheus config
func (i *observabilityInitializer) initProject(ctx *Context, path string, opts InitOpts) (*Result, error) {
	// Create directory
	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("create directory: %w", err)
	}

	// Create example prometheus config file
	exampleContent := `package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

// Production is the Prometheus configuration for production
var Production = prometheus.PrometheusConfig{
	Global: &prometheus.GlobalConfig{
		ScrapeInterval:     prometheus.Duration("15s"),
		EvaluationInterval: prometheus.Duration("15s"),
		ExternalLabels: map[string]string{
			"environment": "production",
		},
	},
	ScrapeConfigs: []*prometheus.ScrapeConfig{
		{
			JobName: "prometheus",
			StaticConfigs: []*prometheus.StaticConfig{
				{
					Targets: []string{"localhost:9090"},
				},
			},
		},
	},
}
`
	examplePath := filepath.Join(path, "prometheus.go")
	if err := os.WriteFile(examplePath, []byte(exampleContent), 0644); err != nil {
		return nil, fmt.Errorf("write example: %w", err)
	}

	return NewResult(fmt.Sprintf("Created %s with example Prometheus config", examplePath)), nil
}

// observabilityValidator implements domain.Validator
type observabilityValidator struct{}

func (v *observabilityValidator) Validate(ctx *Context, path string, opts ValidateOpts) (*Result, error) {
	// For now, validation is the same as lint
	linter := &observabilityLinter{}
	return linter.Lint(ctx, path, LintOpts{})
}

// observabilityLister implements domain.Lister
type observabilityLister struct{}

func (l *observabilityLister) List(ctx *Context, path string, opts ListOpts) (*Result, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("resolve path: %w", err)
	}

	// Discover all resources
	resources, err := discover.Discover(absPath)
	if err != nil {
		return nil, fmt.Errorf("discovery failed: %w", err)
	}

	// Build list
	list := make([]map[string]string, 0)

	for _, ref := range resources.PrometheusConfigs {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "prometheus_config",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.ScrapeConfigs {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "scrape_config",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.GlobalConfigs {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "global_config",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.StaticConfigs {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "static_config",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.AlertmanagerConfigs {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "alertmanager_config",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.RulesFiles {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "rules_file",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.RuleGroups {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "rule_group",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.AlertingRules {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "alerting_rule",
			"file": ref.FilePath,
		})
	}
	for _, ref := range resources.RecordingRules {
		list = append(list, map[string]string{
			"name": ref.Name,
			"type": "recording_rule",
			"file": ref.FilePath,
		})
	}

	return NewResultWithData(fmt.Sprintf("Discovered %d resources", len(list)), list), nil
}

// Helper functions

// resourceRefsToMap converts a slice of ResourceRef to a map keyed by name
func resourceRefsToMap(refs []*discover.ResourceRef) map[string]interface{} {
	result := make(map[string]interface{})
	for _, ref := range refs {
		result[ref.Name] = map[string]interface{}{
			"package":   ref.Package,
			"file_path": ref.FilePath,
			"line":      ref.Line,
		}
	}
	return result
}
