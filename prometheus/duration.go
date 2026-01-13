// Package prometheus provides types for Prometheus configuration synthesis.
package prometheus

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Duration wraps time.Duration with Prometheus-compatible serialization.
// It serializes to compact format like "30s", "1m", "5m30s".
type Duration time.Duration

// Common duration constants for convenience.
const (
	Second = Duration(time.Second)
	Minute = Duration(time.Minute)
	Hour   = Duration(time.Hour)
)

// String returns the Prometheus-compatible string representation.
// Examples: "30s", "1m", "5m30s", "1h30m"
func (d Duration) String() string {
	if d == 0 {
		return "0s"
	}

	dur := time.Duration(d)
	neg := dur < 0
	if neg {
		dur = -dur
	}

	var parts []string

	hours := dur / time.Hour
	dur -= hours * time.Hour
	if hours > 0 {
		parts = append(parts, fmt.Sprintf("%dh", hours))
	}

	minutes := dur / time.Minute
	dur -= minutes * time.Minute
	if minutes > 0 {
		parts = append(parts, fmt.Sprintf("%dm", minutes))
	}

	seconds := dur / time.Second
	dur -= seconds * time.Second
	if seconds > 0 {
		parts = append(parts, fmt.Sprintf("%ds", seconds))
	}

	millis := dur / time.Millisecond
	if millis > 0 {
		parts = append(parts, fmt.Sprintf("%dms", millis))
	}

	if len(parts) == 0 {
		return "0s"
	}

	result := strings.Join(parts, "")
	if neg {
		return "-" + result
	}
	return result
}

// MarshalYAML implements yaml.Marshaler for YAML serialization.
func (d Duration) MarshalYAML() (interface{}, error) {
	return d.String(), nil
}

// UnmarshalYAML implements yaml.Unmarshaler for YAML deserialization.
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	parsed, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// MarshalJSON implements json.Marshaler for JSON serialization.
func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte(`"` + d.String() + `"`), nil
}

// UnmarshalJSON implements json.Unmarshaler for JSON deserialization.
func (d *Duration) UnmarshalJSON(data []byte) error {
	s := strings.Trim(string(data), `"`)
	parsed, err := ParseDuration(s)
	if err != nil {
		return err
	}
	*d = parsed
	return nil
}

// durationComponentRE matches a single duration component (without anchoring to end).
var durationComponentRE = regexp.MustCompile(`^(\d+)(ms|s|m|h|d|w|y)`)

// ParseDuration parses a Prometheus-format duration string.
// Supported units: ms, s, m, h, d, w, y
// Compound durations like "5m30s" or "1h30m15s" are supported.
func ParseDuration(s string) (Duration, error) {
	if s == "" {
		return 0, fmt.Errorf("empty duration string")
	}

	orig := s
	neg := false
	if s[0] == '-' {
		neg = true
		s = s[1:]
	}

	if len(s) == 0 {
		return 0, fmt.Errorf("invalid duration: %q", orig)
	}

	// Handle compound durations like "5m30s" or "1h30m"
	var total time.Duration
	matched := false

	for len(s) > 0 {
		match := durationComponentRE.FindStringSubmatch(s)
		if match == nil {
			return 0, fmt.Errorf("invalid duration: %q", orig)
		}

		matched = true
		valStr := match[1]
		unit := match[2]

		val, err := strconv.ParseInt(valStr, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("invalid duration value: %q", valStr)
		}

		var d time.Duration
		switch unit {
		case "ms":
			d = time.Duration(val) * time.Millisecond
		case "s":
			d = time.Duration(val) * time.Second
		case "m":
			d = time.Duration(val) * time.Minute
		case "h":
			d = time.Duration(val) * time.Hour
		case "d":
			d = time.Duration(val) * 24 * time.Hour
		case "w":
			d = time.Duration(val) * 7 * 24 * time.Hour
		case "y":
			d = time.Duration(val) * 365 * 24 * time.Hour
		default:
			return 0, fmt.Errorf("unknown duration unit: %q", unit)
		}

		total += d
		s = s[len(match[0]):]
	}

	if !matched {
		return 0, fmt.Errorf("invalid duration: %q", orig)
	}

	if neg {
		total = -total
	}
	return Duration(total), nil
}
