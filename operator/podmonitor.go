package operator

import "gopkg.in/yaml.v3"

// PodMonitor represents a Prometheus Operator PodMonitor CRD.
type PodMonitor struct {
	APIVersion string         `yaml:"apiVersion"`
	Kind       string         `yaml:"kind"`
	Metadata   ObjectMeta     `yaml:"metadata"`
	Spec       PodMonitorSpec `yaml:"spec"`

	// Convenience fields (not serialized)
	Name      string            `yaml:"-"`
	Namespace string            `yaml:"-"`
	Labels    map[string]string `yaml:"-"`
}

// PodMonitorSpec contains the PodMonitor specification.
type PodMonitorSpec struct {
	// JobLabel is the label to use for the job name.
	JobLabel string `yaml:"jobLabel,omitempty"`

	// PodTargetLabels transfers labels from the pod to the scraped metrics.
	PodTargetLabels []string `yaml:"podTargetLabels,omitempty"`

	// PodMetricsEndpoints contains endpoint configurations.
	PodMetricsEndpoints []*PodMetricsEndpoint `yaml:"podMetricsEndpoints,omitempty"`

	// Selector selects pods to monitor.
	Selector LabelSelector `yaml:"selector"`

	// NamespaceSelector selects namespaces to discover pods in.
	NamespaceSelector NamespaceSelector `yaml:"namespaceSelector,omitempty"`

	// SampleLimit is the maximum number of samples per scrape.
	SampleLimit uint64 `yaml:"sampleLimit,omitempty"`

	// TargetLimit is the maximum number of targets.
	TargetLimit uint64 `yaml:"targetLimit,omitempty"`
}

// PodMetricsEndpoint defines an endpoint for pod metrics.
type PodMetricsEndpoint struct {
	// Port is the port name or number to scrape.
	Port string `yaml:"port,omitempty"`

	// TargetPort is the target port number when not using named ports.
	TargetPort int `yaml:"targetPort,omitempty"`

	// Path is the HTTP path to scrape (default: /metrics).
	Path string `yaml:"path,omitempty"`

	// Scheme is the URL scheme (http or https).
	Scheme string `yaml:"scheme,omitempty"`

	// Interval is the scrape interval.
	Interval string `yaml:"interval,omitempty"`

	// ScrapeTimeout is the timeout for the scrape.
	ScrapeTimeout string `yaml:"scrapeTimeout,omitempty"`

	// TLSConfig contains TLS configuration.
	TLSConfig *TLSConfig `yaml:"tlsConfig,omitempty"`

	// BearerTokenSecret references a secret containing the bearer token.
	BearerTokenSecret *SecretKeySelector `yaml:"bearerTokenSecret,omitempty"`

	// BasicAuth contains basic auth credentials.
	BasicAuth *BasicAuth `yaml:"basicAuth,omitempty"`

	// HonorLabels controls whether to honor labels from the scraped target.
	HonorLabels bool `yaml:"honorLabels,omitempty"`

	// HonorTimestamps controls whether to honor timestamps from the target.
	HonorTimestamps *bool `yaml:"honorTimestamps,omitempty"`

	// RelabelConfigs contains relabeling configuration.
	RelabelConfigs []*RelabelConfig `yaml:"relabelings,omitempty"`

	// MetricRelabelConfigs contains metric relabeling configuration.
	MetricRelabelConfigs []*RelabelConfig `yaml:"metricRelabelings,omitempty"`

	// Params are URL parameters for the scrape.
	Params map[string][]string `yaml:"params,omitempty"`
}

// NewPodMetricsEndpoint creates a new pod metrics endpoint.
func NewPodMetricsEndpoint(port string) *PodMetricsEndpoint {
	return &PodMetricsEndpoint{
		Port: port,
	}
}

// WithPath sets the scrape path.
func (e *PodMetricsEndpoint) WithPath(path string) *PodMetricsEndpoint {
	e.Path = path
	return e
}

// WithInterval sets the scrape interval.
func (e *PodMetricsEndpoint) WithInterval(interval string) *PodMetricsEndpoint {
	e.Interval = interval
	return e
}

// WithScrapeTimeout sets the scrape timeout.
func (e *PodMetricsEndpoint) WithScrapeTimeout(timeout string) *PodMetricsEndpoint {
	e.ScrapeTimeout = timeout
	return e
}

// WithScheme sets the URL scheme.
func (e *PodMetricsEndpoint) WithScheme(scheme string) *PodMetricsEndpoint {
	e.Scheme = scheme
	return e
}

// WithBearerTokenSecret sets the bearer token secret reference.
func (e *PodMetricsEndpoint) WithBearerTokenSecret(name, key string) *PodMetricsEndpoint {
	e.BearerTokenSecret = &SecretKeySelector{Name: name, Key: key}
	return e
}

// WithTLSConfig sets the TLS configuration.
func (e *PodMetricsEndpoint) WithTLSConfig(insecureSkipVerify bool, caFile, certFile, keyFile string) *PodMetricsEndpoint {
	e.TLSConfig = &TLSConfig{
		InsecureSkipVerify: insecureSkipVerify,
		CAFile:             caFile,
		CertFile:           certFile,
		KeyFile:            keyFile,
	}
	return e
}

