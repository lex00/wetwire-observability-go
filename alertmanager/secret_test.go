package alertmanager

import (
	"fmt"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestSecret_String(t *testing.T) {
	secret := NewSecret("my-secret-value")
	if secret.String() != "<secret>" {
		t.Errorf("String() = %q, want <secret>", secret.String())
	}
}

func TestSecret_String_Empty(t *testing.T) {
	secret := NewSecret("")
	if secret.String() != "" {
		t.Errorf("String() = %q, want empty string", secret.String())
	}
}

func TestSecret_NoLeakInFmt(t *testing.T) {
	secret := NewSecret("super-secret-key")

	// Using %s should not leak the value
	str := fmt.Sprintf("%s", secret)
	if str != "<secret>" {
		t.Errorf("fmt.Sprintf(%%s) = %q, want <secret>", str)
	}

	// Using %v should not leak the value
	str = fmt.Sprintf("%v", secret)
	if str != "<secret>" {
		t.Errorf("fmt.Sprintf(%%v) = %q, want <secret>", str)
	}
}

func TestSecret_MarshalYAML(t *testing.T) {
	secret := NewSecret("test-secret")
	data, err := yaml.Marshal(secret)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	// YAML should contain the actual value
	if string(data) != "test-secret\n" {
		t.Errorf("yaml.Marshal() = %q, want test-secret", string(data))
	}
}

func TestSecret_UnmarshalYAML(t *testing.T) {
	var secret Secret
	if err := yaml.Unmarshal([]byte("my-secret"), &secret); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if string(secret) != "my-secret" {
		t.Errorf("secret = %v, want my-secret", secret)
	}
}

func TestNewSecret(t *testing.T) {
	secret := NewSecret("value")
	if string(secret) != "value" {
		t.Errorf("NewSecret() = %v, want value", secret)
	}
}

func TestSecretFromEnv(t *testing.T) {
	secret := SecretFromEnv("PAGERDUTY_KEY")
	if string(secret) != "${PAGERDUTY_KEY}" {
		t.Errorf("SecretFromEnv() = %v, want ${PAGERDUTY_KEY}", secret)
	}
}

func TestSecretFile_String(t *testing.T) {
	sf := SecretFile("/etc/secrets/key")
	if sf.String() != "<secret-file>" {
		t.Errorf("String() = %q, want <secret-file>", sf.String())
	}
}

func TestSecretFile_String_Empty(t *testing.T) {
	sf := SecretFile("")
	if sf.String() != "" {
		t.Errorf("String() = %q, want empty string", sf.String())
	}
}

func TestSecret_InStruct(t *testing.T) {
	type Config struct {
		APIKey Secret `yaml:"api_key"`
		Name   string `yaml:"name"`
	}

	config := Config{
		APIKey: NewSecret("secret-api-key"),
		Name:   "test",
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if yamlStr != "api_key: secret-api-key\nname: test\n" {
		t.Errorf("yaml.Marshal() = %q", yamlStr)
	}

	// fmt.Sprintf should not leak
	str := fmt.Sprintf("Config: %+v", config)
	if str != "Config: {APIKey:<secret> Name:test}" {
		t.Errorf("fmt.Sprintf() = %q", str)
	}
}

func TestSecret_Comparison(t *testing.T) {
	s1 := NewSecret("test")
	s2 := NewSecret("test")
	s3 := NewSecret("different")

	if s1 != s2 {
		t.Error("Equal secrets should be equal")
	}
	if s1 == s3 {
		t.Error("Different secrets should not be equal")
	}
}
