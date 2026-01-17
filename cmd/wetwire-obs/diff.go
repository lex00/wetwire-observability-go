// Command diff compares two observability configurations.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/lex00/wetwire-observability-go/internal/discover"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

func newDiffCmd() *cobra.Command {
	var (
		format   string
		semantic bool
		color    bool
	)

	cmd := &cobra.Command{
		Use:   "diff <path1> <path2>",
		Short: "Compare two observability configuration directories or files",
		Long: `Diff compares two observability configuration directories or files.

The comparison can be text-based or semantic (ignoring ordering).

Examples:
  wetwire-obs diff old.yaml new.yaml
  wetwire-obs diff ./staging ./production
  wetwire-obs diff --semantic old.yaml new.yaml
  wetwire-obs diff --format json old.yaml new.yaml`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runDiff(args[0], args[1], format, semantic, color)
		},
	}

	cmd.Flags().StringVarP(&format, "format", "f", "text", "Output format: text or json")
	cmd.Flags().BoolVar(&semantic, "semantic", false, "Use semantic comparison (ignores ordering)")
	cmd.Flags().BoolVar(&color, "color", false, "Enable colorized output")

	return cmd
}

// DiffResult holds the comparison result
type DiffResult struct {
	HasDifferences bool        `json:"has_differences"`
	Added          []DiffEntry `json:"added,omitempty"`
	Removed        []DiffEntry `json:"removed,omitempty"`
	Modified       []DiffEntry `json:"modified,omitempty"`
	Summary        DiffSummary `json:"summary"`
}

// DiffEntry represents a single difference
type DiffEntry struct {
	Name    string `json:"name"`
	Type    string `json:"type"`
	Path    string `json:"path,omitempty"`
	Changes string `json:"changes,omitempty"`
}

// DiffSummary summarizes the differences
type DiffSummary struct {
	Added    int `json:"added"`
	Removed  int `json:"removed"`
	Modified int `json:"modified"`
	Total    int `json:"total"`
}

func runDiff(path1, path2, format string, semantic, useColor bool) error {
	result, err := comparePaths(path1, path2, semantic)
	if err != nil {
		return err
	}

	switch format {
	case "json":
		data, _ := json.MarshalIndent(result, "", "  ")
		fmt.Println(string(data))
	default:
		printDiffText(result, useColor)
	}

	if result.HasDifferences {
		os.Exit(1)
	}
	return nil
}

func comparePaths(path1, path2 string, semantic bool) (*DiffResult, error) {
	// Check if paths are files or directories
	info1, err := os.Stat(path1)
	if err != nil {
		return nil, fmt.Errorf("path1: %w", err)
	}
	info2, err := os.Stat(path2)
	if err != nil {
		return nil, fmt.Errorf("path2: %w", err)
	}

	result := &DiffResult{}

	if info1.IsDir() && info2.IsDir() {
		// Compare directories using discover
		res1, err := discover.Discover(path1)
		if err != nil {
			return nil, fmt.Errorf("discover path1: %w", err)
		}
		res2, err := discover.Discover(path2)
		if err != nil {
			return nil, fmt.Errorf("discover path2: %w", err)
		}

		// Build maps of resources by name
		map1 := buildResourceMap(res1)
		map2 := buildResourceMap(res2)

		// Find added, removed, modified
		for name, entry := range map2 {
			if _, exists := map1[name]; !exists {
				result.Added = append(result.Added, entry)
			}
		}
		for name, entry := range map1 {
			if _, exists := map2[name]; !exists {
				result.Removed = append(result.Removed, entry)
			}
		}

	} else if !info1.IsDir() && !info2.IsDir() {
		// Compare files
		data1, err := os.ReadFile(path1)
		if err != nil {
			return nil, fmt.Errorf("read path1: %w", err)
		}
		data2, err := os.ReadFile(path2)
		if err != nil {
			return nil, fmt.Errorf("read path2: %w", err)
		}

		if semantic {
			// Parse as YAML and compare
			var obj1, obj2 interface{}
			if err := yaml.Unmarshal(data1, &obj1); err != nil {
				return nil, fmt.Errorf("parse path1: %w", err)
			}
			if err := yaml.Unmarshal(data2, &obj2); err != nil {
				return nil, fmt.Errorf("parse path2: %w", err)
			}

			// Compare JSON representations
			json1, _ := json.Marshal(obj1)
			json2, _ := json.Marshal(obj2)

			if string(json1) != string(json2) {
				result.Modified = append(result.Modified, DiffEntry{
					Name:    path1,
					Type:    "file",
					Changes: "content differs",
				})
			}
		} else {
			// Simple text comparison
			if string(data1) != string(data2) {
				result.Modified = append(result.Modified, DiffEntry{
					Name:    path1,
					Type:    "file",
					Changes: "content differs",
				})
			}
		}
	} else {
		return nil, fmt.Errorf("cannot compare directory with file")
	}

	// Build summary
	result.Summary.Added = len(result.Added)
	result.Summary.Removed = len(result.Removed)
	result.Summary.Modified = len(result.Modified)
	result.Summary.Total = result.Summary.Added + result.Summary.Removed + result.Summary.Modified
	result.HasDifferences = result.Summary.Total > 0

	return result, nil
}

