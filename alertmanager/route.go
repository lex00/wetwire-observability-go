package alertmanager

// Route defines a node in the Alertmanager routing tree.
// Routes form a tree structure where each node can specify matching rules
// and child routes for further routing.
type Route struct {
	// Receiver is the name of the receiver to send matching alerts to.
	// This is required for the root route.
	Receiver string `yaml:"receiver,omitempty"`

	// GroupBy specifies labels to group alerts by.
	// Alerts with the same values for these labels are batched together.
	GroupBy []string `yaml:"group_by,omitempty"`

	// GroupWait is how long to wait before sending the first notification
	// for a new group of alerts.
	GroupWait Duration `yaml:"group_wait,omitempty"`

	// GroupInterval is how long to wait before sending notifications
	// for new alerts added to a group that already sent an initial notification.
	GroupInterval Duration `yaml:"group_interval,omitempty"`

	// RepeatInterval is how long to wait before resending a notification
	// for a group that has already sent a notification.
	RepeatInterval Duration `yaml:"repeat_interval,omitempty"`

	// Matchers define which alerts this route matches.
	// Alerts must match ALL matchers to be routed here.
	Matchers []*Matcher `yaml:"matchers,omitempty"`

	// Continue indicates whether to continue matching child routes after
	// finding a match. If false, routing stops at the first matching child.
	Continue bool `yaml:"continue,omitempty"`

	// Routes are child routes for more specific matching.
	// Processed in order; first matching route wins (unless Continue is true).
	Routes []*Route `yaml:"routes,omitempty"`

	// MuteTimeIntervals is a list of mute time interval names
	// during which this route will be muted.
	MuteTimeIntervals []string `yaml:"mute_time_intervals,omitempty"`

	// ActiveTimeIntervals is a list of time interval names
	// during which this route will be active.
	ActiveTimeIntervals []string `yaml:"active_time_intervals,omitempty"`
}

// NewRoute creates a new Route with the given receiver.
func NewRoute(receiver string) *Route {
	return &Route{
		Receiver: receiver,
	}
}

// WithGroupBy sets the labels to group by.
func (r *Route) WithGroupBy(labels ...string) *Route {
	r.GroupBy = labels
	return r
}

// WithGroupWait sets the group wait duration.
func (r *Route) WithGroupWait(d Duration) *Route {
	r.GroupWait = d
	return r
}

// WithGroupInterval sets the group interval duration.
func (r *Route) WithGroupInterval(d Duration) *Route {
	r.GroupInterval = d
	return r
}

// WithRepeatInterval sets the repeat interval duration.
func (r *Route) WithRepeatInterval(d Duration) *Route {
	r.RepeatInterval = d
	return r
}

// WithMatchers sets the matchers for this route.
func (r *Route) WithMatchers(matchers ...*Matcher) *Route {
	r.Matchers = matchers
	return r
}

// WithContinue enables continuing to child routes after a match.
func (r *Route) WithContinue(cont bool) *Route {
	r.Continue = cont
	return r
}

// WithRoutes sets the child routes.
func (r *Route) WithRoutes(routes ...*Route) *Route {
	r.Routes = routes
	return r
}

// AddRoute adds a child route.
func (r *Route) AddRoute(route *Route) *Route {
	r.Routes = append(r.Routes, route)
	return r
}

// WithMuteTimeIntervals sets the mute time intervals.
func (r *Route) WithMuteTimeIntervals(intervals ...string) *Route {
	r.MuteTimeIntervals = intervals
	return r
}

// WithActiveTimeIntervals sets the active time intervals.
func (r *Route) WithActiveTimeIntervals(intervals ...string) *Route {
	r.ActiveTimeIntervals = intervals
	return r
}

// Severity creates a matcher for the severity label.
func (r *Route) Severity(severity string) *Route {
	r.Matchers = append(r.Matchers, Eq("severity", severity))
	return r
}

// Team creates a matcher for the team label.
func (r *Route) Team(team string) *Route {
	r.Matchers = append(r.Matchers, Eq("team", team))
	return r
}

// Service creates a matcher for the service label.
func (r *Route) Service(service string) *Route {
	r.Matchers = append(r.Matchers, Eq("service", service))
	return r
}

// Environment creates a matcher for the env/environment label.
func (r *Route) Environment(env string) *Route {
	r.Matchers = append(r.Matchers, Eq("env", env))
	return r
}
