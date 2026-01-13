package operator

import "gopkg.in/yaml.v3"

// AlertmanagerConfig represents a Prometheus Operator AlertmanagerConfig CRD.
// This is a namespace-scoped resource for configuring Alertmanager.
type AlertmanagerConfig struct {
	APIVersion string                   `yaml:"apiVersion"`
	Kind       string                   `yaml:"kind"`
	Metadata   ObjectMeta               `yaml:"metadata"`
	Spec       AlertmanagerConfigSpec   `yaml:"spec"`

	// Convenience fields (not serialized)
	Name      string            `yaml:"-"`
	Namespace string            `yaml:"-"`
	Labels    map[string]string `yaml:"-"`
}

// AlertmanagerConfigSpec contains the AlertmanagerConfig specification.
type AlertmanagerConfigSpec struct {
	Route             *AMRoute             `yaml:"route,omitempty"`
	Receivers         []*AMReceiver        `yaml:"receivers,omitempty"`
	InhibitRules      []*AMInhibitRule     `yaml:"inhibitRules,omitempty"`
	MuteTimeIntervals []*AMMuteTimeInterval `yaml:"muteTimeIntervals,omitempty"`
}

// AMRoute represents a route in the AlertmanagerConfig.
type AMRoute struct {
	Receiver          string       `yaml:"receiver,omitempty"`
	GroupBy           []string     `yaml:"groupBy,omitempty"`
	GroupWait         string       `yaml:"groupWait,omitempty"`
	GroupInterval     string       `yaml:"groupInterval,omitempty"`
	RepeatInterval    string       `yaml:"repeatInterval,omitempty"`
	Matchers          []*AMMatcher `yaml:"matchers,omitempty"`
	Continue          bool         `yaml:"continue,omitempty"`
	Routes            []*AMRoute   `yaml:"routes,omitempty"`
	MuteTimeIntervals []string     `yaml:"muteTimeIntervals,omitempty"`
	ActiveTimeIntervals []string   `yaml:"activeTimeIntervals,omitempty"`
}

// AMMatcher represents a label matcher for AlertmanagerConfig.
type AMMatcher struct {
	Name      string `yaml:"name"`
	Value     string `yaml:"value"`
	MatchType string `yaml:"matchType,omitempty"`
	Regex     bool   `yaml:"regex,omitempty"`
}

// AMReceiver represents a receiver in the AlertmanagerConfig.
type AMReceiver struct {
	Name             string              `yaml:"name"`
	SlackConfigs     []*AMSlackConfig    `yaml:"slackConfigs,omitempty"`
	PagerDutyConfigs []*AMPagerDutyConfig `yaml:"pagerdutyConfigs,omitempty"`
	EmailConfigs     []*AMEmailConfig    `yaml:"emailConfigs,omitempty"`
	WebhookConfigs   []*AMWebhookConfig  `yaml:"webhookConfigs,omitempty"`
	OpsGenieConfigs  []*AMOpsGenieConfig `yaml:"opsgenieConfigs,omitempty"`
}

// AMSlackConfig configures Slack notifications for AlertmanagerConfig.
type AMSlackConfig struct {
	SendResolved *bool              `yaml:"sendResolved,omitempty"`
	APIURL       *SecretKeySelector `yaml:"apiURL,omitempty"`
	Channel      string             `yaml:"channel,omitempty"`
	Username     string             `yaml:"username,omitempty"`
	IconEmoji    string             `yaml:"iconEmoji,omitempty"`
	IconURL      string             `yaml:"iconURL,omitempty"`
	Title        string             `yaml:"title,omitempty"`
	TitleLink    string             `yaml:"titleLink,omitempty"`
	Text         string             `yaml:"text,omitempty"`
	Color        string             `yaml:"color,omitempty"`
	Footer       string             `yaml:"footer,omitempty"`
	Pretext      string             `yaml:"pretext,omitempty"`
	Fallback     string             `yaml:"fallback,omitempty"`
	HTTPConfig   *AMHTTPConfig      `yaml:"httpConfig,omitempty"`
}

