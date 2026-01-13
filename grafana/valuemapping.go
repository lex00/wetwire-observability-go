package grafana

// Mapping type constants.
const (
	MappingTypeValue   = "value"
	MappingTypeRange   = "range"
	MappingTypeRegex   = "regex"
	MappingTypeSpecial = "special"
)

// Special mapping match constants.
const (
	MappingSpecialNull   = "null"
	MappingSpecialNaN    = "nan"
	MappingSpecialTrue   = "true"
	MappingSpecialFalse  = "false"
	MappingSpecialEmpty  = "empty"
)

// ValueMapping represents a Grafana value mapping.
type ValueMapping struct {
	Type    string               `json:"type"`
	Options ValueMappingOptions `json:"options"`
}

// ValueMappingOptions contains mapping options.
type ValueMappingOptions struct {
	// Match is used for value and special mappings
	Match string `json:"match,omitempty"`

	// From and To are used for range mappings
	From *float64 `json:"from,omitempty"`
	To   *float64 `json:"to,omitempty"`

	// Pattern is used for regex mappings
	Pattern string `json:"pattern,omitempty"`

	// Result is the mapping result
	Result MappingResult `json:"result"`
}

// MappingResult contains the result of a mapping.
type MappingResult struct {
	Text  string `json:"text"`
	Color string `json:"color,omitempty"`
	Index int    `json:"index,omitempty"`
}

// ValueMap creates a value mapping (exact match).
func ValueMap(value, text string) *ValueMapping {
	return &ValueMapping{
		Type: MappingTypeValue,
		Options: ValueMappingOptions{
			Match: value,
			Result: MappingResult{
				Text: text,
			},
		},
	}
}

// WithColor sets the result color.
func (m *ValueMapping) WithColor(color string) *ValueMapping {
	m.Options.Result.Color = color
	return m
}

// WithIndex sets the result index.
func (m *ValueMapping) WithIndex(index int) *ValueMapping {
	m.Options.Result.Index = index
	return m
}

// RangeMap creates a range mapping (from-to).
func RangeMap(from, to float64, text string) *ValueMapping {
	return &ValueMapping{
		Type: MappingTypeRange,
		Options: ValueMappingOptions{
			From: &from,
			To:   &to,
			Result: MappingResult{
				Text: text,
			},
		},
	}
}

// RegexMap creates a regex mapping.
func RegexMap(pattern, text string) *ValueMapping {
	return &ValueMapping{
		Type: MappingTypeRegex,
		Options: ValueMappingOptions{
			Pattern: pattern,
			Result: MappingResult{
				Text: text,
			},
		},
	}
}

// specialMapping creates a special mapping.
func specialMapping(match, text string) *ValueMapping {
	return &ValueMapping{
		Type: MappingTypeSpecial,
		Options: ValueMappingOptions{
			Match: match,
			Result: MappingResult{
				Text: text,
			},
		},
	}
}

// NullMapping creates a mapping for null values.
func NullMapping(text string) *ValueMapping {
	return specialMapping(MappingSpecialNull, text)
}

// NaNMapping creates a mapping for NaN values.
func NaNMapping(text string) *ValueMapping {
	return specialMapping(MappingSpecialNaN, text)
}

// TrueMapping creates a mapping for true values.
func TrueMapping(text string) *ValueMapping {
	return specialMapping(MappingSpecialTrue, text)
}

// FalseMapping creates a mapping for false values.
func FalseMapping(text string) *ValueMapping {
	return specialMapping(MappingSpecialFalse, text)
}

// EmptyMapping creates a mapping for empty values.
func EmptyMapping(text string) *ValueMapping {
	return specialMapping(MappingSpecialEmpty, text)
}

// Mappings creates a slice of value mappings.
func Mappings(mappings ...*ValueMapping) []*ValueMapping {
	return mappings
}

// StatusMappings creates common status value mappings.
// Arguments should be triplets of (value, text, color).
func StatusMappings(valueTextColor ...string) []*ValueMapping {
	if len(valueTextColor)%3 != 0 {
		return nil
	}

	var mappings []*ValueMapping
	for i := 0; i < len(valueTextColor); i += 3 {
		mappings = append(mappings, ValueMap(valueTextColor[i], valueTextColor[i+1]).WithColor(valueTextColor[i+2]))
	}
	return mappings
}

// BooleanMappings creates mappings for boolean values.
func BooleanMappings(trueText, falseText string, trueColor, falseColor string) []*ValueMapping {
	return []*ValueMapping{
		TrueMapping(trueText).WithColor(trueColor),
		FalseMapping(falseText).WithColor(falseColor),
	}
}
