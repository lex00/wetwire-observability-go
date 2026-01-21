package monitoring

import (
	"github.com/lex00/wetwire-observability-go/promql"
	"github.com/lex00/wetwire-observability-go/rules"
)

// Error Rate Alerts

// HighErrorRateWarning fires when error rate exceeds 1%.
var HighErrorRateWarning = rules.AlertingRule{
	Alert: "HighErrorRate",
	Expr:  promql.GT(ErrorRateExpr, promql.Scalar(0.01)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
	},
	Annotations: map[string]string{
		"summary":     "Error rate elevated for {{ $labels.service }}",
		"description": "Error rate is {{ $value | humanizePercentage }} (threshold: 1%)",
		"dashboard":   "https://grafana.example.com/d/multi-service?var-service={{ $labels.service }}",
	},
}

// HighErrorRateCritical fires when error rate exceeds 5%.
var HighErrorRateCritical = rules.AlertingRule{
	Alert: "HighErrorRate",
	Expr:  promql.GT(ErrorRateExpr, promql.Scalar(0.05)).String(),
	For:   2 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
	},
	Annotations: map[string]string{
		"summary":     "Critical error rate for {{ $labels.service }}",
		"description": "Error rate is {{ $value | humanizePercentage }} (threshold: 5%)",
		"runbook_url": "https://runbooks.example.com/high-error-rate",
		"dashboard":   "https://grafana.example.com/d/multi-service?var-service={{ $labels.service }}",
	},
}

// Latency Alerts

// HighLatencyWarning fires when P99 latency exceeds 500ms.
var HighLatencyWarning = rules.AlertingRule{
	Alert: "HighLatency",
	Expr:  promql.GT(LatencyP99Expr, promql.Scalar(0.5)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
	},
	Annotations: map[string]string{
		"summary":     "High latency for {{ $labels.service }}",
		"description": "P99 latency is {{ $value | humanizeDuration }} (threshold: 500ms)",
	},
}

// HighLatencyCritical fires when P99 latency exceeds 1s.
var HighLatencyCritical = rules.AlertingRule{
	Alert: "HighLatency",
	Expr:  promql.GT(LatencyP99Expr, promql.Scalar(1.0)).String(),
	For:   2 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
	},
	Annotations: map[string]string{
		"summary":     "Critical latency for {{ $labels.service }}",
		"description": "P99 latency is {{ $value | humanizeDuration }} (threshold: 1s)",
		"runbook_url": "https://runbooks.example.com/high-latency",
	},
}

// SLO-Based Alerts

// SLOFastBurn fires on fast error budget consumption.
var SLOFastBurn = rules.AlertingRule{
	Alert: "SLOFastBurn",
	Expr: promql.And(
		promql.GT(FastBurnExpr, promql.Scalar(14.4)),
		promql.GT(SlowBurnExpr, promql.Scalar(14.4)),
	).String(),
	For: 2 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"slo":      "availability",
	},
	Annotations: map[string]string{
		"summary":     "Fast error budget burn for {{ $labels.service }}",
		"description": "Error budget is being consumed at {{ $value }}x the sustainable rate",
		"runbook_url": "https://runbooks.example.com/slo-fast-burn",
	},
}

// SLOSlowBurn fires on slow but sustained error budget consumption.
var SLOSlowBurn = rules.AlertingRule{
	Alert: "SLOSlowBurn",
	Expr: promql.And(
		promql.GT(FastBurnExpr, promql.Scalar(6)),
		promql.GT(SlowBurnExpr, promql.Scalar(6)),
	).String(),
	For: 15 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"slo":      "availability",
	},
	Annotations: map[string]string{
		"summary":     "Slow error budget burn for {{ $labels.service }}",
		"description": "Error budget is being consumed at {{ $value }}x the sustainable rate",
	},
}

// Availability Alerts

// ServiceDown fires when a service has no metrics.
var ServiceDown = rules.AlertingRule{
	Alert: "ServiceDown",
	Expr:  "up{job=\"microservices\"} == 0",
	For:   2 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
	},
	Annotations: map[string]string{
		"summary":     "Service {{ $labels.service }} is down",
		"description": "Instance {{ $labels.instance }} has been down for more than 2 minutes",
		"runbook_url": "https://runbooks.example.com/service-down",
	},
}

// LowRequestRate fires when request rate drops significantly.
// Uses raw PromQL for complex subquery expression.
var LowRequestRate = rules.AlertingRule{
	Alert: "LowRequestRate",
	Expr:  "sum by (service) (rate(http_requests_total[5m])) < 0.5 * avg_over_time(sum by (service) (rate(http_requests_total[5m]))[1h:5m])",
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
	},
	Annotations: map[string]string{
		"summary":     "Low request rate for {{ $labels.service }}",
		"description": "Request rate has dropped to {{ $value }} rps, which is below 50% of the 1h average",
	},
}

// Saturation Alerts

// HighQueueDepth fires when request queue is backing up.
var HighQueueDepth = rules.AlertingRule{
	Alert: "HighQueueDepth",
	Expr:  promql.GT(RequestQueueDepthExpr, promql.Scalar(100)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
	},
	Annotations: map[string]string{
		"summary":     "High request queue depth for {{ $labels.service }}",
		"description": "Queue depth is {{ $value }} requests",
	},
}

// AlertingRules is the list of all alerting rules.
var AlertingRules = []rules.AlertingRule{
	HighErrorRateWarning,
	HighErrorRateCritical,
	HighLatencyWarning,
	HighLatencyCritical,
	SLOFastBurn,
	SLOSlowBurn,
	ServiceDown,
	LowRequestRate,
	HighQueueDepth,
}