func buildResourceMap(res *discover.DiscoveryResult) map[string]DiffEntry {
	m := make(map[string]DiffEntry)

	for _, ref := range res.PrometheusConfigs {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "prometheus_config", Path: ref.FilePath}
	}
	for _, ref := range res.ScrapeConfigs {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "scrape_config", Path: ref.FilePath}
	}
	for _, ref := range res.GlobalConfigs {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "global_config", Path: ref.FilePath}
	}
	for _, ref := range res.AlertmanagerConfigs {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "alertmanager_config", Path: ref.FilePath}
	}
	for _, ref := range res.RulesFiles {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "rules_file", Path: ref.FilePath}
	}
	for _, ref := range res.RuleGroups {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "rule_group", Path: ref.FilePath}
	}
	for _, ref := range res.AlertingRules {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "alerting_rule", Path: ref.FilePath}
	}
	for _, ref := range res.RecordingRules {
		m[ref.Name] = DiffEntry{Name: ref.Name, Type: "recording_rule", Path: ref.FilePath}
	}

	return m
}

func printDiffText(result *DiffResult, useColor bool) {
	if !result.HasDifferences {
		fmt.Println("No differences found")
		return
	}

	colorRed := ""
	colorGreen := ""
	colorYellow := ""
	colorReset := ""

	if useColor {
		colorRed = "\033[31m"
		colorGreen = "\033[32m"
		colorYellow = "\033[33m"
		colorReset = "\033[0m"
	}

	if len(result.Added) > 0 {
		fmt.Println("=== Added ===")
		for _, e := range result.Added {
			fmt.Printf("%s+ %s (%s)%s\n", colorGreen, e.Name, e.Type, colorReset)
		}
		fmt.Println()
	}

	if len(result.Removed) > 0 {
		fmt.Println("=== Removed ===")
		for _, e := range result.Removed {
			fmt.Printf("%s- %s (%s)%s\n", colorRed, e.Name, e.Type, colorReset)
		}
		fmt.Println()
	}

	if len(result.Modified) > 0 {
		fmt.Println("=== Modified ===")
		for _, e := range result.Modified {
			fmt.Printf("%s~ %s (%s)%s\n", colorYellow, e.Name, e.Type, colorReset)
			if e.Changes != "" {
				fmt.Printf("    %s\n", e.Changes)
			}
		}
		fmt.Println()
	}

	fmt.Println(strings.Repeat("-", 40))
	fmt.Printf("Summary: %d added, %d removed, %d modified\n",
		result.Summary.Added, result.Summary.Removed, result.Summary.Modified)
}
