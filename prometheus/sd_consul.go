package prometheus

// ConsulSD configures Consul service discovery.
// It discovers targets from a Consul catalog based on services and tags.
//
// Example usage:
//
//	var ConsulDiscovery = prometheus.NewConsulSD().
//	    WithServer("consul.example.com:8500").
//	    WithServices("web", "api").
//	    WithTags("production")
type ConsulSD struct {
	// Server is the Consul server address.
	// The address should include the port (e.g., "localhost:8500").
	Server string `yaml:"server,omitempty"`

	// Token is the Consul ACL token.
	Token string `yaml:"token,omitempty"`

	// Datacenter is the Consul datacenter to query.
	// If not set, the datacenter of the Consul agent is used.
	Datacenter string `yaml:"datacenter,omitempty"`

	// Namespace is the Consul namespace to query (Enterprise only).
	Namespace string `yaml:"namespace,omitempty"`

	// Partition is the Consul admin partition to query (Enterprise only).
	Partition string `yaml:"partition,omitempty"`

	// Scheme is the protocol scheme to use for requests.
	// Defaults to "http".
	Scheme string `yaml:"scheme,omitempty"`

	// Services is the list of services for which targets are retrieved.
	// If empty, all services are scraped.
	Services []string `yaml:"services,omitempty"`

	// Tags is the list of tags used to filter nodes for a given service.
	// Services must contain all tags in the list.
	Tags []string `yaml:"tags,omitempty"`

	// NodeMeta is metadata key/value pairs used to filter nodes for a service.
	NodeMeta map[string]string `yaml:"node_meta,omitempty"`

	// TagSeparator is the string by which Consul tags are joined into the __meta_consul_tags label.
	// Defaults to ",".
	TagSeparator string `yaml:"tag_separator,omitempty"`

	// AllowStale allows any Consul server (non-leader) to service a read.
	// This means reads can be arbitrarily stale.
	AllowStale bool `yaml:"allow_stale,omitempty"`

	// RefreshInterval is the time after which the service list is refreshed.
	// Defaults to 30s.
	RefreshInterval Duration `yaml:"refresh_interval,omitempty"`

	// TLSConfig configures TLS settings for connecting to Consul.
	TLSConfig *TLSConfig `yaml:"tls_config,omitempty"`

	// BasicAuth configures basic authentication for connecting to Consul.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty"`

	// ProxyURL is the URL of a proxy to use for connecting to Consul.
	ProxyURL string `yaml:"proxy_url,omitempty"`

	// FollowRedirects controls whether to follow HTTP redirects.
	FollowRedirects *bool `yaml:"follow_redirects,omitempty"`

	// EnableHTTP2 controls whether to use HTTP/2.
	EnableHTTP2 *bool `yaml:"enable_http2,omitempty"`
}

// NewConsulSD creates a new Consul service discovery configuration.
func NewConsulSD() *ConsulSD {
	return &ConsulSD{}
}

// WithServer sets the Consul server address.
func (c *ConsulSD) WithServer(server string) *ConsulSD {
	c.Server = server
	return c
}

// WithToken sets the Consul ACL token.
func (c *ConsulSD) WithToken(token string) *ConsulSD {
	c.Token = token
	return c
}

// WithDatacenter sets the Consul datacenter.
func (c *ConsulSD) WithDatacenter(datacenter string) *ConsulSD {
	c.Datacenter = datacenter
	return c
}

// WithNamespace sets the Consul namespace (Enterprise only).
func (c *ConsulSD) WithNamespace(namespace string) *ConsulSD {
	c.Namespace = namespace
	return c
}

// WithPartition sets the Consul admin partition (Enterprise only).
func (c *ConsulSD) WithPartition(partition string) *ConsulSD {
	c.Partition = partition
	return c
}

// WithScheme sets the protocol scheme (http or https).
func (c *ConsulSD) WithScheme(scheme string) *ConsulSD {
	c.Scheme = scheme
	return c
}

// WithServices sets the services to discover.
func (c *ConsulSD) WithServices(services ...string) *ConsulSD {
	c.Services = services
	return c
}

// WithTags sets the tags used to filter services.
func (c *ConsulSD) WithTags(tags ...string) *ConsulSD {
	c.Tags = tags
	return c
}

// WithNodeMeta sets the node metadata key/value pairs for filtering.
func (c *ConsulSD) WithNodeMeta(meta map[string]string) *ConsulSD {
	c.NodeMeta = meta
	return c
}

// WithTagSeparator sets the tag separator string.
func (c *ConsulSD) WithTagSeparator(separator string) *ConsulSD {
	c.TagSeparator = separator
	return c
}

// WithAllowStale enables stale reads from any Consul server.
func (c *ConsulSD) WithAllowStale(allow bool) *ConsulSD {
	c.AllowStale = allow
	return c
}

// WithRefreshInterval sets the service list refresh interval.
func (c *ConsulSD) WithRefreshInterval(d Duration) *ConsulSD {
	c.RefreshInterval = d
	return c
}

// WithTLSConfig sets the TLS configuration.
func (c *ConsulSD) WithTLSConfig(tls *TLSConfig) *ConsulSD {
	c.TLSConfig = tls
	return c
}

// WithBasicAuth sets the basic authentication configuration.
func (c *ConsulSD) WithBasicAuth(auth *BasicAuth) *ConsulSD {
	c.BasicAuth = auth
	return c
}

// WithProxyURL sets the proxy URL.
func (c *ConsulSD) WithProxyURL(url string) *ConsulSD {
	c.ProxyURL = url
	return c
}
