# CLAUDE.md

This file provides context for AI assistants working with this repository.

## Overview

wetwire-observability-go is a unified wetwire domain package for Prometheus, Alertmanager, and Grafana configuration synthesis. It generates both standalone configs and Prometheus Operator CRDs from typed Go structs.

## Key Concepts

### Shared PromQL Types

The `promql` package provides expression builders used across both alerting rules and Grafana dashboard queries:

```go
// Same expression works in alerts and dashboards
var expr = promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service")

// In alert
var alert = rules.AlertingRule{Expr: expr}

// In dashboard
var target = grafana.PrometheusTarget{Expr: expr}
```

### Dual Output Mode

```bash
wetwire-obs build . --mode=standalone  # prometheus.yml, alertmanager.yml, dashboards/*.json
wetwire-obs build . --mode=operator    # ServiceMonitor, PrometheusRule, etc.
wetwire-obs build . --mode=both        # All outputs
```

### Row-Based Dashboard Layout

Grafana panels are auto-positioned from row definitions:

```go
var Dashboard = grafana.Dashboard{
    Rows: []grafana.Row{
        {Panels: []any{Panel1, Panel2}},  // Side by side
        {Panels: []any{Panel3}},           // Full width below
    },
}
```

## Package Structure

```
prometheus/     # Prometheus config types (ScrapeConfig, GlobalConfig)
alertmanager/   # Alertmanager config types (Route, Receiver)
rules/          # AlertingRule, RecordingRule, RuleGroup
grafana/        # Dashboard, Panel types, Targets
promql/         # Shared PromQL expression builders
operator/       # Prometheus Operator CRD types
internal/       # discover, serialize, lint, importer
```

## Lint Rules

Prefix: `WOB` (Wetwire OBservability)

| Range | Category |
|-------|----------|
| WOB001-019 | Core wetwire patterns |
| WOB020-049 | Prometheus config |
| WOB050-079 | Alertmanager |
| WOB080-099 | Alerting/recording rules |
| WOB100-119 | PromQL patterns |
| WOB120-149 | Grafana dashboards |
| WOB200-219 | Security |

## CLI Commands

```bash
wetwire-obs build .              # Generate all outputs
wetwire-obs lint .               # Check patterns
wetwire-obs import prometheus.yml # Convert to Go
wetwire-obs validate .           # Run external validators
wetwire-obs list .               # Show resources
wetwire-obs mcp                  # Start MCP server
```

## Related Documentation

- [wetwire spec](https://github.com/lex00/wetwire/docs/WETWIRE_SPEC.md)
- [Feature matrix](https://github.com/lex00/wetwire/docs/FEATURE_MATRIX.md)
