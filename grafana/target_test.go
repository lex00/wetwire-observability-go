package grafana

import (
	"testing"

	"github.com/lex00/wetwire-observability-go/promql"
)

func TestPromTarget(t *testing.T) {
	target := PromTarget("sum(rate(http_requests_total[5m]))")
	if target == nil {
		t.Fatal("PromTarget() returned nil")
	}
	if target.Expr != "sum(rate(http_requests_total[5m]))" {
		t.Errorf("Expr = %v", target.Expr)
	}
	if target.RefID != "A" {
		t.Errorf("RefID = %v, want A", target.RefID)
	}
}

func TestPromTarget_WithRefID(t *testing.T) {
	target := PromTarget("up").WithRefID("B")
	if target.RefID != "B" {
		t.Errorf("RefID = %v, want B", target.RefID)
	}
}

func TestPromTarget_WithLegendFormat(t *testing.T) {
	target := PromTarget("up").WithLegendFormat("{{instance}}")
	if target.LegendFormat != "{{instance}}" {
		t.Errorf("LegendFormat = %v", target.LegendFormat)
	}
}

func TestPromTarget_WithInterval(t *testing.T) {
	target := PromTarget("up").WithInterval("$__rate_interval")
	if target.Interval != "$__rate_interval" {
		t.Errorf("Interval = %v", target.Interval)
	}
}

func TestPromTarget_Instant(t *testing.T) {
	target := PromTarget("up").Instant()
	if !target.IsInstant {
		t.Error("IsInstant should be true")
	}
}

func TestPromTarget_Range(t *testing.T) {
	target := PromTarget("up").Instant().Range()
	if target.IsInstant {
		t.Error("IsInstant should be false")
	}
}

func TestPromTarget_WithTypedExpr(t *testing.T) {
	expr := promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m"))).By("service")
	target := PromTargetExpr(expr)
	expected := "sum by (service) (rate(http_requests_total[5m]))"
	if target.Expr != expected {
		t.Errorf("Expr = %v, want %v", target.Expr, expected)
	}
}

func TestPromTarget_FluentAPI(t *testing.T) {
	target := PromTarget("up").
		WithRefID("A").
		WithLegendFormat("{{job}} - {{instance}}").
		WithInterval("$__rate_interval").
		Instant()

	if target.RefID != "A" {
		t.Errorf("RefID = %v", target.RefID)
	}
	if target.LegendFormat != "{{job}} - {{instance}}" {
		t.Errorf("LegendFormat = %v", target.LegendFormat)
	}
	if !target.IsInstant {
		t.Error("IsInstant should be true")
	}
}

func TestPromTarget_Hide(t *testing.T) {
	target := PromTarget("up").Hide()
	if !target.Hidden {
		t.Error("Hidden should be true")
	}
}

func TestPromTarget_WithDatasource(t *testing.T) {
	target := PromTarget("up").WithDatasource("$datasource")
	if target.Datasource != "$datasource" {
		t.Errorf("Datasource = %v", target.Datasource)
	}
}

func TestLokiTarget(t *testing.T) {
	target := LokiTarget(`{job="api"} |= "error"`)
	if target == nil {
		t.Fatal("LokiTarget() returned nil")
	}
	if target.Expr != `{job="api"} |= "error"` {
		t.Errorf("Expr = %v", target.Expr)
	}
	if target.RefID != "A" {
		t.Errorf("RefID = %v, want A", target.RefID)
	}
}

func TestLokiTarget_WithMaxLines(t *testing.T) {
	target := LokiTarget(`{job="api"}`).WithMaxLines(100)
	if target.MaxLines != 100 {
		t.Errorf("MaxLines = %d, want 100", target.MaxLines)
	}
}
