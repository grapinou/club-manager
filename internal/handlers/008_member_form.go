package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/views"
)

func MemberFormHandler(cfg config.Config) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		data := views.MemberFormData{
			SiteName: cfg.SiteName,
			Title:    "Ajouter un membre - " + cfg.SiteName,
		}

		if err := views.RenderMemberForm(w, data); err != nil {
			http.Error(
				w,
				"impossible d'afficher le formulaire membre",
				http.StatusInternalServerError,
			)
			return
		}
	}
}
