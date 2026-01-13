package prometheus

import (
	"encoding/json"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestDuration_String(t *testing.T) {
	tests := []struct {
		name     string
		duration Duration
		want     string
	}{
		{"zero", 0, "0s"},
		{"one second", Duration(time.Second), "1s"},
		{"thirty seconds", Duration(30 * time.Second), "30s"},
		{"one minute", Duration(time.Minute), "1m"},
		{"five minutes", Duration(5 * time.Minute), "5m"},
		{"one hour", Duration(time.Hour), "1h"},
		{"five minutes thirty seconds", Duration(5*time.Minute + 30*time.Second), "5m30s"},
		{"one hour thirty minutes", Duration(time.Hour + 30*time.Minute), "1h30m"},
		{"complex duration", Duration(2*time.Hour + 15*time.Minute + 30*time.Second), "2h15m30s"},
		{"milliseconds", Duration(500 * time.Millisecond), "500ms"},
		{"seconds and milliseconds", Duration(time.Second + 500*time.Millisecond), "1s500ms"},
		{"negative duration", Duration(-30 * time.Second), "-30s"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.duration.String()
			if got != tt.want {
				t.Errorf("Duration.String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestParseDuration(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    Duration
		wantErr bool
	}{
		{"seconds", "30s", Duration(30 * time.Second), false},
		{"minutes", "5m", Duration(5 * time.Minute), false},
		{"hours", "1h", Duration(time.Hour), false},
		{"days", "1d", Duration(24 * time.Hour), false},
		{"weeks", "1w", Duration(7 * 24 * time.Hour), false},
		{"years", "1y", Duration(365 * 24 * time.Hour), false},
		{"milliseconds", "500ms", Duration(500 * time.Millisecond), false},
		{"compound", "5m30s", Duration(5*time.Minute + 30*time.Second), false},
		{"complex", "1h30m15s", Duration(time.Hour + 30*time.Minute + 15*time.Second), false},
		{"negative", "-30s", Duration(-30 * time.Second), false},
		{"empty", "", 0, true},
		{"invalid", "abc", 0, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseDuration(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("ParseDuration() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && got != tt.want {
				t.Errorf("ParseDuration() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDuration_MarshalYAML(t *testing.T) {
	type testStruct struct {
		Interval Duration `yaml:"interval"`
	}

	ts := testStruct{Interval: Duration(30 * time.Second)}
	data, err := yaml.Marshal(&ts)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	expected := "interval: 30s\n"
	if string(data) != expected {
		t.Errorf("yaml.Marshal() = %q, want %q", string(data), expected)
	}
}

func TestDuration_UnmarshalYAML(t *testing.T) {
	type testStruct struct {
		Interval Duration `yaml:"interval"`
	}

	input := "interval: 5m30s\n"
	var ts testStruct
	if err := yaml.Unmarshal([]byte(input), &ts); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	want := Duration(5*time.Minute + 30*time.Second)
	if ts.Interval != want {
		t.Errorf("yaml.Unmarshal() Interval = %v, want %v", ts.Interval, want)
	}
}

func TestDuration_MarshalJSON(t *testing.T) {
	type testStruct struct {
		Interval Duration `json:"interval"`
	}

	ts := testStruct{Interval: Duration(30 * time.Second)}
	data, err := json.Marshal(&ts)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}

	expected := `{"interval":"30s"}`
	if string(data) != expected {
		t.Errorf("json.Marshal() = %q, want %q", string(data), expected)
	}
}

func TestDuration_UnmarshalJSON(t *testing.T) {
	type testStruct struct {
		Interval Duration `json:"interval"`
	}

	input := `{"interval":"5m30s"}`
	var ts testStruct
	if err := json.Unmarshal([]byte(input), &ts); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}

	want := Duration(5*time.Minute + 30*time.Second)
	if ts.Interval != want {
		t.Errorf("json.Unmarshal() Interval = %v, want %v", ts.Interval, want)
	}
}

func TestDuration_Constants(t *testing.T) {
	if Second != Duration(time.Second) {
		t.Errorf("Second = %v, want %v", Second, Duration(time.Second))
	}
	if Minute != Duration(time.Minute) {
		t.Errorf("Minute = %v, want %v", Minute, Duration(time.Minute))
	}
	if Hour != Duration(time.Hour) {
		t.Errorf("Hour = %v, want %v", Hour, Duration(time.Hour))
	}
}
