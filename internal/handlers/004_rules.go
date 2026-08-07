package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func RulesHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.RulesData{
			SiteName:    cfg.SiteName,
			Title:       "Règlement - " + cfg.SiteName,
			Heading:     "Règlement intérieur",
			Description: "Tous les membres doivent respecter le règlement intérieur suivant :",
		}

		err := views.RenderRules(w, data)

		if err != nil {
			http.Error(
				w,
				"Erreur interne du serveur",
				http.StatusInternalServerError,
			)
		}
	}

}
