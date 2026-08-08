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
			Title:       "Où - " + cfg.SiteName,
			Heading:     "Où nous trouver ?",
			Description: "Nous sommes au 3 rue de Perlinpinpin",
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
