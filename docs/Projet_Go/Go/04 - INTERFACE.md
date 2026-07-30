Tag : #interface 

---


# Les interfaces en Go

Les interfaces permettent de décrire **ce qu'un objet sait faire**, sans se préoccuper de son type exact.

Une interface répond à la question :

> "Quelles sont les méthodes disponibles ?"

et non :

> "Quel est le type de cet objet ?"

---

# Une analogie

Imaginons une télécommande.

Tu peux contrôler :

- une télévision ;
- un vidéoprojecteur ;
- une chaîne Hi-Fi.

La télécommande ne se demande pas :

> "Es-tu une télévision ?"

Elle demande simplement :

> "Sais-tu répondre aux boutons que je possède ?"

Si l'appareil comprend :

- Allumer
- Éteindre
- Volume +

alors il fonctionne avec cette télécommande.

En Go, une interface joue exactement ce rôle.

---

# Une première interface

```go
type Speaker interface {
    Speak()
}
```

Cette interface signifie :

> Tout objet possédant une méthode `Speak()` est un `Speaker`.

Elle ne contient aucune donnée.

Seulement des méthodes.

---

# Deux structures différentes

```go
type Dog struct{}

func (Dog) Speak() {
    fmt.Println("Wouf")
}
```

```go
type Cat struct{}

func (Cat) Speak() {
    fmt.Println("Miaou")
}
```

Aucune de ces structures ne déclare :

```go
implements Speaker
```

Pourtant…

Les deux sont des `Speaker`.

Pourquoi ?

Parce qu'elles possèdent la méthode demandée.

En Go, une interface est satisfaite automatiquement.

---

# Utilisation

```go
func SaySomething(s Speaker) {
    s.Speak()
}
```

On peut alors écrire :

```go
SaySomething(Dog{})
SaySomething(Cat{})
```

La fonction fonctionne avec les deux.

Elle ne connaît pas leur type.

Elle sait seulement qu'ils savent parler.

---

# Pourquoi utiliser une interface ?

Sans interface :

```go
func SaySomething(d Dog)
```

La fonction ne fonctionne qu'avec un chien.

Avec une interface :

```go
func SaySomething(s Speaker)
```

Elle fonctionne avec n'importe quel objet possédant `Speak()`.

Le code devient plus souple.

---

# L'interface `http.ResponseWriter`

Dans notre handler :

```go
func homeHandler(w http.ResponseWriter, r *http.Request)
```

`w` est une interface.

Pourquoi ?

Parce que Go n'a pas besoin de connaître le type exact qui écrit la réponse.

Il lui suffit de savoir que cet objet sait écrire des données vers le navigateur.

Le type concret est caché.

Notre code travaille uniquement avec les capacités offertes par l'interface.

---

# Les avantages

Les interfaces permettent :

- d'écrire du code plus générique ;
- de remplacer facilement une implémentation par une autre ;
- de faciliter les tests ;
- de réduire les dépendances entre les composants.

---

# Les interfaces célèbres de Go

Au fil du projet, nous rencontrerons notamment :

| Interface | Rôle |
|-----------|------|
| `error` | Représente une erreur |
| `io.Reader` | Lit des données |
| `io.Writer` | Écrit des données |
| `http.ResponseWriter` | Écrit une réponse HTTP |

Nous compléterons cette fiche lorsqu'elles apparaîtront.

---

# À retenir

Une interface ne décrit pas **ce qu'est** un objet.

Elle décrit **ce qu'il sait faire**.

C'est l'un des principes les plus importants de Go.


Voir : 

- [[04.02 - Interface utilisation]]
- [[04.01 - Interface composée]]