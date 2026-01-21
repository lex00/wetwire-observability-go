package monitoring

import "github.com/lex00/wetwire-observability-go/alertmanager"

// MuteIntervals define time periods when alerts should be muted.

// OutsideBusinessHours mutes alerts outside business hours (Mon-Fri 9am-6pm).
var OutsideBusinessHours = &alertmanager.MuteTimeInterval{
	Name: "outside-business-hours",
	TimeIntervals: []alertmanager.TimeInterval{
		{
			// Weekends
			Weekdays: []alertmanager.WeekdayRange{"saturday", "sunday"},
		},
		{
			// Before 9am on weekdays
			Weekdays: []alertmanager.WeekdayRange{"monday:friday"},
			Times: []alertmanager.TimeRange{
				{StartTime: "00:00", EndTime: "09:00"},
			},
		},
		{
			// After 6pm on weekdays
			Weekdays: []alertmanager.WeekdayRange{"monday:friday"},
			Times: []alertmanager.TimeRange{
				{StartTime: "18:00", EndTime: "24:00"},
			},
		},
	},
}

// BusinessHours defines business hours (Mon-Fri 9am-6pm).
var BusinessHours = &alertmanager.MuteTimeInterval{
	Name: "business-hours",
	TimeIntervals: []alertmanager.TimeInterval{
		{
			Weekdays: []alertmanager.WeekdayRange{"monday:friday"},
			Times: []alertmanager.TimeRange{
				{StartTime: "09:00", EndTime: "18:00"},
			},
		},
	},
}

// Weekends defines weekend hours.
var Weekends = &alertmanager.MuteTimeInterval{
	Name: "weekends",
	TimeIntervals: []alertmanager.TimeInterval{
		{
			Weekdays: []alertmanager.WeekdayRange{"saturday", "sunday"},
		},
	},
}

// MaintenanceWindow defines a recurring maintenance window (Sundays 2am-6am).
var MaintenanceWindow = &alertmanager.MuteTimeInterval{
	Name: "maintenance-window",
	TimeIntervals: []alertmanager.TimeInterval{
		{
			Weekdays: []alertmanager.WeekdayRange{"sunday"},
			Times: []alertmanager.TimeRange{
				{StartTime: "02:00", EndTime: "06:00"},
			},
		},
	},
}

// QuarterlyFreeze defines quarterly release freeze periods.
var QuarterlyFreeze = &alertmanager.MuteTimeInterval{
	Name: "quarterly-freeze",
	TimeIntervals: []alertmanager.TimeInterval{
		{
			// Last week of each quarter
			Months:      []alertmanager.MonthRange{"march", "june", "september", "december"},
			DaysOfMonth: []alertmanager.DayOfMonthRange{"25:31"},
		},
	},
}

// MuteIntervals is the list of all mute time intervals.
var MuteIntervals = []*alertmanager.MuteTimeInterval{
	OutsideBusinessHours,
	BusinessHours,
	Weekends,
	MaintenanceWindow,
	QuarterlyFreeze,
}
