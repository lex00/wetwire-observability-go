package prometheus

import (
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestPrometheusConfig_MarshalYAML(t *testing.T) {
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
				StaticConfigs: []*StaticConfig{
					{
						Targets: []string{"localhost:9090"},
					},
				},
			},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Verify it contains expected fields
	yamlStr := string(data)
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
		"static_configs:",
		"targets:",
		"localhost:9090",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestPrometheusConfig_UnmarshalYAML(t *testing.T) {
	input := `
global:
  scrape_interval: 30s
  scrape_timeout: 10s
  evaluation_interval: 1m
  external_labels:
    environment: staging
scrape_configs:
  - job_name: node
    scrape_interval: 15s
    static_configs:
      - targets:
          - "node1:9100"
          - "node2:9100"
        labels:
          group: nodes
`
	var config PrometheusConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	// Check global config
	if config.Global == nil {
		t.Fatal("Global config is nil")
	}
	if config.Global.ScrapeInterval != Duration(30*time.Second) {
		t.Errorf("ScrapeInterval = %v, want %v", config.Global.ScrapeInterval, Duration(30*time.Second))
	}
	if config.Global.EvaluationInterval != Duration(time.Minute) {
		t.Errorf("EvaluationInterval = %v, want %v", config.Global.EvaluationInterval, Duration(time.Minute))
	}
	if config.Global.ExternalLabels["environment"] != "staging" {
		t.Errorf("ExternalLabels[environment] = %v, want staging", config.Global.ExternalLabels["environment"])
	}

	// Check scrape configs
	if len(config.ScrapeConfigs) != 1 {
		t.Fatalf("len(ScrapeConfigs) = %d, want 1", len(config.ScrapeConfigs))
	}
	sc := config.ScrapeConfigs[0]
	if sc.JobName != "node" {
		t.Errorf("JobName = %v, want node", sc.JobName)
	}
	if sc.ScrapeInterval != Duration(15*time.Second) {
		t.Errorf("ScrapeInterval = %v, want %v", sc.ScrapeInterval, Duration(15*time.Second))
	}
	if len(sc.StaticConfigs) != 1 {
		t.Fatalf("len(StaticConfigs) = %d, want 1", len(sc.StaticConfigs))
	}
	if len(sc.StaticConfigs[0].Targets) != 2 {
		t.Errorf("len(Targets) = %d, want 2", len(sc.StaticConfigs[0].Targets))
	}
}

func TestGlobalConfig_Empty(t *testing.T) {
	config := &GlobalConfig{}
	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// Empty global config should result in empty YAML object
	if string(data) != "{}\n" {
		t.Errorf("Empty GlobalConfig = %q, want {}", string(data))
	}
}

func TestAlertingConfig_MarshalYAML(t *testing.T) {
	config := &PrometheusConfig{
		Alerting: &AlertingConfig{
			Alertmanagers: []*AlertmanagerConfig{
				{
					StaticConfigs: []*StaticConfig{
						{Targets: []string{"alertmanager:9093"}},
					},
					Scheme:     "http",
					PathPrefix: "/alertmanager",
				},
			},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"alerting:",
		"alertmanagers:",
		"static_configs:",
		"alertmanager:9093",
		"scheme: http",
		"path_prefix: /alertmanager",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsHelper(s, substr))
}

func containsHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
