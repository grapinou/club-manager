
Tag : #go #package #architecture

---

# Les packages en Go

## Définition

Un package est un ensemble de fichiers Go ayant la même responsabilité.

Exemple :

```
club-manager/
│
├── cmd/
│   └── server/
│       └── main.go
│
└── internal/
    └── handlers/
        ├── home.go
        ├── club.go
        ├── contact.go
        └── rules.go

```

Ici, il existe deux packages :

- main
- handlers

---

# Déclaration d'un package

Chaque fichier Go commence par la déclaration de son package.

Exemple :

```go
package handlers
```

Tous les fichiers d'un même dossier appartiennent généralement au même package.

Par exemple :

home.go

```go
package handlers
```

club.go

```go
package handlers
```

contact.go

```go
package handlers
```

Ces fichiers sont compilés ensemble.

---

# Pourquoi utiliser des packages ?

Les packages permettent de séparer les responsabilités.

Exemple :

- main → démarre l'application.
- handlers → répond aux requêtes HTTP.
- models → contiendra les structures métier.
- router → configurera les routes.
- repository → accédera à la base de données.

Chaque package a un rôle bien défini.

---

# Pourquoi plusieurs fichiers ?

Un package peut contenir plusieurs fichiers.

Par exemple :
```
handlers/
	├── home.go
	├── club.go
	├── contact.go
	└── rules.go
```

Tous ces fichiers appartiennent au package `handlers`.

Ils peuvent utiliser les fonctions exportées du même package sans avoir besoin d'import.

---

# Une responsabilité par package

Une bonne pratique consiste à regrouper les éléments ayant le même rôle.

On n'organise pas un projet parce qu'il est grand.

On l'organise pour qu'il puisse grandir sans devenir difficile à comprendre.

---

# À retenir

- Un package regroupe des fichiers ayant la même responsabilité.
- Un package peut contenir plusieurs fichiers.
- Tous les fichiers d'un package sont compilés ensemble.
- On importe toujours un package, jamais un fichier.


# [[02.01 - Package import]]