package main

import (
	"testing"
)

func TestTestCmd_Help(t *testing.T) {
	code := testCmd([]string{"--help"})
	if code != 0 {
		t.Errorf("testCmd(--help) = %d, want 0", code)
	}
}

func TestTestCmd_ListPersonas(t *testing.T) {
	code := testCmd([]string{"--list-personas"})
	if code != 0 {
		t.Errorf("testCmd(--list-personas) = %d, want 0", code)
	}
}

func TestTestCmd_NoPath(t *testing.T) {
	code := testCmd([]string{})
	if code == 0 {
		t.Error("expected non-zero exit code when no path provided")
	}
}

func TestTestCmd_WithPersona(t *testing.T) {
	// Run with beginner persona on nonexistent path
	code := testCmd([]string{
		"--persona", "beginner",
		"/nonexistent/path",
	})
	// Should not crash, may return error for path
	_ = code
}

func TestTestCmd_JSONOutput(t *testing.T) {
	code := testCmd([]string{
		"--format", "json",
		"--persona", "beginner",
		"/nonexistent/path",
	})
	_ = code
}

func TestTestCmd_InvalidPersona(t *testing.T) {
	code := testCmd([]string{
		"--persona", "invalid",
		"/some/path",
	})
	if code == 0 {
		t.Error("expected non-zero exit code for invalid persona")
	}
}

func TestTestCmd_AllPersonas(t *testing.T) {
	code := testCmd([]string{
		"--all",
		"/nonexistent/path",
	})
	_ = code
}
