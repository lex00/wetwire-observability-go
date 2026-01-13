package prometheus

import (
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestScrapeConfig_MarshalYAML(t *testing.T) {
	config := &ScrapeConfig{
		JobName:        "api-server",
		ScrapeInterval: Duration(30 * time.Second),
		ScrapeTimeout:  Duration(10 * time.Second),
		MetricsPath:    "/metrics",
		Scheme:         "https",
		HonorLabels:    true,
		StaticConfigs: []*StaticConfig{
			{
				Targets: []string{"api1:8080", "api2:8080"},
				Labels: map[string]string{
					"env":  "production",
					"team": "platform",
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
		"job_name: api-server",
		"scrape_interval: 30s",
		"scrape_timeout: 10s",
		"metrics_path: /metrics",
		"scheme: https",
		"honor_labels: true",
		"static_configs:",
		"targets:",
		"api1:8080",
		"api2:8080",
		"labels:",
		"env: production",
		"team: platform",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestScrapeConfig_UnmarshalYAML(t *testing.T) {
	input := `
job_name: node-exporter
scrape_interval: 15s
scrape_timeout: 5s
metrics_path: /metrics
scheme: http
honor_labels: false
params:
  module:
    - http_2xx
static_configs:
  - targets:
      - "localhost:9100"
    labels:
      instance: local
`
	var config ScrapeConfig
	if err := yaml.Unmarshal([]byte(input), &config); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if config.JobName != "node-exporter" {
		t.Errorf("JobName = %v, want node-exporter", config.JobName)
	}
	if config.ScrapeInterval != Duration(15*time.Second) {
		t.Errorf("ScrapeInterval = %v, want 15s", config.ScrapeInterval)
	}
	if config.ScrapeTimeout != Duration(5*time.Second) {
		t.Errorf("ScrapeTimeout = %v, want 5s", config.ScrapeTimeout)
	}
	if config.MetricsPath != "/metrics" {
		t.Errorf("MetricsPath = %v, want /metrics", config.MetricsPath)
	}
	if config.HonorLabels != false {
		t.Errorf("HonorLabels = %v, want false", config.HonorLabels)
	}
	if len(config.Params) != 1 {
		t.Errorf("len(Params) = %d, want 1", len(config.Params))
	}
	if config.Params["module"][0] != "http_2xx" {
		t.Errorf("Params[module] = %v, want http_2xx", config.Params["module"])
	}
}

func TestStaticConfig_MarshalYAML(t *testing.T) {
	config := &StaticConfig{
		Targets: []string{"host1:9090", "host2:9090", "host3:9090"},
		Labels: map[string]string{
			"region": "us-west-2",
			"az":     "us-west-2a",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"targets:",
		"host1:9090",
		"host2:9090",
		"host3:9090",
		"labels:",
		"region: us-west-2",
		"az: us-west-2a",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestNewScrapeConfig(t *testing.T) {
	sc := NewScrapeConfig("my-job")
	if sc.JobName != "my-job" {
		t.Errorf("JobName = %v, want my-job", sc.JobName)
	}
}

func TestScrapeConfig_WithInterval(t *testing.T) {
	sc := NewScrapeConfig("test").WithInterval(Duration(30 * time.Second))
	if sc.ScrapeInterval != Duration(30*time.Second) {
		t.Errorf("ScrapeInterval = %v, want 30s", sc.ScrapeInterval)
	}
}

func TestScrapeConfig_WithTimeout(t *testing.T) {
	sc := NewScrapeConfig("test").WithTimeout(Duration(10 * time.Second))
	if sc.ScrapeTimeout != Duration(10*time.Second) {
		t.Errorf("ScrapeTimeout = %v, want 10s", sc.ScrapeTimeout)
	}
}

func TestScrapeConfig_WithStaticTargets(t *testing.T) {
	sc := NewScrapeConfig("test").
		WithStaticTargets("host1:9090", "host2:9090")

	if len(sc.StaticConfigs) != 1 {
		t.Fatalf("len(StaticConfigs) = %d, want 1", len(sc.StaticConfigs))
	}
	if len(sc.StaticConfigs[0].Targets) != 2 {
		t.Errorf("len(Targets) = %d, want 2", len(sc.StaticConfigs[0].Targets))
	}
}

func TestNewStaticConfig(t *testing.T) {
	sc := NewStaticConfig("host1:9090", "host2:9090")
	if len(sc.Targets) != 2 {
		t.Errorf("len(Targets) = %d, want 2", len(sc.Targets))
	}
}

func TestStaticConfig_WithLabels(t *testing.T) {
	sc := NewStaticConfig("localhost:9090").
		WithLabels(map[string]string{"env": "dev"})

	if sc.Labels["env"] != "dev" {
		t.Errorf("Labels[env] = %v, want dev", sc.Labels["env"])
	}
}

func TestScrapeConfig_FluentAPI(t *testing.T) {
	sc := NewScrapeConfig("api").
		WithInterval(Duration(30 * time.Second)).
		WithTimeout(Duration(10 * time.Second)).
		WithStaticTargets("api:8080")

	if sc.JobName != "api" {
		t.Errorf("JobName = %v, want api", sc.JobName)
	}
	if sc.ScrapeInterval != Duration(30*time.Second) {
		t.Errorf("ScrapeInterval = %v, want 30s", sc.ScrapeInterval)
	}
	if sc.ScrapeTimeout != Duration(10*time.Second) {
		t.Errorf("ScrapeTimeout = %v, want 10s", sc.ScrapeTimeout)
	}
	if len(sc.StaticConfigs) != 1 || len(sc.StaticConfigs[0].Targets) != 1 {
		t.Errorf("StaticConfigs not set correctly")
	}
}

func TestBasicAuth_MarshalYAML(t *testing.T) {
	config := &ScrapeConfig{
		JobName: "secure",
		BasicAuth: &BasicAuth{
			Username: "admin",
			Password: "secret",
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"basic_auth:",
		"username: admin",
		"password: secret",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestTLSConfig_MarshalYAML(t *testing.T) {
	config := &ScrapeConfig{
		JobName: "secure",
		TLSConfig: &TLSConfig{
			CAFile:             "/etc/prom/ca.crt",
			CertFile:           "/etc/prom/client.crt",
			KeyFile:            "/etc/prom/client.key",
			InsecureSkipVerify: false,
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"tls_config:",
		"ca_file: /etc/prom/ca.crt",
		"cert_file: /etc/prom/client.crt",
		"key_file: /etc/prom/client.key",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestRelabelConfig_MarshalYAML(t *testing.T) {
	config := &ScrapeConfig{
		JobName: "relabel-test",
		RelabelConfigs: []*RelabelConfig{
			{
				SourceLabels: []string{"__meta_kubernetes_pod_label_app"},
				TargetLabel:  "app",
				Action:       "replace",
			},
			{
				SourceLabels: []string{"__meta_kubernetes_pod_annotation_prometheus_io_scrape"},
				Regex:        "true",
				Action:       "keep",
			},
		},
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"relabel_configs:",
		"source_labels:",
		"__meta_kubernetes_pod_label_app",
		"target_label: app",
		"action: replace",
		"regex: \"true\"",
		"action: keep",
	}

	for _, exp := range expectations {
		if !contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() output missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}
