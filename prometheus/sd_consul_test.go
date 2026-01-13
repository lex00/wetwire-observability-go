package prometheus

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestConsulSD_Serialize_Basic(t *testing.T) {
	sd := &ConsulSD{
		Server:     "consul.example.com:8500",
		Datacenter: "dc1",
		Services:   []string{"web", "api"},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"server: consul.example.com:8500",
		"datacenter: dc1",
		"services:",
		"- web",
		"- api",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestConsulSD_Serialize_WithTags(t *testing.T) {
	sd := &ConsulSD{
		Server:   "localhost:8500",
		Services: []string{"myservice"},
		Tags:     []string{"production", "v2"},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"tags:",
		"- production",
		"- v2",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestConsulSD_Serialize_WithNodeMeta(t *testing.T) {
	sd := &ConsulSD{
		Server: "localhost:8500",
		NodeMeta: map[string]string{
			"rack": "rack-1",
			"zone": "us-west-2a",
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "node_meta:") {
		t.Errorf("yaml.Marshal() missing node_meta\nGot:\n%s", yamlStr)
	}
}

func TestConsulSD_Serialize_WithTLS(t *testing.T) {
	sd := &ConsulSD{
		Server: "consul.example.com:8501",
		Scheme: "https",
		TLSConfig: &TLSConfig{
			CAFile:   "/etc/consul/ca.crt",
			CertFile: "/etc/consul/client.crt",
			KeyFile:  "/etc/consul/client.key",
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"scheme: https",
		"tls_config:",
		"ca_file:",
		"cert_file:",
		"key_file:",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestConsulSD_Serialize_WithBasicAuth(t *testing.T) {
	sd := &ConsulSD{
		Server: "localhost:8500",
		BasicAuth: &BasicAuth{
			Username: "admin",
			Password: "secret",
		},
	}

	data, err := yaml.Marshal(sd)
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
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestConsulSD_Serialize_WithToken(t *testing.T) {
	sd := &ConsulSD{
		Server: "localhost:8500",
		Token:  "my-consul-token",
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "token: my-consul-token") {
		t.Errorf("yaml.Marshal() missing token\nGot:\n%s", yamlStr)
	}
}

func TestConsulSD_Serialize_WithRefreshInterval(t *testing.T) {
	sd := &ConsulSD{
		Server:          "localhost:8500",
		RefreshInterval: Duration(30 * time.Second),
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "refresh_interval: 30s") {
		t.Errorf("yaml.Marshal() missing refresh_interval\nGot:\n%s", yamlStr)
	}
}

func TestNewConsulSD(t *testing.T) {
	sd := NewConsulSD()
	if sd == nil {
		t.Error("NewConsulSD() returned nil")
	}
}

func TestConsulSD_FluentAPI(t *testing.T) {
	sd := NewConsulSD().
		WithServer("consul.example.com:8500").
		WithToken("my-token").
		WithDatacenter("dc1").
		WithServices("web", "api").
		WithTags("production", "v2").
		WithNodeMeta(map[string]string{"rack": "rack-1"}).
		WithTagSeparator(";").
		WithRefreshInterval(Duration(60 * time.Second))

	if sd.Server != "consul.example.com:8500" {
		t.Errorf("Server = %v, want consul.example.com:8500", sd.Server)
	}
	if sd.Token != "my-token" {
		t.Errorf("Token = %v, want my-token", sd.Token)
	}
	if sd.Datacenter != "dc1" {
		t.Errorf("Datacenter = %v, want dc1", sd.Datacenter)
	}
	if len(sd.Services) != 2 {
		t.Errorf("len(Services) = %d, want 2", len(sd.Services))
	}
	if len(sd.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(sd.Tags))
	}
	if sd.NodeMeta["rack"] != "rack-1" {
		t.Errorf("NodeMeta[rack] = %v, want rack-1", sd.NodeMeta["rack"])
	}
	if sd.TagSeparator != ";" {
		t.Errorf("TagSeparator = %v, want ;", sd.TagSeparator)
	}
	if sd.RefreshInterval != Duration(60*time.Second) {
		t.Errorf("RefreshInterval = %v, want 60s", sd.RefreshInterval)
	}
}

func TestConsulSD_WithScheme(t *testing.T) {
	sd := NewConsulSD().WithScheme("https")
	if sd.Scheme != "https" {
		t.Errorf("Scheme = %v, want https", sd.Scheme)
	}
}

func TestConsulSD_WithNamespace(t *testing.T) {
	sd := NewConsulSD().WithNamespace("team-a")
	if sd.Namespace != "team-a" {
		t.Errorf("Namespace = %v, want team-a", sd.Namespace)
	}
}

func TestConsulSD_WithPartition(t *testing.T) {
	sd := NewConsulSD().WithPartition("partition-1")
	if sd.Partition != "partition-1" {
		t.Errorf("Partition = %v, want partition-1", sd.Partition)
	}
}

func TestConsulSD_WithAllowStale(t *testing.T) {
	sd := NewConsulSD().WithAllowStale(true)
	if !sd.AllowStale {
		t.Error("AllowStale not set to true")
	}
}

func TestConsulSD_WithTLSConfig(t *testing.T) {
	tls := &TLSConfig{CAFile: "/path/to/ca.crt"}
	sd := NewConsulSD().WithTLSConfig(tls)
	if sd.TLSConfig != tls {
		t.Error("TLSConfig not set correctly")
	}
}

func TestConsulSD_WithBasicAuth(t *testing.T) {
	auth := &BasicAuth{Username: "user", Password: "pass"}
	sd := NewConsulSD().WithBasicAuth(auth)
	if sd.BasicAuth != auth {
		t.Error("BasicAuth not set correctly")
	}
}

func TestConsulSD_WithProxyURL(t *testing.T) {
	sd := NewConsulSD().WithProxyURL("http://proxy.example.com:8080")
	if sd.ProxyURL != "http://proxy.example.com:8080" {
		t.Errorf("ProxyURL = %v, want http://proxy.example.com:8080", sd.ProxyURL)
	}
}

func TestScrapeConfig_WithConsulSD(t *testing.T) {
	sc := NewScrapeConfig("consul-services").
		WithConsulSD(NewConsulSD().
			WithServer("consul.example.com:8500").
			WithServices("web"))

	if len(sc.ConsulSDConfigs) != 1 {
		t.Errorf("len(ConsulSDConfigs) = %d, want 1", len(sc.ConsulSDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"job_name: consul-services",
		"consul_sd_configs:",
		"server: consul.example.com:8500",
		"services:",
		"- web",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestConsulSD_Unmarshal(t *testing.T) {
	input := `
server: consul.example.com:8500
token: my-acl-token
datacenter: dc1
services:
  - web
  - api
tags:
  - production
node_meta:
  rack: rack-1
tag_separator: ";"
refresh_interval: 30s
`
	var sd ConsulSD
	if err := yaml.Unmarshal([]byte(input), &sd); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if sd.Server != "consul.example.com:8500" {
		t.Errorf("Server = %v, want consul.example.com:8500", sd.Server)
	}
	if sd.Token != "my-acl-token" {
		t.Errorf("Token = %v, want my-acl-token", sd.Token)
	}
	if sd.Datacenter != "dc1" {
		t.Errorf("Datacenter = %v, want dc1", sd.Datacenter)
	}
	if len(sd.Services) != 2 {
		t.Errorf("len(Services) = %d, want 2", len(sd.Services))
	}
	if len(sd.Tags) != 1 {
		t.Errorf("len(Tags) = %d, want 1", len(sd.Tags))
	}
	if sd.NodeMeta["rack"] != "rack-1" {
		t.Errorf("NodeMeta[rack] = %v, want rack-1", sd.NodeMeta["rack"])
	}
	if sd.TagSeparator != ";" {
		t.Errorf("TagSeparator = %v, want ;", sd.TagSeparator)
	}
	if sd.RefreshInterval != Duration(30*time.Second) {
		t.Errorf("RefreshInterval = %v, want 30s", sd.RefreshInterval)
	}
}
