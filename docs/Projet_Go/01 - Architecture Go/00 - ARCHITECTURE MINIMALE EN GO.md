
Tag : #serveur ; #go ; #projet

---

# Étape pratique 1

Depuis le dossier du projet :

```
mkdir -p cmd/server
mkdir -p internal/handlers
mkdir -p internal/config
mkdir -p web/templates
```

Remarque : [[00.01 - mkdir -p]]

Puis vérifier :

```
tree
```

Si `tree` n'est pas installé :

```
sudo apt install tree
```

# Principe

Un projet Go professionnel sépare :

- le lancement de l'application ;
- le code interne ;
- les ressources web.

```
club-manager/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── handlers/
│   │
│   └── config/
│
├── web/
│   └── templates/
│
├── go.mod
├── go.sum
├── README.md
└── .gitignore
```


---

## Dossiers principaux

| Dossier   | Rôle                                |
| --------- | ----------------------------------- |
| cmd/      | Contient les programmes exécutables |
| internal/ | Code privé de l'application         |
| web/      | Fichiers utilisés par le navigateur |

---

## cmd/

Exemple :

```
cmd/
└── server/
    └── main.go
```

Contient le point d'entrée.

C'est ici que démarre notre serveur.

Commande :

```
go run ./cmd/server
```


De manière générale : 

Go va :

1. Chercher le package `main`.
2. Compiler toutes les dépendances.
3. Exécuter `main()`.

Le dossier `cmd/server` représente donc **un exécutable**.

Plus tard, ton projet pourrait ressembler à ceci :

```
club-manager/
│
├── cmd/
│   ├── server/
│   │   └── main.go
│   │
│   ├── migrate/
│   │   └── main.go
│   │
│   └── seed/
│       └── main.go
│
├── internal/
│
└── ...
```

Tu aurais alors plusieurs programmes :

- `server` → lance le site.
- `migrate` → met à jour la base de données.
- `seed` → insère des données de démonstration.

Tous partageraient le même code situé dans `internal`.

C'est une architecture très courante dans les projets Go.

---

## internal/

Contient la logique de l'application.

Exemple futur :

```
internal/
├── handlers/
├── service/
├── repository/
└── model/
```

---

### handlers

Responsable des requêtes HTTP.

Exemple :

```
GET /
```

Retourne une page d'accueil.

---

### service

Contient la logique métier.

Exemple :

- calcul d'une cotisation ;
- validation d'une inscription.

---

### repository

Communication avec la base de données.

Exemple futur :

```
PostgreSQL
     |
repository
     |
service
     |
handler
```

---

## web/

Contient les fichiers visibles par le navigateur.

Exemple :

```
web/
└── templates/
    └── home.html
```

Plus tard :

```
web/
├── templates/
├── static/
│   ├── css/
│   └── images/
```

# Étape 2 — Premier SERVEUR GO

Créer :

```
cmd/server/main.go
```

Avec :

```go
package main

import (
	"fmt"
	"net/http"
)

func homeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Club Manager")
}

func main() {

	http.HandleFunc("/", homeHandler)

	fmt.Println("Serveur lancé sur http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
```

---

Lancement :

```
go run ./cmd/server
```

Résultat attendu :

```
Serveur lancé sur http://localhost:8080
```

Puis navigateur :

```
http://localhost:8080
```

Affichage :

```
Bienvenue sur Club Manager
```
