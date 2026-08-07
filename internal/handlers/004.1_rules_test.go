package handlers

import "testing"

func TestRulesHandler(t *testing.T) {
	testHandler(t, "/rules", "Règlement intérieur", RulesHandler)
}
