package grafana

// Data source type constants.
const (
	DataSourceTypePrometheus   = "prometheus"
	DataSourceTypeLoki         = "loki"
	DataSourceTypeJaeger       = "jaeger"
	DataSourceTypeTempo        = "tempo"
	DataSourceTypeElasticsearch = "elasticsearch"
	DataSourceTypeInfluxDB     = "influxdb"
	DataSourceTypeGraphite     = "graphite"
	DataSourceTypeMySQL        = "mysql"
	DataSourceTypePostgreSQL   = "postgres"
	DataSourceTypeCloudWatch   = "cloudwatch"
)

// DataSource represents a Grafana data source configuration.
type DataSource struct {
	// Name is the data source name (must be unique).
	Name string `json:"name" yaml:"name"`

	// UID is the unique identifier for referencing in dashboards.
	UID string `json:"uid,omitempty" yaml:"uid,omitempty"`

	// Type is the data source type (prometheus, loki, etc.).
	Type string `json:"type" yaml:"type"`

	// URL is the data source URL.
	URL string `json:"url" yaml:"url"`

	// Access mode (proxy or direct).
	Access string `json:"access,omitempty" yaml:"access,omitempty"`

	// IsDefault marks this as the default data source.
	IsDefault bool `json:"isDefault,omitempty" yaml:"isDefault,omitempty"`

	// IsEditable controls whether the data source can be modified in UI.
	IsEditable bool `json:"editable,omitempty" yaml:"editable,omitempty"`

	// BasicAuth enables basic authentication.
	BasicAuth bool `json:"basicAuth,omitempty" yaml:"basicAuth,omitempty"`

	// BasicAuthUser is the basic auth username.
	BasicAuthUser string `json:"basicAuthUser,omitempty" yaml:"basicAuthUser,omitempty"`

	// JSONData contains additional configuration.
	JSONData map[string]any `json:"jsonData,omitempty" yaml:"jsonData,omitempty"`

	// SecureJSONData contains sensitive configuration.
	SecureJSONData map[string]string `json:"secureJsonData,omitempty" yaml:"secureJsonData,omitempty"`

	// OrgID is the organization ID.
	OrgID int `json:"orgId,omitempty" yaml:"orgId,omitempty"`

	// Version is the data source version.
	Version int `json:"version,omitempty" yaml:"version,omitempty"`
}

// DataSourceRef is a reference to a data source.
type DataSourceRef struct {
	Type string `json:"type"`
	UID  string `json:"uid"`
}

// NewDataSource creates a new data source.
func NewDataSource(name, uid, dsType string) *DataSource {
	return &DataSource{
		Name:   name,
		UID:    uid,
		Type:   dsType,
		Access: "proxy",
	}
}

// PrometheusDataSource creates a Prometheus data source.
func PrometheusDataSource(name, url string) *DataSource {
	return &DataSource{
		Name:   name,
		UID:    name,
		Type:   DataSourceTypePrometheus,
		URL:    url,
		Access: "proxy",
	}
}

// LokiDataSource creates a Loki data source.
func LokiDataSource(name, url string) *DataSource {
	return &DataSource{
		Name:   name,
		UID:    name,
		Type:   DataSourceTypeLoki,
		URL:    url,
		Access: "proxy",
	}
}

// JaegerDataSource creates a Jaeger data source.
func JaegerDataSource(name, url string) *DataSource {
	return &DataSource{
		Name:   name,
		UID:    name,
		Type:   DataSourceTypeJaeger,
		URL:    url,
		Access: "proxy",
	}
}

// TempoDataSource creates a Tempo data source.
func TempoDataSource(name, url string) *DataSource {
	return &DataSource{
		Name:   name,
		UID:    name,
		Type:   DataSourceTypeTempo,
		URL:    url,
		Access: "proxy",
	}
}

// WithURL sets the data source URL.
func (ds *DataSource) WithURL(url string) *DataSource {
	ds.URL = url
	return ds
}

// WithUID sets the data source UID.
func (ds *DataSource) WithUID(uid string) *DataSource {
	ds.UID = uid
	return ds
}

// AsDefault marks this data source as default.
func (ds *DataSource) AsDefault() *DataSource {
	ds.IsDefault = true
	return ds
}

// Editable makes the data source editable in the UI.
func (ds *DataSource) Editable() *DataSource {
	ds.IsEditable = true
	return ds
}

// ReadOnly makes the data source read-only in the UI.
func (ds *DataSource) ReadOnly() *DataSource {
	ds.IsEditable = false
	return ds
}

// WithBasicAuth enables basic authentication.
func (ds *DataSource) WithBasicAuth(user, password string) *DataSource {
	ds.BasicAuth = true
	ds.BasicAuthUser = user
	if ds.SecureJSONData == nil {
		ds.SecureJSONData = make(map[string]string)
	}
	ds.SecureJSONData["basicAuthPassword"] = password
	return ds
}

// WithJSONData sets additional JSON configuration.
func (ds *DataSource) WithJSONData(data map[string]any) *DataSource {
	ds.JSONData = data
	return ds
}

// AddJSONData adds a JSON data field.
func (ds *DataSource) AddJSONData(key string, value any) *DataSource {
	if ds.JSONData == nil {
		ds.JSONData = make(map[string]any)
	}
	ds.JSONData[key] = value
	return ds
}

// WithSecureJSONData sets sensitive JSON configuration.
func (ds *DataSource) WithSecureJSONData(data map[string]string) *DataSource {
	ds.SecureJSONData = data
	return ds
}

// AddSecureJSONData adds a secure JSON data field.
func (ds *DataSource) AddSecureJSONData(key, value string) *DataSource {
	if ds.SecureJSONData == nil {
		ds.SecureJSONData = make(map[string]string)
	}
	ds.SecureJSONData[key] = value
	return ds
}

// WithOrgID sets the organization ID.
func (ds *DataSource) WithOrgID(orgID int) *DataSource {
	ds.OrgID = orgID
	return ds
}

// Ref returns a reference to this data source.
func (ds *DataSource) Ref() DataSourceRef {
	return DataSourceRef{
		Type: ds.Type,
		UID:  ds.Name,
	}
}

// DirectAccess sets direct access mode (browser connects directly).
func (ds *DataSource) DirectAccess() *DataSource {
	ds.Access = "direct"
	return ds
}

// ProxyAccess sets proxy access mode (default).
func (ds *DataSource) ProxyAccess() *DataSource {
	ds.Access = "proxy"
	return ds
}
