
---


## Objectif

Cette fiche présente une nouvelle étape dans l’architecture de **Club Manager**.

Elle ne remplace pas la fiche `00 - Architecture minimale en Go` et n’a pas vocation à évoluer.

Elle représente un **état précis du projet**, dans lequel les différentes responsabilités commencent à être séparées.

L’objectif est de comprendre :

- pourquoi l’architecture minimale était suffisante au départ ;
    
- pourquoi le code a ensuite été séparé ;
    
- quel est le rôle de chaque package ;
    
- comment une requête HTTP traverse l’application.
    

---

## Point de départ : l’architecture minimale

Dans une architecture minimale, une grande partie de l’application peut être regroupée dans `main.go`.

```go
package main

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Accueil")
}

func ContactHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Contact")
}

func main() {
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/contact", ContactHandler)

	fmt.Println("Serveur lancé sur http://localhost:8080")

	http.ListenAndServe(":8080", nil)
}
```

Cette organisation est adaptée pour :

- découvrir le fonctionnement d’un serveur HTTP ;
    
- créer rapidement une première application ;
    
- comprendre les routes et les handlers ;
    
- tester une idée avec peu de code.
    

Tout est visible au même endroit.

---

## Limite de l’architecture minimale

Lorsque l’application grandit, `main.go` commence à cumuler plusieurs responsabilités.

Il peut devoir :

- démarrer le serveur ;
    
- créer le routeur ;
    
- déclarer les routes ;
    
- contenir les handlers ;
    
- initialiser la configuration ;
    
- préparer l’accès à la base de données.
    

Le fichier devient alors plus long et plus difficile à lire.

Le problème n’est pas que le programme ne fonctionne plus.

Le problème est que plusieurs rôles différents sont mélangés dans une même partie du code.

---

## Nouvelle architecture

Dans Club Manager, les premières responsabilités ont été séparées.

```text
club-manager/
├── cmd/
│   └── server/
│       └── main.go
│
├── internal/
│   ├── handlers/
│   │   └── ...
│   │
│   └── router/
│       └── router.go
│
├── go.mod
├── README.md
└── .gitignore
```

L’application est maintenant organisée autour de trois éléments principaux :

```text
main.go
    ↓
router
    ↓
handlers
```

---

## Responsabilité de `main.go`

Le fichier `main.go` est le point d’entrée de l’application.

```go
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
```

Son rôle est volontairement limité.

Il doit principalement :

1. préparer l’application ;
    
2. démarrer le serveur HTTP.
    

Il ne contient plus la déclaration détaillée des routes.

Il ne contient pas non plus le traitement des différentes pages.

```text
main.go
└── démarre l’application
```

---

## Responsabilité du package `router`

Le package `router` organise les routes de l’application.

```go
package router

import (
	"net/http"

	"github.com/grapinou/club-manager/internal/handlers"
)

func New() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/", handlers.HomeHandler)
	mux.HandleFunc("/club", handlers.ClubHandler)
	mux.HandleFunc("/contact", handlers.ContactHandler)
	mux.HandleFunc("/rules", handlers.RulesHandler)

	return mux
}
```

Son rôle est d’associer une URL à un handler.

```text
URL
↓
Handler correspondant
```

Par exemple :

```text
/          → HomeHandler
/club      → ClubHandler
/contact   → ContactHandler
/rules     → RulesHandler
```

Le routeur ne démarre pas le serveur.

Il ne produit pas directement le contenu des pages.

```text
router
└── associe les routes aux handlers
```

---

## Responsabilité du package `handlers`

Le package `handlers` contient les fonctions qui traitent les requêtes HTTP.

Exemple :

```go
package handlers

import (
	"fmt"
	"net/http"
)

func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Accueil")
}
```

Un handler reçoit :

```go
w http.ResponseWriter
r *http.Request
```

Il utilise la requête reçue et construit la réponse à envoyer au navigateur.

