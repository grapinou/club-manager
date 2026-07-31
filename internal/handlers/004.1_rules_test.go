package handlers

import "testing"

func TestRulesHandler(t *testing.T) {
	testHandler(t, "/rules", "Réglement", RulesHandler)
}
