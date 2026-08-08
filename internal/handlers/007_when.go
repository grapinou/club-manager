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
			Heading:     "Quand nous trouver ?",
			Description: "Nous sommes disponible tous les soir de 20h à 22h. Pour davantage d'information, nous contacter.",
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
