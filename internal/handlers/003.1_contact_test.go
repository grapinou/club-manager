package handlers

import "testing"

func TestContactHandler(t *testing.T) {
	testHandler(t, "/contact", "Contact", ContactHandler)
}
