package prometheus

// ScrapeConfig configures a scrape configuration for a set of targets.
type ScrapeConfig struct {
	// JobName is the job name for the scraped metrics.
	// The value is included as a `job` label on all scraped metrics.
	JobName string `yaml:"job_name"`

	// ScrapeInterval is the scrape interval for this job.
	// Overrides the global default if set.
	ScrapeInterval Duration `yaml:"scrape_interval,omitempty"`

	// ScrapeTimeout is the scrape timeout for this job.
	// Overrides the global default if set.
	ScrapeTimeout Duration `yaml:"scrape_timeout,omitempty"`

	// MetricsPath is the HTTP resource path for fetching metrics.
	// Defaults to '/metrics' if not specified.
	MetricsPath string `yaml:"metrics_path,omitempty"`

	// Scheme is the protocol scheme to use for requests.
	// Defaults to 'http' if not specified.
	Scheme string `yaml:"scheme,omitempty"`

	// HonorLabels controls how conflicts between server and scrape labels are handled.
	// If true, label conflicts are resolved by keeping label values from the scraped data.
	HonorLabels bool `yaml:"honor_labels,omitempty"`

	// HonorTimestamps controls whether Prometheus respects timestamps in scraped data.
	HonorTimestamps *bool `yaml:"honor_timestamps,omitempty"`

	// Params are optional HTTP URL parameters.
	Params map[string][]string `yaml:"params,omitempty"`

	// StaticConfigs define static targets with labels.
	StaticConfigs []*StaticConfig `yaml:"static_configs,omitempty"`

	// RelabelConfigs define relabeling rules applied before scraping.
	RelabelConfigs []*RelabelConfig `yaml:"relabel_configs,omitempty"`

	// MetricRelabelConfigs define relabeling rules applied to scraped samples before ingestion.
	MetricRelabelConfigs []*RelabelConfig `yaml:"metric_relabel_configs,omitempty"`

	// SampleLimit is the per-scrape limit on number of scraped samples.
	// If more samples are present after metric relabeling, the entire scrape is treated as failed.
	SampleLimit uint `yaml:"sample_limit,omitempty"`

	// TargetLimit is the per-scrape-config limit on number of unique targets.
	TargetLimit uint `yaml:"target_limit,omitempty"`

	// LabelLimit is the per-scrape limit on the number of labels.
	LabelLimit uint `yaml:"label_limit,omitempty"`

	// LabelNameLengthLimit is the per-scrape limit on the length of label names.
	LabelNameLengthLimit uint `yaml:"label_name_length_limit,omitempty"`

	// LabelValueLengthLimit is the per-scrape limit on the length of label values.
	LabelValueLengthLimit uint `yaml:"label_value_length_limit,omitempty"`

	// BasicAuth configures HTTP Basic authentication.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty"`

	// TLSConfig configures TLS settings.
	TLSConfig *TLSConfig `yaml:"tls_config,omitempty"`

	// ProxyURL is the proxy URL.
	ProxyURL string `yaml:"proxy_url,omitempty"`

	// Service discovery configurations will be added in Phase 2.
	// - KubernetesSDConfigs
	// - ConsulSDConfigs
	// - EC2SDConfigs
	// - FileSDConfigs
	// - DNSSDConfigs
}

// StaticConfig represents a static target group with an optional set of labels.
type StaticConfig struct {
	// Targets is the list of hosts to scrape.
	// Each target should be of the form host:port.
	Targets []string `yaml:"targets"`

	// Labels are labels assigned to all metrics scraped from the targets.
	Labels map[string]string `yaml:"labels,omitempty"`
}

// BasicAuth configures basic authentication.
type BasicAuth struct {
	// Username for basic authentication.
	Username string `yaml:"username,omitempty"`

	// Password for basic authentication.
	// Mutually exclusive with PasswordFile.
	Password string `yaml:"password,omitempty"`

	// PasswordFile is a file containing the password for basic authentication.
	// Mutually exclusive with Password.
	PasswordFile string `yaml:"password_file,omitempty"`
}

// TLSConfig configures TLS connections.
type TLSConfig struct {
	// CAFile is the path to the CA certificate file.
	CAFile string `yaml:"ca_file,omitempty"`

	// CertFile is the path to the client certificate file.
	CertFile string `yaml:"cert_file,omitempty"`

	// KeyFile is the path to the client key file.
	KeyFile string `yaml:"key_file,omitempty"`

	// ServerName is used to verify the hostname on the returned certificates.
	ServerName string `yaml:"server_name,omitempty"`

	// InsecureSkipVerify skips verifying the server's certificate chain and hostname.
	InsecureSkipVerify bool `yaml:"insecure_skip_verify,omitempty"`
}

// NewScrapeConfig creates a new ScrapeConfig with the given job name.
func NewScrapeConfig(jobName string) *ScrapeConfig {
	return &ScrapeConfig{
		JobName: jobName,
	}
}

// WithInterval sets the scrape interval.
func (s *ScrapeConfig) WithInterval(d Duration) *ScrapeConfig {
	s.ScrapeInterval = d
	return s
}

// WithTimeout sets the scrape timeout.
func (s *ScrapeConfig) WithTimeout(d Duration) *ScrapeConfig {
	s.ScrapeTimeout = d
	return s
}

// WithStaticTargets adds static targets to the scrape config.
func (s *ScrapeConfig) WithStaticTargets(targets ...string) *ScrapeConfig {
	s.StaticConfigs = append(s.StaticConfigs, &StaticConfig{
		Targets: targets,
	})
	return s
}

// NewStaticConfig creates a new StaticConfig with the given targets.
func NewStaticConfig(targets ...string) *StaticConfig {
	return &StaticConfig{
		Targets: targets,
	}
}

// WithLabels adds labels to the static config.
func (s *StaticConfig) WithLabels(labels map[string]string) *StaticConfig {
	s.Labels = labels
	return s
}
