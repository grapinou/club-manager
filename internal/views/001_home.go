package views

import (
	"embed"
	"html/template"
	"io"
)

// homeFiles contient les fichiers intégrés correspondant
// au template de la page d'accueil.

//go:embed templates/layouts/base.html
//go:embed templates/pages/001_home.html
var homeFiles embed.FS

type HomeData struct {
	Title       string
	Heading     string
	Description string
}

// homeTemplate représente le fichier HTML analysé par Go.
//
// template.Must arrête immédiatement le programme si le template
// contient une erreur de syntaxe.
var homeTemplate = template.Must(
	template.ParseFS(homeFiles,
		"templates/layouts/base.html",
		"templates/pages/001_home.html"),
)

// RenderHome exécute le template de la page d'accueil
// et écrit le résultat dans la destination reçue.
func RenderHome(w io.Writer, data HomeData) error {
	return homeTemplate.ExecuteTemplate(
		w,
		"base", // donne le nom défni par {{ define }} et non le nom du fichier
		data,
	)
}
