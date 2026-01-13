// Package main provides the MCP server for wetwire-obs.
package main

import (
	"fmt"
	"os"

	"github.com/lex00/wetwire-observability-go/mcp"
)

// Version is set by the build process
var Version = "dev"

func main() {
	server := mcp.NewServer("wetwire-obs", Version)
	mcp.RegisterTools(server)

	if err := server.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
