package prometheus

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestPrometheusConfig_Serialize(t *testing.T) {
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval:     Duration(30 * time.Second),
			ScrapeTimeout:      Duration(10 * time.Second),
			EvaluationInterval: Duration(30 * time.Second),
			ExternalLabels: map[string]string{
				"environment": "production",
				"cluster":     "main",
			},
		},
		ScrapeConfigs: []*ScrapeConfig{
			{
				JobName:        "prometheus",
				ScrapeInterval: Duration(15 * time.Second),
				MetricsPath:    "/metrics",
				StaticConfigs: []*StaticConfig{
					{
						Targets: []string{"localhost:9090"},
					},
				},
			},
			{
				JobName: "node",
				StaticConfigs: []*StaticConfig{
					{
						Targets: []string{"node1:9100", "node2:9100"},
						Labels: map[string]string{
							"group": "nodes",
						},
					},
				},
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)

	// Verify key fields are present
	expectations := []string{
		"global:",
		"scrape_interval: 30s",
		"scrape_timeout: 10s",
		"evaluation_interval: 30s",
		"external_labels:",
		"environment: production",
		"cluster: main",
		"scrape_configs:",
		"job_name: prometheus",
		"scrape_interval: 15s",
		"metrics_path: /metrics",
		"job_name: node",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Serialize() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPrometheusConfig_RoundTrip(t *testing.T) {
	original := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval:     Duration(15 * time.Second),
			ScrapeTimeout:      Duration(10 * time.Second),
			EvaluationInterval: Duration(1 * time.Minute),
			ExternalLabels: map[string]string{
				"env": "staging",
			},
		},
		ScrapeConfigs: []*ScrapeConfig{
			{
				JobName:        "api",
				ScrapeInterval: Duration(30 * time.Second),
				ScrapeTimeout:  Duration(5 * time.Second),
				MetricsPath:    "/metrics",
				Scheme:         "https",
				HonorLabels:    true,
				StaticConfigs: []*StaticConfig{
					{
						Targets: []string{"api1:8080", "api2:8080"},
						Labels: map[string]string{
							"team": "platform",
						},
					},
				},
			},
		},
		RuleFiles: []string{
			"rules/*.yml",
			"alerts/*.yml",
		},
	}

	// Serialize
	data, err := original.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Deserialize
	var restored PrometheusConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	// Compare key fields
	if restored.Global.ScrapeInterval != original.Global.ScrapeInterval {
		t.Errorf("ScrapeInterval = %v, want %v", restored.Global.ScrapeInterval, original.Global.ScrapeInterval)
	}
	if restored.Global.ExternalLabels["env"] != "staging" {
		t.Errorf("ExternalLabels[env] = %v, want staging", restored.Global.ExternalLabels["env"])
	}
	if len(restored.ScrapeConfigs) != 1 {
		t.Fatalf("len(ScrapeConfigs) = %d, want 1", len(restored.ScrapeConfigs))
	}
	if restored.ScrapeConfigs[0].JobName != "api" {
		t.Errorf("JobName = %v, want api", restored.ScrapeConfigs[0].JobName)
	}
	if len(restored.RuleFiles) != 2 {
		t.Errorf("len(RuleFiles) = %d, want 2", len(restored.RuleFiles))
	}
}

func TestPrometheusConfig_SerializeToFile(t *testing.T) {
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval: Duration(30 * time.Second),
		},
		ScrapeConfigs: []*ScrapeConfig{
			{
				JobName: "test",
				StaticConfigs: []*StaticConfig{
					{Targets: []string{"localhost:9090"}},
				},
			},
		},
	}

	// Create temp directory
	tmpDir := t.TempDir()
	path := filepath.Join(tmpDir, "prometheus.yml")

	// Write to file
	if err := config.SerializeToFile(path); err != nil {
		t.Fatalf("SerializeToFile() error = %v", err)
	}

	// Read back and verify
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "job_name: test") {
		t.Errorf("File content missing expected data\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "scrape_interval: 30s") {
		t.Errorf("File content missing scrape_interval\nGot:\n%s", yamlStr)
	}
}

func TestPrometheusConfig_EmptyConfig(t *testing.T) {
	config := &PrometheusConfig{}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Empty config should produce minimal YAML
	yamlStr := string(data)
	if yamlStr != "{}\n" {
		t.Errorf("Empty config = %q, want {}", yamlStr)
	}
}

