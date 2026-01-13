package prometheus

import (
	"strings"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestEC2SD_Serialize_Basic(t *testing.T) {
	sd := &EC2SD{
		Region: "us-west-2",
		Port:   9100,
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"region: us-west-2",
		"port: 9100",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestEC2SD_Serialize_WithCredentials(t *testing.T) {
	sd := &EC2SD{
		Region:    "us-east-1",
		AccessKey: "AKIAIOSFODNN7EXAMPLE",
		SecretKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"access_key: AKIAIOSFODNN7EXAMPLE",
		"secret_key: wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestEC2SD_Serialize_WithProfile(t *testing.T) {
	sd := &EC2SD{
		Region:  "us-west-2",
		Profile: "production",
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "profile: production") {
		t.Errorf("yaml.Marshal() missing profile\nGot:\n%s", yamlStr)
	}
}

func TestEC2SD_Serialize_WithRoleARN(t *testing.T) {
	sd := &EC2SD{
		Region:  "us-west-2",
		RoleARN: "arn:aws:iam::123456789012:role/prometheus-discovery",
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "role_arn: arn:aws:iam::123456789012:role/prometheus-discovery") {
		t.Errorf("yaml.Marshal() missing role_arn\nGot:\n%s", yamlStr)
	}
}

func TestEC2SD_Serialize_WithFilters(t *testing.T) {
	sd := &EC2SD{
		Region: "us-west-2",
		Filters: []EC2Filter{
			{Name: "tag:Environment", Values: []string{"production"}},
			{Name: "instance-state-name", Values: []string{"running"}},
		},
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"filters:",
		"name: tag:Environment",
		"values:",
		"- production",
		"name: instance-state-name",
		"- running",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestEC2SD_Serialize_WithRefreshInterval(t *testing.T) {
	sd := &EC2SD{
		Region:          "us-west-2",
		RefreshInterval: Duration(60 * time.Second),
	}

	data, err := yaml.Marshal(sd)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "refresh_interval: 1m") {
		t.Errorf("yaml.Marshal() missing refresh_interval\nGot:\n%s", yamlStr)
	}
}

func TestNewEC2SD(t *testing.T) {
	sd := NewEC2SD()
	if sd == nil {
		t.Error("NewEC2SD() returned nil")
	}
}

func TestEC2SD_FluentAPI(t *testing.T) {
	sd := NewEC2SD().
		WithRegion("us-west-2").
		WithCredentials("access", "secret").
		WithProfile("prod").
		WithRoleARN("arn:aws:iam::123456789012:role/test").
		WithPort(9100).
		WithRefreshInterval(Duration(2 * time.Minute))

	if sd.Region != "us-west-2" {
		t.Errorf("Region = %v, want us-west-2", sd.Region)
	}
	if sd.AccessKey != "access" {
		t.Errorf("AccessKey = %v, want access", sd.AccessKey)
	}
	if sd.SecretKey != "secret" {
		t.Errorf("SecretKey = %v, want secret", sd.SecretKey)
	}
	if sd.Profile != "prod" {
		t.Errorf("Profile = %v, want prod", sd.Profile)
	}
	if sd.RoleARN != "arn:aws:iam::123456789012:role/test" {
		t.Errorf("RoleARN = %v, want arn:aws:iam::123456789012:role/test", sd.RoleARN)
	}
	if sd.Port != 9100 {
		t.Errorf("Port = %v, want 9100", sd.Port)
	}
	if sd.RefreshInterval != Duration(2*time.Minute) {
		t.Errorf("RefreshInterval = %v, want 2m", sd.RefreshInterval)
	}
}

func TestEC2SD_WithAccessKey(t *testing.T) {
	sd := NewEC2SD().WithAccessKey("AKIAEXAMPLE")
	if sd.AccessKey != "AKIAEXAMPLE" {
		t.Errorf("AccessKey = %v, want AKIAEXAMPLE", sd.AccessKey)
	}
}

func TestEC2SD_WithSecretKey(t *testing.T) {
	sd := NewEC2SD().WithSecretKey("mysecret")
	if sd.SecretKey != "mysecret" {
		t.Errorf("SecretKey = %v, want mysecret", sd.SecretKey)
	}
}

func TestEC2SD_WithEndpoint(t *testing.T) {
	sd := NewEC2SD().WithEndpoint("http://localhost:4566")
	if sd.Endpoint != "http://localhost:4566" {
		t.Errorf("Endpoint = %v, want http://localhost:4566", sd.Endpoint)
	}
}

func TestEC2SD_WithFilters(t *testing.T) {
	filters := []EC2Filter{
		{Name: "tag:app", Values: []string{"web"}},
	}
	sd := NewEC2SD().WithFilters(filters...)
	if len(sd.Filters) != 1 {
		t.Errorf("len(Filters) = %d, want 1", len(sd.Filters))
	}
}

func TestEC2SD_WithFilter(t *testing.T) {
	sd := NewEC2SD().
		WithFilter("tag:Environment", "production", "staging").
		WithFilter("instance-state-name", "running")

	if len(sd.Filters) != 2 {
		t.Errorf("len(Filters) = %d, want 2", len(sd.Filters))
	}
	if sd.Filters[0].Name != "tag:Environment" {
		t.Errorf("Filters[0].Name = %v, want tag:Environment", sd.Filters[0].Name)
	}
	if len(sd.Filters[0].Values) != 2 {
		t.Errorf("len(Filters[0].Values) = %d, want 2", len(sd.Filters[0].Values))
	}
}

func TestEC2SD_WithTagFilter(t *testing.T) {
	sd := NewEC2SD().WithTagFilter("Environment", "production")

	if len(sd.Filters) != 1 {
		t.Errorf("len(Filters) = %d, want 1", len(sd.Filters))
	}
	if sd.Filters[0].Name != "tag:Environment" {
		t.Errorf("Filters[0].Name = %v, want tag:Environment", sd.Filters[0].Name)
	}
}

func TestScrapeConfig_WithEC2SD(t *testing.T) {
	sc := NewScrapeConfig("ec2-instances").
		WithEC2SD(NewEC2SD().
			WithRegion("us-west-2").
			WithPort(9100).
			WithTagFilter("Environment", "production"))

	if len(sc.EC2SDConfigs) != 1 {
		t.Errorf("len(EC2SDConfigs) = %d, want 1", len(sc.EC2SDConfigs))
	}

	data, err := yaml.Marshal(sc)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"job_name: ec2-instances",
		"ec2_sd_configs:",
		"region: us-west-2",
		"port: 9100",
		"filters:",
		"name: tag:Environment",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestEC2SD_Unmarshal(t *testing.T) {
	input := `
region: us-west-2
access_key: AKIAEXAMPLE
secret_key: secretexample
profile: production
role_arn: arn:aws:iam::123456789012:role/prometheus
port: 9100
refresh_interval: 2m
filters:
  - name: tag:Environment
    values:
      - production
  - name: instance-state-name
    values:
      - running
`
	var sd EC2SD
	if err := yaml.Unmarshal([]byte(input), &sd); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if sd.Region != "us-west-2" {
		t.Errorf("Region = %v, want us-west-2", sd.Region)
	}
	if sd.AccessKey != "AKIAEXAMPLE" {
		t.Errorf("AccessKey = %v, want AKIAEXAMPLE", sd.AccessKey)
	}
	if sd.SecretKey != "secretexample" {
		t.Errorf("SecretKey = %v, want secretexample", sd.SecretKey)
	}
	if sd.Profile != "production" {
		t.Errorf("Profile = %v, want production", sd.Profile)
	}
	if sd.RoleARN != "arn:aws:iam::123456789012:role/prometheus" {
		t.Errorf("RoleARN = %v, want arn:aws:iam::123456789012:role/prometheus", sd.RoleARN)
	}
	if sd.Port != 9100 {
		t.Errorf("Port = %v, want 9100", sd.Port)
	}
	if sd.RefreshInterval != Duration(2*time.Minute) {
		t.Errorf("RefreshInterval = %v, want 2m", sd.RefreshInterval)
	}
	if len(sd.Filters) != 2 {
		t.Errorf("len(Filters) = %d, want 2", len(sd.Filters))
	}
	if sd.Filters[0].Name != "tag:Environment" {
		t.Errorf("Filters[0].Name = %v, want tag:Environment", sd.Filters[0].Name)
	}
}
