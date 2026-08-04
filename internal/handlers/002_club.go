package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/views"
)

func ClubHandler(w http.ResponseWriter, r *http.Request) {

	data := views.ClubData{
		Title:       "Le club - Club Manager",
		Heading:     "Présentation du club",
		Description: "Découvrez l'association, son histoire et ses valeurs.",
	}

	err := views.RenderClub(w, data)

	if err != nil {
		http.Error(
			w,
			"Erreur interne du serveur",
			http.StatusInternalServerError,
		)
	}

}
