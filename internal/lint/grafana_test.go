package lint

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLintDashboard_MissingTitle(t *testing.T) {
	// Create a temp directory with a dashboard file
	tmpDir := t.TempDir()
	dashboardFile := filepath.Join(tmpDir, "dashboard.go")

	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/grafana"

var MyDashboard = grafana.Dashboard{
	UID: "my-dashboard",
	Rows: []*grafana.Row{},
}
`
	if err := os.WriteFile(dashboardFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	linter := NewLinter()
	result, err := linter.LintAll(tmpDir)
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	// Should have WOB120 error for missing Title
	found := false
	for _, issue := range result.Issues {
		if issue.RuleID == RuleGrafanaDashboardTitle {
			found = true
			if issue.Severity != "error" {
				t.Errorf("expected severity 'error', got %q", issue.Severity)
			}
		}
	}
	if !found {
		t.Error("expected WOB120 issue for missing Dashboard Title")
	}
}

func TestLintDashboard_MissingUID(t *testing.T) {
	tmpDir := t.TempDir()
	dashboardFile := filepath.Join(tmpDir, "dashboard.go")

	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/grafana"

var MyDashboard = grafana.Dashboard{
	Title: "My Dashboard",
	Rows: []*grafana.Row{},
}
`
	if err := os.WriteFile(dashboardFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	linter := NewLinter()
	result, err := linter.LintAll(tmpDir)
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	// Should have WOB122 warning for missing UID
	found := false
	for _, issue := range result.Issues {
		if issue.RuleID == RuleGrafanaDashboardUID {
			found = true
			if issue.Severity != "warning" {
				t.Errorf("expected severity 'warning', got %q", issue.Severity)
			}
		}
	}
	if !found {
		t.Error("expected WOB122 issue for missing Dashboard UID")
	}
}

func TestLintDashboard_MissingRows(t *testing.T) {
	tmpDir := t.TempDir()
	dashboardFile := filepath.Join(tmpDir, "dashboard.go")

	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/grafana"

var MyDashboard = grafana.Dashboard{
	Title: "My Dashboard",
	UID:   "my-dashboard",
}
`
	if err := os.WriteFile(dashboardFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	linter := NewLinter()
	result, err := linter.LintAll(tmpDir)
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	// Should have WOB123 warning for missing Rows
	found := false
	for _, issue := range result.Issues {
		if issue.RuleID == RuleGrafanaDashboardRows {
			found = true
			if issue.Severity != "warning" {
				t.Errorf("expected severity 'warning', got %q", issue.Severity)
			}
		}
	}
	if !found {
		t.Error("expected WOB123 issue for missing Dashboard Rows")
	}
}

func TestLintDashboard_ValidDashboard(t *testing.T) {
	tmpDir := t.TempDir()
	dashboardFile := filepath.Join(tmpDir, "dashboard.go")

	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/grafana"

var Row1 = &grafana.Row{
	Title: "Metrics",
}

var MyDashboard = grafana.Dashboard{
	Title: "My Dashboard",
	UID:   "my-dashboard",
	Rows:  []*grafana.Row{Row1},
}
`
	if err := os.WriteFile(dashboardFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	linter := NewLinter()
	result, err := linter.LintAll(tmpDir)
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	// Should have no Grafana lint issues
	for _, issue := range result.Issues {
		if issue.RuleID == RuleGrafanaDashboardTitle ||
			issue.RuleID == RuleGrafanaDashboardUID ||
			issue.RuleID == RuleGrafanaDashboardRows {
			t.Errorf("unexpected issue %s: %s", issue.RuleID, issue.Message)
		}
	}
}

func TestLintDashboard_DisabledRule(t *testing.T) {
	tmpDir := t.TempDir()
	dashboardFile := filepath.Join(tmpDir, "dashboard.go")

	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/grafana"

var MyDashboard = grafana.Dashboard{
	UID: "my-dashboard",
}
`
	if err := os.WriteFile(dashboardFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	// Disable WOB120
	linter := NewLinterWithOptions(LintOptions{
		DisabledRules: []string{"WOB120"},
	})
	result, err := linter.LintAll(tmpDir)
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	// Should NOT have WOB120 error (disabled)
	for _, issue := range result.Issues {
		if issue.RuleID == RuleGrafanaDashboardTitle {
			t.Error("WOB120 should be disabled but was reported")
		}
	}
}

func TestLintDashboard_PointerDashboard(t *testing.T) {
	tmpDir := t.TempDir()
	dashboardFile := filepath.Join(tmpDir, "dashboard.go")

	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/grafana"

var MyDashboard = &grafana.Dashboard{
	UID: "my-dashboard",
}
`
	if err := os.WriteFile(dashboardFile, []byte(content), 0644); err != nil {
		t.Fatalf("failed to write test file: %v", err)
	}

	linter := NewLinter()
	result, err := linter.LintAll(tmpDir)
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	// Should have WOB120 error for missing Title (even with pointer syntax)
	found := false
	for _, issue := range result.Issues {
		if issue.RuleID == RuleGrafanaDashboardTitle {
			found = true
		}
	}
	if !found {
		t.Error("expected WOB120 issue for missing Dashboard Title with pointer syntax")
	}
}
