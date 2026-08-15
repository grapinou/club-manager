package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/database"
	"github.com/grapinou/club-manager/internal/views"
)

func MembersHandler(cfg config.Config, queries database.Queries) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		members, err := queries.ListMembers(r.Context())
		if err != nil {
			http.Error(
				w,
				"impossible de récupérer la liste des membres",
				http.StatusInternalServerError,
			)
			return
		}

		var memberData []views.MemberData

		for _, member := range members {
			memberData = append(memberData, views.MemberData{
				ID:        member.ID,
				FirstName: member.FirstName,
				LastName:  member.LastName,
				BirthDate: member.BirthDate.Time.Format("02/01/2006"),
				Email:     member.Email.String,
			})
		}

		data := views.MembersData{
			SiteName: cfg.SiteName,
			Title:    "Membres - " + cfg.SiteName,
			Members:  memberData,
		}

		if err := views.RenderMembersList(w, data); err != nil {
			http.Error(
				w,
				"impossible d'afficher la liste des adhérents",
				http.StatusInternalServerError,
			)
			return
		}

	}
}
