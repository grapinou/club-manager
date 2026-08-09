package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func RulesHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.PageData{
			SiteName:    cfg.SiteName,
			Title:       cfg.Rules.Title + " - " + cfg.SiteName,
			Heading:     cfg.Rules.Heading,
			Description: cfg.Rules.Description,
		}

		err := views.RenderPage(w, data)

		if err != nil {
			http.Error(
				w,
				"Erreur interne du serveur",
				http.StatusInternalServerError,
			)
		}
	}

}
