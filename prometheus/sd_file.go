package prometheus

// FileSD configures file-based service discovery.
// It reads target groups from JSON or YAML files and watches for changes.
//
// Example usage:
//
//	var FileDiscovery = prometheus.NewFileSD().
//	    WithFiles("/etc/prometheus/targets/*.json").
//	    WithRefreshInterval(prometheus.Duration(5 * time.Minute))
type FileSD struct {
	// Files is the list of file patterns to discover targets from.
	// Files may be JSON or YAML format.
	// Supports glob patterns (e.g., "/etc/prometheus/targets/*.json").
	Files []string `yaml:"files"`

	// RefreshInterval is the time after which the files are re-read.
	// Defaults to 5m.
	RefreshInterval Duration `yaml:"refresh_interval,omitempty"`
}

// NewFileSD creates a new file-based service discovery configuration.
func NewFileSD() *FileSD {
	return &FileSD{}
}

// WithFiles sets the file patterns to discover targets from.
func (f *FileSD) WithFiles(files ...string) *FileSD {
	f.Files = files
	return f
}

// WithRefreshInterval sets the file re-read interval.
func (f *FileSD) WithRefreshInterval(d Duration) *FileSD {
	f.RefreshInterval = d
	return f
}
