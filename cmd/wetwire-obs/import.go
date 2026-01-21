package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/lex00/wetwire-observability-go/internal/importer"
)

func importCmd(args []string) int {
	fs := flag.NewFlagSet("import", flag.ExitOnError)
	output := fs.String("output", "", "Output file path (default: stdout)")
	pkg := fs.String("package", "monitoring", "Go package name for generated code")
	configType := fs.String("type", "", "Config type: prometheus, alertmanager (auto-detected if not specified)")
	help := fs.Bool("help", false, "Show help")

	fs.Usage = func() {
		fmt.Println("wetwire-obs import - Convert config files to Go code")
		fmt.Println()
		fmt.Println("Usage:")
		fmt.Println("  wetwire-obs import [options] <config-file>")
		fmt.Println()
		fmt.Println("Supported formats:")
		fmt.Println("  - prometheus.yml    (Prometheus configuration)")
		fmt.Println("  - alertmanager.yml  (Alertmanager configuration)")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-obs import prometheus.yml --output monitoring/prometheus.go")
		fmt.Println("  wetwire-obs import alertmanager.yml --output monitoring/alertmanager.go")
		fmt.Println("  wetwire-obs import config.yml --type=prometheus --package infra")
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

	// Auto-detect config type from filename if not specified
	detectedType := *configType
	if detectedType == "" {
		detectedType = detectConfigType(inputPath)
	}

	var code []byte
	var err error

	switch detectedType {
	case "prometheus":
		code, err = importPrometheus(inputPath, *pkg)
	case "alertmanager":
		code, err = importAlertmanager(inputPath, *pkg)
	default:
		fmt.Fprintf(os.Stderr, "Error: unknown config type %q. Use --type to specify.\n", detectedType)
		return 1
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Output
	if *output == "" {
		fmt.Print(string(code))
	} else {
		dir := filepath.Dir(*output)
		if dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				fmt.Fprintf(os.Stderr, "Error creating directory %s: %v\n", dir, err)
				return 1
			}
		}

		if err := os.WriteFile(*output, code, 0644); err != nil {
			fmt.Fprintf(os.Stderr, "Error writing %s: %v\n", *output, err)
			return 1
		}
		fmt.Fprintf(os.Stderr, "Generated %s\n", *output)
	}

	return 0
}

// detectConfigType guesses the config type from the filename.
func detectConfigType(path string) string {
	base := strings.ToLower(filepath.Base(path))

	if strings.Contains(base, "alertmanager") {
		return "alertmanager"
	}
	if strings.Contains(base, "prometheus") {
		return "prometheus"
	}

	// Default to prometheus for .yml files
	if strings.HasSuffix(base, ".yml") || strings.HasSuffix(base, ".yaml") {
		return "prometheus"
	}

	return ""
}

func importPrometheus(inputPath, pkg string) ([]byte, error) {
	config, err := importer.ParsePrometheusConfig(inputPath)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", inputPath, err)
	}

	warnings := importer.ValidatePrometheusConfig(config)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", w)
	}

	return importer.GenerateGoCode(config, pkg)
}

func importAlertmanager(inputPath, pkg string) ([]byte, error) {
	config, err := importer.ParseAlertmanagerConfig(inputPath)
	if err != nil {
		return nil, fmt.Errorf("parsing %s: %w", inputPath, err)
	}

	warnings := importer.ValidateAlertmanagerConfig(config)
	for _, w := range warnings {
		fmt.Fprintf(os.Stderr, "Warning: %s\n", w)
	}

	return importer.GenerateAlertmanagerGoCode(config, pkg)
}
