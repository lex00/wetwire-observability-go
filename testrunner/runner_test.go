package testrunner

import (
	"testing"
)

func TestNewRunner(t *testing.T) {
	r := NewRunner()
	if r == nil {
		t.Fatal("NewRunner returned nil")
	}
}

func TestRunner_WithPersona(t *testing.T) {
	r := NewRunner().WithPersona("sre")
	if len(r.personas) != 1 {
		t.Errorf("len(personas) = %d, want 1", len(r.personas))
	}
}

func TestRunner_WithAllPersonas(t *testing.T) {
	r := NewRunner().WithAllPersonas()
	if len(r.personas) < 3 {
		t.Errorf("len(personas) = %d, want >= 3", len(r.personas))
	}
}

func TestRunner_Evaluate(t *testing.T) {
	r := NewRunner().WithPersona("beginner")

	// Evaluate empty path (will find no resources)
	result, err := r.Evaluate("/nonexistent/path")
	if err == nil && result != nil {
		// Should either error or return empty result
		if result.TotalScore > 0 {
			t.Error("expected 0 score for empty evaluation")
		}
	}
}

func TestResult_Score(t *testing.T) {
	result := &Result{
		PersonaResults: []PersonaResult{
			{
				Persona: "sre",
				Score:   80,
				MaxScore: 100,
			},
			{
				Persona: "developer",
				Score:   60,
				MaxScore: 100,
			},
		},
	}

	result.TotalScore = (result.PersonaResults[0].Score + result.PersonaResults[1].Score) / 2
	result.MaxScore = 100

	if result.TotalScore != 70 {
		t.Errorf("TotalScore = %d, want 70", result.TotalScore)
	}
}

func TestCriterionResult(t *testing.T) {
	cr := CriterionResult{
		ID:       "test",
		Name:     "Test Criterion",
		Status:   StatusPass,
		Score:    10,
		MaxScore: 10,
		Message:  "All good",
	}

	if cr.Status != StatusPass {
		t.Errorf("Status = %s, want pass", cr.Status)
	}
}

func TestStatus_Values(t *testing.T) {
	statuses := []Status{StatusPass, StatusPartial, StatusFail, StatusSkip}
	for _, s := range statuses {
		if s == "" {
			t.Error("status should not be empty")
		}
	}
}

func TestResult_PassRate(t *testing.T) {
	result := &Result{
		PersonaResults: []PersonaResult{
			{
				Criteria: []CriterionResult{
					{Status: StatusPass},
					{Status: StatusPass},
					{Status: StatusFail},
					{Status: StatusSkip},
				},
			},
		},
	}

	// Count pass rate
	var passed, total int
	for _, pr := range result.PersonaResults {
		for _, cr := range pr.Criteria {
			if cr.Status != StatusSkip {
				total++
				if cr.Status == StatusPass {
					passed++
				}
			}
		}
	}

	rate := float64(passed) / float64(total) * 100
	if rate < 66 || rate > 67 {
		t.Errorf("pass rate = %.2f, want ~66.67", rate)
	}
}
