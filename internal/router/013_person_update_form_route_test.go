package router

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestUpdatePersonFormRoute(t *testing.T) {
	cfg := config.Config{
		SiteName: "Club Manager",
	}

	queries := &FakeQueries{
		// ou ton fake adapté au package router
	}

	mux := New(cfg, queries)

	request := httptest.NewRequest(
		http.MethodGet,
		"/persons/42/edit",
		nil,
	)

	response := httptest.NewRecorder()

	mux.ServeHTTP(response, request)

	if response.Code != http.StatusOK {
		t.Errorf(
			"statut obtenu : %d, statut attendu : %d",
			response.Code,
			http.StatusOK,
		)
	}
}
