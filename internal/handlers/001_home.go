package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func HomeHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.HomeData{
			SiteName:    cfg.SiteName,
			Title:       cfg.SiteName,
			Heading:     "Bienvenue sur " + cfg.SiteName,
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

}
