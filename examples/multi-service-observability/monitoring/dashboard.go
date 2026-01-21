package monitoring

import (
	"github.com/lex00/wetwire-observability-go/grafana"
	"github.com/lex00/wetwire-observability-go/promql"
)

// Dashboard Variables

// ServiceVariable allows filtering by service.
var ServiceVariable = grafana.QueryVar("service", "label_values(up{job=\"microservices\"}, service)").
	WithLabel("Service").
	MultiSelect().
	IncludeAll()

// EnvironmentVariable allows filtering by environment.
var EnvironmentVariable = grafana.QueryVar("environment", "label_values(up{job=\"microservices\"}, env)").
	WithLabel("Environment").
	IncludeAll()

// IntervalVariable allows selecting aggregation interval.
var IntervalVariable = grafana.IntervalVar("interval", "1m", "5m", "15m", "1h").
	WithLabel("Interval").
	WithDefault("5m")

// Overview Panels

// ErrorRateStatPanel shows current error rate.
var ErrorRateStatPanel = grafana.Stat("Error Rate").
	WithTargets(grafana.PromTargetExpr(
		promql.Avg(promql.Vector("service:http_errors:ratio5m",
			promql.MatchRegex("service", "$service"))),
	).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// RequestRateStatPanel shows current request rate.
var RequestRateStatPanel = grafana.Stat("Request Rate").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Vector("service:http_requests:rate5m",
			promql.MatchRegex("service", "$service"))),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// LatencyP99StatPanel shows current P99 latency.
var LatencyP99StatPanel = grafana.Stat("P99 Latency").
	WithTargets(grafana.PromTargetExpr(
		promql.Max(promql.Vector("service:http_latency_p99:5m",
			promql.MatchRegex("service", "$service"))),
	).WithRefID("A")).
	WithUnit(grafana.UnitSeconds)

// AvailabilityStatPanel shows current availability.
var AvailabilityStatPanel = grafana.Stat("Availability").
	WithTargets(grafana.PromTargetExpr(
		promql.Avg(promql.Vector("service:availability:5m",
			promql.MatchRegex("service", "$service"))),
	).WithRefID("A")).
	WithUnit(grafana.UnitPercentUnit)

// Time Series Panels

// ErrorRatePanel shows error rate over time.
var ErrorRatePanel = grafana.TimeSeries("Error Rate").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval",
			promql.MatchRegex("status", "5.."),
			promql.MatchRegex("service", "$service")))).By("service"),
	).WithRefID("A").WithLegendFormat("{{ service }}")).
	WithUnit(grafana.UnitShort)

// RequestRatePanel shows request rate over time.
var RequestRatePanel = grafana.TimeSeries("Request Rate").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval",
			promql.MatchRegex("service", "$service")))).By("service"),
	).WithRefID("A").WithLegendFormat("{{ service }}")).
	WithUnit(grafana.UnitShort)

// LatencyPanel shows latency percentiles over time.
var LatencyPanel = grafana.TimeSeries("Request Latency").
	WithTargets(
		grafana.PromTargetExpr(
			promql.HistogramQuantile(0.50,
				promql.Sum(promql.Rate(promql.RangeVector("http_request_duration_seconds_bucket", "$__rate_interval",
					promql.MatchRegex("service", "$service")))).By("le", "service")),
		).WithRefID("A").WithLegendFormat("p50 {{ service }}"),
		grafana.PromTargetExpr(
			promql.HistogramQuantile(0.90,
				promql.Sum(promql.Rate(promql.RangeVector("http_request_duration_seconds_bucket", "$__rate_interval",
					promql.MatchRegex("service", "$service")))).By("le", "service")),
		).WithRefID("B").WithLegendFormat("p90 {{ service }}"),
		grafana.PromTargetExpr(
			promql.HistogramQuantile(0.99,
				promql.Sum(promql.Rate(promql.RangeVector("http_request_duration_seconds_bucket", "$__rate_interval",
					promql.MatchRegex("service", "$service")))).By("le", "service")),
		).WithRefID("C").WithLegendFormat("p99 {{ service }}"),
	).
	WithUnit(grafana.UnitSeconds)

// SLO Panels

