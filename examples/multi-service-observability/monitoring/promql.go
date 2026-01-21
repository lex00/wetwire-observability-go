// Package monitoring provides a complete observability configuration
// for a multi-service microservices architecture.
//
// This example demonstrates integration of all three packages:
// - Prometheus (scrape configs, recording rules)
// - Alertmanager (receivers, routing, inhibition)
// - Grafana (dashboards with variables and drill-down)
package monitoring

import "github.com/lex00/wetwire-observability-go/promql"

// Service names for consistency
const (
	APIGateway   = "api-gateway"
	UserService  = "user-service"
	OrderService = "order-service"
)

// Request Rate Expressions

// RequestRateExpr calculates requests per second by service.
var RequestRateExpr = promql.Sum(
	promql.Rate(promql.RangeVector("http_requests_total", "5m")),
).By("service")

// RequestRateByMethodExpr calculates requests per second by service and method.
var RequestRateByMethodExpr = promql.Sum(
	promql.Rate(promql.RangeVector("http_requests_total", "5m")),
).By("service", "method")

// Error Rate Expressions

// ErrorRateExpr calculates the error rate (5xx responses / total).
var ErrorRateExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m",
		promql.MatchRegex("status", "5..")))).By("service"),
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m"))).By("service"),
)

// ErrorCountExpr counts errors per second by service.
var ErrorCountExpr = promql.Sum(
	promql.Rate(promql.RangeVector("http_requests_total", "5m",
		promql.MatchRegex("status", "5..")))).By("service")

// ClientErrorRateExpr calculates client error rate (4xx responses / total).
var ClientErrorRateExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m",
		promql.MatchRegex("status", "4..")))).By("service"),
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m"))).By("service"),
)

// Latency Expressions

// LatencyP50Expr calculates P50 latency by service.
var LatencyP50Expr = promql.HistogramQuantile(0.50,
	promql.Sum(promql.Rate(
		promql.RangeVector("http_request_duration_seconds_bucket", "5m"),
	)).By("le", "service"),
)

// LatencyP90Expr calculates P90 latency by service.
var LatencyP90Expr = promql.HistogramQuantile(0.90,
	promql.Sum(promql.Rate(
		promql.RangeVector("http_request_duration_seconds_bucket", "5m"),
	)).By("le", "service"),
)

// LatencyP99Expr calculates P99 latency by service.
var LatencyP99Expr = promql.HistogramQuantile(0.99,
	promql.Sum(promql.Rate(
		promql.RangeVector("http_request_duration_seconds_bucket", "5m"),
	)).By("le", "service"),
)

// SLO Expressions

// AvailabilityExpr calculates availability (1 - error rate).
var AvailabilityExpr = promql.Sub(
	promql.Scalar(1),
	ErrorRateExpr,
)

// ErrorBudgetBurnRateExpr calculates the error budget burn rate.
// A value > 1 means burning faster than sustainable for the SLO window.
var ErrorBudgetBurnRateExpr = promql.Div(
	ErrorRateExpr,
	promql.Scalar(0.001), // 99.9% SLO = 0.1% error budget
)

// FastBurnExpr detects fast error budget consumption (5m window).
var FastBurnExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m",
		promql.MatchRegex("status", "5..")))).By("service"),
	promql.Mul(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "5m"))).By("service"),
		promql.Scalar(0.001),
	),
)

// SlowBurnExpr detects slow error budget consumption (30m window).
var SlowBurnExpr = promql.Div(
	promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "30m",
		promql.MatchRegex("status", "5..")))).By("service"),
	promql.Mul(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "30m"))).By("service"),
		promql.Scalar(0.001),
	),
)

// Service Health Expressions

// ServiceUpExpr returns 1 if the service is up (has metrics).
var ServiceUpExpr = promql.Metric("up")

// ServiceInstanceCountExpr counts running instances per service.
var ServiceInstanceCountExpr = promql.Count(
	promql.Vector("up", promql.Match("job", "microservices")),
).By("service")

// ActiveConnectionsExpr counts active connections by service.
var ActiveConnectionsExpr = promql.Sum(
	promql.Metric("http_connections_active"),
).By("service")

// Saturation Expressions

// RequestQueueDepthExpr measures request queue depth.
var RequestQueueDepthExpr = promql.Avg(
	promql.Metric("http_request_queue_length"),
).By("service")

// ThreadPoolUtilizationExpr measures thread pool utilization.
var ThreadPoolUtilizationExpr = promql.Div(
	promql.Metric("thread_pool_active_threads"),
	promql.Metric("thread_pool_max_threads"),
)
