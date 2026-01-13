package prometheus

// Secret represents a secret value that should be handled securely.
// When serialized, it outputs as a string, but can be loaded from various sources.
type Secret string

// SecretFromEnv creates a Secret reference from an environment variable.
// The actual value should be loaded at runtime.
func SecretFromEnv(envVar string) Secret {
	return Secret("${" + envVar + "}")
}

// QueueConfig configures the queue for remote write.
type QueueConfig struct {
	// Capacity is the number of samples to buffer per shard before dropping.
	Capacity int `yaml:"capacity,omitempty"`

	// MaxShards is the maximum number of shards.
	MaxShards int `yaml:"max_shards,omitempty"`

	// MinShards is the minimum number of shards.
	MinShards int `yaml:"min_shards,omitempty"`

	// MaxSamplesPerSend is the maximum number of samples per send.
	MaxSamplesPerSend int `yaml:"max_samples_per_send,omitempty"`

	// BatchSendDeadline is the maximum time a sample will wait in buffer.
	BatchSendDeadline Duration `yaml:"batch_send_deadline,omitempty"`

	// MinBackoff is the initial retry delay.
	MinBackoff Duration `yaml:"min_backoff,omitempty"`

	// MaxBackoff is the maximum retry delay.
	MaxBackoff Duration `yaml:"max_backoff,omitempty"`

	// RetryOnRateLimit enables retrying requests on HTTP 429 responses.
	RetryOnRateLimit bool `yaml:"retry_on_http_429,omitempty"`
}

// MetadataConfig configures how metadata is sent during remote write.
type MetadataConfig struct {
	// Send enables sending metric metadata.
	Send bool `yaml:"send"`

	// SendInterval is how frequently metric metadata is sent.
	SendInterval Duration `yaml:"send_interval,omitempty"`

	// MaxSamplesPerSend is the maximum number of metadata samples per send.
	MaxSamplesPerSend int `yaml:"max_samples_per_send,omitempty"`
}

// NewRemoteWrite creates a new RemoteWriteConfig with the given URL.
func NewRemoteWrite(url string) *RemoteWriteConfig {
	return &RemoteWriteConfig{
		URL: url,
	}
}

// WithName sets the name for this remote write configuration.
func (r *RemoteWriteConfig) WithName(name string) *RemoteWriteConfig {
	r.Name = name
	return r
}

// WithTimeout sets the remote timeout.
func (r *RemoteWriteConfig) WithTimeout(d Duration) *RemoteWriteConfig {
	r.RemoteTimeout = d
	return r
}

// WithHeaders adds custom headers to remote write requests.
func (r *RemoteWriteConfig) WithHeaders(headers map[string]string) *RemoteWriteConfig {
	r.Headers = headers
	return r
}

// WithWriteRelabelConfigs sets the write relabel configurations.
func (r *RemoteWriteConfig) WithWriteRelabelConfigs(configs ...*RelabelConfig) *RemoteWriteConfig {
	r.WriteRelabelConfigs = configs
	return r
}

// WithBasicAuth sets basic authentication for remote write.
func (r *RemoteWriteConfig) WithBasicAuth(username string, password Secret) *RemoteWriteConfig {
	r.BasicAuth = &BasicAuth{
		Username: username,
		Password: string(password),
	}
	return r
}

// WithBearerToken sets bearer token authentication.
func (r *RemoteWriteConfig) WithBearerToken(token Secret) *RemoteWriteConfig {
	r.BearerToken = token
	return r
}

// WithBearerTokenFile sets the file path for bearer token authentication.
func (r *RemoteWriteConfig) WithBearerTokenFile(path string) *RemoteWriteConfig {
	r.BearerTokenFile = path
	return r
}

// WithTLSConfig sets the TLS configuration.
func (r *RemoteWriteConfig) WithTLSConfig(tls *TLSConfig) *RemoteWriteConfig {
	r.TLSConfig = tls
	return r
}

// WithQueueConfig sets the queue configuration for buffering.
func (r *RemoteWriteConfig) WithQueueConfig(q *QueueConfig) *RemoteWriteConfig {
	r.QueueConfig = q
	return r
}

// WithMetadataConfig sets the metadata configuration.
func (r *RemoteWriteConfig) WithMetadataConfig(m *MetadataConfig) *RemoteWriteConfig {
	r.MetadataConfig = m
	return r
}

// WithProxyURL sets the proxy URL for remote write requests.
func (r *RemoteWriteConfig) WithProxyURL(url string) *RemoteWriteConfig {
	r.ProxyURL = url
	return r
}

