package views

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/layouts/base.html
//go:embed templates/pages/008_member_form.html
var memberFormFiles embed.FS

type MemberFormData struct {
	SiteName string
	Title    string
}

var memberFormTemplate = template.Must(
	template.ParseFS(
		memberFormFiles,
		"templates/layouts/base.html",
		"templates/pages/008_member_form.html",
	),
)

func RenderMemberForm(w io.Writer, data MemberFormData) error {
	return memberFormTemplate.ExecuteTemplate(
		w,
		"base",
		data,
	)
}
