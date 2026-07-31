package handlers

import "testing"

func TestHomeHandler(t *testing.T) {
	testHandler(t, "/", "Bienvenue", HomeHandler)
}
