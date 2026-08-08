package views

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/layouts/base.html
//go:embed templates/pages/006_where.html
var whereFiles embed.FS

type WhereData struct {
	SiteName    string
	Title       string
	Heading     string
	Description string
}

var whereTemplate = template.Must(
	template.ParseFS(whereFiles,
		"templates/layouts/base.html",
		"templates/pages/006_where.html"),
)

func RenderWhere(w io.Writer, data WhereData) error {
	return whereTemplate.ExecuteTemplate(
		w,
		"base",
		data)
}
