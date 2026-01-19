package monitoring

import (
	"github.com/lex00/wetwire-observability-go/grafana"
	"github.com/lex00/wetwire-observability-go/promql"
)

// ErrorRatePanel displays the current error rate as a stat.
var ErrorRatePanel = grafana.Stat("Error Rate").
	WithTargets(grafana.PromTargetExpr(ErrorRateExpr).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// RequestRatePanel displays requests per second over time.
var RequestRatePanel = grafana.TimeSeries("Request Rate").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval"))).By("service"),
	).WithRefID("A").WithLegendFormat("{{ service }}")).
	WithUnit(grafana.UnitShort)

// LatencyPanel displays P50, P90, P99 latency over time.
var LatencyPanel = grafana.TimeSeries("Request Latency").
	WithTargets(
		grafana.PromTargetExpr(
			promql.HistogramQuantile(0.50,
				promql.Sum(promql.Rate(promql.RangeVector("http_request_duration_seconds_bucket", "$__rate_interval"))).By("le")),
		).WithRefID("A").WithLegendFormat("p50"),
		grafana.PromTargetExpr(
			promql.HistogramQuantile(0.90,
				promql.Sum(promql.Rate(promql.RangeVector("http_request_duration_seconds_bucket", "$__rate_interval"))).By("le")),
		).WithRefID("B").WithLegendFormat("p90"),
		grafana.PromTargetExpr(
			promql.HistogramQuantile(0.99,
				promql.Sum(promql.Rate(promql.RangeVector("http_request_duration_seconds_bucket", "$__rate_interval"))).By("le")),
		).WithRefID("C").WithLegendFormat("p99"),
	).
	WithUnit(grafana.UnitSeconds)

// APIDashboard is the main API service monitoring dashboard.
var APIDashboard = grafana.NewDashboard("api-metrics", "API Service Metrics").
	WithTags("api", "platform").
	WithRows(
		grafana.NewRow("Overview").WithPanels(ErrorRatePanel, RequestRatePanel),
		grafana.NewRow("Latency").WithPanels(LatencyPanel),
	)
