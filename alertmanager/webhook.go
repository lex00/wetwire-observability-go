package alertmanager

// WebhookConfig configures webhook notifications.
type WebhookConfig struct {
	// SendResolved determines if resolved alerts should be sent.
	SendResolved *bool `yaml:"send_resolved,omitempty"`

	// URL is the endpoint to send HTTP requests to.
	URL string `yaml:"url,omitempty"`

	// URLFile is a file containing the webhook URL.
	URLFile string `yaml:"url_file,omitempty"`

	// MaxAlerts is the maximum number of alerts to include in a single webhook.
	MaxAlerts *int `yaml:"max_alerts,omitempty"`

	// HTTPConfig configures HTTP client settings.
	HTTPConfig *HTTPConfig `yaml:"http_config,omitempty"`
}

// NewWebhookConfig creates a new WebhookConfig.
func NewWebhookConfig() *WebhookConfig {
	return &WebhookConfig{}
}

// WithSendResolved sets whether to send resolved alerts.
func (w *WebhookConfig) WithSendResolved(send bool) *WebhookConfig {
	w.SendResolved = &send
	return w
}

// WithURL sets the webhook URL.
func (w *WebhookConfig) WithURL(url string) *WebhookConfig {
	w.URL = url
	return w
}

// WithURLFile sets the file containing the webhook URL.
func (w *WebhookConfig) WithURLFile(path string) *WebhookConfig {
	w.URLFile = path
	return w
}

// WithMaxAlerts sets the maximum number of alerts per webhook.
func (w *WebhookConfig) WithMaxAlerts(max int) *WebhookConfig {
	w.MaxAlerts = &max
	return w
}

// WithHTTPConfig sets HTTP client configuration.
func (w *WebhookConfig) WithHTTPConfig(config *HTTPConfig) *WebhookConfig {
	w.HTTPConfig = config
	return w
}
