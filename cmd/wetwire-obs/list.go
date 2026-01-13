package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"text/tabwriter"

	"github.com/lex00/wetwire-observability-go/internal/discover"
)

// listCmd handles the list command
func listCmd(args []string) int {
	fs := flag.NewFlagSet("list", flag.ExitOnError)
	format := fs.String("format", "table", "Output format: table or json")
	fs.Usage = func() {
		fmt.Println("Usage: wetwire-obs list [options] [directory]")
		fmt.Println()
		fmt.Println("List discovered wetwire resources.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-obs list                     # List from current directory")
		fmt.Println("  wetwire-obs list ./monitoring        # List from specific directory")
		fmt.Println("  wetwire-obs list --format json       # Output as JSON")
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

	// Validate format
	if *format != "table" && *format != "json" {
		fmt.Fprintf(os.Stderr, "Error: invalid format %q, must be table or json\n", *format)
		return 2
	}

	// Discover resources
	result, err := discover.Discover(srcDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering resources: %v\n", err)
		return 1
	}

	// Output results
	switch *format {
	case "json":
		return listJSON(result)
	default:
		return listTable(result, srcDir)
	}
}

// listTable outputs resources in table format
func listTable(result *discover.DiscoveryResult, baseDir string) int {
	if result.TotalCount() == 0 {
		fmt.Println("No resources discovered")
		return 0
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	fmt.Fprintln(w, "TYPE\tNAME\tPACKAGE\tFILE\tLINE")
	fmt.Fprintln(w, "----\t----\t-------\t----\t----")

	for _, ref := range result.All() {
		// Make file path relative for readability
		relPath, err := filepath.Rel(baseDir, ref.FilePath)
		if err != nil {
			relPath = ref.FilePath
		}

		fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%d\n",
			ref.Type,
			ref.Name,
			ref.Package,
			relPath,
			ref.Line,
		)
	}

	w.Flush()

	fmt.Println()
	fmt.Printf("Total: %d resources\n", result.TotalCount())

	// Print summary by type
	if len(result.PrometheusConfigs) > 0 {
		fmt.Printf("  PrometheusConfig: %d\n", len(result.PrometheusConfigs))
	}
	if len(result.ScrapeConfigs) > 0 {
		fmt.Printf("  ScrapeConfig: %d\n", len(result.ScrapeConfigs))
	}
	if len(result.GlobalConfigs) > 0 {
		fmt.Printf("  GlobalConfig: %d\n", len(result.GlobalConfigs))
	}
	if len(result.StaticConfigs) > 0 {
		fmt.Printf("  StaticConfig: %d\n", len(result.StaticConfigs))
	}

	// Print warnings
	if len(result.Errors) > 0 {
		fmt.Println()
		fmt.Printf("Warnings: %d files could not be parsed\n", len(result.Errors))
		for _, e := range result.Errors {
			fmt.Printf("  - %s\n", e)
		}
	}

	return 0
}

// listJSON outputs resources in JSON format
func listJSON(result *discover.DiscoveryResult) int {
	output := struct {
		Resources []*discover.ResourceRef `json:"resources"`
		Summary   struct {
			Total            int `json:"total"`
			PrometheusConfig int `json:"prometheus_config"`
			ScrapeConfig     int `json:"scrape_config"`
			GlobalConfig     int `json:"global_config"`
			StaticConfig     int `json:"static_config"`
		} `json:"summary"`
		Errors []string `json:"errors,omitempty"`
	}{
		Resources: result.All(),
		Errors:    result.Errors,
	}

	output.Summary.Total = result.TotalCount()
	output.Summary.PrometheusConfig = len(result.PrometheusConfigs)
	output.Summary.ScrapeConfig = len(result.ScrapeConfigs)
	output.Summary.GlobalConfig = len(result.GlobalConfigs)
	output.Summary.StaticConfig = len(result.StaticConfigs)

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	fmt.Println(string(data))
	return 0
}
