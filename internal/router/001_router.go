package router

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/handlers"
)

func New(cfg config.Config) *http.ServeMux {

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.HomeHandler(cfg))
	mux.HandleFunc("GET /club", handlers.ClubHandler(cfg))
	mux.HandleFunc("GET /contact", handlers.ContactHandler(cfg))
	mux.HandleFunc("GET /where", handlers.WhereHandler(cfg))
	mux.HandleFunc("GET /when", handlers.WhenHandler(cfg))
	mux.HandleFunc("GET /rules", handlers.RulesHandler(cfg))

	return mux

}
