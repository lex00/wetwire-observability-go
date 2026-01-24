<picture>
  <source media="(prefers-color-scheme: dark)" srcset="docs/wetwire-dark.svg">
  <img src="docs/wetwire-light.svg" width="100" height="67" align="right">
</picture>


# wetwire-observability (Go)

[![CI](https://github.com/lex00/wetwire-observability-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-observability-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-observability-go/branch/main/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-observability-go)
[![Go](https://img.shields.io/badge/Go-1.23-blue?logo=go)](https://golang.org/)
[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-observability-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-observability-go)
[![Go Report Card](https://goreportcard.com/badge/github.com/lex00/wetwire-observability-go)](https://goreportcard.com/report/github.com/lex00/wetwire-observability-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Prometheus, Alertmanager, and Grafana configuration synthesis using Go struct literals.

## Installation

```bash
go install github.com/lex00/wetwire-observability-go/cmd/wetwire-obs@latest
```

## Quick Example

```go
package monitoring

import (
    "time"
    "github.com/lex00/wetwire-observability-go/prometheus"
    "github.com/lex00/wetwire-observability-go/alertmanager"
    "github.com/lex00/wetwire-observability-go/rules"
    "github.com/lex00/wetwire-observability-go/grafana"
    "github.com/lex00/wetwire-observability-go/promql"
)

// Shared PromQL expression used in both alert and dashboard
var ErrorRateExpr = promql.GT(
    promql.Div(
        promql.Sum(promql.Rate(promql.Vector("http_requests_total",
            promql.Match("status", "5..")), "$__rate_interval"), "service"),
        promql.Sum(promql.Rate(promql.Vector("http_requests_total"),
            "$__rate_interval"), "service"),
    ),
    promql.Scalar(0.05),
)

// Alert using shared expression
var HighErrorRate = rules.AlertingRule{
    Alert:    "HighErrorRate",
    Expr:     ErrorRateExpr,
    For:      5 * time.Minute,
    Severity: rules.Critical,
}

// Dashboard panel using same expression
var ErrorRatePanel = grafana.StatPanel{
    Title:   "Error Rate",
    Targets: []any{grafana.PrometheusTarget{RefID: "A", Expr: ErrorRateExpr}},
}
```

```bash
# Generate standalone configs
wetwire-obs build . --mode=standalone
# Output: prometheus.yml, alertmanager.yml, rules/*.yml, dashboards/*.json

# Generate Prometheus Operator CRDs
wetwire-obs build . --mode=operator
# Output: manifests/*.yaml (ServiceMonitor, PrometheusRule, etc.)
```

## Features

- **Unified observability stack** - Prometheus, Alertmanager, and Grafana in one package
- **Shared PromQL types** - Same expression builders for alerts and dashboards
- **Dual output mode** - Standalone configs or Prometheus Operator CRDs
- **Row-based layout** - Auto-positioned Grafana panels
- **Type-safe references** - Direct variable references, IDE autocomplete
- **Lint enforcement** - WOB rules ensure consistent patterns

## AI-Assisted Design

Create observability configuration interactively with AI:

```bash
# No API key required - uses Claude CLI
wetwire-obs design "Add monitoring for my API service"

# Automated testing with personas
wetwire-obs test --persona beginner "Create error rate dashboard"
```

Uses [Claude CLI](https://claude.ai/download) by default (no API key required). Falls back to Anthropic API if Claude CLI is not installed. See [CLI Reference](docs/CLI.md#design) for details.

## Documentation

**Getting Started:**
- [Quick Start](docs/QUICK_START.md) - 5-minute tutorial
- [FAQ](docs/FAQ.md) - Common questions

**Reference:**
- [CLI Reference](docs/CLI.md) - All commands including design, test, import
- [Lint Rules](docs/LINT_RULES.md) - WOB rule reference

**Advanced:**
- [Internals](docs/INTERNALS.md) - Architecture and extension points
- [Adoption Guide](docs/ADOPTION.md) - Team migration strategies
- [Import Workflow](docs/IMPORT_WORKFLOW.md) - Migrate existing configs

## Development

```bash
git clone https://github.com/lex00/wetwire-observability-go.git
cd wetwire-observability-go
go mod download
go test ./...
```

## License

MIT - See [LICENSE](LICENSE) for details.
