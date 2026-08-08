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
			Title:       cfg.Home.Title + " - " + cfg.SiteName,
			Heading:     cfg.Home.Heading,
			Description: cfg.Home.Description,
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
