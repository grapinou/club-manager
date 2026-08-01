package handlers

import (
	"fmt"
	"net/http"
)

func RulesHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintln(w, "Règlement intérieur")
}
