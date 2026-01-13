// Package operator provides Prometheus Operator CRD types.
package operator

// Endpoint represents a ServiceMonitor endpoint configuration.
type Endpoint struct {
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

	// BearerTokenFile is the file containing the bearer token.
	BearerTokenFile string `yaml:"bearerTokenFile,omitempty"`

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

// TLSConfig contains TLS configuration.
type TLSConfig struct {
	// CAFile is the path to the CA certificate file.
	CAFile string `yaml:"caFile,omitempty"`

	// CA is a reference to a secret containing the CA cert.
	CA *SecretOrConfigMap `yaml:"ca,omitempty"`

	// CertFile is the path to the client certificate file.
	CertFile string `yaml:"certFile,omitempty"`

	// Cert is a reference to a secret containing the client cert.
	Cert *SecretOrConfigMap `yaml:"cert,omitempty"`

	// KeyFile is the path to the client key file.
	KeyFile string `yaml:"keyFile,omitempty"`

	// KeySecret is a reference to a secret containing the client key.
	KeySecret *SecretKeySelector `yaml:"keySecret,omitempty"`

	// ServerName is the expected server name for TLS verification.
	ServerName string `yaml:"serverName,omitempty"`

	// InsecureSkipVerify skips TLS certificate verification.
	InsecureSkipVerify bool `yaml:"insecureSkipVerify,omitempty"`
}

// SecretOrConfigMap references a secret or configmap.
type SecretOrConfigMap struct {
	Secret    *SecretKeySelector    `yaml:"secret,omitempty"`
	ConfigMap *ConfigMapKeySelector `yaml:"configMap,omitempty"`
}

// SecretKeySelector selects a key from a secret.
type SecretKeySelector struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

// ConfigMapKeySelector selects a key from a configmap.
type ConfigMapKeySelector struct {
	Name string `yaml:"name"`
	Key  string `yaml:"key"`
}

// BasicAuth contains basic auth credentials.
type BasicAuth struct {
	Username SecretKeySelector `yaml:"username"`
	Password SecretKeySelector `yaml:"password"`
}

// RelabelConfig contains relabeling configuration.
type RelabelConfig struct {
	SourceLabels []string `yaml:"sourceLabels,omitempty"`
	Separator    string   `yaml:"separator,omitempty"`
	TargetLabel  string   `yaml:"targetLabel,omitempty"`
	Regex        string   `yaml:"regex,omitempty"`
	Modulus      uint64   `yaml:"modulus,omitempty"`
	Replacement  string   `yaml:"replacement,omitempty"`
	Action       string   `yaml:"action,omitempty"`
}

// NewEndpoint creates a new endpoint with the given port.
func NewEndpoint(port string) *Endpoint {
	return &Endpoint{
		Port: port,
	}
}

// WithPath sets the scrape path.
func (e *Endpoint) WithPath(path string) *Endpoint {
	e.Path = path
	return e
}

// WithInterval sets the scrape interval.
func (e *Endpoint) WithInterval(interval string) *Endpoint {
	e.Interval = interval
	return e
}

// WithScrapeTimeout sets the scrape timeout.
func (e *Endpoint) WithScrapeTimeout(timeout string) *Endpoint {
	e.ScrapeTimeout = timeout
	return e
}

// WithScheme sets the URL scheme.
func (e *Endpoint) WithScheme(scheme string) *Endpoint {
	e.Scheme = scheme
	return e
}

// WithBearerTokenFile sets the bearer token file path.
func (e *Endpoint) WithBearerTokenFile(file string) *Endpoint {
	e.BearerTokenFile = file
	return e
}

// WithBearerTokenSecret sets the bearer token secret reference.
func (e *Endpoint) WithBearerTokenSecret(name, key string) *Endpoint {
	e.BearerTokenSecret = &SecretKeySelector{Name: name, Key: key}
	return e
}

// WithTLSConfig sets the TLS configuration.
func (e *Endpoint) WithTLSConfig(insecureSkipVerify bool, caFile, certFile, keyFile string) *Endpoint {
	e.TLSConfig = &TLSConfig{
		InsecureSkipVerify: insecureSkipVerify,
		CAFile:             caFile,
		CertFile:           certFile,
		KeyFile:            keyFile,
	}
	return e
}

// WithBasicAuth sets the basic auth credentials.
func (e *Endpoint) WithBasicAuth(usernameSecret, usernameKey, passwordSecret, passwordKey string) *Endpoint {
	e.BasicAuth = &BasicAuth{
		Username: SecretKeySelector{Name: usernameSecret, Key: usernameKey},
		Password: SecretKeySelector{Name: passwordSecret, Key: passwordKey},
	}
	return e
}

// HonorLabelsOn honors labels from the scraped target.
func (e *Endpoint) HonorLabelsOn() *Endpoint {
	e.HonorLabels = true
	return e
}

// AddRelabeling adds a relabeling configuration.
func (e *Endpoint) AddRelabeling(r *RelabelConfig) *Endpoint {
	e.RelabelConfigs = append(e.RelabelConfigs, r)
	return e
}

// AddMetricRelabeling adds a metric relabeling configuration.
func (e *Endpoint) AddMetricRelabeling(r *RelabelConfig) *Endpoint {
	e.MetricRelabelConfigs = append(e.MetricRelabelConfigs, r)
	return e
}

// WithParam adds a URL parameter.
func (e *Endpoint) WithParam(key string, values ...string) *Endpoint {
	if e.Params == nil {
		e.Params = make(map[string][]string)
	}
	e.Params[key] = values
	return e
}

// Relabeling helper functions

// KeepLabel creates a keep action relabel config.
func KeepLabel(sourceLabel string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{sourceLabel},
		Action:       "keep",
	}
}

// DropLabel creates a drop action relabel config.
func DropLabel(regex string) *RelabelConfig {
	return &RelabelConfig{
		Regex:  regex,
		Action: "drop",
	}
}

// ReplaceLabel creates a replace action relabel config.
func ReplaceLabel(sourceLabel, targetLabel string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{sourceLabel},
		TargetLabel:  targetLabel,
		Action:       "replace",
	}
}

// LabelMap creates a labelmap action relabel config.
func LabelMap(regex, replacement string) *RelabelConfig {
	return &RelabelConfig{
		Regex:       regex,
		Replacement: replacement,
		Action:      "labelmap",
	}
}

// DropMetric creates a drop action for metrics matching regex.
func DropMetric(regex string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{"__name__"},
		Regex:        regex,
		Action:       "drop",
	}
}

// KeepMetric creates a keep action for metrics matching regex.
func KeepMetric(regex string) *RelabelConfig {
	return &RelabelConfig{
		SourceLabels: []string{"__name__"},
		Regex:        regex,
		Action:       "keep",
	}
}
