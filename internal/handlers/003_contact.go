package handlers

import (
	"fmt"
	"net/http"
)

func ContactHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Contact")
}
