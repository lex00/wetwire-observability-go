---
title: "FAQ"
---

Frequently asked questions about wetwire-observability-go for Prometheus, Alertmanager, and Grafana configuration.

---

## Getting Started

<details>
<summary>How do I install wetwire-obs?</summary>

```bash
go install github.com/lex00/wetwire-observability-go/cmd/wetwire-obs@latest
```

See the [Quick Start]({{< relref "/quick-start" >}}) for complete setup instructions.
</details>

<details>
<summary>How do I create a new project?</summary>

```bash
wetwire-obs init my-monitoring
cd my-monitoring
```
</details>

<details>
<summary>How do I build configuration files?</summary>

```bash
# Standalone configs
wetwire-obs build ./monitoring --mode=standalone

# Prometheus Operator CRDs
wetwire-obs build ./monitoring --mode=operator
```
</details>

---

## Multi-Backend Support

<details>
<summary>How do I target multiple observability backends?</summary>

Use the `--mode` flag to generate configurations for different backends from the same Go source:

```bash
# Generate standalone Prometheus/Alertmanager/Grafana configs
wetwire-obs build ./monitoring --mode=standalone

# Generate Prometheus Operator CRDs for Kubernetes
wetwire-obs build ./monitoring --mode=operator

# Generate both formats simultaneously
wetwire-obs build ./monitoring --mode=both
```

The same Go struct definitions produce appropriate output for each backend:

