package router

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func testRoute(
	t *testing.T,
	route,
	expectedContent string,
) {

	t.Helper()

	// indépendance du test vis à vis de load.
	// construction de la donnée nécessaire pour tester la route

	homeCfg := config.HomeConfig{
		Title: "Accueil",
	}

	clubCfg := config.ClubConfig{
		Title: "Le Club",
	}

	contactCfg := config.ContactConfig{
		Title: "Contact",
	}

	rulesCfg := config.RulesConfig{
		Title: "Règlement",
	}
	whereCfg := config.WhereConfig{
		Title: "Où",
	}

	whenCfg := config.WhenConfig{
		Title: "Quand",
	}
	cfg := config.Config{
		SiteName: "Club Manager",
		Home:     homeCfg,
		Club:     clubCfg,
		Contact:  contactCfg,
		Rules:    rulesCfg,
		Where:    whereCfg,
		When:     whenCfg,
	}

	mux := New(cfg)

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
