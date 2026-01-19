# wetwire-observability-go FAQ

This FAQ covers questions specific to the Go implementation of wetwire for Prometheus, Alertmanager, and Grafana. For general wetwire questions, see the [central FAQ](https://github.com/lex00/wetwire/blob/main/docs/FAQ.md).

---

## Getting Started

### How do I install wetwire-obs?

See the [README](../README.md#installation) for installation instructions.

### How do I create a new project?

```bash
wetwire-obs init my-monitoring
cd my-monitoring
```

### How do I build configuration files?

```bash
# Standalone configs
wetwire-obs build ./monitoring --mode=standalone

# Prometheus Operator CRDs
wetwire-obs build ./monitoring --mode=operator
```

---

## Syntax

### How do I define a Prometheus scrape config?

```go
var APIServer = prometheus.ScrapeConfig{
    JobName:        "api-server",
    ScrapeInterval: prometheus.Duration("10s"),
    StaticConfigs: []*prometheus.StaticConfig{
        {Targets: []string{"api:8080"}},
    },
}
```

### How do I create a PromQL expression?

Use the `promql` package builders:

```go
// rate(http_requests_total[5m])
expr := promql.Rate(promql.Vector("http_requests_total"), "5m")

// sum by (service) (rate(http_requests_total[5m]))
expr := promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service")

// http_requests_total{status=~"5.."}
expr := promql.Vector("http_requests_total", promql.Match("status", "5.."))
```

### How do I share PromQL expressions between alerts and dashboards?

Define the expression as a variable and reference it:

```go
// Define once
var ErrorRateExpr = promql.Div(
    promql.Sum(promql.Rate(promql.Vector("http_requests_total", promql.Match("status", "5..")), "5m"), "service"),
    promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service"),
)

// Use in alert
var HighErrorRate = rules.AlertingRule{Expr: ErrorRateExpr}

// Use in dashboard
var Panel = grafana.StatPanel{Targets: []any{grafana.PrometheusTarget{Expr: ErrorRateExpr}}}
```

### How do I use Grafana variables in PromQL?

Use the standard Grafana variable syntax:

```go
// $__rate_interval is resolved by Grafana
expr := promql.Rate(promql.Vector("http_requests_total"), "$__rate_interval")

// Custom variables work too
expr := promql.Vector("http_requests_total", promql.Match("service", "$service"))
```

---

## Output Modes

### What's the difference between standalone and operator mode?

**Standalone mode** generates traditional config files:
- `prometheus.yml` - Prometheus configuration
- `alertmanager.yml` - Alertmanager configuration
- `rules/*.yml` - Rule files
- `dashboards/*.json` - Grafana dashboard JSON

**Operator mode** generates Kubernetes CRDs for Prometheus Operator:
- `ServiceMonitor` - Kubernetes-native scrape config
- `PrometheusRule` - Kubernetes-native rules
- `GrafanaDashboard` - Dashboard as ConfigMap

### Can I use both modes?

Yes:

```bash
wetwire-obs build ./monitoring --mode=both
```

---

## Grafana Dashboards

### How does panel auto-positioning work?

Panels are automatically positioned based on row definitions:

```go
var Dashboard = grafana.Dashboard{
    Rows: []grafana.Row{
        {Panels: []any{Panel1, Panel2}},  // Side by side (50% each)
        {Panels: []any{Panel3}},           // Full width below
        {Panels: []any{P4, P5, P6}},       // Three panels (33% each)
    },
}
```

### How do I set panel dimensions explicitly?

Use the `GridPos` field:

```go
var MyPanel = grafana.TimeseriesPanel{
    Title: "Custom Size",
    GridPos: grafana.GridPos{
        X: 0, Y: 0,
        W: 12, H: 8,  // Width 12 (half), Height 8
    },
}
```

### What panel types are supported?

- `TimeseriesPanel` - Time series graphs
- `StatPanel` - Single stat display
- `GaugePanel` - Gauge visualization
- `BarGaugePanel` - Bar gauge
- `TablePanel` - Table view
- `HeatmapPanel` - Heatmap
- `PieChartPanel` - Pie chart
- `LogsPanel` - Log viewer
- `TextPanel` - Text/markdown

---

## Alertmanager

### How do I define alert receivers?

```go
var SlackReceiver = alertmanager.Receiver{
    Name: "slack-notifications",
    SlackConfigs: []*alertmanager.SlackConfig{
        {
            Channel:  "#alerts",
            Username: "alertmanager",
        },
    },
}

var PagerDutyReceiver = alertmanager.Receiver{
    Name: "pagerduty",
    PagerDutyConfigs: []*alertmanager.PagerDutyConfig{
        {ServiceKey: alertmanager.SecretRef("pagerduty-key")},
    },
}
```

### How do I route alerts to different receivers?

```go
var AlertRouting = alertmanager.Route{
    Receiver: "default",
    Routes: []*alertmanager.Route{
        {
            Match:    map[string]string{"severity": "critical"},
            Receiver: "pagerduty",
        },
        {
            Match:    map[string]string{"team": "platform"},
            Receiver: "slack-notifications",
        },
    },
}
```

---

## Troubleshooting

### "cannot find package" errors

Ensure your `go.mod` has the correct dependencies:

```bash
go mod tidy
```

### "undefined: prometheus" or similar import errors

Add the missing import:

```go
import "github.com/lex00/wetwire-observability-go/prometheus"
```

### Build produces empty output

Check that:
1. Resources are declared as package-level `var` statements
2. Resources have the correct type
3. The package path is correct in the build command

### PromQL expression is invalid

Use the `promql` package builders instead of raw strings:

```go
// Bad - raw string
Expr: "rate(http_requests_total[5m])"

// Good - type-safe builder
Expr: promql.Rate(promql.Vector("http_requests_total"), "5m")
```

### Dashboard panels overlap

Use row-based layout or explicit `GridPos`:

```go
// Row-based (recommended)
Rows: []grafana.Row{
    {Panels: []any{Panel1, Panel2}},
}

// Or explicit positioning
GridPos: grafana.GridPos{X: 0, Y: 0, W: 12, H: 8}
```

---

## See Also

- [Wetwire Specification](https://github.com/lex00/wetwire/blob/main/docs/WETWIRE_SPEC.md)
- [CLI Reference](CLI.md)
- [Quick Start](QUICK_START.md)
- [Lint Rules](LINT_RULES.md)
