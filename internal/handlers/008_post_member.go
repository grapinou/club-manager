package handlers

import (
	"fmt"
	"net/http"
)

func PostMemberHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "post membre")
	}
}
