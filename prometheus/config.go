package prometheus

// PrometheusConfig represents the top-level Prometheus configuration.
// This is the root structure that serializes to prometheus.yml.
type PrometheusConfig struct {
	// Global configures values that are shared across all configuration contexts.
	Global *GlobalConfig `yaml:"global,omitempty"`

	// ScrapeConfigs defines the scrape configurations for targets.
	ScrapeConfigs []*ScrapeConfig `yaml:"scrape_configs,omitempty"`

	// RuleFiles specifies a list of globs for rule file paths.
	RuleFiles []string `yaml:"rule_files,omitempty"`

	// AlertingConfig configures alerting and Alertmanager communication.
	Alerting *AlertingConfig `yaml:"alerting,omitempty"`

	// RemoteWrite specifies settings for remote write.
	RemoteWrite []*RemoteWriteConfig `yaml:"remote_write,omitempty"`

	// RemoteRead specifies settings for remote read.
	RemoteRead []*RemoteReadConfig `yaml:"remote_read,omitempty"`
}

// GlobalConfig configures values that apply to all other configuration contexts.
type GlobalConfig struct {
	// ScrapeInterval is the default scrape interval for targets.
	// Defaults to 1m if not specified.
	ScrapeInterval Duration `yaml:"scrape_interval,omitempty"`

	// ScrapeTimeout is how long until a scrape request times out.
	// Defaults to 10s if not specified.
	ScrapeTimeout Duration `yaml:"scrape_timeout,omitempty"`

	// EvaluationInterval is how frequently to evaluate rules.
	// Defaults to 1m if not specified.
	EvaluationInterval Duration `yaml:"evaluation_interval,omitempty"`

	// ExternalLabels are labels to add to any time series or alerts
	// when communicating with external systems.
	ExternalLabels map[string]string `yaml:"external_labels,omitempty"`
}

// AlertingConfig configures alerting and Alertmanager communication.
type AlertingConfig struct {
	// Alertmanagers defines Alertmanager instances.
	Alertmanagers []*AlertmanagerConfig `yaml:"alertmanagers,omitempty"`
}

// AlertmanagerConfig configures an Alertmanager instance for alert delivery.
type AlertmanagerConfig struct {
	// StaticConfigs defines static Alertmanager targets.
	StaticConfigs []*StaticConfig `yaml:"static_configs,omitempty"`

	// Scheme defaults to 'http'.
	Scheme string `yaml:"scheme,omitempty"`

	// PathPrefix defines a path prefix to use when communicating with Alertmanager.
	PathPrefix string `yaml:"path_prefix,omitempty"`

	// Timeout is the timeout for requests to Alertmanager.
	Timeout Duration `yaml:"timeout,omitempty"`

	// APIVersion is the Alertmanager API version to use.
	APIVersion string `yaml:"api_version,omitempty"`
}

// RemoteWriteConfig configures remote write.
// Placeholder - will be fully implemented in Phase 2.
type RemoteWriteConfig struct {
	// URL specifies the endpoint to send samples to.
	URL string `yaml:"url"`

	// Name is an optional identifier for this remote write.
	Name string `yaml:"name,omitempty"`

	// RemoteTimeout is the timeout for requests to the remote write endpoint.
	RemoteTimeout Duration `yaml:"remote_timeout,omitempty"`

	// WriteRelabelConfigs is the list of relabel configurations for remote write.
	WriteRelabelConfigs []*RelabelConfig `yaml:"write_relabel_configs,omitempty"`
}

// RemoteReadConfig configures remote read.
// Placeholder - will be fully implemented in Phase 2.
type RemoteReadConfig struct {
	// URL specifies the endpoint to read samples from.
	URL string `yaml:"url"`

	// Name is an optional identifier for this remote read.
	Name string `yaml:"name,omitempty"`

	// RemoteTimeout is the timeout for requests to the remote read endpoint.
	RemoteTimeout Duration `yaml:"remote_timeout,omitempty"`

	// ReadRecent determines if remote read should only return data from the recent time range.
	ReadRecent bool `yaml:"read_recent,omitempty"`
}

// RelabelConfig configures target relabeling.
// Placeholder - will be fully implemented in Phase 2.
type RelabelConfig struct {
	// SourceLabels select values from existing labels for relabeling.
	SourceLabels []string `yaml:"source_labels,omitempty"`

	// Separator placed between concatenated source label values.
	Separator string `yaml:"separator,omitempty"`

	// Regex against which the extracted value is matched.
	Regex string `yaml:"regex,omitempty"`

	// Modulus to take of the hash of the source label values.
	Modulus uint64 `yaml:"modulus,omitempty"`

	// TargetLabel to which the resulting value is written.
	TargetLabel string `yaml:"target_label,omitempty"`

	// Replacement value against which a regex replace is performed.
	Replacement string `yaml:"replacement,omitempty"`

	// Action to perform based on regex matching.
	Action string `yaml:"action,omitempty"`
}
