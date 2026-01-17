package main

import (
	"bytes"
	"testing"
)

func TestTestCmd_Help(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{"--help"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Execute()
}

func TestTestCmd_ListPersonas(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{"--list-personas", "."})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("testCmd(--list-personas) error = %v", err)
	}
}

func TestTestCmd_NoPath(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no path provided")
	}
}

func TestTestCmd_WithPersona(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{
		"--persona", "beginner",
		"/nonexistent/path",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	// Should not crash, may return error for path
	_ = cmd.Execute()
}

func TestTestCmd_JSONOutput(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{
		"--format", "json",
		"--persona", "beginner",
		"/nonexistent/path",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Execute()
}

func TestTestCmd_InvalidPersona(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{
		"--persona", "invalid",
		"/some/path",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid persona")
	}
}

func TestTestCmd_AllPersonas(t *testing.T) {
	cmd := newTestCmd()
	cmd.SetArgs([]string{
		"--all",
		"/nonexistent/path",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	_ = cmd.Execute()
}
