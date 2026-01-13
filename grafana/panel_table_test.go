package grafana

import (
	"testing"
)

func TestTable(t *testing.T) {
	p := Table("Requests")
	if p == nil {
		t.Fatal("Table() returned nil")
	}
	if p.Title != "Requests" {
		t.Errorf("Title = %v, want Requests", p.Title)
	}
	if p.Type != "table" {
		t.Errorf("Type = %v, want table", p.Type)
	}
}

func TestTable_WithDescription(t *testing.T) {
	p := Table("Requests").WithDescription("HTTP requests table")
	if p.Description != "HTTP requests table" {
		t.Errorf("Description = %v", p.Description)
	}
}

func TestTable_WithSize(t *testing.T) {
	p := Table("Requests").WithSize(24, 10)
	if p.GridPos.W != 24 {
		t.Errorf("Width = %d, want 24", p.GridPos.W)
	}
	if p.GridPos.H != 10 {
		t.Errorf("Height = %d, want 10", p.GridPos.H)
	}
}

func TestTable_ShowHeader(t *testing.T) {
	p := Table("Requests").ShowHeader()
	if !p.Options.ShowHeader {
		t.Error("ShowHeader should be true")
	}
}

func TestTable_HideHeader(t *testing.T) {
	p := Table("Requests").ShowHeader().HideHeader()
	if p.Options.ShowHeader {
		t.Error("ShowHeader should be false")
	}
}

func TestTable_WithFooter(t *testing.T) {
	p := Table("Requests").WithFooter(true)
	if !p.Options.Footer.Show {
		t.Error("Footer.Show should be true")
	}
}

func TestTable_FooterCalcs(t *testing.T) {
	p := Table("Requests").WithFooter(true).WithFooterCalcs("sum", "mean")
	if len(p.Options.Footer.Reducer) != 2 {
		t.Errorf("len(Footer.Reducer) = %d, want 2", len(p.Options.Footer.Reducer))
	}
}

func TestTable_SortByColumn(t *testing.T) {
	p := Table("Requests").SortByColumn("Value", true)
	if p.Options.SortBy == nil || len(p.Options.SortBy) == 0 {
		t.Fatal("SortBy is nil or empty")
	}
	if p.Options.SortBy[0].DisplayName != "Value" {
		t.Errorf("SortBy[0].DisplayName = %v", p.Options.SortBy[0].DisplayName)
	}
	if !p.Options.SortBy[0].Desc {
		t.Error("SortBy[0].Desc should be true")
	}
}

func TestTable_WithColumnWidth(t *testing.T) {
	p := Table("Requests").WithColumnWidth("Time", 150)
	if len(p.FieldConfig.Overrides) == 0 {
		t.Fatal("Overrides is empty")
	}
}

func TestTable_HideColumn(t *testing.T) {
	p := Table("Requests").HideColumn("__name__")
	if len(p.FieldConfig.Overrides) == 0 {
		t.Fatal("Overrides is empty")
	}
}

func TestTable_EnableFiltering(t *testing.T) {
	p := Table("Requests").EnableFiltering()
	if !p.Options.EnableFiltering {
		t.Error("EnableFiltering should be true")
	}
}

func TestTable_EnablePagination(t *testing.T) {
	p := Table("Requests").EnablePagination()
	if !p.Options.EnablePagination {
		t.Error("EnablePagination should be true")
	}
}

func TestTable_FluentAPI(t *testing.T) {
	p := Table("HTTP Requests").
		WithDescription("Recent HTTP requests").
		WithDatasource("Prometheus").
		WithSize(24, 12).
		ShowHeader().
		EnableFiltering().
		SortByColumn("Time", true).
		HideColumn("__name__")

	if p.Title != "HTTP Requests" {
		t.Errorf("Title = %v", p.Title)
	}
	if !p.Options.ShowHeader {
		t.Error("ShowHeader should be true")
	}
	if !p.Options.EnableFiltering {
		t.Error("EnableFiltering should be true")
	}
}
