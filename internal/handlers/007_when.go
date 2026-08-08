package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func WhenHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.WhenData{
			SiteName:    cfg.SiteName,
			Title:       "Quand - " + cfg.SiteName,
			Heading:     cfg.When.Heading,
			Description: cfg.When.Description,
		}

		err := views.RenderWhen(w, data)

		if err != nil {
			http.Error(
				w,
				"Erreur interne du serveur",
				http.StatusInternalServerError,
			)
		}
	}
}
