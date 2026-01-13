package grafana

import "testing"

func TestValueMapping(t *testing.T) {
	m := ValueMap("0", "Offline")
	if m.Type != MappingTypeValue {
		t.Errorf("Type = %q, want value", m.Type)
	}
	if m.Options.Match != "0" {
		t.Errorf("Options.Match = %q, want 0", m.Options.Match)
	}
	if m.Options.Result.Text != "Offline" {
		t.Errorf("Options.Result.Text = %q, want Offline", m.Options.Result.Text)
	}
}

func TestValueMapping_WithColor(t *testing.T) {
	m := ValueMap("1", "Online").WithColor("green")
	if m.Options.Result.Color != "green" {
		t.Errorf("Options.Result.Color = %q, want green", m.Options.Result.Color)
	}
}

func TestValueMapping_WithIndex(t *testing.T) {
	m := ValueMap("0", "Down").WithIndex(0)
	if m.Options.Result.Index != 0 {
		t.Errorf("Options.Result.Index = %d, want 0", m.Options.Result.Index)
	}
}

func TestRangeMapping(t *testing.T) {
	m := RangeMap(0, 50, "Low")
	if m.Type != MappingTypeRange {
		t.Errorf("Type = %q, want range", m.Type)
	}
	if *m.Options.From != 0 {
		t.Errorf("Options.From = %v, want 0", *m.Options.From)
	}
	if *m.Options.To != 50 {
		t.Errorf("Options.To = %v, want 50", *m.Options.To)
	}
	if m.Options.Result.Text != "Low" {
		t.Errorf("Options.Result.Text = %q, want Low", m.Options.Result.Text)
	}
}

func TestRangeMapping_WithColor(t *testing.T) {
	m := RangeMap(50, 100, "High").WithColor("red")
	if m.Options.Result.Color != "red" {
		t.Errorf("Options.Result.Color = %q, want red", m.Options.Result.Color)
	}
}

func TestRegexMapping(t *testing.T) {
	m := RegexMap("error.*", "Error")
	if m.Type != MappingTypeRegex {
		t.Errorf("Type = %q, want regex", m.Type)
	}
	if m.Options.Pattern != "error.*" {
		t.Errorf("Options.Pattern = %q, want error.*", m.Options.Pattern)
	}
	if m.Options.Result.Text != "Error" {
		t.Errorf("Options.Result.Text = %q, want Error", m.Options.Result.Text)
	}
}

func TestRegexMapping_WithColor(t *testing.T) {
	m := RegexMap("success.*", "Success").WithColor("green")
	if m.Options.Result.Color != "green" {
		t.Errorf("Options.Result.Color = %q, want green", m.Options.Result.Color)
	}
}

func TestSpecialMapping(t *testing.T) {
	tests := []struct {
		name     string
		mapping  *ValueMapping
		expected string
	}{
		{"NullMapping", NullMapping("N/A"), MappingSpecialNull},
		{"NaNMapping", NaNMapping("Invalid"), MappingSpecialNaN},
		{"TrueMapping", TrueMapping("Yes"), MappingSpecialTrue},
		{"FalseMapping", FalseMapping("No"), MappingSpecialFalse},
		{"EmptyMapping", EmptyMapping("-"), MappingSpecialEmpty},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mapping.Type != MappingTypeSpecial {
				t.Errorf("Type = %q, want special", tt.mapping.Type)
			}
			if tt.mapping.Options.Match != tt.expected {
				t.Errorf("Options.Match = %q, want %q", tt.mapping.Options.Match, tt.expected)
			}
		})
	}
}

func TestMappings(t *testing.T) {
	mappings := Mappings(
		ValueMap("0", "Off").WithColor("red"),
		ValueMap("1", "On").WithColor("green"),
	)

	if len(mappings) != 2 {
		t.Errorf("len(mappings) = %d, want 2", len(mappings))
	}
}

func TestStatusMappings(t *testing.T) {
	mappings := StatusMappings(
		"0", "Offline", "red",
		"1", "Online", "green",
		"2", "Maintenance", "yellow",
	)

	if len(mappings) != 3 {
		t.Errorf("len(mappings) = %d, want 3", len(mappings))
	}
	if mappings[0].Options.Match != "0" {
		t.Errorf("mappings[0].Options.Match = %q, want 0", mappings[0].Options.Match)
	}
	if mappings[0].Options.Result.Text != "Offline" {
		t.Errorf("mappings[0].Options.Result.Text = %q, want Offline", mappings[0].Options.Result.Text)
	}
	if mappings[0].Options.Result.Color != "red" {
		t.Errorf("mappings[0].Options.Result.Color = %q, want red", mappings[0].Options.Result.Color)
	}
}
