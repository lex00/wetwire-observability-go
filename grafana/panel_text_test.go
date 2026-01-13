package grafana

import "testing"

func TestText(t *testing.T) {
	p := Text("Instructions")
	if p.Title != "Instructions" {
		t.Errorf("Title = %q, want Instructions", p.Title)
	}
	if p.Type != "text" {
		t.Errorf("Type = %q, want text", p.Type)
	}
}

func TestText_WithDescription(t *testing.T) {
	p := Text("Test").WithDescription("A test text panel")
	if p.Description != "A test text panel" {
		t.Errorf("Description = %q, want A test text panel", p.Description)
	}
}

func TestText_WithSize(t *testing.T) {
	p := Text("Test").WithSize(24, 4)
	if p.GridPos.W != 24 || p.GridPos.H != 4 {
		t.Errorf("GridPos = %+v, want W=24 H=4", p.GridPos)
	}
}

func TestText_WithContent(t *testing.T) {
	p := Text("Test").WithContent("# Hello World")
	if p.Options.Content != "# Hello World" {
		t.Errorf("Content = %q, want # Hello World", p.Options.Content)
	}
}

func TestText_Markdown(t *testing.T) {
	p := Text("Test").Markdown()
	if p.Options.Mode != "markdown" {
		t.Errorf("Mode = %q, want markdown", p.Options.Mode)
	}
}

func TestText_HTML(t *testing.T) {
	p := Text("Test").HTML()
	if p.Options.Mode != "html" {
		t.Errorf("Mode = %q, want html", p.Options.Mode)
	}
}

func TestText_Code(t *testing.T) {
	p := Text("Test").Code()
	if p.Options.Mode != "code" {
		t.Errorf("Mode = %q, want code", p.Options.Mode)
	}
}

func TestText_WithCodeLanguage(t *testing.T) {
	p := Text("Test").Code().WithCodeLanguage("go")
	if p.Options.Code.Language != "go" {
		t.Errorf("Code.Language = %q, want go", p.Options.Code.Language)
	}
}

func TestText_ShowLineNumbers(t *testing.T) {
	p := Text("Test").Code().ShowLineNumbers()
	if !p.Options.Code.ShowLineNumbers {
		t.Error("ShowLineNumbers should be true")
	}
}

func TestText_FluentAPI(t *testing.T) {
	p := Text("Info").
		WithDescription("Important info").
		WithSize(12, 6).
		Markdown().
		WithContent("## Overview\n\nThis is a test.")

	if p.Title != "Info" {
		t.Error("Fluent API should preserve title")
	}
	if p.Options.Mode != "markdown" {
		t.Error("Fluent API should set mode")
	}
	if p.Options.Content != "## Overview\n\nThis is a test." {
		t.Error("Fluent API should set content")
	}
}

func TestMarkdownText(t *testing.T) {
	p := MarkdownText("# Title", "Welcome")
	if p.Title != "Welcome" {
		t.Errorf("Title = %q, want Welcome", p.Title)
	}
	if p.Options.Mode != "markdown" {
		t.Errorf("Mode = %q, want markdown", p.Options.Mode)
	}
	if p.Options.Content != "# Title" {
		t.Errorf("Content = %q, want # Title", p.Options.Content)
	}
}

func TestHTMLText(t *testing.T) {
	p := HTMLText("<h1>Title</h1>", "Welcome")
	if p.Title != "Welcome" {
		t.Errorf("Title = %q, want Welcome", p.Title)
	}
	if p.Options.Mode != "html" {
		t.Errorf("Mode = %q, want html", p.Options.Mode)
	}
	if p.Options.Content != "<h1>Title</h1>" {
		t.Errorf("Content = %q, want <h1>Title</h1>", p.Options.Content)
	}
}