func TestPrometheusConfig_NilFields(t *testing.T) {
	config := &PrometheusConfig{
		Global: nil,
		ScrapeConfigs: []*ScrapeConfig{
			{
				JobName: "minimal",
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	if strings.Contains(yamlStr, "global:") {
		t.Errorf("Nil global should not appear in output\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "job_name: minimal") {
		t.Errorf("Missing job_name\nGot:\n%s", yamlStr)
	}
}

func TestPrometheusConfig_SpecialCharacters(t *testing.T) {
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ExternalLabels: map[string]string{
				"path":        "/api/v1",
				"description": "Test with \"quotes\"",
				"multiline":   "line1\nline2",
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Verify it's valid YAML by unmarshaling
	var restored PrometheusConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if restored.Global.ExternalLabels["path"] != "/api/v1" {
		t.Errorf("path = %v, want /api/v1", restored.Global.ExternalLabels["path"])
	}
	if restored.Global.ExternalLabels["description"] != "Test with \"quotes\"" {
		t.Errorf("description not preserved correctly")
	}
}

func TestPrometheusConfig_MustSerialize(t *testing.T) {
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval: Duration(30 * time.Second),
		},
	}

	// Should not panic
	data := config.MustSerialize()
	if len(data) == 0 {
		t.Error("MustSerialize() returned empty data")
	}
}

func TestGlobalConfig_Serialize(t *testing.T) {
	config := &GlobalConfig{
		ScrapeInterval:     Duration(30 * time.Second),
		EvaluationInterval: Duration(1 * time.Minute),
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "scrape_interval: 30s") {
		t.Errorf("Missing scrape_interval\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "evaluation_interval: 1m") {
		t.Errorf("Missing evaluation_interval\nGot:\n%s", yamlStr)
	}
}

func TestScrapeConfig_Serialize(t *testing.T) {
	config := &ScrapeConfig{
		JobName:        "api",
		ScrapeInterval: Duration(15 * time.Second),
		StaticConfigs: []*StaticConfig{
			{
				Targets: []string{"localhost:8080"},
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "job_name: api") {
		t.Errorf("Missing job_name\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "scrape_interval: 15s") {
		t.Errorf("Missing scrape_interval\nGot:\n%s", yamlStr)
	}
}

func TestPrometheusConfig_AlertingConfig(t *testing.T) {
	config := &PrometheusConfig{
		Alerting: &AlertingConfig{
			Alertmanagers: []*AlertmanagerConfig{
				{
					StaticConfigs: []*StaticConfig{
						{Targets: []string{"alertmanager:9093"}},
					},
					Scheme:     "http",
					PathPrefix: "/alertmanager",
					Timeout:    Duration(10 * time.Second),
				},
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"alerting:",
		"alertmanagers:",
		"alertmanager:9093",
		"scheme: http",
		"path_prefix: /alertmanager",
		"timeout: 10s",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPrometheusConfig_RemoteWrite(t *testing.T) {
	config := &PrometheusConfig{
		RemoteWrite: []*RemoteWriteConfig{
			{
				URL:           "http://remote:9090/api/v1/write",
				Name:          "remote1",
				RemoteTimeout: Duration(30 * time.Second),
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"remote_write:",
		"url: http://remote:9090/api/v1/write",
		"name: remote1",
		"remote_timeout: 30s",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("Missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPrometheusConfig_CompleteExample(t *testing.T) {
	// Create a complete, realistic Prometheus configuration
	config := &PrometheusConfig{
		Global: &GlobalConfig{
			ScrapeInterval:     Duration(15 * time.Second),
			ScrapeTimeout:      Duration(10 * time.Second),
			EvaluationInterval: Duration(15 * time.Second),
			ExternalLabels: map[string]string{
				"monitor":     "prometheus",
				"environment": "production",
			},
		},
		RuleFiles: []string{
			"rules/recording.yml",
			"rules/alerting.yml",
		},
		Alerting: &AlertingConfig{
			Alertmanagers: []*AlertmanagerConfig{
				{
					StaticConfigs: []*StaticConfig{
						{Targets: []string{"alertmanager:9093"}},
					},
				},
			},
		},
		ScrapeConfigs: []*ScrapeConfig{
			{
				JobName: "prometheus",
				StaticConfigs: []*StaticConfig{
					{Targets: []string{"localhost:9090"}},
				},
			},
			{
				JobName:        "node-exporter",
				ScrapeInterval: Duration(30 * time.Second),
				StaticConfigs: []*StaticConfig{
					{
						Targets: []string{"node1:9100", "node2:9100", "node3:9100"},
						Labels: map[string]string{
							"group": "infrastructure",
						},
					},
				},
				RelabelConfigs: []*RelabelConfig{
					{
						SourceLabels: []string{"__address__"},
						TargetLabel:  "instance",
						Regex:        "(.*):\\d+",
						Replacement:  "$1",
					},
				},
			},
			{
				JobName:     "api-servers",
				MetricsPath: "/metrics",
				Scheme:      "https",
				TLSConfig: &TLSConfig{
					CAFile:             "/etc/prometheus/ca.crt",
					InsecureSkipVerify: false,
				},
				StaticConfigs: []*StaticConfig{
					{
						Targets: []string{"api1:8443", "api2:8443"},
						Labels: map[string]string{
							"group": "api",
							"tier":  "frontend",
						},
					},
				},
			},
		},
	}

	data, err := config.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	// Verify it's valid YAML
	var restored PrometheusConfig
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	// Verify structure
	if len(restored.ScrapeConfigs) != 3 {
		t.Errorf("len(ScrapeConfigs) = %d, want 3", len(restored.ScrapeConfigs))
	}
	if len(restored.RuleFiles) != 2 {
		t.Errorf("len(RuleFiles) = %d, want 2", len(restored.RuleFiles))
	}

	// Log output for visual inspection
	t.Logf("Generated prometheus.yml:\n%s", string(data))
}
