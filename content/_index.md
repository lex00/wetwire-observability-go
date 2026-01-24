---
title: "Wetwire Observability"
---

[![Go Reference](https://pkg.go.dev/badge/github.com/lex00/wetwire-observability-go.svg)](https://pkg.go.dev/github.com/lex00/wetwire-observability-go)
[![CI](https://github.com/lex00/wetwire-observability-go/actions/workflows/ci.yml/badge.svg)](https://github.com/lex00/wetwire-observability-go/actions/workflows/ci.yml)
[![codecov](https://codecov.io/gh/lex00/wetwire-observability-go/graph/badge.svg)](https://codecov.io/gh/lex00/wetwire-observability-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

Generate multi-backend observability configurations from Go structs with AI-assisted design.

## Philosophy

Wetwire uses typed constraints to reduce the model capability required for accurate code generation.

**Core hypothesis:** Typed input + smaller model ≈ Semantic input + larger model

The type system and lint rules act as a force multiplier — cheaper models produce quality output when guided by schema-generated types and iterative lint feedback.

## Documentation

| Document | Description |
|----------|-------------|
| [CLI Reference]({{< relref "/cli" >}}) | Command-line interface |
| [Quick Start]({{< relref "/quick-start" >}}) | Get started in 5 minutes |
| [Examples]({{< relref "/examples" >}}) | Sample observability projects |
| [FAQ]({{< relref "/faq" >}}) | Frequently asked questions |

## Installation

```bash
go install github.com/lex00/wetwire-observability-go@latest
```

## Quick Example

```go
var LatencyDashboard = dashboard.Dashboard{
    Title: "Service Latency",
    Panels: []dashboard.Panel{P99Panel, ErrorRatePanel},
}
```
