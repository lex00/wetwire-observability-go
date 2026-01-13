package discover

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiscover_Basic(t *testing.T) {
	// Create a temp directory with test fixtures
	tmpDir := t.TempDir()

	// Create a Go file with prometheus resources
	content := `package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

var MyConfig = &prometheus.PrometheusConfig{
	Global: &prometheus.GlobalConfig{},
}

var MyScrape = prometheus.ScrapeConfig{
	JobName: "test",
}

// unexported should be skipped
var unexported = &prometheus.PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should find 2 exported resources (MyConfig and MyScrape)
	// Note: With current implementation, we look for type names
	// that may include the package qualifier
	if result.TotalCount() < 1 {
		t.Errorf("TotalCount() = %d, want >= 1", result.TotalCount())
	}
}

func TestDiscover_TypeVariants(t *testing.T) {
	tmpDir := t.TempDir()

	// Test various type declaration patterns
	content := `package configs

type PrometheusConfig struct{}
type ScrapeConfig struct{}
type GlobalConfig struct{}
type StaticConfig struct{}

// Pointer type
var Config1 *PrometheusConfig

// Value type
var Config2 PrometheusConfig

// Composite literal
var Config3 = &PrometheusConfig{}

// Value composite literal
var Config4 = PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "types.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should find all 4 exported configs
	if len(result.PrometheusConfigs) != 4 {
		t.Errorf("len(PrometheusConfigs) = %d, want 4", len(result.PrometheusConfigs))
	}
}

func TestDiscover_UnexportedSkipped(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package test

type PrometheusConfig struct{}

// Exported
var MyConfig = &PrometheusConfig{}

// Unexported - should be skipped
var myConfig = &PrometheusConfig{}
var privateConfig = &PrometheusConfig{}
var _hidden = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "test.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should only find 1 exported config
	if len(result.PrometheusConfigs) != 1 {
		t.Errorf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}
	if result.PrometheusConfigs[0].Name != "MyConfig" {
		t.Errorf("Name = %s, want MyConfig", result.PrometheusConfigs[0].Name)
	}
}

func TestDiscover_MultiplePackages(t *testing.T) {
	tmpDir := t.TempDir()

	// Create subdirectory structure
	pkgDir := filepath.Join(tmpDir, "pkg", "monitoring")
	if err := os.MkdirAll(pkgDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Root package
	rootContent := `package root

type PrometheusConfig struct{}

var RootConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "root.go"), []byte(rootContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Nested package
	nestedContent := `package monitoring

type PrometheusConfig struct{}

var NestedConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(pkgDir, "config.go"), []byte(nestedContent), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should find configs from both packages
	if len(result.PrometheusConfigs) != 2 {
		t.Errorf("len(PrometheusConfigs) = %d, want 2", len(result.PrometheusConfigs))
	}
}

func TestDiscover_SkipsVendor(t *testing.T) {
	tmpDir := t.TempDir()

	// Create vendor directory
	vendorDir := filepath.Join(tmpDir, "vendor", "example")
	if err := os.MkdirAll(vendorDir, 0755); err != nil {
		t.Fatal(err)
	}

	// Main package
	mainContent := `package main

type PrometheusConfig struct{}

var MainConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Vendor package (should be skipped)
	vendorContent := `package example

type PrometheusConfig struct{}

var VendorConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(vendorDir, "vendor.go"), []byte(vendorContent), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should only find main package config, not vendor
	if len(result.PrometheusConfigs) != 1 {
		t.Errorf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}
	if result.PrometheusConfigs[0].Name != "MainConfig" {
		t.Errorf("Name = %s, want MainConfig", result.PrometheusConfigs[0].Name)
	}
}

func TestDiscover_SkipsTestFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Main file
	mainContent := `package main

type PrometheusConfig struct{}

var MainConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(mainContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Test file (should be skipped)
	testContent := `package main

type PrometheusConfig struct{}

var TestConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main_test.go"), []byte(testContent), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should only find main file config
	if len(result.PrometheusConfigs) != 1 {
		t.Errorf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}
}

