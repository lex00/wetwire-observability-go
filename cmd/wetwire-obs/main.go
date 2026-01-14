// Package main provides the wetwire-obs CLI for observability configuration synthesis.
package main

import (
	"fmt"
	"os"
)

// Version is set by the build process
var Version = "dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) < 1 {
		printUsage()
		return 0
	}

	cmd := args[0]
	cmdArgs := args[1:]

	switch cmd {
	case "build":
		return buildCmd(cmdArgs)
	case "lint":
		return lintCmd(cmdArgs)
	case "list":
		return listCmd(cmdArgs)
	case "import":
		return importCmd(cmdArgs)
	case "validate":
		return validateCmd(cmdArgs)
	case "design":
		return designCmd(cmdArgs)
	case "test":
		return testCmd(cmdArgs)
	case "version":
		fmt.Printf("wetwire-obs %s\n", Version)
		return 0
	case "help", "-h", "--help":
		printUsage()
		return 0
	default:
		fmt.Fprintf(os.Stderr, "Unknown command: %s\n", cmd)
		fmt.Fprintln(os.Stderr, "Run 'wetwire-obs help' for usage.")
		return 1
	}
}

func printUsage() {
	fmt.Println("wetwire-obs - Prometheus, Alertmanager, and Grafana configuration synthesis")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  wetwire-obs <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  build     Generate configuration files (prometheus.yml, etc.)")
	fmt.Println("  lint      Check code quality and best practices")
	fmt.Println("  list      List discovered resources")
	fmt.Println("  import    Convert existing configs (prometheus.yml) to Go code")
	fmt.Println("  validate  Run external validators (promtool, amtool, kubeconform)")
	fmt.Println("  design    AI-assisted observability configuration generation")
	fmt.Println("  test      Evaluate configurations against personas")
	fmt.Println("  version   Show version information")
	fmt.Println()
	fmt.Println("Run 'wetwire-obs <command> --help' for command-specific help.")
}
