package router

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/handlers"
)

func New(cfg config.Config) *http.ServeMux {

	_ = cfg

	mux := http.NewServeMux()

	mux.HandleFunc("GET /{$}", handlers.HomeHandler)
	mux.HandleFunc("GET /club", handlers.ClubHandler)
	mux.HandleFunc("GET /contact", handlers.ContactHandler)
	mux.HandleFunc("GET /rules", handlers.RulesHandler)

	return mux

}