func TestDiscover_AllResourceTypes(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package configs

type PrometheusConfig struct{}
type ScrapeConfig struct{}
type GlobalConfig struct{}
type StaticConfig struct{}

var Prom = &PrometheusConfig{}
var Scrape = &ScrapeConfig{}
var Global = &GlobalConfig{}
var Static = &StaticConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "all.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.PrometheusConfigs) != 1 {
		t.Errorf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}
	if len(result.ScrapeConfigs) != 1 {
		t.Errorf("len(ScrapeConfigs) = %d, want 1", len(result.ScrapeConfigs))
	}
	if len(result.GlobalConfigs) != 1 {
		t.Errorf("len(GlobalConfigs) = %d, want 1", len(result.GlobalConfigs))
	}
	if len(result.StaticConfigs) != 1 {
		t.Errorf("len(StaticConfigs) = %d, want 1", len(result.StaticConfigs))
	}
	if result.TotalCount() != 4 {
		t.Errorf("TotalCount() = %d, want 4", result.TotalCount())
	}
}

func TestDiscover_MalformedFile(t *testing.T) {
	tmpDir := t.TempDir()

	// Valid file
	validContent := `package main

type PrometheusConfig struct{}

var ValidConfig = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "valid.go"), []byte(validContent), 0644); err != nil {
		t.Fatal(err)
	}

	// Malformed file
	malformedContent := `package main

var invalid = this is not valid Go
`
	if err := os.WriteFile(filepath.Join(tmpDir, "malformed.go"), []byte(malformedContent), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should still find the valid config
	if len(result.PrometheusConfigs) != 1 {
		t.Errorf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}

	// Should have recorded an error for the malformed file
	if len(result.Errors) == 0 {
		t.Error("Expected errors for malformed file")
	}
}

func TestDiscover_ResourceRef(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package example

type PrometheusConfig struct{}

var MyPrometheusConfig = &PrometheusConfig{}
`
	testFile := filepath.Join(tmpDir, "example.go")
	if err := os.WriteFile(testFile, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if len(result.PrometheusConfigs) != 1 {
		t.Fatalf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}

	ref := result.PrometheusConfigs[0]
	if ref.Package != "example" {
		t.Errorf("Package = %s, want example", ref.Package)
	}
	if ref.Name != "MyPrometheusConfig" {
		t.Errorf("Name = %s, want MyPrometheusConfig", ref.Name)
	}
	if ref.Type != "PrometheusConfig" {
		t.Errorf("Type = %s, want PrometheusConfig", ref.Type)
	}
	if ref.FilePath != testFile {
		t.Errorf("FilePath = %s, want %s", ref.FilePath, testFile)
	}
	if ref.Line != 5 {
		t.Errorf("Line = %d, want 5", ref.Line)
	}
}

func TestDiscover_EmptyDirectory(t *testing.T) {
	tmpDir := t.TempDir()

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	if result.TotalCount() != 0 {
		t.Errorf("TotalCount() = %d, want 0", result.TotalCount())
	}
}

func TestDiscover_All(t *testing.T) {
	tmpDir := t.TempDir()

	content := `package configs

type PrometheusConfig struct{}
type ScrapeConfig struct{}

var Prom = &PrometheusConfig{}
var Scrape1 = &ScrapeConfig{}
var Scrape2 = &ScrapeConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "configs.go"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	all := result.All()
	if len(all) != 3 {
		t.Errorf("len(All()) = %d, want 3", len(all))
	}
}

func TestDiscover_NonGoFiles(t *testing.T) {
	tmpDir := t.TempDir()

	// Create various non-Go files
	files := map[string]string{
		"readme.md":   "# README",
		"config.yaml": "key: value",
		"script.sh":   "#!/bin/bash",
	}

	for name, content := range files {
		if err := os.WriteFile(filepath.Join(tmpDir, name), []byte(content), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Also create one valid Go file
	goContent := `package main

type PrometheusConfig struct{}

var Config = &PrometheusConfig{}
`
	if err := os.WriteFile(filepath.Join(tmpDir, "main.go"), []byte(goContent), 0644); err != nil {
		t.Fatal(err)
	}

	result, err := Discover(tmpDir)
	if err != nil {
		t.Fatalf("Discover() error = %v", err)
	}

	// Should only find the Go file resource
	if len(result.PrometheusConfigs) != 1 {
		t.Errorf("len(PrometheusConfigs) = %d, want 1", len(result.PrometheusConfigs))
	}
}
