package alertmanager

import (
	"testing"

	"gopkg.in/yaml.v3"
)

func TestMatch(t *testing.T) {
	m := Match("severity", "critical")
	if m.Label != "severity" {
		t.Errorf("Label = %v, want severity", m.Label)
	}
	if m.Op != MatchEqual {
		t.Errorf("Op = %v, want =", m.Op)
	}
	if m.Value != "critical" {
		t.Errorf("Value = %v, want critical", m.Value)
	}
}

func TestEq(t *testing.T) {
	m := Eq("team", "backend")
	if m.Label != "team" || m.Op != MatchEqual || m.Value != "backend" {
		t.Errorf("Eq() = %v, want team=\"backend\"", m)
	}
}

func TestNotEq(t *testing.T) {
	m := NotEq("env", "test")
	if m.Label != "env" || m.Op != MatchNotEqual || m.Value != "test" {
		t.Errorf("NotEq() = %v, want env!=\"test\"", m)
	}
}

func TestRegex(t *testing.T) {
	m := Regex("alertname", "High.*")
	if m.Label != "alertname" || m.Op != MatchRegex || m.Value != "High.*" {
		t.Errorf("Regex() = %v, want alertname=~\"High.*\"", m)
	}
}

func TestNotRegex(t *testing.T) {
	m := NotRegex("service", "test-.*")
	if m.Label != "service" || m.Op != MatchNotRegex || m.Value != "test-.*" {
		t.Errorf("NotRegex() = %v, want service!~\"test-.*\"", m)
	}
}

func TestMatcher_String(t *testing.T) {
	tests := []struct {
		matcher *Matcher
		want    string
	}{
		{Match("severity", "critical"), `severity="critical"`},
		{NotEq("env", "test"), `env!="test"`},
		{Regex("alertname", "High.*"), `alertname=~"High.*"`},
		{NotRegex("service", "test-.*"), `service!~"test-.*"`},
	}

	for _, tt := range tests {
		got := tt.matcher.String()
		if got != tt.want {
			t.Errorf("String() = %q, want %q", got, tt.want)
		}
	}
}

func TestMatcher_MarshalYAML(t *testing.T) {
	m := Match("severity", "critical")
	data, err := yaml.Marshal(m)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	want := `severity="critical"`
	got := string(data)
	// YAML adds a newline
	if got != want+"\n" {
		t.Errorf("yaml.Marshal() = %q, want %q", got, want)
	}
}

func TestMatcher_UnmarshalYAML(t *testing.T) {
	tests := []struct {
		input string
		want  *Matcher
	}{
		{`severity="critical"`, Match("severity", "critical")},
		{`env!="test"`, NotEq("env", "test")},
		{`alertname=~"High.*"`, Regex("alertname", "High.*")},
		{`service!~"test-.*"`, NotRegex("service", "test-.*")},
	}

	for _, tt := range tests {
		var m Matcher
		if err := yaml.Unmarshal([]byte(tt.input), &m); err != nil {
			t.Errorf("yaml.Unmarshal(%q) error = %v", tt.input, err)
			continue
		}
		if m.Label != tt.want.Label || m.Op != tt.want.Op || m.Value != tt.want.Value {
			t.Errorf("yaml.Unmarshal(%q) = %v, want %v", tt.input, m, tt.want)
		}
	}
}

func TestParseMatcher(t *testing.T) {
	tests := []struct {
		input string
		want  *Matcher
	}{
		{`severity="critical"`, Match("severity", "critical")},
		{`env!="test"`, NotEq("env", "test")},
		{`alertname=~"High.*"`, Regex("alertname", "High.*")},
		{`service!~"test-.*"`, NotRegex("service", "test-.*")},
	}

	for _, tt := range tests {
		got, err := ParseMatcher(tt.input)
		if err != nil {
			t.Errorf("ParseMatcher(%q) error = %v", tt.input, err)
			continue
		}
		if got.Label != tt.want.Label || got.Op != tt.want.Op || got.Value != tt.want.Value {
			t.Errorf("ParseMatcher(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}

func TestParseMatcher_Invalid(t *testing.T) {
	inputs := []string{
		"invalid",
		"no-operator",
		"",
	}

	for _, input := range inputs {
		_, err := ParseMatcher(input)
		if err == nil {
			t.Errorf("ParseMatcher(%q) should return error", input)
		}
	}
}

func TestSeverityHelpers(t *testing.T) {
	tests := []struct {
		name    string
		matcher *Matcher
		want    string
	}{
		{"critical", SeverityCritical(), "critical"},
		{"warning", SeverityWarning(), "warning"},
		{"info", SeverityInfo(), "info"},
	}

	for _, tt := range tests {
		if tt.matcher.Label != "severity" {
			t.Errorf("%s: Label = %v, want severity", tt.name, tt.matcher.Label)
		}
		if tt.matcher.Value != tt.want {
			t.Errorf("%s: Value = %v, want %v", tt.name, tt.matcher.Value, tt.want)
		}
	}
}

func TestLabelHelpers(t *testing.T) {
	tests := []struct {
		name    string
		matcher *Matcher
		label   string
		value   string
	}{
		{"Team", Team("backend"), "team", "backend"},
		{"Service", Service("api"), "service", "api"},
		{"Environment", Environment("production"), "env", "production"},
		{"Alertname", Alertname("HighCPU"), "alertname", "HighCPU"},
	}

	for _, tt := range tests {
		if tt.matcher.Label != tt.label {
			t.Errorf("%s: Label = %v, want %v", tt.name, tt.matcher.Label, tt.label)
		}
		if tt.matcher.Value != tt.value {
			t.Errorf("%s: Value = %v, want %v", tt.name, tt.matcher.Value, tt.value)
		}
	}
}

func TestAlertnameRegex(t *testing.T) {
	m := AlertnameRegex("High.*")
	if m.Label != "alertname" {
		t.Errorf("Label = %v, want alertname", m.Label)
	}
	if m.Op != MatchRegex {
		t.Errorf("Op = %v, want =~", m.Op)
	}
	if m.Value != "High.*" {
		t.Errorf("Value = %v, want High.*", m.Value)
	}
}

func TestMatcherInSlice(t *testing.T) {
	matchers := []*Matcher{
		SeverityCritical(),
		Team("backend"),
		Environment("production"),
	}

	data, err := yaml.Marshal(matchers)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	expected := `- severity="critical"
- team="backend"
- env="production"
`
	if string(data) != expected {
		t.Errorf("yaml.Marshal() = %q, want %q", string(data), expected)
	}

	// Unmarshal back
	var restored []*Matcher
	if err := yaml.Unmarshal(data, &restored); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if len(restored) != 3 {
		t.Errorf("len(restored) = %d, want 3", len(restored))
	}

	for i, m := range restored {
		if m.Label != matchers[i].Label || m.Value != matchers[i].Value {
			t.Errorf("restored[%d] = %v, want %v", i, m, matchers[i])
		}
	}
}
