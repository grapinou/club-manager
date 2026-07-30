MOC : [[00 - BASE GIT]]

---

## Rôle

La commande `git diff` affiche les différences entre les fichiers de votre dossier de travail et la **zone de préparation**.

Elle permet de relire les modifications avant de les ajouter avec `git add`.

---

## Syntaxe

```bash
git diff
```

---

## Quand l'utiliser ?

Après avoir modifié des fichiers et avant de les ajouter dans la **zone de préparation**.

Workflow classique :

```text
Modifier les fichiers
        │
        ▼
git status
        │
        ▼
git diff
        │
        ▼
git add .
```

---

## Exemples

### Voir toutes les modifications

```bash
git diff
```

### Voir les modifications d'un fichier

```bash
git diff main.go
```

---

## Bonnes pratiques

- Relire les modifications avant de faire un `git add`.
- Vérifier qu'aucun code de test ou de débogage n'est présent.
- Utiliser `git diff` pour comprendre précisément ce qui a changé.

---

## 💡 Comprendre

`git diff` compare deux versions d'un fichier.

Dans son utilisation la plus courante, il compare :

```text
Dossier de travail
        │
        │ git diff
        ▼
Zone de préparation
```

Il affiche :

- les lignes supprimées (précédées d'un `-`) ;
- les lignes ajoutées (précédées d'un `+`).

Une fois les modifications ajoutées avec `git add`, elles ne sont plus affichées par `git diff`.

---

## À retenir

- `git diff` ne modifie jamais le dépôt.
- Il permet de relire les modifications avant un commit.
- Il compare les fichiers modifiés avec la **zone de préparation**.
- C'est une excellente habitude à prendre avant d'utiliser `git add`.

---

## Voir aussi

- git status
- git add
- git commit

### Étape suivante

Une fois les modifications vérifiées, prépare-les avec git add.