package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/database"
)

func PostPersonHandler(queries database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	}
}
