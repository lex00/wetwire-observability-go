package prometheus

// RelabelAction represents a relabeling action type.
type RelabelAction string

// Relabel actions supported by Prometheus.
const (
	// RelabelReplace replaces the value of target_label with the result
	// of matching against regex and using replacement string.
	RelabelReplace RelabelAction = "replace"

	// RelabelKeep drops targets where regex does not match.
	RelabelKeep RelabelAction = "keep"

	// RelabelDrop drops targets where regex matches.
	RelabelDrop RelabelAction = "drop"

	// RelabelHashMod sets target_label to the modulus of a hash of source labels.
	RelabelHashMod RelabelAction = "hashmod"

	// RelabelLabelMap copies labels matching regex to new labels.
	RelabelLabelMap RelabelAction = "labelmap"

	// RelabelLabelDrop drops labels matching regex.
	RelabelLabelDrop RelabelAction = "labeldrop"

	// RelabelLabelKeep keeps only labels matching regex.
	RelabelLabelKeep RelabelAction = "labelkeep"

	// RelabelLowercase converts label values to lowercase.
	RelabelLowercase RelabelAction = "lowercase"

	// RelabelUppercase converts label values to uppercase.
	RelabelUppercase RelabelAction = "uppercase"

	// RelabelKeepEqual keeps targets where source and target labels are equal.
	RelabelKeepEqual RelabelAction = "keepequal"

	// RelabelDropEqual drops targets where source and target labels are equal.
	RelabelDropEqual RelabelAction = "dropequal"
)

// NewRelabelConfig creates a new RelabelConfig with the given action.
func NewRelabelConfig(action RelabelAction) *RelabelConfig {
	return &RelabelConfig{
		Action: string(action),
	}
}

// WithSourceLabels sets the source labels for the relabel config.
func (r *RelabelConfig) WithSourceLabels(labels ...string) *RelabelConfig {
	r.SourceLabels = labels
	return r
}

// WithSeparator sets the separator for concatenating source labels.
func (r *RelabelConfig) WithSeparator(sep string) *RelabelConfig {
	r.Separator = sep
	return r
}

// WithRegex sets the regex pattern.
func (r *RelabelConfig) WithRegex(regex string) *RelabelConfig {
	r.Regex = regex
	return r
}

// WithTargetLabel sets the target label.
func (r *RelabelConfig) WithTargetLabel(label string) *RelabelConfig {
	r.TargetLabel = label
	return r
}

// WithReplacement sets the replacement string.
func (r *RelabelConfig) WithReplacement(replacement string) *RelabelConfig {
	r.Replacement = replacement
	return r
}

// WithModulus sets the modulus for hashmod action.
func (r *RelabelConfig) WithModulus(mod uint64) *RelabelConfig {
	r.Modulus = mod
	return r
}

// KeepByLabel creates a relabel config that keeps targets where the label matches the regex.
//
// Example:
//
//	// Keep only targets with env=production
//	KeepByLabel("env", "production")
func KeepByLabel(label, regex string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{label},
		Regex:        regex,
		Action:       string(RelabelKeep),
	}
}

// DropByLabel creates a relabel config that drops targets where the label matches the regex.
//
// Example:
//
//	// Drop targets with status=disabled
//	DropByLabel("status", "disabled")
func DropByLabel(label, regex string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{label},
		Regex:        regex,
		Action:       string(RelabelDrop),
	}
}

// LabelFromMeta creates a relabel config that extracts a label from __meta_* labels.
// This is commonly used to extract Kubernetes annotations/labels.
//
// Example:
//
//	// Extract app name from Kubernetes pod annotation
//	LabelFromMeta("__meta_kubernetes_pod_annotation_app", "app")
func LabelFromMeta(metaLabel, targetLabel string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{metaLabel},
		TargetLabel:  targetLabel,
		Action:       string(RelabelReplace),
	}
}

// RenameLabel creates a relabel config that renames a label.
//
// Example:
//
//	// Rename "kubernetes_pod_name" to "pod"
//	RenameLabel("kubernetes_pod_name", "pod")
func RenameLabel(from, to string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{from},
		TargetLabel:  to,
		Action:       string(RelabelReplace),
	}
}

// DropLabels creates a relabel config that drops labels matching the regex.
//
// Example:
//
//	// Drop all __meta_ labels
//	DropLabels("__meta_.*")
func DropLabels(regex string) *RelabelConfig {
	return &RelabelConfig{
		Regex:  regex,
		Action: string(RelabelLabelDrop),
	}
}

