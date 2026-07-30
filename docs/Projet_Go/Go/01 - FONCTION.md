
Tag : #go #fonction

---

- [[01.01 - Fonction exportation]]
- [[01.02 - Fonction New]]



---

# Les fonctions en Go

## Définition

Une fonction est un bloc de code réalisant une tâche précise.

Elle permet d'éviter de répéter le même code et de découper un programme en petites unités faciles à comprendre.

---

# Déclarer une fonction

La syntaxe générale est :

```go
func NomDeLaFonction(paramètres) typeDeRetour {

}
```

Exemple :

```go
func Bonjour() {
	fmt.Println("Bonjour")
}
```

---

# Appeler une fonction

Une fonction est exécutée lorsqu'on l'appelle.

Exemple :

```go
Bonjour()
```

---

# Les paramètres

Une fonction peut recevoir des informations appelées **paramètres**.

Exemple :

```go
func DireBonjour(prenom string) {
	fmt.Println("Bonjour", prenom)
}
```

Appel :

```go
DireBonjour("Rémi")
```

Résultat :

```
Bonjour Rémi
```

---

# Les valeurs de retour

Une fonction peut renvoyer une valeur.

Exemple :

```go
func Addition(a int, b int) int {
	return a + b
}
```

Utilisation :

```go
resultat := Addition(2, 3)
```

`resultat` vaut :

```
5
```

Une fonction peut également retourner plusieurs valeurs.

Exemple :

```go
func Division(a, b float64) (float64, error) {

}
```

---

# Une fonction peut ne rien retourner

Exemple :

```go
func AfficherMessage() {
	fmt.Println("Bienvenue")
}
```

Le type de retour est alors omis.

---

# Une fonction est un type

En Go, une fonction peut être :

- appelée ;
- stockée dans une variable ;
- passée en argument à une autre fonction ;
- retournée par une fonction.

Exemple :

```go
http.HandleFunc("/", handlers.HomeHandler)
```

Ici, `HomeHandler` est passée en argument à `HandleFunc`.

Go l'appellera plus tard lorsqu'une requête arrivera.

Ce mécanisme est appelé un **callback**.

---

# Les fonctions exportées

Une fonction dont le nom commence par une majuscule est **exportée**.

Exemple :

```go
func HomeHandler() {

}
```

Elle est accessible depuis les autres packages.

À l'inverse :

```go
func homeHandler() {

}
```

n'est accessible que dans son propre package.

---

# Les méthodes

Une méthode est une fonction associée à un type.

Exemple :

```go
type Club struct {
	Name string
}

func (c Club) Display() {
	fmt.Println(c.Name)
}
```

Une méthode appartient à un type.

Une fonction n'appartient à aucun type.

---

# Les fonctions dans Club Manager

Nous utilisons déjà plusieurs fonctions :

```go
func main()
```

Point d'entrée du programme.

---

```go
func HomeHandler(w http.ResponseWriter, r *http.Request)
```

Répond à la route `/`.

---

```go
func ClubHandler(w http.ResponseWriter, r *http.Request)
```

Répond à la route `/club`.

---

```go
func New() *http.ServeMux
```

Construit et retourne un nouveau routeur.

---

# Bonnes pratiques

Une fonction devrait :

- avoir une responsabilité unique ;
- être courte et facile à lire ;
- avoir un nom décrivant son rôle ;
- recevoir uniquement les paramètres dont elle a besoin.

---

# À retenir

- Une fonction réalise une tâche précise.
- Elle peut recevoir des paramètres.
- Elle peut retourner une ou plusieurs valeurs.
- Elle peut être passée en argument à une autre fonction.
- Une majuscule rend une fonction exportée.
- Une méthode est une fonction associée à un type.