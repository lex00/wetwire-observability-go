package promql

import "testing"

func TestAdd(t *testing.T) {
	expr := Add(Metric("a"), Metric("b"))
	expected := "(a + b)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestSub(t *testing.T) {
	expr := Sub(Metric("a"), Metric("b"))
	expected := "(a - b)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestMul(t *testing.T) {
	expr := Mul(Metric("a"), Metric("b"))
	expected := "(a * b)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestDiv(t *testing.T) {
	expr := Div(Metric("a"), Metric("b"))
	expected := "(a / b)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestGT(t *testing.T) {
	expr := GT(Metric("cpu"), Scalar(90))
	expected := "(cpu > 90)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestLT(t *testing.T) {
	expr := LT(Metric("memory"), Scalar(100))
	expected := "(memory < 100)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestGTE(t *testing.T) {
	expr := GTE(Metric("cpu"), Scalar(80))
	expected := "(cpu >= 80)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestLTE(t *testing.T) {
	expr := LTE(Metric("memory"), Scalar(50))
	expected := "(memory <= 50)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestEq(t *testing.T) {
	expr := Eq(Metric("up"), Scalar(1))
	expected := "(up == 1)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestNeq(t *testing.T) {
	expr := Neq(Metric("up"), Scalar(0))
	expected := "(up != 0)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestComplexExpression(t *testing.T) {
	// Error rate: sum(rate(errors[5m])) / sum(rate(requests[5m]))
	errors := Sum(Rate(RangeVector("http_errors_total", "5m"))).By("service")
	requests := Sum(Rate(RangeVector("http_requests_total", "5m"))).By("service")
	errorRate := Div(errors, requests)

	expected := "(sum by (service) (rate(http_errors_total[5m])) / sum by (service) (rate(http_requests_total[5m])))"
	if errorRate.String() != expected {
		t.Errorf("String() = %v, want %v", errorRate.String(), expected)
	}
}

func TestAlertingExpression(t *testing.T) {
	// Error rate > 5%
	errors := Sum(Rate(RangeVector("http_errors_total", "5m"))).By("service")
	requests := Sum(Rate(RangeVector("http_requests_total", "5m"))).By("service")
	errorRate := Div(errors, requests)
	alert := GT(errorRate, Scalar(0.05))

	expected := "((sum by (service) (rate(http_errors_total[5m])) / sum by (service) (rate(http_requests_total[5m]))) > 0.05)"
	if alert.String() != expected {
		t.Errorf("String() = %v, want %v", alert.String(), expected)
	}
}

func TestBinaryOp_OnLabels(t *testing.T) {
	// Test vector matching with on() modifier
	expr := Div(
		Vector("requests", Match("job", "api")),
		Vector("errors", Match("job", "api")),
	).On("job", "instance")

	expected := `(requests{job="api"} / on (job,instance) errors{job="api"})`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestBinaryOp_Ignoring(t *testing.T) {
	expr := Div(
		Metric("requests"),
		Metric("errors"),
	).Ignoring("status")

	expected := "(requests / ignoring (status) errors)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}
