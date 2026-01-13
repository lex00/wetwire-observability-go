package grafana

// LogsPanel represents a Grafana logs panel.
type LogsPanel struct {
	BasePanel
	Options LogsOptions `json:"options,omitempty"`
}

// LogsOptions contains logs panel options.
type LogsOptions struct {
	ShowTime            bool   `json:"showTime,omitempty"`
	ShowLabels          bool   `json:"showLabels,omitempty"`
	ShowCommonLabels    bool   `json:"showCommonLabels,omitempty"`
	WrapLogMessage      bool   `json:"wrapLogMessage,omitempty"`
	PrettifyLogMessage  bool   `json:"prettifyLogMessage,omitempty"`
	EnableLogDetails    bool   `json:"enableLogDetails,omitempty"`
	SortOrder           string `json:"sortOrder,omitempty"`
	DedupStrategy       string `json:"dedupStrategy,omitempty"`
}

// Logs creates a new LogsPanel.
func Logs(title string) *LogsPanel {
	return &LogsPanel{
		BasePanel: BasePanel{
			Type:  "logs",
			Title: title,
			GridPos: GridPos{
				W: 24,
				H: 12,
			},
		},
		Options: LogsOptions{
			ShowTime:         true,
			WrapLogMessage:   true,
			EnableLogDetails: true,
			SortOrder:        "Descending",
		},
	}
}

// WithDescription sets the panel description.
func (p *LogsPanel) WithDescription(desc string) *LogsPanel {
	p.Description = desc
	return p
}

// WithDatasource sets the data source.
func (p *LogsPanel) WithDatasource(ds string) *LogsPanel {
	p.Datasource = ds
	return p
}

// WithSize sets the panel size.
func (p *LogsPanel) WithSize(w, h int) *LogsPanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *LogsPanel) WithPosition(x, y int) *LogsPanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithTargets sets the query targets.
func (p *LogsPanel) WithTargets(targets ...any) *LogsPanel {
	p.Targets = targets
	return p
}

// AddTarget adds a query target.
func (p *LogsPanel) AddTarget(target any) *LogsPanel {
	p.Targets = append(p.Targets, target)
	return p
}

// ShowTime shows the timestamp.
func (p *LogsPanel) ShowTime() *LogsPanel {
	p.Options.ShowTime = true
	return p
}

// HideTime hides the timestamp.
func (p *LogsPanel) HideTime() *LogsPanel {
	p.Options.ShowTime = false
	return p
}

// ShowLabels shows the labels.
func (p *LogsPanel) ShowLabels() *LogsPanel {
	p.Options.ShowLabels = true
	return p
}

// HideLabels hides the labels.
func (p *LogsPanel) HideLabels() *LogsPanel {
	p.Options.ShowLabels = false
	return p
}

// ShowCommonLabels shows common labels.
func (p *LogsPanel) ShowCommonLabels() *LogsPanel {
	p.Options.ShowCommonLabels = true
	return p
}

// HideCommonLabels hides common labels.
func (p *LogsPanel) HideCommonLabels() *LogsPanel {
	p.Options.ShowCommonLabels = false
	return p
}

// WrapLines enables line wrapping.
func (p *LogsPanel) WrapLines() *LogsPanel {
	p.Options.WrapLogMessage = true
	return p
}

// NoWrap disables line wrapping.
func (p *LogsPanel) NoWrap() *LogsPanel {
	p.Options.WrapLogMessage = false
	return p
}

// PrettifyJSON enables JSON prettification.
func (p *LogsPanel) PrettifyJSON() *LogsPanel {
	p.Options.PrettifyLogMessage = true
	return p
}

// EnableLogDetails enables log details view.
func (p *LogsPanel) EnableLogDetails() *LogsPanel {
	p.Options.EnableLogDetails = true
	return p
}

// DisableLogDetails disables log details view.
func (p *LogsPanel) DisableLogDetails() *LogsPanel {
	p.Options.EnableLogDetails = false
	return p
}

// SortDescending sorts logs in descending order (newest first).
func (p *LogsPanel) SortDescending() *LogsPanel {
	p.Options.SortOrder = "Descending"
	return p
}

// SortAscending sorts logs in ascending order (oldest first).
func (p *LogsPanel) SortAscending() *LogsPanel {
	p.Options.SortOrder = "Ascending"
	return p
}

// WithDedupStrategy sets the deduplication strategy.
func (p *LogsPanel) WithDedupStrategy(strategy string) *LogsPanel {
	p.Options.DedupStrategy = strategy
	return p
}

// DedupNone disables deduplication.
func (p *LogsPanel) DedupNone() *LogsPanel {
	p.Options.DedupStrategy = "none"
	return p
}

// DedupExact deduplicates exact matches.
func (p *LogsPanel) DedupExact() *LogsPanel {
	p.Options.DedupStrategy = "exact"
	return p
}

// DedupNumbers deduplicates lines that differ only in numbers.
func (p *LogsPanel) DedupNumbers() *LogsPanel {
	p.Options.DedupStrategy = "numbers"
	return p
}

// DedupSignature deduplicates by signature.
func (p *LogsPanel) DedupSignature() *LogsPanel {
	p.Options.DedupStrategy = "signature"
	return p
}
