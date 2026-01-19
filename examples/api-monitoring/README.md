# API Monitoring Example

This example demonstrates a complete API service monitoring setup using wetwire-observability patterns.

## What's Included

| File | Description |
|------|-------------|
| `promql.go` | Shared PromQL expressions used in alerts and dashboards |
| `scrapes.go` | Prometheus scrape configurations for API, database, and Redis |
| `alerts.go` | Alerting rules for error rate, latency, and request rate |
| `recording.go` | Recording rules for pre-computed metrics |
| `dashboard.go` | Grafana dashboard with error, request, and latency panels |

## Key Patterns

### 1. Shared PromQL Expressions

Define expressions once, use in both alerts and dashboards:

```go
// promql.go
var ErrorRateExpr = promql.Div(...)

// alerts.go - used in alert
var HighErrorRate = rules.AlertingRule{
    Expr: promql.GT(ErrorRateExpr, promql.Scalar(0.05)),
}

// dashboard.go - used in panel
var ErrorRatePanel = grafana.StatPanel{
    Targets: []any{grafana.PrometheusTarget{Expr: ErrorRateExpr}},
}
```

### 2. Flat Variable Declarations

All resources are top-level variables with direct references:

```go
var APIScrape = prometheus.ScrapeConfig{...}
var HighErrorRate = rules.AlertingRule{...}
var APIDashboard = grafana.Dashboard{...}
```

### 3. Row-Based Dashboard Layout

Panels are auto-positioned in rows:

```go
var APIDashboard = grafana.Dashboard{
    Rows: []grafana.Row{
        {Panels: []any{Panel1, Panel2}},  // Side by side
        {Panels: []any{Panel3}},           // Full width
    },
}
```

## Build Commands

```bash
# Generate standalone configs
wetwire-obs build . --mode=standalone

# Generate Prometheus Operator CRDs
wetwire-obs build . --mode=operator

# Lint for issues
wetwire-obs lint .

# List all resources
wetwire-obs list .
```

## Output Files

### Standalone Mode
- `prometheus.yml` - Prometheus configuration with scrape configs
- `rules/alerts.yml` - Alerting rules
- `rules/recording.yml` - Recording rules
- `dashboards/api-service-metrics.json` - Grafana dashboard

### Operator Mode
- `servicemonitor.yaml` - Prometheus Operator ServiceMonitor
- `prometheusrule.yaml` - Prometheus Operator PrometheusRule
- `grafanadashboard.yaml` - Grafana Dashboard ConfigMap
