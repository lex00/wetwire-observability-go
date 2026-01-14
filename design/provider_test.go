package design

import (
	"context"
	"errors"
	"testing"
)

func TestNewAnthropicProvider(t *testing.T) {
	p := NewAnthropicProvider("")
	if p == nil {
		t.Fatal("NewAnthropicProvider returned nil")
	}
}

func TestAnthropicProvider_Name(t *testing.T) {
	p := NewAnthropicProvider("")
	if p.Name() != "anthropic" {
		t.Errorf("Name() = %q, want anthropic", p.Name())
	}
}

func TestAnthropicProvider_Generate_NoAPIKey(t *testing.T) {
	p := NewAnthropicProvider("")
	_, err := p.Generate(context.Background(), "system", "user")
	if err == nil {
		t.Error("expected error when API key not set")
	}
}

func TestMockProvider(t *testing.T) {
	p := &MockProvider{
		Response: "generated code",
	}

	result, err := p.Generate(context.Background(), "system", "user")
	if err != nil {
		t.Errorf("Generate() error = %v", err)
	}
	if result != "generated code" {
		t.Errorf("Generate() = %q", result)
	}
}

func TestMockProvider_Error(t *testing.T) {
	p := &MockProvider{
		Err: errors.New("mock error"),
	}

	_, err := p.Generate(context.Background(), "system", "user")
	if err == nil {
		t.Error("expected error from mock")
	}
}

func TestProviderConfig(t *testing.T) {
	cfg := DefaultProviderConfig()
	if cfg.Model == "" {
		t.Error("default config should have model set")
	}
	if cfg.MaxTokens == 0 {
		t.Error("default config should have max tokens set")
	}
}

func TestProviderConfig_WithModel(t *testing.T) {
	cfg := DefaultProviderConfig().WithModel("claude-3-sonnet")
	if cfg.Model != "claude-3-sonnet" {
		t.Errorf("Model = %q", cfg.Model)
	}
}

func TestProviderConfig_WithMaxTokens(t *testing.T) {
	cfg := DefaultProviderConfig().WithMaxTokens(8192)
	if cfg.MaxTokens != 8192 {
		t.Errorf("MaxTokens = %d", cfg.MaxTokens)
	}
}
