package alertmanager

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Serialize converts an AlertmanagerConfig to YAML bytes.
// The output is a valid alertmanager.yml configuration.
func (c *AlertmanagerConfig) Serialize() ([]byte, error) {
	return yaml.Marshal(c)
}

// SerializeToFile writes an AlertmanagerConfig to a file in YAML format.
// The file is created with 0644 permissions.
func (c *AlertmanagerConfig) SerializeToFile(path string) error {
	data, err := c.Serialize()
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}

// MustSerialize converts an AlertmanagerConfig to YAML bytes.
// It panics if serialization fails.
func (c *AlertmanagerConfig) MustSerialize() []byte {
	data, err := c.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}

// Serialize converts a GlobalConfig to YAML bytes.
func (c *GlobalConfig) Serialize() ([]byte, error) {
	return yaml.Marshal(c)
}

// Serialize converts a Route to YAML bytes.
func (r *Route) Serialize() ([]byte, error) {
	return yaml.Marshal(r)
}

// Serialize converts a Receiver to YAML bytes.
func (r *Receiver) Serialize() ([]byte, error) {
	return yaml.Marshal(r)
}

// Serialize converts an InhibitRule to YAML bytes.
func (i *InhibitRule) Serialize() ([]byte, error) {
	return yaml.Marshal(i)
}

// Serialize converts a MuteTimeInterval to YAML bytes.
func (m *MuteTimeInterval) Serialize() ([]byte, error) {
	return yaml.Marshal(m)
}
