package operator

import "gopkg.in/yaml.v3"

// ServiceMonitor represents a Prometheus Operator ServiceMonitor CRD.
type ServiceMonitor struct {
	APIVersion string            `yaml:"apiVersion"`
	Kind       string            `yaml:"kind"`
	Metadata   ObjectMeta        `yaml:"metadata"`
	Spec       ServiceMonitorSpec `yaml:"spec"`

	// Convenience fields (not serialized)
	Name      string            `yaml:"-"`
	Namespace string            `yaml:"-"`
	Labels    map[string]string `yaml:"-"`
}

// ObjectMeta contains Kubernetes object metadata.
type ObjectMeta struct {
	Name        string            `yaml:"name"`
	Namespace   string            `yaml:"namespace,omitempty"`
	Labels      map[string]string `yaml:"labels,omitempty"`
	Annotations map[string]string `yaml:"annotations,omitempty"`
}

// ServiceMonitorSpec contains the ServiceMonitor specification.
type ServiceMonitorSpec struct {
	// JobLabel is the label to use for the job name.
	JobLabel string `yaml:"jobLabel,omitempty"`

	// TargetLabels transfers labels from the service to the scraped metrics.
	TargetLabels []string `yaml:"targetLabels,omitempty"`

	// PodTargetLabels transfers labels from the pod to the scraped metrics.
	PodTargetLabels []string `yaml:"podTargetLabels,omitempty"`

	// Endpoints contains endpoint configurations.
	Endpoints []*Endpoint `yaml:"endpoints,omitempty"`

	// Selector selects services to monitor.
	Selector LabelSelector `yaml:"selector"`

	// NamespaceSelector selects namespaces to discover services in.
	NamespaceSelector NamespaceSelector `yaml:"namespaceSelector,omitempty"`

	// SampleLimit is the maximum number of samples per scrape.
	SampleLimit uint64 `yaml:"sampleLimit,omitempty"`

	// TargetLimit is the maximum number of targets.
	TargetLimit uint64 `yaml:"targetLimit,omitempty"`
}

// LabelSelector selects Kubernetes objects by labels.
type LabelSelector struct {
	MatchLabels      map[string]string    `yaml:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorReq `yaml:"matchExpressions,omitempty"`
}

// LabelSelectorReq is a label selector requirement.
type LabelSelectorReq struct {
	Key      string   `yaml:"key"`
	Operator string   `yaml:"operator"`
	Values   []string `yaml:"values,omitempty"`
}

// NamespaceSelector selects namespaces.
type NamespaceSelector struct {
	Any        bool     `yaml:"any,omitempty"`
	MatchNames []string `yaml:"matchNames,omitempty"`
}

// ServiceMon creates a new ServiceMonitor.
func ServiceMon(name, namespace string) *ServiceMonitor {
	return &ServiceMonitor{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "ServiceMonitor",
		Metadata: ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Name:      name,
		Namespace: namespace,
	}
}

// WithLabels sets the ServiceMonitor labels.
func (sm *ServiceMonitor) WithLabels(labels map[string]string) *ServiceMonitor {
	sm.Labels = labels
	sm.Metadata.Labels = labels
	return sm
}

// AddLabel adds a label to the ServiceMonitor.
func (sm *ServiceMonitor) AddLabel(key, value string) *ServiceMonitor {
	if sm.Labels == nil {
		sm.Labels = make(map[string]string)
	}
	sm.Labels[key] = value
	sm.Metadata.Labels = sm.Labels
	return sm
}

// WithAnnotations sets the ServiceMonitor annotations.
func (sm *ServiceMonitor) WithAnnotations(annotations map[string]string) *ServiceMonitor {
	sm.Metadata.Annotations = annotations
	return sm
}

// SelectService sets the service selector with a single label.
func (sm *ServiceMonitor) SelectService(labelKey, labelValue string) *ServiceMonitor {
	if sm.Spec.Selector.MatchLabels == nil {
		sm.Spec.Selector.MatchLabels = make(map[string]string)
	}
	sm.Spec.Selector.MatchLabels[labelKey] = labelValue
	return sm
}

// SelectServiceByLabels sets the service selector with multiple labels.
func (sm *ServiceMonitor) SelectServiceByLabels(labels map[string]string) *ServiceMonitor {
	sm.Spec.Selector.MatchLabels = labels
	return sm
}

// SelectServiceByExpression adds a label selector expression.
func (sm *ServiceMonitor) SelectServiceByExpression(key, operator string, values ...string) *ServiceMonitor {
	sm.Spec.Selector.MatchExpressions = append(sm.Spec.Selector.MatchExpressions, LabelSelectorReq{
		Key:      key,
		Operator: operator,
		Values:   values,
	})
	return sm
}

// InNamespace sets the namespace selector to a single namespace.
func (sm *ServiceMonitor) InNamespace(namespace string) *ServiceMonitor {
	sm.Spec.NamespaceSelector = NamespaceSelector{
		MatchNames: []string{namespace},
	}
	return sm
}

// InNamespaces sets the namespace selector to multiple namespaces.
func (sm *ServiceMonitor) InNamespaces(namespaces ...string) *ServiceMonitor {
	sm.Spec.NamespaceSelector = NamespaceSelector{
		MatchNames: namespaces,
	}
	return sm
}

// InAllNamespaces selects services in all namespaces.
func (sm *ServiceMonitor) InAllNamespaces() *ServiceMonitor {
	sm.Spec.NamespaceSelector = NamespaceSelector{
		Any: true,
	}
	return sm
}

// WithEndpoint sets a single endpoint.
func (sm *ServiceMonitor) WithEndpoint(endpoint *Endpoint) *ServiceMonitor {
	sm.Spec.Endpoints = []*Endpoint{endpoint}
	return sm
}

// AddEndpoint adds an endpoint.
func (sm *ServiceMonitor) AddEndpoint(endpoint *Endpoint) *ServiceMonitor {
	sm.Spec.Endpoints = append(sm.Spec.Endpoints, endpoint)
	return sm
}

// WithJobLabel sets the job label.
func (sm *ServiceMonitor) WithJobLabel(label string) *ServiceMonitor {
	sm.Spec.JobLabel = label
	return sm
}

// WithTargetLabels sets the labels to transfer from services.
func (sm *ServiceMonitor) WithTargetLabels(labels ...string) *ServiceMonitor {
	sm.Spec.TargetLabels = labels
	return sm
}

// WithPodTargetLabels sets the labels to transfer from pods.
func (sm *ServiceMonitor) WithPodTargetLabels(labels ...string) *ServiceMonitor {
	sm.Spec.PodTargetLabels = labels
	return sm
}

// WithSampleLimit sets the sample limit per scrape.
func (sm *ServiceMonitor) WithSampleLimit(limit uint64) *ServiceMonitor {
	sm.Spec.SampleLimit = limit
	return sm
}

// WithTargetLimit sets the target limit.
func (sm *ServiceMonitor) WithTargetLimit(limit uint64) *ServiceMonitor {
	sm.Spec.TargetLimit = limit
	return sm
}

// Serialize converts the ServiceMonitor to YAML bytes.
func (sm *ServiceMonitor) Serialize() ([]byte, error) {
	return yaml.Marshal(sm)
}

// MustSerialize converts the ServiceMonitor to YAML bytes, panicking on error.
func (sm *ServiceMonitor) MustSerialize() []byte {
	data, err := sm.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}
