package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func ContactHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		data := views.ContactData{
			SiteName:     cfg.SiteName,
			Title:        "Contact - " + cfg.SiteName,
			Heading:      "Comment nous contacter ?",
			Description:  "Vous pouvez nous joindre par téléphone ou par mail",
			EmailAddress: "clubmanager@mail.com",
			PhoneNumber:  "07-00-00-00-07",
		}

		err := views.RenderContact(w, data)

		if err != nil {
			http.Error(
				w,
				"Erreur interne du serveur",
				http.StatusInternalServerError,
			)
		}
	}

}
