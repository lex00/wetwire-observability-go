package alertmanager

// WithSlackConfigs adds Slack notification configurations.
func (r *Receiver) WithSlackConfigs(configs ...*SlackConfig) *Receiver {
	r.SlackConfigs = configs
	return r
}

// WithPagerDutyConfigs adds PagerDuty notification configurations.
func (r *Receiver) WithPagerDutyConfigs(configs ...*PagerDutyConfig) *Receiver {
	r.PagerDutyConfigs = configs
	return r
}

// WithEmailConfigs adds email notification configurations.
func (r *Receiver) WithEmailConfigs(configs ...*EmailConfig) *Receiver {
	r.EmailConfigs = configs
	return r
}

// WithWebhookConfigs adds webhook notification configurations.
func (r *Receiver) WithWebhookConfigs(configs ...*WebhookConfig) *Receiver {
	r.WebhookConfigs = configs
	return r
}

// WithOpsGenieConfigs adds OpsGenie notification configurations.
func (r *Receiver) WithOpsGenieConfigs(configs ...*OpsGenieConfig) *Receiver {
	r.OpsGenieConfigs = configs
	return r
}

// SlackReceiver creates a Receiver with a single Slack configuration.
func SlackReceiver(name, channel string) *Receiver {
	return NewReceiver(name).WithSlackConfigs(
		NewSlackConfig().WithChannel(channel),
	)
}

// PagerDutyReceiver creates a Receiver with a single PagerDuty configuration.
func PagerDutyReceiver(name string, routingKey Secret) *Receiver {
	return NewReceiver(name).WithPagerDutyConfigs(
		NewPagerDutyConfig().WithRoutingKey(routingKey),
	)
}

// PagerDutyReceiverFromFile creates a Receiver with PagerDuty config reading key from file.
func PagerDutyReceiverFromFile(name, keyFile string) *Receiver {
	return NewReceiver(name).WithPagerDutyConfigs(
		NewPagerDutyConfig().WithRoutingKeyFile(keyFile),
	)
}
