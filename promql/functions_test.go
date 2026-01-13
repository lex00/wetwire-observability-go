package promql

import "testing"

func TestRate(t *testing.T) {
	expr := Rate(RangeVector("http_requests_total", "5m"))
	expected := "rate(http_requests_total[5m])"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestRate_WithLabels(t *testing.T) {
	expr := Rate(RangeVector("http_requests_total", "5m", Match("job", "api")))
	expected := `rate(http_requests_total{job="api"}[5m])`
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestIrate(t *testing.T) {
	expr := Irate(RangeVector("http_requests_total", "5m"))
	expected := "irate(http_requests_total[5m])"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestIncrease(t *testing.T) {
	expr := Increase(RangeVector("http_requests_total", "1h"))
	expected := "increase(http_requests_total[1h])"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestDelta(t *testing.T) {
	expr := Delta(RangeVector("temperature", "1h"))
	expected := "delta(temperature[1h])"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestSum(t *testing.T) {
	expr := Sum(Metric("http_requests_total"))
	expected := "sum(http_requests_total)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestSum_By(t *testing.T) {
	expr := Sum(Metric("http_requests_total")).By("service", "status")
	expected := "sum by (service,status) (http_requests_total)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestSum_Without(t *testing.T) {
	expr := Sum(Metric("http_requests_total")).Without("instance")
	expected := "sum without (instance) (http_requests_total)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestAvg(t *testing.T) {
	expr := Avg(Metric("cpu_usage"))
	expected := "avg(cpu_usage)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestAvg_By(t *testing.T) {
	expr := Avg(Metric("cpu_usage")).By("job")
	expected := "avg by (job) (cpu_usage)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestMin(t *testing.T) {
	expr := Min(Metric("memory_bytes"))
	expected := "min(memory_bytes)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestMax(t *testing.T) {
	expr := Max(Metric("memory_bytes"))
	expected := "max(memory_bytes)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestCount(t *testing.T) {
	expr := Count(Metric("up"))
	expected := "count(up)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestCount_By(t *testing.T) {
	expr := Count(Metric("up")).By("job")
	expected := "count by (job) (up)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestStddev(t *testing.T) {
	expr := Stddev(Metric("latency"))
	expected := "stddev(latency)"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestHistogramQuantile(t *testing.T) {
	rateExpr := Rate(RangeVector("http_request_duration_seconds_bucket", "5m"))
	expr := HistogramQuantile(0.99, Sum(rateExpr).By("le"))
	expected := "histogram_quantile(0.99,sum by (le) (rate(http_request_duration_seconds_bucket[5m])))"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestP99(t *testing.T) {
	rateExpr := Rate(RangeVector("http_request_duration_seconds_bucket", "5m"))
	expr := P99(Sum(rateExpr).By("le"))
	expected := "histogram_quantile(0.99,sum by (le) (rate(http_request_duration_seconds_bucket[5m])))"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestP95(t *testing.T) {
	rateExpr := Rate(RangeVector("http_request_duration_seconds_bucket", "5m"))
	expr := P95(Sum(rateExpr).By("le"))
	expected := "histogram_quantile(0.95,sum by (le) (rate(http_request_duration_seconds_bucket[5m])))"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestP90(t *testing.T) {
	rateExpr := Rate(RangeVector("http_request_duration_seconds_bucket", "5m"))
	expr := P90(Sum(rateExpr).By("le"))
	expected := "histogram_quantile(0.9,sum by (le) (rate(http_request_duration_seconds_bucket[5m])))"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}

func TestP50(t *testing.T) {
	rateExpr := Rate(RangeVector("http_request_duration_seconds_bucket", "5m"))
	expr := P50(Sum(rateExpr).By("le"))
	expected := "histogram_quantile(0.5,sum by (le) (rate(http_request_duration_seconds_bucket[5m])))"
	if expr.String() != expected {
		t.Errorf("String() = %v, want %v", expr.String(), expected)
	}
}
