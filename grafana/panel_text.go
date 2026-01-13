package grafana

// Text mode options.
const (
	TextModeMarkdown = "markdown"
	TextModeHTML     = "html"
	TextModeCode     = "code"
)

// TextPanel represents a Grafana text panel.
type TextPanel struct {
	BasePanel
	Options TextOptions `json:"options,omitempty"`
}

// TextOptions contains text panel options.
type TextOptions struct {
	Mode    string      `json:"mode,omitempty"`
	Content string      `json:"content,omitempty"`
	Code    CodeOptions `json:"code,omitempty"`
}

// CodeOptions contains code display options.
type CodeOptions struct {
	Language        string `json:"language,omitempty"`
	ShowLineNumbers bool   `json:"showLineNumbers,omitempty"`
	ShowMiniMap     bool   `json:"showMiniMap,omitempty"`
}

// Text creates a new TextPanel.
func Text(title string) *TextPanel {
	return &TextPanel{
		BasePanel: BasePanel{
			Type:  "text",
			Title: title,
			GridPos: GridPos{
				W: 12,
				H: 6,
			},
		},
		Options: TextOptions{
			Mode: TextModeMarkdown,
		},
	}
}

// MarkdownText creates a text panel with markdown content.
func MarkdownText(content, title string) *TextPanel {
	return &TextPanel{
		BasePanel: BasePanel{
			Type:  "text",
			Title: title,
			GridPos: GridPos{
				W: 12,
				H: 6,
			},
		},
		Options: TextOptions{
			Mode:    TextModeMarkdown,
			Content: content,
		},
	}
}

// HTMLText creates a text panel with HTML content.
func HTMLText(content, title string) *TextPanel {
	return &TextPanel{
		BasePanel: BasePanel{
			Type:  "text",
			Title: title,
			GridPos: GridPos{
				W: 12,
				H: 6,
			},
		},
		Options: TextOptions{
			Mode:    TextModeHTML,
			Content: content,
		},
	}
}

// WithDescription sets the panel description.
func (p *TextPanel) WithDescription(desc string) *TextPanel {
	p.Description = desc
	return p
}

// WithSize sets the panel size.
func (p *TextPanel) WithSize(w, h int) *TextPanel {
	p.GridPos.W = w
	p.GridPos.H = h
	return p
}

// WithPosition sets the panel position.
func (p *TextPanel) WithPosition(x, y int) *TextPanel {
	p.GridPos.X = x
	p.GridPos.Y = y
	return p
}

// WithContent sets the text content.
func (p *TextPanel) WithContent(content string) *TextPanel {
	p.Options.Content = content
	return p
}

// Markdown sets the mode to markdown.
func (p *TextPanel) Markdown() *TextPanel {
	p.Options.Mode = TextModeMarkdown
	return p
}

// HTML sets the mode to HTML.
func (p *TextPanel) HTML() *TextPanel {
	p.Options.Mode = TextModeHTML
	return p
}

// Code sets the mode to code.
func (p *TextPanel) Code() *TextPanel {
	p.Options.Mode = TextModeCode
	return p
}

// WithCodeLanguage sets the code language for syntax highlighting.
func (p *TextPanel) WithCodeLanguage(lang string) *TextPanel {
	p.Options.Code.Language = lang
	return p
}

// ShowLineNumbers shows line numbers in code mode.
func (p *TextPanel) ShowLineNumbers() *TextPanel {
	p.Options.Code.ShowLineNumbers = true
	return p
}

// HideLineNumbers hides line numbers in code mode.
func (p *TextPanel) HideLineNumbers() *TextPanel {
	p.Options.Code.ShowLineNumbers = false
	return p
}

// ShowMiniMap shows the mini map in code mode.
func (p *TextPanel) ShowMiniMap() *TextPanel {
	p.Options.Code.ShowMiniMap = true
	return p
}

// HideMiniMap hides the mini map in code mode.
func (p *TextPanel) HideMiniMap() *TextPanel {
	p.Options.Code.ShowMiniMap = false
	return p
}

// MakeTransparent makes the panel background transparent.
func (p *TextPanel) MakeTransparent() *TextPanel {
	p.Transparent = true
	return p
}
