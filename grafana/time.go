// Package grafana provides types for Grafana dashboard configuration synthesis.
package grafana

// TimeRange represents the time range for a dashboard.
type TimeRange struct {
	// From is the start of the time range (e.g., "now-1h", "now-6h").
	From string `json:"from"`
	// To is the end of the time range (e.g., "now").
	To string `json:"to"`
}

// NewTimeRange creates a new TimeRange with the specified from and to values.
func NewTimeRange(from, to string) *TimeRange {
	return &TimeRange{From: from, To: to}
}

// LastHour returns a time range for the last hour.
func LastHour() *TimeRange {
	return NewTimeRange("now-1h", "now")
}

// Last6Hours returns a time range for the last 6 hours.
func Last6Hours() *TimeRange {
	return NewTimeRange("now-6h", "now")
}

// Last24Hours returns a time range for the last 24 hours.
func Last24Hours() *TimeRange {
	return NewTimeRange("now-24h", "now")
}

// Last7Days returns a time range for the last 7 days.
func Last7Days() *TimeRange {
	return NewTimeRange("now-7d", "now")
}

// Last30Days returns a time range for the last 30 days.
func Last30Days() *TimeRange {
	return NewTimeRange("now-30d", "now")
}

// Today returns a time range for today.
func Today() *TimeRange {
	return NewTimeRange("now/d", "now/d")
}

// ThisWeek returns a time range for this week.
func ThisWeek() *TimeRange {
	return NewTimeRange("now/w", "now/w")
}

// ThisMonth returns a time range for this month.
func ThisMonth() *TimeRange {
	return NewTimeRange("now/M", "now/M")
}

// DefaultTimeRange returns the default time range (last 6 hours).
func DefaultTimeRange() *TimeRange {
	return Last6Hours()
}
