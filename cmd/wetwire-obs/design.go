// Command design generates observability configurations using AI.
package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/lex00/wetwire-core-go/kiro"
	"github.com/lex00/wetwire-observability-go/design"
	"github.com/spf13/cobra"
)

func newDesignCmd() *cobra.Command {
	var (
		provider    string
		focus       string
		dryRun      bool
		model       string
		maxTokens   int
		contextFile string
		output      string
		timeout     time.Duration
	)

	cmd := &cobra.Command{
		Use:   "design <request>",
		Short: "Generate observability configurations using AI",
		Long: `Design generates observability configurations using AI providers.

Providers:
  - anthropic: Direct Anthropic API (requires ANTHROPIC_API_KEY)
  - kiro: Kiro CLI with MCP server (interactive)

Focus areas:
  - prometheus: Prometheus configuration
  - alertmanager: Alertmanager configuration
  - grafana: Grafana dashboards
  - rules: Alerting and recording rules

Examples:
  wetwire-obs design 'Add monitoring for an API server'
  wetwire-obs design --focus prometheus 'Add kubernetes discovery'
  wetwire-obs design --provider kiro 'Create SLO dashboard'
  wetwire-obs design --dry-run 'Show prompts only'`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			request := args[0]
			return runDesign(request, provider, focus, dryRun, model, maxTokens, contextFile, output, timeout)
		},
	}

	cmd.Flags().StringVar(&provider, "provider", "anthropic", "AI provider: anthropic or kiro")
	cmd.Flags().StringVar(&focus, "focus", "", "Focus area: prometheus, alertmanager, grafana, rules")
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "Show prompts without calling API")
	cmd.Flags().StringVar(&model, "model", "", "Model to use (default: claude-sonnet-4-20250514)")
	cmd.Flags().IntVar(&maxTokens, "max-tokens", 4096, "Maximum tokens in response")
	cmd.Flags().StringVar(&contextFile, "context", "", "File containing existing code context")
	cmd.Flags().StringVarP(&output, "output", "o", "", "Output file for generated code")
	cmd.Flags().DurationVar(&timeout, "timeout", 2*time.Minute, "API request timeout")

	return cmd
}

func runDesign(request, provider, focus string, dryRun bool, model string, maxTokens int, contextFile, output string, timeout time.Duration) error {
	// Build prompt
	pb := design.NewPromptBuilder()

	// Apply focus
	if focus != "" {
		switch focus {
		case "prometheus":
			pb.ForPrometheus()
		case "alertmanager":
			pb.ForAlertmanager()
		case "grafana":
			pb.ForGrafana()
		case "rules":
			pb.ForRules()
		default:
			return fmt.Errorf("invalid focus %q (use: prometheus, alertmanager, grafana, rules)", focus)
		}
	}

	// Add context from file
	if contextFile != "" {
		content, err := os.ReadFile(contextFile)
		if err != nil {
			return fmt.Errorf("reading context file: %w", err)
		}
		pb.WithContext("existing code", string(content))
	}

	systemPrompt := pb.SystemPrompt()
	userPrompt := pb.BuildUserPrompt(request)

	// Dry run: show prompts and exit
	if dryRun {
		fmt.Println("=== System Prompt ===")
		fmt.Println(systemPrompt)
		fmt.Println()
		fmt.Println("=== User Prompt ===")
		fmt.Println(userPrompt)
		return nil
	}

	// Handle provider selection
	switch provider {
	case "kiro":
		return runDesignKiro(systemPrompt, userPrompt, timeout)
	case "anthropic":
		return runDesignAnthropic(systemPrompt, userPrompt, model, maxTokens, timeout, output)
	default:
		return fmt.Errorf("unknown provider %q (use: anthropic, kiro)", provider)
	}
}

func runDesignKiro(systemPrompt, userPrompt string, timeout time.Duration) error {
	// Check if kiro is available
	if !kiro.KiroAvailable() {
		return fmt.Errorf("kiro-cli not found in PATH\nInstall from: https://github.com/aws/amazon-q-developer-cli")
	}

	// Configure kiro with MCP server
	config := kiro.Config{
		AgentName:   "wetwire-observability",
		AgentPrompt: systemPrompt,
		MCPCommand:  "wetwire-obs",
		MCPArgs:     []string{"mcp"},
	}

	// Get working directory
	if workDir, err := os.Getwd(); err == nil {
		config.WorkDir = workDir
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Fprintln(os.Stderr, "Launching Kiro with wetwire-obs MCP server...")

	return kiro.Launch(ctx, config, userPrompt)
}

func runDesignAnthropic(systemPrompt, userPrompt, model string, maxTokens int, timeout time.Duration, output string) error {
	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY environment variable not set\nUse --dry-run to preview prompts or --provider kiro for Kiro sessions")
	}

	// Configure provider
	config := design.DefaultProviderConfig().WithAPIKey(apiKey)
	if model != "" {
		config.WithModel(model)
	}
	if maxTokens > 0 {
		config.WithMaxTokens(maxTokens)
	}

	anthropicProvider := design.NewAnthropicProviderWithConfig(config)

	// Generate
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Fprintln(os.Stderr, "Generating configuration...")

	result, err := anthropicProvider.Generate(ctx, systemPrompt, userPrompt)
	if err != nil {
		return err
	}

	// Output result
	if output != "" {
		if err := os.WriteFile(output, []byte(result), 0644); err != nil {
			return fmt.Errorf("writing output file: %w", err)
		}
		fmt.Fprintf(os.Stderr, "Generated code written to %s\n", output)
	} else {
		fmt.Println(result)
	}

	return nil
}
