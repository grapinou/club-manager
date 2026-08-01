package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func testRoute(
	t *testing.T,
	route,
	expectedContent string,
) {

	t.Helper()

	mux := New()

	request := httptest.NewRequest(http.MethodGet, route, nil)
	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"statut obtenu : %d, statut attendu : %d",
			response.Code,
			http.StatusOK,
		)
	}

	body := response.Body.String()

	if !strings.Contains(body, expectedContent) {
		t.Errorf(
			"la réponse ne contient pas le texte %q ; contenu obtenu : %q",
			expectedContent,
			body,
		)
	}

}
