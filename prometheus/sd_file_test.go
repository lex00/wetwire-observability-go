package prometheus

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestFileSD_Serialize_Basic(t *testing.T) {
	sd := &FileSD{
		Files: []string{"/etc/prometheus/targets/*.json"},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"files:",
		"- /etc/prometheus/targets/*.json",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestFileSD_Serialize_MultipleFiles(t *testing.T) {
	sd := &FileSD{
		Files: []string{
			"/etc/prometheus/targets/web.json",
			"/etc/prometheus/targets/api.json",
			"/etc/prometheus/targets/*.yaml",
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"files:",
		"- /etc/prometheus/targets/web.json",
		"- /etc/prometheus/targets/api.json",
		"- /etc/prometheus/targets/*.yaml",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestFileSD_Serialize_WithRefreshInterval(t *testing.T) {
	sd := &FileSD{
		Files:           []string{"/etc/prometheus/targets/*.json"},
		RefreshInterval: Duration(5 * time.Minute),
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "refresh_interval: 5m") {
		t.Errorf("yaml.Marshal() missing refresh_interval\nGot:\n%s", yamlStr)
	}
}

func TestNewFileSD(t *testing.T) {
	sd := NewFileSD()
	if sd == nil {
		t.Error("NewFileSD() returned nil")
	}
}

func TestFileSD_FluentAPI(t *testing.T) {
	sd := NewFileSD().
		WithFiles("/etc/prometheus/targets/*.json", "/etc/prometheus/targets/*.yaml").
		WithRefreshInterval(Duration(10 * time.Minute))

	if len(sd.Files) != 2 {
		t.Errorf("len(Files) = %d, want 2", len(sd.Files))
	}
	if sd.Files[0] != "/etc/prometheus/targets/*.json" {
		t.Errorf("Files[0] = %v, want /etc/prometheus/targets/*.json", sd.Files[0])
	}
	if sd.Files[1] != "/etc/prometheus/targets/*.yaml" {
		t.Errorf("Files[1] = %v, want /etc/prometheus/targets/*.yaml", sd.Files[1])
	}
	if sd.RefreshInterval != Duration(10*time.Minute) {
		t.Errorf("RefreshInterval = %v, want 10m", sd.RefreshInterval)
	}
}

func TestFileSD_WithFiles(t *testing.T) {
	sd := NewFileSD().WithFiles("/path/to/targets.json")
	if len(sd.Files) != 1 {
		t.Errorf("len(Files) = %d, want 1", len(sd.Files))
	}
	if sd.Files[0] != "/path/to/targets.json" {
		t.Errorf("Files[0] = %v, want /path/to/targets.json", sd.Files[0])
	}
}

func TestFileSD_WithRefreshInterval(t *testing.T) {
	sd := NewFileSD().WithRefreshInterval(Duration(30 * time.Second))
	if sd.RefreshInterval != Duration(30*time.Second) {
		t.Errorf("RefreshInterval = %v, want 30s", sd.RefreshInterval)
	}
}

func TestScrapeConfig_WithFileSD(t *testing.T) {
	sc := NewScrapeConfig("file-targets").
		WithFileSD(NewFileSD().
			WithFiles("/etc/prometheus/targets/*.json").
			WithRefreshInterval(Duration(5 * time.Minute)))

	if len(sc.FileSDConfigs) != 1 {
		t.Errorf("len(FileSDConfigs) = %d, want 1", len(sc.FileSDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"job_name: file-targets",
		"file_sd_configs:",
		"files:",
		"- /etc/prometheus/targets/*.json",
		"refresh_interval: 5m",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestFileSD_Unmarshal(t *testing.T) {
	input := `
files:
  - /etc/prometheus/targets/*.json
  - /etc/prometheus/targets/*.yaml
refresh_interval: 5m
`
	var sd FileSD
	if err := yaml.Unmarshal([]byte(input), &sd); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(sd.Files) != 2 {
		t.Errorf("len(Files) = %d, want 2", len(sd.Files))
	}
	if sd.Files[0] != "/etc/prometheus/targets/*.json" {
		t.Errorf("Files[0] = %v, want /etc/prometheus/targets/*.json", sd.Files[0])
	}
	if sd.Files[1] != "/etc/prometheus/targets/*.yaml" {
		t.Errorf("Files[1] = %v, want /etc/prometheus/targets/*.yaml", sd.Files[1])
	}
	if sd.RefreshInterval != Duration(5*time.Minute) {
		t.Errorf("RefreshInterval = %v, want 5m", sd.RefreshInterval)
	}
}

func TestScrapeConfig_MultipleFileSD(t *testing.T) {
	sc := NewScrapeConfig("multi-file").
		WithFileSD(NewFileSD().WithFiles("/targets/web/*.json")).
		WithFileSD(NewFileSD().WithFiles("/targets/api/*.json"))

	if len(sc.FileSDConfigs) != 2 {
		t.Errorf("len(FileSDConfigs) = %d, want 2", len(sc.FileSDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "file_sd_configs:") {
		t.Errorf("yaml.Marshal() missing file_sd_configs\nGot:\n%s", yamlStr)
	}
}
