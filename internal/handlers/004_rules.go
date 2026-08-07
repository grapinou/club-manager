package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/views"
)

func RulesHandler(w http.ResponseWriter, r *http.Request) {

	data := views.RulesData{
		Title:       "Règlement",
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
