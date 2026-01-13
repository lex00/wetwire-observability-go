package prometheus

// KubernetesSD configures Kubernetes service discovery.
// It discovers targets from the Kubernetes API server based on the configured role.
//
// Example usage:
//
//	var PodDiscovery = &KubernetesSD{
//	    Role: KubernetesRolePod,
//	    Namespaces: &KubernetesNamespaceDiscovery{
//	        Names: []string{"production", "staging"},
//	    },
//	}
type KubernetesSD struct {
	// Role specifies the type of Kubernetes resource to discover.
	Role KubernetesRole `yaml:"role"`

	// Namespaces configures namespace filtering.
	// If not set, all namespaces are used.
	Namespaces *KubernetesNamespaceDiscovery `yaml:"namespaces,omitempty"`

	// Selectors specifies label and field selectors to filter resources.
	Selectors []KubernetesSelector `yaml:"selectors,omitempty"`

	// APIServer is the URL of the Kubernetes API server.
	// If not set, the in-cluster config is used.
	APIServer string `yaml:"api_server,omitempty"`

	// KubeConfigFile is the path to the kubeconfig file.
	// If not set, the in-cluster config or KUBECONFIG env var is used.
	KubeConfigFile string `yaml:"kubeconfig_file,omitempty"`

	// TLSConfig configures TLS settings for connecting to the API server.
	TLSConfig *TLSConfig `yaml:"tls_config,omitempty"`

	// BasicAuth configures basic auth for connecting to the API server.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty"`

	// BearerToken is the bearer token to use for authentication.
	BearerToken string `yaml:"bearer_token,omitempty"`

	// BearerTokenFile is the path to a file containing the bearer token.
	BearerTokenFile string `yaml:"bearer_token_file,omitempty"`

	// ProxyURL is the URL of a proxy to use for connecting to the API server.
	ProxyURL string `yaml:"proxy_url,omitempty"`

	// FollowRedirects controls whether to follow HTTP redirects.
	FollowRedirects *bool `yaml:"follow_redirects,omitempty"`

	// EnableHTTP2 controls whether to use HTTP/2.
	EnableHTTP2 *bool `yaml:"enable_http2,omitempty"`
}

// KubernetesRole represents the type of Kubernetes resource to discover.
type KubernetesRole string

// Kubernetes discovery roles.
const (
	// KubernetesRoleNode discovers one target per cluster node.
	// The __address__ label is set to the node's Kubelet address.
	KubernetesRoleNode KubernetesRole = "node"

	// KubernetesRolePod discovers all pods in the cluster.
	// Each pod container port is discovered as a target.
	KubernetesRolePod KubernetesRole = "pod"

	// KubernetesRoleService discovers all services in the cluster.
	// Each service port is discovered as a target.
	KubernetesRoleService KubernetesRole = "service"

	// KubernetesRoleEndpoints discovers endpoints of services.
	// For each endpoint address, one target is discovered per port.
	KubernetesRoleEndpoints KubernetesRole = "endpoints"

	// KubernetesRoleEndpointSlice discovers endpoints slices.
	// Similar to endpoints but uses the newer EndpointSlice API.
	KubernetesRoleEndpointSlice KubernetesRole = "endpointslice"

	// KubernetesRoleIngress discovers ingress resources.
	// Each ingress rule path is discovered as a target.
	KubernetesRoleIngress KubernetesRole = "ingress"
)

// KubernetesNamespaceDiscovery configures which namespaces to discover resources in.
type KubernetesNamespaceDiscovery struct {
	// OwnNamespace limits discovery to the namespace the Prometheus is running in.
	OwnNamespace bool `yaml:"own_namespace,omitempty"`

	// Names is the list of namespaces to discover resources in.
	// If not set and OwnNamespace is false, all namespaces are used.
	Names []string `yaml:"names,omitempty"`
}

// KubernetesSelector configures label and field selectors for a specific role.
type KubernetesSelector struct {
	// Role is the role that this selector applies to.
	Role KubernetesRole `yaml:"role"`

	// Label is a label selector string (e.g., "app=nginx,version!=v1").
	Label string `yaml:"label,omitempty"`

	// Field is a field selector string (e.g., "metadata.name=my-pod").
	Field string `yaml:"field,omitempty"`
}

// NewKubernetesSD creates a new Kubernetes service discovery configuration with the given role.
func NewKubernetesSD(role KubernetesRole) *KubernetesSD {
	return &KubernetesSD{
		Role: role,
	}
}

// WithNamespaces sets the namespaces to discover resources in.
func (k *KubernetesSD) WithNamespaces(namespaces ...string) *KubernetesSD {
	k.Namespaces = &KubernetesNamespaceDiscovery{
		Names: namespaces,
	}
	return k
}

// WithOwnNamespace limits discovery to the Prometheus namespace.
func (k *KubernetesSD) WithOwnNamespace() *KubernetesSD {
	k.Namespaces = &KubernetesNamespaceDiscovery{
		OwnNamespace: true,
	}
	return k
}

// WithLabelSelector adds a label selector for the specified role.
func (k *KubernetesSD) WithLabelSelector(role KubernetesRole, selector string) *KubernetesSD {
	k.Selectors = append(k.Selectors, KubernetesSelector{
		Role:  role,
		Label: selector,
	})
	return k
}

// WithFieldSelector adds a field selector for the specified role.
func (k *KubernetesSD) WithFieldSelector(role KubernetesRole, selector string) *KubernetesSD {
	k.Selectors = append(k.Selectors, KubernetesSelector{
		Role:  role,
		Field: selector,
	})
	return k
}

// WithAPIServer sets the Kubernetes API server URL.
func (k *KubernetesSD) WithAPIServer(url string) *KubernetesSD {
	k.APIServer = url
	return k
}

// WithKubeConfigFile sets the kubeconfig file path.
func (k *KubernetesSD) WithKubeConfigFile(path string) *KubernetesSD {
	k.KubeConfigFile = path
	return k
}

// WithBearerTokenFile sets the bearer token file for authentication.
func (k *KubernetesSD) WithBearerTokenFile(path string) *KubernetesSD {
	k.BearerTokenFile = path
	return k
}

// WithTLSConfig sets the TLS configuration.
func (k *KubernetesSD) WithTLSConfig(tls *TLSConfig) *KubernetesSD {
	k.TLSConfig = tls
	return k
}
