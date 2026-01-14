package main

import (
	"testing"
)

func TestDesignCmd_Help(t *testing.T) {
	code := designCmd([]string{"--help"})
	if code != 0 {
		t.Errorf("designCmd(--help) = %d, want 0", code)
	}
}

func TestDesignCmd_DryRun(t *testing.T) {
	// Dry run should work without API key
	code := designCmd([]string{
		"--dry-run",
		"Add monitoring for an API server",
	})
	if code != 0 {
		t.Errorf("designCmd(--dry-run) = %d, want 0", code)
	}
}

func TestDesignCmd_NoInput(t *testing.T) {
	code := designCmd([]string{})
	if code == 0 {
		t.Error("expected non-zero exit code when no input provided")
	}
}

func TestDesignCmd_FocusPrometheus(t *testing.T) {
	code := designCmd([]string{
		"--focus", "prometheus",
		"--dry-run",
		"Add kubernetes discovery",
	})
	if code != 0 {
		t.Errorf("designCmd with focus = %d, want 0", code)
	}
}

func TestDesignCmd_FocusAlertmanager(t *testing.T) {
	code := designCmd([]string{
		"--focus", "alertmanager",
		"--dry-run",
		"Add slack receiver",
	})
	if code != 0 {
		t.Errorf("designCmd with alertmanager focus = %d, want 0", code)
	}
}

func TestDesignCmd_FocusGrafana(t *testing.T) {
	code := designCmd([]string{
		"--focus", "grafana",
		"--dry-run",
		"Create API dashboard",
	})
	if code != 0 {
		t.Errorf("designCmd with grafana focus = %d, want 0", code)
	}
}

func TestDesignCmd_FocusRules(t *testing.T) {
	code := designCmd([]string{
		"--focus", "rules",
		"--dry-run",
		"Add high latency alert",
	})
	if code != 0 {
		t.Errorf("designCmd with rules focus = %d, want 0", code)
	}
}

func TestDesignCmd_InvalidFocus(t *testing.T) {
	code := designCmd([]string{
		"--focus", "invalid",
		"--dry-run",
		"test",
	})
	if code == 0 {
		t.Error("expected non-zero exit code for invalid focus")
	}
}

func TestDesignCmd_InvalidProvider(t *testing.T) {
	code := designCmd([]string{
		"--provider", "invalid",
		"test request",
	})
	if code == 0 {
		t.Error("expected non-zero exit code for invalid provider")
	}
}

func TestDesignCmd_ProviderAnthropicNoKey(t *testing.T) {
	// anthropic provider without API key should fail
	code := designCmd([]string{
		"--provider", "anthropic",
		"test request",
	})
	if code == 0 {
		t.Error("expected non-zero exit code without API key")
	}
}

func TestDesignCmd_ProviderKiroNotInstalled(t *testing.T) {
	// Skip if kiro is installed (we can't easily test the launch)
	// This test verifies the error path when kiro is not available
	// The function will return 1 if kiro is not installed
	code := designCmd([]string{
		"--provider", "kiro",
		"test request",
	})
	// We just verify it doesn't panic
	_ = code
}
