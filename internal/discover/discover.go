// Package discover provides AST-based discovery of wetwire resources.
package discover

import (
	"go/parser"
	"go/token"

	coreast "github.com/lex00/wetwire-core-go/ast"
	corediscover "github.com/lex00/wetwire-core-go/discover"
)

// ResourceRef represents a discovered wetwire resource.
type ResourceRef struct {
	// Package is the full package path.
	Package string `json:"package"`
	// Name is the variable/constant name.
	Name string `json:"name"`
	// Type is the resource type (e.g., "PrometheusConfig", "ScrapeConfig").
	Type string `json:"type"`
	// FilePath is the absolute path to the source file.
	FilePath string `json:"file_path"`
	// Line is the line number where the resource is declared.
	Line int `json:"line"`
}

// DiscoveryResult contains all discovered wetwire resources.
type DiscoveryResult struct {
	// PrometheusConfigs are discovered Prometheus configuration resources.
	PrometheusConfigs []*ResourceRef `json:"prometheus_configs,omitempty"`
	// ScrapeConfigs are discovered scrape configuration resources.
	ScrapeConfigs []*ResourceRef `json:"scrape_configs,omitempty"`
	// GlobalConfigs are discovered global configuration resources.
	GlobalConfigs []*ResourceRef `json:"global_configs,omitempty"`
	// StaticConfigs are discovered static configuration resources.
	StaticConfigs []*ResourceRef `json:"static_configs,omitempty"`
	// AlertmanagerConfigs are discovered Alertmanager configuration resources.
	AlertmanagerConfigs []*ResourceRef `json:"alertmanager_configs,omitempty"`
	// RulesFiles are discovered rule file resources.
	RulesFiles []*ResourceRef `json:"rules_files,omitempty"`
	// RuleGroups are discovered rule group resources.
	RuleGroups []*ResourceRef `json:"rule_groups,omitempty"`
	// AlertingRules are discovered alerting rule resources.
	AlertingRules []*ResourceRef `json:"alerting_rules,omitempty"`
	// RecordingRules are discovered recording rule resources.
	RecordingRules []*ResourceRef `json:"recording_rules,omitempty"`
	// Dashboards are discovered Grafana dashboard resources.
	Dashboards []*ResourceRef `json:"dashboards,omitempty"`
	// Errors encountered during discovery (non-fatal).
	Errors []string `json:"errors,omitempty"`
}

// observabilityTypes are the types we look for in wetwire packages.
var observabilityTypes = map[string]bool{
	// Prometheus types
	"PrometheusConfig": true,
	"GlobalConfig":     true,
	"ScrapeConfig":     true,
	"StaticConfig":     true,
	// Alertmanager types
	"AlertmanagerConfig": true,
	// Rules types
	"RulesFile":     true,
	"RuleGroup":     true,
	"AlertingRule":  true,
	"RecordingRule": true,
	// Grafana types
	"Dashboard": true,
}

// observabilityTypeMatcher creates a TypeMatcher for observability types.
func observabilityTypeMatcher(pkgName, typeName string, imports map[string]string) (string, bool) {
	// Check if this is an observability type
	if observabilityTypes[typeName] {
		return typeName, true
	}
	return "", false
}

// Discover finds all wetwire resources in the given directory.
// It recursively searches all Go files and parses them for resource declarations.
func Discover(dir string) (*DiscoveryResult, error) {
	// Use core discover package
	coreResult, err := corediscover.Discover(corediscover.DiscoverOptions{
		Packages:    []string{dir},
		TypeMatcher: observabilityTypeMatcher,
	})
	if err != nil {
		return nil, err
	}

	// Convert core result to observability-specific result
	result := &DiscoveryResult{}

	// Extract package name from files
	packageCache := make(map[string]string)

	for _, resource := range coreResult.Resources {
		// Skip unexported variables (wetwire only considers exported resources)
		if !isExported(resource.Name) {
			continue
		}

		// Get the package name from file if not cached
		pkg, ok := packageCache[resource.File]
		if !ok {
			// Extract package name by parsing the file using core AST utilities
			pkg = extractPackageFromFile(resource.File)
			packageCache[resource.File] = pkg
		}

		ref := &ResourceRef{
			Package:  pkg,
			Name:     resource.Name,
			Type:     resource.Type,
			FilePath: resource.File,
			Line:     resource.Line,
		}

		// Categorize by type
		switch resource.Type {
		case "PrometheusConfig":
			result.PrometheusConfigs = append(result.PrometheusConfigs, ref)
		case "ScrapeConfig":
			result.ScrapeConfigs = append(result.ScrapeConfigs, ref)
		case "GlobalConfig":
			result.GlobalConfigs = append(result.GlobalConfigs, ref)
		case "StaticConfig":
			result.StaticConfigs = append(result.StaticConfigs, ref)
		case "AlertmanagerConfig":
			result.AlertmanagerConfigs = append(result.AlertmanagerConfigs, ref)
		case "RulesFile":
			result.RulesFiles = append(result.RulesFiles, ref)
		case "RuleGroup":
			result.RuleGroups = append(result.RuleGroups, ref)
		case "AlertingRule":
			result.AlertingRules = append(result.AlertingRules, ref)
		case "RecordingRule":
			result.RecordingRules = append(result.RecordingRules, ref)
		case "Dashboard":
			result.Dashboards = append(result.Dashboards, ref)
		}
	}

	// Convert errors
	for _, err := range coreResult.Errors {
		result.Errors = append(result.Errors, err.Error())
	}

	return result, nil
}

// extractPackageFromFile extracts the package name from a Go file.
// This uses the core AST utilities for parsing.
func extractPackageFromFile(path string) string {
	// Use core ast.ParseFile to get the package name
	astFile, _, err := coreast.ParseFile(path)
	if err != nil {
		// Fallback: try direct parsing with minimal mode
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, path, nil, parser.PackageClauseOnly)
		if err != nil {
			return ""
		}
		return file.Name.Name
	}

	if astFile != nil && astFile.Name != nil {
		return astFile.Name.Name
	}

	return ""
}

// TotalCount returns the total number of discovered resources.
func (r *DiscoveryResult) TotalCount() int {
	return len(r.PrometheusConfigs) + len(r.ScrapeConfigs) +
		len(r.GlobalConfigs) + len(r.StaticConfigs) +
		len(r.AlertmanagerConfigs) +
		len(r.RulesFiles) + len(r.RuleGroups) +
		len(r.AlertingRules) + len(r.RecordingRules) +
		len(r.Dashboards)
}

// All returns all discovered resources as a flat slice.
func (r *DiscoveryResult) All() []*ResourceRef {
	var all []*ResourceRef
	all = append(all, r.PrometheusConfigs...)
	all = append(all, r.ScrapeConfigs...)
	all = append(all, r.GlobalConfigs...)
	all = append(all, r.StaticConfigs...)
	all = append(all, r.AlertmanagerConfigs...)
	all = append(all, r.RulesFiles...)
	all = append(all, r.RuleGroups...)
	all = append(all, r.AlertingRules...)
	all = append(all, r.RecordingRules...)
	all = append(all, r.Dashboards...)
	return all
}

// isExported reports whether name starts with an upper-case letter.
func isExported(name string) bool {
	if len(name) == 0 {
		return false
	}
	r := rune(name[0])
	return r >= 'A' && r <= 'Z'
}

