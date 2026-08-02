package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestUnknownRoute(t *testing.T) {
	mux := New()

	request := httptest.NewRequest(
		http.MethodGet,
		"/unknown",
		nil,
	)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusNotFound {
		t.Errorf(
			"statut obtenu : %d, statut attendu : %d",
			response.Code,
			http.StatusNotFound,
		)
	}
}
