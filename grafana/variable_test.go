package grafana

import (
	"testing"
)

func TestQueryVar(t *testing.T) {
	v := QueryVar("namespace", "label_values(kube_pod_info, namespace)")
	if v == nil {
		t.Fatal("QueryVar() returned nil")
	}
	if v.Name != "namespace" {
		t.Errorf("Name = %v, want namespace", v.Name)
	}
	if v.Type != "query" {
		t.Errorf("Type = %v, want query", v.Type)
	}
	if v.Query != "label_values(kube_pod_info, namespace)" {
		t.Errorf("Query = %v", v.Query)
	}
}

func TestQueryVar_WithDatasource(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").
		WithDatasource("Prometheus")
	if v.Datasource != "Prometheus" {
		t.Errorf("Datasource = %v, want Prometheus", v.Datasource)
	}
}

func TestQueryVar_WithRegex(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").
		WithRegex("^(?!kube-).*")
	if v.Regex != "^(?!kube-).*" {
		t.Errorf("Regex = %v", v.Regex)
	}
}

func TestQueryVar_WithSort(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").
		WithSort(SortAlphabetical)
	if v.Sort != SortAlphabetical {
		t.Errorf("Sort = %v, want %v", v.Sort, SortAlphabetical)
	}
}

func TestQueryVar_WithRefresh(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").
		WithRefresh(RefreshOnTimeRangeChange)
	if v.Refresh != RefreshOnTimeRangeChange {
		t.Errorf("Refresh = %v, want %v", v.Refresh, RefreshOnTimeRangeChange)
	}
}

func TestQueryVar_WithLabel(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").
		WithLabel("Namespace")
	if v.Label != "Namespace" {
		t.Errorf("Label = %v, want Namespace", v.Label)
	}
}

func TestQueryVar_WithDescription(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").
		WithDescription("Filter by namespace")
	if v.Description != "Filter by namespace" {
		t.Errorf("Description = %v", v.Description)
	}
}

func TestQueryVar_MultiSelect(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").MultiSelect()
	if !v.Multi {
		t.Error("Multi should be true")
	}
}

func TestQueryVar_IncludeAll(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").IncludeAll()
	if !v.IncludeAllOption {
		t.Error("IncludeAllOption should be true")
	}
}

func TestQueryVar_FluentAPI(t *testing.T) {
	v := QueryVar("namespace", "label_values(kube_pod_info, namespace)").
		WithDatasource("Prometheus").
		WithLabel("Namespace").
		WithRegex("^(?!kube-).*").
		WithSort(SortAlphabetical).
		WithRefresh(RefreshOnDashboardLoad).
		MultiSelect().
		IncludeAll()

	if v.Name != "namespace" {
		t.Errorf("Name = %v", v.Name)
	}
	if v.Datasource != "Prometheus" {
		t.Errorf("Datasource = %v", v.Datasource)
	}
	if !v.Multi {
		t.Error("Multi should be true")
	}
	if !v.IncludeAllOption {
		t.Error("IncludeAllOption should be true")
	}
}

func TestCustomVar(t *testing.T) {
	v := CustomVar("env", "dev", "staging", "prod")
	if v == nil {
		t.Fatal("CustomVar() returned nil")
	}
	if v.Name != "env" {
		t.Errorf("Name = %v, want env", v.Name)
	}
	if v.Type != "custom" {
		t.Errorf("Type = %v, want custom", v.Type)
	}
	if v.Query != "dev,staging,prod" {
		t.Errorf("Query = %v, want dev,staging,prod", v.Query)
	}
}

func TestCustomVar_WithDefault(t *testing.T) {
	v := CustomVar("env", "dev", "staging", "prod").WithDefault("prod")
	if v.Current != "prod" {
		t.Errorf("Current = %v, want prod", v.Current)
	}
}

func TestIntervalVar(t *testing.T) {
	v := IntervalVar("interval", "1m", "5m", "15m", "1h")
	if v == nil {
		t.Fatal("IntervalVar() returned nil")
	}
	if v.Name != "interval" {
		t.Errorf("Name = %v, want interval", v.Name)
	}
	if v.Type != "interval" {
		t.Errorf("Type = %v, want interval", v.Type)
	}
}

func TestIntervalVar_AutoOption(t *testing.T) {
	v := IntervalVar("interval", "1m", "5m").AutoOption(10, "10s")
	if !v.Auto {
		t.Error("Auto should be true")
	}
	if v.AutoCount != 10 {
		t.Errorf("AutoCount = %d, want 10", v.AutoCount)
	}
	if v.AutoMin != "10s" {
		t.Errorf("AutoMin = %v, want 10s", v.AutoMin)
	}
}

func TestDatasourceVar(t *testing.T) {
	v := DatasourceVar("datasource", "prometheus")
	if v == nil {
		t.Fatal("DatasourceVar() returned nil")
	}
	if v.Name != "datasource" {
		t.Errorf("Name = %v, want datasource", v.Name)
	}
	if v.Type != "datasource" {
		t.Errorf("Type = %v, want datasource", v.Type)
	}
	if v.Query != "prometheus" {
		t.Errorf("Query = %v, want prometheus", v.Query)
	}
}

func TestDatasourceVar_WithRegex(t *testing.T) {
	v := DatasourceVar("datasource", "prometheus").WithRegex("^prod-.*")
	if v.Regex != "^prod-.*" {
		t.Errorf("Regex = %v", v.Regex)
	}
}

func TestTextboxVar(t *testing.T) {
	v := TextboxVar("filter", "default-value")
	if v == nil {
		t.Fatal("TextboxVar() returned nil")
	}
	if v.Name != "filter" {
		t.Errorf("Name = %v, want filter", v.Name)
	}
	if v.Type != "textbox" {
		t.Errorf("Type = %v, want textbox", v.Type)
	}
	if v.Query != "default-value" {
		t.Errorf("Query = %v, want default-value", v.Query)
	}
}

func TestConstantVar(t *testing.T) {
	v := ConstantVar("cluster", "production")
	if v == nil {
		t.Fatal("ConstantVar() returned nil")
	}
	if v.Name != "cluster" {
		t.Errorf("Name = %v, want cluster", v.Name)
	}
	if v.Type != "constant" {
		t.Errorf("Type = %v, want constant", v.Type)
	}
	if v.Query != "production" {
		t.Errorf("Query = %v, want production", v.Query)
	}
}

func TestVariable_Hide(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").Hide()
	if v.HideOption != HideVariable {
		t.Errorf("HideOption = %v, want %v", v.HideOption, HideVariable)
	}
}

func TestVariable_HideLabel(t *testing.T) {
	v := QueryVar("namespace", "label_values(namespace)").HideLabel()
	if v.HideOption != HideLabelOnly {
		t.Errorf("HideOption = %v, want %v", v.HideOption, HideLabelOnly)
	}
}
