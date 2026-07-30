MOC : [[00 - BASE GIT]]

---

## Rôle

La commande `git commit` enregistre les fichiers présents dans la **zone de préparation** dans l'historique du projet.

Chaque commit représente un **point de sauvegarde** accompagné d'un message décrivant les modifications réalisées.

---

## Syntaxe

```bash
git commit -m "Message du commit"
```

Exemple :

```bash
git commit -m "Ajout de la gestion des membres"
```

L'option `-m` permet d'écrire directement le message du commit.

---

## Quand l'utiliser ?

Une fois que les fichiers ont été ajoutés dans la **zone de préparation** avec `git add`.

Workflow classique :

```text
Modifier les fichiers
        │
        ▼
git status
        │
        ▼
git add .
        │
        ▼
git commit -m "Description des modifications"
```

---

## Exemples

### Créer un commit

```bash
git commit -m "Correction du calcul des cotisations"
```

### Ajouter une nouvelle fonctionnalité

```bash
git commit -m "Ajout de l'authentification"
```

---

## Bonnes pratiques

- Le message doit être court et décrire clairement la modification.
- Un commit doit idéalement correspondre à une seule évolution ou correction.
- Il est préférable de réaliser plusieurs petits commits qu'un très gros.

Exemples de bons messages :

```text
Ajout de la page de connexion
Correction du calcul des licences
Suppression du code inutile
```

---

## À retenir

- `git commit` crée un nouveau point dans l'historique Git.
- Seuls les fichiers présents dans la **zone de préparation** sont enregistrés.
- Les modifications non ajoutées avec `git add` ne seront pas incluses dans le commit.

---

## Comprendre

Un commit est comme une photo de ton projet à un instant donné.

Git ne sauvegarde pas tout le projet à chaque fois, mais il mémorise les différences entre les versions.

Grâce aux commits, il est possible de retrouver une ancienne version du projet ou de comprendre l'évolution du code.

---

## Voir aussi

- git add
- git status
- git push
- git log