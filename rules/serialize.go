package rules

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// Serialize converts the RulesFile to YAML bytes.
func (f *RulesFile) Serialize() ([]byte, error) {
	return yaml.Marshal(f)
}

// SerializeToFile writes the RulesFile to a YAML file.
func (f *RulesFile) SerializeToFile(path string) error {
	data, err := f.Serialize()
	if err != nil {
		return err
	}

	// Ensure the directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// MustSerialize converts the RulesFile to YAML bytes, panicking on error.
func (f *RulesFile) MustSerialize() []byte {
	data, err := f.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}

// Serialize converts a single RuleGroup to YAML bytes, wrapping it in a RulesFile.
func (g *RuleGroup) Serialize() ([]byte, error) {
	file := &RulesFile{Groups: []*RuleGroup{g}}
	return file.Serialize()
}

// SerializeToFile writes the RuleGroup to a YAML file.
func (g *RuleGroup) SerializeToFile(path string) error {
	file := &RulesFile{Groups: []*RuleGroup{g}}
	return file.SerializeToFile(path)
}

// MustSerialize converts the RuleGroup to YAML bytes, panicking on error.
func (g *RuleGroup) MustSerialize() []byte {
	data, err := g.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}
