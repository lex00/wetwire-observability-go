package alertmanager

import (
	"strings"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestNewMuteTimeInterval(t *testing.T) {
	mti := NewMuteTimeInterval("maintenance")
	if mti == nil {
		t.Error("NewMuteTimeInterval() returned nil")
	}
	if mti.Name != "maintenance" {
		t.Errorf("Name = %v, want maintenance", mti.Name)
	}
}

func TestMuteTimeInterval_WithTimeIntervals(t *testing.T) {
	mti := NewMuteTimeInterval("weekends").
		WithTimeIntervals(
			NewTimeInterval().
				WithWeekdays(Saturday, Sunday),
		)

	if len(mti.TimeIntervals) != 1 {
		t.Errorf("len(TimeIntervals) = %d, want 1", len(mti.TimeIntervals))
	}
}

func TestNewTimeInterval(t *testing.T) {
	ti := NewTimeInterval()
	if ti == nil {
		t.Error("NewTimeInterval() returned nil")
	}
}

func TestTimeInterval_WithTimes(t *testing.T) {
	ti := NewTimeInterval().
		WithTimes(
			NewTimeRange("09:00", "17:00"),
		)

	if len(ti.Times) != 1 {
		t.Errorf("len(Times) = %d, want 1", len(ti.Times))
	}
	if ti.Times[0].StartTime != "09:00" {
		t.Errorf("StartTime = %v", ti.Times[0].StartTime)
	}
}

func TestTimeInterval_WithWeekdays(t *testing.T) {
	ti := NewTimeInterval().
		WithWeekdays(Monday, Tuesday, Wednesday)

	if len(ti.Weekdays) != 3 {
		t.Errorf("len(Weekdays) = %d, want 3", len(ti.Weekdays))
	}
}

func TestTimeInterval_WithWeekdayRange(t *testing.T) {
	ti := NewTimeInterval().
		WithWeekdayRange(Monday, Friday)

	if len(ti.Weekdays) != 1 {
		t.Errorf("len(Weekdays) = %d, want 1", len(ti.Weekdays))
	}
	if string(ti.Weekdays[0]) != "monday:friday" {
		t.Errorf("Weekdays[0] = %v, want monday:friday", ti.Weekdays[0])
	}
}

func TestTimeInterval_WithDaysOfMonth(t *testing.T) {
	ti := NewTimeInterval().
		WithDaysOfMonth(DayOfMonth(1), DayOfMonth(15), DayOfMonthEnd(-1))

	if len(ti.DaysOfMonth) != 3 {
		t.Errorf("len(DaysOfMonth) = %d, want 3", len(ti.DaysOfMonth))
	}
}

func TestTimeInterval_WithDayOfMonthRange(t *testing.T) {
	ti := NewTimeInterval().
		WithDayOfMonthRange(1, 7)

	if len(ti.DaysOfMonth) != 1 {
		t.Errorf("len(DaysOfMonth) = %d, want 1", len(ti.DaysOfMonth))
	}
	if string(ti.DaysOfMonth[0]) != "1:7" {
		t.Errorf("DaysOfMonth[0] = %v, want 1:7", ti.DaysOfMonth[0])
	}
}

func TestTimeInterval_WithMonths(t *testing.T) {
	ti := NewTimeInterval().
		WithMonths(January, February, December)

	if len(ti.Months) != 3 {
		t.Errorf("len(Months) = %d, want 3", len(ti.Months))
	}
}

func TestTimeInterval_WithMonthRange(t *testing.T) {
	ti := NewTimeInterval().
		WithMonthRange(January, March)

	if len(ti.Months) != 1 {
		t.Errorf("len(Months) = %d, want 1", len(ti.Months))
	}
	if string(ti.Months[0]) != "january:march" {
		t.Errorf("Months[0] = %v, want january:march", ti.Months[0])
	}
}

func TestTimeInterval_WithYears(t *testing.T) {
	ti := NewTimeInterval().
		WithYears(Year(2024), Year(2025))

	if len(ti.Years) != 2 {
		t.Errorf("len(Years) = %d, want 2", len(ti.Years))
	}
}

func TestTimeInterval_WithYearRange(t *testing.T) {
	ti := NewTimeInterval().
		WithYearRange(2024, 2026)

	if len(ti.Years) != 1 {
		t.Errorf("len(Years) = %d, want 1", len(ti.Years))
	}
	if string(ti.Years[0]) != "2024:2026" {
		t.Errorf("Years[0] = %v, want 2024:2026", ti.Years[0])
	}
}

func TestMuteTimeInterval_Serialize(t *testing.T) {
	mti := NewMuteTimeInterval("maintenance").
		WithTimeIntervals(
			NewTimeInterval().
				WithWeekdays(Saturday, Sunday).
				WithTimes(NewTimeRange("00:00", "23:59")),
		)

	data, err := yaml.Marshal(mti)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	expectations := []string{
		"name: maintenance",
		"time_intervals:",
		"weekdays:",
		"- saturday",
		"- sunday",
		"times:",
		"start_time: \"00:00\"",
		"end_time: \"23:59\"",
	}

	for _, exp := range expectations {
		if !strings.Contains(yamlStr, exp) {
			t.Errorf("yaml.Marshal() missing %q\nGot:\n%s", exp, yamlStr)
		}
	}
}

func TestMuteTimeInterval_SerializeComplex(t *testing.T) {
	mti := NewMuteTimeInterval("business-hours").
		WithTimeIntervals(
			NewTimeInterval().
				WithWeekdayRange(Monday, Friday).
				WithTimes(NewTimeRange("09:00", "17:00")),
		)

	data, err := yaml.Marshal(mti)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "monday:friday") {
		t.Errorf("yaml.Marshal() missing monday:friday\nGot:\n%s", yamlStr)
	}
}

