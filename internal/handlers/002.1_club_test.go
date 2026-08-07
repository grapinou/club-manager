package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestClubHandler(t *testing.T) {

	// pour vérifier que la transmission de données provient de ce cfg
	// on le remplace par un autre nom
	cfg := config.Config{
		SiteName: "test sur /club",
	}

	testHandler(t, "/club", "Le club - "+cfg.SiteName, ClubHandler(cfg))
}