// AMPagerDutyConfig configures PagerDuty notifications for AlertmanagerConfig.
type AMPagerDutyConfig struct {
	SendResolved *bool              `yaml:"sendResolved,omitempty"`
	RoutingKey   *SecretKeySelector `yaml:"routingKey,omitempty"`
	ServiceKey   *SecretKeySelector `yaml:"serviceKey,omitempty"`
	URL          string             `yaml:"url,omitempty"`
	Client       string             `yaml:"client,omitempty"`
	ClientURL    string             `yaml:"clientURL,omitempty"`
	Description  string             `yaml:"description,omitempty"`
	Severity     string             `yaml:"severity,omitempty"`
	Class        string             `yaml:"class,omitempty"`
	Group        string             `yaml:"group,omitempty"`
	Component    string             `yaml:"component,omitempty"`
	Details      []AMKeyValue       `yaml:"details,omitempty"`
	HTTPConfig   *AMHTTPConfig      `yaml:"httpConfig,omitempty"`
}

// AMEmailConfig configures email notifications for AlertmanagerConfig.
type AMEmailConfig struct {
	SendResolved *bool         `yaml:"sendResolved,omitempty"`
	To           string        `yaml:"to,omitempty"`
	From         string        `yaml:"from,omitempty"`
	Hello        string        `yaml:"hello,omitempty"`
	Smarthost    string        `yaml:"smarthost,omitempty"`
	AuthUsername string        `yaml:"authUsername,omitempty"`
	AuthPassword *SecretKeySelector `yaml:"authPassword,omitempty"`
	AuthSecret   *SecretKeySelector `yaml:"authSecret,omitempty"`
	AuthIdentity string        `yaml:"authIdentity,omitempty"`
	RequireTLS   *bool         `yaml:"requireTLS,omitempty"`
	TLSConfig    *AMTLSConfig  `yaml:"tlsConfig,omitempty"`
	HTML         string        `yaml:"html,omitempty"`
	Text         string        `yaml:"text,omitempty"`
	Headers      []AMKeyValue  `yaml:"headers,omitempty"`
}

// AMWebhookConfig configures webhook notifications for AlertmanagerConfig.
type AMWebhookConfig struct {
	SendResolved *bool              `yaml:"sendResolved,omitempty"`
	URL          *string            `yaml:"url,omitempty"`
	URLSecret    *SecretKeySelector `yaml:"urlSecret,omitempty"`
	HTTPConfig   *AMHTTPConfig      `yaml:"httpConfig,omitempty"`
	MaxAlerts    int                `yaml:"maxAlerts,omitempty"`
}

// AMOpsGenieConfig configures OpsGenie notifications for AlertmanagerConfig.
type AMOpsGenieConfig struct {
	SendResolved *bool              `yaml:"sendResolved,omitempty"`
	APIKey       *SecretKeySelector `yaml:"apiKey,omitempty"`
	APIURL       string             `yaml:"apiURL,omitempty"`
	Message      string             `yaml:"message,omitempty"`
	Description  string             `yaml:"description,omitempty"`
	Source       string             `yaml:"source,omitempty"`
	Tags         string             `yaml:"tags,omitempty"`
	Note         string             `yaml:"note,omitempty"`
	Priority     string             `yaml:"priority,omitempty"`
	Details      []AMKeyValue       `yaml:"details,omitempty"`
	Responders   []AMOpsGenieResponder `yaml:"responders,omitempty"`
	HTTPConfig   *AMHTTPConfig      `yaml:"httpConfig,omitempty"`
}

// AMOpsGenieResponder represents an OpsGenie responder.
type AMOpsGenieResponder struct {
	ID       string `yaml:"id,omitempty"`
	Name     string `yaml:"name,omitempty"`
	Username string `yaml:"username,omitempty"`
	Type     string `yaml:"type,omitempty"`
}

// AMKeyValue represents a key-value pair.
type AMKeyValue struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

// AMHTTPConfig configures HTTP client settings.
type AMHTTPConfig struct {
	BasicAuth       *AMBasicAuth       `yaml:"basicAuth,omitempty"`
	BearerTokenSecret *SecretKeySelector `yaml:"bearerTokenSecret,omitempty"`
	TLSConfig       *AMTLSConfig       `yaml:"tlsConfig,omitempty"`
	ProxyURL        string             `yaml:"proxyURL,omitempty"`
}

