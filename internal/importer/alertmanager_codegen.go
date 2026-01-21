package importer

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"time"

	"github.com/lex00/wetwire-observability-go/alertmanager"
)

// GenerateAlertmanagerGoCode generates Go source code from an AlertmanagerConfig.
func GenerateAlertmanagerGoCode(config *alertmanager.AlertmanagerConfig, packageName string) ([]byte, error) {
	gen := &alertmanagerCodeGenerator{
		packageName: packageName,
		config:      config,
	}
	return gen.generate()
}

type alertmanagerCodeGenerator struct {
	packageName string
	config      *alertmanager.AlertmanagerConfig
}

func (g *alertmanagerCodeGenerator) generate() ([]byte, error) {
	var buf bytes.Buffer

	// Write package and imports
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"github.com/lex00/wetwire-observability-go/alertmanager\"\n")
	buf.WriteString(")\n\n")

	// Generate global config if present
	if g.config.Global != nil {
		if err := g.writeGlobalConfig(&buf); err != nil {
			return nil, err
		}
	}

	// Generate receivers
	for _, r := range g.config.Receivers {
		if err := g.writeReceiver(&buf, r); err != nil {
			return nil, err
		}
	}

	// Generate routes
	if g.config.Route != nil {
		if err := g.writeRoute(&buf, g.config.Route, "RootRoute"); err != nil {
			return nil, err
		}
	}

	// Generate inhibit rules
	for i, ir := range g.config.InhibitRules {
		if err := g.writeInhibitRule(&buf, ir, i); err != nil {
			return nil, err
		}
	}

	// Generate mute time intervals
	for _, mti := range g.config.MuteTimeIntervals {
		if err := g.writeMuteTimeInterval(&buf, mti); err != nil {
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

func (g *alertmanagerCodeGenerator) writeGlobalConfig(buf *bytes.Buffer) error {
	gc := g.config.Global

	buf.WriteString("// GlobalConfig defines global Alertmanager settings.\n")
	buf.WriteString("var GlobalConfig = &alertmanager.GlobalConfig{\n")

	if gc.SMTPSmarthost != "" {
		buf.WriteString(fmt.Sprintf("\tSMTPSmarthost: %q,\n", gc.SMTPSmarthost))
	}
	if gc.SMTPFrom != "" {
		buf.WriteString(fmt.Sprintf("\tSMTPFrom: %q,\n", gc.SMTPFrom))
	}
	if gc.SMTPAuthUsername != "" {
		buf.WriteString(fmt.Sprintf("\tSMTPAuthUsername: %q,\n", gc.SMTPAuthUsername))
	}
	if gc.SlackAPIURL != "" {
		buf.WriteString(fmt.Sprintf("\tSlackAPIURL: %q,\n", gc.SlackAPIURL))
	}
	if gc.PagerDutyURL != "" {
		buf.WriteString(fmt.Sprintf("\tPagerDutyURL: %q,\n", gc.PagerDutyURL))
	}
	if gc.OpsGenieAPIURL != "" {
		buf.WriteString(fmt.Sprintf("\tOpsGenieAPIURL: %q,\n", gc.OpsGenieAPIURL))
	}
	if gc.ResolveTimeout != 0 {
		buf.WriteString(fmt.Sprintf("\tResolveTimeout: %s,\n", g.formatDuration(gc.ResolveTimeout)))
	}

	buf.WriteString("}\n\n")
	return nil
}

func (g *alertmanagerCodeGenerator) writeReceiver(buf *bytes.Buffer, r *alertmanager.Receiver) error {
	varName := g.sanitizeVarName(r.Name) + "Receiver"

	buf.WriteString(fmt.Sprintf("// %s defines the %s receiver.\n", varName, r.Name))
	buf.WriteString(fmt.Sprintf("var %s = &alertmanager.Receiver{\n", varName))
	buf.WriteString(fmt.Sprintf("\tName: %q,\n", r.Name))

	// Email configs
	if len(r.EmailConfigs) > 0 {
		buf.WriteString("\tEmailConfigs: []*alertmanager.EmailConfig{\n")
		for _, ec := range r.EmailConfigs {
			buf.WriteString("\t\t{\n")
			if ec.To != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tTo: %q,\n", ec.To))
			}
			if ec.From != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tFrom: %q,\n", ec.From))
			}
			if ec.Smarthost != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tSmarthost: %q,\n", ec.Smarthost))
			}
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	// Slack configs
	if len(r.SlackConfigs) > 0 {
		buf.WriteString("\tSlackConfigs: []*alertmanager.SlackConfig{\n")
		for _, sc := range r.SlackConfigs {
			buf.WriteString("\t\t{\n")
			if sc.Channel != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tChannel: %q,\n", sc.Channel))
			}
			if sc.APIURL != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tAPIURL: %q,\n", sc.APIURL))
			}
			if sc.Username != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tUsername: %q,\n", sc.Username))
			}
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	// PagerDuty configs
	if len(r.PagerDutyConfigs) > 0 {
		buf.WriteString("\tPagerDutyConfigs: []*alertmanager.PagerDutyConfig{\n")
		for _, pd := range r.PagerDutyConfigs {
			buf.WriteString("\t\t{\n")
			if pd.RoutingKey != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tRoutingKey: %q,\n", pd.RoutingKey))
			}
			if pd.ServiceKey != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tServiceKey: %q,\n", pd.ServiceKey))
			}
			if pd.Severity != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tSeverity: %q,\n", pd.Severity))
			}
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	// Webhook configs
	if len(r.WebhookConfigs) > 0 {
		buf.WriteString("\tWebhookConfigs: []*alertmanager.WebhookConfig{\n")
		for _, wc := range r.WebhookConfigs {
			buf.WriteString("\t\t{\n")
			if wc.URL != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tURL: %q,\n", wc.URL))
			}
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	// OpsGenie configs
	if len(r.OpsGenieConfigs) > 0 {
		buf.WriteString("\tOpsGenieConfigs: []*alertmanager.OpsGenieConfig{\n")
		for _, og := range r.OpsGenieConfigs {
			buf.WriteString("\t\t{\n")
			if og.APIKey != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tAPIKey: %q,\n", og.APIKey))
			}
			if og.APIURL != "" {
				buf.WriteString(fmt.Sprintf("\t\t\tAPIURL: %q,\n", og.APIURL))
			}
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	buf.WriteString("}\n\n")
	return nil
}

func (g *alertmanagerCodeGenerator) writeRoute(buf *bytes.Buffer, r *alertmanager.Route, varName string) error {
	buf.WriteString(fmt.Sprintf("// %s defines the routing configuration.\n", varName))
	buf.WriteString(fmt.Sprintf("var %s = &alertmanager.Route{\n", varName))

	if r.Receiver != "" {
		buf.WriteString(fmt.Sprintf("\tReceiver: %q,\n", r.Receiver))
	}

	if len(r.GroupBy) > 0 {
		groupBy := make([]string, len(r.GroupBy))
		for i, g := range r.GroupBy {
			groupBy[i] = fmt.Sprintf("%q", g)
		}
		buf.WriteString(fmt.Sprintf("\tGroupBy: []string{%s},\n", strings.Join(groupBy, ", ")))
	}

	if r.GroupWait != 0 {
		buf.WriteString(fmt.Sprintf("\tGroupWait: %s,\n", g.formatDuration(r.GroupWait)))
	}
	if r.GroupInterval != 0 {
		buf.WriteString(fmt.Sprintf("\tGroupInterval: %s,\n", g.formatDuration(r.GroupInterval)))
	}
	if r.RepeatInterval != 0 {
		buf.WriteString(fmt.Sprintf("\tRepeatInterval: %s,\n", g.formatDuration(r.RepeatInterval)))
	}

	if r.Continue {
		buf.WriteString("\tContinue: true,\n")
	}

	// Matchers
	if len(r.Matchers) > 0 {
		buf.WriteString("\tMatchers: []*alertmanager.Matcher{\n")
		for _, m := range r.Matchers {
			buf.WriteString(fmt.Sprintf("\t\t{Label: %q, Op: alertmanager.MatchOp(%q), Value: %q},\n",
				m.Label, m.Op, m.Value))
		}
		buf.WriteString("\t},\n")
	}

	// Child routes (inline for simplicity)
	if len(r.Routes) > 0 {
		buf.WriteString("\tRoutes: []*alertmanager.Route{\n")
		for _, child := range r.Routes {
			g.writeInlineRoute(buf, child, "\t\t")
		}
		buf.WriteString("\t},\n")
	}

	if len(r.MuteTimeIntervals) > 0 {
		mti := make([]string, len(r.MuteTimeIntervals))
		for i, m := range r.MuteTimeIntervals {
			mti[i] = fmt.Sprintf("%q", m)
		}
		buf.WriteString(fmt.Sprintf("\tMuteTimeIntervals: []string{%s},\n", strings.Join(mti, ", ")))
	}

	buf.WriteString("}\n\n")
	return nil
}

func (g *alertmanagerCodeGenerator) writeInlineRoute(buf *bytes.Buffer, r *alertmanager.Route, indent string) {
	buf.WriteString(indent + "{\n")

	if r.Receiver != "" {
		buf.WriteString(fmt.Sprintf("%s\tReceiver: %q,\n", indent, r.Receiver))
	}

	if len(r.GroupBy) > 0 {
		groupBy := make([]string, len(r.GroupBy))
		for i, g := range r.GroupBy {
			groupBy[i] = fmt.Sprintf("%q", g)
		}
		buf.WriteString(fmt.Sprintf("%s\tGroupBy: []string{%s},\n", indent, strings.Join(groupBy, ", ")))
	}

	if r.Continue {
		buf.WriteString(indent + "\tContinue: true,\n")
	}

	if len(r.Matchers) > 0 {
		buf.WriteString(indent + "\tMatchers: []*alertmanager.Matcher{\n")
		for _, m := range r.Matchers {
			buf.WriteString(fmt.Sprintf("%s\t\t{Label: %q, Op: alertmanager.MatchOp(%q), Value: %q},\n",
				indent, m.Label, m.Op, m.Value))
		}
		buf.WriteString(indent + "\t},\n")
	}

	// Recursive child routes
	if len(r.Routes) > 0 {
		buf.WriteString(indent + "\tRoutes: []*alertmanager.Route{\n")
		for _, child := range r.Routes {
			g.writeInlineRoute(buf, child, indent+"\t\t")
		}
		buf.WriteString(indent + "\t},\n")
	}

	buf.WriteString(indent + "},\n")
}

func (g *alertmanagerCodeGenerator) writeInhibitRule(buf *bytes.Buffer, ir *alertmanager.InhibitRule, index int) error {
	varName := fmt.Sprintf("InhibitRule%d", index)

	buf.WriteString(fmt.Sprintf("// %s defines an inhibition rule.\n", varName))
	buf.WriteString(fmt.Sprintf("var %s = &alertmanager.InhibitRule{\n", varName))

	if len(ir.SourceMatchers) > 0 {
		buf.WriteString("\tSourceMatchers: []*alertmanager.Matcher{\n")
		for _, m := range ir.SourceMatchers {
			buf.WriteString(fmt.Sprintf("\t\t{Label: %q, Op: alertmanager.MatchOp(%q), Value: %q},\n",
				m.Label, m.Op, m.Value))
		}
		buf.WriteString("\t},\n")
	}

	if len(ir.TargetMatchers) > 0 {
		buf.WriteString("\tTargetMatchers: []*alertmanager.Matcher{\n")
		for _, m := range ir.TargetMatchers {
			buf.WriteString(fmt.Sprintf("\t\t{Label: %q, Op: alertmanager.MatchOp(%q), Value: %q},\n",
				m.Label, m.Op, m.Value))
		}
		buf.WriteString("\t},\n")
	}

	if len(ir.Equal) > 0 {
		equal := make([]string, len(ir.Equal))
		for i, e := range ir.Equal {
			equal[i] = fmt.Sprintf("%q", e)
		}
		buf.WriteString(fmt.Sprintf("\tEqual: []string{%s},\n", strings.Join(equal, ", ")))
	}

	buf.WriteString("}\n\n")
	return nil
}

func (g *alertmanagerCodeGenerator) writeMuteTimeInterval(buf *bytes.Buffer, mti *alertmanager.MuteTimeInterval) error {
	varName := g.sanitizeVarName(mti.Name) + "MuteTime"

	buf.WriteString(fmt.Sprintf("// %s defines a mute time interval.\n", varName))
	buf.WriteString(fmt.Sprintf("var %s = &alertmanager.MuteTimeInterval{\n", varName))
	buf.WriteString(fmt.Sprintf("\tName: %q,\n", mti.Name))

	if len(mti.TimeIntervals) > 0 {
		buf.WriteString("\tTimeIntervals: []alertmanager.TimeInterval{\n")
		for range mti.TimeIntervals {
			buf.WriteString("\t\t{\n")
			// TODO: Write TimeInterval fields
			buf.WriteString("\t\t},\n")
		}
		buf.WriteString("\t},\n")
	}

	buf.WriteString("}\n\n")
	return nil
}

func (g *alertmanagerCodeGenerator) writeMainConfig(buf *bytes.Buffer) error {
	buf.WriteString("// Config is the main Alertmanager configuration.\n")
	buf.WriteString("var Config = &alertmanager.AlertmanagerConfig{\n")

	if g.config.Global != nil {
		buf.WriteString("\tGlobal: GlobalConfig,\n")
	}

	if g.config.Route != nil {
		buf.WriteString("\tRoute: RootRoute,\n")
	}

	if len(g.config.Receivers) > 0 {
		buf.WriteString("\tReceivers: []*alertmanager.Receiver{\n")
		for _, r := range g.config.Receivers {
			varName := g.sanitizeVarName(r.Name) + "Receiver"
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t},\n")
	}

	if len(g.config.InhibitRules) > 0 {
		buf.WriteString("\tInhibitRules: []*alertmanager.InhibitRule{\n")
		for i := range g.config.InhibitRules {
			varName := fmt.Sprintf("InhibitRule%d", i)
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t},\n")
	}

	if len(g.config.MuteTimeIntervals) > 0 {
		buf.WriteString("\tMuteTimeIntervals: []*alertmanager.MuteTimeInterval{\n")
		for _, mti := range g.config.MuteTimeIntervals {
			varName := g.sanitizeVarName(mti.Name) + "MuteTime"
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t},\n")
	}

	if len(g.config.Templates) > 0 {
		buf.WriteString("\tTemplates: []string{\n")
		for _, t := range g.config.Templates {
			buf.WriteString(fmt.Sprintf("\t\t%q,\n", t))
		}
		buf.WriteString("\t},\n")
	}

	buf.WriteString("}\n")
	return nil
}

func (g *alertmanagerCodeGenerator) formatDuration(d alertmanager.Duration) string {
	dur := time.Duration(d)

	if dur%time.Hour == 0 && dur >= time.Hour {
		hours := dur / time.Hour
		if hours == 1 {
			return "alertmanager.Hour"
		}
		return fmt.Sprintf("%d * alertmanager.Hour", hours)
	}
	if dur%time.Minute == 0 && dur >= time.Minute {
		mins := dur / time.Minute
		if mins == 1 {
			return "alertmanager.Minute"
		}
		return fmt.Sprintf("%d * alertmanager.Minute", mins)
	}
	if dur%time.Second == 0 {
		secs := dur / time.Second
		if secs == 1 {
			return "alertmanager.Second"
		}
		return fmt.Sprintf("%d * alertmanager.Second", secs)
	}

	ms := dur / time.Millisecond
	return fmt.Sprintf("alertmanager.Duration(%d * time.Millisecond)", ms)
}

func (g *alertmanagerCodeGenerator) sanitizeVarName(name string) string {
	result := strings.Builder{}
	capitalize := true

	for _, c := range name {
		if c == '-' || c == '_' || c == '.' || c == '/' || c == ' ' {
			capitalize = true
			continue
		}
		if capitalize {
			result.WriteRune(toUpperRune(c))
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

func toUpperRune(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}
