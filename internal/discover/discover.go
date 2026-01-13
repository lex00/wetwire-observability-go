// Package discover provides AST-based discovery of wetwire resources.
package discover

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
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
}

// normalizeTypeName extracts the simple type name from a potentially qualified type.
// e.g., "prometheus.PrometheusConfig" -> "PrometheusConfig"
func normalizeTypeName(typeName string) string {
	if idx := strings.LastIndex(typeName, "."); idx != -1 {
		return typeName[idx+1:]
	}
	return typeName
}

// Discover finds all wetwire resources in the given directory.
// It recursively searches all Go files and parses them for resource declarations.
func Discover(dir string) (*DiscoveryResult, error) {
	result := &DiscoveryResult{}

	// Walk the directory tree
	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
			return nil // Continue walking
		}

		// Skip hidden directories and vendor
		if info.IsDir() {
			name := info.Name()
			if name == "vendor" || strings.HasPrefix(name, ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// Only process .go files
		if !strings.HasSuffix(path, ".go") {
			return nil
		}

		// Skip test files
		if strings.HasSuffix(path, "_test.go") {
			return nil
		}

		// Parse the file
		refs, errs := discoverFile(path)
		for _, ref := range refs {
			switch ref.Type {
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
			}
		}
		result.Errors = append(result.Errors, errs...)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return result, nil
}

// discoverFile parses a single Go file and returns discovered resources.
func discoverFile(path string) ([]*ResourceRef, []string) {
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, path, nil, parser.ParseComments)
	if err != nil {
		return nil, []string{err.Error()}
	}

	var refs []*ResourceRef
	packageName := file.Name.Name

	// Walk the AST looking for top-level variable declarations
	for _, decl := range file.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if genDecl.Tok != token.VAR {
			continue
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				// Skip unexported variables
				if !ast.IsExported(name.Name) {
					continue
				}

				// Try to determine the type
				typeName := extractTypeName(valueSpec, i)
				if typeName == "" {
					continue
				}

				// Normalize the type name (strip package prefix)
				simpleType := normalizeTypeName(typeName)

				// Check if it's an observability type we care about
				if !observabilityTypes[simpleType] {
					continue
				}

				pos := fset.Position(name.Pos())
				refs = append(refs, &ResourceRef{
					Package:  packageName,
					Name:     name.Name,
					Type:     simpleType, // Use the normalized type name
					FilePath: pos.Filename,
					Line:     pos.Line,
				})
			}
		}
	}

	return refs, nil
}

// extractTypeName extracts the type name from a ValueSpec.
func extractTypeName(spec *ast.ValueSpec, index int) string {
	// First try the explicit type if present
	if spec.Type != nil {
		return typeToString(spec.Type)
	}

	// Try to infer from the value
	if index < len(spec.Values) {
		return inferTypeFromValue(spec.Values[index])
	}

	return ""
}

// typeToString converts an ast.Expr type to a string.
func typeToString(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		// Pointer type, get the underlying type
		return typeToString(t.X)
	case *ast.SelectorExpr:
		// Package-qualified type (e.g., prometheus.PrometheusConfig)
		if ident, ok := t.X.(*ast.Ident); ok {
			return ident.Name + "." + t.Sel.Name
		}
		return t.Sel.Name
	case *ast.ArrayType:
		// Array/slice type
		elemType := typeToString(t.Elt)
		return "[]" + elemType
	}
	return ""
}

// inferTypeFromValue tries to infer the type from a value expression.
func inferTypeFromValue(expr ast.Expr) string {
	switch v := expr.(type) {
	case *ast.UnaryExpr:
		// Handle &Type{} expressions
		if v.Op == token.AND {
			return inferTypeFromValue(v.X)
		}
	case *ast.CompositeLit:
		// Handle Type{} composite literals
		return typeToString(v.Type)
	}
	return ""
}

// TotalCount returns the total number of discovered resources.
func (r *DiscoveryResult) TotalCount() int {
	return len(r.PrometheusConfigs) + len(r.ScrapeConfigs) +
		len(r.GlobalConfigs) + len(r.StaticConfigs) +
		len(r.AlertmanagerConfigs) +
		len(r.RulesFiles) + len(r.RuleGroups) +
		len(r.AlertingRules) + len(r.RecordingRules)
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
	return all
}
