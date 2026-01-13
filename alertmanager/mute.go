package alertmanager

import "fmt"

// Weekday constants.
const (
	Sunday    WeekdayRange = "sunday"
	Monday    WeekdayRange = "monday"
	Tuesday   WeekdayRange = "tuesday"
	Wednesday WeekdayRange = "wednesday"
	Thursday  WeekdayRange = "thursday"
	Friday    WeekdayRange = "friday"
	Saturday  WeekdayRange = "saturday"
)

// Month constants.
const (
	January   MonthRange = "january"
	February  MonthRange = "february"
	March     MonthRange = "march"
	April     MonthRange = "april"
	May       MonthRange = "may"
	June      MonthRange = "june"
	July      MonthRange = "july"
	August    MonthRange = "august"
	September MonthRange = "september"
	October   MonthRange = "october"
	November  MonthRange = "november"
	December  MonthRange = "december"
)

// NewMuteTimeInterval creates a new MuteTimeInterval with the given name.
func NewMuteTimeInterval(name string) *MuteTimeInterval {
	return &MuteTimeInterval{Name: name}
}

// WithTimeIntervals sets the time intervals.
func (m *MuteTimeInterval) WithTimeIntervals(intervals ...*TimeInterval) *MuteTimeInterval {
	result := make([]TimeInterval, len(intervals))
	for i, ti := range intervals {
		result[i] = *ti
	}
	m.TimeIntervals = result
	return m
}

// NewTimeInterval creates a new TimeInterval.
func NewTimeInterval() *TimeInterval {
	return &TimeInterval{}
}

// WithTimes sets the time ranges for this interval.
func (t *TimeInterval) WithTimes(times ...*TimeRange) *TimeInterval {
	result := make([]TimeRange, len(times))
	for i, tr := range times {
		result[i] = *tr
	}
	t.Times = result
	return t
}

// WithWeekdays sets specific weekdays for this interval.
func (t *TimeInterval) WithWeekdays(days ...WeekdayRange) *TimeInterval {
	t.Weekdays = days
	return t
}

// WithWeekdayRange sets a range of weekdays (e.g., Monday to Friday).
func (t *TimeInterval) WithWeekdayRange(start, end WeekdayRange) *TimeInterval {
	t.Weekdays = []WeekdayRange{WeekdayRange(fmt.Sprintf("%s:%s", start, end))}
	return t
}

// WithDaysOfMonth sets specific days of the month.
func (t *TimeInterval) WithDaysOfMonth(days ...DayOfMonthRange) *TimeInterval {
	t.DaysOfMonth = days
	return t
}

// WithDayOfMonthRange sets a range of days (e.g., 1 to 7).
func (t *TimeInterval) WithDayOfMonthRange(start, end int) *TimeInterval {
	t.DaysOfMonth = []DayOfMonthRange{DayOfMonthRange(fmt.Sprintf("%d:%d", start, end))}
	return t
}

// WithMonths sets specific months for this interval.
func (t *TimeInterval) WithMonths(months ...MonthRange) *TimeInterval {
	t.Months = months
	return t
}

// WithMonthRange sets a range of months (e.g., January to March).
func (t *TimeInterval) WithMonthRange(start, end MonthRange) *TimeInterval {
	t.Months = []MonthRange{MonthRange(fmt.Sprintf("%s:%s", start, end))}
	return t
}

// WithYears sets specific years for this interval.
func (t *TimeInterval) WithYears(years ...YearRange) *TimeInterval {
	t.Years = years
	return t
}

// WithYearRange sets a range of years (e.g., 2024 to 2026).
func (t *TimeInterval) WithYearRange(start, end int) *TimeInterval {
	t.Years = []YearRange{YearRange(fmt.Sprintf("%d:%d", start, end))}
	return t
}

// NewTimeRange creates a new TimeRange.
func NewTimeRange(start, end string) *TimeRange {
	return &TimeRange{
		StartTime: start,
		EndTime:   end,
	}
}

// DayOfMonth creates a DayOfMonthRange for a specific day (1-31).
func DayOfMonth(day int) DayOfMonthRange {
	return DayOfMonthRange(fmt.Sprintf("%d", day))
}

// DayOfMonthEnd creates a DayOfMonthRange for a day from the end (-1 = last day).
func DayOfMonthEnd(offset int) DayOfMonthRange {
	return DayOfMonthRange(fmt.Sprintf("%d", offset))
}

// Year creates a YearRange for a specific year.
func Year(year int) YearRange {
	return YearRange(fmt.Sprintf("%d", year))
}

// WeekendsMuteInterval creates a mute interval for weekends (Saturday and Sunday).
func WeekendsMuteInterval() *MuteTimeInterval {
	return NewMuteTimeInterval("weekends").
		WithTimeIntervals(
			NewTimeInterval().WithWeekdays(Saturday, Sunday),
		)
}

// BusinessHoursMuteInterval creates a mute interval for business hours
// (Monday-Friday 09:00-17:00).
func BusinessHoursMuteInterval() *MuteTimeInterval {
	return NewMuteTimeInterval("business-hours").
		WithTimeIntervals(
			NewTimeInterval().
				WithWeekdayRange(Monday, Friday).
				WithTimes(NewTimeRange("09:00", "17:00")),
		)
}

// OutsideBusinessHoursMuteInterval creates a mute interval for outside business hours.
// This includes evenings (17:00-09:00 on weekdays) and weekends.
func OutsideBusinessHoursMuteInterval() *MuteTimeInterval {
	return NewMuteTimeInterval("outside-business-hours").
		WithTimeIntervals(
			// Weekday evenings
			NewTimeInterval().
				WithWeekdayRange(Monday, Friday).
				WithTimes(NewTimeRange("17:00", "23:59")),
			// Weekday mornings
			NewTimeInterval().
				WithWeekdayRange(Monday, Friday).
				WithTimes(NewTimeRange("00:00", "09:00")),
			// Weekends
			NewTimeInterval().
				WithWeekdays(Saturday, Sunday),
		)
}

// NightsMuteInterval creates a mute interval for nighttime (22:00-06:00).
func NightsMuteInterval() *MuteTimeInterval {
	return NewMuteTimeInterval("nights").
		WithTimeIntervals(
			NewTimeInterval().
				WithTimes(NewTimeRange("22:00", "23:59")),
			NewTimeInterval().
				WithTimes(NewTimeRange("00:00", "06:00")),
		)
}