```text
handlers
└── traitent les requêtes et produisent les réponses
```

Les handlers ne démarrent pas le serveur.

Ils ne décident pas non plus eux-mêmes de l’URL à laquelle ils sont associés.

---

## Parcours d’une requête HTTP

Lorsqu’un utilisateur visite une page, la requête traverse plusieurs parties de l’application.

Exemple avec la page de contact :

```text
Navigateur
    ↓
GET /contact
    ↓
Serveur HTTP
    ↓
ServeMux
    ↓
handlers.ContactHandler
    ↓
Réponse HTTP
    ↓
Navigateur
```

Le serveur reçoit la requête.

Le `ServeMux` cherche la route correspondante :

```go
mux.HandleFunc("/contact", handlers.ContactHandler)
```

Il transmet ensuite la requête à :

```go
handlers.ContactHandler
```

Le handler construit enfin la réponse.

---

## Comparaison des deux architectures

### Architecture minimale

```text
main.go
├── démarre le serveur
├── déclare les routes
└── contient ou appelle les handlers
```

Tout est regroupé dans une même partie de l’application.

### Architecture avec séparation des responsabilités

```text
main.go
└── démarre le serveur

router
└── déclare les routes

handlers
└── traitent les requêtes
```

Chaque partie possède maintenant un rôle plus précis.

---

## Avant

```text
main.go
├── serveur
├── routes
└── handlers
```

## Après

```text
main.go
└── serveur

router
└── routes

handlers
└── traitement des requêtes
```

Le fonctionnement général reste similaire.

Ce n’est pas le serveur HTTP qui a changé.

C’est principalement l’organisation du code.

---

## Pourquoi séparer les responsabilités ?

### Améliorer la lisibilité

Chaque fichier possède un objectif identifiable.

Lorsque nous cherchons une route, nous savons qu’elle se trouve dans le package `router`.

Lorsque nous cherchons le traitement d’une page, nous savons qu’il se trouve dans le package `handlers`.

### Faciliter les modifications

Une modification dans un handler ne demande pas de modifier le démarrage du serveur.

L’ajout d’une route se fait dans le routeur sans surcharger `main.go`.

### Préparer la croissance du projet

Club Manager aura progressivement besoin de nouvelles fonctionnalités :

- templates HTML ;
    
- configuration ;
    
- accès à une base de données ;
    
- gestion des membres ;
    
- authentification ;
    
- rôles et autorisations.
    

La séparation actuelle prépare ces évolutions sans les ajouter trop tôt.

### Faciliter les tests

Une partie du programme ayant une responsabilité précise est généralement plus simple à tester indépendamment.

Le routeur peut être testé sans lancer réellement un serveur sur un port.

Les handlers peuvent également être testés séparément.

---

## Ce que cette architecture ne cherche pas à faire

Cette architecture n’a pas pour objectif de multiplier les dossiers ou les abstractions.

Elle ne cherche pas à reproduire immédiatement l’organisation d’une très grande application.

Le découpage est apparu parce que trois responsabilités concrètes ont été identifiées :

1. démarrer l’application ;
    
2. organiser les routes ;
    
3. traiter les requêtes.
    

Chaque nouveau package doit répondre à un besoin réel.

Créer davantage de fichiers ne rend pas automatiquement une architecture meilleure.

---

## Une séparation progressive

L’architecture n’a pas été entièrement décidée à l’avance.

Elle évolue en même temps que les besoins du projet.

```text
Architecture minimale
        ↓
Identification de plusieurs responsabilités
        ↓
Séparation du routeur
        ↓
Séparation des handlers
```

Cette progression permet d’éviter deux extrêmes.

### Tout conserver dans `main.go`

Cette solution reste simple au début, mais devient difficile à maintenir lorsque le projet grandit.

### Créer une architecture très complexe dès le départ

Cette solution ajoute des concepts et des abstractions avant de savoir s’ils sont réellement nécessaires.

