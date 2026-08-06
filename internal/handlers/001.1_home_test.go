package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestHomeHandler(t *testing.T) {

	cfg := config.Config{
		SiteName: "Club Manager",
	}

	testHandler(t, "/", "Bienvenue", HomeHandler(cfg))
}
