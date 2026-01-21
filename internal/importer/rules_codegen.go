package importer

import (
	"bytes"
	"fmt"
	"go/format"
	"strings"
	"time"
)

// GenerateRulesGoCode generates Go source code from a RulesFile.
func GenerateRulesGoCode(rf *RulesFile, packageName string) ([]byte, error) {
	gen := &rulesCodeGenerator{
		packageName: packageName,
		rulesFile:   rf,
	}
	return gen.generate()
}

type rulesCodeGenerator struct {
	packageName string
	rulesFile   *RulesFile
}

func (g *rulesCodeGenerator) generate() ([]byte, error) {
	var buf bytes.Buffer

	// Write package and imports
	buf.WriteString(fmt.Sprintf("package %s\n\n", g.packageName))
	buf.WriteString("import (\n")
	buf.WriteString("\t\"github.com/lex00/wetwire-observability-go/rules\"\n")
	buf.WriteString(")\n\n")

	// Track rule variable names per group
	groupRuleVars := make(map[string][]string)

	// Generate individual rules first
	for _, group := range g.rulesFile.Groups {
		var ruleVars []string

		for i, rule := range group.Rules {
			varName := g.getRuleVarName(group.Name, rule, i)
			if err := g.writeRule(&buf, rule, varName); err != nil {
				return nil, err
			}
			ruleVars = append(ruleVars, varName)
		}

		groupRuleVars[group.Name] = ruleVars
	}

	// Generate rule groups
	for _, group := range g.rulesFile.Groups {
		if err := g.writeRuleGroup(&buf, group, groupRuleVars[group.Name]); err != nil {
			return nil, err
		}
	}

	// Generate the main rules file
	if err := g.writeRulesFile(&buf); err != nil {
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

func (g *rulesCodeGenerator) getRuleVarName(groupName string, rule RuleInput, index int) string {
	var name string
	if rule.IsAlertingRule() {
		name = rule.Alert
	} else {
		name = rule.Record
	}

	if name != "" {
		return g.sanitizeVarName(name)
	}

	return fmt.Sprintf("%sRule%d", g.sanitizeVarName(groupName), index)
}

func (g *rulesCodeGenerator) writeRule(buf *bytes.Buffer, rule RuleInput, varName string) error {
	if rule.IsAlertingRule() {
		return g.writeAlertingRule(buf, rule, varName)
	}
	return g.writeRecordingRule(buf, rule, varName)
}

func (g *rulesCodeGenerator) writeAlertingRule(buf *bytes.Buffer, rule RuleInput, varName string) error {
	buf.WriteString(fmt.Sprintf("// %s is an alerting rule.\n", varName))
	buf.WriteString(fmt.Sprintf("var %s = rules.NewAlertingRule(%q)", varName, rule.Alert))

	// Expression
	buf.WriteString(fmt.Sprintf(".\n\tWithExpr(%q)", rule.Expr))

	// For duration
	if rule.For != "" {
		dur := g.formatDuration(rule.For)
		buf.WriteString(fmt.Sprintf(".\n\tWithFor(%s)", dur))
	}

	// KeepFiringFor duration
	if rule.KeepFiringFor != "" {
		dur := g.formatDuration(rule.KeepFiringFor)
		buf.WriteString(fmt.Sprintf(".\n\tWithKeepFiringFor(%s)", dur))
	}

	// Labels - use convenience methods for severity
	if sev, ok := rule.Labels["severity"]; ok {
		switch sev {
		case "critical":
			buf.WriteString(".\n\tCritical()")
		case "warning":
			buf.WriteString(".\n\tWarning()")
		case "info":
			buf.WriteString(".\n\tInfo()")
		default:
			// Add as custom label below
		}
	}

	// Other labels (excluding severity if already handled)
	otherLabels := make(map[string]string)
	for k, v := range rule.Labels {
		if k != "severity" || (k == "severity" && v != "critical" && v != "warning" && v != "info") {
			otherLabels[k] = v
		}
	}
	if len(otherLabels) > 0 {
		buf.WriteString(".\n\tWithLabels(map[string]string{")
		for k, v := range otherLabels {
			buf.WriteString(fmt.Sprintf("\n\t\t%q: %q,", k, v))
		}
		buf.WriteString("\n\t})")
	}

	// Annotations - use convenience methods for common ones
	hasAnnotations := false
	if summary, ok := rule.Annotations["summary"]; ok {
		buf.WriteString(fmt.Sprintf(".\n\tWithSummary(%q)", summary))
		hasAnnotations = true
	}
	if desc, ok := rule.Annotations["description"]; ok {
		buf.WriteString(fmt.Sprintf(".\n\tWithDescription(%q)", desc))
		hasAnnotations = true
	}
	if runbook, ok := rule.Annotations["runbook_url"]; ok {
		buf.WriteString(fmt.Sprintf(".\n\tWithRunbook(%q)", runbook))
		hasAnnotations = true
	}

	// Other annotations
	otherAnnotations := make(map[string]string)
	for k, v := range rule.Annotations {
		if k != "summary" && k != "description" && k != "runbook_url" {
			otherAnnotations[k] = v
		}
	}
	if len(otherAnnotations) > 0 {
		if hasAnnotations {
			buf.WriteString(".\n\tWithAnnotations(map[string]string{")
		} else {
			buf.WriteString(".\n\tWithAnnotations(map[string]string{")
		}
		for k, v := range otherAnnotations {
			buf.WriteString(fmt.Sprintf("\n\t\t%q: %q,", k, v))
		}
		buf.WriteString("\n\t})")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *rulesCodeGenerator) writeRecordingRule(buf *bytes.Buffer, rule RuleInput, varName string) error {
	buf.WriteString(fmt.Sprintf("// %s is a recording rule.\n", varName))
	buf.WriteString(fmt.Sprintf("var %s = rules.NewRecordingRule(%q)", varName, rule.Record))

	// Expression
	buf.WriteString(fmt.Sprintf(".\n\tWithExpr(%q)", rule.Expr))

	// Labels
	if len(rule.Labels) > 0 {
		buf.WriteString(".\n\tWithLabels(map[string]string{")
		for k, v := range rule.Labels {
			buf.WriteString(fmt.Sprintf("\n\t\t%q: %q,", k, v))
		}
		buf.WriteString("\n\t})")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *rulesCodeGenerator) writeRuleGroup(buf *bytes.Buffer, group RuleGroupInput, ruleVars []string) error {
	varName := g.sanitizeVarName(group.Name) + "Group"

	buf.WriteString(fmt.Sprintf("// %s contains rules for %s.\n", varName, group.Name))
	buf.WriteString(fmt.Sprintf("var %s = rules.NewRuleGroup(%q)", varName, group.Name))

	// Interval
	if group.Interval != "" {
		dur := g.formatDuration(group.Interval)
		buf.WriteString(fmt.Sprintf(".\n\tWithInterval(%s)", dur))
	}

	// Limit
	if group.Limit > 0 {
		buf.WriteString(fmt.Sprintf(".\n\tWithLimit(%d)", group.Limit))
	}

	// Rules
	if len(ruleVars) > 0 {
		buf.WriteString(".\n\tWithRules(\n")
		for _, rv := range ruleVars {
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", rv))
		}
		buf.WriteString("\t)")
	}

	buf.WriteString("\n\n")
	return nil
}

func (g *rulesCodeGenerator) writeRulesFile(buf *bytes.Buffer) error {
	buf.WriteString("// RulesConfig is the main rules configuration.\n")
	buf.WriteString("var RulesConfig = rules.NewRulesFile()")

	if len(g.rulesFile.Groups) > 0 {
		buf.WriteString(".\n\tWithGroups(\n")
		for _, group := range g.rulesFile.Groups {
			varName := g.sanitizeVarName(group.Name) + "Group"
			buf.WriteString(fmt.Sprintf("\t\t%s,\n", varName))
		}
		buf.WriteString("\t)")
	}

	buf.WriteString("\n")
	return nil
}

func (g *rulesCodeGenerator) formatDuration(s string) string {
	if s == "" {
		return "0"
	}

	// Parse the duration
	d, err := time.ParseDuration(s)
	if err != nil {
		// Return raw string as fallback
		return fmt.Sprintf("rules.Duration(%q)", s)
	}

	// Format nicely
	if d%time.Hour == 0 && d >= time.Hour {
		hours := d / time.Hour
		if hours == 1 {
			return "rules.Hour"
		}
		return fmt.Sprintf("%d * rules.Hour", hours)
	}
	if d%time.Minute == 0 && d >= time.Minute {
		mins := d / time.Minute
		if mins == 1 {
			return "rules.Minute"
		}
		return fmt.Sprintf("%d * rules.Minute", mins)
	}
	if d%time.Second == 0 {
		secs := d / time.Second
		if secs == 1 {
			return "rules.Second"
		}
		return fmt.Sprintf("%d * rules.Second", secs)
	}

	ms := d / time.Millisecond
	return fmt.Sprintf("rules.Duration(%d * time.Millisecond)", ms)
}

func (g *rulesCodeGenerator) sanitizeVarName(name string) string {
	result := strings.Builder{}
	capitalize := true

	for _, c := range name {
		if c == '-' || c == '_' || c == '.' || c == '/' || c == ' ' || c == ':' {
			capitalize = true
			continue
		}
		if capitalize {
			result.WriteRune(rulesUpperRune(c))
			capitalize = false
		} else {
			result.WriteRune(c)
		}
	}

	s := result.String()
	if s == "" {
		return "Rules"
	}
	return s
}

func rulesUpperRune(r rune) rune {
	if r >= 'a' && r <= 'z' {
		return r - 'a' + 'A'
	}
	return r
}
