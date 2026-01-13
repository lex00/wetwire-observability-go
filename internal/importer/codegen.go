package importer

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"text/template"
	"time"

	"github.com/lex00/wetwire-observability-go/prometheus"
)

// GenerateGoCode generates Go source code from a PrometheusConfig.
func GenerateGoCode(config *prometheus.PrometheusConfig, packageName string) ([]byte, error) {
	gen := &codeGenerator{
		packageName: packageName,
		config:      config,
	}
	return gen.generate()
}

type codeGenerator struct {
	packageName string
	config      *prometheus.PrometheusConfig
}

func (g *codeGenerator) generate() ([]byte, error) {
	var buf bytes.Buffer

	// Write package and imports
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"github.com/lex00/wetwire-observability-go/prometheus\"\n")
	buf.WriteString(")\n\n")

	// Generate global config if present
	if g.config.Global != nil {
		if err := g.writeGlobalConfig(&buf); err != nil {
			return nil, err
		}
	}

	// Generate scrape configs
	for _, sc := range g.config.ScrapeConfigs {
		if err := g.writeScrapeConfig(&buf, sc); err != nil {
			return nil, err
		}
	}

	// Generate remote write configs
	for i, rw := range g.config.RemoteWrite {
		if err := g.writeRemoteWrite(&buf, rw, i); err != nil {
			return nil, err
		}
	}

	// Generate remote read configs
	for i, rr := range g.config.RemoteRead {
		if err := g.writeRemoteRead(&buf, rr, i); err != nil {
			return nil, err
		}
	}

	// Generate the main config
	if err := g.writeMainConfig(&buf); err != nil {
		return nil, err
	}

	// Format the generated code
	formatted, err := format.Source(buf.Bytes())
	if err != nil {
		// Return unformatted code with error for debugging
		return buf.Bytes(), fmt.Errorf("failed to format generated code: %w", err)
	}

	return formatted, nil
}

func (g *codeGenerator) writeGlobalConfig(buf *bytes.Buffer) error {
	gc := g.config.Global

	buf.WriteString("// GlobalConfig defines global settings for Prometheus.\n")
	buf.WriteString("var GlobalConfig = &prometheus.GlobalConfig{\n")

	if gc.ScrapeInterval != 0 {
		buf.WriteString(fmt.Sprintf("\tScrapeInterval: %s,\n", g.formatDuration(gc.ScrapeInterval)))
	}
	if gc.ScrapeTimeout != 0 {
		buf.WriteString(fmt.Sprintf("\tScrapeTimeout: %s,\n", g.formatDuration(gc.ScrapeTimeout)))
	}
	if gc.EvaluationInterval != 0 {
		buf.WriteString(fmt.Sprintf("\tEvaluationInterval: %s,\n", g.formatDuration(gc.EvaluationInterval)))
	}
	if len(gc.ExternalLabels) > 0 {
		buf.WriteString("\tExternalLabels: map[string]string{\n")
		for k, v := range gc.ExternalLabels {
			buf.WriteString(fmt.Sprintf("\t\t%q: %q,\n", k, v))
		}
		buf.WriteString("\t},\n")
	}

	buf.WriteString("}\n\n")
	return nil
}

