package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/lex00/wetwire-observability-go/internal/importer"
)

func importCmd(args []string) int {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	output := fs.String("output", "", "Output file path (default: stdout)")
	pkg := fs.String("package", "monitoring", "Go package name for generated code")
	help := fs.Bool("help", false, "Show help")

	fs.Usage = func() {
		fmt.Println("wetwire-obs import - Convert prometheus.yml to Go code")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  wetwire-obs import [options] <prometheus.yml>")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Example:")
		fmt.Println("  wetwire-obs import prometheus.yml --output monitoring/config.go")
		fmt.Println("  wetwire-obs import /etc/prometheus/prometheus.yml --package infra")
	}

	if err := fs.Parse(args); err != nil {
		return 1
	}

	if *help {
		fs.Usage()
		return 0
	}

	if fs.NArg() < 1 {
		fmt.Fprintln(os.Stderr, "Error: no input file specified")
		fs.Usage()
		return 1
	}

	inputPath := fs.Arg(0)

	// Parse the prometheus.yml
	config, err := importer.ParsePrometheusConfig(inputPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing %s: %v\n", inputPath, err)
		return 1
	}

	// Validate and show warnings
	warnings := importer.ValidatePrometheusConfig(config)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", w)
	}

	// Generate Go code
	code, err := importer.GenerateGoCode(config, *pkg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error generating code: %v\n", err)
		return 1
	}

	// Output
	if *output == "" {
		// Write to stdout
		fmt.Print(string(code))
	} else {
		// Create output directory if needed
		dir := filepath.Dir(*output)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
				return 1
			}
		}

		// Write to file
		if err := os.WriteFile(*output, code, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *output, err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "Generated %s\n", *output)
	}

	return 0
}
