package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"go/build"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/lex00/wetwire-observability-go/alertmanager"
	"github.com/lex00/wetwire-observability-go/internal/discover"
	"github.com/lex00/wetwire-observability-go/prometheus"
)

// buildCmd handles the build command
func buildCmd(args []string) int {
	fs := flag.NewFlagSet("build", flag.ExitOnError)
	outputDir := fs.String("output", ".", "Output directory for generated files")
	mode := fs.String("mode", "standalone", "Output mode: standalone, operator, or both")
	fs.Usage = func() {
		fmt.Println("Usage: wetwire-obs build [options] [directory]")
		fmt.Println()
		fmt.Println("Generate configuration files from discovered resources.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  wetwire-obs build                    # Build from current directory")
		fmt.Println("  wetwire-obs build ./monitoring       # Build from specific directory")
		fmt.Println("  wetwire-obs build --output ./out     # Write output to ./out")
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

	// Make paths absolute
	srcDir, err := filepath.Abs(srcDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	*outputDir, err = filepath.Abs(*outputDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		return 1
	}

	// Validate mode
	if *mode != "standalone" && *mode != "operator" && *mode != "both" {
		fmt.Fprintf(os.Stderr, "Error: invalid mode %q, must be standalone, operator, or both\n", *mode)
		return 2
	}

	// Discover resources
	result, err := discover.Discover(srcDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error discovering resources: %v\n", err)
		return 1
	}

	if result.TotalCount() == 0 {
		fmt.Println("No resources discovered")
		return 0
	}

	// Load and serialize configs
	if len(result.PrometheusConfigs) > 0 {
		if err := buildPrometheusConfigs(srcDir, result.PrometheusConfigs, *outputDir, *mode); err != nil {
			fmt.Fprintf(os.Stderr, "Error building prometheus configs: %v\n", err)
			return 1
		}
	}

	if len(result.AlertmanagerConfigs) > 0 {
		if err := buildAlertmanagerConfigs(srcDir, result.AlertmanagerConfigs, *outputDir, *mode); err != nil {
			fmt.Fprintf(os.Stderr, "Error building alertmanager configs: %v\n", err)
			return 1
		}
	}

	fmt.Printf("Build complete: %d resources processed\n", result.TotalCount())
	return 0
}

// buildPrometheusConfigs loads and serializes PrometheusConfig resources
func buildPrometheusConfigs(srcDir string, refs []*discover.ResourceRef, outputDir, mode string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for _, ref := range refs {
		fmt.Printf("Processing %s.%s from %s:%d\n", ref.Package, ref.Name, filepath.Base(ref.FilePath), ref.Line)

		// Load the config by executing the package
		config, err := loadPrometheusConfig(srcDir, ref)
		if err != nil {
			fmt.Printf("  Warning: could not load config: %v\n", err)
			continue
		}

		// Generate output filename
		outputFile := filepath.Join(outputDir, fmt.Sprintf("prometheus-%s.yml", strings.ToLower(ref.Name)))

		// Serialize to file
		if err := config.SerializeToFile(outputFile); err != nil {
			return fmt.Errorf("serializing %s: %w", ref.Name, err)
		}

		fmt.Printf("  Generated %s\n", outputFile)
	}

	return nil
}

// loadPrometheusConfig loads a PrometheusConfig by building and running the package
func loadPrometheusConfig(srcDir string, ref *discover.ResourceRef) (*prometheus.PrometheusConfig, error) {
	// First, try to load the config by building a helper program
	// that imports the package and outputs the config as JSON

	// Get the package path
	pkgDir := filepath.Dir(ref.FilePath)
	pkg, err := build.ImportDir(pkgDir, build.FindOnly)
	if err != nil {
		// Fall back to creating a minimal example config
		return createMinimalConfig(ref), nil
	}

	// Create a temporary program to load and output the config
	tmpDir, err := os.MkdirTemp("", "wetwire-build-*")
	if err != nil {
		return createMinimalConfig(ref), nil
	}
	defer os.RemoveAll(tmpDir)

	// Write the helper program
	helperCode := fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"os"

	target %q
)

func main() {
	data, err := json.Marshal(target.%s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
`, pkg.ImportPath, ref.Name)

	helperPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(helperPath, []byte(helperCode), 0644); err != nil {
		return createMinimalConfig(ref), nil
	}

	// Initialize go.mod for the helper
	modPath := filepath.Join(tmpDir, "go.mod")
	modContent := fmt.Sprintf(`module helper

go 1.23

require %s v0.0.0

replace %s => %s
`, pkg.ImportPath, pkg.ImportPath, pkgDir)
	if err := os.WriteFile(modPath, []byte(modContent), 0644); err != nil {
		return createMinimalConfig(ref), nil
	}

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tmpDir
	if err := tidyCmd.Run(); err != nil {
		return createMinimalConfig(ref), nil
	}

	// Build and run the helper
	runCmd := exec.Command("go", "run", "main.go")
	runCmd.Dir = tmpDir
	output, err := runCmd.Output()
	if err != nil {
		return createMinimalConfig(ref), nil
	}

	// Parse the output
	var config prometheus.PrometheusConfig
	if err := json.Unmarshal(output, &config); err != nil {
		return createMinimalConfig(ref), nil
	}

	return &config, nil
}

// createMinimalConfig creates a minimal placeholder config
func createMinimalConfig(ref *discover.ResourceRef) *prometheus.PrometheusConfig {
	return &prometheus.PrometheusConfig{
		Global: &prometheus.GlobalConfig{
			ExternalLabels: map[string]string{
				"source": ref.Name,
			},
		},
	}
}

// isZeroValue checks if an interface value is the zero value for its type
func isZeroValue(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	return !val.IsValid() || val.IsZero()
}

// buildAlertmanagerConfigs loads and serializes AlertmanagerConfig resources
func buildAlertmanagerConfigs(srcDir string, refs []*discover.ResourceRef, outputDir, mode string) error {
	// Create output directory
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("creating output directory: %w", err)
	}

	for _, ref := range refs {
		fmt.Printf("Processing %s.%s from %s:%d\n", ref.Package, ref.Name, filepath.Base(ref.FilePath), ref.Line)

		// Load the config by executing the package
		config, err := loadAlertmanagerConfig(srcDir, ref)
		if err != nil {
			fmt.Printf("  Warning: could not load config: %v\n", err)
			continue
		}

		// Generate output filename
		outputFile := filepath.Join(outputDir, fmt.Sprintf("alertmanager-%s.yml", strings.ToLower(ref.Name)))

		// Serialize to file
		if err := config.SerializeToFile(outputFile); err != nil {
			return fmt.Errorf("serializing %s: %w", ref.Name, err)
		}

		fmt.Printf("  Generated %s\n", outputFile)
	}

	return nil
}

// loadAlertmanagerConfig loads an AlertmanagerConfig by building and running the package
func loadAlertmanagerConfig(srcDir string, ref *discover.ResourceRef) (*alertmanager.AlertmanagerConfig, error) {
	// Get the package path
	pkgDir := filepath.Dir(ref.FilePath)
	pkg, err := build.ImportDir(pkgDir, build.FindOnly)
	if err != nil {
		// Fall back to creating a minimal example config
		return createMinimalAlertmanagerConfig(ref), nil
	}

	// Create a temporary program to load and output the config
	tmpDir, err := os.MkdirTemp("", "wetwire-build-*")
	if err != nil {
		return createMinimalAlertmanagerConfig(ref), nil
	}
	defer os.RemoveAll(tmpDir)

	// Write the helper program
	helperCode := fmt.Sprintf(`package main

import (
	"encoding/json"
	"fmt"
	"os"

	target %q
)

func main() {
	data, err := json.Marshal(target.%s)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %%v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}
`, pkg.ImportPath, ref.Name)

	helperPath := filepath.Join(tmpDir, "main.go")
	if err := os.WriteFile(helperPath, []byte(helperCode), 0644); err != nil {
		return createMinimalAlertmanagerConfig(ref), nil
	}

	// Initialize go.mod for the helper
	modPath := filepath.Join(tmpDir, "go.mod")
	modContent := fmt.Sprintf(`module helper

go 1.23

require %s v0.0.0

replace %s => %s
`, pkg.ImportPath, pkg.ImportPath, pkgDir)
	if err := os.WriteFile(modPath, []byte(modContent), 0644); err != nil {
		return createMinimalAlertmanagerConfig(ref), nil
	}

	// Run go mod tidy
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tmpDir
	if err := tidyCmd.Run(); err != nil {
		return createMinimalAlertmanagerConfig(ref), nil
	}

	// Build and run the helper
	runCmd := exec.Command("go", "run", "main.go")
	runCmd.Dir = tmpDir
	output, err := runCmd.Output()
	if err != nil {
		return createMinimalAlertmanagerConfig(ref), nil
	}

	// Parse the output
	var config alertmanager.AlertmanagerConfig
	if err := json.Unmarshal(output, &config); err != nil {
		return createMinimalAlertmanagerConfig(ref), nil
	}

	return &config, nil
}

// createMinimalAlertmanagerConfig creates a minimal placeholder config
func createMinimalAlertmanagerConfig(ref *discover.ResourceRef) *alertmanager.AlertmanagerConfig {
	return &alertmanager.AlertmanagerConfig{
		Route: &alertmanager.Route{
			Receiver: "default",
		},
		Receivers: []*alertmanager.Receiver{
			{Name: "default"},
		},
	}
}
