package main

import (
	"fmt"
	"net/http"

	"github.com/grapinou/club-manager/internal/router"
)

func main() {

	mux := router.New()

	fmt.Println("Serveur lancé sur http://localhost:8080")

	http.ListenAndServe(":8080", mux)
}