| Mode | Output |
|------|--------|
| `standalone` | prometheus.yml, alertmanager.yml, rules/*.yml, dashboards/*.json |
| `operator` | ServiceMonitor, PrometheusRule, AlertmanagerConfig CRDs |
</details>

<details>
<summary>Can I import existing dashboard configurations?</summary>

Yes, use the `import` command to convert existing configurations to Go code:

```bash
# Import Prometheus config
wetwire-obs import prometheus.yml -o ./monitoring/

# Import Alertmanager config
wetwire-obs import alertmanager.yml -o ./monitoring/

# Import Grafana dashboard
wetwire-obs import dashboard.json -o ./monitoring/

# Import rule files
wetwire-obs import rules/*.yml -o ./monitoring/
```

The importer converts YAML/JSON to typed Go structs, including PromQL expressions. See [Import Workflow]({{< relref "/import-workflow" >}}) for detailed migration workflows.
</details>

---

## Syntax

<details>
<summary>How do I define a Prometheus scrape config?</summary>

```go
var APIServer = prometheus.ScrapeConfig{
    JobName:        "api-server",
    ScrapeInterval: prometheus.Duration("10s"),
    StaticConfigs: []*prometheus.StaticConfig{
        {Targets: []string{"api:8080"}},
    },
}
```
</details>

<details>
<summary>How do I create a PromQL expression?</summary>

Use the `promql` package builders:

```go
// rate(http_requests_total[5m])
expr := promql.Rate(promql.Vector("http_requests_total"), "5m")

// sum by (service) (rate(http_requests_total[5m]))
expr := promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service")

// http_requests_total{status=~"5.."}
expr := promql.Vector("http_requests_total", promql.Match("status", "5.."))
```
</details>

<details>
<summary>How do I share PromQL expressions between alerts and dashboards?</summary>

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
</details>

<details>
<summary>How do I use Grafana variables in PromQL?</summary>

Use the standard Grafana variable syntax:

```go
// $__rate_interval is resolved by Grafana
expr := promql.Rate(vector, "$__rate_interval")

// Custom variables work too
expr := promql.Vector("http_requests_total", promql.Match("service", "$service"))
```
</details>

---

## Output Modes

<details>
<summary>What's the difference between standalone and operator mode?</summary>

**Standalone mode** generates traditional config files:
- `prometheus.yml` - Prometheus configuration
- `alertmanager.yml` - Alertmanager configuration
- `rules/*.yml` - Rule files
- `dashboards/*.json` - Grafana dashboard JSON

**Operator mode** generates Kubernetes CRDs for Prometheus Operator:
- `ServiceMonitor` - Kubernetes-native scrape config
- `PrometheusRule` - Kubernetes-native rules
- `GrafanaDashboard` - Dashboard as ConfigMap
</details>

<details>
<summary>Can I use both modes?</summary>

Yes:

```bash
wetwire-obs build ./monitoring --mode=both
```
</details>

---

## Linting and Validation

<details>
<summary>How does the linter help catch errors?</summary>

The linter enforces best practices and catches common mistakes before deployment:

```bash
wetwire-obs lint ./monitoring
```

Key checks include:

| Rule | Description |
|------|-------------|
| WOB022 | Require job_name in ScrapeConfig |
| WOB080 | Require alert name |
| WOB082 | Require severity label on alerts |
| WOB101 | Validate PromQL syntax |
| WOB120 | Require dashboard title |
| WOB200 | Detect hardcoded secrets |

The linter runs automatically during `wetwire-obs build` and can be integrated into CI/CD pipelines. See [Lint Rules]({{< relref "/lint-rules" >}}) for the complete rule reference.
</details>

---

## Project Structure

<details>
<summary>What's the recommended project structure?</summary>

Organize by concern for maintainability:

```
monitoring/
├── prometheus.go     # Global Prometheus config
├── scrape.go         # Scrape configurations
├── alerts.go         # Alerting rules
├── recording.go      # Recording rules
├── dashboards.go     # Grafana dashboards
├── promql.go         # Shared PromQL expressions
└── alertmanager.go   # Alertmanager config and receivers
```

Alternatively, organize by service:

```
monitoring/
├── api/
│   ├── scrape.go
│   ├── alerts.go
│   └── dashboard.go
├── database/
│   ├── scrape.go
│   └── alerts.go
└── shared/
    └── promql.go
```

Both approaches work well. Choose based on team preference and project size.
</details>

---

## Alert Routing

<details>
<summary>How do I handle alert routing?</summary>

Define receivers and routes in Alertmanager configuration:

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

Use `alertmanager.SecretRef()` for sensitive values like API keys and webhook URLs.
</details>

---

## Grafana Dashboards

<details>
<summary>How does panel auto-positioning work?</summary>

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
</details>

<details>
<summary>How do I set panel dimensions explicitly?</summary>

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
</details>

<details>
<summary>What panel types are supported?</summary>

- `TimeseriesPanel` - Time series graphs
- `StatPanel` - Single stat display
- `GaugePanel` - Gauge visualization
- `BarGaugePanel` - Bar gauge
- `TablePanel` - Table view
- `HeatmapPanel` - Heatmap
- `PieChartPanel` - Pie chart
- `LogsPanel` - Log viewer
- `TextPanel` - Text/markdown
</details>

---

## Troubleshooting

<details>
<summary>"cannot find package" errors</summary>

Ensure your `go.mod` has the correct dependencies:

```bash
go mod tidy
```
</details>

<details>
<summary>"undefined: prometheus" or similar import errors</summary>

Add the missing import:

```go
import "github.com/lex00/wetwire-observability-go/prometheus"
```
</details>

<details>
<summary>Build produces empty output</summary>

Check that:
1. Resources are declared as package-level `var` statements
2. Resources have the correct type
3. The package path is correct in the build command
</details>

<details>
<summary>PromQL expression is invalid</summary>

Use the `promql` package builders instead of raw strings:

```go
// Bad - raw string
Expr: "rate(http_requests_total[5m])"

// Good - type-safe builder
Expr: promql.Rate(promql.Vector("http_requests_total"), "5m")
```
</details>

<details>
<summary>Dashboard panels overlap</summary>

Use row-based layout or explicit `GridPos`:

```go
// Row-based (recommended)
Rows: []grafana.Row{
    {Panels: []any{Panel1, Panel2}},
}

// Or explicit positioning
GridPos: grafana.GridPos{X: 0, Y: 0, W: 12, H: 8}
```
</details>

---

## See Also

- [CLI Reference]({{< relref "/cli" >}})
- [Quick Start]({{< relref "/quick-start" >}})
- [Lint Rules]({{< relref "/lint-rules" >}})
