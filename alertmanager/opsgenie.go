package alertmanager

// OpsGenieConfig configures OpsGenie notifications.
type OpsGenieConfig struct {
	// SendResolved determines if resolved alerts should be sent.
	SendResolved *bool `yaml:"send_resolved,omitempty"`

	// APIKey is the OpsGenie API key.
	APIKey Secret `yaml:"api_key,omitempty"`

	// APIKeyFile is a file containing the API key.
	APIKeyFile string `yaml:"api_key_file,omitempty"`

	// APIURL is the OpsGenie API URL.
	APIURL string `yaml:"api_url,omitempty"`

	// Message is the alert message.
	Message string `yaml:"message,omitempty"`

	// Description is the alert description.
	Description string `yaml:"description,omitempty"`

	// Source is the source of the alert.
	Source string `yaml:"source,omitempty"`

	// Details contains custom key/value pairs.
	Details map[string]string `yaml:"details,omitempty"`

	// Responders defines who to notify.
	Responders []*OpsGenieResponder `yaml:"responders,omitempty"`

	// Tags are alert tags.
	Tags []string `yaml:"tags,omitempty"`

	// Note is a note to add to the alert.
	Note string `yaml:"note,omitempty"`

	// Priority is the alert priority (P1-P5).
	Priority string `yaml:"priority,omitempty"`

	// Entity is the entity the alert is related to.
	Entity string `yaml:"entity,omitempty"`

	// Actions are custom actions.
	Actions []string `yaml:"actions,omitempty"`

	// UpdateAlerts updates existing alerts instead of creating new ones.
	UpdateAlerts *bool `yaml:"update_alerts,omitempty"`

	// HTTPConfig configures HTTP client settings.
	HTTPConfig *HTTPConfig `yaml:"http_config,omitempty"`
}

// OpsGenieResponder defines a responder for OpsGenie alerts.
type OpsGenieResponder struct {
	// Type is the responder type (team, user, escalation, schedule).
	Type string `yaml:"type"`

	// ID is the responder ID.
	ID string `yaml:"id,omitempty"`

	// Name is the responder name.
	Name string `yaml:"name,omitempty"`

	// Username is the responder username (for user type).
	Username string `yaml:"username,omitempty"`
}

// OpsGenie priority levels.
const (
	OpsGeniePriorityP1 = "P1"
	OpsGeniePriorityP2 = "P2"
	OpsGeniePriorityP3 = "P3"
	OpsGeniePriorityP4 = "P4"
	OpsGeniePriorityP5 = "P5"
)

// NewOpsGenieConfig creates a new OpsGenieConfig.
func NewOpsGenieConfig() *OpsGenieConfig {
	return &OpsGenieConfig{}
}

// WithSendResolved sets whether to send resolved alerts.
func (o *OpsGenieConfig) WithSendResolved(send bool) *OpsGenieConfig {
	o.SendResolved = &send
	return o
}

// WithAPIKey sets the OpsGenie API key.
func (o *OpsGenieConfig) WithAPIKey(key Secret) *OpsGenieConfig {
	o.APIKey = key
	return o
}

// WithAPIKeyFile sets the file containing the API key.
func (o *OpsGenieConfig) WithAPIKeyFile(path string) *OpsGenieConfig {
	o.APIKeyFile = path
	return o
}

// WithAPIURL sets the OpsGenie API URL.
func (o *OpsGenieConfig) WithAPIURL(url string) *OpsGenieConfig {
	o.APIURL = url
	return o
}

// WithMessage sets the alert message.
func (o *OpsGenieConfig) WithMessage(message string) *OpsGenieConfig {
	o.Message = message
	return o
}

// WithDescription sets the alert description.
func (o *OpsGenieConfig) WithDescription(desc string) *OpsGenieConfig {
	o.Description = desc
	return o
}

// WithSource sets the alert source.
func (o *OpsGenieConfig) WithSource(source string) *OpsGenieConfig {
	o.Source = source
	return o
}

// WithDetails sets custom key/value pairs.
func (o *OpsGenieConfig) WithDetails(details map[string]string) *OpsGenieConfig {
	o.Details = details
	return o
}

// WithResponders sets the responders.
func (o *OpsGenieConfig) WithResponders(responders ...*OpsGenieResponder) *OpsGenieConfig {
	o.Responders = responders
	return o
}

// WithTags sets alert tags.
func (o *OpsGenieConfig) WithTags(tags ...string) *OpsGenieConfig {
	o.Tags = tags
	return o
}

// WithNote sets a note for the alert.
func (o *OpsGenieConfig) WithNote(note string) *OpsGenieConfig {
	o.Note = note
	return o
}

// WithPriority sets the alert priority.
func (o *OpsGenieConfig) WithPriority(priority string) *OpsGenieConfig {
	o.Priority = priority
	return o
}

// WithEntity sets the entity the alert is related to.
func (o *OpsGenieConfig) WithEntity(entity string) *OpsGenieConfig {
	o.Entity = entity
	return o
}

// WithActions sets custom actions.
func (o *OpsGenieConfig) WithActions(actions ...string) *OpsGenieConfig {
	o.Actions = actions
	return o
}

// WithUpdateAlerts sets whether to update existing alerts.
func (o *OpsGenieConfig) WithUpdateAlerts(update bool) *OpsGenieConfig {
	o.UpdateAlerts = &update
	return o
}

// WithHTTPConfig sets HTTP client configuration.
func (o *OpsGenieConfig) WithHTTPConfig(config *HTTPConfig) *OpsGenieConfig {
	o.HTTPConfig = config
	return o
}

// NewOpsGenieResponder creates a new OpsGenieResponder.
func NewOpsGenieResponder(responderType, name string) *OpsGenieResponder {
	return &OpsGenieResponder{
		Type: responderType,
		Name: name,
	}
}

// WithID sets the responder ID.
func (r *OpsGenieResponder) WithID(id string) *OpsGenieResponder {
	r.ID = id
	return r
}

// WithUsername sets the responder username.
func (r *OpsGenieResponder) WithUsername(username string) *OpsGenieResponder {
	r.Username = username
	return r
}
