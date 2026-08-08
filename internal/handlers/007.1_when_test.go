package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestWhenHandler(t *testing.T) {

	cfg := config.Config{
		SiteName: "Club Manager",
	}

	testHandler(t, "/", "Quand - "+cfg.SiteName, WhenHandler(cfg))
}