// ErrorBudgetBurnPanel shows error budget burn rate.
var ErrorBudgetBurnPanel = grafana.TimeSeries("Error Budget Burn Rate").
	WithTargets(grafana.PromTargetExpr(
		promql.Vector("service:error_budget_burn_rate:5m",
			promql.MatchRegex("service", "$service")),
	).WithRefID("A").WithLegendFormat("{{ service }}")).
	WithUnit(grafana.UnitShort)

// AvailabilityPanel shows availability over time.
var AvailabilityPanel = grafana.TimeSeries("Availability").
	WithTargets(grafana.PromTargetExpr(
		promql.Vector("service:availability:5m",
			promql.MatchRegex("service", "$service")),
	).WithRefID("A").WithLegendFormat("{{ service }}")).
	WithUnit(grafana.UnitPercentUnit)

// Service Health Panels

// ServiceUpPanel shows service instances.
var ServiceUpPanel = grafana.TimeSeries("Service Instances").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Vector("up",
			promql.Match("job", "microservices"),
			promql.MatchRegex("service", "$service"))).By("service"),
	).WithRefID("A").WithLegendFormat("{{ service }}")).
	WithUnit(grafana.UnitShort)

// InstanceCountPanel shows instance count by service.
var InstanceCountPanel = grafana.Stat("Instances").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Vector("up",
			promql.Match("job", "microservices"),
			promql.MatchRegex("service", "$service"))),
	).WithRefID("A")).
	WithUnit(grafana.UnitShort)

// Per-Service Panels

// APIGatewayRequestsPanel shows API Gateway requests by method.
var APIGatewayRequestsPanel = grafana.TimeSeries("API Gateway Requests by Method").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval",
			promql.Match("service", APIGateway)))).By("method"),
	).WithRefID("A").WithLegendFormat("{{ method }}")).
	WithUnit(grafana.UnitShort)

// APIGatewayStatusPanel shows API Gateway response status distribution.
var APIGatewayStatusPanel = grafana.TimeSeries("API Gateway Response Status").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval",
			promql.Match("service", APIGateway)))).By("status"),
	).WithRefID("A").WithLegendFormat("{{ status }}")).
	WithUnit(grafana.UnitShort)

// UserServicePanel shows User Service metrics.
var UserServicePanel = grafana.TimeSeries("User Service Operations").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval",
			promql.Match("service", UserService)))).By("endpoint"),
	).WithRefID("A").WithLegendFormat("{{ endpoint }}")).
	WithUnit(grafana.UnitShort)

// OrderServicePanel shows Order Service metrics.
var OrderServicePanel = grafana.TimeSeries("Order Service Operations").
	WithTargets(grafana.PromTargetExpr(
		promql.Sum(promql.Rate(promql.RangeVector("http_requests_total", "$__rate_interval",
			promql.Match("service", OrderService)))).By("endpoint"),
	).WithRefID("A").WithLegendFormat("{{ endpoint }}")).
	WithUnit(grafana.UnitShort)

// Dashboard Definition

// MultiServiceDashboard is the main multi-service observability dashboard.
var MultiServiceDashboard = grafana.NewDashboard("multi-service", "Multi-Service Observability").
	WithDescription("Complete observability for microservices platform").
	WithTags("microservices", "slo", "platform").
	WithRefresh("30s").
	WithVariables(
		ServiceVariable,
		EnvironmentVariable,
		IntervalVariable,
	).
	WithRows(
		grafana.NewRow("Overview").WithPanels(
			ErrorRateStatPanel,
			RequestRateStatPanel,
			LatencyP99StatPanel,
			AvailabilityStatPanel,
			InstanceCountPanel,
		),
		grafana.NewRow("Request Flow").WithPanels(
			RequestRatePanel,
			ErrorRatePanel,
		),
		grafana.NewRow("Latency").WithPanels(
			LatencyPanel,
		),
		grafana.NewRow("SLO").WithPanels(
			AvailabilityPanel,
			ErrorBudgetBurnPanel,
		),
		grafana.NewRow("Service Health").WithPanels(
			ServiceUpPanel,
		),
		grafana.NewRow("API Gateway").WithPanels(
			APIGatewayRequestsPanel,
			APIGatewayStatusPanel,
		).Collapsed(),
		grafana.NewRow("User Service").WithPanels(
			UserServicePanel,
		).Collapsed(),
		grafana.NewRow("Order Service").WithPanels(
			OrderServicePanel,
		).Collapsed(),
	)
