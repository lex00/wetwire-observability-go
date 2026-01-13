package grafana

import "testing"

func TestLogs(t *testing.T) {
	p := Logs("Application Logs")
	if p.Title != "Application Logs" {
		t.Errorf("Title = %q, want Application Logs", p.Title)
	}
	if p.Type != "logs" {
		t.Errorf("Type = %q, want logs", p.Type)
	}
}

func TestLogs_WithDescription(t *testing.T) {
	p := Logs("Test").WithDescription("A test logs panel")
	if p.Description != "A test logs panel" {
		t.Errorf("Description = %q, want A test logs panel", p.Description)
	}
}

func TestLogs_WithSize(t *testing.T) {
	p := Logs("Test").WithSize(24, 12)
	if p.GridPos.W != 24 || p.GridPos.H != 12 {
		t.Errorf("GridPos = %+v, want W=24 H=12", p.GridPos)
	}
}

func TestLogs_ShowTime(t *testing.T) {
	p := Logs("Test").ShowTime()
	if !p.Options.ShowTime {
		t.Error("ShowTime should be true")
	}
}

func TestLogs_HideTime(t *testing.T) {
	p := Logs("Test").ShowTime().HideTime()
	if p.Options.ShowTime {
		t.Error("ShowTime should be false")
	}
}

func TestLogs_WrapLines(t *testing.T) {
	p := Logs("Test").WrapLines()
	if !p.Options.WrapLogMessage {
		t.Error("WrapLogMessage should be true")
	}
}

func TestLogs_NoWrap(t *testing.T) {
	p := Logs("Test").WrapLines().NoWrap()
	if p.Options.WrapLogMessage {
		t.Error("WrapLogMessage should be false")
	}
}

func TestLogs_ShowLabels(t *testing.T) {
	p := Logs("Test").ShowLabels()
	if !p.Options.ShowLabels {
		t.Error("ShowLabels should be true")
	}
}

func TestLogs_ShowCommonLabels(t *testing.T) {
	p := Logs("Test").ShowCommonLabels()
	if !p.Options.ShowCommonLabels {
		t.Error("ShowCommonLabels should be true")
	}
}

func TestLogs_EnableLogDetails(t *testing.T) {
	p := Logs("Test").EnableLogDetails()
	if !p.Options.EnableLogDetails {
		t.Error("EnableLogDetails should be true")
	}
}

func TestLogs_SortOrder(t *testing.T) {
	p := Logs("Test").SortDescending()
	if p.Options.SortOrder != "Descending" {
		t.Errorf("SortOrder = %q, want Descending", p.Options.SortOrder)
	}

	p = Logs("Test").SortAscending()
	if p.Options.SortOrder != "Ascending" {
		t.Errorf("SortOrder = %q, want Ascending", p.Options.SortOrder)
	}
}

func TestLogs_PrettifyJSON(t *testing.T) {
	p := Logs("Test").PrettifyJSON()
	if !p.Options.PrettifyLogMessage {
		t.Error("PrettifyLogMessage should be true")
	}
}

func TestLogs_FluentAPI(t *testing.T) {
	p := Logs("Test").
		WithDescription("desc").
		WithSize(24, 16).
		ShowTime().
		WrapLines().
		ShowLabels().
		SortDescending().
		EnableLogDetails()

	if p.Title != "Test" {
		t.Error("Fluent API should preserve title")
	}
	if !p.Options.ShowTime {
		t.Error("Fluent API should set ShowTime")
	}
	if !p.Options.WrapLogMessage {
		t.Error("Fluent API should set WrapLogMessage")
	}
}
