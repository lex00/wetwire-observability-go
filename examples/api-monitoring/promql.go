// Package monitoring demonstrates API service observability patterns.
//
// This example shows shared PromQL expressions that can be used
// in both alerting rules and Grafana dashboard panels.
package monitoring

import "github.com/lex00/wetwire-observability-go/promql"

// ErrorRateExpr calculates the percentage of 5xx responses.
// Used by both HighErrorRate alert and ErrorRatePanel dashboard.
var ErrorRateExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m",
		promql.MatchRegex("status", "5..")))).By("service"),
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m"))).By("service"),
)

// LatencyP99Expr calculates the 99th percentile latency.
// Used by HighLatency alert and LatencyPanel dashboard.
var LatencyP99Expr = promql.HistogramQuantile(0.99,
	promql.Sum(promql.Rate(
		promql.RangeVector("http_request_duration_seconds_bucket", "5m"),
	)).By("le", "service"),
)

// RequestRateExpr calculates requests per second by service.
var RequestRateExpr = promql.Sum(
	promql.Rate(promql.RangeVector("http_requests_total", "5m")),
).By("service")
