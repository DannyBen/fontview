package report

import (
	"errors"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"

	"golang.org/x/image/font/sfnt"
)

var fontExts = map[string]bool{
	".ttf":   true,
	".otf":   true,
	".woff":  true,
	".woff2": true,
}

func loadFonts(inputs []string) ([]Font, error) {
	paths, err := fontPaths(inputs)
	if err != nil {
		return nil, err
	}
	if len(paths) == 0 {
		return nil, errors.New("no font files found; pass .ttf, .otf, .woff, or .woff2 files")
	}

	fonts := make([]Font, 0, len(paths))
	for i, path := range paths {
		font, err := inspectFont(path, i)
		if err != nil {
			fmt.Fprintf(os.Stderr, "fontview: skipping %s: %v\n", path, err)
			continue
		}
		fonts = append(fonts, font)
	}
	if len(fonts) == 0 {
		return nil, errors.New("no readable fonts found")
	}
	return fonts, nil
}

func fontPaths(inputs []string) ([]string, error) {
	if len(inputs) == 0 {
		return discoverFonts(".")
	}

	var paths []string
	for _, input := range inputs {
		info, err := os.Stat(input)
		if err != nil {
			return nil, err
		}
		if !info.IsDir() {
			if fontExts[strings.ToLower(filepath.Ext(input))] {
				paths = append(paths, input)
			}
			continue
		}
		found, err := discoverFonts(input)
		if err != nil {
			return nil, err
		}
		paths = append(paths, found...)
	}
	sort.Strings(paths)
	return paths, nil
}

func discoverFonts(root string) ([]string, error) {
	var paths []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if path != "." && strings.HasPrefix(filepath.Base(path), ".") {
				return filepath.SkipDir
			}
			return nil
		}
		if fontExts[strings.ToLower(filepath.Ext(path))] {
			paths = append(paths, path)
		}
		return nil
	})
	sort.Strings(paths)
	return paths, err
}

func inspectFont(path string, index int) (Font, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Font{}, err
	}
	font, err := sfnt.Parse(data)
	if err != nil {
		return Font{}, err
	}

	var buf sfnt.Buffer
	name := strings.TrimSpace(fontName(font, &buf))
	if name == "" {
		name = filepath.Base(path)
	}

	report := Font{
		ID:   fmt.Sprintf("font-%d", index),
		Name: name,
		Path: filepath.ToSlash(path),
	}

	for r := rune(0); r <= unicode.MaxRune; r++ {
		if !utf8.ValidRune(r) {
			continue
		}
		idx, err := font.GlyphIndex(&buf, r)
		if err != nil {
			return Font{}, err
		}
		if idx == 0 {
			continue
		}
		report.Glyphs = append(report.Glyphs, Glyph{
			Char:     string(r),
			Code:     fmt.Sprintf("U+%04X", r),
			CodeInt:  int(r),
			Category: categorizeRune(r),
			Name:     runeName(r),
		})
	}

	return report, nil
}

func fontName(font *sfnt.Font, buf *sfnt.Buffer) string {
	for _, id := range []sfnt.NameID{sfnt.NameIDTypographicFamily, sfnt.NameIDFamily} {
		name, err := font.Name(buf, id)
		if err == nil && strings.TrimSpace(name) != "" {
			return name
		}
	}
	return ""
}

func categorizeRune(r rune) string {
	switch {
	case isEmoji(r):
		return "Emoji"
	case unicode.In(r, unicode.Co):
		return "Private Use / Icons"
	case unicode.In(r, unicode.So):
		return "Symbols"
	case unicode.IsLetter(r):
		return "Letters"
	case unicode.IsNumber(r):
		return "Numbers"
	case unicode.IsPunct(r):
		return "Punctuation"
	case unicode.IsSpace(r):
		return "Spacing"
	case unicode.IsMark(r):
		return "Marks"
	default:
		return "Other"
	}
}

func isEmoji(r rune) bool {
	return inRanges(r, [][2]rune{
		{0x1F000, 0x1FAFF},
		{0x2600, 0x27BF},
		{0x2300, 0x23FF},
	})
}

func inRanges(r rune, ranges [][2]rune) bool {
	for _, rr := range ranges {
		if r >= rr[0] && r <= rr[1] {
			return true
		}
	}
	return false
}

func runeName(r rune) string {
	switch {
	case r == ' ':
		return "SPACE"
	case unicode.IsControl(r):
		return "CONTROL"
	case unicode.In(r, unicode.Co):
		return "PRIVATE USE"
	default:
		return ""
	}
}
