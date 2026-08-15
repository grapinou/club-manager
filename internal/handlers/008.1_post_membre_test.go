package handlers

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/grapinou/club-manager/internal/database/dbsqlc"
)

// recordingQueries enregistre les paramètres reçus par CreateMember
// afin de vérifier ce que le handler transmet à la couche base de données.
type recordingQueries struct {
	CreateMemberParams dbsqlc.CreateMemberParams
}

func (q *recordingQueries) CreateMember(ctx context.Context, arg dbsqlc.CreateMemberParams) (dbsqlc.Member, error) {

	q.CreateMemberParams = arg

	return dbsqlc.Member{}, nil
}

func (q *recordingQueries) ListMembers(ctx context.Context) ([]dbsqlc.Member, error) {
	return nil, nil
}

func TestPostMemberHandler(t *testing.T) {

	form := url.Values{}

	firstName := "Robin"
	lastName := "Des Bois"
	birthdate := "1990-05-12"
	email := "robin.desbois@example.com"

	form.Set("FirstName", firstName)
	form.Set("LastName", lastName)
	form.Set("Birthdate", birthdate)
	form.Set("Email", email)

	// attention, la request est en post
	request := httptest.NewRequest(http.MethodPost, "/members", strings.NewReader(form.Encode()))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	queries := &recordingQueries{}

	PostMemberHandler(queries)(response, request)

	if response.Code != http.StatusSeeOther {
		t.Errorf(
			"statut obtenu : %d, statut attendu : %d",
			response.Code,
			http.StatusSeeOther,
		)
	}

	location := response.Header().Get("Location")

	if location != "/members/new" {
		t.Errorf(
			"redirection obtenue %q ; attendue : %q",
			location,
			"/members/new",
		)
	}

	if queries.CreateMemberParams.FirstName != firstName {
		t.Errorf(
			"prénom obtenu : %q, attendu : %q",
			queries.CreateMemberParams.FirstName,
			firstName,
		)
	}

	if queries.CreateMemberParams.LastName != lastName {
		t.Errorf(
			"nom obtenu : %q, attendu : %q",
			queries.CreateMemberParams.LastName,
			lastName,
		)
	}

	if !queries.CreateMemberParams.BirthDate.Valid {
		t.Error("la date de naissance devrait être valide")
	}

	if queries.CreateMemberParams.BirthDate.Time.Format("2006-01-02") != birthdate {
		t.Errorf(
			"date obtenue : %q, attendue : %q",
			queries.CreateMemberParams.BirthDate.Time.Format("2006-01-02"),
			birthdate,
		)
	}

	if !queries.CreateMemberParams.Email.Valid {
		t.Error("l'email devrait être valide")
	}

	if queries.CreateMemberParams.Email.String != email {
		t.Errorf(
			"email obtenu : %q, attendu : %q",
			queries.CreateMemberParams.Email.String,
			email,
		)
	}
}
