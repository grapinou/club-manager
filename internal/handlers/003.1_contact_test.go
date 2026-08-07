package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestContactHandler(t *testing.T) {

	cfg := config.Config{
		SiteName: "test sur /contact",
	}

	testHandler(t, "/contact", "Contact - "+cfg.SiteName, ContactHandler(cfg))
}
