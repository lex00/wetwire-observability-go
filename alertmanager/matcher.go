package alertmanager

import "fmt"

// MatchOp represents a matcher operator.
type MatchOp string

// Matcher operators.
const (
	MatchEqual    MatchOp = "="
	MatchNotEqual MatchOp = "!="
	MatchRegex    MatchOp = "=~"
	MatchNotRegex MatchOp = "!~"
)

// Matcher represents a label matcher in Alertmanager.
// Matchers are used in routes and inhibit rules to match alerts.
type Matcher struct {
	Label string  `yaml:"-"`
	Op    MatchOp `yaml:"-"`
	Value string  `yaml:"-"`
}

// MarshalYAML implements yaml.Marshaler.
// Matchers serialize to the format: label="value" or label=~"regex"
func (m *Matcher) MarshalYAML() (interface{}, error) {
	return m.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler.
func (m *Matcher) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}

	// Parse the matcher string
	parsed, err := ParseMatcher(s)
	if err != nil {
		return err
	}
	*m = *parsed
	return nil
}

// String returns the matcher in Alertmanager format.
func (m *Matcher) String() string {
	return fmt.Sprintf("%s%s%q", m.Label, m.Op, m.Value)
}

// ParseMatcher parses a matcher string in Alertmanager format.
// Supported formats: label="value", label!="value", label=~"regex", label!~"regex"
func ParseMatcher(s string) (*Matcher, error) {
	// Find the operator
	ops := []MatchOp{MatchNotRegex, MatchNotEqual, MatchRegex, MatchEqual}
	for _, op := range ops {
		opStr := string(op)
		for i := 0; i < len(s)-len(opStr); i++ {
			if s[i:i+len(opStr)] == opStr {
				label := s[:i]
				value := s[i+len(opStr):]
				// Remove quotes from value
				if len(value) >= 2 && value[0] == '"' && value[len(value)-1] == '"' {
					value = value[1 : len(value)-1]
				}
				return &Matcher{
					Label: label,
					Op:    op,
					Value: value,
				}, nil
			}
		}
	}
	return nil, fmt.Errorf("invalid matcher: %s", s)
}

// Match creates an equality matcher (label="value").
func Match(label, value string) *Matcher {
	return &Matcher{
		Label: label,
		Op:    MatchEqual,
		Value: value,
	}
}

// Eq is an alias for Match.
func Eq(label, value string) *Matcher {
	return Match(label, value)
}

// NotEq creates a not-equal matcher (label!="value").
func NotEq(label, value string) *Matcher {
	return &Matcher{
		Label: label,
		Op:    MatchNotEqual,
		Value: value,
	}
}

// Regex creates a regex matcher (label=~"regex").
func Regex(label, pattern string) *Matcher {
	return &Matcher{
		Label: label,
		Op:    MatchRegex,
		Value: pattern,
	}
}

// NotRegex creates a negative regex matcher (label!~"regex").
func NotRegex(label, pattern string) *Matcher {
	return &Matcher{
		Label: label,
		Op:    MatchNotRegex,
		Value: pattern,
	}
}

// Severity creates an equality matcher for the severity label.
func Severity(level string) *Matcher {
	return Eq("severity", level)
}

// SeverityCritical creates a matcher for critical severity.
func SeverityCritical() *Matcher {
	return Severity("critical")
}

// SeverityWarning creates a matcher for warning severity.
func SeverityWarning() *Matcher {
	return Severity("warning")
}

// SeverityInfo creates a matcher for info severity.
func SeverityInfo() *Matcher {
	return Severity("info")
}

// Team creates an equality matcher for the team label.
func Team(name string) *Matcher {
	return Eq("team", name)
}

// Service creates an equality matcher for the service label.
func Service(name string) *Matcher {
	return Eq("service", name)
}

// Environment creates an equality matcher for the env label.
func Environment(env string) *Matcher {
	return Eq("env", env)
}

// Alertname creates an equality matcher for the alertname label.
func Alertname(name string) *Matcher {
	return Eq("alertname", name)
}

// AlertnameRegex creates a regex matcher for the alertname label.
func AlertnameRegex(pattern string) *Matcher {
	return Regex("alertname", pattern)
}