func (g *codeGenerator) writeScrapeConfig(buf *bytes.Buffer, sc *prometheus.ScrapeConfig) error {
	varName := g.sanitizeVarName(sc.JobName) + "Scrape"

	buf.WriteString(fmt.Sprintf("// %s configures scraping for the %s job.\n", varName, sc.JobName))
	buf.WriteString(fmt.Sprintf("var %s = prometheus.NewScrapeConfig(%q)", varName, sc.JobName))

	// Chain fluent methods
	if sc.ScrapeInterval != 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithInterval(%s)", g.formatDuration(sc.ScrapeInterval)))
	}
	if sc.ScrapeTimeout != 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithTimeout(%s)", g.formatDuration(sc.ScrapeTimeout)))
	}

	// Static configs
	for _, static := range sc.StaticConfigs {
		if len(static.Targets) > 0 {
			targets := make([]string, len(static.Targets))
			for i, t := range static.Targets {
				targets[i] = fmt.Sprintf("%q", t)
			}
			buf.WriteString(fmt.Sprintf(".\n\tWithStaticTargets(%s)", strings.Join(targets, ", ")))
		}
	}

	// Kubernetes SD
	for _, k8s := range sc.KubernetesSDConfigs {
		buf.WriteString(".\n\tWithKubernetesSD(")
		buf.WriteString(g.formatKubernetesSD(k8s))
		buf.WriteString(")")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *codeGenerator) formatKubernetesSD(sd *prometheus.KubernetesSD) string {
	var buf bytes.Buffer
	buf.WriteString(fmt.Sprintf("prometheus.NewKubernetesSD(prometheus.KubernetesRole%s)", g.capitalizeRole(string(sd.Role))))

	if sd.Namespaces != nil && len(sd.Namespaces.Names) > 0 {
		names := make([]string, len(sd.Namespaces.Names))
		for i, n := range sd.Namespaces.Names {
			names[i] = fmt.Sprintf("%q", n)
		}
		buf.WriteString(fmt.Sprintf(".\n\t\tWithNamespaces(%s)", strings.Join(names, ", ")))
	}

	if sd.Namespaces != nil && sd.Namespaces.OwnNamespace {
		buf.WriteString(".\n\t\tWithOwnNamespace()")
	}

	for _, sel := range sd.Selectors {
		if sel.Label != "" {
			buf.WriteString(fmt.Sprintf(".\n\t\tWithLabelSelector(prometheus.KubernetesRole%s, %q)", g.capitalizeRole(string(sel.Role)), sel.Label))
		}
		if sel.Field != "" {
			buf.WriteString(fmt.Sprintf(".\n\t\tWithFieldSelector(prometheus.KubernetesRole%s, %q)", g.capitalizeRole(string(sel.Role)), sel.Field))
		}
	}

	if sd.BearerTokenFile != "" {
		buf.WriteString(fmt.Sprintf(".\n\t\tWithBearerTokenFile(%q)", sd.BearerTokenFile))
	}

	return buf.String()
}

func (g *codeGenerator) capitalizeRole(role string) string {
	switch role {
	case "pod":
		return "Pod"
	case "node":
		return "Node"
	case "service":
		return "Service"
	case "endpoints":
		return "Endpoints"
	case "endpointslice":
		return "EndpointSlice"
	case "ingress":
		return "Ingress"
	default:
		return strings.Title(role)
	}
}

