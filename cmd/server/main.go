package main

import (
	"fmt"
	"net/http"

	"github.com/grapinou/club-manager/internal/handlers"
)

func main() {

	http.HandleFunc("/", handlers.HomeHandler)
	http.HandleFunc("/club", handlers.ClubHandler)
	http.HandleFunc("/contact", handlers.ContactHandler)
	http.HandleFunc("/rules", handlers.RulesHandler)

	fmt.Println("Serveur lancé sur http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
