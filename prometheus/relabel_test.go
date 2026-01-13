package prometheus

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestRelabelAction_Constants(t *testing.T) {
	tests := []struct {
		action RelabelAction
		want   string
	}{
		{RelabelReplace, "replace"},
		{RelabelKeep, "keep"},
		{RelabelDrop, "drop"},
		{RelabelHashMod, "hashmod"},
		{RelabelLabelMap, "labelmap"},
		{RelabelLabelDrop, "labeldrop"},
		{RelabelLabelKeep, "labelkeep"},
		{RelabelLowercase, "lowercase"},
		{RelabelUppercase, "uppercase"},
		{RelabelKeepEqual, "keepequal"},
		{RelabelDropEqual, "dropequal"},
	}

	for _, tt := range tests {
		if string(tt.action) != tt.want {
			t.Errorf("RelabelAction %v = %q, want %q", tt.action, string(tt.action), tt.want)
		}
	}
}

func TestKeepByLabel(t *testing.T) {
	rc := KeepByLabel("env", "production")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "env" {
		t.Errorf("SourceLabels = %v, want [env]", rc.SourceLabels)
	}
	if rc.Regex != "production" {
		t.Errorf("Regex = %v, want production", rc.Regex)
	}
	if rc.Action != "keep" {
		t.Errorf("Action = %v, want keep", rc.Action)
	}
}

func TestDropByLabel(t *testing.T) {
	rc := DropByLabel("status", "disabled")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "status" {
		t.Errorf("SourceLabels = %v, want [status]", rc.SourceLabels)
	}
	if rc.Regex != "disabled" {
		t.Errorf("Regex = %v, want disabled", rc.Regex)
	}
	if rc.Action != "drop" {
		t.Errorf("Action = %v, want drop", rc.Action)
	}
}

func TestLabelFromMeta(t *testing.T) {
	rc := LabelFromMeta("__meta_kubernetes_pod_annotation_app", "app")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__meta_kubernetes_pod_annotation_app" {
		t.Errorf("SourceLabels = %v, want [__meta_kubernetes_pod_annotation_app]", rc.SourceLabels)
	}
	if rc.TargetLabel != "app" {
		t.Errorf("TargetLabel = %v, want app", rc.TargetLabel)
	}
	if rc.Action != "replace" {
		t.Errorf("Action = %v, want replace", rc.Action)
	}
}

func TestRenameLabel(t *testing.T) {
	rc := RenameLabel("kubernetes_pod_name", "pod")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "kubernetes_pod_name" {
		t.Errorf("SourceLabels = %v, want [kubernetes_pod_name]", rc.SourceLabels)
	}
	if rc.TargetLabel != "pod" {
		t.Errorf("TargetLabel = %v, want pod", rc.TargetLabel)
	}
	if rc.Action != "replace" {
		t.Errorf("Action = %v, want replace", rc.Action)
	}
}

func TestDropLabels(t *testing.T) {
	rc := DropLabels("__meta_.*")

	if rc.Regex != "__meta_.*" {
		t.Errorf("Regex = %v, want __meta_.*", rc.Regex)
	}
	if rc.Action != "labeldrop" {
		t.Errorf("Action = %v, want labeldrop", rc.Action)
	}
}

func TestKeepLabels(t *testing.T) {
	rc := KeepLabels("job|instance")

	if rc.Regex != "job|instance" {
		t.Errorf("Regex = %v, want job|instance", rc.Regex)
	}
	if rc.Action != "labelkeep" {
		t.Errorf("Action = %v, want labelkeep", rc.Action)
	}
}

func TestHashMod(t *testing.T) {
	rc := HashMod("__address__", "__shard", 3)

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__address__" {
		t.Errorf("SourceLabels = %v, want [__address__]", rc.SourceLabels)
	}
	if rc.TargetLabel != "__shard" {
		t.Errorf("TargetLabel = %v, want __shard", rc.TargetLabel)
	}
	if rc.Modulus != 3 {
		t.Errorf("Modulus = %v, want 3", rc.Modulus)
	}
	if rc.Action != "hashmod" {
		t.Errorf("Action = %v, want hashmod", rc.Action)
	}
}

