package report

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestFontPathsRecursesAndSkipsHiddenDirectories(t *testing.T) {
	root := t.TempDir()
	writeTestFile(t, filepath.Join(root, "Root.TTF"))
	writeTestFile(t, filepath.Join(root, "nested", "icon.otf"))
	writeTestFile(t, filepath.Join(root, "nested", "web.woff2"))
	writeTestFile(t, filepath.Join(root, "nested", "ignore.txt"))
	writeTestFile(t, filepath.Join(root, ".hidden", "hidden.ttf"))

	got, err := fontPaths([]string{root})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		filepath.Join(root, "Root.TTF"),
		filepath.Join(root, "nested", "icon.otf"),
		filepath.Join(root, "nested", "web.woff2"),
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("fontPaths() = %#v, want %#v", got, want)
	}
}

func TestFontPathsAcceptsOnlyFontFileInputs(t *testing.T) {
	root := t.TempDir()
	font := filepath.Join(root, "font.ttf")
	text := filepath.Join(root, "notes.txt")
	writeTestFile(t, font)
	writeTestFile(t, text)

	got, err := fontPaths([]string{text, font})
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, []string{font}) {
		t.Fatalf("fontPaths() = %#v, want only %#v", got, font)
	}
}

func TestCategorizeRune(t *testing.T) {
	tests := map[rune]string{
		'A':      "Letters",
		'7':      "Numbers",
		'.':      "Punctuation",
		' ':      "Spacing",
		'\uE000': "Private Use / Icons",
		'℃':      "Symbols",
		'😀':      "Emoji",
	}

	for r, want := range tests {
		if got := categorizeRune(r); got != want {
			t.Fatalf("categorizeRune(%U) = %q, want %q", r, got, want)
		}
	}
}

func TestRuneName(t *testing.T) {
	tests := map[rune]string{
		' ':      "SPACE",
		'\n':     "CONTROL",
		'\uE000': "PRIVATE USE",
		'A':      "",
	}

	for r, want := range tests {
		if got := runeName(r); got != want {
			t.Fatalf("runeName(%U) = %q, want %q", r, got, want)
		}
	}
}

func writeTestFile(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte("x"), 0644); err != nil {
		t.Fatal(err)
	}
}