// AMBasicAuth configures basic HTTP authentication.
type AMBasicAuth struct {
	Username SecretKeySelector `yaml:"username"`
	Password SecretKeySelector `yaml:"password"`
}

// AMTLSConfig configures TLS settings.
type AMTLSConfig struct {
	CA                 SecretOrConfigMap  `yaml:"ca,omitempty"`
	Cert               SecretOrConfigMap  `yaml:"cert,omitempty"`
	KeySecret          *SecretKeySelector `yaml:"keySecret,omitempty"`
	ServerName         string             `yaml:"serverName,omitempty"`
	InsecureSkipVerify bool               `yaml:"insecureSkipVerify,omitempty"`
}

// AMInhibitRule defines an inhibition rule.
type AMInhibitRule struct {
	SourceMatch   []*AMMatcher `yaml:"sourceMatch,omitempty"`
	TargetMatch   []*AMMatcher `yaml:"targetMatch,omitempty"`
	Equal         []string     `yaml:"equal,omitempty"`
}

// AMMuteTimeInterval defines a named time interval for muting.
type AMMuteTimeInterval struct {
	Name          string          `yaml:"name"`
	TimeIntervals []AMTimeInterval `yaml:"timeIntervals,omitempty"`
}

// AMTimeInterval defines a time interval.
type AMTimeInterval struct {
	Times       []AMTimeRange `yaml:"times,omitempty"`
	Weekdays    []string      `yaml:"weekdays,omitempty"`
	DaysOfMonth []string      `yaml:"daysOfMonth,omitempty"`
	Months      []string      `yaml:"months,omitempty"`
	Years       []string      `yaml:"years,omitempty"`
}

// AMTimeRange defines a time range.
type AMTimeRange struct {
	StartTime string `yaml:"startTime,omitempty"`
	EndTime   string `yaml:"endTime,omitempty"`
}

// AMConfig creates a new AlertmanagerConfig.
func AMConfig(name, namespace string) *AlertmanagerConfig {
	return &AlertmanagerConfig{
		APIVersion: "monitoring.coreos.com/v1alpha1",
		Kind:       "AlertmanagerConfig",
		Metadata: ObjectMeta{
			Name:      name,
			Namespace: namespace,
		},
		Name:      name,
		Namespace: namespace,
	}
}

// WithLabels sets the AlertmanagerConfig labels.
func (am *AlertmanagerConfig) WithLabels(labels map[string]string) *AlertmanagerConfig {
	am.Labels = labels
	am.Metadata.Labels = labels
	return am
}

// AddLabel adds a label to the AlertmanagerConfig.
func (am *AlertmanagerConfig) AddLabel(key, value string) *AlertmanagerConfig {
	if am.Labels == nil {
		am.Labels = make(map[string]string)
	}
	am.Labels[key] = value
	am.Metadata.Labels = am.Labels
	return am
}

// WithAnnotations sets the AlertmanagerConfig annotations.
func (am *AlertmanagerConfig) WithAnnotations(annotations map[string]string) *AlertmanagerConfig {
	am.Metadata.Annotations = annotations
	return am
}

// WithRoute sets the route configuration.
func (am *AlertmanagerConfig) WithRoute(route *AMRoute) *AlertmanagerConfig {
	am.Spec.Route = route
	return am
}

// WithReceivers sets the receivers.
func (am *AlertmanagerConfig) WithReceivers(receivers ...*AMReceiver) *AlertmanagerConfig {
	am.Spec.Receivers = receivers
	return am
}

// AddReceiver adds a receiver.
func (am *AlertmanagerConfig) AddReceiver(receiver *AMReceiver) *AlertmanagerConfig {
	am.Spec.Receivers = append(am.Spec.Receivers, receiver)
	return am
}

// AddInhibitRule adds an inhibit rule.
func (am *AlertmanagerConfig) AddInhibitRule(rule *AMInhibitRule) *AlertmanagerConfig {
	am.Spec.InhibitRules = append(am.Spec.InhibitRules, rule)
	return am
}

