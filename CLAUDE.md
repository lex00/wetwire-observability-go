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

Uses the `WOB` prefix (Wetwire OBservability). See [LINT_RULES.md](docs/LINT_RULES.md) for the complete rule reference with categories WOB001-WOB219.

## CLI Commands

```bash
wetwire-obs build .              # Generate all outputs
wetwire-obs lint .               # Check patterns
wetwire-obs import prometheus.yml # Convert to Go
wetwire-obs validate .           # Run external validators
wetwire-obs list .               # Show resources
wetwire-obs mcp                  # Start MCP server
```


## Diff

Compare Prometheus/Alertmanager configs semantically:

```bash
# Compare two files
wetwire-obs diff file1 file2

# JSON output for CI/CD
wetwire-obs diff file1 file2 -f json

# Ignore array ordering differences
wetwire-obs diff file1 file2 --ignore-order
```

The diff command performs semantic comparison by resource name, detecting:
- Added resources
- Removed resources
- Modified resources (with property-level change details)

Exit code is 1 if differences are found, enabling CI pipeline validation.

## Related Documentation

- [wetwire spec](https://github.com/lex00/wetwire/docs/WETWIRE_SPEC.md)
- [Feature matrix](https://github.com/lex00/wetwire/docs/FEATURE_MATRIX.md)
