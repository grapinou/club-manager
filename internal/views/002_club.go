package views

import (
	"embed"
	"html/template"
	"io"
)

// clubFiles contient le layout commun et le template
// correspondant au contenu de la page du club.
//
//go:embed templates/layouts/base.html
//go:embed templates/pages/002_club.html
var clubFiles embed.FS

// ClubData contient les données nécessaires
// à l'affichage de la page du club.
type ClubData struct {
	SiteName    string
	Title       string
	Heading     string
	Description string
}

// clubTemplate rassemble le layout commun et le contenu
// propre à la page du club.
//
// template.Must arrête immédiatement le programme si
// l'un des templates contient une erreur de syntaxe.
var clubTemplate = template.Must(
	template.ParseFS(
		clubFiles,
		"templates/layouts/base.html",
		"templates/pages/002_club.html",
	),
)

// RenderClub exécute le layout commun avec les données
// propres à la page du club.
func RenderClub(w io.Writer, data ClubData) error {
	return clubTemplate.ExecuteTemplate(
		w,
		"base",
		data,
	)
}
