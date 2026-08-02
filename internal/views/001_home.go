package views

import (
	"embed"
	"html/template"
	"io"
)

// homeFiles contient les fichiers intégrés correspondant
// au template de la page d'accueil.
//
//go:embed templates/home.html
var homeFiles embed.FS

// homeTemplate représente le fichier HTML analysé par Go.
//
// template.Must arrête immédiatement le programme si le template
// contient une erreur de syntaxe.
var homeTemplate = template.Must(
	template.ParseFS(homeFiles, "templates/home.html"),
)

// RenderHome exécute le template de la page d'accueil
// et écrit le résultat dans la destination reçue.
func RenderHome(w io.Writer) error {
	return homeTemplate.ExecuteTemplate(
		w,
		"home.html",
		nil,
	)
}
