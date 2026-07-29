package handlers

import (
	"fmt"
	"net/http"
)

func ClubHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Présentation du club")
}
