MOC : [[00 - BASE GIT]]

---

## Rôle

La commande `git log` permet de consulter l'historique des commits d'un dépôt Git.

Elle affiche les différents points de sauvegarde du projet dans l'ordre chronologique inverse (du plus récent au plus ancien).

---

## Syntaxe

### Historique simplifié

```bash
git log --oneline
```

Chaque commit est affiché sur une seule ligne.

---

## Quand l'utiliser ?

- Pour retrouver un ancien commit.
- Pour suivre l'évolution du projet.
- Pour vérifier qu'un commit a bien été créé.

Workflow classique :

```text
git commit
      │
      ▼
git log --oneline
```

---

## Exemples

### Afficher l'historique

```bash
git log --oneline
```

Exemple de résultat :

```text
8f7d2a1 Ajout de la gestion des licences
2c91ab4 Correction de l'authentification
f13b5e7 Premier commit
```

---

## Bonnes pratiques

- Préférer `git log --oneline` à `git log` pour une lecture rapide.
- Écrire des messages de commit explicites afin que l'historique reste compréhensible.
- Consulter régulièrement l'historique pour suivre l'évolution du projet.

---

## 💡 Comprendre

Chaque commit possède un identifiant unique appelé **hash**.

Remarque : 
	**Chaque commit reçoit une référence unique générée par Git.**
	
	Cette référence permet à Git de retrouver précisément un commit, même plusieurs années plus tard.

Par exemple :

```text
8f7d2a1 Ajout de la gestion des licences
```

- `8f7d2a1` est le début du hash du commit.
- `Ajout de la gestion des licences` est le message que vous avez écrit avec `git commit`.

L'historique Git raconte l'évolution de votre projet.

Chaque commit représente une étape importante.

---

## À retenir

- `git log --oneline` affiche l'historique de manière compacte.
- Les commits sont affichés du plus récent au plus ancien.
- Chaque commit possède un identifiant unique (hash).
- La qualité de l'historique dépend de la qualité des messages de commit.

---

## Voir aussi

- git commit
- git diff

### Étape suivante

Après avoir retrouvé un commit, vous pouvez poursuivre votre développement ou explorer d'autres fonctionnalités de Git.