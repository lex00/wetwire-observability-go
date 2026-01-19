# Internals

This document covers the internal architecture of wetwire-observability-go for contributors and maintainers.

**Contents:**
- [AST Discovery](#ast-discovery) - How resource discovery works
- [Config Generation](#config-generation) - How configs are built
- [PromQL Builder](#promql-builder) - Shared expression types
- [Dual Output Mode](#dual-output-mode) - Standalone vs Operator
- [Linter Architecture](#linter-architecture) - How lint rules work

---

## AST Discovery

wetwire-observability uses Go's `go/ast` package to discover configuration declarations without executing user code.

### How It Works

When you define a scrape config as a package-level variable:

```go
var APIScrape = prometheus.ScrapeConfig{
    JobName: "api",
    StaticConfigs: []prometheus.StaticConfig{
        {Targets: []string{"localhost:8080"}},
    },
}
```

The discovery phase:
1. Parses Go source files using `go/parser`
2. Walks the AST looking for `var` declarations
3. Identifies composite literals with prometheus/alertmanager/grafana types
4. Extracts metadata: name, type, file, line, dependencies

### Discovery API

```go
import "github.com/lex00/wetwire-observability-go/internal/discover"

resources, err := discover.DiscoverAll("./monitoring/...")

// Access discovered resources
for _, scrape := range resources.ScrapeConfigs {
    fmt.Printf("%s: %s at %s:%d\n", scrape.Name, scrape.JobName, scrape.File, scrape.Line)
}
```

### What Gets Discovered

| Type | Example | Package |
|------|---------|---------|
| ScrapeConfig | `var APIScrape = prometheus.ScrapeConfig{...}` | prometheus |
| AlertingRule | `var HighErrorRate = rules.AlertingRule{...}` | rules |
| RecordingRule | `var ErrorRate = rules.RecordingRule{...}` | rules |
| Route | `var DefaultRoute = alertmanager.Route{...}` | alertmanager |
| Receiver | `var SlackReceiver = alertmanager.Receiver{...}` | alertmanager |
| Dashboard | `var APIDashboard = grafana.Dashboard{...}` | grafana |
| Panel | `var ErrorRatePanel = grafana.StatPanel{...}` | grafana |

---

## Config Generation

The builder constructs configuration files from discovered resources through a three-stage pipeline.

### Generation Pipeline

```
┌─────────────────────────────────────────────────────────────┐
│                   Config Generation Pipeline                 │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────┐  │
│  │   DISCOVER  │ ─▶ │  SERIALIZE  │ ─▶ │     OUTPUT      │  │
│  └─────────────┘    └─────────────┘    └─────────────────┘  │
│                                                              │
│  AST parsing        Convert structs    Write config files   │
│  finds configs      to YAML/JSON                            │
└─────────────────────────────────────────────────────────────┘
```

### Build Process

```go
import "github.com/lex00/wetwire-observability-go/internal/builder"

// Build all configs
output, err := builder.Build(resources, builder.Options{
    Mode: builder.Standalone,  // or Operator
})

// output.Prometheus - prometheus.yml content
// output.Alertmanager - alertmanager.yml content
// output.Rules - map[filename]content
// output.Dashboards - map[filename]content
```

### Type Mapping

#### Prometheus

| Go Type | YAML Output |
|---------|-------------|
| `prometheus.Config` | Root prometheus.yml |
| `prometheus.ScrapeConfig` | `scrape_configs` entry |
| `prometheus.StaticConfig` | `static_configs` entry |
| `prometheus.RemoteWrite` | `remote_write` entry |

#### Alertmanager

| Go Type | YAML Output |
|---------|-------------|
| `alertmanager.Config` | Root alertmanager.yml |
| `alertmanager.Route` | `route` section |
| `alertmanager.Receiver` | `receivers` entry |

#### Rules

| Go Type | YAML Output |
|---------|-------------|
| `rules.AlertingRule` | Alert in rule group |
| `rules.RecordingRule` | Recording rule in group |
| `rules.RuleGroup` | Rule group file |

#### Grafana

| Go Type | JSON Output |
|---------|-------------|
| `grafana.Dashboard` | Dashboard JSON |
| `grafana.StatPanel` | Panel with type "stat" |
| `grafana.GraphPanel` | Panel with type "graph" |
| `grafana.Row` | Row container |

### Dependency Resolution

Resources can reference each other. The builder resolves these:

```go
// PromQL expression used in both alert and dashboard
var ErrorRateExpr = promql.GT(...)

var HighErrorRate = rules.AlertingRule{
    Expr: ErrorRateExpr,  // Reference
}

var ErrorRatePanel = grafana.StatPanel{
    Targets: []any{grafana.PrometheusTarget{Expr: ErrorRateExpr}},  // Same reference
}
```

---

## PromQL Builder

The `promql` package provides typed expression builders shared across alerts and dashboards.

### Expression Types

```go
// Vector selector
promql.Vector("http_requests_total", promql.Match("status", "5.."))

// Functions
promql.Rate(vector, "5m")
promql.Sum(expr, "service")
promql.Avg(expr)

// Operators
promql.Div(expr1, expr2)
promql.GT(expr, promql.Scalar(0.05))
promql.And(expr1, expr2)
```

### Serialization

Expressions serialize to PromQL strings:

```go
expr := promql.Sum(promql.Rate(promql.Vector("requests_total"), "5m"), "service")
// Serializes to: sum by (service) (rate(requests_total[5m]))
```

### Dashboard Variables

For Grafana dashboard variables:

```go
// Use $__rate_interval for dashboard-aware intervals
promql.Rate(vector, "$__rate_interval")
```

---

## Dual Output Mode

wetwire-observability supports two output modes.

### Standalone Mode

Generates traditional configuration files:

```bash
wetwire-obs build . --mode=standalone
```

Output:
- `prometheus.yml` - Prometheus configuration
- `alertmanager.yml` - Alertmanager configuration
- `rules/*.yml` - Rule files
- `dashboards/*.json` - Grafana dashboards

### Operator Mode

Generates Prometheus Operator CRDs:

```bash
wetwire-obs build . --mode=operator
```

Output:
- `manifests/servicemonitor-*.yaml` - ServiceMonitor resources
- `manifests/prometheusrule-*.yaml` - PrometheusRule resources
- `manifests/alertmanagerconfig-*.yaml` - AlertmanagerConfig resources

### Implementation

```go
switch mode {
case Standalone:
    return serializeStandalone(resources)
case Operator:
    return serializeOperator(resources)
case Both:
    standalone := serializeStandalone(resources)
    operator := serializeOperator(resources)
    return merge(standalone, operator)
}
```

---

## Row-Based Dashboard Layout

Grafana panels are auto-positioned from row definitions:

```go
var Dashboard = grafana.Dashboard{
    Title: "API Metrics",
    Rows: []grafana.Row{
        {Panels: []any{Panel1, Panel2}},  // Side by side, y=0
        {Panels: []any{Panel3}},           // Full width, y=8
    },
}
```

The serializer calculates x/y positions automatically based on row index and panel count.

---

## Linter Architecture

The linter checks Go source for style issues and potential problems.

### Rule Structure

Each rule has:
- **ID**: `WOB001` through `WOB219`
- **Severity**: error, warning, or info
- **Check function**: Analyzes discovered resources

### Rule Categories

See [LINT_RULES.md](LINT_RULES.md) for the complete rule reference with categories WOB001-WOB219.

### Key Rules

| ID | Description |
|----|-------------|
| WOB020 | Scrape config must have job_name |
| WOB050 | Route must have receiver |
| WOB080 | Alert must have severity label |
| WOB100 | Prefer typed PromQL over strings |
| WOB120 | Dashboard must have title |
| WOB200 | No secrets in config |

### Running the Linter

```go
import "github.com/lex00/wetwire-observability-go/internal/lint"

results := lint.Lint(resources)

for _, r := range results {
    fmt.Printf("%s:%d [%s] %s\n", r.File, r.Line, r.Rule, r.Message)
}
```

---

## Files Reference

| File | Purpose |
|------|---------|
| `prometheus/config.go` | Prometheus config types |
| `alertmanager/config.go` | Alertmanager config types |
| `rules/rules.go` | AlertingRule, RecordingRule |
| `grafana/dashboard.go` | Dashboard, Panel types |
| `promql/promql.go` | PromQL expression builders |
| `operator/types.go` | Prometheus Operator CRD types |
| `internal/discover/` | AST-based discovery |
| `internal/serialize/` | Config serialization |
| `internal/lint/` | Lint rules |
| `internal/importer/` | Config importers |

---

## See Also

- [Developer Guide](DEVELOPERS.md) - Development workflow
- [Lint Rules](LINT_RULES.md) - Complete rule reference
- [CLI Reference](CLI.md) - CLI commands
