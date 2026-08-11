
---

# Struct et méthodes

## Pourquoi une struct ?

Une **struct** est simplement un moyen de **regrouper plusieurs informations qui décrivent une même chose**.

Au lieu d'avoir plusieurs variables séparées :

```go
name := "Rex"
age := 5
breed := "Berger"
````

on les rassemble :

```go
type Dog struct {
    Name  string
    Age   int
    Breed string
}
```

On obtient alors un seul objet :

```go
dog := Dog{
    Name: "Rex",
    Age: 5,
    Breed: "Berger",
}
```

---

## Une struct représente un objet

Une struct décrit un objet du monde réel ou du programme.

Exemples :

- Un utilisateur
- Un fichier
- Une fenêtre
- Un processus
- Un paquet réseau
- Un chien

Elle contient uniquement des **données**.

---

## Pourquoi ajouter des méthodes ?

Les données seules ne suffisent pas.

Un chien peut :

- aboyer
- courir
- manger

Un fichier peut :

- s'ouvrir
- se fermer
- être lu

Un processus peut :

- démarrer
- s'arrêter
- être suspendu

Les **méthodes** représentent les actions que l'objet sait réaliser.

---

## Exemple

```go
type Dog struct {
    Name string
}
```

On ajoute une méthode :

```go
func (d Dog) Speak() {
    fmt.Println("Wouf")
}
```

Utilisation :

```go
dog := Dog{Name: "Rex"}

dog.Speak()
```

Affiche :

```go
Wouf
```

---

## Comment lire une méthode ?

```go
func (d Dog) Speak() {
```

se lit :

> La struct **Dog** possède une méthode appelée **Speak**.

Le `(d Dog)` indique sur quel type cette méthode est définie.

`d` représente simplement l'objet courant.

On pourrait écrire :

```go
func (chien Dog) Speak()
```

ou

```go
func (x Dog) Speak()
```

mais on utilise généralement une lettre courte.

---

## Une struct vide

Une struct peut être vide :

```go
type Dog struct{}
```

Cela signifie qu'elle ne contient aucune donnée.

Pourtant elle peut avoir des méthodes :

```go
func (d Dog) Speak() {
    fmt.Println("Wouf")
}
```

C'est parfaitement valide.

---

## Les méthodes utilisent souvent les données

Le plus souvent, une méthode utilise les champs de la struct.

```go
type Dog struct {
    Name string
}
```

```go
func (d Dog) Speak() {
    fmt.Println(d.Name + " : Wouf")
}
```

Résultat :

```go
Rex : Wouf
```

---

## Une bonne analogie

Une struct est comme une fiche d'identité.

```go
Dog

Nom : Rex
Age : 5
Race : Berger
```

Une méthode correspond aux actions que cette fiche permet de réaliser.

```go
Dog

Nom : Rex

Actions :

- Speak()
- Run()
- Eat()
```

Les données et les comportements restent regroupés.

---

## Lien avec les interfaces

Une interface décrit des comportements.

```go
type Speaker interface {
    Speak()
}
```

Une struct devient automatiquement un `Speaker` si elle possède cette méthode.

```go
type Dog struct{}

func (d Dog) Speak() {}
```

Aucun mot-clé n'est nécessaire.

Dog satisfait automatiquement l'interface `Speaker`.

---

## À retenir

- Une struct regroupe des données.
- Une méthode représente une action.
- Une méthode appartient à un type.
- Une struct peut être vide.
- Une struct satisfait automatiquement une interface lorsqu'elle possède les méthodes demandées.

````

---

### Pourquoi cette notion est importante pour un système d'exploitation

Quand tu développeras un OS, tu retrouveras ce schéma partout :

```text
Process
├── pid
├── state
├── priority
└── ...
````

avec des opérations comme :

```
Start()
Stop()
Sleep()
WakeUp()
```

ou encore :

```
File
├── name
├── size
├── permissions
└── ...

Open()
Read()
Write()
Close()
```