package operator

import (
	"gopkg.in/yaml.v3"

	"github.com/lex00/wetwire-observability-go/rules"
)

// PrometheusRule represents a Prometheus Operator PrometheusRule CRD.
type PrometheusRule struct {
	APIVersion string             `yaml:"apiVersion"`
	Kind       string             `yaml:"kind"`
	Metadata   ObjectMeta         `yaml:"metadata"`
	Spec       PrometheusRuleSpec `yaml:"spec"`

	// Convenience fields (not serialized)
	Name      string            `yaml:"-"`
	Namespace string            `yaml:"-"`
	Labels    map[string]string `yaml:"-"`
}

// PrometheusRuleSpec contains the PrometheusRule specification.
type PrometheusRuleSpec struct {
	Groups []*rules.RuleGroup `yaml:"groups"`
}

// PromRule creates a new PrometheusRule.
func PromRule(name, namespace string) *PrometheusRule {
	return &PrometheusRule{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "PrometheusRule",
		Metadata: ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Name:      name,
		Namespace: namespace,
	}
}

// PromRuleFromRulesFile creates a PrometheusRule from an existing RulesFile.
func PromRuleFromRulesFile(name, namespace string, rf *rules.RulesFile) *PrometheusRule {
	pr := PromRule(name, namespace)
	pr.Spec.Groups = rf.Groups
	return pr
}

// WithLabels sets the PrometheusRule labels.
func (pr *PrometheusRule) WithLabels(labels map[string]string) *PrometheusRule {
	pr.Labels = labels
	pr.Metadata.Labels = labels
	return pr
}

// AddLabel adds a label to the PrometheusRule.
func (pr *PrometheusRule) AddLabel(key, value string) *PrometheusRule {
	if pr.Labels == nil {
		pr.Labels = make(map[string]string)
	}
	pr.Labels[key] = value
	pr.Metadata.Labels = pr.Labels
	return pr
}

// WithAnnotations sets the PrometheusRule annotations.
func (pr *PrometheusRule) WithAnnotations(annotations map[string]string) *PrometheusRule {
	pr.Metadata.Annotations = annotations
	return pr
}

// AddAnnotation adds an annotation to the PrometheusRule.
func (pr *PrometheusRule) AddAnnotation(key, value string) *PrometheusRule {
	if pr.Metadata.Annotations == nil {
		pr.Metadata.Annotations = make(map[string]string)
	}
	pr.Metadata.Annotations[key] = value
	return pr
}

// WithRuleGroups sets the rule groups.
func (pr *PrometheusRule) WithRuleGroups(groups ...*rules.RuleGroup) *PrometheusRule {
	pr.Spec.Groups = groups
	return pr
}

// AddRuleGroup adds a rule group.
func (pr *PrometheusRule) AddRuleGroup(group *rules.RuleGroup) *PrometheusRule {
	pr.Spec.Groups = append(pr.Spec.Groups, group)
	return pr
}

// Serialize converts the PrometheusRule to YAML bytes.
func (pr *PrometheusRule) Serialize() ([]byte, error) {
	return yaml.Marshal(pr)
}

// MustSerialize converts the PrometheusRule to YAML bytes, panicking on error.
func (pr *PrometheusRule) MustSerialize() []byte {
	data, err := pr.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}
