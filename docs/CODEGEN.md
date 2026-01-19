# Code Generation

This document describes how wetwire-observability-go generates configuration files from Go declarations.

---

## Overview

wetwire-observability uses hand-crafted Go types that mirror Prometheus, Alertmanager, and Grafana configurations. The code generation pipeline converts these typed Go structs into the appropriate output formats (YAML for Prometheus/Alertmanager, JSON for Grafana).

---

## Directory Structure

```
wetwire-observability-go/
├── prometheus/           # Prometheus config types
│   ├── config.go        # PrometheusConfig, ScrapeConfig
│   ├── scrape.go        # ScrapeConfig details
│   └── remote.go        # RemoteWrite, RemoteRead
│
├── alertmanager/        # Alertmanager config types
│   ├── config.go        # AlertmanagerConfig
│   ├── route.go         # Route, routing tree
│   └── receiver.go      # Receiver types
│
├── rules/               # Alerting and recording rules
│   ├── rules.go         # AlertingRule, RecordingRule
│   └── group.go         # RuleGroup
│
├── grafana/             # Grafana dashboard types
│   ├── dashboard.go     # Dashboard
│   ├── panel.go         # Panel types (Stat, Graph, etc.)
│   └── target.go        # PrometheusTarget, etc.
│
├── promql/              # Shared PromQL builders
│   ├── promql.go        # Expression types
│   ├── functions.go     # Rate, Sum, Avg, etc.
│   └── operators.go     # GT, LT, And, Or, etc.
│
├── operator/            # Prometheus Operator CRDs
│   ├── servicemonitor.go
│   └── prometheusrule.go
│
└── internal/
    ├── discover/        # AST-based discovery
    └── serialize/       # YAML/JSON serialization
```

---

## Generation Pipeline

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

### Stage 1: Discover

The discovery phase uses Go's AST package:

```go
import "github.com/lex00/wetwire-observability-go/internal/discover"

resources, err := discover.DiscoverAll("./monitoring/...")
```

### Stage 2: Serialize

The serialization phase converts Go structs to output formats:

```go
import "github.com/lex00/wetwire-observability-go/internal/serialize"

// Prometheus YAML
yaml, err := serialize.PrometheusConfig(config)

// Grafana JSON
json, err := serialize.GrafanaDashboard(dashboard)
```

### Stage 3: Output

Write configuration files:

```bash
wetwire-obs build . --mode=standalone -o ./output/
# Creates:
#   output/prometheus.yml
#   output/alertmanager.yml
#   output/rules/*.yml
#   output/dashboards/*.json
```

---

## Type Mapping

### Prometheus

| Go Type | YAML Output |
|---------|-------------|
| `prometheus.Config` | Root prometheus.yml |
| `prometheus.ScrapeConfig` | `scrape_configs` entry |
| `prometheus.StaticConfig` | `static_configs` entry |
| `prometheus.RemoteWrite` | `remote_write` entry |

### Alertmanager

| Go Type | YAML Output |
|---------|-------------|
| `alertmanager.Config` | Root alertmanager.yml |
| `alertmanager.Route` | `route` section |
| `alertmanager.Receiver` | `receivers` entry |

### Rules

| Go Type | YAML Output |
|---------|-------------|
| `rules.AlertingRule` | Alert in rule group |
| `rules.RecordingRule` | Recording rule in group |
| `rules.RuleGroup` | Rule group file |

### Grafana

| Go Type | JSON Output |
|---------|-------------|
| `grafana.Dashboard` | Dashboard JSON |
| `grafana.StatPanel` | Panel with type "stat" |
| `grafana.GraphPanel` | Panel with type "graph" |
| `grafana.Row` | Row container |

---

## PromQL Serialization

PromQL expressions serialize to strings:

```go
// Go
var expr = promql.Sum(
    promql.Rate(promql.Vector("http_requests_total"), "5m"),
    "service",
)

// Output string
// sum by (service) (rate(http_requests_total[5m]))
```

### Expression Types

| Go Builder | PromQL Output |
|------------|---------------|
| `promql.Vector("metric")` | `metric` |
| `promql.Rate(v, "5m")` | `rate(v[5m])` |
| `promql.Sum(e, "label")` | `sum by (label) (e)` |
| `promql.GT(e, promql.Scalar(0.5))` | `e > 0.5` |

---

## Dual Mode Output

### Standalone Mode

Generates traditional configuration files:

```go
// Input
var APIScrape = prometheus.ScrapeConfig{
    JobName: "api",
    StaticConfigs: []prometheus.StaticConfig{
        {Targets: []string{"localhost:8080"}},
    },
}
```

```yaml
# Output: prometheus.yml
scrape_configs:
  - job_name: api
    static_configs:
      - targets:
          - localhost:8080
```

### Operator Mode

Generates Kubernetes CRDs:

```go
// Same input generates ServiceMonitor
```

```yaml
# Output: servicemonitor-api.yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: api
spec:
  endpoints:
    - port: http
  selector:
    matchLabels:
      app: api
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

## Validation

After generation, verify the output:

```bash
# Check syntax
wetwire-obs lint ./monitoring/...

# Build and review
wetwire-obs build ./monitoring/... -o ./output/

# Validate with external tools
promtool check config output/prometheus.yml
amtool check-config output/alertmanager.yml
```

---

## See Also

- [Developer Guide](DEVELOPERS.md) - Development workflow
- [Internals](INTERNALS.md) - Architecture details
- [CLI Reference](CLI.md) - Build command options
