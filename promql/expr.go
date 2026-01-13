// Package promql provides typed PromQL expression builders.
package promql

import (
	"fmt"
	"strconv"
	"strings"
)

// Expr is the interface for all PromQL expressions.
type Expr interface {
	String() string
}

// Raw represents a raw PromQL expression string.
// Use this for complex expressions that aren't easily built with typed helpers.
type Raw string

// String returns the raw expression string.
func (r Raw) String() string {
	return string(r)
}

// ScalarExpr represents a scalar value.
type ScalarExpr struct {
	value float64
}

// String returns the scalar value as a string.
func (s *ScalarExpr) String() string {
	// Use strconv.FormatFloat to avoid trailing zeros
	return strconv.FormatFloat(s.value, 'f', -1, 64)
}

// Scalar creates a scalar expression.
func Scalar(v float64) *ScalarExpr {
	return &ScalarExpr{value: v}
}

// VectorExpr represents an instant vector selector.
type VectorExpr struct {
	metric   string
	matchers []LabelMatcher
	offset   string
}

// String returns the vector selector as a PromQL string.
func (v *VectorExpr) String() string {
	var sb strings.Builder
	sb.WriteString(v.metric)

	if len(v.matchers) > 0 {
		sb.WriteByte('{')
		for i, m := range v.matchers {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(m.String())
		}
		sb.WriteByte('}')
	}

	if v.offset != "" {
		sb.WriteString(" offset ")
		sb.WriteString(v.offset)
	}

	return sb.String()
}

// WithOffset adds an offset modifier.
func (v *VectorExpr) WithOffset(offset string) *VectorExpr {
	v.offset = offset
	return v
}

// Metric creates a simple metric selector without labels.
func Metric(name string) *VectorExpr {
	return &VectorExpr{metric: name}
}

// Vector creates a vector selector with optional label matchers.
func Vector(metric string, matchers ...LabelMatcher) *VectorExpr {
	return &VectorExpr{
		metric:   metric,
		matchers: matchers,
	}
}

// RangeVectorExpr represents a range vector selector.
type RangeVectorExpr struct {
	metric   string
	matchers []LabelMatcher
	duration string
	offset   string
}

// String returns the range vector selector as a PromQL string.
func (r *RangeVectorExpr) String() string {
	var sb strings.Builder
	sb.WriteString(r.metric)

	if len(r.matchers) > 0 {
		sb.WriteByte('{')
		for i, m := range r.matchers {
			if i > 0 {
				sb.WriteByte(',')
			}
			sb.WriteString(m.String())
		}
		sb.WriteByte('}')
	}

	sb.WriteByte('[')
	sb.WriteString(r.duration)
	sb.WriteByte(']')

	if r.offset != "" {
		sb.WriteString(" offset ")
		sb.WriteString(r.offset)
	}

	return sb.String()
}

// WithOffset adds an offset modifier.
func (r *RangeVectorExpr) WithOffset(offset string) *RangeVectorExpr {
	r.offset = offset
	return r
}

// RangeVector creates a range vector selector.
func RangeVector(metric string, duration string, matchers ...LabelMatcher) *RangeVectorExpr {
	return &RangeVectorExpr{
		metric:   metric,
		duration: duration,
		matchers: matchers,
	}
}

// LabelMatcher represents a label matching condition.
type LabelMatcher struct {
	name  string
	op    string
	value string
}

// String returns the label matcher as a string.
func (l LabelMatcher) String() string {
	return fmt.Sprintf(`%s%s"%s"`, l.name, l.op, l.value)
}

// Match creates an equality matcher (=).
func Match(name, value string) LabelMatcher {
	return LabelMatcher{name: name, op: "=", value: value}
}

// NotMatch creates a non-equality matcher (!=).
func NotMatch(name, value string) LabelMatcher {
	return LabelMatcher{name: name, op: "!=", value: value}
}

// MatchRegex creates a regex matcher (=~).
func MatchRegex(name, regex string) LabelMatcher {
	return LabelMatcher{name: name, op: "=~", value: regex}
}

// NotMatchRegex creates a negative regex matcher (!~).
func NotMatchRegex(name, regex string) LabelMatcher {
	return LabelMatcher{name: name, op: "!~", value: regex}
}