func (g *codeGenerator) writeRemoteWrite(buf *bytes.Buffer, rw *prometheus.RemoteWriteConfig, index int) error {
	varName := "RemoteWrite"
	if rw.Name != "" {
		varName = g.sanitizeVarName(rw.Name) + "RemoteWrite"
	} else if index > 0 {
		varName = fmt.Sprintf("RemoteWrite%d", index)
	}

	buf.WriteString(fmt.Sprintf("// %s configures remote write to %s.\n", varName, rw.URL))
	buf.WriteString(fmt.Sprintf("var %s = prometheus.NewRemoteWrite(%q)", varName, rw.URL))

	if rw.Name != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithName(%q)", rw.Name))
	}
	if rw.RemoteTimeout != 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithTimeout(%s)", g.formatDuration(rw.RemoteTimeout)))
	}
	if rw.BearerTokenFile != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithBearerTokenFile(%q)", rw.BearerTokenFile))
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *codeGenerator) writeRemoteRead(buf *bytes.Buffer, rr *prometheus.RemoteReadConfig, index int) error {
	varName := "RemoteRead"
	if rr.Name != "" {
		varName = g.sanitizeVarName(rr.Name) + "RemoteRead"
	} else if index > 0 {
		varName = fmt.Sprintf("RemoteRead%d", index)
	}

	buf.WriteString(fmt.Sprintf("// %s configures remote read from %s.\n", varName, rr.URL))
	buf.WriteString(fmt.Sprintf("var %s = prometheus.NewRemoteRead(%q)", varName, rr.URL))

	if rr.Name != "" {
		buf.WriteString(fmt.Sprintf(".\n\tWithName(%q)", rr.Name))
	}
	if rr.RemoteTimeout != 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithTimeout(%s)", g.formatDuration(rr.RemoteTimeout)))
	}
	if rr.ReadRecent {
		buf.WriteString(".\n\tWithReadRecent(true)")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *codeGenerator) writeMainConfig(buf *bytes.Buffer) error {
	buf.WriteString("// Config is the main Prometheus configuration.\n")
	buf.WriteString("var Config = &prometheus.PrometheusConfig{\n")

	if g.config.Global != nil {
		buf.WriteString("\tGlobal: GlobalConfig,\n")
	}

	if len(g.config.ScrapeConfigs) > 0 {
		buf.WriteString("\tScrapeConfigs: []*prometheus.ScrapeConfig{\n")
		for _, sc := range g.config.ScrapeConfigs {
			varName := g.sanitizeVarName(sc.JobName) + "Scrape"
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t},\n")
	}

	if len(g.config.RemoteWrite) > 0 {
		buf.WriteString("\tRemoteWrite: []*prometheus.RemoteWriteConfig{\n")
		for i, rw := range g.config.RemoteWrite {
			varName := "RemoteWrite"
			if rw.Name != "" {
				varName = g.sanitizeVarName(rw.Name) + "RemoteWrite"
			} else if i > 0 {
				varName = fmt.Sprintf("RemoteWrite%d", i)
			}
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t},\n")
	}

	if len(g.config.RemoteRead) > 0 {
		buf.WriteString("\tRemoteRead: []*prometheus.RemoteReadConfig{\n")
		for i, rr := range g.config.RemoteRead {
			varName := "RemoteRead"
			if rr.Name != "" {
				varName = g.sanitizeVarName(rr.Name) + "RemoteRead"
			} else if i > 0 {
				varName = fmt.Sprintf("RemoteRead%d", i)
			}
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t},\n")
	}

	if len(g.config.RuleFiles) > 0 {
		buf.WriteString("\tRuleFiles: []string{\n")
		for _, rf := range g.config.RuleFiles {
			buf.WriteString(fmt.Sprintf("\t\t%q,\n", rf))
		}
		buf.WriteString("\t},\n")
	}

	buf.WriteString("}\n")
	return nil
}

func (g *codeGenerator) formatDuration(d prometheus.Duration) string {
	dur := time.Duration(d)

	// Use convenience constants where possible
	if dur%time.Hour == 0 && dur >= time.Hour {
		hours := dur / time.Hour
		if hours == 1 {
			return "prometheus.Hour"
		}
		return fmt.Sprintf("%d * prometheus.Hour", hours)
	}
	if dur%time.Minute == 0 && dur >= time.Minute {
		mins := dur / time.Minute
		if mins == 1 {
			return "prometheus.Minute"
		}
		return fmt.Sprintf("%d * prometheus.Minute", mins)
	}
	if dur%time.Second == 0 {
		secs := dur / time.Second
		if secs == 1 {
			return "prometheus.Second"
		}
		return fmt.Sprintf("%d * prometheus.Second", secs)
	}

	// For milliseconds
	ms := dur / time.Millisecond
	return fmt.Sprintf("prometheus.Duration(%d * time.Millisecond)", ms)
}

func (g *codeGenerator) sanitizeVarName(name string) string {
	// Convert job name to valid Go identifier
	result := strings.Builder{}
	capitalize := true

	for _, c := range name {
		if c == '-' || c == '_' || c == '.' || c == '/' {
			capitalize = true
			continue
		}
		if capitalize {
			result.WriteRune(toUpper(c))
			capitalize = false
		} else {
			result.WriteRune(c)
		}
	}

	s := result.String()
	if s == "" {
		return "Config"
	}
	return s
}

func toUpper(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}

// templateFuncs provides helper functions for templates.
var templateFuncs = template.FuncMap{
	"quote": func(s string) string { return fmt.Sprintf("%q", s) },
}
