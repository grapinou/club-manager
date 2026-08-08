package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestRulesHandler(t *testing.T) {

	rulesCfg := config.RulesConfig{
		Title: "Règlement de test",
	}

	cfg := config.Config{
		SiteName: "test sur /rules",
		Rules:    rulesCfg,
	}

	testHandler(t, "/rules", "Règlement de test", RulesHandler(cfg))
}
