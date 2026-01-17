// Package main provides the MCP server for wetwire-obs.
package main

import (
	"context"
	"fmt"
	"os"

	"github.com/lex00/wetwire-core-go/domain"
	observability "github.com/lex00/wetwire-observability-go/domain"
)

// Version is set by the build process
var Version = "dev"

func main() {
	// Set version before creating server
	observability.Version = Version

	// Build MCP server using domain.BuildMCPServer()
	server := domain.BuildMCPServer(&observability.ObservabilityDomain{})

	// Start the server
	ctx := context.Background()
	if err := server.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
