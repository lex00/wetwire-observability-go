package alertmanager

// Secret represents a secret value that should be handled securely.
// It protects sensitive data from being accidentally logged or exposed.
type Secret string

// String returns a redacted representation to prevent accidental logging.
func (s Secret) String() string {
	if s == "" {
		return ""
	}
	return "<secret>"
}

// MarshalYAML implements yaml.Marshaler.
// It outputs the actual value for serialization to YAML.
func (s Secret) MarshalYAML() (interface{}, error) {
	return string(s), nil
}

// NewSecret creates a new Secret from a string value.
func NewSecret(value string) Secret {
	return Secret(value)
}

// SecretFromEnv creates a Secret reference from an environment variable.
// The actual value should be loaded at runtime.
func SecretFromEnv(envVar string) Secret {
	return Secret("${" + envVar + "}")
}

// SecretFromFile creates a secret reference to a file path.
// Used when Alertmanager should read the secret from a file.
type SecretFile string

// String returns a redacted representation.
func (s SecretFile) String() string {
	if s == "" {
		return ""
	}
	return "<secret-file>"
}
