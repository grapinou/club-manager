package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestRulesHandler(t *testing.T) {

	cfg := config.Config{
		SiteName: "test sur /rules",
	}

	testHandler(t, "/rules", "Règlement - "+cfg.SiteName, RulesHandler(cfg))
}
