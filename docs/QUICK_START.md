# Quick Start

Get started with `wetwire-obs` in 5 minutes.

## Installation

```bash
go install github.com/lex00/wetwire-observability-go/cmd/wetwire-obs@latest
```

Or add to your project:

```bash
go get github.com/lex00/wetwire-observability-go
```

---

## Your First Project

Create a monitoring configuration:

```
mymonitoring/
├── go.mod
└── monitoring/
    └── prometheus.go
```

**monitoring/prometheus.go:**
```go
package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

// Production is the Prometheus configuration
var Production = prometheus.PrometheusConfig{
    Global: &prometheus.GlobalConfig{
        ScrapeInterval:     prometheus.Duration("15s"),
        EvaluationInterval: prometheus.Duration("15s"),
    },
    ScrapeConfigs: []*prometheus.ScrapeConfig{
        {
            JobName: "prometheus",
            StaticConfigs: []*prometheus.StaticConfig{
                {Targets: []string{"localhost:9090"}},
            },
        },
    },
}
```

**Generate configuration:**
```bash
wetwire-obs build ./monitoring
```

---

## Adding Scrape Configs

Define service discovery and scrape targets:

**monitoring/scrape.go:**
```go
package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

// APIServer scrapes the API service
var APIServer = prometheus.ScrapeConfig{
    JobName:        "api-server",
    ScrapeInterval: prometheus.Duration("10s"),
    StaticConfigs: []*prometheus.StaticConfig{
        {
            Targets: []string{"api:8080"},
            Labels: map[string]string{
                "service": "api",
                "env":     "production",
            },
        },
    },
}

// KubernetesNodes uses Kubernetes service discovery
var KubernetesNodes = prometheus.ScrapeConfig{
    JobName: "kubernetes-nodes",
    KubernetesSDConfigs: []*prometheus.KubernetesSDConfig{
        {Role: "node"},
    },
}
```

---

## Adding Alerting Rules

Define alerts using the `rules` package:

**monitoring/alerts.go:**
```go
package monitoring

import (
    "time"
    "github.com/lex00/wetwire-observability-go/rules"
    "github.com/lex00/wetwire-observability-go/promql"
)

// HighErrorRate alerts when error rate exceeds 5%
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
        "description": "Error rate is {{ $value | humanizePercentage }}",
    },
}

// ErrorRateExpr is the shared PromQL expression
var ErrorRateExpr = promql.Div(
    promql.Sum(promql.Rate(promql.Vector("http_requests_total",
        promql.Match("status", "5..")), "5m"), "service"),
    promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service"),
)
```

---

## Adding Recording Rules

Pre-compute expensive queries:

**monitoring/recording.go:**
```go
package monitoring

import (
    "github.com/lex00/wetwire-observability-go/rules"
    "github.com/lex00/wetwire-observability-go/promql"
)

// ErrorRatio5m is a pre-computed error ratio
var ErrorRatio5m = rules.RecordingRule{
    Record: "service:http_error_ratio:5m",
    Expr:   ErrorRateExpr,
    Labels: map[string]string{
        "aggregation": "5m",
    },
}

// RequestRate5m is a pre-computed request rate
var RequestRate5m = rules.RecordingRule{
    Record: "service:http_requests:rate5m",
    Expr: promql.Sum(
        promql.Rate(promql.Vector("http_requests_total"), "5m"),
        "service",
    ),
}
```

---

## Adding Grafana Dashboards

Define dashboards with auto-positioned panels:

**monitoring/dashboards.go:**
```go
package monitoring

import (
    "github.com/lex00/wetwire-observability-go/grafana"
    "github.com/lex00/wetwire-observability-go/promql"
)

// APIDashboard shows API service metrics
var APIDashboard = grafana.Dashboard{
    Title: "API Service",
    Rows: []grafana.Row{
        {Panels: []any{RequestRatePanel, ErrorRatePanel}},  // Side by side
        {Panels: []any{LatencyPanel}},                       // Full width
    },
}

// RequestRatePanel shows requests per second
var RequestRatePanel = grafana.TimeseriesPanel{
    Title: "Request Rate",
    Targets: []any{
        grafana.PrometheusTarget{
            RefID: "A",
            Expr:  promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "$__rate_interval"), "service"),
        },
    },
}

// ErrorRatePanel shows error percentage
var ErrorRatePanel = grafana.StatPanel{
    Title: "Error Rate",
    Targets: []any{
        grafana.PrometheusTarget{RefID: "A", Expr: ErrorRateExpr},
    },
    Unit: "percentunit",
}

// LatencyPanel shows request latency histogram
var LatencyPanel = grafana.HeatmapPanel{
    Title: "Request Latency",
    Targets: []any{
        grafana.PrometheusTarget{
            RefID: "A",
            Expr:  promql.Sum(promql.Rate(promql.Vector("http_request_duration_seconds_bucket"), "$__rate_interval"), "le"),
        },
    },
}
```

---

## Shared PromQL Expressions

The same PromQL expression works in both alerts and dashboards:

```go
// Define once
var ErrorRateExpr = promql.GT(
    promql.Div(
        promql.Sum(promql.Rate(promql.Vector("http_requests_total",
            promql.Match("status", "5..")), "$__rate_interval"), "service"),
        promql.Sum(promql.Rate(promql.Vector("http_requests_total"),
            "$__rate_interval"), "service"),
    ),
    promql.Scalar(0.05),
)

// Use in alert
var HighErrorRate = rules.AlertingRule{
    Expr: ErrorRateExpr,
}

// Use in dashboard panel
var ErrorRatePanel = grafana.StatPanel{
    Targets: []any{grafana.PrometheusTarget{Expr: ErrorRateExpr}},
}
```

---

## Building Output

### Standalone Configs

```bash
wetwire-obs build ./monitoring --mode=standalone
```

Generates:
- `prometheus.yml` - Main Prometheus config
- `alertmanager.yml` - Alertmanager config
- `rules/*.yml` - Alert and recording rules
- `dashboards/*.json` - Grafana dashboards

### Prometheus Operator CRDs

```bash
wetwire-obs build ./monitoring --mode=operator
```

Generates Kubernetes manifests:
- `ServiceMonitor` - Scrape configs
- `PrometheusRule` - Alert and recording rules
- `GrafanaDashboard` - Dashboard ConfigMaps

---

## Multi-File Organization

Split resources across files by concern:

```
monitoring/
├── prometheus.go   # Global config
├── scrape.go       # Scrape configs
├── alerts.go       # Alerting rules
├── recording.go    # Recording rules
├── dashboards.go   # Grafana dashboards
└── promql.go       # Shared PromQL expressions
```

---

## Next Steps

- See the full [CLI Reference](CLI.md)
- Read the [Lint Rules](LINT_RULES.md) for best practices
- Check the [FAQ](FAQ.md) for common questions
