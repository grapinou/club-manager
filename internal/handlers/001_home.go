package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/views"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	data := views.HomeData{
		Title:       "Club Manager",
		Heading:     "Bienvenue sur Club Manager",
		Description: "Une application destinée à faciliter la gestion d'une association.",
	}

	err := views.RenderHome(w, data)

	if err != nil {
		http.Error(
			w,
			"Erreur interne du serveur",
			http.StatusInternalServerError,
		)
	}

}
