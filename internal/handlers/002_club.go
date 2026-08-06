package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func ClubHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.ClubData{
			Title:       "Le club - " + cfg.SiteName,
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
}
