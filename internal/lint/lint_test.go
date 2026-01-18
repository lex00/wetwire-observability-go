package lint

import (
	"testing"
)

func TestNewLinter(t *testing.T) {
	linter := NewLinter()
	if linter == nil {
		t.Fatal("expected non-nil linter")
	}
}

func TestNewLinterWithOptions(t *testing.T) {
	opts := LintOptions{
		DisabledRules: []string{"WOB001", "WOB002"},
		Fix:           true,
	}
	linter := NewLinterWithOptions(opts)
	if linter == nil {
		t.Fatal("expected non-nil linter")
	}
	if len(linter.options.DisabledRules) != 2 {
		t.Errorf("expected 2 disabled rules, got %d", len(linter.options.DisabledRules))
	}
	if !linter.options.Fix {
		t.Error("expected Fix to be true")
	}
}

func TestIsRuleDisabled(t *testing.T) {
	tests := []struct {
		name          string
		disabledRules []string
		ruleID        string
		want          bool
	}{
		{
			name:          "rule is disabled",
			disabledRules: []string{"WOB001", "WOB002"},
			ruleID:        "WOB001",
			want:          true,
		},
		{
			name:          "rule is not disabled",
			disabledRules: []string{"WOB001", "WOB002"},
			ruleID:        "WOB003",
			want:          false,
		},
		{
			name:          "case insensitive match",
			disabledRules: []string{"wob001"},
			ruleID:        "WOB001",
			want:          true,
		},
		{
			name:          "no disabled rules",
			disabledRules: []string{},
			ruleID:        "WOB001",
			want:          false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			linter := NewLinterWithOptions(LintOptions{DisabledRules: tt.disabledRules})
			got := linter.IsRuleDisabled(tt.ruleID)
			if got != tt.want {
				t.Errorf("IsRuleDisabled(%q) = %v, want %v", tt.ruleID, got, tt.want)
			}
		})
	}
}

func TestLintResultFieldsWithOptions(t *testing.T) {
	opts := LintOptions{
		DisabledRules: []string{"WOB001"},
		Fix:           true,
	}
	linter := NewLinterWithOptions(opts)

	// Lint the current directory (may not have resources, but should not error)
	result, err := linter.LintAll(".")
	if err != nil {
		t.Fatalf("LintAll failed: %v", err)
	}

	if !result.FixRequested {
		t.Error("expected FixRequested to be true")
	}
	if len(result.DisabledRules) != 1 {
		t.Errorf("expected 1 disabled rule, got %d", len(result.DisabledRules))
	}
	if result.DisabledRules[0] != "WOB001" {
		t.Errorf("expected disabled rule WOB001, got %s", result.DisabledRules[0])
	}
}

func TestLintAllWithOptions(t *testing.T) {
	opts := LintOptions{
		DisabledRules: []string{"WOB001", "WOB002"},
		Fix:           false,
	}

	result, err := LintAllWithOptions(".", opts)
	if err != nil {
		t.Fatalf("LintAllWithOptions failed: %v", err)
	}

	if result.FixRequested {
		t.Error("expected FixRequested to be false")
	}
	if len(result.DisabledRules) != 2 {
		t.Errorf("expected 2 disabled rules, got %d", len(result.DisabledRules))
	}
}
