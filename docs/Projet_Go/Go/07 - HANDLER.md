

Tag : #go #callback #http #handler #serveur

---

# Comment un Handler est-il appelé ?

Code de base : 

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

Lorsque l'on écrit :

```go
http.HandleFunc("/", handlers.HomeHandler)
```

on pourrait croire que `HomeHandler` est exécutée immédiatement.

Ce n'est pas le cas.

`HandleFunc` enregistre simplement la fonction qui devra être appelée lorsqu'une requête arrivera sur cette route.

---

# Qui appelle le Handler ?

Ce n'est jamais notre code.

C'est le serveur HTTP de Go.

Lorsque l'on démarre le serveur :

```go
http.ListenAndServe(":8080", nil)
```

Go :

- écoute le port 8080 ;
- attend qu'un client se connecte ;
- reçoit une requête HTTP ;
- recherche le Handler associé à la route ;
- appelle automatiquement notre fonction.

Le schéma est le suivant :

```
Navigateur
     │
     ▼
ListenAndServe
     │
     ▼
Recherche du Handler
     │
     ▼
HomeHandler(w, r)
```

---

# D'où vient la Request ?

Lorsque le navigateur envoie une requête :

```
GET /contact HTTP/1.1
Host: localhost:8080
...
```

Go construit une structure :

```go
http.Request
```

Cette structure contient toutes les informations de la requête :

- l'URL ;
- la méthode (GET, POST...) ;
- les en-têtes HTTP ;
- les paramètres ;
- les cookies ;
- etc.

Go passe ensuite un pointeur vers cette structure au Handler :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request)
```

Le Handler n'a rien à créer.

Tout est préparé par Go.

---

# D'où vient ResponseWriter ?

Avant d'appeler le Handler, Go prépare également un objet permettant d'envoyer une réponse au navigateur.

Cet objet implémente l'interface :

```go
http.ResponseWriter
```

Il est passé au Handler dans le paramètre :

```go
w
```

Lorsque l'on écrit :

```go
fmt.Fprintln(w, "Bonjour")
```

le texte est envoyé au navigateur.

Le Handler n'a pas besoin de savoir comment cette réponse est transmise.

Il utilise simplement l'interface `ResponseWriter`.

Voir :
- [[06.01 - fmt.Fprintln]]
	- Permet d'envoyer sur w qui est un io.writer du texte
- Le serveur écrit la réponse dans `w`, puis cette réponse est envoyée au navigateur, qui l'affiche.

---

# Pourquoi ResponseWriter est-il une interface ?

Une interface décrit ce qu'un objet sait faire.

`http.ResponseWriter` indique qu'un objet est capable :

- d'écrire une réponse HTTP ;
- de définir des en-têtes HTTP ;
- de choisir un code de statut (200, 404, 500...).

Le Handler travaille avec l'interface et ignore l'implémentation réelle.

Cette abstraction permet au même Handler de fonctionner avec différents types de serveurs ou d'outils de test.

---

# Le rôle du Handler

Le Handler a une responsabilité très simple :

- lire les informations de la requête (`r`) ;
- construire la réponse en écrivant dans (`w`).

Exemple :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Club Manager")
}
```

---

# Une analogie

Imagine un restaurant.

Le serveur :

- accueille le client ;
- prend la commande ;
- apporte une assiette vide au cuisinier.

Le cuisinier :

- ne crée ni le client ni l'assiette ;
- prépare simplement le plat.

Le Handler joue le rôle du cuisinier.

Go lui fournit :

- la requête (`r`) ;
- un moyen de répondre (`w`).

Le Handler se contente de traiter la demande.

---

# Le principe du callback

Le Handler est un callback.

Un callback est une fonction que l'on fournit à une autre partie du programme afin qu'elle soit appelée automatiquement lorsqu'un événement se produit.

Ici :

- on enregistre le Handler avec `http.HandleFunc()`;
- Go le mémorise ;
- lorsque le navigateur effectue une requête, Go appelle automatiquement cette fonction.

Ce n'est donc jamais notre code qui appelle directement :

```go
HomeHandler(w, r)
```

C'est le serveur HTTP qui le fait.

---

# Résumé

Lorsque tu écris :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Club Manager")
}
```

on peut le lire presque comme une phrase :

1. Le navigateur envoie une requête.
2. Go construit `r`.
3. Go prépare `w`.
4. Ton handler lit éventuellement `r`.
5. Ton handler écrit dans `w`.
6. Go envoie le contenu de `w` au navigateur.
7. Le navigateur interprète la réponse et l'affiche.

---

# À retenir

- `HandleFunc()` enregistre un Handler.
- `ListenAndServe()` attend les requêtes des clients.
- Go construit automatiquement `http.Request`.
- Go fournit un `http.ResponseWriter`.
- Le Handler lit `r` et écrit dans `w`.
- Le Handler est un callback : il est appelé automatiquement par le serveur HTTP.