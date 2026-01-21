package importer

import (
	"strings"
	"testing"
)

func TestParseGrafanaDashboardFromBytes(t *testing.T) {
	input := `{
		"uid": "test-dashboard",
		"title": "Test Dashboard",
		"description": "A test dashboard",
		"tags": ["test", "example"],
		"editable": true,
		"schemaVersion": 39,
		"panels": [
			{
				"id": 1,
				"type": "row",
				"title": "Overview",
				"collapsed": false,
				"gridPos": {"x": 0, "y": 0, "w": 24, "h": 1}
			},
			{
				"id": 2,
				"type": "timeseries",
				"title": "Request Rate",
				"gridPos": {"x": 0, "y": 1, "w": 12, "h": 8},
				"targets": [
					{
						"refId": "A",
						"expr": "rate(http_requests_total[5m])",
						"legendFormat": "{{method}}"
					}
				]
			},
			{
				"id": 3,
				"type": "stat",
				"title": "Total Requests",
				"gridPos": {"x": 12, "y": 1, "w": 12, "h": 8}
			}
		],
		"templating": {
			"list": [
				{
					"name": "namespace",
					"type": "query",
					"query": "label_values(namespace)",
					"refresh": 1,
					"multi": true
				}
			]
		}
	}`

	dashboard, err := ParseGrafanaDashboardFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse dashboard: %v", err)
	}

	if dashboard.UID != "test-dashboard" {
		t.Errorf("expected UID 'test-dashboard', got %q", dashboard.UID)
	}

	if dashboard.Title != "Test Dashboard" {
		t.Errorf("expected title 'Test Dashboard', got %q", dashboard.Title)
	}

	if len(dashboard.Panels) != 3 {
		t.Errorf("expected 3 panels, got %d", len(dashboard.Panels))
	}

	if dashboard.Panels[0].Type != "row" {
		t.Errorf("expected first panel type 'row', got %q", dashboard.Panels[0].Type)
	}

	if dashboard.Panels[1].Type != "timeseries" {
		t.Errorf("expected second panel type 'timeseries', got %q", dashboard.Panels[1].Type)
	}

	if len(dashboard.Panels[1].Targets) != 1 {
		t.Errorf("expected 1 target, got %d", len(dashboard.Panels[1].Targets))
	}

	if dashboard.Templating == nil || len(dashboard.Templating.List) != 1 {
		t.Error("expected 1 variable")
	}
}

func TestValidateGrafanaDashboard(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		wantWarn []string
	}{
		{
			name: "valid dashboard",
			input: `{
				"uid": "test",
				"title": "Test",
				"panels": [
					{"type": "row", "title": "Row"},
					{"type": "timeseries", "title": "Panel"}
				]
			}`,
			wantWarn: nil,
		},
		{
			name:     "no title",
			input:    `{"uid": "test", "panels": []}`,
			wantWarn: []string{"dashboard has no title", "dashboard has no panels"},
		},
		{
			name:     "no UID",
			input:    `{"title": "Test", "panels": []}`,
			wantWarn: []string{"dashboard has no UID", "dashboard has no panels"},
		},
		{
			name: "no rows",
			input: `{
				"uid": "test",
				"title": "Test",
				"panels": [{"type": "timeseries", "title": "Panel"}]
			}`,
			wantWarn: []string{"dashboard has no rows"},
		},
		{
			name: "unsupported panel type",
			input: `{
				"uid": "test",
				"title": "Test",
				"panels": [
					{"type": "row", "title": "Row"},
					{"type": "unknown", "title": "Unknown Panel"}
				]
			}`,
			wantWarn: []string{"unsupported panel type"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dashboard, err := ParseGrafanaDashboardFromBytes([]byte(tt.input))
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}

			warnings := ValidateGrafanaDashboard(dashboard)

			if len(tt.wantWarn) == 0 && len(warnings) > 0 {
				t.Errorf("expected no warnings, got %v", warnings)
			}

			for _, want := range tt.wantWarn {
				found := false
				for _, w := range warnings {
					if strings.Contains(w, want) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("expected warning containing %q, got %v", want, warnings)
				}
			}
		})
	}
}

