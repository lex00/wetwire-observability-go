// Package differ provides semantic comparison of observability configurations.
package differ

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
	"sort"

	coredomain "github.com/lex00/wetwire-core-go/domain"
	"github.com/lex00/wetwire-observability-go/internal/discover"
	"gopkg.in/yaml.v3"
)

// ObservabilityDiffer implements coredomain.Differ for observability configs.
type ObservabilityDiffer struct{}

// Compile-time check that ObservabilityDiffer implements Differ.
var _ coredomain.Differ = (*ObservabilityDiffer)(nil)

// New creates a new observability differ.
func New() *ObservabilityDiffer {
	return &ObservabilityDiffer{}
}

// Diff compares two observability configurations and returns differences.
func (d *ObservabilityDiffer) Diff(ctx *coredomain.Context, file1, file2 string, opts coredomain.DiffOpts) (*coredomain.DiffResult, error) {
	info1, err := os.Stat(file1)
	if err != nil {
		return nil, fmt.Errorf("path1: %w", err)
	}
	info2, err := os.Stat(file2)
	if err != nil {
		return nil, fmt.Errorf("path2: %w", err)
	}

	result := &coredomain.DiffResult{}

	if info1.IsDir() && info2.IsDir() {
		// Compare directories using discover
		res1, err := discover.Discover(file1)
		if err != nil {
			return nil, fmt.Errorf("discover path1: %w", err)
		}
		res2, err := discover.Discover(file2)
		if err != nil {
			return nil, fmt.Errorf("discover path2: %w", err)
		}

		// Build maps of resources by name
		map1 := buildResourceMap(res1)
		map2 := buildResourceMap(res2)

		// Find added resources (in dir2 but not in dir1)
		for name, entry := range map2 {
			if _, exists := map1[name]; !exists {
				result.Entries = append(result.Entries, coredomain.DiffEntry{
					Resource: name,
					Type:     entry.Type,
					Action:   "added",
				})
			}
		}

		// Find removed resources (in dir1 but not in dir2)
		for name, entry := range map1 {
			if _, exists := map2[name]; !exists {
				result.Entries = append(result.Entries, coredomain.DiffEntry{
					Resource: name,
					Type:     entry.Type,
					Action:   "removed",
				})
			}
		}

		// Find modified resources
		for name, entry1 := range map1 {
			if entry2, exists := map2[name]; exists {
				if entry1.Type != entry2.Type {
					result.Entries = append(result.Entries, coredomain.DiffEntry{
						Resource: name,
						Type:     entry1.Type,
						Action:   "modified",
						Changes:  []string{fmt.Sprintf("type changed: %s → %s", entry1.Type, entry2.Type)},
					})
				}
			}
		}

	} else if !info1.IsDir() && !info2.IsDir() {
		// Compare files
		changes, err := compareFiles(file1, file2, opts)
		if err != nil {
			return nil, err
		}
		result.Entries = changes
	} else {
		return nil, fmt.Errorf("cannot compare directory with file")
	}

	// Sort entries for consistent output
	sort.Slice(result.Entries, func(i, j int) bool {
		if result.Entries[i].Action != result.Entries[j].Action {
			order := map[string]int{"added": 0, "modified": 1, "removed": 2}
			return order[result.Entries[i].Action] < order[result.Entries[j].Action]
		}
		return result.Entries[i].Resource < result.Entries[j].Resource
	})

	// Calculate summary
	for _, e := range result.Entries {
		switch e.Action {
		case "added":
			result.Summary.Added++
		case "removed":
			result.Summary.Removed++
		case "modified":
			result.Summary.Modified++
		}
	}
	result.Summary.Total = result.Summary.Added + result.Summary.Removed + result.Summary.Modified

	return result, nil
}

// resourceEntry holds info about a discovered resource
type resourceEntry struct {
	Name string
	Type string
	Path string
}

// buildResourceMap builds a map of resources from discovery results.
func buildResourceMap(res *discover.DiscoveryResult) map[string]resourceEntry {
	m := make(map[string]resourceEntry)

	for _, ref := range res.PrometheusConfigs {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "prometheus_config", Path: ref.FilePath}
	}
	for _, ref := range res.ScrapeConfigs {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "scrape_config", Path: ref.FilePath}
	}
	for _, ref := range res.GlobalConfigs {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "global_config", Path: ref.FilePath}
	}
	for _, ref := range res.StaticConfigs {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "static_config", Path: ref.FilePath}
	}
	for _, ref := range res.AlertmanagerConfigs {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "alertmanager_config", Path: ref.FilePath}
	}
	for _, ref := range res.RulesFiles {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "rules_file", Path: ref.FilePath}
	}
	for _, ref := range res.RuleGroups {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "rule_group", Path: ref.FilePath}
	}
	for _, ref := range res.AlertingRules {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "alerting_rule", Path: ref.FilePath}
	}
	for _, ref := range res.RecordingRules {
		m[ref.Name] = resourceEntry{Name: ref.Name, Type: "recording_rule", Path: ref.FilePath}
	}

	return m
}

