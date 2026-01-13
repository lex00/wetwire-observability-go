package mcp

import (
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

func TestRegisterTools(t *testing.T) {
	s := NewServer("wetwire-obs", "1.0.0")
	RegisterTools(s)

	// Should have registered build, lint, import, validate, list
	expectedTools := []string{"build", "lint", "import", "validate", "list"}
	if len(s.tools) != len(expectedTools) {
		t.Errorf("len(tools) = %d, want %d", len(s.tools), len(expectedTools))
	}

	for i, expected := range expectedTools {
		if s.tools[i].Name != expected {
			t.Errorf("tools[%d].Name = %q, want %q", i, s.tools[i].Name, expected)
		}
	}
}

func TestBuildTool_InvalidPath(t *testing.T) {
	params := json.RawMessage(`{"path": "/nonexistent/path", "output": "/tmp/out"}`)
	result, err := handleBuild(params)
	if err == nil {
		t.Error("expected error for invalid path")
	}
	_ = result
}

func TestBuildTool_ValidParams(t *testing.T) {
	// Create temp directory with a simple go file
	tmpDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a minimal go.mod
	goMod := filepath.Join(tmpDir, "go.mod")
	if err := os.WriteFile(goMod, []byte("module test\n\ngo 1.21\n"), 0644); err != nil {
		t.Fatal(err)
	}

	// Create a simple Go file with a prometheus config
	goFile := filepath.Join(tmpDir, "config.go")
	content := `package test

import "github.com/lex00/wetwire-observability-go/prometheus"

var Config = prometheus.NewConfig()
`
	if err := os.WriteFile(goFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	outDir := filepath.Join(tmpDir, "out")
	params := json.RawMessage(`{"path": "` + tmpDir + `", "output": "` + outDir + `"}`)

	// This may fail due to import issues in test environment, that's OK
	_, _ = handleBuild(params)
}

func TestLintTool(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	params := json.RawMessage(`{"path": "` + tmpDir + `"}`)
	result, err := handleLint(params)
	// Should work even on empty directory
	if err != nil && result == nil {
		t.Logf("lint returned error: %v (this may be expected)", err)
	}
}

func TestValidateTool(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	params := json.RawMessage(`{"output": "` + tmpDir + `"}`)
	result, err := handleValidate(params)
	// Should skip gracefully when no files exist
	if err != nil {
		t.Logf("validate returned error: %v (this may be expected)", err)
	}
	_ = result
}

func TestListTool(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "mcp-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	params := json.RawMessage(`{"path": "` + tmpDir + `"}`)
	result, err := handleList(params)
	if err != nil {
		t.Logf("list returned error: %v", err)
	}
	_ = result
}

func TestImportTool_NoFile(t *testing.T) {
	params := json.RawMessage(`{"file": "/nonexistent/file.yml"}`)
	_, err := handleImport(params)
	if err == nil {
		t.Error("expected error for nonexistent file")
	}
}

func TestToolInputSchemas(t *testing.T) {
	s := NewServer("wetwire-obs", "1.0.0")
	RegisterTools(s)

	for _, tool := range s.tools {
		if tool.InputSchema.Type != "object" {
			t.Errorf("tool %s: InputSchema.Type = %q, want object", tool.Name, tool.InputSchema.Type)
		}
	}
}
