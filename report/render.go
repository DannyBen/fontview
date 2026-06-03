package report

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"path/filepath"
)

//go:embed templates/page.html templates/styles/*.css templates/assets/*.js
var templateFS embed.FS

type pageData struct {
	Title     string
	Version   string
	Fonts     template.JS
	Styles    template.CSS
	Script    template.JS
	Generated string
}

func render(fonts []Font, version string) ([]byte, error) {
	payload, err := json.Marshal(fonts)
	if err != nil {
		return nil, err
	}

	styles, err := templateFS.ReadFile("templates/styles/app.css")
	if err != nil {
		return nil, err
	}
	script, err := templateFS.ReadFile("templates/assets/app.js")
	if err != nil {
		return nil, err
	}
	tmpl, err := template.ParseFS(templateFS, "templates/page.html")
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = tmpl.Execute(&buf, pageData{
		Title:   "Font Viewer",
		Version: version,
		Fonts:   template.JS(payload),
		Styles:  template.CSS(styles),
		Script:  template.JS(script),
	})
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func serverFonts(fonts []Font) []Font {
	server := make([]Font, len(fonts))
	copy(server, fonts)
	for i := range server {
		server[i].Path = "/font/" + filepath.ToSlash(server[i].Path)
	}
	return server
}
