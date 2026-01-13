package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestValidateCmd_Help(t *testing.T) {
	// Just test that help works
	code := validateCmd([]string{"--help"})
	if code != 0 {
		t.Errorf("validateCmd(--help) = %d, want 0", code)
	}
}

func TestValidateCmd_EmptyDir(t *testing.T) {
	// Create temp directory
	tmpDir, err := os.MkdirTemp("", "validate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Run validate on empty dir - should gracefully skip all validators
	code := validateCmd([]string{"-output", tmpDir})
	// Should be 0 since all validators are skipped gracefully
	if code != 0 {
		t.Errorf("validateCmd on empty dir = %d, want 0", code)
	}
}

func TestValidateCmd_SkipFlags(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Test all skip flags work
	code := validateCmd([]string{
		"-output", tmpDir,
		"-skip-kubeconform",
		"-skip-promtool",
		"-skip-amtool",
	})
	if code != 0 {
		t.Errorf("validateCmd with skip flags = %d, want 0", code)
	}
}

func TestValidationResults_HasErrors(t *testing.T) {
	// No errors when all skipped
	r := &validationResults{
		kubeconform: &validationResult{skipped: true},
		promtool:    &validationResult{skipped: true},
		amtool:      &validationResult{skipped: true},
	}
	if r.hasErrors() {
		t.Error("hasErrors() = true for all skipped, want false")
	}

	// No errors when all pass
	r = &validationResults{
		kubeconform: &validationResult{success: true},
		promtool:    &validationResult{success: true},
		amtool:      &validationResult{success: true},
	}
	if r.hasErrors() {
		t.Error("hasErrors() = true for all success, want false")
	}

	// Has errors when one fails
	r = &validationResults{
		kubeconform: &validationResult{success: true},
		promtool:    &validationResult{success: false},
		amtool:      &validationResult{skipped: true},
	}
	if !r.hasErrors() {
		t.Error("hasErrors() = false with failure, want true")
	}
}

func TestValidateKubeconform_SkippedWhenNotInstalled(t *testing.T) {
	// Save PATH and restore after test
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	// Set PATH to empty to simulate kubeconform not being installed
	os.Setenv("PATH", "")

	result := validateKubeconform("nonexistent", false)
	if !result.skipped {
		t.Error("expected kubeconform to be skipped when not installed")
	}
}

func TestValidatePromtool_SkippedWhenNotInstalled(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	os.Setenv("PATH", "")

	result := validatePromtool("nonexistent", false)
	if !result.skipped {
		t.Error("expected promtool to be skipped when not installed")
	}
}

func TestValidateAmtool_SkippedWhenNotInstalled(t *testing.T) {
	oldPath := os.Getenv("PATH")
	defer os.Setenv("PATH", oldPath)

	os.Setenv("PATH", "")

	result := validateAmtool("nonexistent", false)
	if !result.skipped {
		t.Error("expected amtool to be skipped when not installed")
	}
}

func TestValidateKubeconform_NoManifests(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create empty operator dir
	operatorDir := filepath.Join(tmpDir, "operator")
	if err := os.MkdirAll(operatorDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Should skip gracefully when no manifests found
	result := validateKubeconform(tmpDir, false)
	if !result.skipped {
		// Either skipped (no manifests) or skipped (kubeconform not installed)
		if result.success {
			t.Error("expected validation to skip or fail gracefully with no manifests")
		}
	}
}

func TestValidatePromtool_NoConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	result := validatePromtool(tmpDir, false)
	if !result.skipped {
		t.Error("expected promtool validation to be skipped when no config files exist")
	}
}

func TestValidateAmtool_NoConfig(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "validate-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	result := validateAmtool(tmpDir, false)
	if !result.skipped {
		t.Error("expected amtool validation to be skipped when no alertmanager.yml exists")
	}
}
