package monitoring

import "github.com/lex00/wetwire-observability-go/rules"

// ErrorRatio5m is a pre-computed error ratio for dashboards.
var ErrorRatio5m = rules.RecordingRule{
	Record: "service:http_error_ratio:5m",
	Expr:   ErrorRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// RequestRate5m is a pre-computed request rate for dashboards.
var RequestRate5m = rules.RecordingRule{
	Record: "service:http_requests:rate5m",
	Expr:   RequestRateExpr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}

// LatencyP995m is a pre-computed P99 latency for dashboards.
var LatencyP995m = rules.RecordingRule{
	Record: "service:http_latency_p99:5m",
	Expr:   LatencyP99Expr.String(),
	Labels: map[string]string{
		"aggregation": "5m",
	},
}