// AddMuteTimeInterval adds a mute time interval.
func (am *AlertmanagerConfig) AddMuteTimeInterval(interval *AMMuteTimeInterval) *AlertmanagerConfig {
	am.Spec.MuteTimeIntervals = append(am.Spec.MuteTimeIntervals, interval)
	return am
}

// Serialize converts the AlertmanagerConfig to YAML bytes.
func (am *AlertmanagerConfig) Serialize() ([]byte, error) {
	return yaml.Marshal(am)
}

// MustSerialize converts the AlertmanagerConfig to YAML bytes, panicking on error.
func (am *AlertmanagerConfig) MustSerialize() []byte {
	data, err := am.Serialize()
	if err != nil {
		panic(err)
	}
	return data
}

// Route helpers

// NewAMRoute creates a new AMRoute.
func NewAMRoute(receiver string) *AMRoute {
	return &AMRoute{
		Receiver: receiver,
	}
}

// WithGroupBy sets the labels to group by.
func (r *AMRoute) WithGroupBy(labels ...string) *AMRoute {
	r.GroupBy = labels
	return r
}

// WithGroupWait sets the group wait duration.
func (r *AMRoute) WithGroupWait(duration string) *AMRoute {
	r.GroupWait = duration
	return r
}

// WithGroupInterval sets the group interval duration.
func (r *AMRoute) WithGroupInterval(duration string) *AMRoute {
	r.GroupInterval = duration
	return r
}

// WithRepeatInterval sets the repeat interval duration.
func (r *AMRoute) WithRepeatInterval(duration string) *AMRoute {
	r.RepeatInterval = duration
	return r
}

// WithMatchers sets the matchers.
func (r *AMRoute) WithMatchers(matchers ...*AMMatcher) *AMRoute {
	r.Matchers = matchers
	return r
}

// AddMatcher adds a matcher.
func (r *AMRoute) AddMatcher(matcher *AMMatcher) *AMRoute {
	r.Matchers = append(r.Matchers, matcher)
	return r
}

// WithContinue sets whether to continue to child routes.
func (r *AMRoute) WithContinue(cont bool) *AMRoute {
	r.Continue = cont
	return r
}

// WithRoutes sets child routes.
func (r *AMRoute) WithRoutes(routes ...*AMRoute) *AMRoute {
	r.Routes = routes
	return r
}

// AddRoute adds a child route.
func (r *AMRoute) AddRoute(route *AMRoute) *AMRoute {
	r.Routes = append(r.Routes, route)
	return r
}

// WithMuteTimeIntervals sets mute time interval names.
func (r *AMRoute) WithMuteTimeIntervals(intervals ...string) *AMRoute {
	r.MuteTimeIntervals = intervals
	return r
}

// Matcher helpers

// AMMatcherEq creates an equality matcher.
func AMMatcherEq(name, value string) *AMMatcher {
	return &AMMatcher{
		Name:      name,
		Value:     value,
		MatchType: "=",
	}
}

// AMMatcherNeq creates a not-equal matcher.
func AMMatcherNeq(name, value string) *AMMatcher {
	return &AMMatcher{
		Name:      name,
		Value:     value,
		MatchType: "!=",
	}
}

// AMMatcherRe creates a regex matcher.
func AMMatcherRe(name, value string) *AMMatcher {
	return &AMMatcher{
		Name:  name,
		Value: value,
		Regex: true,
	}
}

// Receiver helpers

// NewAMReceiver creates a new AMReceiver.
func NewAMReceiver(name string) *AMReceiver {
	return &AMReceiver{
		Name: name,
	}
}

// WithSlackConfig adds a Slack configuration.
func (r *AMReceiver) WithSlackConfig(config *AMSlackConfig) *AMReceiver {
	r.SlackConfigs = append(r.SlackConfigs, config)
	return r
}

// WithPagerDutyConfig adds a PagerDuty configuration.
func (r *AMReceiver) WithPagerDutyConfig(config *AMPagerDutyConfig) *AMReceiver {
	r.PagerDutyConfigs = append(r.PagerDutyConfigs, config)
	return r
}

// WithEmailConfig adds an email configuration.
func (r *AMReceiver) WithEmailConfig(config *AMEmailConfig) *AMReceiver {
	r.EmailConfigs = append(r.EmailConfigs, config)
	return r
}

