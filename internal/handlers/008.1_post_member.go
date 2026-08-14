package handlers

import (
	"net/http"
	"time"

	"github.com/grapinou/club-manager/internal/database"
	"github.com/grapinou/club-manager/internal/database/dbsqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

func PostMemberHandler(queries database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		birthDate, err := time.Parse("2006-01-02", r.FormValue("Birthdate"))
		if err != nil {
			http.Error(w, "date de naissance invalide", http.StatusBadRequest)
			return
		}

		datas := dbsqlc.CreateMemberParams{
			FirstName: r.FormValue("FirstName"),
			LastName:  r.FormValue("LastName"),

			BirthDate: pgtype.Date{
				Time:  birthDate,
				Valid: true,
			},
			Email: pgtype.Text{
				String: r.FormValue("Email"),
				Valid:  true,
			},
		}

		_, err = queries.CreateMember(r.Context(), datas)
		if err != nil {
			http.Error(w, "db insert non fait", http.StatusInternalServerError)
			return
		}

		http.Redirect(
			w,
			r,
			"/members/new",
			http.StatusSeeOther,
		)
	}
}
