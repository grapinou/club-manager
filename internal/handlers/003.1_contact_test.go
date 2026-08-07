package handlers

import "testing"

func TestContactHandler(t *testing.T) {
	testHandler(t, "/contact", "Comment nous contacter ?", ContactHandler)
}