func TestGenerateGrafanaGoCode(t *testing.T) {
	input := `{
		"uid": "api-metrics",
		"title": "API Metrics",
		"description": "Dashboard for API metrics",
		"tags": ["api", "monitoring"],
		"editable": true,
		"schemaVersion": 39,
		"time": {"from": "now-1h", "to": "now"},
		"refresh": "30s",
		"panels": [
			{
				"id": 1,
				"type": "row",
				"title": "Overview",
				"gridPos": {"x": 0, "y": 0, "w": 24, "h": 1}
			},
			{
				"id": 2,
				"type": "timeseries",
				"title": "Request Rate",
				"gridPos": {"x": 0, "y": 1, "w": 12, "h": 8},
				"targets": [
					{
						"refId": "A",
						"expr": "rate(http_requests_total[5m])",
						"legendFormat": "{{method}}"
					}
				],
				"fieldConfig": {
					"defaults": {
						"unit": "reqps"
					}
				}
			},
			{
				"id": 3,
				"type": "stat",
				"title": "Total Requests",
				"gridPos": {"x": 12, "y": 1, "w": 12, "h": 8}
			}
		],
		"templating": {
			"list": [
				{
					"name": "namespace",
					"type": "query",
					"query": "label_values(namespace)",
					"label": "Namespace",
					"refresh": 1,
					"multi": true,
					"includeAll": true
				}
			]
		}
	}`

	dashboard, err := ParseGrafanaDashboardFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse dashboard: %v", err)
	}

	code, err := GenerateGrafanaGoCode(dashboard, "monitoring")
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	codeStr := string(code)

	// Check package declaration
	if !strings.Contains(codeStr, "package monitoring") {
		t.Error("expected 'package monitoring' in generated code")
	}

	// Check import
	if !strings.Contains(codeStr, "github.com/lex00/wetwire-observability-go/grafana") {
		t.Error("expected grafana import in generated code")
	}

	// Check dashboard
	if !strings.Contains(codeStr, "grafana.NewDashboard") {
		t.Error("expected NewDashboard in generated code")
	}

	// Check title
	if !strings.Contains(codeStr, `"API Metrics"`) {
		t.Error("expected dashboard title in generated code")
	}

	// Check panels
	if !strings.Contains(codeStr, "grafana.TimeSeries") {
		t.Error("expected TimeSeries panel in generated code")
	}

	if !strings.Contains(codeStr, "grafana.Stat") {
		t.Error("expected Stat panel in generated code")
	}

	// Check variable
	if !strings.Contains(codeStr, "grafana.QueryVar") {
		t.Error("expected QueryVar in generated code")
	}

	// Check method chains
	if !strings.Contains(codeStr, "WithDescription") {
		t.Error("expected WithDescription in generated code")
	}

	if !strings.Contains(codeStr, "WithTags") {
		t.Error("expected WithTags in generated code")
	}

	// Check targets
	if !strings.Contains(codeStr, "grafana.PromTarget") {
		t.Error("expected PromTarget in generated code")
	}

	if !strings.Contains(codeStr, `rate(http_requests_total[5m])`) {
		t.Error("expected PromQL expression in generated code")
	}
}

func TestGenerateGrafanaGoCodeCollapsedRow(t *testing.T) {
	input := `{
		"uid": "test",
		"title": "Test",
		"panels": [
			{
				"type": "row",
				"title": "Collapsed Section",
				"collapsed": true,
				"panels": [
					{"type": "timeseries", "title": "Hidden Panel"}
				]
			}
		]
	}`

	dashboard, err := ParseGrafanaDashboardFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse dashboard: %v", err)
	}

	code, err := GenerateGrafanaGoCode(dashboard, "monitoring")
	if err != nil {
		t.Fatalf("failed to generate code: %v", err)
	}

	codeStr := string(code)

	if !strings.Contains(codeStr, "Collapsed()") {
		t.Error("expected Collapsed() method call for collapsed row")
	}

	if !strings.Contains(codeStr, "Hidden Panel") {
		t.Error("expected nested panel in collapsed row")
	}
}

func TestConvertToWetwire(t *testing.T) {
	input := `{
		"uid": "test-dashboard",
		"title": "Test Dashboard",
		"panels": [
			{"type": "row", "title": "Row 1"},
			{"type": "timeseries", "title": "Panel 1"},
			{"type": "row", "title": "Row 2"},
			{"type": "stat", "title": "Panel 2"}
		]
	}`

	gd, err := ParseGrafanaDashboardFromBytes([]byte(input))
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	dashboard := ConvertToWetwire(gd)

	if dashboard.UID != "test-dashboard" {
		t.Errorf("expected UID 'test-dashboard', got %q", dashboard.UID)
	}

	if len(dashboard.Rows) != 2 {
		t.Errorf("expected 2 rows, got %d", len(dashboard.Rows))
	}

	if dashboard.Rows[0].Title != "Row 1" {
		t.Errorf("expected first row title 'Row 1', got %q", dashboard.Rows[0].Title)
	}

	if len(dashboard.Rows[0].Panels) != 1 {
		t.Errorf("expected 1 panel in first row, got %d", len(dashboard.Rows[0].Panels))
	}
}
