package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/views"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {

	data := views.ContactData{
		Title:        "Contact",
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
