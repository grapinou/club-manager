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

	if cfg.SiteName != expectedSiteName {
		t.Errorf(
			"nom du site obtenu : %q, nom du site attendu : %q",
			cfg.SiteName,
			expectedSiteName,
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
