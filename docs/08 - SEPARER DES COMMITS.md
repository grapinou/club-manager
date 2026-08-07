MOC : [[00 - BASE GIT]]

---

Le principe est simplement : **on ne fait pas `git add .` tout de suite**. On choisit ce qui entre dans chaque commit.

### 1. Commit du code

Tu peux ajouter uniquement les fichiers applicatifs :

```
git add internal/handlers internal/views
```

Puis vérifier ce qui va réellement être commité :

```
git status
```

et, encore mieux :

```
git diff --staged
```

Tu devrais retrouver uniquement Contact + Rules, leurs vues, templates et tests.

Ensuite :

```
git commit -m "Render contact and rules pages with templates"
```

### 2. Commit de la documentation

Il ne reste alors normalement que les fiches Obsidian :

```
git status
```

Tu peux les ajouter :

```
git add docs
```

Vérification :

```
git diff --staged
```

Puis :

```
git commit -m "Document configuration and closures"
```

Enfin :

```
git push
```

### Si tu avais déjà fait `git add .`

Ce n'est pas grave. Tu peux simplement retirer les fichiers de la zone de staging **sans supprimer tes modifications** :

```
git restore --staged .
```

Puis reprendre :

```
git add internal/handlers internal/views
git commit -m "Render contact and rules pages with templates"

git add docs
git commit -m "Document configuration and closures"

git push
```

C'est un excellent exemple pour comprendre que `git add` ne sert pas seulement à « préparer un commit » : il sert aussi à **choisir précisément quelles modifications appartiennent au prochain commit**.

```
Working directory
├── modifications code ──git add──→ commit 1
└── modifications docs ──git add──→ commit 2
```

Et dans ton diff actuel, la séparation `internal/` d'un côté et `docs/` de l'autre rend l'opération particulièrement simple.

