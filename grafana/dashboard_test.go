package grafana

import (
	"testing"
)

func TestNewDashboard(t *testing.T) {
	d := NewDashboard("test-uid", "Test Dashboard")
	if d == nil {
		t.Fatal("NewDashboard() returned nil")
	}
	if d.UID != "test-uid" {
		t.Errorf("UID = %v, want test-uid", d.UID)
	}
	if d.Title != "Test Dashboard" {
		t.Errorf("Title = %v, want Test Dashboard", d.Title)
	}
}

func TestDashboard_WithTags(t *testing.T) {
	d := NewDashboard("test", "Test").WithTags("monitoring", "sre")
	if len(d.Tags) != 2 {
		t.Errorf("len(Tags) = %d, want 2", len(d.Tags))
	}
	if d.Tags[0] != "monitoring" || d.Tags[1] != "sre" {
		t.Errorf("Tags = %v, want [monitoring, sre]", d.Tags)
	}
}

func TestDashboard_WithTime(t *testing.T) {
	d := NewDashboard("test", "Test").WithTime("now-1h", "now")
	if d.Time == nil {
		t.Fatal("Time is nil")
	}
	if d.Time.From != "now-1h" {
		t.Errorf("Time.From = %v, want now-1h", d.Time.From)
	}
	if d.Time.To != "now" {
		t.Errorf("Time.To = %v, want now", d.Time.To)
	}
}

func TestDashboard_WithRefresh(t *testing.T) {
	d := NewDashboard("test", "Test").WithRefresh("30s")
	if d.Refresh != "30s" {
		t.Errorf("Refresh = %v, want 30s", d.Refresh)
	}
}

func TestDashboard_WithDescription(t *testing.T) {
	d := NewDashboard("test", "Test").WithDescription("A test dashboard")
	if d.Description != "A test dashboard" {
		t.Errorf("Description = %v, want 'A test dashboard'", d.Description)
	}
}

func TestDashboard_WithRows(t *testing.T) {
	row1 := NewRow("Row 1")
	row2 := NewRow("Row 2")
	d := NewDashboard("test", "Test").WithRows(row1, row2)
	if len(d.Rows) != 2 {
		t.Errorf("len(Rows) = %d, want 2", len(d.Rows))
	}
}

func TestDashboard_AddRow(t *testing.T) {
	d := NewDashboard("test", "Test")
	d.AddRow(NewRow("Row 1"))
	d.AddRow(NewRow("Row 2"))
	if len(d.Rows) != 2 {
		t.Errorf("len(Rows) = %d, want 2", len(d.Rows))
	}
}

func TestDashboard_Editable(t *testing.T) {
	d := NewDashboard("test", "Test").Editable()
	if !d.IsEditable {
		t.Error("IsEditable should be true")
	}
}

func TestDashboard_ReadOnly(t *testing.T) {
	d := NewDashboard("test", "Test").ReadOnly()
	if d.IsEditable {
		t.Error("IsEditable should be false")
	}
}

func TestDashboard_WithTimezone(t *testing.T) {
	d := NewDashboard("test", "Test").WithTimezone("UTC")
	if d.Timezone != "UTC" {
		t.Errorf("Timezone = %v, want UTC", d.Timezone)
	}
}

func TestDashboard_FluentAPI(t *testing.T) {
	d := NewDashboard("api-overview", "API Overview").
		WithDescription("Overview of API metrics").
		WithTags("api", "overview").
		WithTime("now-6h", "now").
		WithRefresh("1m").
		WithTimezone("browser").
		Editable().
		WithRows(
			NewRow("Overview"),
			NewRow("Details"),
		)

	if d.UID != "api-overview" {
		t.Errorf("UID = %v", d.UID)
	}
	if d.Title != "API Overview" {
		t.Errorf("Title = %v", d.Title)
	}
	if len(d.Tags) != 2 {
		t.Errorf("len(Tags) = %d", len(d.Tags))
	}
	if d.Refresh != "1m" {
		t.Errorf("Refresh = %v", d.Refresh)
	}
	if len(d.Rows) != 2 {
		t.Errorf("len(Rows) = %d", len(d.Rows))
	}
}
