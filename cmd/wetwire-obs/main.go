// Package main provides the wetwire-obs CLI for observability configuration synthesis.
package main

import (
	"fmt"
	"os"

	"github.com/lex00/wetwire-observability-go/domain"
)

// Version is set by the build process
var Version = "dev"

func main() {
	// Set version before creating command
	domain.Version = Version

	// Create domain and root command
	d := &domain.ObservabilityDomain{}
	cmd := domain.CreateRootCommand(d)

	// Add observability-specific commands
	cmd.AddCommand(newDesignCmd())
	cmd.AddCommand(newTestCmd())
	cmd.AddCommand(newDiffCmd())
	cmd.AddCommand(newWatchCmd())
	cmd.AddCommand(newMCPCmd())

	// Execute
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
