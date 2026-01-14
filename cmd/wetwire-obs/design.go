package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/lex00/wetwire-core-go/kiro"
	"github.com/lex00/wetwire-observability-go/design"
)

func designCmd(args []string) int {
	fs := flag.NewFlagSet("design", flag.ContinueOnError)
	provider := fs.String("provider", "anthropic", "AI provider: anthropic or kiro")
	focus := fs.String("focus", "", "Focus area: prometheus, alertmanager, grafana, rules")
	dryRun := fs.Bool("dry-run", false, "Show prompts without calling API")
	model := fs.String("model", "", "Model to use (default: claude-sonnet-4-20250514)")
	maxTokens := fs.Int("max-tokens", 4096, "Maximum tokens in response")
	contextFile := fs.String("context", "", "File containing existing code context")
	output := fs.String("output", "", "Output file for generated code (default: stdout)")
	timeout := fs.Duration("timeout", 2*time.Minute, "API request timeout")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-obs design [options] <request>")
		fmt.Println()
		fmt.Println("Generate observability configurations using AI.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Providers:")
		fmt.Println("  anthropic  Direct Anthropic API (requires ANTHROPIC_API_KEY)")
		fmt.Println("  kiro       Kiro CLI with MCP server (interactive)")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-obs design 'Add monitoring for an API server'")
		fmt.Println("  wetwire-obs design --focus prometheus 'Add kubernetes discovery'")
		fmt.Println("  wetwire-obs design --provider kiro 'Create SLO dashboard'")
		fmt.Println("  wetwire-obs design --dry-run 'Show prompts only'")
		fmt.Println()
		fmt.Println("Environment:")
		fmt.Println("  ANTHROPIC_API_KEY  Required for anthropic provider")
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	// Get request from remaining args
	request := strings.Join(fs.Args(), " ")
	if request == "" {
		fmt.Fprintln(os.Stderr, "Error: no request provided")
		fmt.Fprintln(os.Stderr, "Usage: wetwire-obs design [options] <request>")
		return 1
	}

	// Build prompt
	pb := design.NewPromptBuilder()

	// Apply focus
	if *focus != "" {
		switch *focus {
		case "prometheus":
			pb.ForPrometheus()
		case "alertmanager":
			pb.ForAlertmanager()
		case "grafana":
			pb.ForGrafana()
		case "rules":
			pb.ForRules()
		default:
			fmt.Fprintf(os.Stderr, "Error: invalid focus %q (use: prometheus, alertmanager, grafana, rules)\n", *focus)
			return 1
		}
	}

	// Add context from file
	if *contextFile != "" {
		content, err := os.ReadFile(*contextFile)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading context file: %v\n", err)
			return 1
		}
		pb.WithContext("existing code", string(content))
	}

	systemPrompt := pb.SystemPrompt()
	userPrompt := pb.BuildUserPrompt(request)

	// Dry run: show prompts and exit
	if *dryRun {
		fmt.Println("=== System Prompt ===")
		fmt.Println(systemPrompt)
		fmt.Println()
		fmt.Println("=== User Prompt ===")
		fmt.Println(userPrompt)
		return 0
	}

	// Handle provider selection
	switch *provider {
	case "kiro":
		return runKiroProvider(systemPrompt, userPrompt, *timeout)
	case "anthropic":
		return runAnthropicProvider(systemPrompt, userPrompt, *model, *maxTokens, *timeout, *output)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown provider %q (use: anthropic, kiro)\n", *provider)
		return 1
	}
}

func runKiroProvider(systemPrompt, userPrompt string, timeout time.Duration) int {
	// Check if kiro is available
	if !kiro.KiroAvailable() {
		fmt.Fprintln(os.Stderr, "Error: kiro-cli not found in PATH")
		fmt.Fprintln(os.Stderr, "Install from: https://github.com/aws/amazon-q-developer-cli")
		return 1
	}

	// Configure kiro with MCP server
	config := kiro.Config{
		AgentName:   "wetwire-observability",
		AgentPrompt: systemPrompt,
		MCPCommand:  "wetwire-obs-mcp",
	}

	// Get working directory
	workDir, err := os.Getwd()
	if err == nil {
		config.WorkDir = workDir
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	fmt.Fprintln(os.Stderr, "Launching Kiro with wetwire-obs MCP server...")

	if err := kiro.Launch(ctx, config, userPrompt); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	return 0
}

func runAnthropicProvider(systemPrompt, userPrompt, model string, maxTokens int, timeout time.Duration, output string) int {
	// Check for API key
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: ANTHROPIC_API_KEY environment variable not set")
		fmt.Fprintln(os.Stderr, "Use --dry-run to preview prompts without API calls")
		fmt.Fprintln(os.Stderr, "Or use --provider kiro for interactive Kiro sessions")
		return 1
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
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Output result
	if output != "" {
		if err := os.WriteFile(output, []byte(result), 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing output file: %v\n", err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "Generated code written to %s\n", output)
	} else {
		fmt.Println(result)
	}

	return 0
}
