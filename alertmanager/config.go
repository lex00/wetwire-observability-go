// Package alertmanager provides types for Alertmanager configuration synthesis.
package alertmanager

import "github.com/lex00/wetwire-observability-go/prometheus"

// Duration is an alias for prometheus.Duration for consistency.
type Duration = prometheus.Duration

// Convenience duration constants.
const (
	Second = prometheus.Second
	Minute = prometheus.Minute
	Hour   = prometheus.Hour
)

// AlertmanagerConfig represents the top-level Alertmanager configuration.
// This is the root structure that serializes to alertmanager.yml.
type AlertmanagerConfig struct {
	// Global contains global configuration options.
	Global *GlobalConfig `yaml:"global,omitempty"`

	// Route is the root of the routing tree.
	Route *Route `yaml:"route,omitempty"`

	// Receivers defines the notification integrations.
	Receivers []*Receiver `yaml:"receivers,omitempty"`

	// InhibitRules defines inhibition rules.
	InhibitRules []*InhibitRule `yaml:"inhibit_rules,omitempty"`

	// MuteTimeIntervals defines named time intervals for muting.
	MuteTimeIntervals []*MuteTimeInterval `yaml:"mute_time_intervals,omitempty"`

	// Templates defines paths to template files.
	Templates []string `yaml:"templates,omitempty"`
}

// GlobalConfig contains global alertmanager settings.
type GlobalConfig struct {
	// SMTPSmarthost is the SMTP server address.
	SMTPSmarthost string `yaml:"smtp_smarthost,omitempty"`

	// SMTPFrom is the sender address.
	SMTPFrom string `yaml:"smtp_from,omitempty"`

	// SMTPAuthUsername is the SMTP auth username.
	SMTPAuthUsername string `yaml:"smtp_auth_username,omitempty"`

	// SMTPAuthPassword is the SMTP auth password.
	SMTPAuthPassword string `yaml:"smtp_auth_password,omitempty"`

	// SMTPAuthSecret is the SMTP auth secret.
	SMTPAuthSecret string `yaml:"smtp_auth_secret,omitempty"`

	// SMTPAuthIdentity is the SMTP auth identity.
	SMTPAuthIdentity string `yaml:"smtp_auth_identity,omitempty"`

	// SMTPRequireTLS requires TLS.
	SMTPRequireTLS *bool `yaml:"smtp_require_tls,omitempty"`

	// SlackAPIURL is the default Slack API URL.
	SlackAPIURL string `yaml:"slack_api_url,omitempty"`

	// PagerDutyURL is the PagerDuty URL.
	PagerDutyURL string `yaml:"pagerduty_url,omitempty"`

	// OpsGenieAPIURL is the OpsGenie API URL.
	OpsGenieAPIURL string `yaml:"opsgenie_api_url,omitempty"`

	// OpsGenieAPIKey is the OpsGenie API key.
	OpsGenieAPIKey string `yaml:"opsgenie_api_key,omitempty"`

	// HTTPConfig sets HTTP client configuration.
	HTTPConfig *HTTPConfig `yaml:"http_config,omitempty"`

	// ResolveTimeout is the default resolve timeout.
	ResolveTimeout Duration `yaml:"resolve_timeout,omitempty"`
}

// HTTPConfig configures HTTP client settings.
type HTTPConfig struct {
	// BasicAuth configures basic authentication.
	BasicAuth *BasicAuth `yaml:"basic_auth,omitempty"`

	// BearerToken is the bearer token for authentication.
	BearerToken string `yaml:"bearer_token,omitempty"`

	// BearerTokenFile is the path to a file containing the bearer token.
	BearerTokenFile string `yaml:"bearer_token_file,omitempty"`

	// TLSConfig configures TLS settings.
	TLSConfig *TLSConfig `yaml:"tls_config,omitempty"`

	// ProxyURL is the HTTP proxy URL.
	ProxyURL string `yaml:"proxy_url,omitempty"`
}

// BasicAuth configures basic HTTP authentication.
type BasicAuth struct {
	Username     string `yaml:"username,omitempty"`
	Password     string `yaml:"password,omitempty"`
	PasswordFile string `yaml:"password_file,omitempty"`
}

// TLSConfig configures TLS settings.
type TLSConfig struct {
	CAFile             string `yaml:"ca_file,omitempty"`
	CertFile           string `yaml:"cert_file,omitempty"`
	KeyFile            string `yaml:"key_file,omitempty"`
	ServerName         string `yaml:"server_name,omitempty"`
	InsecureSkipVerify bool   `yaml:"insecure_skip_verify,omitempty"`
}