// AddRelabeling adds a relabeling configuration.
func (e *PodMetricsEndpoint) AddRelabeling(r *RelabelConfig) *PodMetricsEndpoint {
	e.RelabelConfigs = append(e.RelabelConfigs, r)
	return e
}

// AddMetricRelabeling adds a metric relabeling configuration.
func (e *PodMetricsEndpoint) AddMetricRelabeling(r *RelabelConfig) *PodMetricsEndpoint {
	e.MetricRelabelConfigs = append(e.MetricRelabelConfigs, r)
	return e
}

// PodMon creates a new PodMonitor.
func PodMon(name, namespace string) *PodMonitor {
	return &PodMonitor{
		APIVersion: "monitoring.coreos.com/v1",
		Kind:       "PodMonitor",
		Metadata: ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Name:      name,
		Namespace: namespace,
	}
}

// WithLabels sets the PodMonitor labels.
func (pm *PodMonitor) WithLabels(labels map[string]string) *PodMonitor {
	pm.Labels = labels
	pm.Metadata.Labels = labels
	return pm
}

// AddLabel adds a label to the PodMonitor.
func (pm *PodMonitor) AddLabel(key, value string) *PodMonitor {
	if pm.Labels == nil {
		pm.Labels = make(map[string]string)
	}
	pm.Labels[key] = value
	pm.Metadata.Labels = pm.Labels
	return pm
}

// WithAnnotations sets the PodMonitor annotations.
func (pm *PodMonitor) WithAnnotations(annotations map[string]string) *PodMonitor {
	pm.Metadata.Annotations = annotations
	return pm
}

// SelectPods sets the pod selector with a single label.
func (pm *PodMonitor) SelectPods(labelKey, labelValue string) *PodMonitor {
	if pm.Spec.Selector.MatchLabels == nil {
		pm.Spec.Selector.MatchLabels = make(map[string]string)
	}
	pm.Spec.Selector.MatchLabels[labelKey] = labelValue
	return pm
}

// SelectPodsByLabels sets the pod selector with multiple labels.
func (pm *PodMonitor) SelectPodsByLabels(labels map[string]string) *PodMonitor {
	pm.Spec.Selector.MatchLabels = labels
	return pm
}

// SelectPodsByExpression adds a label selector expression.
func (pm *PodMonitor) SelectPodsByExpression(key, operator string, values ...string) *PodMonitor {
	pm.Spec.Selector.MatchExpressions = append(pm.Spec.Selector.MatchExpressions, LabelSelectorReq{
		Key:      key,
		Operator: operator,
		Values:   values,
	})
	return pm
}

// InNamespace sets the namespace selector to a single namespace.
func (pm *PodMonitor) InNamespace(namespace string) *PodMonitor {
	pm.Spec.NamespaceSelector = NamespaceSelector{
		MatchNames: []string{namespace},
	}
	return pm
}

// InNamespaces sets the namespace selector to multiple namespaces.
func (pm *PodMonitor) InNamespaces(namespaces ...string) *PodMonitor {
	pm.Spec.NamespaceSelector = NamespaceSelector{
		MatchNames: namespaces,
	}
	return pm
}

// InAllNamespaces selects pods in all namespaces.
func (pm *PodMonitor) InAllNamespaces() *PodMonitor {
	pm.Spec.NamespaceSelector = NamespaceSelector{
		Any: true,
	}
	return pm
}

// WithPodMetricsEndpoint sets a single endpoint.
func (pm *PodMonitor) WithPodMetricsEndpoint(endpoint *PodMetricsEndpoint) *PodMonitor {
	pm.Spec.PodMetricsEndpoints = []*PodMetricsEndpoint{endpoint}
	return pm
}

// AddPodMetricsEndpoint adds an endpoint.
func (pm *PodMonitor) AddPodMetricsEndpoint(endpoint *PodMetricsEndpoint) *PodMonitor {
	pm.Spec.PodMetricsEndpoints = append(pm.Spec.PodMetricsEndpoints, endpoint)
	return pm
}

// WithJobLabel sets the job label.
func (pm *PodMonitor) WithJobLabel(label string) *PodMonitor {
	pm.Spec.JobLabel = label
	return pm
}

// WithPodTargetLabels sets the labels to transfer from pods.
func (pm *PodMonitor) WithPodTargetLabels(labels ...string) *PodMonitor {
	pm.Spec.PodTargetLabels = labels
	return pm
}

// WithSampleLimit sets the sample limit per scrape.
func (pm *PodMonitor) WithSampleLimit(limit uint64) *PodMonitor {
	pm.Spec.SampleLimit = limit
	return pm
}

// WithTargetLimit sets the target limit.
func (pm *PodMonitor) WithTargetLimit(limit uint64) *PodMonitor {
	pm.Spec.TargetLimit = limit
	return pm
}

// Serialize converts the PodMonitor to YAML bytes.
func (pm *PodMonitor) Serialize() ([]byte, error) {
	return yaml.Marshal(pm)
}

// MustSerialize converts the PodMonitor to YAML bytes, panicking on error.
func (pm *PodMonitor) MustSerialize() []byte {
	data, err := pm.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}
