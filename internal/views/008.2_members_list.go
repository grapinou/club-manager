package views

import (
	"embed"
	"html/template"
	"io"
)

//go:embed templates/layouts/base.html
//go:embed templates/pages/008.2_members_list.html
var membersListFiles embed.FS

type MemberData struct {
	ID        int32
	FirstName string
	LastName  string
	BirthDate string
	Email     string
	CreatedAt string
}

type MembersData struct {
	SiteName string
	Title    string
	Members  []MemberData
}

var membersListTemplate = template.Must(
	template.ParseFS(
		membersListFiles,
		"templates/layouts/base.html",
		"templates/pages/008.2_members_list.html",
	),
)

func RenderMembersList(w io.Writer, data MembersData) error {
	return membersListTemplate.ExecuteTemplate(
		w,
		"base",
		data,
	)
}
