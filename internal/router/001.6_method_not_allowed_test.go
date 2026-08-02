package router

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestMethodNotAllowed(t *testing.T) {
	mux := New()

	request := httptest.NewRequest(
		http.MethodPost,
		"/contact",
		nil,
	)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusMethodNotAllowed {
		t.Errorf(
			"statut obtenu : %d, statut attendu : %d",
			response.Code,
			http.StatusMethodNotAllowed,
		)
	}
}
