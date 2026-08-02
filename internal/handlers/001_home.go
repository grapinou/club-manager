package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/views"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	err := views.RenderHome(w)

	if err != nil {
		http.Error(
			w,
			"Erreur interne du serveur",
			http.StatusInternalServerError,
		)
	}

}
