package prometheus

// DNSSDType represents the type of DNS record to query.
type DNSSDType string

// DNS service discovery record types.
const (
	// DNSSDTypeSRV queries SRV records for service discovery.
	// SRV records provide host, port, and priority information.
	DNSSDTypeSRV DNSSDType = "SRV"

	// DNSSDTypeA queries A records for service discovery.
	// A records provide IPv4 addresses. Port must be specified separately.
	DNSSDTypeA DNSSDType = "A"

	// DNSSDTypeAAAA queries AAAA records for service discovery.
	// AAAA records provide IPv6 addresses. Port must be specified separately.
	DNSSDTypeAAAA DNSSDType = "AAAA"

	// DNSSDTypeMX queries MX records for service discovery.
	// MX records provide mail exchange hosts.
	DNSSDTypeMX DNSSDType = "MX"

	// DNSSDTypeNS queries NS records for service discovery.
	// NS records provide nameserver hosts.
	DNSSDTypeNS DNSSDType = "NS"
)

// DNSSD configures DNS-based service discovery.
// It discovers targets by querying DNS records.
//
// Example usage:
//
//	var DNSDiscovery = prometheus.NewDNSSD().
//	    WithNames("_prometheus._tcp.example.com").
//	    WithType(prometheus.DNSSDTypeSRV)
type DNSSD struct {
	// Names is the list of DNS names to query.
	Names []string `yaml:"names"`

	// Type is the type of DNS record to query.
	// Defaults to "SRV".
	Type DNSSDType `yaml:"type,omitempty"`

	// Port is the port to use for targets discovered via A/AAAA records.
	// Required when Type is A or AAAA.
	// Ignored when Type is SRV (port from SRV record is used).
	Port int `yaml:"port,omitempty"`

	// RefreshInterval is the time after which DNS records are re-queried.
	// Defaults to 30s.
	RefreshInterval Duration `yaml:"refresh_interval,omitempty"`
}

// NewDNSSD creates a new DNS-based service discovery configuration.
func NewDNSSD() *DNSSD {
	return &DNSSD{}
}

// WithNames sets the DNS names to query.
func (d *DNSSD) WithNames(names ...string) *DNSSD {
	d.Names = names
	return d
}

// WithType sets the DNS record type to query.
func (d *DNSSD) WithType(t DNSSDType) *DNSSD {
	d.Type = t
	return d
}

// WithPort sets the port for A/AAAA record queries.
func (d *DNSSD) WithPort(port int) *DNSSD {
	d.Port = port
	return d
}

// WithRefreshInterval sets the DNS re-query interval.
func (d *DNSSD) WithRefreshInterval(dur Duration) *DNSSD {
	d.RefreshInterval = dur
	return d
}
