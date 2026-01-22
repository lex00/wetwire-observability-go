<picture>
  <source media="(prefers-color-scheme: dark)" srcset="./wetwire-dark.svg">
  <img src="./wetwire-light.svg" width="100" height="67">
</picture>

Comprehensive guide for developers working on wetwire-observability-go.

## Table of Contents

- [Development Setup](#development-setup)
- [Project Structure](#project-structure)
- [Running Tests](#running-tests)
- [Adding Features](#adding-features)
- [Contributing](#contributing)

---

## Development Setup

### Prerequisites

- **Go 1.23+** (required)
- **git** (version control)
- **golangci-lint** (optional, for linting)

### Clone and Setup

```bash
# Clone repository
git clone https://github.com/lex00/wetwire-observability-go.git
cd wetwire-observability-go

# Download dependencies
go mod download

# Build CLI
go build -o wetwire-obs ./cmd/wetwire-obs

# Verify installation
./wetwire-obs version
```

### Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test -v ./internal/lint/...
```

---

## Project Structure

```
wetwire-observability-go/
├── cmd/wetwire-obs/         # CLI application
│   ├── main.go              # Entry point, command registration
│   ├── build.go             # build command
│   ├── lint.go              # lint command
│   ├── import.go            # import command
│   ├── validate.go          # validate command
│   ├── design.go            # design command (AI-assisted)
│   ├── test.go              # test command (persona testing)
│   └── mcp.go               # MCP server
│
├── prometheus/              # Prometheus config types
│   ├── config.go            # PrometheusConfig
│   ├── scrape.go            # ScrapeConfig
│   └── remote.go            # RemoteWrite, RemoteRead
│
├── alertmanager/            # Alertmanager config types
│   ├── config.go            # AlertmanagerConfig
│   ├── route.go             # Route
│   └── receiver.go          # Receiver types
│
├── rules/                   # Alerting and recording rules
│   ├── rules.go             # AlertingRule, RecordingRule
│   └── group.go             # RuleGroup
│
├── grafana/                 # Grafana dashboard types
│   ├── dashboard.go         # Dashboard
│   ├── panel.go             # Panel types
│   ├── row.go               # Row layout
│   └── target.go            # Data source targets
│
├── promql/                  # Shared PromQL builders
│   ├── promql.go            # Expression types
│   ├── functions.go         # Functions (Rate, Sum, etc.)
│   └── operators.go         # Operators (GT, LT, etc.)
│
├── operator/                # Prometheus Operator CRDs
│   ├── servicemonitor.go    # ServiceMonitor
│   └── prometheusrule.go    # PrometheusRule
│
├── internal/
│   ├── discover/            # AST-based resource discovery
│   ├── serialize/           # YAML/JSON serialization
│   ├── lint/                # Lint rules (WOB001-WOB219)
│   ├── importer/            # Config importers
│   └── builder/             # Build pipeline
│
├── examples/                # Example configurations
├── testdata/                # Test fixtures
└── docs/                    # Documentation
```

---

## Running Tests

```bash
# Run all tests
go test -v ./...

# Run with coverage
go test -cover ./...

# Run specific test
go test -v ./internal/lint/... -run TestWOB020

# Run with race detection
go test -race ./...

# Generate coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Test Organization

- **Unit tests**: `*_test.go` files next to source
- **Integration tests**: `internal/*/integration_test.go`
- **E2E tests**: `cmd/wetwire-obs/e2e_test.go`

---

## Adding Features

### Adding a New Panel Type

1. Add the type definition in `grafana/panel.go`:

```go
// HeatmapPanel represents a Grafana heatmap panel.
type HeatmapPanel struct {
    BasePanel
    ColorMode string
    // Heatmap-specific options...
}
```

2. Add serialization in `internal/serialize/grafana.go`:

```go
func serializeHeatmapPanel(p *grafana.HeatmapPanel) map[string]any {
    // Convert to Grafana JSON format
}
```

3. Add discovery support in `internal/discover/grafana.go`

4. Add tests

### Adding a New Lint Rule

1. Choose rule code in appropriate range (WOB001-WOB219)

2. Add rule in `internal/lint/rules.go`:

```go
// WOB125 checks for dashboard panels without titles.
func (l *Linter) checkWOB125(panel *grafana.Panel) []LintResult {
    if panel.Title == "" {
        return []LintResult{{
            Rule:     "WOB125",
            Severity: "warning",
            Message:  "Panel should have a title",
        }}
    }
    return nil
}
```

3. Register in the linter

4. Add tests in `internal/lint/rules_test.go`

5. Document in `docs/LINT_RULES.md`

### Adding a New PromQL Function

1. Add builder in `promql/functions.go`:

```go
// Clamp_max returns an expression that clamps values to a maximum.
func Clamp_max(expr Expr, max float64) Expr {
    return &clampMaxExpr{expr: expr, max: max}
}

type clampMaxExpr struct {
    expr Expr
    max  float64
}

func (e *clampMaxExpr) String() string {
    return fmt.Sprintf("clamp_max(%s, %g)", e.expr.String(), e.max)
}
```

2. Add tests in `promql/functions_test.go`

---

## Code Style

### Go Formatting

- Use `gofmt` for all Go code
- Use `goimports` to organize imports
- Follow [Effective Go](https://golang.org/doc/effective_go) guidelines

```bash
# Format all files
gofmt -w .

# Check formatting
gofmt -d .
```

### Error Handling

```go
// Good: Return wrapped errors with context
if err != nil {
    return nil, fmt.Errorf("failed to parse config: %w", err)
}

// Good: Use sentinel errors for expected conditions
var ErrInvalidConfig = errors.New("invalid configuration")
```

### Documentation

- Document all exported types and functions
- Include examples in doc comments where helpful

```go
// AlertingRule represents a Prometheus alerting rule.
// It defines conditions that trigger alerts when met.
type AlertingRule struct {
    // Alert is the name of the alert.
    Alert string

    // Expr is the PromQL expression to evaluate.
    Expr Expr

    // For is the duration to wait before firing.
    For time.Duration
    // ...
}
```

---

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for:
- Code style guidelines
- Commit message format
- Pull request process
- Adding new features

---

## Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/spf13/cobra` | CLI framework |
| `gopkg.in/yaml.v3` | YAML parsing/generation |
| `github.com/stretchr/testify` | Test assertions |
| `github.com/lex00/wetwire-core-go` | AI orchestration |
| `github.com/modelcontextprotocol/go-sdk` | MCP server |

---

## See Also

- [Quick Start](QUICK_START.md) - Getting started
- [CLI Reference](CLI.md) - CLI commands
- [Internals](INTERNALS.md) - Architecture details
- [Lint Rules](LINT_RULES.md) - All WOB rules
