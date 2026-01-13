package promql

import "testing"

func TestRaw(t *testing.T) {
	expr := Raw("sum(rate(http_requests_total[5m])) by (service)")
	if expr.String() != "sum(rate(http_requests_total[5m])) by (service)" {
		t.Errorf("String() = %v", expr.String())
	}
}

func TestMetric(t *testing.T) {
	expr := Metric("up")
	if expr.String() != "up" {
		t.Errorf("String() = %v", expr.String())
	}
}

func TestVector(t *testing.T) {
	expr := Vector("http_requests_total")
	if expr.String() != "http_requests_total" {
		t.Errorf("String() = %v", expr.String())
	}
}

func TestVector_WithLabels(t *testing.T) {
	expr := Vector("http_requests_total", Match("job", "api"), Match("status", "200"))
	expected := `http_requests_total{job="api",status="200"}`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestVector_WithRegex(t *testing.T) {
	expr := Vector("http_requests_total", MatchRegex("status", "5.."))
	expected := `http_requests_total{status=~"5.."}`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestVector_WithNotMatch(t *testing.T) {
	expr := Vector("http_requests_total", NotMatch("job", "test"))
	expected := `http_requests_total{job!="test"}`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestVector_WithNotRegex(t *testing.T) {
	expr := Vector("http_requests_total", NotMatchRegex("status", "4.."))
	expected := `http_requests_total{status!~"4.."}`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestVector_MixedMatchers(t *testing.T) {
	expr := Vector("http_requests_total",
		Match("job", "api"),
		MatchRegex("status", "5.."),
		NotMatch("env", "test"),
	)
	expected := `http_requests_total{job="api",status=~"5..",env!="test"}`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestScalar(t *testing.T) {
	tests := []struct {
		value float64
		want  string
	}{
		{0, "0"},
		{1, "1"},
		{0.5, "0.5"},
		{0.95, "0.95"},
		{100, "100"},
		{-1, "-1"},
	}

	for _, tt := range tests {
		expr := Scalar(tt.value)
		if expr.String() != tt.want {
			t.Errorf("Scalar(%v).String() = %v, want %v", tt.value, expr.String(), tt.want)
		}
	}
}

func TestRangeVector(t *testing.T) {
	expr := RangeVector("http_requests_total", "5m")
	expected := "http_requests_total[5m]"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestRangeVector_WithLabels(t *testing.T) {
	expr := RangeVector("http_requests_total", "1h", Match("job", "api"))
	expected := `http_requests_total{job="api"}[1h]`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestVector_WithOffset(t *testing.T) {
	expr := Vector("http_requests_total").WithOffset("1h")
	expected := "http_requests_total offset 1h"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestRangeVector_WithOffset(t *testing.T) {
	expr := RangeVector("http_requests_total", "5m").WithOffset("1h")
	expected := "http_requests_total[5m] offset 1h"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestLabelMatcher_String(t *testing.T) {
	tests := []struct {
		matcher LabelMatcher
		want    string
	}{
		{Match("job", "api"), `job="api"`},
		{NotMatch("job", "test"), `job!="test"`},
		{MatchRegex("status", "5.."), `status=~"5.."`},
		{NotMatchRegex("status", "2.."), `status!~"2.."`},
	}

	for _, tt := range tests {
		if tt.matcher.String() != tt.want {
			t.Errorf("String() = %v, want %v", tt.matcher.String(), tt.want)
		}
	}
}

func TestExpr_InAlert(t *testing.T) {
	// Example of using PromQL expressions in alerting rules
	errorRate := Raw("sum(rate(http_requests_total{status=~\"5..\"}[5m])) / sum(rate(http_requests_total[5m]))")

	// The expression string can be used directly in AlertingRule.Expr
	if errorRate.String() != `sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m]))` {
		t.Errorf("errorRate.String() = %v", errorRate.String())
	}
}
