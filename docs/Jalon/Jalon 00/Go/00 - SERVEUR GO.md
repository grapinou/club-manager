Tag : #pointeur ; #serveur ; #interface

---


# Les briques d'un serveur web Go

Cette fiche regroupe les fonctions et concepts rencontrés lors du développement de Club Manager.

---

# package main

```go
package main
```

## Rôle

Déclare le paquet principal de l'application.

Un programme Go exécutable doit obligatoirement appartenir au package `main`.

---

# func main()

```go
func main() {

}
```

## Rôle

Point d'entrée du programme.

Lorsque l'on lance :

```bash
go run ./cmd/server
```

Go commence toujours par exécuter `main()`.

---

# import

```go
import (
    "fmt"
    "net/http"
)
```

## Rôle

Importer les bibliothèques nécessaires.

### fmt

Permet principalement :

- afficher du texte
- formater des chaînes

### net/http

Bibliothèque standard permettant de créer un serveur web.

---

# http.HandleFunc()

```go
http.HandleFunc("/", homeHandler)
```

## Rôle

Associe une URL à une fonction.

Ici :

```
/
```

sera traitée par :

```go
homeHandler
```

On appelle cela une **route**.

---

# func homeHandler()

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {

}
```

Remarque : 

- voir la note interface pour http.ResponseWriter

- voir la note pointeur pour \*http.Request
## Rôle

Fonction appelée lorsqu'un navigateur visite une route.

### Paramètres

#### http.ResponseWriter

Permet d'écrire la réponse envoyée au navigateur.

Exemple :

```go
fmt.Fprintln(w, "Bonjour")
```

#### \*http.Request

Contient les informations de la requête :

- URL
- méthode (GET, POST...)
- paramètres
- cookies
- etc.

---

# fmt.Fprintln()

```go
fmt.Fprintln(w, "Bonjour")
```

## Rôle

Écrit une ligne dans la réponse HTTP.

Le navigateur recevra :

```
Bonjour
```

---

# http.ListenAndServe()

```go
http.ListenAndServe(":8080", nil)
```

## Rôle

Démarre le serveur web.

### Paramètres

```
:8080
```

Port utilisé.

```
nil
```

Utilise le routeur par défaut de Go.

---

# nil

`nil` est l'équivalent de `null` dans d'autres langages.

Il représente l'absence de valeur pour certains types (pointeurs, interfaces, slices, maps, fonctions, canaux...).

Dans notre exemple :

```go
http.ListenAndServe(":8080", nil)
```

Go utilise automatiquement le routeur HTTP par défaut.