// WithWebhookConfig adds a webhook configuration.
func (r *AMReceiver) WithWebhookConfig(config *AMWebhookConfig) *AMReceiver {
	r.WebhookConfigs = append(r.WebhookConfigs, config)
	return r
}

// WithOpsGenieConfig adds an OpsGenie configuration.
func (r *AMReceiver) WithOpsGenieConfig(config *AMOpsGenieConfig) *AMReceiver {
	r.OpsGenieConfigs = append(r.OpsGenieConfigs, config)
	return r
}

// Slack config helpers

// NewAMSlackConfig creates a new AMSlackConfig.
func NewAMSlackConfig() *AMSlackConfig {
	return &AMSlackConfig{}
}

// WithChannel sets the Slack channel.
func (s *AMSlackConfig) WithChannel(channel string) *AMSlackConfig {
	s.Channel = channel
	return s
}

// WithAPIURLSecret sets the API URL from a Kubernetes secret.
func (s *AMSlackConfig) WithAPIURLSecret(name, key string) *AMSlackConfig {
	s.APIURL = &SecretKeySelector{Name: name, Key: key}
	return s
}

// WithUsername sets the bot username.
func (s *AMSlackConfig) WithUsername(username string) *AMSlackConfig {
	s.Username = username
	return s
}

// WithIconEmoji sets the bot icon emoji.
func (s *AMSlackConfig) WithIconEmoji(emoji string) *AMSlackConfig {
	s.IconEmoji = emoji
	return s
}

// WithTitle sets the message title.
func (s *AMSlackConfig) WithTitle(title string) *AMSlackConfig {
	s.Title = title
	return s
}

// WithText sets the message text.
func (s *AMSlackConfig) WithText(text string) *AMSlackConfig {
	s.Text = text
	return s
}

// WithColor sets the attachment color.
func (s *AMSlackConfig) WithColor(color string) *AMSlackConfig {
	s.Color = color
	return s
}

// WithSendResolved sets whether to send resolved alerts.
func (s *AMSlackConfig) WithSendResolved(send bool) *AMSlackConfig {
	s.SendResolved = &send
	return s
}

// PagerDuty config helpers

// NewAMPagerDutyConfig creates a new AMPagerDutyConfig.
func NewAMPagerDutyConfig() *AMPagerDutyConfig {
	return &AMPagerDutyConfig{}
}

// WithRoutingKeySecret sets the routing key from a Kubernetes secret.
func (p *AMPagerDutyConfig) WithRoutingKeySecret(name, key string) *AMPagerDutyConfig {
	p.RoutingKey = &SecretKeySelector{Name: name, Key: key}
	return p
}

// WithServiceKeySecret sets the service key from a Kubernetes secret.
func (p *AMPagerDutyConfig) WithServiceKeySecret(name, key string) *AMPagerDutyConfig {
	p.ServiceKey = &SecretKeySelector{Name: name, Key: key}
	return p
}

// WithSeverity sets the PagerDuty severity.
func (p *AMPagerDutyConfig) WithSeverity(severity string) *AMPagerDutyConfig {
	p.Severity = severity
	return p
}

// WithDescription sets the description.
func (p *AMPagerDutyConfig) WithDescription(description string) *AMPagerDutyConfig {
	p.Description = description
	return p
}

// WithClass sets the class.
func (p *AMPagerDutyConfig) WithClass(class string) *AMPagerDutyConfig {
	p.Class = class
	return p
}

// WithSendResolved sets whether to send resolved alerts.
func (p *AMPagerDutyConfig) WithSendResolved(send bool) *AMPagerDutyConfig {
	p.SendResolved = &send
	return p
}

// Email config helpers

// NewAMEmailConfig creates a new AMEmailConfig.
func NewAMEmailConfig() *AMEmailConfig {
	return &AMEmailConfig{}
}

// WithTo sets the recipient email.
func (e *AMEmailConfig) WithTo(to string) *AMEmailConfig {
	e.To = to
	return e
}

// WithFrom sets the sender email.
func (e *AMEmailConfig) WithFrom(from string) *AMEmailConfig {
	e.From = from
	return e
}

