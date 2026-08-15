package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/database"
)

func MembersHandler(queries database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, err := queries.ListMembers(r.Context())
		if err != nil {
			http.Error(
				w,
				"impossible de récupérer la liste des membres",
				http.StatusInternalServerError,
			)
			return
		}
	}
}
