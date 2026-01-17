package main

import (
	"bytes"
	"testing"
)

func TestDesignCmd_Help(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{"--help"})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	// Help returns nil, not an error
	_ = cmd.Execute()
}

func TestDesignCmd_DryRun(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--dry-run",
		"Add monitoring for an API server",
	})
	out := &bytes.Buffer{}
	cmd.SetOut(out)
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("designCmd(--dry-run) error = %v", err)
	}
}

func TestDesignCmd_NoInput(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error when no input provided")
	}
}

func TestDesignCmd_FocusPrometheus(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--focus", "prometheus",
		"--dry-run",
		"Add kubernetes discovery",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("designCmd with focus = %v", err)
	}
}

func TestDesignCmd_FocusAlertmanager(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--focus", "alertmanager",
		"--dry-run",
		"Add slack receiver",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("designCmd with alertmanager focus = %v", err)
	}
}

func TestDesignCmd_FocusGrafana(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--focus", "grafana",
		"--dry-run",
		"Create API dashboard",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("designCmd with grafana focus = %v", err)
	}
}

func TestDesignCmd_FocusRules(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--focus", "rules",
		"--dry-run",
		"Add high latency alert",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err != nil {
		t.Errorf("designCmd with rules focus = %v", err)
	}
}

func TestDesignCmd_InvalidFocus(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--focus", "invalid",
		"--dry-run",
		"test",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid focus")
	}
}

func TestDesignCmd_InvalidProvider(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--provider", "invalid",
		"test request",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error for invalid provider")
	}
}

func TestDesignCmd_ProviderAnthropicNoKey(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--provider", "anthropic",
		"test request",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	err := cmd.Execute()
	if err == nil {
		t.Error("expected error without API key")
	}
}

func TestDesignCmd_ProviderKiroNotInstalled(t *testing.T) {
	cmd := newDesignCmd()
	cmd.SetArgs([]string{
		"--provider", "kiro",
		"test request",
	})
	cmd.SetOut(&bytes.Buffer{})
	cmd.SetErr(&bytes.Buffer{})
	// We just verify it doesn't panic
	_ = cmd.Execute()
}
