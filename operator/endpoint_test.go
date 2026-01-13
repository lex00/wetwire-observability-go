package operator

import "testing"

func TestNewEndpoint(t *testing.T) {
	e := NewEndpoint("metrics")
	if e.Port != "metrics" {
		t.Errorf("Port = %q, want metrics", e.Port)
	}
}

func TestEndpoint_WithPath(t *testing.T) {
	e := NewEndpoint("metrics").WithPath("/metrics")
	if e.Path != "/metrics" {
		t.Errorf("Path = %q, want /metrics", e.Path)
	}
}

func TestEndpoint_WithInterval(t *testing.T) {
	e := NewEndpoint("metrics").WithInterval("30s")
	if e.Interval != "30s" {
		t.Errorf("Interval = %q, want 30s", e.Interval)
	}
}

func TestEndpoint_WithScrapeTimeout(t *testing.T) {
	e := NewEndpoint("metrics").WithScrapeTimeout("10s")
	if e.ScrapeTimeout != "10s" {
		t.Errorf("ScrapeTimeout = %q, want 10s", e.ScrapeTimeout)
	}
}

func TestEndpoint_WithScheme(t *testing.T) {
	e := NewEndpoint("metrics").WithScheme("https")
	if e.Scheme != "https" {
		t.Errorf("Scheme = %q, want https", e.Scheme)
	}
}

func TestEndpoint_WithBearerToken(t *testing.T) {
	e := NewEndpoint("metrics").WithBearerTokenFile("/var/run/secrets/token")
	if e.BearerTokenFile != "/var/run/secrets/token" {
		t.Errorf("BearerTokenFile = %q", e.BearerTokenFile)
	}
}

func TestEndpoint_WithTLSConfig(t *testing.T) {
	e := NewEndpoint("metrics").WithTLSConfig(true, "/ca.crt", "/cert.crt", "/key.key")
	if !e.TLSConfig.InsecureSkipVerify {
		t.Error("InsecureSkipVerify should be true")
	}
	if e.TLSConfig.CAFile != "/ca.crt" {
		t.Errorf("CAFile = %q, want /ca.crt", e.TLSConfig.CAFile)
	}
}

func TestEndpoint_AddRelabeling(t *testing.T) {
	e := NewEndpoint("metrics").
		AddRelabeling(KeepLabel("__meta_kubernetes_service_label_app"))
	if len(e.RelabelConfigs) != 1 {
		t.Errorf("len(RelabelConfigs) = %d, want 1", len(e.RelabelConfigs))
	}
}

func TestEndpoint_AddMetricRelabeling(t *testing.T) {
	e := NewEndpoint("metrics").
		AddMetricRelabeling(DropLabel("go_.*"))
	if len(e.MetricRelabelConfigs) != 1 {
		t.Errorf("len(MetricRelabelConfigs) = %d, want 1", len(e.MetricRelabelConfigs))
	}
}

func TestEndpoint_FluentAPI(t *testing.T) {
	e := NewEndpoint("metrics").
		WithPath("/metrics").
		WithInterval("30s").
		WithScrapeTimeout("10s").
		WithScheme("https")

	if e.Port != "metrics" {
		t.Error("Fluent API should preserve port")
	}
	if e.Path != "/metrics" {
		t.Error("Fluent API should set path")
	}
	if e.Interval != "30s" {
		t.Error("Fluent API should set interval")
	}
}

func TestKeepLabel(t *testing.T) {
	r := KeepLabel("__meta_kubernetes_namespace")
	if r.Action != "keep" {
		t.Errorf("Action = %q, want keep", r.Action)
	}
}

func TestDropLabel(t *testing.T) {
	r := DropLabel("go_.*")
	if r.Action != "drop" {
		t.Errorf("Action = %q, want drop", r.Action)
	}
}

func TestReplaceLabel(t *testing.T) {
	r := ReplaceLabel("__meta_kubernetes_namespace", "namespace")
	if r.Action != "replace" {
		t.Errorf("Action = %q, want replace", r.Action)
	}
	if r.TargetLabel != "namespace" {
		t.Errorf("TargetLabel = %q, want namespace", r.TargetLabel)
	}
}