func TestMuteTimeInterval_Unmarshal(t *testing.T) {
	input := `
name: maintenance
time_intervals:
  - weekdays:
      - saturday
      - sunday
    times:
      - start_time: "00:00"
        end_time: "23:59"
`
	var mti MuteTimeInterval
	if err := yaml.Unmarshal([]byte(input), &mti); err != nil {
		t.Fatalf("yaml.Unmarshal() error = %v", err)
	}

	if mti.Name != "maintenance" {
		t.Errorf("Name = %v", mti.Name)
	}
	if len(mti.TimeIntervals) != 1 {
		t.Errorf("len(TimeIntervals) = %d, want 1", len(mti.TimeIntervals))
	}
}

func TestWeekdayConstants(t *testing.T) {
	tests := []struct {
		constant WeekdayRange
		want     string
	}{
		{Sunday, "sunday"},
		{Monday, "monday"},
		{Tuesday, "tuesday"},
		{Wednesday, "wednesday"},
		{Thursday, "thursday"},
		{Friday, "friday"},
		{Saturday, "saturday"},
	}

	for _, tt := range tests {
		if string(tt.constant) != tt.want {
			t.Errorf("constant = %v, want %v", tt.constant, tt.want)
		}
	}
}

func TestMonthConstants(t *testing.T) {
	tests := []struct {
		constant MonthRange
		want     string
	}{
		{January, "january"},
		{February, "february"},
		{March, "march"},
		{April, "april"},
		{May, "may"},
		{June, "june"},
		{July, "july"},
		{August, "august"},
		{September, "september"},
		{October, "october"},
		{November, "november"},
		{December, "december"},
	}

	for _, tt := range tests {
		if string(tt.constant) != tt.want {
			t.Errorf("constant = %v, want %v", tt.constant, tt.want)
		}
	}
}

func TestWeekendsMuteInterval(t *testing.T) {
	mti := WeekendsMuteInterval()
	if mti.Name != "weekends" {
		t.Errorf("Name = %v, want weekends", mti.Name)
	}
	if len(mti.TimeIntervals) != 1 {
		t.Errorf("len(TimeIntervals) = %d, want 1", len(mti.TimeIntervals))
	}
}

func TestBusinessHoursMuteInterval(t *testing.T) {
	mti := BusinessHoursMuteInterval()
	if mti.Name != "business-hours" {
		t.Errorf("Name = %v, want business-hours", mti.Name)
	}
}

func TestOutsideBusinessHoursMuteInterval(t *testing.T) {
	mti := OutsideBusinessHoursMuteInterval()
	if mti.Name != "outside-business-hours" {
		t.Errorf("Name = %v, want outside-business-hours", mti.Name)
	}
}

func TestMuteTimeIntervalInConfig(t *testing.T) {
	config := NewAlertmanagerConfig().
		WithMuteTimeIntervals(
			WeekendsMuteInterval(),
			BusinessHoursMuteInterval(),
		)

	data, err := yaml.Marshal(config)
	if err != nil {
		t.Fatalf("yaml.Marshal() error = %v", err)
	}

	yamlStr := string(data)
	if !strings.Contains(yamlStr, "mute_time_intervals:") {
		t.Errorf("yaml.Marshal() missing mute_time_intervals:\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "name: weekends") {
		t.Errorf("yaml.Marshal() missing weekends\nGot:\n%s", yamlStr)
	}
	if !strings.Contains(yamlStr, "name: business-hours") {
		t.Errorf("yaml.Marshal() missing business-hours\nGot:\n%s", yamlStr)
	}
}
