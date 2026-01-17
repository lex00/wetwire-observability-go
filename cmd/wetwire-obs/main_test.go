package main

import (
	"testing"

	"github.com/lex00/wetwire-observability-go/domain"
)

func TestMain_Version(t *testing.T) {
	// Test that we can create the domain
	d := &domain.ObservabilityDomain{}
	if d.Name() != "observability" {
		t.Errorf("expected domain name 'observability', got %q", d.Name())
	}
}

func TestMain_CreateCommand(t *testing.T) {
	// Test that we can create the root command
	d := &domain.ObservabilityDomain{}
	cmd := domain.CreateRootCommand(d)
	if cmd == nil {
		t.Fatal("expected root command, got nil")
	}
}

// Legacy command tests removed - domain now handles command dispatch.
// Command functionality is tested in the domain package.
