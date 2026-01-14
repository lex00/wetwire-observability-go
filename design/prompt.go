// Package design provides AI-assisted observability configuration design.
package design

import (
	"fmt"
	"strings"
)

// PromptBuilder builds prompts for AI-assisted design.
type PromptBuilder struct {
	focus   string
	context map[string]string
}

// NewPromptBuilder creates a new PromptBuilder.
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{
		context: make(map[string]string),
	}
}

// ForPrometheus focuses the prompt on Prometheus configuration.
func (pb *PromptBuilder) ForPrometheus() *PromptBuilder {
	pb.focus = "prometheus"
	return pb
}

// ForAlertmanager focuses the prompt on Alertmanager configuration.
func (pb *PromptBuilder) ForAlertmanager() *PromptBuilder {
	pb.focus = "alertmanager"
	return pb
}

// ForGrafana focuses the prompt on Grafana dashboards.
func (pb *PromptBuilder) ForGrafana() *PromptBuilder {
	pb.focus = "grafana"
	return pb
}

// ForRules focuses the prompt on alerting and recording rules.
func (pb *PromptBuilder) ForRules() *PromptBuilder {
	pb.focus = "rules"
	return pb
}

// WithContext adds context to the prompt.
func (pb *PromptBuilder) WithContext(key, value string) *PromptBuilder {
	pb.context[key] = value
	return pb
}

// SystemPrompt returns the system prompt for the AI.
func (pb *PromptBuilder) SystemPrompt() string {
	var parts []string

	// Base system prompt
	parts = append(parts, baseSystemPrompt)

	// Focus-specific additions
	switch pb.focus {
	case "prometheus":
		parts = append(parts, prometheusPrompt)
	case "alertmanager":
		parts = append(parts, alertmanagerPrompt)
	case "grafana":
		parts = append(parts, grafanaPrompt)
	case "rules":
		parts = append(parts, rulesPrompt)
	}

	return strings.Join(parts, "\n\n")
}

// BuildUserPrompt builds the user prompt with context.
func (pb *PromptBuilder) BuildUserPrompt(request string) string {
	var parts []string

	// Add context
	if len(pb.context) > 0 {
		parts = append(parts, "## Context")
		for key, value := range pb.context {
			parts = append(parts, fmt.Sprintf("### %s\n```go\n%s\n```", key, value))
		}
		parts = append(parts, "")
	}

	// Add request
	parts = append(parts, "## Request")
	parts = append(parts, request)

	return strings.Join(parts, "\n")
}

const baseSystemPrompt = `You are an expert in observability configuration using the wetwire pattern.
You help users design Prometheus, Alertmanager, Grafana, and alerting rule configurations.

## The wetwire Pattern

wetwire is a Go-based configuration synthesis approach where:
- Configurations are defined as Go variables using fluent builder APIs
- All declarations are flat (package-level var) for AST discoverability
- References between resources are direct Go references
- The wetwire-obs CLI generates prometheus.yml, alertmanager.yml, dashboards, etc.

## Key Packages

- github.com/lex00/wetwire-observability-go/prometheus - Prometheus configuration
- github.com/lex00/wetwire-observability-go/alertmanager - Alertmanager configuration
- github.com/lex00/wetwire-observability-go/rules - Alerting and recording rules
- github.com/lex00/wetwire-observability-go/grafana - Grafana dashboards
- github.com/lex00/wetwire-observability-go/promql - PromQL expression builder
- github.com/lex00/wetwire-observability-go/operator - Kubernetes Operator CRDs

## Output Format

When generating configurations:
1. Use the fluent builder APIs (NewX().WithY().AddZ())
2. Define each resource as a package-level var
3. Include appropriate imports
4. Follow Go naming conventions (CamelCase for exports)
5. Add brief comments explaining the purpose

Generate Go code that compiles and follows best practices.`

