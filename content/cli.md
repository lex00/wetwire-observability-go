---
title: "Cli"
---
<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

The `wetwire-obs` command generates Prometheus, Alertmanager, and Grafana configurations from Go code.

## Quick Reference

| Command | Description |
|---------|-------------|
| `wetwire-obs build` | Generate configuration files from Go source |
| `wetwire-obs lint` | Lint code for issues |
| `wetwire-obs init` | Initialize a new project |
| `wetwire-obs import` | Convert existing configs to Go code |
| `wetwire-obs validate` | Validate resources |
| `wetwire-obs list` | List discovered resources |
| `wetwire-obs design` | AI-assisted config design |
| `wetwire-obs test` | Test with simulated personas |
| `wetwire-obs mcp` | Start MCP server |

```bash
wetwire-obs --help     # Show help
```

---

## build

Generate Prometheus/Alertmanager/Grafana configuration from Go source files.

```bash
# Generate JSON to stdout
wetwire-obs build ./monitoring

# Generate with pretty formatting
wetwire-obs build ./monitoring --format pretty

# Output to file
wetwire-obs build ./monitoring -o config.json

# Build only specific resource type
wetwire-obs build ./monitoring --type prometheus
```

### Options

| Option | Description |
|--------|-------------|
| `PATH` | Directory containing Go source files |
| `--format, -f {json,pretty}` | Output format (default: json) |
| `--output, -o FILE` | Output file (default: stdout) |
| `--type {prometheus,alertmanager,rules,grafana}` | Filter by resource type |
| `--dry-run` | Show what would be generated without writing |

### Output Modes

The `--mode` flag controls output format:

```bash
# Standalone configs (default)
wetwire-obs build . --mode=standalone
# Output: prometheus.yml, alertmanager.yml, rules/*.yml, dashboards/*.json

# Prometheus Operator CRDs
wetwire-obs build . --mode=operator
# Output: manifests/*.yaml (ServiceMonitor, PrometheusRule, etc.)

# Both formats
wetwire-obs build . --mode=both
```

### How It Works

1. Parses Go source files using `go/ast`
2. Discovers resource declarations (PrometheusConfig, ScrapeConfig, AlertingRule, Dashboard, etc.)
3. Extracts resource dependencies
4. Serializes to output format (YAML for configs, JSON for dashboards)

---

## lint

Lint wetwire-obs code for issues.

```bash
# Lint a directory
wetwire-obs lint ./monitoring

# Lint a single file
wetwire-obs lint ./monitoring/alerts.go
```

### Options

| Option | Description |
|--------|-------------|
| `PATH` | File or directory to lint |

### What It Checks

1. **Resource discovery**: Validates resources can be parsed from source
2. **Reference validity**: Checks that referenced resources exist
3. **PromQL syntax**: Validates PromQL expressions
4. **Best practices**: Enforces observability patterns

### Output Examples

**Linting passed:**
```
Linted 12 resources: no issues found
```

**Issues found:**
```
./monitoring/alerts.go:15: invalid PromQL expression
./monitoring/dashboards.go:23: panel references undefined query
```

### Lint Rules

See [LINT_RULES.md](LINT_RULES.md) for the complete rule reference.

---

## init

Initialize a new wetwire-obs project.

```bash
# Create a new project
wetwire-obs init mymonitoring
```

### Arguments

| Argument | Description |
|----------|-------------|
| `path` | Path for the new project (required) |

### Generated Structure

```
mymonitoring/
└── prometheus.go   # Example Prometheus configuration
```

**prometheus.go:**
```go
package monitoring

import "github.com/lex00/wetwire-observability-go/prometheus"

var Production = prometheus.PrometheusConfig{
    Global: &prometheus.GlobalConfig{
        ScrapeInterval:     prometheus.Duration("15s"),
        EvaluationInterval: prometheus.Duration("15s"),
        ExternalLabels: map[string]string{
            "environment": "production",
        },
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

---

## import

Convert existing Prometheus, Alertmanager, or Grafana configurations to Go code.

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

### Options

| Option | Description |
|--------|-------------|
| `FILE` | Configuration file to import |
| `--output, -o DIR` | Output directory for generated Go files |
| `--package, -p NAME` | Package name for generated code (default: monitoring) |

### Supported Formats

| Format | Description |
|--------|-------------|
| `prometheus.yml` | Prometheus configuration |
| `alertmanager.yml` | Alertmanager configuration |
| `*.json` | Grafana dashboard JSON |
| `rules/*.yml` | Prometheus alert/recording rules |

See [IMPORT_WORKFLOW.md](IMPORT_WORKFLOW.md) for detailed migration workflows.

---

## validate

Validate resources and check dependencies.

```bash
wetwire-obs validate ./monitoring
```

### Checks Performed

- **Reference validity**: All resource references point to defined resources
- **PromQL validity**: All PromQL expressions are syntactically correct
- **Dashboard integrity**: Panel references exist

---

## list

List discovered resources in a package.

```bash
wetwire-obs list ./monitoring
```

### Output

```
Discovered 8 resources:
  prometheus_config  Production         monitoring/prometheus.go
  scrape_config      APIServer          monitoring/scrape.go
  alerting_rule      HighErrorRate      monitoring/alerts.go
  recording_rule     ErrorRatio5m       monitoring/recording.go
  rule_group         APIAlerts          monitoring/alerts.go
```

---

## mcp

Start the MCP (Model Context Protocol) server for AI assistant integration.

```bash
wetwire-obs mcp
```

### Available Tools

| Tool | Description |
|------|-------------|
| `wetwire_build` | Generate configurations |
| `wetwire_lint` | Lint code for issues |
| `wetwire_init` | Initialize a new project |
| `wetwire_list` | List discovered resources |

---

## design

Start an AI-assisted design session to create observability configurations.

```bash
# Interactive design session
wetwire-obs design --provider anthropic "Create monitoring for a REST API"

# Using Kiro provider
wetwire-obs design --provider kiro "Create alerts for payment service"
```

### Options

| Option | Description |
|--------|-------------|
| `PROMPT` | Description of what to create |
| `--provider {anthropic,kiro}` | AI provider to use |
| `--model MODEL` | Model to use (default: claude-sonnet-4) |

See [OBSERVABILITY-KIRO-CLI.md](OBSERVABILITY-KIRO-CLI.md) for Kiro integration details.

---

## test

Test configurations with simulated user personas.

```bash
# Test with expert persona
wetwire-obs test --provider anthropic --persona expert "Create error rate alerts"

# Test with novice persona
wetwire-obs test --provider anthropic --persona novice "Set up monitoring"
```

### Options

| Option | Description |
|--------|-------------|
| `PROMPT` | Test prompt to run |
| `--provider {anthropic,kiro}` | AI provider to use |
| `--persona {expert,novice,adversarial}` | User persona to simulate |

---

## Typical Workflow

### Development

```bash
# Lint before generating
wetwire-obs lint ./monitoring

# Generate configs
wetwire-obs build ./monitoring --mode=standalone

# Preview without writing
wetwire-obs build ./monitoring --dry-run
```

### CI/CD

```bash
#!/bin/bash
# ci.sh

# Lint first
wetwire-obs lint ./monitoring || exit 1

# Generate standalone configs
wetwire-obs build ./monitoring --mode=standalone -o ./generated/

# Or generate Kubernetes manifests
wetwire-obs build ./monitoring --mode=operator -o ./manifests/
```

---

## See Also

- [Quick Start](QUICK_START.md) - Create your first project
- [Lint Rules](LINT_RULES.md) - WOB rule reference
- [FAQ](FAQ.md) - Common questions
