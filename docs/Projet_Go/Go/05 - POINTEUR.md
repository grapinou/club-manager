Tag : #pointeur 

---

# Les pointeurs en Go

Les pointeurs permettent de manipuler une variable **sans en faire une copie**.

Ils contiennent l'adresse mémoire d'une valeur.

---

# Pourquoi utiliser un pointeur ?

Imaginons :

```go
type Member struct {
    Name string
}
```

Puis :

```go
m := Member{Name: "Rémi"}
```

Si une fonction reçoit `m` :

```go
func hello(m Member) {

}
```

Go crée une **copie** de `m`.

Le contenu est identique, mais ce n'est plus le même objet en mémoire.

---

# Recevoir un pointeur

```go
func hello(m *Member) {

}
```

Ici, `m` n'est plus une copie.

La fonction reçoit l'adresse mémoire de `Member`.

Toutes les modifications effectuées dans la fonction modifient l'objet original.

---

# Le symbole *

Dans une déclaration :

```go
var p *Member
```

`*` signifie :

> « p est un pointeur vers un Member. »

---

# Le symbole &

Pour obtenir l'adresse d'une variable :

```go
m := Member{Name: "Rémi"}

p := &m
```

`&` signifie :

> « Donne-moi l'adresse mémoire de cette variable. »

---

# Accéder à la valeur

Si :

```go
p := &m
```

alors :

```go
(*p).Name
```

accède au champ `Name`.

En Go, on peut écrire plus simplement :

```go
p.Name
```

Go effectue automatiquement le déréférencement.

---

# Pourquoi *http.Request ?

Notre handler est écrit ainsi :

```go
func homeHandler(w http.ResponseWriter, r *http.Request) {

}
```

`r` est un pointeur vers la requête HTTP.

Pourquoi ?

Parce qu'une requête contient beaucoup d'informations :

- URL
- méthode HTTP
- en-têtes
- cookies
- paramètres
- contexte
- etc.

La copier à chaque appel serait inutilement coûteuse.

Go transmet donc un pointeur vers la requête.

---

# À retenir

| Symbole | Signification |
|----------|---------------|
| `*T` | pointeur vers un `T` |
| `&x` | adresse mémoire de `x` |
| `*p` | valeur pointée par `p` |

---

# Ce que nous verrons plus tard

Nous rencontrerons souvent les pointeurs :

- `*http.Request`
- `*sql.DB`
- `*template.Template`
- `*os.File`
- `context.Context`
- nos propres structures (`*Member`, `*Association`, etc.)

À chaque fois, nous compléterons cette fiche avec des exemples concrets.