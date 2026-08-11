MOC : [[00 - BASE GIT]]

---

## Rôle

La commande `git add` ajoute un ou plusieurs fichiers dans **la zone de préparation** (*staging area ou index*).

Les fichiers ajoutés seront inclus dans le prochain commit.

---

## Syntaxe

### Ajouter un fichier

```bash
git add nom-du-fichier
```

Exemple :

```bash
git add main.go
```

---

### Ajouter tous les fichiers

```bash
git add .
```

Le point (`.`) représente le dossier courant. Git ajoute tous les fichiers modifiés, créés ou supprimés présents dans ce dossier et ses sous-dossiers.

---

## Quand l'utiliser ?

Après avoir terminé une ou plusieurs modifications et avant de créer un commit.

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

### Ajouter un seul fichier

```bash
git add README.md
```

Pratique lorsque seules certaines modifications doivent être enregistrées.

---

### Ajouter tous les fichiers du projet

```bash
git add .
```

C'est la commande la plus utilisée pour les projets personnels.

---

## Bonnes pratiques

- Vérifier l'état du dépôt avec `git status` avant d'utiliser `git add`.
- Si possible, utiliser `git diff` pour relire les modifications avant de les préparer.
- Après un `git add`, exécuter de nouveau `git status` pour vérifier que les bons fichiers sont prêts à être commités.

---

## À retenir

- `git add` **n'enregistre rien** dans l'historique.
- Il prépare simplement les fichiers pour le prochain commit.
- Tant qu'un commit n'a pas été créé, il est possible de continuer à modifier les fichiers.

---

## Voir aussi

- git status
- git commit
- git diff

