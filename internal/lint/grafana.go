package lint

import (
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/lex00/wetwire-observability-go/internal/discover"
)

// Grafana lint rules (WOB120-149)
const (
	RuleGrafanaDashboardTitle = "WOB120"
	RuleGrafanaDashboardUID   = "WOB122"
	RuleGrafanaDashboardRows  = "WOB123"
	RuleGrafanaPanelTitle     = "WOB124"
	RuleGrafanaPanelTarget    = "WOB125"
)

// lintDashboards runs Grafana-specific lint rules on discovered dashboards.
func (l *Linter) lintDashboards(dashboards []*discover.ResourceRef) []LintIssue {
	var issues []LintIssue

	for _, ref := range dashboards {
		dashboardIssues := l.lintDashboardFile(ref)
		issues = append(issues, dashboardIssues...)
	}

	return issues
}

// lintDashboardFile parses a file and lints Dashboard declarations.
func (l *Linter) lintDashboardFile(ref *discover.ResourceRef) []LintIssue {
	var issues []LintIssue

	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, ref.FilePath, nil, parser.ParseComments)
	if err != nil {
		return issues
	}

	// Find the specific declaration
	ast.Inspect(file, func(n ast.Node) bool {
		genDecl, ok := n.(*ast.GenDecl)
		if !ok || genDecl.Tok != token.VAR {
			return true
		}

		for _, spec := range genDecl.Specs {
			valueSpec, ok := spec.(*ast.ValueSpec)
			if !ok {
				continue
			}

			for i, name := range valueSpec.Names {
				if name.Name != ref.Name {
					continue
				}

				if i >= len(valueSpec.Values) {
					continue
				}

				// Check if this is a Dashboard composite literal
				compLit, ok := valueSpec.Values[i].(*ast.CompositeLit)
				if !ok {
					// Could be a pointer: &grafana.Dashboard{}
					unaryExpr, ok := valueSpec.Values[i].(*ast.UnaryExpr)
					if ok && unaryExpr.Op == token.AND {
						compLit, ok = unaryExpr.X.(*ast.CompositeLit)
						if !ok {
							continue
						}
					} else {
						continue
					}
				}

				// Verify it's a Dashboard type
				if !isDashboardType(compLit.Type) {
					continue
				}

				// Run lint checks on this dashboard
				issues = append(issues, l.checkDashboardLiteral(compLit, ref, fset)...)
			}
		}
		return true
	})

	return issues
}

// isDashboardType checks if the type expression is a Dashboard.
func isDashboardType(expr ast.Expr) bool {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name == "Dashboard"
	case *ast.SelectorExpr:
		return t.Sel.Name == "Dashboard"
	case *ast.StarExpr:
		return isDashboardType(t.X)
	}
	return false
}

// checkDashboardLiteral runs lint rules on a Dashboard composite literal.
func (l *Linter) checkDashboardLiteral(lit *ast.CompositeLit, ref *discover.ResourceRef, fset *token.FileSet) []LintIssue {
	var issues []LintIssue

	// Extract field values
	fields := extractCompositeLitFields(lit)

	// WOB120: Dashboard must have Title
	if !l.IsRuleDisabled(RuleGrafanaDashboardTitle) {
		if _, hasTitle := fields["Title"]; !hasTitle {
			issues = append(issues, LintIssue{
				RuleID:   RuleGrafanaDashboardTitle,
				Severity: "error",
				Message:  "Dashboard must have a Title",
				File:     ref.FilePath,
				Line:     ref.Line,
				Fixable:  false,
			})
		} else if isEmptyStringLiteral(fields["Title"]) {
			issues = append(issues, LintIssue{
				RuleID:   RuleGrafanaDashboardTitle,
				Severity: "error",
				Message:  "Dashboard Title must not be empty",
				File:     ref.FilePath,
				Line:     ref.Line,
				Fixable:  false,
			})
		}
	}

	// WOB122: Dashboard should have UID
	if !l.IsRuleDisabled(RuleGrafanaDashboardUID) {
		if _, hasUID := fields["UID"]; !hasUID {
			issues = append(issues, LintIssue{
				RuleID:   RuleGrafanaDashboardUID,
				Severity: "warning",
				Message:  "Dashboard should have a UID for stable identification",
				File:     ref.FilePath,
				Line:     ref.Line,
				Fixable:  false,
			})
		} else if isEmptyStringLiteral(fields["UID"]) {
			issues = append(issues, LintIssue{
				RuleID:   RuleGrafanaDashboardUID,
				Severity: "warning",
				Message:  "Dashboard UID should not be empty",
				File:     ref.FilePath,
				Line:     ref.Line,
				Fixable:  false,
			})
		}
	}

	// WOB123: Dashboard should have Rows
	if !l.IsRuleDisabled(RuleGrafanaDashboardRows) {
		rowsExpr, hasRows := fields["Rows"]
		if !hasRows {
			issues = append(issues, LintIssue{
				RuleID:   RuleGrafanaDashboardRows,
				Severity: "warning",
				Message:  "Dashboard should have at least one Row",
				File:     ref.FilePath,
				Line:     ref.Line,
				Fixable:  false,
			})
		} else if isEmptySliceLiteral(rowsExpr) {
			issues = append(issues, LintIssue{
				RuleID:   RuleGrafanaDashboardRows,
				Severity: "warning",
				Message:  "Dashboard Rows should not be empty",
				File:     ref.FilePath,
				Line:     ref.Line,
				Fixable:  false,
			})
		}
	}

	return issues
}

// extractCompositeLitFields extracts field names and values from a composite literal.
func extractCompositeLitFields(lit *ast.CompositeLit) map[string]ast.Expr {
	fields := make(map[string]ast.Expr)
	for _, elt := range lit.Elts {
		kv, ok := elt.(*ast.KeyValueExpr)
		if !ok {
			continue
		}
		key, ok := kv.Key.(*ast.Ident)
		if !ok {
			continue
		}
		fields[key.Name] = kv.Value
	}
	return fields
}

// isEmptyStringLiteral checks if an expression is an empty string literal.
func isEmptyStringLiteral(expr ast.Expr) bool {
	lit, ok := expr.(*ast.BasicLit)
	if !ok || lit.Kind != token.STRING {
		return false
	}
	// String literals are quoted, so empty is `""` or ``
	return lit.Value == `""` || lit.Value == "``"
}

// isEmptySliceLiteral checks if an expression is an empty slice literal.
func isEmptySliceLiteral(expr ast.Expr) bool {
	compLit, ok := expr.(*ast.CompositeLit)
	if !ok {
		return false
	}
	return len(compLit.Elts) == 0
}

