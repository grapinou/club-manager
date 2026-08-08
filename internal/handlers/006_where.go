package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func WhereHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.WhereData{
			SiteName:    cfg.SiteName,
			Title:       cfg.Where.Title + " - " + cfg.SiteName,
			Heading:     cfg.Where.Heading,
			Description: cfg.Where.Description,
		}

		err := views.RenderWhere(w, data)

		if err != nil {
			http.Error(
				w,
				"Erreur interne du serveur",
				http.StatusInternalServerError,
			)
		}

	}

}
