package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Club Manager")
}

func clubPresentationHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Présentation de Club Manager")
}

func contactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Pour contacter Club Manager")
}

func rulesHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Réglement de Club Manager")
}

func main() {

	http.HandleFunc("/", homeHandler)
	http.HandleFunc("/clubPresentation", clubPresentationHandler)
	http.HandleFunc("/contact", contactHandler)
	http.HandleFunc("/rules", rulesHandler)

	fmt.Println("Serveur lancé sur http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
