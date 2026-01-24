---
title: "Lint Rules"
---

wetwire-observability-go includes lint rules (WOB prefix) to enforce best practices for Prometheus, Alertmanager, and Grafana configuration.

**Note:** Some rules are planned but not yet implemented. Run `wetwire-obs lint --help` to see currently active rules.

## Quick Start

```bash
# Lint your monitoring code
wetwire-obs lint ./monitoring

# Output in JSON format
wetwire-obs lint ./monitoring -f json
```

## Rule Index

| Rule | Description | Severity | Category |
|------|-------------|----------|----------|
| WOB001 | Use typed structs instead of map literals | warning | Core |
| WOB002 | Extract inline structs to named variables | warning | Core |
| WOB003 | Detect duplicate resource names | error | Core |
| WOB004 | Split large files | warning | Core |
| WOB020 | Use Duration type for intervals | warning | Prometheus |
| WOB021 | Validate scrape interval bounds | warning | Prometheus |
| WOB022 | Require job_name in ScrapeConfig | error | Prometheus |
| WOB050 | Validate receiver names | warning | Alertmanager |
| WOB051 | Require default receiver | error | Alertmanager |
| WOB080 | Require alert name | error | Rules |
| WOB081 | Require for duration on alerts | warning | Rules |
| WOB082 | Require severity label | warning | Rules |
| WOB100 | Use promql builders | warning | PromQL |
| WOB101 | Validate PromQL syntax | error | PromQL |
| WOB120 | Require dashboard title | error | Grafana |
| WOB121 | Use row-based layout | warning | Grafana |
| WOB200 | Detect hardcoded secrets | error | Security |

## Rule Categories

### Core Rules (WOB001-019)

General wetwire patterns that apply across all resource types.

### Prometheus Rules (WOB020-049)

Rules specific to Prometheus configuration.

### Alertmanager Rules (WOB050-079)

Rules specific to Alertmanager configuration.

### Alerting/Recording Rules (WOB080-099)

Rules for alert and recording rule definitions.

### PromQL Rules (WOB100-119)

Rules for PromQL expression patterns.

### Grafana Rules (WOB120-149)

Rules for Grafana dashboard definitions.

### Security Rules (WOB200-219)

Security-focused rules for detecting sensitive data.

---

## Rule Details

### WOB001: Use Typed Structs

**Description:** Use typed structs instead of `map[string]any`.

**Severity:** warning

#### Bad

```go
var Config = map[string]any{
    "global": map[string]any{
        "scrape_interval": "15s",
    },
}
```

#### Good

```go
var Config = prometheus.PrometheusConfig{
    Global: &prometheus.GlobalConfig{
        ScrapeInterval: prometheus.Duration("15s"),
    },
}
```

---

### WOB002: Extract Inline Structs

**Description:** Extract inline structs to named variables.

**Severity:** warning

#### Bad

```go
var Config = prometheus.PrometheusConfig{
    ScrapeConfigs: []*prometheus.ScrapeConfig{
        {JobName: "api", StaticConfigs: []*prometheus.StaticConfig{{Targets: []string{"api:8080"}}}},
    },
}
```

#### Good

```go
var APITargets = prometheus.StaticConfig{
    Targets: []string{"api:8080"},
}

var APIScrape = prometheus.ScrapeConfig{
    JobName:       "api",
    StaticConfigs: []*prometheus.StaticConfig{&APITargets},
}

var Config = prometheus.PrometheusConfig{
    ScrapeConfigs: []*prometheus.ScrapeConfig{&APIScrape},
}
```

---

### WOB020: Use Duration Type

**Description:** Use the `prometheus.Duration` type for time intervals.

**Severity:** warning

#### Bad

```go
ScrapeInterval: "15s",
```

#### Good

```go
ScrapeInterval: prometheus.Duration("15s"),
```

---

### WOB021: Validate Scrape Interval

**Description:** Scrape intervals should be between 5s and 5m.

**Severity:** warning

Very short intervals can overload targets; very long intervals may miss data.

#### Bad

