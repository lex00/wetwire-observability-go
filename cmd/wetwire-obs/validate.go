package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func validateCmd(args []string) int {
	fs := flag.NewFlagSet("validate", flag.ContinueOnError)
	outputDir := fs.String("output", "output", "Directory containing generated files")
	skipKubeconform := fs.Bool("skip-kubeconform", false, "Skip kubeconform validation")
	skipPromtool := fs.Bool("skip-promtool", false, "Skip promtool validation")
	skipAmtool := fs.Bool("skip-amtool", false, "Skip amtool validation")
	verbose := fs.Bool("v", false, "Verbose output")

	fs.Usage = func() {
		fmt.Println("Usage: wetwire-obs validate [options]")
		fmt.Println()
		fmt.Println("Validates generated configuration files using external tools.")
		fmt.Println()
		fmt.Println("Options:")
		fs.PrintDefaults()
		fmt.Println()
		fmt.Println("Validators:")
		fmt.Println("  kubeconform  Validates Kubernetes manifests (ServiceMonitor, etc.)")
		fmt.Println("  promtool     Validates prometheus.yml and rule files")
		fmt.Println("  amtool       Validates alertmanager.yml")
		fmt.Println()
		fmt.Println("Validators are skipped gracefully if not installed.")
	}

	if err := fs.Parse(args); err != nil {
		if err == flag.ErrHelp {
			return 0
		}
		return 1
	}

	results := &validationResults{}

	// Validate Kubernetes manifests with kubeconform
	if !*skipKubeconform {
		results.kubeconform = validateKubeconform(*outputDir, *verbose)
	}

	// Validate Prometheus config with promtool
	if !*skipPromtool {
		results.promtool = validatePromtool(*outputDir, *verbose)
	}

	// Validate Alertmanager config with amtool
	if !*skipAmtool {
		results.amtool = validateAmtool(*outputDir, *verbose)
	}

	// Print summary
	printValidationSummary(results)

	if results.hasErrors() {
		return 1
	}
	return 0
}

type validationResults struct {
	kubeconform *validationResult
	promtool    *validationResult
	amtool      *validationResult
}

type validationResult struct {
	tool     string
	skipped  bool
	skipMsg  string
	success  bool
	messages []string
}

func (r *validationResults) hasErrors() bool {
	for _, result := range []*validationResult{r.kubeconform, r.promtool, r.amtool} {
		if result != nil && !result.skipped && !result.success {
			return true
		}
	}
	return false
}

func validateKubeconform(outputDir string, verbose bool) *validationResult {
	result := &validationResult{tool: "kubeconform"}

	// Check if kubeconform is installed
	if _, err := exec.LookPath("kubeconform"); err != nil {
		result.skipped = true
		result.skipMsg = "kubeconform not found in PATH (install: go install github.com/yannh/kubeconform/cmd/kubeconform@latest)"
		return result
	}

	// Find Kubernetes manifests
	var manifests []string
	operatorDir := filepath.Join(outputDir, "operator")
	if _, err := os.Stat(operatorDir); err == nil {
		files, err := filepath.Glob(filepath.Join(operatorDir, "*.yaml"))
		if err == nil {
			manifests = append(manifests, files...)
		}
	}

	if len(manifests) == 0 {
		result.skipped = true
		result.skipMsg = "no Kubernetes manifests found in " + operatorDir
		return result
	}

	// Run kubeconform with prometheus-operator schemas
	args := []string{
		"-strict",
		"-summary",
		"-schema-location", "default",
		"-schema-location", "https://raw.githubusercontent.com/datreeio/CRDs-catalog/main/{{.Group}}/{{.ResourceKind}}_{{.ResourceAPIVersion}}.json",
	}
	args = append(args, manifests...)

	cmd := exec.Command("kubeconform", args...)
	output, err := cmd.CombinedOutput()

	if verbose && len(output) > 0 {
		result.messages = append(result.messages, string(output))
	}

	if err != nil {
		result.success = false
		result.messages = append(result.messages, fmt.Sprintf("validation failed: %v", err))
		if len(output) > 0 && !verbose {
			result.messages = append(result.messages, string(output))
		}
	} else {
		result.success = true
		result.messages = append(result.messages, fmt.Sprintf("validated %d manifests", len(manifests)))
	}

	return result
}

