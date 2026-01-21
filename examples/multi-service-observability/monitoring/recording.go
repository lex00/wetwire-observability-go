package monitoring

import "github.com/lex00/wetwire-observability-go/rules"

// Request Rate Recording Rules

// RequestRate5m pre-computes request rate by service.
var RequestRate5m = rules.RecordingRule{
	Record: "service:http_requests:rate5m",
	Expr:   RequestRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// RequestRateByMethod5m pre-computes request rate by service and method.
var RequestRateByMethod5m = rules.RecordingRule{
	Record: "service:http_requests_by_method:rate5m",
	Expr:   RequestRateByMethodExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// Error Rate Recording Rules

// ErrorRatio5m pre-computes error ratio by service.
var ErrorRatio5m = rules.RecordingRule{
	Record: "service:http_errors:ratio5m",
	Expr:   ErrorRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// ErrorCount5m pre-computes error count by service.
var ErrorCount5m = rules.RecordingRule{
	Record: "service:http_errors:rate5m",
	Expr:   ErrorCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// ClientErrorRatio5m pre-computes client error ratio.
var ClientErrorRatio5m = rules.RecordingRule{
	Record: "service:http_client_errors:ratio5m",
	Expr:   ClientErrorRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// Latency Recording Rules

// LatencyP50_5m pre-computes P50 latency by service.
var LatencyP50_5m = rules.RecordingRule{
	Record: "service:http_latency_p50:5m",
	Expr:   LatencyP50Expr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"quantile":    "0.50",
	},
}

// LatencyP90_5m pre-computes P90 latency by service.
var LatencyP90_5m = rules.RecordingRule{
	Record: "service:http_latency_p90:5m",
	Expr:   LatencyP90Expr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"quantile":    "0.90",
	},
}

// LatencyP99_5m pre-computes P99 latency by service.
var LatencyP99_5m = rules.RecordingRule{
	Record: "service:http_latency_p99:5m",
	Expr:   LatencyP99Expr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"quantile":    "0.99",
	},
}

// SLO Recording Rules

// Availability5m pre-computes availability by service.
var Availability5m = rules.RecordingRule{
	Record: "service:availability:5m",
	Expr:   AvailabilityExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// ErrorBudgetBurnRate5m pre-computes error budget burn rate.
var ErrorBudgetBurnRate5m = rules.RecordingRule{
	Record: "service:error_budget_burn_rate:5m",
	Expr:   ErrorBudgetBurnRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
		"slo":         "99.9",
	},
}

// FastBurn5m pre-computes fast burn rate.
var FastBurn5m = rules.RecordingRule{
	Record: "service:slo_burn_rate:fast",
	Expr:   FastBurnExpr.String(),
	Labels: map[string]string{
		"window": "5m",
	},
}

// SlowBurn30m pre-computes slow burn rate.
var SlowBurn30m = rules.RecordingRule{
	Record: "service:slo_burn_rate:slow",
	Expr:   SlowBurnExpr.String(),
	Labels: map[string]string{
		"window": "30m",
	},
}

// Saturation Recording Rules

// InstanceCount pre-computes instance count by service.
var InstanceCount = rules.RecordingRule{
	Record: "service:instance_count:count",
	Expr:   ServiceInstanceCountExpr.String(),
	Labels: map[string]string{
		"aggregation": "instant",
	},
}

// RecordingRules is the list of all recording rules.
var RecordingRules = []rules.RecordingRule{
	RequestRate5m,
	RequestRateByMethod5m,
	ErrorRatio5m,
	ErrorCount5m,
	ClientErrorRatio5m,
	LatencyP50_5m,
	LatencyP90_5m,
	LatencyP99_5m,
	Availability5m,
	ErrorBudgetBurnRate5m,
	FastBurn5m,
	SlowBurn30m,
	InstanceCount,
}