// NewRemoteRead creates a new RemoteReadConfig with the given URL.
func NewRemoteRead(url string) *RemoteReadConfig {
	return &RemoteReadConfig{
		URL: url,
	}
}

// WithName sets the name for this remote read configuration.
func (r *RemoteReadConfig) WithName(name string) *RemoteReadConfig {
	r.Name = name
	return r
}

// WithTimeout sets the remote timeout.
func (r *RemoteReadConfig) WithTimeout(d Duration) *RemoteReadConfig {
	r.RemoteTimeout = d
	return r
}

// WithHeaders adds custom headers to remote read requests.
func (r *RemoteReadConfig) WithHeaders(headers map[string]string) *RemoteReadConfig {
	r.Headers = headers
	return r
}

// WithReadRecent enables reading only recent data.
func (r *RemoteReadConfig) WithReadRecent(recent bool) *RemoteReadConfig {
	r.ReadRecent = recent
	return r
}

// WithBasicAuth sets basic authentication for remote read.
func (r *RemoteReadConfig) WithBasicAuth(username string, password Secret) *RemoteReadConfig {
	r.BasicAuth = &BasicAuth{
		Username: username,
		Password: string(password),
	}
	return r
}

// WithBearerToken sets bearer token authentication.
func (r *RemoteReadConfig) WithBearerToken(token Secret) *RemoteReadConfig {
	r.BearerToken = token
	return r
}

// WithBearerTokenFile sets the file path for bearer token authentication.
func (r *RemoteReadConfig) WithBearerTokenFile(path string) *RemoteReadConfig {
	r.BearerTokenFile = path
	return r
}

// WithTLSConfig sets the TLS configuration.
func (r *RemoteReadConfig) WithTLSConfig(tls *TLSConfig) *RemoteReadConfig {
	r.TLSConfig = tls
	return r
}

// WithRequiredMatchers sets required label matchers that must be present in all queries.
func (r *RemoteReadConfig) WithRequiredMatchers(matchers map[string]string) *RemoteReadConfig {
	r.RequiredMatchers = matchers
	return r
}

// WithProxyURL sets the proxy URL for remote read requests.
func (r *RemoteReadConfig) WithProxyURL(url string) *RemoteReadConfig {
	r.ProxyURL = url
	return r
}

// WithFilterExternalLabels enables filtering of external labels.
func (r *RemoteReadConfig) WithFilterExternalLabels(filter bool) *RemoteReadConfig {
	r.FilterExternalLabels = filter
	return r
}

// NewQueueConfig creates a new QueueConfig with sensible defaults.
func NewQueueConfig() *QueueConfig {
	return &QueueConfig{}
}

// WithCapacity sets the queue capacity.
func (q *QueueConfig) WithCapacity(capacity int) *QueueConfig {
	q.Capacity = capacity
	return q
}

// WithMaxShards sets the maximum number of shards.
func (q *QueueConfig) WithMaxShards(shards int) *QueueConfig {
	q.MaxShards = shards
	return q
}

// WithMinShards sets the minimum number of shards.
func (q *QueueConfig) WithMinShards(shards int) *QueueConfig {
	q.MinShards = shards
	return q
}

// WithMaxSamplesPerSend sets the maximum samples per send.
func (q *QueueConfig) WithMaxSamplesPerSend(samples int) *QueueConfig {
	q.MaxSamplesPerSend = samples
	return q
}

// WithBatchSendDeadline sets the batch send deadline.
func (q *QueueConfig) WithBatchSendDeadline(d Duration) *QueueConfig {
	q.BatchSendDeadline = d
	return q
}

// WithBackoff sets the min and max backoff durations.
func (q *QueueConfig) WithBackoff(min, max Duration) *QueueConfig {
	q.MinBackoff = min
	q.MaxBackoff = max
	return q
}

// WithRetryOnRateLimit enables retrying on HTTP 429 responses.
func (q *QueueConfig) WithRetryOnRateLimit(retry bool) *QueueConfig {
	q.RetryOnRateLimit = retry
	return q
}

// NewMetadataConfig creates a new MetadataConfig with sending enabled.
func NewMetadataConfig() *MetadataConfig {
	return &MetadataConfig{
		Send: true,
	}
}

// WithSendInterval sets the send interval for metadata.
func (m *MetadataConfig) WithSendInterval(d Duration) *MetadataConfig {
	m.SendInterval = d
	return m
}

// WithMaxSamplesPerSend sets the maximum metadata samples per send.
func (m *MetadataConfig) WithMaxSamplesPerSend(samples int) *MetadataConfig {
	m.MaxSamplesPerSend = samples
	return m
}
