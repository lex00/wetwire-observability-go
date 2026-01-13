package alertmanager

// PagerDutyConfig configures notifications to PagerDuty.
type PagerDutyConfig struct {
	// SendResolved determines if resolved alerts should be sent.
	SendResolved *bool `yaml:"send_resolved,omitempty"`

	// RoutingKey is the PagerDuty integration key (Events API v2).
	RoutingKey Secret `yaml:"routing_key,omitempty"`

	// RoutingKeyFile is a file containing the routing key.
	RoutingKeyFile string `yaml:"routing_key_file,omitempty"`

	// ServiceKey is the PagerDuty service key (Events API v1).
	// Deprecated: Use RoutingKey instead.
	ServiceKey Secret `yaml:"service_key,omitempty"`

	// ServiceKeyFile is a file containing the service key.
	// Deprecated: Use RoutingKeyFile instead.
	ServiceKeyFile string `yaml:"service_key_file,omitempty"`

	// URL is the PagerDuty API URL.
	URL string `yaml:"url,omitempty"`

	// Client is the name of the monitoring client.
	Client string `yaml:"client,omitempty"`

	// ClientURL is a URL to link to from PagerDuty.
	ClientURL string `yaml:"client_url,omitempty"`

	// Description is the incident description.
	Description string `yaml:"description,omitempty"`

	// Severity is the incident severity (critical, error, warning, info).
	Severity string `yaml:"severity,omitempty"`

	// Class is the incident class/type.
	Class string `yaml:"class,omitempty"`

	// Component is the affected component.
	Component string `yaml:"component,omitempty"`

	// Group is the logical grouping.
	Group string `yaml:"group,omitempty"`

	// Details contains custom key/value pairs.
	Details map[string]string `yaml:"details,omitempty"`

	// Images are images to include in the incident.
	Images []*PagerDutyImage `yaml:"images,omitempty"`

	// Links are links to include in the incident.
	Links []*PagerDutyLink `yaml:"links,omitempty"`

	// HTTPConfig configures HTTP client settings.
	HTTPConfig *HTTPConfig `yaml:"http_config,omitempty"`
}

// PagerDutyImage represents an image in a PagerDuty incident.
type PagerDutyImage struct {
	// Src is the image URL.
	Src string `yaml:"src"`

	// Href is a link URL for the image.
	Href string `yaml:"href,omitempty"`

	// Alt is the alt text for the image.
	Alt string `yaml:"alt,omitempty"`
}

// PagerDutyLink represents a link in a PagerDuty incident.
type PagerDutyLink struct {
	// Href is the link URL.
	Href string `yaml:"href"`

	// Text is the link text.
	Text string `yaml:"text,omitempty"`
}

// PagerDuty severity levels.
const (
	PagerDutySeverityCritical = "critical"
	PagerDutySeverityError    = "error"
	PagerDutySeverityWarning  = "warning"
	PagerDutySeverityInfo     = "info"
)

// NewPagerDutyConfig creates a new PagerDutyConfig.
func NewPagerDutyConfig() *PagerDutyConfig {
	return &PagerDutyConfig{}
}

// WithRoutingKey sets the PagerDuty routing key.
func (p *PagerDutyConfig) WithRoutingKey(key Secret) *PagerDutyConfig {
	p.RoutingKey = key
	return p
}

// WithRoutingKeyFile sets the file containing the routing key.
func (p *PagerDutyConfig) WithRoutingKeyFile(path string) *PagerDutyConfig {
	p.RoutingKeyFile = path
	return p
}

// WithServiceKey sets the PagerDuty service key (deprecated).
func (p *PagerDutyConfig) WithServiceKey(key Secret) *PagerDutyConfig {
	p.ServiceKey = key
	return p
}

// WithClient sets the monitoring client name.
func (p *PagerDutyConfig) WithClient(client string) *PagerDutyConfig {
	p.Client = client
	return p
}

// WithClientURL sets the monitoring client URL.
func (p *PagerDutyConfig) WithClientURL(url string) *PagerDutyConfig {
	p.ClientURL = url
	return p
}

// WithDescription sets the incident description.
func (p *PagerDutyConfig) WithDescription(desc string) *PagerDutyConfig {
	p.Description = desc
	return p
}

// WithSeverity sets the incident severity.
func (p *PagerDutyConfig) WithSeverity(severity string) *PagerDutyConfig {
	p.Severity = severity
	return p
}

// WithClass sets the incident class/type.
func (p *PagerDutyConfig) WithClass(class string) *PagerDutyConfig {
	p.Class = class
	return p
}

// WithComponent sets the affected component.
func (p *PagerDutyConfig) WithComponent(component string) *PagerDutyConfig {
	p.Component = component
	return p
}

// WithGroup sets the logical grouping.
func (p *PagerDutyConfig) WithGroup(group string) *PagerDutyConfig {
	p.Group = group
	return p
}

// WithDetails sets custom key/value pairs.
func (p *PagerDutyConfig) WithDetails(details map[string]string) *PagerDutyConfig {
	p.Details = details
	return p
}

// WithSendResolved sets whether to send resolved alerts.
func (p *PagerDutyConfig) WithSendResolved(send bool) *PagerDutyConfig {
	p.SendResolved = &send
	return p
}

// WithImages adds images to the incident.
func (p *PagerDutyConfig) WithImages(images ...*PagerDutyImage) *PagerDutyConfig {
	p.Images = images
	return p
}

// WithLinks adds links to the incident.
func (p *PagerDutyConfig) WithLinks(links ...*PagerDutyLink) *PagerDutyConfig {
	p.Links = links
	return p
}

// NewPagerDutyImage creates a new PagerDutyImage.
func NewPagerDutyImage(src string) *PagerDutyImage {
	return &PagerDutyImage{Src: src}
}

// WithHref sets the link URL for the image.
func (i *PagerDutyImage) WithHref(href string) *PagerDutyImage {
	i.Href = href
	return i
}

// WithAlt sets the alt text for the image.
func (i *PagerDutyImage) WithAlt(alt string) *PagerDutyImage {
	i.Alt = alt
	return i
}

// NewPagerDutyLink creates a new PagerDutyLink.
func NewPagerDutyLink(href string) *PagerDutyLink {
	return &PagerDutyLink{Href: href}
}

// WithText sets the link text.
func (l *PagerDutyLink) WithText(text string) *PagerDutyLink {
	l.Text = text
	return l
}