// Receiver defines a notification integration.
// Placeholder - will be expanded in receiver files.
type Receiver struct {
	// Name is the unique identifier for this receiver.
	Name string `yaml:"name"`

	// EmailConfigs defines email notification targets.
	EmailConfigs []*EmailConfig `yaml:"email_configs,omitempty"`

	// SlackConfigs defines Slack notification targets.
	SlackConfigs []*SlackConfig `yaml:"slack_configs,omitempty"`

	// PagerDutyConfigs defines PagerDuty notification targets.
	PagerDutyConfigs []*PagerDutyConfig `yaml:"pagerduty_configs,omitempty"`

	// WebhookConfigs defines webhook notification targets.
	WebhookConfigs []*WebhookConfig `yaml:"webhook_configs,omitempty"`

	// OpsGenieConfigs defines OpsGenie notification targets.
	OpsGenieConfigs []*OpsGenieConfig `yaml:"opsgenie_configs,omitempty"`
}

// InhibitRule defines a rule for muting alerts.
// Placeholder - will be expanded later.
type InhibitRule struct {
	SourceMatch   map[string]string `yaml:"source_match,omitempty"`
	SourceMatchers []*Matcher       `yaml:"source_matchers,omitempty"`
	TargetMatch   map[string]string `yaml:"target_match,omitempty"`
	TargetMatchers []*Matcher       `yaml:"target_matchers,omitempty"`
	Equal         []string          `yaml:"equal,omitempty"`
}

// MuteTimeInterval defines a named time interval for muting.
// Placeholder - will be expanded later.
type MuteTimeInterval struct {
	Name          string         `yaml:"name"`
	TimeIntervals []TimeInterval `yaml:"time_intervals,omitempty"`
}

// TimeInterval defines a time interval.
type TimeInterval struct {
	Times       []TimeRange       `yaml:"times,omitempty"`
	Weekdays    []WeekdayRange    `yaml:"weekdays,omitempty"`
	DaysOfMonth []DayOfMonthRange `yaml:"days_of_month,omitempty"`
	Months      []MonthRange      `yaml:"months,omitempty"`
	Years       []YearRange       `yaml:"years,omitempty"`
}

// Placeholder range types - will be expanded later.
type TimeRange struct {
	StartTime string `yaml:"start_time,omitempty"`
	EndTime   string `yaml:"end_time,omitempty"`
}

type WeekdayRange string
type DayOfMonthRange string
type MonthRange string
type YearRange string

// NewAlertmanagerConfig creates a new AlertmanagerConfig.
func NewAlertmanagerConfig() *AlertmanagerConfig {
	return &AlertmanagerConfig{}
}

// WithGlobal sets the global configuration.
func (c *AlertmanagerConfig) WithGlobal(g *GlobalConfig) *AlertmanagerConfig {
	c.Global = g
	return c
}

// WithRoute sets the root route.
func (c *AlertmanagerConfig) WithRoute(r *Route) *AlertmanagerConfig {
	c.Route = r
	return c
}

// WithReceivers sets the receivers.
func (c *AlertmanagerConfig) WithReceivers(receivers ...*Receiver) *AlertmanagerConfig {
	c.Receivers = receivers
	return c
}

// WithInhibitRules sets the inhibit rules.
func (c *AlertmanagerConfig) WithInhibitRules(rules ...*InhibitRule) *AlertmanagerConfig {
	c.InhibitRules = rules
	return c
}

// WithMuteTimeIntervals sets the mute time intervals.
func (c *AlertmanagerConfig) WithMuteTimeIntervals(intervals ...*MuteTimeInterval) *AlertmanagerConfig {
	c.MuteTimeIntervals = intervals
	return c
}

// WithTemplates sets the template paths.
func (c *AlertmanagerConfig) WithTemplates(templates ...string) *AlertmanagerConfig {
	c.Templates = templates
	return c
}

// NewGlobalConfig creates a new GlobalConfig.
func NewGlobalConfig() *GlobalConfig {
	return &GlobalConfig{}
}

// WithSMTP sets SMTP settings.
func (g *GlobalConfig) WithSMTP(smarthost, from string) *GlobalConfig {
	g.SMTPSmarthost = smarthost
	g.SMTPFrom = from
	return g
}

// WithSMTPAuth sets SMTP authentication.
func (g *GlobalConfig) WithSMTPAuth(username, password string) *GlobalConfig {
	g.SMTPAuthUsername = username
	g.SMTPAuthPassword = password
	return g
}

// WithSlackAPIURL sets the default Slack API URL.
func (g *GlobalConfig) WithSlackAPIURL(url string) *GlobalConfig {
	g.SlackAPIURL = url
	return g
}

// WithResolveTimeout sets the resolve timeout.
func (g *GlobalConfig) WithResolveTimeout(d Duration) *GlobalConfig {
	g.ResolveTimeout = d
	return g
}

// NewReceiver creates a new Receiver with the given name.
func NewReceiver(name string) *Receiver {
	return &Receiver{Name: name}
}
