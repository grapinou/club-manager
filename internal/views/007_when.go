package views

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/layouts/base.html
//go:embed templates/pages/007_when.html
var whenFiles embed.FS

type WhenData struct {
	SiteName    string
	Title       string
	Heading     string
	Description string
}

var whenTemplate = template.Must(
	template.ParseFS(whenFiles,
		"templates/layouts/base.html",
		"templates/pages/007_when.html"),
)

func RenderWhen(w io.Writer, data WhenData) error {
	return whenTemplate.ExecuteTemplate(
		w,
		"base",
		data,
	)
}