```go
ScrapeInterval: prometheus.Duration("1s"),   // Too frequent
ScrapeInterval: prometheus.Duration("10m"),  // Too infrequent
```

#### Good

```go
ScrapeInterval: prometheus.Duration("15s"),
```

---

### WOB080: Require Alert Name

**Description:** Alerting rules must have an alert name.

**Severity:** error

#### Bad

```go
var MyAlert = rules.AlertingRule{
    Expr: promql.GT(promql.Vector("up"), promql.Scalar(0)),
}
```

#### Good

```go
var MyAlert = rules.AlertingRule{
    Alert: "TargetDown",
    Expr:  promql.GT(promql.Vector("up"), promql.Scalar(0)),
}
```

---

### WOB081: Require For Duration

**Description:** Alerting rules should have a `For` duration to avoid flapping.

**Severity:** warning

#### Bad

```go
var MyAlert = rules.AlertingRule{
    Alert: "HighLatency",
    Expr:  latencyExpr,
    // No For duration - will fire immediately
}
```

#### Good

```go
var MyAlert = rules.AlertingRule{
    Alert: "HighLatency",
    Expr:  latencyExpr,
    For:   5 * time.Minute,
}
```

---

### WOB082: Require Severity Label

**Description:** Alerting rules should have a severity label.

**Severity:** warning

#### Bad

```go
var MyAlert = rules.AlertingRule{
    Alert: "HighLatency",
    Expr:  latencyExpr,
}
```

#### Good

```go
var MyAlert = rules.AlertingRule{
    Alert:    "HighLatency",
    Expr:     latencyExpr,
    Severity: rules.Critical,  // or rules.Warning, rules.Info
}
```

---

### WOB100: Use PromQL Builders

**Description:** Use promql package builders instead of raw strings.

**Severity:** warning

#### Bad

```go
Expr: "sum(rate(http_requests_total[5m])) by (service)",
```

#### Good

```go
Expr: promql.Sum(promql.Rate(promql.Vector("http_requests_total"), "5m"), "service"),
```

---

### WOB101: Validate PromQL Syntax

**Description:** PromQL expressions must be syntactically valid.

**Severity:** error

Catches syntax errors before deployment.

---

### WOB120: Require Dashboard Title

**Description:** Dashboards must have a title.

**Severity:** error

#### Bad

```go
var MyDashboard = grafana.Dashboard{
    Rows: []grafana.Row{...},
}
```

#### Good

```go
var MyDashboard = grafana.Dashboard{
    Title: "API Service Metrics",
    Rows:  []grafana.Row{...},
}
```

---

### WOB121: Use Row-Based Layout

**Description:** Use row-based layout instead of manual GridPos.

**Severity:** warning

Row-based layout ensures consistent panel positioning.

#### Bad

```go
var Dashboard = grafana.Dashboard{
    Panels: []any{
        grafana.TimeseriesPanel{GridPos: grafana.GridPos{X: 0, Y: 0, W: 12, H: 8}},
        grafana.TimeseriesPanel{GridPos: grafana.GridPos{X: 12, Y: 0, W: 12, H: 8}},
    },
}
```

#### Good

```go
var Dashboard = grafana.Dashboard{
    Rows: []grafana.Row{
        {Panels: []any{Panel1, Panel2}},
    },
}
```

---

### WOB200: Detect Hardcoded Secrets

**Description:** Detect hardcoded secrets, API keys, and credentials.

**Severity:** error

**Detected patterns:**
- API keys and tokens in string literals
- Passwords in configuration fields
- Webhook URLs with embedded tokens

#### Bad

```go
var SlackConfig = alertmanager.SlackConfig{
    APIURL: "https://hooks.slack.com/services/T00/B00/XXXXXXXXXX",
}
```

#### Good

```go
var SlackConfig = alertmanager.SlackConfig{
    APIURL: alertmanager.SecretRef("slack-webhook-url"),
}
```

---

## See Also

- [CLI Reference](CLI.md)
- [Quick Start](QUICK_START.md)
- [FAQ](FAQ.md)
