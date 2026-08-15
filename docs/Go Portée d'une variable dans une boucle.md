
---

# Go — Portée d'une variable dans une boucle

## Idée générale

L'endroit où une variable est déclarée détermine notamment **où elle est accessible** et à quel moment sa valeur est initialisée.

Dans une boucle, cela peut avoir une conséquence importante :

```go
email := "Aucun mail"

for _, member := range members {
	// ...
}
```

n'a pas le même comportement que :

```go
for _, member := range members {
	email := "Aucun mail"

	// ...
}
```

---

## Exemple rencontré dans Club Manager

On souhaite transformer un email provenant de PostgreSQL en une chaîne propre pour la vue.

Un membre peut avoir :

```go
member.Email.Valid == true
```

si un email existe, ou :

```go
member.Email.Valid == false
```

si PostgreSQL contient `NULL`.

On souhaite obtenir :

```text
email présent → adresse email
email absent  → "Aucun mail"
```

---

## Mauvaise version

```go
email := "Aucun mail"

for _, member := range members {

	if member.Email.Valid {
		email = member.Email.String
	}

	memberData = append(memberData, views.MemberData{
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     email,
	})
}
```

### Le problème

`email` est créée **avant la boucle**.

Elle conserve donc sa valeur d'une itération à l'autre.

Imaginons :

```text
Robin → robin@example.com
Jean  → NULL
```

### Premier passage

Au début :

```go
email == "Aucun mail"
```

Robin possède un email :

```go
if member.Email.Valid {
	email = member.Email.String
}
```

La variable devient :

```go
email == "robin@example.com"
```

### Deuxième passage

Jean n'a pas d'email :

```go
member.Email.Valid == false
```

Le `if` n'est donc pas exécuté.

Mais `email` contient toujours :

```go
"robin@example.com"
```

Jean reçoit alors accidentellement l'email de Robin.

---

## Bonne version

```go
for _, member := range members {

	email := "Aucun mail"

	if member.Email.Valid {
		email = member.Email.String
	}

	memberData = append(memberData, views.MemberData{
		FirstName: member.FirstName,
		LastName:  member.LastName,
		Email:     email,
	})
}
```

Ici :

```go
email := "Aucun mail"
```

est exécuté à chaque passage dans la boucle.

On repart donc systématiquement de la valeur par défaut.

---

## Déroulement

Pour Robin :

```text
nouvelle itération
↓
email = "Aucun mail"
↓
Email.Valid == true
↓
email = "robin@example.com"
```

Pour Jean :

```text
nouvelle itération
↓
email = "Aucun mail"
↓
Email.Valid == false
↓
email reste "Aucun mail"
```

Résultat :

```text
Robin → robin@example.com
Jean  → Aucun mail
```

---

## Portée de la variable

### Variable déclarée avant la boucle

```go
email := "Aucun mail"

for _, member := range members {
	// email est accessible ici
}
```

`email` appartient au bloc qui contient la boucle.

Elle continue donc d'exister après chaque itération et conserve la valeur qui lui a été affectée.

Schéma :

```text
email
  │
  ├── tour 1
  ├── tour 2
  ├── tour 3
  └── ...
```

---

### Variable déclarée dans la boucle

```go
for _, member := range members {

	email := "Aucun mail"

	// utilisation de email
}
```

La déclaration appartient au bloc de la boucle.

À chaque exécution du corps, on repart de :

```go
email := "Aucun mail"
```

Schéma :

```text
tour 1
└── email = "Aucun mail"

tour 2
└── email = "Aucun mail"

tour 3
└── email = "Aucun mail"
```

---

## `:=` et `=`

Cette situation permet également de revoir une distinction importante.

### `:=`

```go
email := "Aucun mail"
```

déclare une variable et lui donne une première valeur.

### `=`

```go
email = member.Email.String
```

modifie une variable qui existe déjà.

Dans notre code :

```go
for _, member := range members {

	email := "Aucun mail"

	if member.Email.Valid {
		email = member.Email.String
	}
}
```

on peut lire :

```text
:= → créer email avec une valeur par défaut

=  → remplacer cette valeur si un email existe
```

---

## Valeur par défaut puis exception

Le code suit un modèle très utile :

```go
email := "Aucun mail"

if member.Email.Valid {
	email = member.Email.String
}
```

On commence par définir le cas général :

```text
pas d'email
```

puis on traite l'exception :

```text
un email existe
```

Cela évite par exemple :

```go
var email string

if member.Email.Valid {
	email = member.Email.String
} else {
	email = "Aucun mail"
}
```

Les deux fonctionnent, mais la première version exprime simplement :

```text
par défaut → Aucun mail
sinon      → utiliser le vrai mail
```

---

## Lien avec la séparation des responsabilités

Dans Club Manager, PostgreSQL conserve l'état réel :

```text
email existant → valeur TEXT
email absent   → NULL
```

`sqlc` représente cela avec :

```go
pgtype.Text
```

et notamment :

```go
Valid
```

Le handler transforme ensuite cette donnée technique en donnée d'affichage :

```text
NULL
↓
"Aucun mail"
```

La vue reçoit donc simplement :

```go
Email string
```

Elle n'a pas besoin de connaître :

```go
pgtype.Text
```

ni :

```go
Valid
```

On obtient :

```text
PostgreSQL
    ↓
NULL
    ↓
dbsqlc.Member
    ↓
pgtype.Text{Valid: false}
    ↓
handler
    ↓
"Aucun mail"
    ↓
views.MemberData
    ↓
template
```

---

## Comprendre et retenir

> Une variable déclarée avant une boucle peut conserver une valeur modifiée entre les différentes itérations.

Si chaque élément doit repartir d'une valeur initiale indépendante, déclarer ou réinitialiser cette valeur **dans la boucle**.

Dans notre cas :

```go
for _, member := range members {
	email := "Aucun mail"
}
```

signifie :

> Pour chaque membre, recommencer avec `"Aucun mail"`.

---

## Exemple minimal

### Problématique

```go
value := "défaut"

for _, n := range []int{1, 2, 3} {
	if n == 1 {
		value = "modifié"
	}

	fmt.Println(value)
}
```

Résultat :

```text
modifié
modifié
modifié
```

La modification survit aux itérations suivantes.

### Correction

```go
for _, n := range []int{1, 2, 3} {

	value := "défaut"

	if n == 1 {
		value = "modifié"
	}

	fmt.Println(value)
}
```

Résultat :

```text
modifié
défaut
défaut
```

Chaque itération repart de sa propre valeur initiale.

---

## Exercice

Que va afficher ce code ?

```go
status := "inconnu"

users := []string{"Alice", "", "Bob"}

for _, user := range users {

	if user != "" {
		status = user
	}

	fmt.Println(status)
}
```

Puis modifier le code pour obtenir :

```text
Alice
inconnu
Bob
```

### Indice

Réfléchir à l'endroit où :

```go
status := "inconnu"
```

doit être exécuté.

