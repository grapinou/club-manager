package config

import "testing"

func TestLoad(t *testing.T) {

	cfg, err := Load("testdata/config.json")
	if err != nil {
		t.Fatalf(
			"le chargement de la configuration a échoué : %v",
			err,
		)
	}

	expectedSiteName := "Club Manager"
	expectedWhereHeading := "Où nous trouver ?"
	expectedWhenDescription := "Nous sommes disponibles tous les soirs de 20h à 22h. Pour davantage d'informations, nous contacter."

	if cfg.SiteName != expectedSiteName {
		t.Errorf(
			"nom du site obtenu : %q, nom du site attendu : %q",
			cfg.SiteName,
			expectedSiteName,
		)
	}

	if cfg.Where.Heading != expectedWhereHeading {
		t.Errorf(
			"where heading obtenu : %q, where heading attendu : %q",
			cfg.Where.Heading,
			expectedWhereHeading,
		)
	}

	if cfg.When.Description != expectedWhenDescription {
		t.Errorf(
			"when description obtenu : %q, when description attendu : %q",
			cfg.When.Description,
			expectedWhenDescription,
		)
	}

}

func TestLoadReturnsErrorWhenFileDoesNotExist(t *testing.T) {
	_, err := Load("testdata/unknown.json")

	if err == nil {
		t.Fatal(
			"une erreur était attendue pour un fichier inexistant",
		)
	}
}
