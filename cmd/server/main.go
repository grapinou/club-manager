package main

import (
	"log"
	"net/http"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/router"
)

const configPath = "config/config.json"

func main() {

	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatalf("impossible de charger la configuration depuis %q : %v",
			configPath,
			err,
		)
	}

	mux := router.New(cfg)

	log.Println("Serveur lancé sur http://localhost:8080")

	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal(err)
	}

}
