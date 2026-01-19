# Examples Reference

The `examples/` directory contains reference implementations demonstrating wetwire-observability patterns. These serve as:

1. **Learning resources** - See how to define configs, alerts, and dashboards
2. **Test artifacts** - Validate the build pipeline
3. **Starting points** - Copy and modify for your needs

## Directory Structure

```
examples/
├── api-monitoring/         # API service monitoring
│   ├── scrapes.go         # Prometheus scrape configs
│   ├── alerts.go          # Alerting rules
│   ├── recording.go       # Recording rules
│   └── dashboard.go       # Grafana dashboard
│
├── kubernetes/            # Kubernetes monitoring
│   ├── operator.go        # ServiceMonitor definitions
│   └── alerts.go          # K8s-specific alerts
│
└── full-stack/            # Complete example
    ├── prometheus.go      # Full Prometheus config
    ├── alertmanager.go    # Alertmanager config
    └── dashboards/        # Multiple dashboards
```

## API Monitoring Example

A complete API monitoring setup demonstrating core patterns.

### scrapes.go

```go
package monitoring

import (
    "time"
    "github.com/lex00/wetwire-observability-go/prometheus"
)

var APIScrape = prometheus.ScrapeConfig{
    JobName:        "api",
    ScrapeInterval: 15 * time.Second,
    StaticConfigs: []prometheus.StaticConfig{
        {Targets: []string{"api:8080"}},
    },
    MetricsPath: "/metrics",
    Labels: map[string]string{
        "team": "platform",
    },
}
```

### alerts.go

```go
package monitoring

import (
    "time"
    "github.com/lex00/wetwire-observability-go/rules"
    "github.com/lex00/wetwire-observability-go/promql"
)

// Shared PromQL expression for error rate
var ErrorRateExpr = promql.Div(
    promql.Sum(promql.Rate(promql.Vector("http_requests_total",
        promql.Match("status", "5..")), "5m"), "service"),
    promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service"),
)

var HighErrorRate = rules.AlertingRule{
    Alert:    "HighErrorRate",
    Expr:     promql.GT(ErrorRateExpr, promql.Scalar(0.05)),
    For:      5 * time.Minute,
    Severity: rules.Critical,
    Labels: map[string]string{
        "team": "platform",
    },
    Annotations: map[string]string{
        "summary":     "High error rate detected",
        "description": "Error rate is above 5% for {{ $labels.service }}",
    },
}

var HighLatency = rules.AlertingRule{
    Alert: "HighLatency",
    Expr: promql.GT(
        promql.Histogram_quantile(0.99,
            promql.Sum(promql.Rate(
                promql.Vector("http_request_duration_seconds_bucket"),
                "5m",
            ), "le", "service"),
        ),
        promql.Scalar(0.5),
    ),
    For:      5 * time.Minute,
    Severity: rules.Warning,
}
```

### dashboard.go

```go
package monitoring

import (
    "github.com/lex00/wetwire-observability-go/grafana"
)

var ErrorRatePanel = grafana.StatPanel{
    Title: "Error Rate",
    Targets: []any{
        grafana.PrometheusTarget{
            RefID: "A",
            Expr:  ErrorRateExpr,  // Reuse same expression from alerts
        },
    },
    Unit:       "percentunit",
    Thresholds: []grafana.Threshold{
        {Value: 0, Color: "green"},
        {Value: 0.01, Color: "yellow"},
        {Value: 0.05, Color: "red"},
    },
}

var RequestRatePanel = grafana.GraphPanel{
    Title: "Request Rate",
    Targets: []any{
        grafana.PrometheusTarget{
            RefID:  "A",
            Expr:   promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "$__rate_interval"), "service"),
            Legend: "{{ service }}",
        },
    },
    Unit: "reqps",
}

var APIDashboard = grafana.Dashboard{
    Title: "API Metrics",
    Tags:  []string{"api", "platform"},
    Rows: []grafana.Row{
        {Panels: []any{ErrorRatePanel, RequestRatePanel}},  // Side by side
    },
}
```

## Using Examples

### View the generated configs

```bash
cd examples/api-monitoring
wetwire-obs build . --mode=standalone
```

### List resources

```bash
cd examples/api-monitoring
wetwire-obs list .
```

### Copy as starting point

```bash
cp -r examples/api-monitoring ./my-monitoring
cd my-monitoring
# Edit files
wetwire-obs build . --mode=standalone
```

## Notable Examples

| Example | Description |
|---------|-------------|
| `api-monitoring/` | API service with alerts and dashboard |
| `kubernetes/` | K8s monitoring with Prometheus Operator |
| `full-stack/` | Complete Prometheus + Alertmanager + Grafana |

## Key Patterns Demonstrated

### 1. Shared PromQL Expressions

Define once, use in alerts and dashboards:

```go
var ErrorRateExpr = promql.Div(...)  // Define once

var Alert = rules.AlertingRule{Expr: ErrorRateExpr}  // Use in alert
var Panel = grafana.StatPanel{Targets: []any{
    grafana.PrometheusTarget{Expr: ErrorRateExpr},  // Use in dashboard
}}
```

### 2. Row-Based Dashboard Layout

Auto-positioned panels:

```go
var Dashboard = grafana.Dashboard{
    Rows: []grafana.Row{
        {Panels: []any{Panel1, Panel2}},  // Row 1: side by side
        {Panels: []any{Panel3}},           // Row 2: full width
    },
}
```

### 3. Dual Output Mode

Same Go code, different outputs:

```bash
# Standalone configs
wetwire-obs build . --mode=standalone

# Prometheus Operator CRDs
wetwire-obs build . --mode=operator
```

## Building Examples

To build and verify all examples:

```bash
# Build all
for dir in examples/*/; do
    echo "Building $dir..."
    wetwire-obs build "$dir" --mode=standalone || exit 1
done

# Lint all
wetwire-obs lint ./examples/...
```

## Notes

- Examples demonstrate patterns, not production-ready configs
- Adjust scrape intervals, thresholds, etc. for your environment
- Run `wetwire-obs lint ./examples/...` to check for issues

## See Also

- [Quick Start](QUICK_START.md) - Getting started guide
- [Developer Guide](DEVELOPERS.md) - Development workflow
- [Internals](INTERNALS.md) - How discovery works
