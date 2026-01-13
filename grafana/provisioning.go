package grafana

import (
	"gopkg.in/yaml.v3"
)

// DataSourceProvisioning represents Grafana data source provisioning configuration.
type DataSourceProvisioning struct {
	APIVersion int                          `yaml:"apiVersion"`
	DeleteDatasources []DeleteDatasourceEntry `yaml:"deleteDatasources,omitempty"`
	Datasources       []*DataSource           `yaml:"datasources"`
}

// DeleteDatasourceEntry identifies a data source to delete.
type DeleteDatasourceEntry struct {
	Name  string `yaml:"name"`
	OrgID int    `yaml:"orgId,omitempty"`
}

// NewDataSourceProvisioning creates a new data source provisioning config.
func NewDataSourceProvisioning(name string) *DataSourceProvisioning {
	return &DataSourceProvisioning{
		APIVersion:  1,
		Datasources: []*DataSource{},
	}
}

// AddDataSource adds a data source to the provisioning config.
func (p *DataSourceProvisioning) AddDataSource(ds *DataSource) *DataSourceProvisioning {
	p.Datasources = append(p.Datasources, ds)
	return p
}

// DeleteAllExisting marks all existing data sources for deletion.
func (p *DataSourceProvisioning) DeleteAllExisting() *DataSourceProvisioning {
	p.DeleteDatasources = []DeleteDatasourceEntry{{Name: "*"}}
	return p
}

// DeleteDataSource marks a specific data source for deletion.
func (p *DataSourceProvisioning) DeleteDataSource(name string, orgID int) *DataSourceProvisioning {
	p.DeleteDatasources = append(p.DeleteDatasources, DeleteDatasourceEntry{
		Name:  name,
		OrgID: orgID,
	})
	return p
}

// Serialize converts the provisioning config to YAML bytes.
func (p *DataSourceProvisioning) Serialize() ([]byte, error) {
	return yaml.Marshal(p)
}

// DashboardProvisioning represents Grafana dashboard provisioning configuration.
type DashboardProvisioning struct {
	APIVersion int                         `yaml:"apiVersion"`
	Providers  []DashboardProviderConfig `yaml:"providers"`
}

// DashboardProviderConfig represents a dashboard provider configuration.
type DashboardProviderConfig struct {
	Name                  string                   `yaml:"name"`
	OrgID                 int                      `yaml:"orgId,omitempty"`
	Folder                string                   `yaml:"folder,omitempty"`
	FolderUID             string                   `yaml:"folderUid,omitempty"`
	Type                  string                   `yaml:"type"`
	DisableDeletion       bool                     `yaml:"disableDeletion,omitempty"`
	Editable              bool                     `yaml:"editable,omitempty"`
	UpdateIntervalSeconds int                      `yaml:"updateIntervalSeconds,omitempty"`
	AllowUIUpdates        bool                     `yaml:"allowUiUpdates,omitempty"`
	Options               DashboardProviderOptions `yaml:"options"`
}

// DashboardProviderOptions contains dashboard provider options.
type DashboardProviderOptions struct {
	Path            string `yaml:"path"`
	FoldersFromFilesStructure bool `yaml:"foldersFromFilesStructure,omitempty"`
}

// NewDashboardProvisioning creates a new dashboard provisioning config.
func NewDashboardProvisioning(name string) *DashboardProvisioning {
	return &DashboardProvisioning{
		APIVersion: 1,
		Providers: []DashboardProviderConfig{
			{
				Name: name,
				Type: "file",
			},
		},
	}
}

// WithFolder sets the folder for dashboards.
func (p *DashboardProvisioning) WithFolder(folder string) *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].Folder = folder
	}
	return p
}

// WithFolderUID sets the folder UID for dashboards.
func (p *DashboardProvisioning) WithFolderUID(uid string) *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].FolderUID = uid
	}
	return p
}

// WithPath sets the path to dashboard files.
func (p *DashboardProvisioning) WithPath(path string) *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].Options.Path = path
	}
	return p
}

// Editable makes provisioned dashboards editable.
func (p *DashboardProvisioning) Editable() *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].Editable = true
	}
	return p
}

// ReadOnly makes provisioned dashboards read-only.
func (p *DashboardProvisioning) ReadOnly() *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].Editable = false
	}
	return p
}

// DisableDelete prevents deletion of provisioned dashboards.
func (p *DashboardProvisioning) DisableDelete() *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].DisableDeletion = true
	}
	return p
}

// AllowUIUpdates allows UI updates to provisioned dashboards.
func (p *DashboardProvisioning) AllowUIUpdates() *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].AllowUIUpdates = true
	}
	return p
}

// WithUpdateInterval sets the update interval in seconds.
func (p *DashboardProvisioning) WithUpdateInterval(seconds int) *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].UpdateIntervalSeconds = seconds
	}
	return p
}

// WithOrgID sets the organization ID.
func (p *DashboardProvisioning) WithOrgID(orgID int) *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].OrgID = orgID
	}
	return p
}

// FoldersFromFileStructure creates folders based on file structure.
func (p *DashboardProvisioning) FoldersFromFileStructure() *DashboardProvisioning {
	if len(p.Providers) > 0 {
		p.Providers[0].Options.FoldersFromFilesStructure = true
	}
	return p
}

// Serialize converts the provisioning config to YAML bytes.
func (p *DashboardProvisioning) Serialize() ([]byte, error) {
	return yaml.Marshal(p)
}
