package design

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

// Provider is an interface for AI text generation providers.
type Provider interface {
	// Name returns the provider name.
	Name() string

	// Generate generates text given a system prompt and user message.
	Generate(ctx context.Context, system, user string) (string, error)
}

// ProviderConfig contains configuration for AI providers.
type ProviderConfig struct {
	Model     string
	MaxTokens int
	APIKey    string
}

// DefaultProviderConfig returns the default provider configuration.
func DefaultProviderConfig() *ProviderConfig {
	return &ProviderConfig{
		Model:     "claude-sonnet-4-20250514",
		MaxTokens: 4096,
		APIKey:    os.Getenv("ANTHROPIC_API_KEY"),
	}
}

// WithModel sets the model.
func (c *ProviderConfig) WithModel(model string) *ProviderConfig {
	c.Model = model
	return c
}

// WithMaxTokens sets the max tokens.
func (c *ProviderConfig) WithMaxTokens(tokens int) *ProviderConfig {
	c.MaxTokens = tokens
	return c
}

// WithAPIKey sets the API key.
func (c *ProviderConfig) WithAPIKey(key string) *ProviderConfig {
	c.APIKey = key
	return c
}

// AnthropicProvider implements Provider using the Anthropic API.
type AnthropicProvider struct {
	config *ProviderConfig
	client *http.Client
}

// NewAnthropicProvider creates a new AnthropicProvider.
func NewAnthropicProvider(apiKey string) *AnthropicProvider {
	cfg := DefaultProviderConfig()
	if apiKey != "" {
		cfg.APIKey = apiKey
	}
	return &AnthropicProvider{
		config: cfg,
		client: &http.Client{},
	}
}

// NewAnthropicProviderWithConfig creates a new AnthropicProvider with custom config.
func NewAnthropicProviderWithConfig(config *ProviderConfig) *AnthropicProvider {
	return &AnthropicProvider{
		config: config,
		client: &http.Client{},
	}
}

// Name returns the provider name.
func (p *AnthropicProvider) Name() string {
	return "anthropic"
}

// Generate generates text using the Anthropic API.
func (p *AnthropicProvider) Generate(ctx context.Context, system, user string) (string, error) {
	if p.config.APIKey == "" {
		return "", fmt.Errorf("ANTHROPIC_API_KEY not set")
	}

	reqBody := anthropicRequest{
		Model:     p.config.Model,
		MaxTokens: p.config.MaxTokens,
		System:    system,
		Messages: []anthropicMessage{
			{Role: "user", Content: user},
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", p.config.APIKey)
	req.Header.Set("anthropic-version", "2023-06-01")

	resp, err := p.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		var errResp anthropicError
		json.NewDecoder(resp.Body).Decode(&errResp)
		return "", fmt.Errorf("API error (%d): %s", resp.StatusCode, errResp.Error.Message)
	}

	var result anthropicResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(result.Content) == 0 {
		return "", fmt.Errorf("empty response")
	}

	return result.Content[0].Text, nil
}

// Anthropic API types

type anthropicRequest struct {
	Model     string             `json:"model"`
	MaxTokens int                `json:"max_tokens"`
	System    string             `json:"system,omitempty"`
	Messages  []anthropicMessage `json:"messages"`
}

type anthropicMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type anthropicResponse struct {
	ID      string            `json:"id"`
	Type    string            `json:"type"`
	Role    string            `json:"role"`
	Content []anthropicContent `json:"content"`
	Model   string            `json:"model"`
	Usage   anthropicUsage    `json:"usage"`
}

type anthropicContent struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type anthropicUsage struct {
	InputTokens  int `json:"input_tokens"`
	OutputTokens int `json:"output_tokens"`
}

type anthropicError struct {
	Type  string `json:"type"`
	Error struct {
		Type    string `json:"type"`
		Message string `json:"message"`
	} `json:"error"`
}

// MockProvider is a mock provider for testing.
type MockProvider struct {
	Response string
	Err      error
}

// Name returns the provider name.
func (p *MockProvider) Name() string {
	return "mock"
}

// Generate returns the mock response.
func (p *MockProvider) Generate(ctx context.Context, system, user string) (string, error) {
	if p.Err != nil {
		return "", p.Err
	}
	return p.Response, nil
}
