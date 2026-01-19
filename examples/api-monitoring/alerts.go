package monitoring

import (
	"github.com/lex00/wetwire-observability-go/promql"
	"github.com/lex00/wetwire-observability-go/rules"
)

// HighErrorRate fires when error rate exceeds 5% for 5 minutes.
var HighErrorRate = rules.AlertingRule{
	Alert: "HighErrorRate",
	Expr:  promql.GT(ErrorRateExpr, promql.Scalar(0.05)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "critical",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "High error rate detected",
		"description": "Error rate is above 5% for {{ $labels.service }}",
	},
}

// HighLatency fires when P99 latency exceeds 500ms for 5 minutes.
var HighLatency = rules.AlertingRule{
	Alert: "HighLatency",
	Expr:  promql.GT(LatencyP99Expr, promql.Scalar(0.5)).String(),
	For:   5 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "High latency detected",
		"description": "P99 latency is above 500ms for {{ $labels.service }}",
	},
}

// LowRequestRate fires when request rate drops below 10 rps.
var LowRequestRate = rules.AlertingRule{
	Alert: "LowRequestRate",
	Expr:  promql.LT(RequestRateExpr, promql.Scalar(10)).String(),
	For:   10 * rules.Minute,
	Labels: map[string]string{
		"severity": "warning",
		"team":     "platform",
	},
	Annotations: map[string]string{
		"summary":     "Low request rate detected",
		"description": "Request rate has dropped below 10 rps for {{ $labels.service }}",
	},
}
