// Package main provides the wetwire-obs CLI for observability configuration synthesis.
package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("wetwire-obs - Prometheus, Alertmanager, and Grafana configuration synthesis")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  wetwire-obs <command> [options]")
		fmt.Println()
		fmt.Println("Commands:")
		fmt.Println("  build     Generate configuration files")
		fmt.Println("  lint      Check code quality")
		fmt.Println("  import    Convert existing configs to Go")
		fmt.Println("  validate  Run external validators")
		fmt.Println("  list      List discovered resources")
		fmt.Println("  mcp       Start MCP server")
		fmt.Println("  version   Show version information")
		os.Exit(0)
	}

	cmd := os.Args[1]
	switch cmd {
	case "version":
		fmt.Println("wetwire-obs v0.0.1")
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		os.Exit(1)
	}
}
