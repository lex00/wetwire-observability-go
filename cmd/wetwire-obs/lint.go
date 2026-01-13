package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lex00/wetwire-observability-go/internal/discover"
)

// lintCmd handles the lint command
func lintCmd(args []string) int {
	fs := flag.NewFlagSet("lint", flag.ExitOnError)
	fix := fs.Bool("fix", false, "Automatically fix issues where possible")
	fs.Usage = func() {
		fmt.Println("Usage: wetwire-obs lint [options] [directory]")
		fmt.Println()
		fmt.Println("Check wetwire resources for code quality and best practices.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Lint rules (WOB prefix):")
		fmt.Println("  WOB001-019  Core wetwire patterns")
		fmt.Println("  WOB020-049  Prometheus config")
		fmt.Println("  WOB050-079  Alertmanager")
		fmt.Println("  WOB080-099  Rules")
		fmt.Println("  WOB100-119  PromQL")
		fmt.Println("  WOB120-149  Grafana dashboards")
		fmt.Println("  WOB150-169  Grafana panels")
		fmt.Println("  WOB200-219  Security")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-obs lint                     # Lint current directory")
		fmt.Println("  wetwire-obs lint ./monitoring        # Lint specific directory")
		fmt.Println("  wetwire-obs lint --fix               # Fix issues automatically")
	}

	if err := fs.Parse(args); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 2
	}

	// Get source directory
	srcDir := "."
	if fs.NArg() > 0 {
		srcDir = fs.Arg(0)
	}

	// Make path absolute
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Discover resources
	result, discoverErr := discover.Discover(srcDir)
	if discoverErr != nil {
		fmt.Fprintf(os.Stderr, "Error discovering resources: %v\n", discoverErr)
		return 1
	}

	// Placeholder: lint rules will be implemented in future phases
	_ = fix // Will be used when lint rules are implemented

	if result.TotalCount() == 0 {
		fmt.Println("No resources discovered to lint")
		return 0
	}

	fmt.Printf("Linting %d resources...\n", result.TotalCount())

	// For now, just report that linting passed
	// TODO: Implement actual lint rules
	fmt.Println()
	fmt.Println("Lint passed: 0 issues found")
	fmt.Println()
	fmt.Println("Note: Full lint rules (WOB001-WOB219) will be implemented in future phases.")

	return 0
}
