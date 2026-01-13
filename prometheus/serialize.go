package prometheus

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Serialize converts a PrometheusConfig to YAML bytes.
// The output is a valid prometheus.yml configuration.
func (c *PrometheusConfig) Serialize() ([]byte, error) {
	return yaml.Marshal(c)
}

// SerializeToFile writes a PrometheusConfig to a file in YAML format.
// The file is created with 0644 permissions.
func (c *PrometheusConfig) SerializeToFile(path string) error {
	data, err := c.Serialize()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// MustSerialize converts a PrometheusConfig to YAML bytes.
// It panics if serialization fails.
func (c *PrometheusConfig) MustSerialize() []byte {
	data, err := c.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}

// SerializeGlobalConfig converts a GlobalConfig to YAML bytes.
func (c *GlobalConfig) Serialize() ([]byte, error) {
	return yaml.Marshal(c)
}

// SerializeScrapeConfig converts a ScrapeConfig to YAML bytes.
func (c *ScrapeConfig) Serialize() ([]byte, error) {
	return yaml.Marshal(c)
}
