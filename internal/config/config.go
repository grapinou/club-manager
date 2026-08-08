package config

import (
	"encoding/json"
	"fmt"
	"os"
)

// Config contient les paramètres généraux
// chargés au démarrage de l'application.
type Config struct {
	SiteName string      `json:"site_name"`
	Where    WhereConfig `json:"where"`
	When     WhenConfig  `json:"when"`
}

type WhereConfig struct {
	Heading     string `json:"heading"`
	Description string `json:"description"`
}

type WhenConfig struct {
	Heading     string `json:"heading"`
	Description string `json:"description"`
}

// Load ouvre le fichier situé au chemin reçu,
// décode son contenu JSON et retourne la configuration.
func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if err != nil {
		return Config{}, fmt.Errorf(
			"ouvrir le fichier de configuration : %w",
			err,
		)
	}
	defer file.Close()

	var cfg Config

	err = json.NewDecoder(file).Decode(&cfg)
	if err != nil {
		return Config{}, fmt.Errorf(
			"décoder le fichier de configuration : %w",
			err,
		)
	}

	return cfg, nil

}
