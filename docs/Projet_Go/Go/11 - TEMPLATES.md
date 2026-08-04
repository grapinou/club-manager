
---

- [[11.01 - embed]]
- [[11.02 - ParseFS]]
- [[11.03 - Passage de données]]
- [[11.04 - Layout commun et composition des templates]]
## Architecture obtenue

```
club-manager/
├── cmd/
│   └── server/
│       └── main.go
│
└── internal/
    ├── handlers/
    │   └── 001_home.go
    │
    ├── router/
    │   └── 001_router.go
    │
    └── views/
        ├── 001_home.go
        │
        └── templates/
            └── home.html
```

`001_home.go` : la vue

```go

package views

import (
	"embed"
	"html/template"
	"io"
)

// homeFiles contient les fichiers intégrés correspondant
// au template de la page d'accueil.
//
//go:embed templates/home.html
var homeFiles embed.FS

// homeTemplate représente le fichier HTML analysé par Go.
//
// template.Must arrête immédiatement le programme si le template
// contient une erreur de syntaxe.
var homeTemplate = template.Must(
	template.ParseFS(homeFiles, "templates/home.html"),
)

// RenderHome exécute le template de la page d'accueil
// et écrit le résultat dans la destination reçue.
func RenderHome(w io.Writer) error {
	return homeTemplate.ExecuteTemplate(
		w,
		"home.html",
		nil,
	)
}

```

le html : 

```html 
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">

    <title>Club Manager</title>
</head>

<body>
    <h1>Bienvenue sur Club Manager</h1>

    <p>
        Une application destinée à faciliter la gestion d'une association.
    </p>
</body>
</html>

```

le handler :

```go 

package handlers

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/views"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {

	err := views.RenderHome(w)

	if err != nil {
		http.Error(
			w,
			"Erreur interne du serveur",
			http.StatusInternalServerError,
		)
	}

}


```

