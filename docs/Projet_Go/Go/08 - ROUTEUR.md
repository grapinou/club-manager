
Tag : #go #http #routeur #architecture

---

- [[08.01 - Routeur ServeMux]]
# Le routeur

## Définition

Un routeur est un composant chargé d'associer une URL à une fonction.

Lorsqu'une requête HTTP arrive, le routeur détermine quel Handler doit être appelé.

Par exemple :

```
GET /
        │
        ▼
HomeHandler
```

```
GET /contact
        │
        ▼
ContactHandler
```

Le routeur ne traite pas lui-même la requête.

Il choisit simplement le bon Handler.

---

# Le rôle du routeur

Le routeur répond à une question :

> Quelle fonction doit traiter cette URL ?

Il agit comme un répartiteur.

```
                Requête HTTP
                      │
                      ▼
                 Routeur
          ┌───────────┼───────────┐
          ▼           ▼           ▼
     HomeHandler ContactHandler RulesHandler
```

---

# Les routes

Une route associe :

- une URL ;
- un Handler.

Exemple :

```go
mux.HandleFunc("/", handlers.HomeHandler)

mux.HandleFunc("/contact", handlers.ContactHandler)
```

---

# Pourquoi utiliser un routeur ?

Sans routeur, il faudrait analyser manuellement chaque URL.

Le routeur automatise cette tâche.

Il permet également :

- d'organiser les routes ;
- de centraliser leur configuration ;
- de rendre l'application plus facile à maintenir.

---

# Dans Club Manager

Le package `router` est responsable de la configuration des routes.

Le package `main` se contente de demander un routeur prêt à être utilisé.

Exemple :

```go
router := router.New()

http.ListenAndServe(":8080", router)
```

Ainsi, `main.go` ne connaît pas les détails des routes.

---

# À retenir

- Le routeur associe une URL à un Handler.
- Il choisit quelle fonction sera appelée.
- Il ne traite pas lui-même les requêtes.
- Il centralise la configuration des routes.