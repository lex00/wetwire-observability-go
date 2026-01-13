package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun_Help(t *testing.T) {
	exitCode := run([]string{"help"})
	if exitCode != 0 {
		t.Errorf("run(help) = %d, want 0", exitCode)
	}
}

func TestRun_Version(t *testing.T) {
	exitCode := run([]string{"version"})
	if exitCode != 0 {
		t.Errorf("run(version) = %d, want 0", exitCode)
	}
}

func TestRun_UnknownCommand(t *testing.T) {
	exitCode := run([]string{"unknown"})
	if exitCode != 1 {
		t.Errorf("run(unknown) = %d, want 1", exitCode)
	}
}

func TestRun_NoArgs(t *testing.T) {
	exitCode := run([]string{})
	if exitCode != 0 {
		t.Errorf("run() = %d, want 0", exitCode)
	}
}

func TestListCmd_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	exitCode := listCmd([]string{tmpDir})
	if exitCode != 0 {
		t.Errorf("listCmd(empty) = %d, want 0", exitCode)
	}
}

func TestListCmd_WithResources(t *testing.T) {
	tmpDir := t.TempDir()

	// Create a Go file with resources
	content := `package monitoring

type PrometheusConfig struct{}

var MyConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	exitCode := listCmd([]string{tmpDir})
	if exitCode != 0 {
		t.Errorf("listCmd() = %d, want 0", exitCode)
	}
}

func TestListCmd_JSON(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package monitoring

type PrometheusConfig struct{}

var MyConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	exitCode := listCmd([]string{"--format", "json", tmpDir})
	if exitCode != 0 {
		t.Errorf("listCmd(json) = %d, want 0", exitCode)
	}
}

func TestListCmd_InvalidFormat(t *testing.T) {
	exitCode := listCmd([]string{"--format", "invalid", "."})
	if exitCode != 2 {
		t.Errorf("listCmd(invalid format) = %d, want 2", exitCode)
	}
}

func TestLintCmd_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	exitCode := lintCmd([]string{tmpDir})
	if exitCode != 0 {
		t.Errorf("lintCmd(empty) = %d, want 0", exitCode)
	}
}

func TestLintCmd_WithResources(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package monitoring

type PrometheusConfig struct{}

var MyConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	exitCode := lintCmd([]string{tmpDir})
	if exitCode != 0 {
		t.Errorf("lintCmd() = %d, want 0", exitCode)
	}
}

func TestBuildCmd_EmptyDir(t *testing.T) {
	tmpDir := t.TempDir()
	exitCode := buildCmd([]string{tmpDir})
	if exitCode != 0 {
		t.Errorf("buildCmd(empty) = %d, want 0", exitCode)
	}
}

func TestBuildCmd_WithResources(t *testing.T) {
	tmpDir := t.TempDir()
	outDir := filepath.Join(tmpDir, "out")

	content := `package monitoring

type PrometheusConfig struct{}

var MyConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "config.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	exitCode := buildCmd([]string{"--output", outDir, tmpDir})
	if exitCode != 0 {
		t.Errorf("buildCmd() = %d, want 0", exitCode)
	}

	// Check output file was created
	files, err := os.ReadDir(outDir)
	if err != nil {
		t.Fatalf("ReadDir(%s) error = %v", outDir, err)
	}

	found := false
	for _, f := range files {
		if strings.HasPrefix(f.Name(), "prometheus-") && strings.HasSuffix(f.Name(), ".yml") {
			found = true
			break
		}
	}

	if !found {
		t.Error("Expected prometheus-*.yml file to be created")
	}
}

func TestBuildCmd_InvalidMode(t *testing.T) {
	exitCode := buildCmd([]string{"--mode", "invalid", "."})
	if exitCode != 2 {
		t.Errorf("buildCmd(invalid mode) = %d, want 2", exitCode)
	}
}