L’objectif est donc d’ajouter de la structure uniquement lorsqu’elle devient utile.

---

## Principe de responsabilité unique

Cette organisation illustre le principe suivant :

> Une partie du programme devrait avoir une responsabilité principale clairement identifiable.

Dans l’architecture actuelle :

|Élément|Responsabilité principale|
|---|---|
|`main.go`|Démarrer l’application|
|`router`|Associer les routes aux handlers|
|`handlers`|Traiter les requêtes HTTP|

Cela ne signifie pas qu’un fichier ne doit contenir qu’une seule fonction.

Cela signifie que les fonctions regroupées doivent participer à un même rôle général.

---

## Dépendances entre les packages

Les dépendances suivent actuellement ce chemin :

```text
main
  ↓
router
  ↓
handlers
```

`main` utilise le package `router`.

```go
mux := router.New()
```

Le package `router` utilise le package `handlers`.

```go
mux.HandleFunc("/", handlers.HomeHandler)
```

Les handlers n’ont pas besoin de connaître `main`.

Ils n’ont pas besoin non plus de connaître la manière dont le serveur est lancé.

Cette organisation limite les connaissances nécessaires à chaque package.

---

## Pourquoi le dossier `internal` ?

Les packages `router` et `handlers` sont placés dans le dossier :

```text
internal/
```

En Go, un package situé dans `internal` est destiné à être utilisé uniquement à l’intérieur du projet.

```text
internal/
├── handlers/
└── router/
```

Cela indique que ces packages font partie du fonctionnement interne de Club Manager.

Ils ne sont pas conçus comme des bibliothèques destinées à être importées par d’autres projets.

---

## Pourquoi `cmd/server` ?

Le fichier principal se trouve dans :

```text
cmd/server/main.go
```

Le dossier `cmd` contient les points d’entrée exécutables du projet.

Le sous-dossier `server` indique que cet exécutable démarre le serveur web.

```text
cmd/
└── server/
    └── main.go
```

Cette organisation permettrait plus tard d’ajouter un autre exécutable sans mélanger les responsabilités.

Par exemple :

```text
cmd/
├── server/
│   └── main.go
└── migrate/
    └── main.go
```

Ce second exécutable n’est pas nécessaire actuellement.

Cet exemple montre seulement l’intérêt de distinguer le serveur web du reste du projet.

---

## Limites de cette étape

L’architecture actuelle reste volontairement simple.

Les handlers produisent encore directement des réponses basiques.

```go
fmt.Fprintln(w, "Accueil")
```

Il n’y a pas encore :

- de templates HTML ;
    
- de fichiers statiques ;
    
- de configuration ;
    
- de base de données ;
    
- de couche métier ;
    
- d’authentification ;
    
- de gestion des autorisations.
    

Ces éléments seront ajoutés seulement lorsqu’ils répondront à un besoin concret.

---

## Comprendre et retenir

L’architecture minimale permet de comprendre rapidement le fonctionnement général d’une application web en Go.

```text
route
↓
handler
↓
réponse
```

Lorsque l’application grandit, plusieurs responsabilités apparaissent.

```text
démarrer le serveur
organiser les routes
traiter les requêtes
```

Ces responsabilités peuvent alors être séparées :

```text
main.go  → démarrage
router   → routes
handlers → traitement
```

Le but n’est pas d’avoir le plus grand nombre de dossiers possible.

Le but est que chaque partie du programme ait un rôle compréhensible.

> Une architecture doit évoluer parce que le projet en a besoin, et non uniquement pour paraître plus professionnelle.

La séparation des responsabilités rend le code :

- plus lisible ;
    
- plus simple à modifier ;
    
- plus facile à tester ;
    
- mieux préparé à évoluer.
    

Cette fiche constitue un jalon figé de Club Manager.

Elle permet de comparer cette organisation avec l’architecture minimale présentée dans la fiche `00`.