const prometheusPrompt = `## Prometheus Focus

For Prometheus configuration, key types include:

- prometheus.NewConfig() - Main Prometheus configuration
- prometheus.NewScrapeConfig() - Job-level scrape configuration
- prometheus.NewStaticConfig() - Static target configuration
- prometheus.NewKubernetesSD() - Kubernetes service discovery
- prometheus.NewRelabelConfig() - Relabeling rules

### Common Patterns

Scrape configuration:
var APIScrape = prometheus.NewScrapeConfig("api-server").
    WithScheme("https").
    WithMetricsPath("/metrics").
    WithStaticConfigs(
        prometheus.NewStaticConfig().AddTarget("api.example.com:9090"),
    )

Kubernetes service discovery:
var KubeScrape = prometheus.NewScrapeConfig("kubernetes-pods").
    WithKubernetesSD(prometheus.NewKubernetesSD().
        WithRole("pod").
        InNamespaces("production"))

Relabeling:
.AddRelabelConfig(prometheus.NewRelabelConfig().
    WithSourceLabels("__meta_kubernetes_pod_annotation_prometheus_io_scrape").
    WithAction("keep").
    WithRegex("true"))`

const alertmanagerPrompt = `## Alertmanager Focus

For Alertmanager configuration, key types include:

- alertmanager.NewAlertmanagerConfig() - Main configuration
- alertmanager.NewRoute() - Routing tree node
- alertmanager.NewReceiver() - Notification receiver
- alertmanager.SlackReceiver() - Slack notifications
- alertmanager.PagerDutyReceiver() - PagerDuty notifications
- alertmanager.NewInhibitRule() - Alert inhibition
- alertmanager.NewMuteTimeInterval() - Muting schedules

### Common Patterns

Routing tree:
var RootRoute = alertmanager.NewRoute("default").
    WithGroupBy("alertname", "namespace").
    WithGroupWait(30 * alertmanager.Second).
    AddRoute(alertmanager.NewRoute("pagerduty-critical").
        Severity("critical"))

Receivers:
var SlackAlerts = alertmanager.SlackReceiver("slack-alerts", "#alerts").
    WithAPIURL(alertmanager.FromSecret("alertmanager", "slack-url"))

Inhibition:
var CriticalInhibits = alertmanager.NewInhibitRule().
    WithSourceMatcher(alertmanager.Eq("severity", "critical")).
    WithTargetMatcher(alertmanager.Eq("severity", "warning")).
    WithEqual("alertname", "namespace")`

const grafanaPrompt = `## Grafana Focus

For Grafana dashboard configuration, key types include:

- grafana.NewDashboard() - Dashboard container
- grafana.NewRow() - Row of panels
- grafana.NewTimeSeriesPanel() - Time series visualization
- grafana.NewStatPanel() - Single stat display
- grafana.NewTablePanel() - Table visualization
- grafana.NewVariable() - Dashboard variables
- grafana.PrometheusTarget() - Prometheus queries

### Common Patterns

Dashboard structure:
var APIDashboard = grafana.NewDashboard("api-overview", "API Overview").
    WithTags("api", "overview").
    WithTime("now-1h", "now").
    WithRefresh("30s").
    AddRow(grafana.NewRow("Overview").
        AddPanel(RequestRatePanel).
        AddPanel(ErrorRatePanel))

Time series panel:
var RequestRatePanel = grafana.NewTimeSeriesPanel("Request Rate").
    WithDataSource("prometheus").
    AddTarget(grafana.PrometheusTarget{
        Expr: promql.Rate(promql.Vector("http_requests_total"), "$__rate_interval"),
    })

Variables:
var NamespaceVar = grafana.NewQueryVariable("namespace", "Namespace").
    WithDataSource("prometheus").
    WithQuery("label_values(up, namespace)")`

const rulesPrompt = `## Alerting Rules Focus

For alerting and recording rules, key types include:

- rules.NewRuleGroup() - Group of rules
- rules.NewAlertingRule() - Alert definition
- rules.NewRecordingRule() - Recording rule definition
- promql expressions for conditions

### Common Patterns

Alert rule:
var HighErrorRate = rules.NewAlertingRule("HighErrorRate").
    WithExpr("rate(http_errors_total[5m]) / rate(http_requests_total[5m]) > 0.05").
    WithFor(5 * rules.Minute).
    WithSeverity(rules.Critical).
    WithSummary("Error rate above 5%").
    WithDescription("{{ $labels.service }} error rate: {{ $value | printf \"%.2f\" }}%")

Recording rule:
var RequestRate5m = rules.NewRecordingRule("job:http_requests:rate5m").
    WithExpr("sum(rate(http_requests_total[5m])) by (job)")

Rule group:
var APIRules = rules.NewRuleGroup("api.rules").
    WithInterval(rules.Minute).
    AddRule(HighErrorRate).
    AddRule(RequestRate5m)`
