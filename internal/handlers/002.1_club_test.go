package handlers

import "testing"

func TestClubHandler(t *testing.T) {
	testHandler(t, "/club", "Présentation", ClubHandler)
}
