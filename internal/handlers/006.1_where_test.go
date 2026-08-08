package handlers

import (
	"testing"

	"github.com/grapinou/club-manager/internal/config"
)

func TestWhereHandler(t *testing.T) {

	cfg := config.Config{
		SiteName: "Club Manager",
	}

	testHandler(t, "/", "Où - "+cfg.SiteName, WhereHandler(cfg))
}