// KeepLabels creates a relabel config that keeps only labels matching the regex.
//
// Example:
//
//	// Keep only job, instance, and __* labels
//	KeepLabels("job|instance|__.*")
func KeepLabels(regex string) *RelabelConfig {
	return &RelabelConfig{
		Regex:  regex,
		Action: string(RelabelLabelKeep),
	}
}

// HashMod creates a relabel config for consistent hashing/sharding.
// It sets target_label to hash(source_labels) % modulus.
//
// Example:
//
//	// Shard targets across 3 Prometheus instances
//	HashMod("__address__", "__shard", 3)
func HashMod(sourceLabel, targetLabel string, modulus uint64) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{sourceLabel},
		TargetLabel:  targetLabel,
		Modulus:      modulus,
		Action:       string(RelabelHashMod),
	}
}

// Replace creates a relabel config that replaces a label value using regex.
//
// Example:
//
//	// Extract port from address and set as port label
//	Replace([]string{"__address__"}, "port", ".*:(.*)", "$1")
func Replace(sourceLabels []string, targetLabel, regex, replacement string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: sourceLabels,
		TargetLabel:  targetLabel,
		Regex:        regex,
		Replacement:  replacement,
		Action:       string(RelabelReplace),
	}
}

// LabelMap creates a relabel config that copies labels matching regex to new labels.
// The new label name is generated by the replacement string.
//
// Example:
//
//	// Copy __meta_kubernetes_pod_label_<name> to <name>
//	LabelMap("__meta_kubernetes_pod_label_(.+)", "$1")
func LabelMap(regex, replacement string) *RelabelConfig {
	return &RelabelConfig{
		Regex:       regex,
		Replacement: replacement,
		Action:      string(RelabelLabelMap),
	}
}

// KeepByAnnotation creates a relabel config that keeps targets with a specific
// Kubernetes annotation set to a specific value. Common for prometheus.io/scrape.
//
// Example:
//
//	// Keep pods with prometheus.io/scrape=true annotation
//	KeepByAnnotation("prometheus.io/scrape", "true")
func KeepByAnnotation(annotation, value string) *RelabelConfig {
	metaLabel := "__meta_kubernetes_pod_annotation_" + sanitizeAnnotation(annotation)
	return KeepByLabel(metaLabel, value)
}

// SetFromAnnotation creates a relabel config that sets a label from a Kubernetes annotation.
//
// Example:
//
//	// Set metrics_path from prometheus.io/path annotation
//	SetFromAnnotation("prometheus.io/path", "__metrics_path__")
func SetFromAnnotation(annotation, targetLabel string) *RelabelConfig {
	metaLabel := "__meta_kubernetes_pod_annotation_" + sanitizeAnnotation(annotation)
	return LabelFromMeta(metaLabel, targetLabel)
}

// sanitizeAnnotation converts annotation names to the format used in __meta_ labels.
// Replaces dots and slashes with underscores.
func sanitizeAnnotation(annotation string) string {
	result := make([]byte, len(annotation))
	for i := 0; i < len(annotation); i++ {
		c := annotation[i]
		if c == '.' || c == '/' || c == '-' {
			result[i] = '_'
		} else {
			result[i] = c
		}
	}
	return string(result)
}

// KeepByPodLabel creates a relabel config that keeps targets with a specific
// Kubernetes pod label.
//
// Example:
//
//	// Keep pods with app=nginx label
//	KeepByPodLabel("app", "nginx")
func KeepByPodLabel(labelKey, value string) *RelabelConfig {
	metaLabel := "__meta_kubernetes_pod_label_" + sanitizeAnnotation(labelKey)
	return KeepByLabel(metaLabel, value)
}

// SetPort creates a relabel config that sets the scrape port from an annotation.
// This is commonly used with prometheus.io/port annotation.
//
// Example:
//
//	// Set port from prometheus.io/port annotation
//	SetPort("prometheus.io/port")
func SetPort(annotation string) *RelabelConfig {
	metaLabel := "__meta_kubernetes_pod_annotation_" + sanitizeAnnotation(annotation)
	return &RelabelConfig{
		SourceLabels: []string{metaLabel, "__meta_kubernetes_pod_ip"},
		Separator:    ":",
		TargetLabel:  "__address__",
		Regex:        "(.+):(.+);(.+)",
		Replacement:  "$3:$1",
		Action:       string(RelabelReplace),
	}
}
