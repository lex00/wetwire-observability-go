package importer

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/lex00/wetwire-observability-go/prometheus"
)

func TestParsePrometheusConfigFromBytes_Simple(t *testing.T) {
	yaml := `
global:
  scrape_interval: 15s
  evaluation_interval: 15s
  external_labels:
    cluster: production

scrape_configs:
  - job_name: prometheus
    static_configs:
      - targets:
          - localhost:9090
`
	config, err := ParsePrometheusConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePrometheusConfigFromBytes() error = %v", err)
	}

	if config.Global == nil {
		t.Fatal("Global config is nil")
	}
	if config.Global.ScrapeInterval != 15*prometheus.Second {
		t.Errorf("ScrapeInterval = %v, want 15s", config.Global.ScrapeInterval)
	}
	if len(config.ScrapeConfigs) != 1 {
		t.Errorf("len(ScrapeConfigs) = %d, want 1", len(config.ScrapeConfigs))
	}
	if config.ScrapeConfigs[0].JobName != "prometheus" {
		t.Errorf("JobName = %v, want prometheus", config.ScrapeConfigs[0].JobName)
	}
}

func TestParsePrometheusConfigFromBytes_WithKubernetesSD(t *testing.T) {
	yaml := `
scrape_configs:
  - job_name: kubernetes-pods
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - production
            - staging
        selectors:
          - role: pod
            label: "app=nginx"
`
	config, err := ParsePrometheusConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePrometheusConfigFromBytes() error = %v", err)
	}

	if len(config.ScrapeConfigs) != 1 {
		t.Fatalf("len(ScrapeConfigs) = %d, want 1", len(config.ScrapeConfigs))
	}

	sc := config.ScrapeConfigs[0]
	if len(sc.KubernetesSDConfigs) != 1 {
		t.Fatalf("len(KubernetesSDConfigs) = %d, want 1", len(sc.KubernetesSDConfigs))
	}

	k8s := sc.KubernetesSDConfigs[0]
	if k8s.Role != prometheus.KubernetesRolePod {
		t.Errorf("Role = %v, want pod", k8s.Role)
	}
	if len(k8s.Namespaces.Names) != 2 {
		t.Errorf("len(Namespaces.Names) = %d, want 2", len(k8s.Namespaces.Names))
	}
}

func TestParsePrometheusConfigFromBytes_WithRemoteWrite(t *testing.T) {
	yaml := `
remote_write:
  - url: http://thanos-receive:10908/api/v1/receive
    name: thanos
    remote_timeout: 30s
    queue_config:
      capacity: 10000
      max_shards: 50
`
	config, err := ParsePrometheusConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePrometheusConfigFromBytes() error = %v", err)
	}

	if len(config.RemoteWrite) != 1 {
		t.Fatalf("len(RemoteWrite) = %d, want 1", len(config.RemoteWrite))
	}

	rw := config.RemoteWrite[0]
	if rw.URL != "http://thanos-receive:10908/api/v1/receive" {
		t.Errorf("URL = %v, want http://thanos-receive:10908/api/v1/receive", rw.URL)
	}
	if rw.Name != "thanos" {
		t.Errorf("Name = %v, want thanos", rw.Name)
	}
	if rw.QueueConfig == nil || rw.QueueConfig.Capacity != 10000 {
		t.Error("QueueConfig not parsed correctly")
	}
}

func TestParsePrometheusConfigFromBytes_Invalid(t *testing.T) {
	yaml := `
invalid: yaml: content:
  - this is not valid
`
	_, err := ParsePrometheusConfigFromBytes([]byte(yaml))
	if err == nil {
		t.Error("ParsePrometheusConfigFromBytes() should return error for invalid YAML")
	}
}