func TestReplace(t *testing.T) {
	rc := Replace([]string{"__address__"}, "port", ".*:(.*)", "$1")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__address__" {
		t.Errorf("SourceLabels = %v, want [__address__]", rc.SourceLabels)
	}
	if rc.TargetLabel != "port" {
		t.Errorf("TargetLabel = %v, want port", rc.TargetLabel)
	}
	if rc.Regex != ".*:(.*)" {
		t.Errorf("Regex = %v, want .*:(.*)", rc.Regex)
	}
	if rc.Replacement != "$1" {
		t.Errorf("Replacement = %v, want $1", rc.Replacement)
	}
	if rc.Action != "replace" {
		t.Errorf("Action = %v, want replace", rc.Action)
	}
}

func TestLabelMap(t *testing.T) {
	rc := LabelMap("__meta_kubernetes_pod_label_(.+)", "$1")

	if rc.Regex != "__meta_kubernetes_pod_label_(.+)" {
		t.Errorf("Regex = %v, want __meta_kubernetes_pod_label_(.+)", rc.Regex)
	}
	if rc.Replacement != "$1" {
		t.Errorf("Replacement = %v, want $1", rc.Replacement)
	}
	if rc.Action != "labelmap" {
		t.Errorf("Action = %v, want labelmap", rc.Action)
	}
}

func TestKeepByAnnotation(t *testing.T) {
	rc := KeepByAnnotation("prometheus.io/scrape", "true")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__meta_kubernetes_pod_annotation_prometheus_io_scrape" {
		t.Errorf("SourceLabels = %v, want [__meta_kubernetes_pod_annotation_prometheus_io_scrape]", rc.SourceLabels)
	}
	if rc.Regex != "true" {
		t.Errorf("Regex = %v, want true", rc.Regex)
	}
	if rc.Action != "keep" {
		t.Errorf("Action = %v, want keep", rc.Action)
	}
}

func TestSetFromAnnotation(t *testing.T) {
	rc := SetFromAnnotation("prometheus.io/path", "__metrics_path__")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__meta_kubernetes_pod_annotation_prometheus_io_path" {
		t.Errorf("SourceLabels = %v, want [__meta_kubernetes_pod_annotation_prometheus_io_path]", rc.SourceLabels)
	}
	if rc.TargetLabel != "__metrics_path__" {
		t.Errorf("TargetLabel = %v, want __metrics_path__", rc.TargetLabel)
	}
}

func TestKeepByPodLabel(t *testing.T) {
	rc := KeepByPodLabel("app", "nginx")

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__meta_kubernetes_pod_label_app" {
		t.Errorf("SourceLabels = %v, want [__meta_kubernetes_pod_label_app]", rc.SourceLabels)
	}
	if rc.Regex != "nginx" {
		t.Errorf("Regex = %v, want nginx", rc.Regex)
	}
}