// WithSmarthost sets the SMTP server.
func (e *AMEmailConfig) WithSmarthost(smarthost string) *AMEmailConfig {
	e.Smarthost = smarthost
	return e
}

// WithAuthUsername sets the auth username.
func (e *AMEmailConfig) WithAuthUsername(username string) *AMEmailConfig {
	e.AuthUsername = username
	return e
}

// WithAuthPasswordSecret sets the auth password from a Kubernetes secret.
func (e *AMEmailConfig) WithAuthPasswordSecret(name, key string) *AMEmailConfig {
	e.AuthPassword = &SecretKeySelector{Name: name, Key: key}
	return e
}

// WithSendResolved sets whether to send resolved alerts.
func (e *AMEmailConfig) WithSendResolved(send bool) *AMEmailConfig {
	e.SendResolved = &send
	return e
}

// Webhook config helpers

// NewAMWebhookConfig creates a new AMWebhookConfig.
func NewAMWebhookConfig() *AMWebhookConfig {
	return &AMWebhookConfig{}
}

// WithURL sets the webhook URL.
func (w *AMWebhookConfig) WithURL(url string) *AMWebhookConfig {
	w.URL = &url
	return w
}

// WithURLSecret sets the webhook URL from a Kubernetes secret.
func (w *AMWebhookConfig) WithURLSecret(name, key string) *AMWebhookConfig {
	w.URLSecret = &SecretKeySelector{Name: name, Key: key}
	return w
}

// WithMaxAlerts sets the maximum number of alerts.
func (w *AMWebhookConfig) WithMaxAlerts(max int) *AMWebhookConfig {
	w.MaxAlerts = max
	return w
}

// WithSendResolved sets whether to send resolved alerts.
func (w *AMWebhookConfig) WithSendResolved(send bool) *AMWebhookConfig {
	w.SendResolved = &send
	return w
}

// OpsGenie config helpers

// NewAMOpsGenieConfig creates a new AMOpsGenieConfig.
func NewAMOpsGenieConfig() *AMOpsGenieConfig {
	return &AMOpsGenieConfig{}
}

// WithAPIKeySecret sets the API key from a Kubernetes secret.
func (o *AMOpsGenieConfig) WithAPIKeySecret(name, key string) *AMOpsGenieConfig {
	o.APIKey = &SecretKeySelector{Name: name, Key: key}
	return o
}

// WithMessage sets the message.
func (o *AMOpsGenieConfig) WithMessage(message string) *AMOpsGenieConfig {
	o.Message = message
	return o
}

// WithPriority sets the priority.
func (o *AMOpsGenieConfig) WithPriority(priority string) *AMOpsGenieConfig {
	o.Priority = priority
	return o
}

// WithSendResolved sets whether to send resolved alerts.
func (o *AMOpsGenieConfig) WithSendResolved(send bool) *AMOpsGenieConfig {
	o.SendResolved = &send
	return o
}

// Inhibit rule helpers

// NewAMInhibitRule creates a new AMInhibitRule.
func NewAMInhibitRule() *AMInhibitRule {
	return &AMInhibitRule{}
}

// WithSourceMatcher adds a source matcher.
func (r *AMInhibitRule) WithSourceMatcher(matcher *AMMatcher) *AMInhibitRule {
	r.SourceMatch = append(r.SourceMatch, matcher)
	return r
}

// WithTargetMatcher adds a target matcher.
func (r *AMInhibitRule) WithTargetMatcher(matcher *AMMatcher) *AMInhibitRule {
	r.TargetMatch = append(r.TargetMatch, matcher)
	return r
}

// WithEqual sets the labels that must be equal.
func (r *AMInhibitRule) WithEqual(labels ...string) *AMInhibitRule {
	r.Equal = labels
	return r
}

// Mute time interval helpers

// NewAMMuteTimeInterval creates a new AMMuteTimeInterval.
func NewAMMuteTimeInterval(name string) *AMMuteTimeInterval {
	return &AMMuteTimeInterval{
		Name: name,
	}
}

// AddTimeInterval adds a time interval.
func (m *AMMuteTimeInterval) AddTimeInterval(interval AMTimeInterval) *AMMuteTimeInterval {
	m.TimeIntervals = append(m.TimeIntervals, interval)
	return m
}