func TestGenerateGoCode_Simple(t *testing.T) {
	config := &prometheus.PrometheusConfig{
		Global: &prometheus.GlobalConfig{
			ScrapeInterval:     15 * prometheus.Second,
			EvaluationInterval: 15 * prometheus.Second,
		},
		ScrapeConfigs: []*prometheus.ScrapeConfig{
			{
				JobName: "prometheus",
				StaticConfigs: []*prometheus.StaticConfig{
					{Targets: []string{"localhost:9090"}},
				},
			},
		},
	}

	code, err := GenerateGoCode(config, "monitoring")
	if err != nil {
		t.Fatalf("GenerateGoCode() error = %v", err)
	}

	codeStr := string(code)

	// Check for expected elements
	expectations := []string{
		"package monitoring",
		`"github.com/lex00/wetwire-observability-go/prometheus"`,
		"var GlobalConfig = &prometheus.GlobalConfig{",
		"ScrapeInterval:",
		"prometheus.NewScrapeConfig(",
		`"prometheus"`,
		"WithStaticTargets",
		"var Config = &prometheus.PrometheusConfig{",
	}

	for _, exp := range expectations {
		if !strings.Contains(codeStr, exp) {
			t.Errorf("GenerateGoCode() missing %q\nGot:\n%s", exp, codeStr)
		}
	}
}

func TestGenerateGoCode_WithKubernetesSD(t *testing.T) {
	config := &prometheus.PrometheusConfig{
		ScrapeConfigs: []*prometheus.ScrapeConfig{
			{
				JobName: "kubernetes-pods",
				KubernetesSDConfigs: []*prometheus.KubernetesSD{
					{
						Role: prometheus.KubernetesRolePod,
						Namespaces: &prometheus.KubernetesNamespaceDiscovery{
							Names: []string{"production"},
						},
					},
				},
			},
		},
	}

	code, err := GenerateGoCode(config, "monitoring")
	if err != nil {
		t.Fatalf("GenerateGoCode() error = %v", err)
	}

	codeStr := string(code)

	expectations := []string{
		"prometheus.NewKubernetesSD(prometheus.KubernetesRolePod)",
		"WithNamespaces",
		`"production"`,
	}

	for _, exp := range expectations {
		if !strings.Contains(codeStr, exp) {
			t.Errorf("GenerateGoCode() missing %q\nGot:\n%s", exp, codeStr)
		}
	}
}

func TestGenerateGoCode_WithRemoteWrite(t *testing.T) {
	config := &prometheus.PrometheusConfig{
		RemoteWrite: []*prometheus.RemoteWriteConfig{
			{
				URL:           "http://thanos:10908/api/v1/receive",
				Name:          "thanos",
				RemoteTimeout: 30 * prometheus.Second,
			},
		},
	}

	code, err := GenerateGoCode(config, "monitoring")
	if err != nil {
		t.Fatalf("GenerateGoCode() error = %v", err)
	}

	codeStr := string(code)

	expectations := []string{
		"prometheus.NewRemoteWrite",
		`"http://thanos:10908/api/v1/receive"`,
		"WithName",
		`"thanos"`,
		"WithTimeout",
	}

	for _, exp := range expectations {
		if !strings.Contains(codeStr, exp) {
			t.Errorf("GenerateGoCode() missing %q\nGot:\n%s", exp, codeStr)
		}
	}
}

