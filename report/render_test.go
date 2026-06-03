package report

import (
	"strings"
	"testing"
)

func TestRenderIncludesPayloadStylesAndScript(t *testing.T) {
	page, err := render([]Font{
		{
			ID:   "font-0",
			Name: "Example",
			Path: "Example.ttf",
			Glyphs: []Glyph{
				{Char: "A", Code: "U+0041", CodeInt: 65, Category: "Letters"},
			},
		},
	}, "test-version")
	if err != nil {
		t.Fatal(err)
	}

	html := string(page)
	for _, want := range []string{
		"<title>Font Viewer</title>",
		`"name":"Example"`,
		`"code":"U+0041"`,
		".toolbar",
		"function glyphButton",
	} {
		if !strings.Contains(html, want) {
			t.Fatalf("render() output missing %q", want)
		}
	}
}

func TestServerFontsRewritesPathsWithoutMutatingOriginal(t *testing.T) {
	fonts := []Font{{ID: "font-0", Path: "nested/Example.ttf"}}

	got := serverFonts(fonts)
	if got[0].Path != "/font/nested/Example.ttf" {
		t.Fatalf("serverFonts()[0].Path = %q", got[0].Path)
	}
	if fonts[0].Path != "nested/Example.ttf" {
		t.Fatalf("serverFonts() mutated original path to %q", fonts[0].Path)
	}
}
