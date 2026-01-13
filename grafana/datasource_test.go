package grafana

import "testing"

func TestNewDataSource(t *testing.T) {
	ds := NewDataSource("prometheus", "Prometheus", DataSourceTypePrometheus)
	if ds.Name != "prometheus" {
		t.Errorf("Name = %q, want prometheus", ds.Name)
	}
	if ds.UID != "Prometheus" {
		t.Errorf("UID = %q, want Prometheus", ds.UID)
	}
	if ds.Type != DataSourceTypePrometheus {
		t.Errorf("Type = %q, want prometheus", ds.Type)
	}
}

func TestPrometheusDataSource(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090")
	if ds.Type != DataSourceTypePrometheus {
		t.Errorf("Type = %q, want prometheus", ds.Type)
	}
	if ds.URL != "http://prometheus:9090" {
		t.Errorf("URL = %q, want http://prometheus:9090", ds.URL)
	}
}

func TestLokiDataSource(t *testing.T) {
	ds := LokiDataSource("loki", "http://loki:3100")
	if ds.Type != DataSourceTypeLoki {
		t.Errorf("Type = %q, want loki", ds.Type)
	}
	if ds.URL != "http://loki:3100" {
		t.Errorf("URL = %q, want http://loki:3100", ds.URL)
	}
}

func TestJaegerDataSource(t *testing.T) {
	ds := JaegerDataSource("jaeger", "http://jaeger:16686")
	if ds.Type != DataSourceTypeJaeger {
		t.Errorf("Type = %q, want jaeger", ds.Type)
	}
}

func TestTempoDataSource(t *testing.T) {
	ds := TempoDataSource("tempo", "http://tempo:3200")
	if ds.Type != DataSourceTypeTempo {
		t.Errorf("Type = %q, want tempo", ds.Type)
	}
}

func TestDataSource_WithURL(t *testing.T) {
	ds := NewDataSource("ds", "DS", "custom").WithURL("http://custom:8080")
	if ds.URL != "http://custom:8080" {
		t.Errorf("URL = %q, want http://custom:8080", ds.URL)
	}
}

func TestDataSource_AsDefault(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").AsDefault()
	if !ds.IsDefault {
		t.Error("IsDefault should be true")
	}
}

func TestDataSource_Editable(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").Editable()
	if !ds.IsEditable {
		t.Error("IsEditable should be true")
	}
}

func TestDataSource_ReadOnly(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").Editable().ReadOnly()
	if ds.IsEditable {
		t.Error("IsEditable should be false")
	}
}

func TestDataSource_WithBasicAuth(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").
		WithBasicAuth("admin", "password")
	if !ds.BasicAuth {
		t.Error("BasicAuth should be true")
	}
	if ds.BasicAuthUser != "admin" {
		t.Errorf("BasicAuthUser = %q, want admin", ds.BasicAuthUser)
	}
}

func TestDataSource_WithJSONData(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").
		WithJSONData(map[string]any{
			"httpMethod": "POST",
		})
	if ds.JSONData == nil {
		t.Fatal("JSONData is nil")
	}
	if ds.JSONData["httpMethod"] != "POST" {
		t.Errorf("JSONData[httpMethod] = %v, want POST", ds.JSONData["httpMethod"])
	}
}

func TestDataSource_AddJSONData(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").
		AddJSONData("httpMethod", "POST").
		AddJSONData("timeout", "60")
	if ds.JSONData["httpMethod"] != "POST" {
		t.Errorf("JSONData[httpMethod] = %v, want POST", ds.JSONData["httpMethod"])
	}
	if ds.JSONData["timeout"] != "60" {
		t.Errorf("JSONData[timeout] = %v, want 60", ds.JSONData["timeout"])
	}
}

func TestDataSource_WithSecureJSONData(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").
		WithSecureJSONData(map[string]string{
			"basicAuthPassword": "secret",
		})
	if ds.SecureJSONData == nil {
		t.Fatal("SecureJSONData is nil")
	}
	if ds.SecureJSONData["basicAuthPassword"] != "secret" {
		t.Errorf("SecureJSONData[basicAuthPassword] = %v", ds.SecureJSONData["basicAuthPassword"])
	}
}

func TestDataSource_Ref(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090")
	ref := ds.Ref()
	if ref.Type != DataSourceTypePrometheus {
		t.Errorf("Type = %q, want prometheus", ref.Type)
	}
	if ref.UID != "prometheus" {
		t.Errorf("UID = %q, want prometheus", ref.UID)
	}
}

func TestDataSource_FluentAPI(t *testing.T) {
	ds := PrometheusDataSource("prometheus", "http://prometheus:9090").
		AsDefault().
		Editable().
		AddJSONData("httpMethod", "POST").
		AddJSONData("timeInterval", "60s")

	if !ds.IsDefault {
		t.Error("Fluent API should set IsDefault")
	}
	if !ds.IsEditable {
		t.Error("Fluent API should set IsEditable")
	}
	if ds.JSONData["httpMethod"] != "POST" {
		t.Error("Fluent API should set JSONData")
	}
}