func TestGenerateGoCode_Compiles(t *testing.T) {
	config := &prometheus.PrometheusConfig{
		Global: &prometheus.GlobalConfig{
			ScrapeInterval:     15 * prometheus.Second,
			EvaluationInterval: 15 * prometheus.Second,
			ExternalLabels: map[string]string{
				"cluster": "production",
				"env":     "prod",
			},
		},
		ScrapeConfigs: []*prometheus.ScrapeConfig{
			{
				JobName: "prometheus",
				StaticConfigs: []*prometheus.StaticConfig{
					{Targets: []string{"localhost:9090"}},
				},
			},
			{
				JobName: "kubernetes-pods",
				KubernetesSDConfigs: []*prometheus.KubernetesSD{
					{
						Role: prometheus.KubernetesRolePod,
						Namespaces: &prometheus.KubernetesNamespaceDiscovery{
							Names: []string{"production"},
						},
					},
				},
			},
		},
		RemoteWrite: []*prometheus.RemoteWriteConfig{
			{
				URL:           "http://thanos:10908/api/v1/receive",
				Name:          "thanos",
				RemoteTimeout: 30 * prometheus.Second,
			},
		},
	}

	code, err := GenerateGoCode(config, "generated")
	if err != nil {
		t.Fatalf("GenerateGoCode() error = %v", err)
	}

	// Write to temp file and try to compile
	tmpDir, err := os.MkdirTemp("", "codegen-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Write go.mod
	goMod := `module generated

go 1.21

require github.com/lex00/wetwire-observability-go v0.0.0

replace github.com/lex00/wetwire-observability-go => ` + filepath.Join(os.Getenv("PWD"), "../..")
	if err := os.WriteFile(filepath.Join(tmpDir, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("Failed to write go.mod: %v", err)
	}

	// Write generated code
	if err := os.WriteFile(filepath.Join(tmpDir, "config.go"), code, 0644); err != nil {
		t.Fatalf("Failed to write generated code: %v", err)
	}

	// Run go mod tidy first
	tidyCmd := exec.Command("go", "mod", "tidy")
	tidyCmd.Dir = tmpDir
	if output, err := tidyCmd.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy failed: %v\n%s", err, output)
	}

	// Try to build
	cmd := exec.Command("go", "build", ".")
	cmd.Dir = tmpDir
	output, err := cmd.CombinedOutput()
	if err != nil {
		t.Errorf("Generated code does not compile:\n%s\nCode:\n%s", output, code)
	}
}

func TestValidatePrometheusConfig(t *testing.T) {
	config := &prometheus.PrometheusConfig{
		ScrapeConfigs: []*prometheus.ScrapeConfig{
			{JobName: ""}, // Missing job name
			{
				JobName: "k8s-pods",
				KubernetesSDConfigs: []*prometheus.KubernetesSD{
					{Role: ""}, // Missing role
				},
			},
		},
	}

	warnings := ValidatePrometheusConfig(config)
	if len(warnings) != 2 {
		t.Errorf("ValidatePrometheusConfig() returned %d warnings, want 2", len(warnings))
	}
}

func TestRoundTrip(t *testing.T) {
	// Parse YAML
	yaml := `
global:
  scrape_interval: 30s
  evaluation_interval: 30s
  external_labels:
    cluster: test

scrape_configs:
  - job_name: api-server
    scrape_interval: 15s
    static_configs:
      - targets:
          - localhost:8080
          - localhost:8081

  - job_name: kubernetes-pods
    kubernetes_sd_configs:
      - role: pod
        namespaces:
          names:
            - default
`
	config, err := ParsePrometheusConfigFromBytes([]byte(yaml))
	if err != nil {
		t.Fatalf("ParsePrometheusConfigFromBytes() error = %v", err)
	}

	// Generate Go code
	code, err := GenerateGoCode(config, "roundtrip")
	if err != nil {
		t.Fatalf("GenerateGoCode() error = %v", err)
	}

	// Verify key elements are preserved
	codeStr := string(code)
	expectations := []string{
		"30 * prometheus.Second", // scrape_interval
		"api-server",
		"kubernetes-pods",
		"KubernetesRolePod",
		`"default"`,
	}

	for _, exp := range expectations {
		if !strings.Contains(codeStr, exp) {
			t.Errorf("Round-trip missing %q\nGot:\n%s", exp, codeStr)
		}
	}

	t.Logf("Generated code:\n%s", code)
}

func TestSanitizeVarName(t *testing.T) {
	gen := &codeGenerator{}

	tests := []struct {
		input string
		want  string
	}{
		{"prometheus", "Prometheus"},
		{"api-server", "ApiServer"},
		{"kubernetes_pods", "KubernetesPods"},
		{"my.job.name", "MyJobName"},
		{"namespace/service", "NamespaceService"},
		{"", "Config"},
	}

	for _, tt := range tests {
		got := gen.sanitizeVarName(tt.input)
		if got != tt.want {
			t.Errorf("sanitizeVarName(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatDuration(t *testing.T) {
	gen := &codeGenerator{}

	tests := []struct {
		input prometheus.Duration
		want  string
	}{
		{prometheus.Second, "prometheus.Second"},
		{15 * prometheus.Second, "15 * prometheus.Second"},
		{prometheus.Minute, "prometheus.Minute"},
		{5 * prometheus.Minute, "5 * prometheus.Minute"},
		{prometheus.Hour, "prometheus.Hour"},
		{2 * prometheus.Hour, "2 * prometheus.Hour"},
	}

	for _, tt := range tests {
		got := gen.formatDuration(tt.input)
		if got != tt.want {
			t.Errorf("formatDuration(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
