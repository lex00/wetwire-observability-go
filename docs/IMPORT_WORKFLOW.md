# Import Workflow Documentation

This document explains the import workflow used by wetwire-observability-go to convert existing Prometheus, Alertmanager, and Grafana configurations into Go code.

## Overview

The import workflow enables converting existing observability configurations into wetwire Go declarations:

1. **Prometheus configs** - prometheus.yml files
2. **Alertmanager configs** - alertmanager.yml files
3. **Grafana dashboards** - JSON dashboard files
4. **Recording/alerting rules** - rules/*.yml files

## Import Sources

### Prometheus Configuration

Import from prometheus.yml:

```bash
wetwire-obs import prometheus.yml -o prometheus.go
```

Supported sections:
- `global` - Global configuration
- `scrape_configs` - Scrape configurations
- `alerting` - Alertmanager configuration
- `rule_files` - Rule file references
- `remote_write` / `remote_read` - Remote storage

### Alertmanager Configuration

Import from alertmanager.yml:

```bash
wetwire-obs import alertmanager.yml -o alertmanager.go
```

Supported sections:
- `global` - Global settings
- `route` - Routing tree
- `receivers` - Notification receivers
- `inhibit_rules` - Inhibition rules
- `templates` - Template files

### Grafana Dashboards

Import from JSON dashboard files:

```bash
wetwire-obs import dashboard.json -o dashboard.go
```

Supported elements:
- Dashboard metadata (title, tags, etc.)
- Rows and panels
- Variables (template variables)
- Annotations

### Rule Files

Import alerting and recording rules:

```bash
wetwire-obs import rules/alerts.yml -o alerts.go
```

Supported:
- Alerting rules with annotations and labels
- Recording rules
- Rule groups

## Workflow Steps

```
┌─────────────────────────────────────────────────────────────┐
│                    Import Workflow                           │
│                                                              │
│  ┌─────────────┐    ┌─────────────┐    ┌─────────────────┐  │
│  │    PARSE    │ ─▶ │   CONVERT   │ ─▶ │    GENERATE     │  │
│  └─────────────┘    └─────────────┘    └─────────────────┘  │
│                                                              │
│  Parse YAML/JSON    Convert to IR      Generate Go code     │
│  config files       (intermediate)     with wetwire types   │
└─────────────────────────────────────────────────────────────┘
```

### Stage 1: Parse

The parser reads the input source based on file type:
- **YAML files**: prometheus.yml, alertmanager.yml, rules/*.yml
- **JSON files**: Grafana dashboard exports

### Stage 2: Convert

Converts parsed data to intermediate representation, extracting:
- Configuration structure
- PromQL expressions
- Panel definitions
- Routing rules

### Stage 3: Generate

Generates idiomatic Go code using wetwire patterns:

```yaml
# Input: prometheus.yml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'api'
    static_configs:
      - targets: ['localhost:8080']
```

```go
// Output: prometheus.go
package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

var PrometheusConfig = prometheus.Config{
    Global: prometheus.GlobalConfig{
        ScrapeInterval: 15 * time.Second,
    },
    ScrapeConfigs: []prometheus.ScrapeConfig{
        APIScrape,
    },
}

var APIScrape = prometheus.ScrapeConfig{
    JobName: "api",
    StaticConfigs: []prometheus.StaticConfig{
        {Targets: []string{"localhost:8080"}},
    },
}
```

## Usage Examples

### Import Prometheus config

```bash
# Basic import
wetwire-obs import prometheus.yml -o prometheus.go

# With custom package name
wetwire-obs import prometheus.yml --package monitoring -o prometheus.go
```

### Import Alertmanager config

```bash
wetwire-obs import alertmanager.yml -o alertmanager.go
```

### Import Grafana dashboard

```bash
# Import single dashboard
wetwire-obs import dashboard.json -o dashboard.go

# Import all dashboards in directory
wetwire-obs import dashboards/*.json -o dashboards/
```

### Import rule files

```bash
wetwire-obs import rules/*.yml -o rules.go
```

## Output Structure

The importer generates organized Go files following wetwire patterns:

```go
// monitoring.go
package monitoring

import (
    "time"
    "github.com/lex00/wetwire-observability-go/prometheus"
    "github.com/lex00/wetwire-observability-go/rules"
    "github.com/lex00/wetwire-observability-go/promql"
)

// Prometheus configuration
var PrometheusConfig = prometheus.Config{...}

// Scrape configurations
var APIScrape = prometheus.ScrapeConfig{...}

// Alerting rules
var HighErrorRate = rules.AlertingRule{
    Alert:    "HighErrorRate",
    Expr:     promql.GT(ErrorRateExpr, promql.Scalar(0.05)),
    For:      5 * time.Minute,
    Severity: rules.Critical,
}
```

## PromQL Expression Handling

PromQL expressions are converted to typed builders:

```yaml
# Input
expr: sum(rate(http_requests_total{status=~"5.."}[5m])) / sum(rate(http_requests_total[5m])) > 0.05
```

```go
// Output
var ErrorRateExpr = promql.GT(
    promql.Div(
        promql.Sum(promql.Rate(promql.Vector("http_requests_total",
            promql.Match("status", "5..")), "5m")),
        promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m")),
    ),
    promql.Scalar(0.05),
)
```

## Validation

After import, verify the generated code:

```bash
# Check syntax
go build ./...

# Lint for issues
wetwire-obs lint ./...

# Build output and compare
wetwire-obs build ./... -o output/
diff original.yml output/prometheus.yml
```

## Limitations

1. **Complex PromQL** - Some complex expressions may need manual adjustment
2. **Custom Grafana plugins** - Plugin-specific panel options may not be fully supported
3. **Templating** - Some dynamic templating may need manual conversion

## See Also

- [Developer Guide](DEVELOPERS.md) - Development workflow
- [Internals](INTERNALS.md) - Architecture details
- [CLI Reference](CLI.md) - Import command options
