MOC : [[00 - INITIALISATION PROJET]]
Tag : #go #gomod

---


# Qu'est-ce qu'un module Go ?

Avant Go 1.11, les projets Go étaient organisés dans un espace de travail appelé **GOPATH**. Cela fonctionnait, mais présentait plusieurs limites : difficile de gérer plusieurs versions d'une même bibliothèque, dépendance à une structure de dossiers particulière, etc.

Les **Go Modules** ont remplacé ce système.

Un module est simplement un projet Go identifié par un nom unique.

Dans notre cas, nous allons utiliser le nom du dépôt GitHub :

```
github.com/grapinou/club-manager
```

Pourquoi ce nom ?

Parce que si, un jour, quelqu'un veut utiliser une partie de ton projet, il pourra écrire :

```
import "github.com/grapinou/club-manager/..."
```

C'est devenu la convention dans l'écosystème Go.

---

# La commande

Depuis le dossier `club-manager`, exécute :

```
go mod init github.com/grapinou/club-manager
```

---

## Que va-t-il se passer ?

Go va créer un fichier :

```
go.mod
```

Son contenu devrait ressembler à ceci :

```
module github.com/grapinou/club-manager

go 1.26.5
```

_(Le numéro de version dépendra de la version de Go que tu as installée.)_

Ce fichier indique deux choses essentielles :

- le **nom du module** ;
- la **version minimale de Go** utilisée par le projet.

---

# Pourquoi `go.mod` est-il important ?

Il joue un rôle comparable à celui de :

- `package.json` en JavaScript,
- `Cargo.toml` en Rust,
- `pom.xml` en Java,
- `pyproject.toml` en Python.

Il décrit ton projet et ses dépendances.

Par exemple, plus tard, lorsque nous ajouterons PostgreSQL, tu verras apparaître des dépendances dans `go.mod`.

---

# Et `go.sum` ?

Tu remarqueras peut-être qu'il **n'apparaît pas tout de suite**.

C'est normal.

`go.sum` est créé lorsque Go télécharge sa première dépendance externe.

Il sert à garantir l'intégrité des bibliothèques téléchargées (en enregistrant leurs sommes de contrôle).

---

# Vérification

Une fois la commande exécutée, fais :

```
ls
```

Tu devrais voir :

```
README.md
go.mod
```

Puis :

```
cat go.mod
```

Envoie-moi son contenu.

---

## Pour ta fiche Obsidian

Voici un premier résumé que nous enrichirons ensuite.

### `go mod init`

**Objectif**

Initialiser un projet Go en créant un module.

**Syntaxe**

```
go mod init <nom_du_module>
```

**Exemple**

```
go mod init github.com/grapinou/club-manager
```

**Fichier créé**

```
go.mod
```

**Contenu**

- nom du module ;
- version de Go ;
- dépendances (au fur et à mesure du projet).

**Bonnes pratiques**

- Utiliser le chemin GitHub comme nom du module si le projet est hébergé sur GitHub.
- Versionner `go.mod` dans Git.
- Versionner également `go.sum` lorsqu'il sera créé.

---

Une fois `go.mod` créé, nous ferons immédiatement un **commit dédié** (`Initialize Go module`). Ce sera notre dernier commit d'initialisation avant de commencer à concevoir l'architecture de **Club Manager**. 🚀


## Vérifions le contenu

Peux-tu exécuter :

```
cat go.mod
```

Tu devrais obtenir quelque chose de très proche de :

```
module github.com/grapinou/club-manager

go 1.26.5
```

(Le numéro de version peut varier si tu mets Go à jour.)

---

# Comprendre ce que contient `go.mod`

Le fichier est très simple au départ.

```
module github.com/grapinou/club-manager
```

Cette ligne définit **l'identité du projet**.

À partir de maintenant, lorsqu'un fichier Go importera un package de ton projet, il utilisera ce chemin.

Par exemple, si plus tard tu crées :

```
internal/database
```

tu pourras écrire :

```
import "github.com/grapinou/club-manager/internal/database"
```

Go saura immédiatement que ce package appartient à **ton** projet.

---

Ensuite :

```
go 1.26.5
```

Cette ligne indique :

> « Ce module a été créé pour Go 1.26.5. »

Ce n'est **pas** une obligation d'utiliser exactement cette version, mais une indication de la version minimale du langage et des fonctionnalités attendues par le module.

---

# Pourquoi le chemin GitHub ?

C'est une convention très élégante de Go.

L'identifiant du module est **unique au niveau mondial**.

Il y a des millions de projets Go, mais un seul :

```
github.com/grapinou/club-manager
```

Cette convention évite les conflits de noms.

---

# Pourquoi `go.sum` n'est-il pas encore là ?

Tu pourrais te dire :

> « On m'a parlé de `go.sum`, mais il n'existe pas. »

C'est normal.

Pour l'instant ton projet n'a **aucune dépendance externe**.

Le jour où nous installerons, par exemple, le pilote PostgreSQL :

```
go get github.com/jackc/pgx/v5
```

Go téléchargera cette bibliothèque et créera automatiquement :

```
go.sum
```

Ce fichier enregistrera les empreintes cryptographiques des dépendances téléchargées pour vérifier leur intégrité.

---

# Nous arrivons au bout de l'initialisation

Notre checklist est maintenant :

```
✅ Créer le dossier du projet
✅ git init
✅ README
✅ .gitignore
✅ Premier commit
✅ Dépôt GitHub
✅ Configuration SSH
✅ Premier push
✅ go mod init
⬜ Commit de go.mod
```

Il ne reste plus qu'une étape.

---

# Le dernier commit d'initialisation

Vérifie d'abord :

```
git status
```

Tu devrais voir :

```
new file: go.mod
```

Puis :

```
git add go.mod
```

Ensuite :

```
git commit -m "Initialize Go module"
```

Enfin :

```
git push
```

Grâce au `-u` que nous avons utilisé précédemment, un simple `git push` suffira.

---

## Un petit bilan

Je trouve que nous avons pris une bonne décision en avançant lentement.

Beaucoup de tutoriels enchaînent en quelques minutes :

- `git init`
- `go mod init`
- `git push`

Sans expliquer ce qui se passe.

Ici, tu sais maintenant :

- ce qu'est un dépôt Git local ;
- ce qu'est un dépôt distant (`origin`) ;
- pourquoi on utilise une clé SSH ;
- à quoi sert un module Go ;
- pourquoi son nom correspond au dépôt GitHub.

À mon avis, c'est cette compréhension qui fera la différence lorsque tu commenceras à concevoir **Club Manager**. Les commandes deviennent alors des outils que tu maîtrises, plutôt qu'une suite d'instructions à mémoriser.