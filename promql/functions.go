package promql

import (
	"fmt"
	"strings"
)

// FunctionExpr represents a PromQL function call.
type FunctionExpr struct {
	name string
	args []Expr
}

// String returns the function call as a PromQL string.
func (f *FunctionExpr) String() string {
	var sb strings.Builder
	sb.WriteString(f.name)
	sb.WriteByte('(')
	for i, arg := range f.args {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(arg.String())
	}
	sb.WriteByte(')')
	return sb.String()
}

// Rate calculates the per-second rate of increase of a counter.
func Rate(v *RangeVectorExpr) *FunctionExpr {
	return &FunctionExpr{name: "rate", args: []Expr{v}}
}

// Irate calculates the instant per-second rate of increase.
func Irate(v *RangeVectorExpr) *FunctionExpr {
	return &FunctionExpr{name: "irate", args: []Expr{v}}
}

// Increase calculates the increase in value over a time range.
func Increase(v *RangeVectorExpr) *FunctionExpr {
	return &FunctionExpr{name: "increase", args: []Expr{v}}
}

// Delta calculates the difference between first and last value.
func Delta(v *RangeVectorExpr) *FunctionExpr {
	return &FunctionExpr{name: "delta", args: []Expr{v}}
}

// AggregationExpr represents an aggregation function with optional grouping.
type AggregationExpr struct {
	name     string
	expr     Expr
	by       []string
	without  []string
}

// String returns the aggregation as a PromQL string.
func (a *AggregationExpr) String() string {
	var sb strings.Builder
	sb.WriteString(a.name)

	if len(a.by) > 0 {
		sb.WriteString(" by (")
		sb.WriteString(strings.Join(a.by, ","))
		sb.WriteString(") ")
	} else if len(a.without) > 0 {
		sb.WriteString(" without (")
		sb.WriteString(strings.Join(a.without, ","))
		sb.WriteString(") ")
	}

	sb.WriteByte('(')
	sb.WriteString(a.expr.String())
	sb.WriteByte(')')
	return sb.String()
}

// By adds a "by" clause to group results by the specified labels.
func (a *AggregationExpr) By(labels ...string) *AggregationExpr {
	a.by = labels
	a.without = nil
	return a
}

// Without adds a "without" clause to exclude the specified labels from grouping.
func (a *AggregationExpr) Without(labels ...string) *AggregationExpr {
	a.without = labels
	a.by = nil
	return a
}

// Sum aggregates by summing values.
func Sum(expr Expr) *AggregationExpr {
	return &AggregationExpr{name: "sum", expr: expr}
}

// Avg aggregates by averaging values.
func Avg(expr Expr) *AggregationExpr {
	return &AggregationExpr{name: "avg", expr: expr}
}

// Min aggregates by taking minimum values.
func Min(expr Expr) *AggregationExpr {
	return &AggregationExpr{name: "min", expr: expr}
}

// Max aggregates by taking maximum values.
func Max(expr Expr) *AggregationExpr {
	return &AggregationExpr{name: "max", expr: expr}
}

// Count aggregates by counting elements.
func Count(expr Expr) *AggregationExpr {
	return &AggregationExpr{name: "count", expr: expr}
}

// Stddev aggregates by calculating standard deviation.
func Stddev(expr Expr) *AggregationExpr {
	return &AggregationExpr{name: "stddev", expr: expr}
}

// HistogramQuantile calculates a quantile from a histogram.
func HistogramQuantile(quantile float64, expr Expr) *FunctionExpr {
	return &FunctionExpr{
		name: "histogram_quantile",
		args: []Expr{Scalar(quantile), expr},
	}
}

// P99 is a convenience function for the 99th percentile.
func P99(expr Expr) *FunctionExpr {
	return HistogramQuantile(0.99, expr)
}

// P95 is a convenience function for the 95th percentile.
func P95(expr Expr) *FunctionExpr {
	return HistogramQuantile(0.95, expr)
}

// P90 is a convenience function for the 90th percentile.
func P90(expr Expr) *FunctionExpr {
	return HistogramQuantile(0.9, expr)
}

// P50 is a convenience function for the 50th percentile (median).
func P50(expr Expr) *FunctionExpr {
	return HistogramQuantile(0.5, expr)
}

// Abs returns absolute value.
func Abs(expr Expr) *FunctionExpr {
	return &FunctionExpr{name: "abs", args: []Expr{expr}}
}

// Ceil rounds up to the nearest integer.
func Ceil(expr Expr) *FunctionExpr {
	return &FunctionExpr{name: "ceil", args: []Expr{expr}}
}

// Floor rounds down to the nearest integer.
func Floor(expr Expr) *FunctionExpr {
	return &FunctionExpr{name: "floor", args: []Expr{expr}}
}

// Round rounds to the nearest integer.
func Round(expr Expr) *FunctionExpr {
	return &FunctionExpr{name: "round", args: []Expr{expr}}
}

// Clamp clamps values between min and max.
func Clamp(expr Expr, min, max float64) *FunctionExpr {
	return &FunctionExpr{
		name: "clamp",
		args: []Expr{expr, Scalar(min), Scalar(max)},
	}
}

// ClampMin clamps values to a minimum.
func ClampMin(expr Expr, min float64) *FunctionExpr {
	return &FunctionExpr{
		name: "clamp_min",
		args: []Expr{expr, Scalar(min)},
	}
}

// ClampMax clamps values to a maximum.
func ClampMax(expr Expr, max float64) *FunctionExpr {
	return &FunctionExpr{
		name: "clamp_max",
		args: []Expr{expr, Scalar(max)},
	}
}

// Changes returns the number of times the value changed.
func Changes(v *RangeVectorExpr) *FunctionExpr {
	return &FunctionExpr{name: "changes", args: []Expr{v}}
}

// Resets returns the number of counter resets.
func Resets(v *RangeVectorExpr) *FunctionExpr {
	return &FunctionExpr{name: "resets", args: []Expr{v}}
}

// LabelReplace performs label replacement.
func LabelReplace(expr Expr, dst, replacement, src, regex string) *FunctionExpr {
	return &FunctionExpr{
		name: "label_replace",
		args: []Expr{expr, Raw(fmt.Sprintf(`"%s"`, dst)), Raw(fmt.Sprintf(`"%s"`, replacement)), Raw(fmt.Sprintf(`"%s"`, src)), Raw(fmt.Sprintf(`"%s"`, regex))},
	}
}

// LabelJoin joins label values.
func LabelJoin(expr Expr, dst, sep string, srcLabels ...string) *FunctionExpr {
	args := []Expr{expr, Raw(fmt.Sprintf(`"%s"`, dst)), Raw(fmt.Sprintf(`"%s"`, sep))}
	for _, l := range srcLabels {
		args = append(args, Raw(fmt.Sprintf(`"%s"`, l)))
	}
	return &FunctionExpr{name: "label_join", args: args}
}
