package router

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/handlers"
)

func New() *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("/{$}", handlers.HomeHandler)
	mux.HandleFunc("/club", handlers.ClubHandler)
	mux.HandleFunc("/contact", handlers.ContactHandler)
	mux.HandleFunc("/rules", handlers.RulesHandler)

	return mux

}