func validatePromtool(outputDir string, verbose bool) *validationResult {
	result := &validationResult{tool: "promtool"}

	// Check if promtool is installed
	if _, err := exec.LookPath("promtool"); err != nil {
		result.skipped = true
		result.skipMsg = "promtool not found in PATH (install Prometheus)"
		return result
	}

	var validated int

	// Validate prometheus.yml
	prometheusYml := filepath.Join(outputDir, "prometheus.yml")
	if _, err := os.Stat(prometheusYml); err == nil {
		cmd := exec.Command("promtool", "check", "config", prometheusYml)
		output, err := cmd.CombinedOutput()
		if err != nil {
			result.success = false
			result.messages = append(result.messages, fmt.Sprintf("prometheus.yml: %v", err))
			if len(output) > 0 {
				result.messages = append(result.messages, string(output))
			}
			return result
		}
		validated++
		if verbose {
			result.messages = append(result.messages, "prometheus.yml: valid")
		}
	}

	// Validate rule files
	rulesDir := filepath.Join(outputDir, "rules")
	if _, err := os.Stat(rulesDir); err == nil {
		files, err := filepath.Glob(filepath.Join(rulesDir, "*.yml"))
		if err == nil {
			for _, file := range files {
				cmd := exec.Command("promtool", "check", "rules", file)
				output, err := cmd.CombinedOutput()
				if err != nil {
					result.success = false
					result.messages = append(result.messages, fmt.Sprintf("%s: %v", filepath.Base(file), err))
					if len(output) > 0 {
						result.messages = append(result.messages, string(output))
					}
				} else {
					validated++
					if verbose {
						result.messages = append(result.messages, fmt.Sprintf("%s: valid", filepath.Base(file)))
					}
				}
			}
		}
	}

	if validated == 0 {
		result.skipped = true
		result.skipMsg = "no Prometheus config or rules files found"
		return result
	}

	if !result.success {
		return result
	}

	result.success = true
	result.messages = append(result.messages, fmt.Sprintf("validated %d files", validated))
	return result
}

func validateAmtool(outputDir string, verbose bool) *validationResult {
	result := &validationResult{tool: "amtool"}

	// Check if amtool is installed
	if _, err := exec.LookPath("amtool"); err != nil {
		result.skipped = true
		result.skipMsg = "amtool not found in PATH (install Alertmanager)"
		return result
	}

	// Validate alertmanager.yml
	alertmanagerYml := filepath.Join(outputDir, "alertmanager.yml")
	if _, err := os.Stat(alertmanagerYml); err != nil {
		result.skipped = true
		result.skipMsg = "alertmanager.yml not found"
		return result
	}

	cmd := exec.Command("amtool", "check-config", alertmanagerYml)
	output, err := cmd.CombinedOutput()

	if verbose && len(output) > 0 {
		result.messages = append(result.messages, string(output))
	}

	if err != nil {
		result.success = false
		result.messages = append(result.messages, fmt.Sprintf("validation failed: %v", err))
		if len(output) > 0 && !verbose {
			result.messages = append(result.messages, string(output))
		}
	} else {
		result.success = true
		result.messages = append(result.messages, "alertmanager.yml: valid")
	}

	return result
}

func printValidationSummary(results *validationResults) {
	fmt.Println()
	fmt.Println("Validation Summary")
	fmt.Println(strings.Repeat("-", 50))

	for _, result := range []*validationResult{results.kubeconform, results.promtool, results.amtool} {
		if result == nil {
			continue
		}

		status := ""
		if result.skipped {
			status = "SKIPPED"
		} else if result.success {
			status = "PASS"
		} else {
			status = "FAIL"
		}

		fmt.Printf("%-15s %s\n", result.tool+":", status)

		if result.skipped && result.skipMsg != "" {
			fmt.Printf("  %s\n", result.skipMsg)
		}

		for _, msg := range result.messages {
			for _, line := range strings.Split(msg, "\n") {
				if line != "" {
					fmt.Printf("  %s\n", line)
				}
			}
		}
	}

	fmt.Println()
}
