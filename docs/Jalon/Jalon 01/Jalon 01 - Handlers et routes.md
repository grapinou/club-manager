
---

## project: Club Manager  
type: milestone  
milestone: 01  
git-tag: milestone-01-handlers-routes  
git-commit: 72e5e90  
date: 2026-08-02  
status: frozen

# Jalon 01 - Handlers et routes

> [!info] État figé  
> Cette fiche représente Club Manager au tag Git :
> 
> `milestone-01-handlers-routes`
> 
> Elle ne doit pas évoluer avec le projet.

---

## Objectif atteint

Club Manager possède maintenant une première véritable architecture HTTP.

L'application n'est plus entièrement définie dans `main.go`.

Les responsabilités commencent à être séparées entre :

- le point d'entrée ;
    
- le routeur ;
    
- les handlers.
    

---

# Architecture

```text
Navigateur
    │
    │ requête HTTP
    ↓
  Router
    │
    ↓
 Handler
    │
    ↓
Réponse texte
```

Le projet possède notamment :

```text
cmd/server/
└── main.go

internal/
├── handlers/
└── router/
```

---

# Rôle de `main`

`main` crée le routeur :

```go
mux := router.New()
```

puis démarre le serveur :

```go
http.ListenAndServe(":8080", mux)
```

Il ne contient donc plus directement la définition des routes.

---

# Le routeur

Le routeur repose sur :

```go
http.NewServeMux()
```

Les quatre pages disponibles sont :

```go
mux.HandleFunc("GET /{$}", handlers.HomeHandler)
mux.HandleFunc("GET /club", handlers.ClubHandler)
mux.HandleFunc("GET /contact", handlers.ContactHandler)
mux.HandleFunc("GET /rules", handlers.RulesHandler)
```

Une route relie donc :

```text
méthode HTTP + URL
        ↓
     handler
```

Exemple :

```text
GET /club
    ↓
ClubHandler
```

Les routes sont explicitement limitées à la méthode HTTP `GET`.

---

# Les handlers

Chaque page possède son propre handler.

À ce stade, un handler produit encore directement la réponse HTTP.

Exemple :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
    fmt.Fprintln(w, "Bienvenue sur Club Manager")
}
```

Le handler possède donc deux responsabilités :

```text
traiter la requête
       +
produire le contenu
```

Cette architecture est suffisante pour un site minimal, mais cette double responsabilité deviendra bientôt une limite.

---

# Tests HTTP

Le routeur commence également à être testé.

Par exemple, une requête :

```text
POST /contact
```

doit être refusée puisque la route attend :

```text
GET /contact
```

Le test vérifie alors :

```go
http.StatusMethodNotAllowed
```

soit le statut HTTP :

```text
405 Method Not Allowed
```

Les tests commencent ainsi à vérifier le comportement HTTP de l'application et pas seulement l'exécution des fonctions.

---

# Ce que ce jalon apporte

Au départ, on pourrait imaginer :

```text
main.go
│
├── serveur
├── routes
└── contenu
```

Le projet possède maintenant :

```text
main
  │
  ↓
router
  │
  ↓
handlers
```

Une première séparation des responsabilités apparaît.

---

# Limite observée

Les handlers produisent encore directement le contenu :

```go
fmt.Fprintln(...)
```

Cela convient pour quelques lignes de texte.

Cela devient cependant peu adapté lorsque l'on souhaite produire :

```html
<!DOCTYPE html>
<html>
...
</html>
```

La prochaine étape naturelle sera donc de séparer :

```text
traitement HTTP
```

de :

```text
présentation HTML
```

---

# Comprendre et retenir

> **Le routeur choisit le handler à partir de la requête HTTP.**

```text
GET /club
    ↓
Router
    ↓
ClubHandler
```

---

> **Un handler est une fonction capable de traiter une requête HTTP et d'écrire une réponse.**

---

> **À ce stade, les handlers gèrent encore eux-mêmes la présentation.**

Cette limite conduira naturellement au jalon suivant.

---

# État du jalon

```text
Serveur HTTP           ✅
ServeMux               ✅
Routes séparées        ✅
Handlers séparés       ✅
Routes limitées à GET  ✅
Tests HTTP             ✅

Views                   ❌
Templates HTML          ❌
Layout commun           ❌
Configuration           ❌
Base de données         ❌
```

**Ce jalon marque la séparation entre le démarrage de l'application, le routage HTTP et le traitement des requêtes.**