// compareFiles compares two config files and returns differences.
func compareFiles(file1, file2 string, opts coredomain.DiffOpts) ([]coredomain.DiffEntry, error) {
	data1, err := os.ReadFile(file1)
	if err != nil {
		return nil, fmt.Errorf("read file1: %w", err)
	}
	data2, err := os.ReadFile(file2)
	if err != nil {
		return nil, fmt.Errorf("read file2: %w", err)
	}

	// Parse YAML (or JSON)
	var obj1, obj2 interface{}
	if err := yaml.Unmarshal(data1, &obj1); err != nil {
		// Try JSON
		if err := json.Unmarshal(data1, &obj1); err != nil {
			return nil, fmt.Errorf("parse file1: %w", err)
		}
	}
	if err := yaml.Unmarshal(data2, &obj2); err != nil {
		// Try JSON
		if err := json.Unmarshal(data2, &obj2); err != nil {
			return nil, fmt.Errorf("parse file2: %w", err)
		}
	}

	// Compare using deep equal with optional order ignoring
	if deepEqual(obj1, obj2, opts) {
		return nil, nil
	}

	// Try to identify specific changes
	changes := findChanges("", obj1, obj2, opts)

	return []coredomain.DiffEntry{
		{
			Resource: file1,
			Type:     "file",
			Action:   "modified",
			Changes:  changes,
		},
	}, nil
}

// findChanges recursively finds changes between two values.
func findChanges(path string, v1, v2 interface{}, opts coredomain.DiffOpts) []string {
	var changes []string

	if v1 == nil && v2 == nil {
		return nil
	}
	if v1 == nil {
		if path != "" {
			return []string{fmt.Sprintf("%s: added", path)}
		}
		return nil
	}
	if v2 == nil {
		if path != "" {
			return []string{fmt.Sprintf("%s: removed", path)}
		}
		return nil
	}

	// Type mismatch
	if reflect.TypeOf(v1) != reflect.TypeOf(v2) {
		if path != "" {
			return []string{fmt.Sprintf("%s: type changed", path)}
		}
		return nil
	}

	switch val1 := v1.(type) {
	case map[string]interface{}:
		val2 := v2.(map[string]interface{})
		// Check for removed keys
		for k := range val1 {
			if _, exists := val2[k]; !exists {
				subPath := joinPath(path, k)
				changes = append(changes, fmt.Sprintf("%s: removed", subPath))
			}
		}
		// Check for added/modified keys
		for k, v2Val := range val2 {
			subPath := joinPath(path, k)
			if v1Val, exists := val1[k]; !exists {
				changes = append(changes, fmt.Sprintf("%s: added", subPath))
			} else {
				changes = append(changes, findChanges(subPath, v1Val, v2Val, opts)...)
			}
		}

	case []interface{}:
		val2 := v2.([]interface{})
		if opts.IgnoreOrder {
			// Compare as sets
			if !equalSets(val1, val2) {
				changes = append(changes, fmt.Sprintf("%s: changed", path))
			}
		} else {
			// Compare element by element
			if len(val1) != len(val2) {
				changes = append(changes, fmt.Sprintf("%s: length changed (%d → %d)", path, len(val1), len(val2)))
			} else {
				for i := range val1 {
					subPath := fmt.Sprintf("%s[%d]", path, i)
					changes = append(changes, findChanges(subPath, val1[i], val2[i], opts)...)
				}
			}
		}

	default:
		// Scalar comparison
		if !reflect.DeepEqual(v1, v2) {
			if path != "" {
				changes = append(changes, fmt.Sprintf("%s: %v → %v", path, v1, v2))
			}
		}
	}

	return changes
}

// joinPath joins path components.
func joinPath(base, key string) string {
	if base == "" {
		return key
	}
	return base + "." + key
}

// deepEqual compares two values with optional order ignoring.
func deepEqual(a, b interface{}, opts coredomain.DiffOpts) bool {
	if opts.IgnoreOrder {
		a = normalizeValue(a)
		b = normalizeValue(b)
	}
	return reflect.DeepEqual(a, b)
}

// normalizeValue normalizes a value for comparison.
func normalizeValue(v interface{}) interface{} {
	switch val := v.(type) {
	case []interface{}:
		result := make([]interface{}, len(val))
		copy(result, val)
		return result
	case map[string]interface{}:
		result := make(map[string]interface{})
		for k, v := range val {
			result[k] = normalizeValue(v)
		}
		return result
	default:
		return v
	}
}

// equalSets checks if two slices contain the same elements.
func equalSets(a, b []interface{}) bool {
	if len(a) != len(b) {
		return false
	}

	strA := make([]string, len(a))
	strB := make([]string, len(b))
	for i, v := range a {
		strA[i] = fmt.Sprintf("%v", v)
	}
	for i, v := range b {
		strB[i] = fmt.Sprintf("%v", v)
	}
	sort.Strings(strA)
	sort.Strings(strB)

	for i := range strA {
		if strA[i] != strB[i] {
			return false
		}
	}
	return true
}