func TestSanitizeAnnotation(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"prometheus.io/scrape", "prometheus_io_scrape"},
		{"app", "app"},
		{"app.kubernetes.io/name", "app_kubernetes_io_name"},
		{"my-annotation", "my_annotation"},
	}

	for _, tt := range tests {
		got := sanitizeAnnotation(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeAnnotation(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestRelabelConfig_Serialize_KeepByLabel(t *testing.T) {
	rc := KeepByLabel("env", "production")

	data, err := yaml.Marshal(rc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"source_labels:",
		"env",
		"regex: production",
		"action: keep",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRelabelConfig_Serialize_HashMod(t *testing.T) {
	rc := HashMod("__address__", "__shard", 5)

	data, err := yaml.Marshal(rc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"source_labels:",
		"__address__",
		"modulus: 5",
		"target_label: __shard",
		"action: hashmod",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRelabelConfig_Serialize_LabelMap(t *testing.T) {
	rc := LabelMap("__meta_kubernetes_pod_label_(.+)", "$1")

	data, err := yaml.Marshal(rc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"regex: __meta_kubernetes_pod_label_(.+)",
		"replacement: $1",
		"action: labelmap",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestNewRelabelConfig_FluentAPI(t *testing.T) {
	rc := NewRelabelConfig(RelabelReplace).
		WithSourceLabels("__address__").
		WithRegex("(.*):\\d+").
		WithTargetLabel("instance").
		WithReplacement("$1")

	if rc.Action != "replace" {
		t.Errorf("Action = %v, want replace", rc.Action)
	}
	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__address__" {
		t.Errorf("SourceLabels = %v, want [__address__]", rc.SourceLabels)
	}
	if rc.Regex != "(.*):\\d+" {
		t.Errorf("Regex = %v, want (.*):\\d+", rc.Regex)
	}
	if rc.TargetLabel != "instance" {
		t.Errorf("TargetLabel = %v, want instance", rc.TargetLabel)
	}
	if rc.Replacement != "$1" {
		t.Errorf("Replacement = %v, want $1", rc.Replacement)
	}
}

func TestScrapeConfig_WithRelabelConfigs(t *testing.T) {
	sc := NewScrapeConfig("kubernetes-pods").
		WithKubernetesSD(NewKubernetesSD(KubernetesRolePod))

	sc.RelabelConfigs = []*RelabelConfig{
		KeepByAnnotation("prometheus.io/scrape", "true"),
		SetFromAnnotation("prometheus.io/path", "__metrics_path__"),
		LabelMap("__meta_kubernetes_pod_label_(.+)", "$1"),
		DropLabels("__meta_.*"),
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"job_name: kubernetes-pods",
		"relabel_configs:",
		"action: keep",
		"action: replace",
		"action: labelmap",
		"action: labeldrop",
		"__meta_kubernetes_pod_annotation_prometheus_io_scrape",
		"__metrics_path__",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRelabelConfig_Unmarshal(t *testing.T) {
	input := `
source_labels:
  - __address__
target_label: instance
regex: "(.*):\\d+"
replacement: "$1"
action: replace
`
	var rc RelabelConfig
	if err := yaml.Unmarshal([]byte(input), &rc); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(rc.SourceLabels) != 1 || rc.SourceLabels[0] != "__address__" {
		t.Errorf("SourceLabels = %v, want [__address__]", rc.SourceLabels)
	}
	if rc.TargetLabel != "instance" {
		t.Errorf("TargetLabel = %v, want instance", rc.TargetLabel)
	}
	if rc.Regex != "(.*):\\d+" {
		t.Errorf("Regex = %v, want (.*):\\d+", rc.Regex)
	}
	if rc.Replacement != "$1" {
		t.Errorf("Replacement = %v, want $1", rc.Replacement)
	}
	if rc.Action != "replace" {
		t.Errorf("Action = %v, want replace", rc.Action)
	}
}

func TestRelabelConfig_ComplexKubernetesExample(t *testing.T) {
	// This test verifies a complete Kubernetes pod discovery relabeling configuration
	sc := &ScrapeConfig{
		JobName: "kubernetes-pods",
		KubernetesSDConfigs: []*KubernetesSD{
			NewKubernetesSD(KubernetesRolePod),
		},
		RelabelConfigs: []*RelabelConfig{
			// Keep only pods with prometheus.io/scrape=true
			KeepByAnnotation("prometheus.io/scrape", "true"),
			// Set metrics path from annotation
			SetFromAnnotation("prometheus.io/path", "__metrics_path__"),
			// Set scheme from annotation (default http)
			{
				SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_scheme"},
				TargetLabel:  "__scheme__",
				Regex:        "(https?)",
				Action:       string(RelabelReplace),
			},
			// Set address:port from annotation
			Replace(
				[]string{"__address__", "__meta_kubernetes_pod_annotation_prometheus_io_port"},
				"__address__",
				"([^:]+)(?::\\d+)?;(\\d+)",
				"$1:$2",
			),
			// Copy pod labels
			LabelMap("__meta_kubernetes_pod_label_(.+)", "$1"),
			// Set namespace label
			LabelFromMeta("__meta_kubernetes_namespace", "namespace"),
			// Set pod label
			LabelFromMeta("__meta_kubernetes_pod_name", "pod"),
			// Drop all __meta labels
			DropLabels("__meta_.*"),
		},
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Just verify it produces valid YAML
	var restored ScrapeConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(restored.RelabelConfigs) != 8 {
		t.Errorf("len(RelabelConfigs) = %d, want 8", len(restored.RelabelConfigs))
	}

	t.Logf("Generated relabel config:\n%s", string(data))
}
