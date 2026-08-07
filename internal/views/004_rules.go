package views

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/layouts/base.html
//go:embed templates/pages/004_rules.html
var rulesFiles embed.FS

type RulesData struct {
	SiteName    string
	Title       string
	Heading     string
	Description string
}

var rulesTemplate = template.Must(
	template.ParseFS(rulesFiles,
		"templates/layouts/base.html",
		"templates/pages/004_rules.html"),
)

func RenderRules(w io.Writer, data RulesData) error {
	return rulesTemplate.ExecuteTemplate(
		w,
		"base",
		data,
	)
}
