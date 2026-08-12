package handlers

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestPostMemberHandler(t *testing.T) {

	// attention, la request est en post
	request := httptest.NewRequest(http.MethodPost, "/members", nil)

	response := httptest.NewRecorder()

	PostMemberHandler()(response, request)

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
