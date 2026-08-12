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

// fake queries côté handler
type FakeQueries struct{}

func (f FakeQueries) CreateMember(ctx context.Context, arg dbsqlc.CreateMemberParams) (dbsqlc.Member, error) {
	return dbsqlc.Member{}, nil
}

func TestPostMemberHandler(t *testing.T) {

	form := url.Values{}

	form.Set("FirstName", "Robin")
	form.Set("LastName", "Des Bois")
	form.Set("Birthdate", "1990-05-12")
	form.Set("Email", "robin.desbois@example.com")

	// attention, la request est en post
	request := httptest.NewRequest(http.MethodPost, "/members", strings.NewReader(form.Encode()))

	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	response := httptest.NewRecorder()

	queries := FakeQueries{}

	PostMemberHandler(queries)(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"statut obtenu : %d, statut attendu : %d",
			response.Code,
			http.StatusOK,
		)
	}

	body := response.Body.String()

	if !strings.Contains(body, "post membre") {
		t.Errorf(
			"la réponse ne contient pas %q ; contenu obtenu : %q",
			"post membre",
			body,
		)
	}

}
