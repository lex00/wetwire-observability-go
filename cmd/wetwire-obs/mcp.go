// Command mcp runs an MCP server that exposes wetwire-obs tools.
//
// This server implements the Model Context Protocol (MCP) using infrastructure
// from github.com/lex00/wetwire-core-go/domain and automatically generates tools
// from the domain.Domain interface implementation.
//
// Tools are automatically registered based on the domain interface:
//   - wetwire_init: Initialize a new wetwire-obs project
//   - wetwire_build: Generate observability configurations from Go packages
//   - wetwire_lint: Lint Go packages for wetwire-obs issues
//   - wetwire_validate: Validate generated configurations
//   - wetwire_list: List discovered resources
//
// Usage:
//
//	wetwire-obs mcp  # Runs on stdio transport
package main

import (
	"context"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-observability-go/domain"
	"github.com/spf13/cobra"
)

func newMCPCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "mcp",
		Short: "Run MCP server for wetwire-obs tools",
		Long: `Run an MCP (Model Context Protocol) server that exposes wetwire-obs tools.

This command starts an MCP server on stdio transport, automatically providing tools
generated from the domain interface for:
  - Initializing projects (wetwire_init)
  - Building observability configurations (wetwire_build)
  - Linting code (wetwire_lint)
  - Validating configurations (wetwire_validate)
  - Listing resources (wetwire_list)

This is typically called by Claude Code or other MCP clients, not directly by users.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runMCPServer()
		},
	}
}

func runMCPServer() error {
	// Create domain instance
	obsDomain := &domain.ObservabilityDomain{}

	// Build MCP server from domain using auto-generation
	server := coredomain.BuildMCPServer(obsDomain)

	// Start the server on stdio transport
	return server.Start(context.Background())
}
