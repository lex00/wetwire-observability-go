package alertmanager

// EmailConfig configures email notifications.
type EmailConfig struct {
	// SendResolved determines if resolved alerts should be sent.
	SendResolved *bool `yaml:"send_resolved,omitempty"`

	// To is the email address to send notifications to.
	To string `yaml:"to,omitempty"`

	// From is the sender's address.
	From string `yaml:"from,omitempty"`

	// Smarthost is the SMTP server address (host:port).
	Smarthost string `yaml:"smarthost,omitempty"`

	// Hello is the hostname to identify to the SMTP server.
	Hello string `yaml:"hello,omitempty"`

	// AuthUsername is the SMTP AUTH username.
	AuthUsername string `yaml:"auth_username,omitempty"`

	// AuthPassword is the SMTP AUTH password.
	AuthPassword Secret `yaml:"auth_password,omitempty"`

	// AuthSecret is the SMTP AUTH secret.
	AuthSecret Secret `yaml:"auth_secret,omitempty"`

	// AuthIdentity is the SMTP AUTH identity.
	AuthIdentity string `yaml:"auth_identity,omitempty"`

	// RequireTLS requires TLS connection.
	RequireTLS *bool `yaml:"require_tls,omitempty"`

	// TLSConfig configures TLS settings.
	TLSConfig *TLSConfig `yaml:"tls_config,omitempty"`

	// HTML is the HTML body of the email.
	HTML string `yaml:"html,omitempty"`

	// Text is the text body of the email.
	Text string `yaml:"text,omitempty"`

	// Headers contains additional email headers.
	Headers map[string]string `yaml:"headers,omitempty"`
}

// NewEmailConfig creates a new EmailConfig.
func NewEmailConfig() *EmailConfig {
	return &EmailConfig{}
}

// WithSendResolved sets whether to send resolved alerts.
func (e *EmailConfig) WithSendResolved(send bool) *EmailConfig {
	e.SendResolved = &send
	return e
}

// WithTo sets the recipient email address.
func (e *EmailConfig) WithTo(to string) *EmailConfig {
	e.To = to
	return e
}

// WithFrom sets the sender email address.
func (e *EmailConfig) WithFrom(from string) *EmailConfig {
	e.From = from
	return e
}

// WithSmarthost sets the SMTP server address.
func (e *EmailConfig) WithSmarthost(smarthost string) *EmailConfig {
	e.Smarthost = smarthost
	return e
}

// WithHello sets the hostname for SMTP HELLO.
func (e *EmailConfig) WithHello(hello string) *EmailConfig {
	e.Hello = hello
	return e
}

// WithAuthUsername sets the SMTP AUTH username.
func (e *EmailConfig) WithAuthUsername(username string) *EmailConfig {
	e.AuthUsername = username
	return e
}

// WithAuthPassword sets the SMTP AUTH password.
func (e *EmailConfig) WithAuthPassword(password Secret) *EmailConfig {
	e.AuthPassword = password
	return e
}

// WithAuthSecret sets the SMTP AUTH secret.
func (e *EmailConfig) WithAuthSecret(secret Secret) *EmailConfig {
	e.AuthSecret = secret
	return e
}

// WithAuthIdentity sets the SMTP AUTH identity.
func (e *EmailConfig) WithAuthIdentity(identity string) *EmailConfig {
	e.AuthIdentity = identity
	return e
}

// WithRequireTLS sets whether TLS is required.
func (e *EmailConfig) WithRequireTLS(require bool) *EmailConfig {
	e.RequireTLS = &require
	return e
}

// WithTLSConfig sets TLS configuration.
func (e *EmailConfig) WithTLSConfig(config *TLSConfig) *EmailConfig {
	e.TLSConfig = config
	return e
}

// WithHTML sets the HTML body template.
func (e *EmailConfig) WithHTML(html string) *EmailConfig {
	e.HTML = html
	return e
}

// WithText sets the text body template.
func (e *EmailConfig) WithText(text string) *EmailConfig {
	e.Text = text
	return e
}

// WithHeaders sets additional email headers.
func (e *EmailConfig) WithHeaders(headers map[string]string) *EmailConfig {
	e.Headers = headers
	return e
}
