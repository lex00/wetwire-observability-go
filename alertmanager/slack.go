package alertmanager

// SlackConfig configures notifications to Slack.
type SlackConfig struct {
	// SendResolved determines if resolved alerts should be sent.
	SendResolved *bool `yaml:"send_resolved,omitempty"`

	// APIURL is the Slack webhook URL.
	APIURL Secret `yaml:"api_url,omitempty"`

	// APIURLFile is a file containing the Slack webhook URL.
	APIURLFile string `yaml:"api_url_file,omitempty"`

	// Channel is the Slack channel to send messages to.
	// Overrides the channel defined in the webhook.
	Channel string `yaml:"channel,omitempty"`

	// Username is the bot username.
	Username string `yaml:"username,omitempty"`

	// IconEmoji is the emoji to use as the bot icon.
	IconEmoji string `yaml:"icon_emoji,omitempty"`

	// IconURL is a URL to an image to use as the bot icon.
	IconURL string `yaml:"icon_url,omitempty"`

	// Title is the message title.
	Title string `yaml:"title,omitempty"`

	// TitleLink is a URL to link the title to.
	TitleLink string `yaml:"title_link,omitempty"`

	// Pretext is text displayed before the message.
	Pretext string `yaml:"pretext,omitempty"`

	// Text is the main message text.
	Text string `yaml:"text,omitempty"`

	// Fallback is the text displayed in notifications.
	Fallback string `yaml:"fallback,omitempty"`

	// Color is the message attachment color.
	Color string `yaml:"color,omitempty"`

	// ShortFields determines if fields should be displayed side by side.
	ShortFields bool `yaml:"short_fields,omitempty"`

	// Footer is the message footer text.
	Footer string `yaml:"footer,omitempty"`

	// MrkdwnIn specifies which fields should be parsed as markdown.
	MrkdwnIn []string `yaml:"mrkdwn_in,omitempty"`

	// Actions are interactive buttons.
	Actions []*SlackAction `yaml:"actions,omitempty"`

	// Fields are additional fields to display.
	Fields []*SlackField `yaml:"fields,omitempty"`

	// ImageURL is a URL to an image to include in the message.
	ImageURL string `yaml:"image_url,omitempty"`

	// ThumbURL is a URL to a thumbnail image.
	ThumbURL string `yaml:"thumb_url,omitempty"`

	// HTTPConfig configures HTTP client settings.
	HTTPConfig *HTTPConfig `yaml:"http_config,omitempty"`
}

// SlackAction represents an interactive button in a Slack message.
type SlackAction struct {
	// Type is the action type (usually "button").
	Type string `yaml:"type,omitempty"`

	// Text is the button label.
	Text string `yaml:"text,omitempty"`

	// URL is the URL to open when clicked.
	URL string `yaml:"url,omitempty"`

	// Style is the button style (primary, danger).
	Style string `yaml:"style,omitempty"`

	// Name is the action name.
	Name string `yaml:"name,omitempty"`

	// Value is the action value.
	Value string `yaml:"value,omitempty"`
}

// SlackField represents a field in a Slack message attachment.
type SlackField struct {
	// Title is the field title.
	Title string `yaml:"title"`

	// Value is the field value.
	Value string `yaml:"value"`

	// Short determines if the field should be displayed side by side.
	Short *bool `yaml:"short,omitempty"`
}

// NewSlackConfig creates a new SlackConfig.
func NewSlackConfig() *SlackConfig {
	return &SlackConfig{}
}

// WithAPIURL sets the Slack webhook URL.
func (s *SlackConfig) WithAPIURL(url Secret) *SlackConfig {
	s.APIURL = url
	return s
}

// WithAPIURLFile sets the file containing the Slack webhook URL.
func (s *SlackConfig) WithAPIURLFile(path string) *SlackConfig {
	s.APIURLFile = path
	return s
}

// WithChannel sets the Slack channel.
func (s *SlackConfig) WithChannel(channel string) *SlackConfig {
	s.Channel = channel
	return s
}

// WithUsername sets the bot username.
func (s *SlackConfig) WithUsername(username string) *SlackConfig {
	s.Username = username
	return s
}

// WithIconEmoji sets the bot icon emoji.
func (s *SlackConfig) WithIconEmoji(emoji string) *SlackConfig {
	s.IconEmoji = emoji
	return s
}

// WithTitle sets the message title.
func (s *SlackConfig) WithTitle(title string) *SlackConfig {
	s.Title = title
	return s
}

// WithText sets the message text.
func (s *SlackConfig) WithText(text string) *SlackConfig {
	s.Text = text
	return s
}

// WithColor sets the attachment color.
func (s *SlackConfig) WithColor(color string) *SlackConfig {
	s.Color = color
	return s
}

// WithSendResolved sets whether to send resolved alerts.
func (s *SlackConfig) WithSendResolved(send bool) *SlackConfig {
	s.SendResolved = &send
	return s
}

// WithActions adds interactive actions (buttons).
func (s *SlackConfig) WithActions(actions ...*SlackAction) *SlackConfig {
	s.Actions = actions
	return s
}

// WithFields adds fields to the message.
func (s *SlackConfig) WithFields(fields ...*SlackField) *SlackConfig {
	s.Fields = fields
	return s
}

// NewSlackAction creates a new SlackAction button.
func NewSlackAction(text, url string) *SlackAction {
	return &SlackAction{
		Type: "button",
		Text: text,
		URL:  url,
	}
}

// WithStyle sets the button style (primary, danger).
func (a *SlackAction) WithStyle(style string) *SlackAction {
	a.Style = style
	return a
}

// NewSlackField creates a new SlackField.
func NewSlackField(title, value string) *SlackField {
	return &SlackField{
		Title: title,
		Value: value,
	}
}

// WithShort sets whether the field should be displayed side by side.
func (f *SlackField) WithShort(short bool) *SlackField {
	f.Short = &short
	return f
}

// SlackColorGood is the green color for Slack attachments.
const SlackColorGood = "good"

// SlackColorWarning is the yellow color for Slack attachments.
const SlackColorWarning = "warning"

// SlackColorDanger is the red color for Slack attachments.
const SlackColorDanger = "danger"
