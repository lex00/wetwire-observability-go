package prometheus

// EC2SD configures AWS EC2 service discovery.
// It discovers EC2 instances in an AWS region and configures them as targets.
//
// Example usage:
//
//	var EC2Discovery = prometheus.NewEC2SD().
//	    WithRegion("us-west-2").
//	    WithPort(9100).
//	    WithFilters(
//	        prometheus.EC2Filter{Name: "tag:Environment", Values: []string{"production"}},
//	    )
type EC2SD struct {
	// Region is the AWS region to discover EC2 instances in.
	// If empty, the region from the instance metadata is used.
	Region string `yaml:"region,omitempty"`

	// Endpoint is a custom API endpoint to use.
	// If empty, the standard AWS endpoint for the region is used.
	Endpoint string `yaml:"endpoint,omitempty"`

	// AccessKey is the AWS API access key.
	// If empty, the environment and shared credentials are used.
	AccessKey string `yaml:"access_key,omitempty"`

	// SecretKey is the AWS API secret key.
	// If empty, the environment and shared credentials are used.
	SecretKey string `yaml:"secret_key,omitempty"`

	// Profile is the named AWS profile to use for credentials.
	Profile string `yaml:"profile,omitempty"`

	// RoleARN is the AWS role to assume for API calls.
	RoleARN string `yaml:"role_arn,omitempty"`

	// RefreshInterval is the time after which the instance list is refreshed.
	// Defaults to 60s.
	RefreshInterval Duration `yaml:"refresh_interval,omitempty"`

	// Port is the port to scrape metrics from.
	// If not set, no port is used and must be determined via relabeling.
	Port int `yaml:"port,omitempty"`

	// Filters is a list of EC2 filters to apply when discovering instances.
	Filters []EC2Filter `yaml:"filters,omitempty"`
}

// EC2Filter defines an EC2 filter for instance discovery.
type EC2Filter struct {
	// Name is the filter name (e.g., "tag:Environment", "instance-state-name").
	Name string `yaml:"name"`

	// Values is the list of values to filter on.
	Values []string `yaml:"values"`
}

// NewEC2SD creates a new EC2 service discovery configuration.
func NewEC2SD() *EC2SD {
	return &EC2SD{}
}

// WithRegion sets the AWS region.
func (e *EC2SD) WithRegion(region string) *EC2SD {
	e.Region = region
	return e
}

// WithEndpoint sets a custom AWS API endpoint.
func (e *EC2SD) WithEndpoint(endpoint string) *EC2SD {
	e.Endpoint = endpoint
	return e
}

// WithAccessKey sets the AWS access key.
func (e *EC2SD) WithAccessKey(accessKey string) *EC2SD {
	e.AccessKey = accessKey
	return e
}

// WithSecretKey sets the AWS secret key.
func (e *EC2SD) WithSecretKey(secretKey string) *EC2SD {
	e.SecretKey = secretKey
	return e
}

// WithCredentials sets both access key and secret key.
func (e *EC2SD) WithCredentials(accessKey, secretKey string) *EC2SD {
	e.AccessKey = accessKey
	e.SecretKey = secretKey
	return e
}

// WithProfile sets the named AWS profile.
func (e *EC2SD) WithProfile(profile string) *EC2SD {
	e.Profile = profile
	return e
}

// WithRoleARN sets the AWS role ARN to assume.
func (e *EC2SD) WithRoleARN(roleARN string) *EC2SD {
	e.RoleARN = roleARN
	return e
}

// WithRefreshInterval sets the instance list refresh interval.
func (e *EC2SD) WithRefreshInterval(d Duration) *EC2SD {
	e.RefreshInterval = d
	return e
}

// WithPort sets the port to scrape metrics from.
func (e *EC2SD) WithPort(port int) *EC2SD {
	e.Port = port
	return e
}

// WithFilters sets the EC2 filters for instance discovery.
func (e *EC2SD) WithFilters(filters ...EC2Filter) *EC2SD {
	e.Filters = filters
	return e
}

// WithFilter adds a single EC2 filter.
func (e *EC2SD) WithFilter(name string, values ...string) *EC2SD {
	e.Filters = append(e.Filters, EC2Filter{
		Name:   name,
		Values: values,
	})
	return e
}

// WithTagFilter adds an EC2 filter for a specific tag key and values.
func (e *EC2SD) WithTagFilter(tagKey string, values ...string) *EC2SD {
	return e.WithFilter("tag:"+tagKey, values...)
}
