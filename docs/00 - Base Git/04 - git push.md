MOC : [[00 - BASE GIT]]

---

## Rôle

La commande `git push` envoie les commits de votre dépôt local vers un dépôt distant, comme GitHub.

Elle permet de partager vos modifications et de sauvegarder votre travail sur le serveur distant.

---

## Syntaxe

```bash
git push
```

Lors du premier envoi vers une nouvelle branche, Git peut demander de préciser le dépôt distant et la branche :

```bash
git push -u origin main
```

- `origin` : nom du dépôt distant.
- `main` : branche sur laquelle envoyer les commits.
- `-u` : associe la branche locale à la branche distante afin que les prochains `git push` puissent être exécutés simplement avec `git push`.

---

## Quand l'utiliser ?

Après avoir créé un ou plusieurs commits que vous souhaitez envoyer sur le dépôt distant.

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
        │
        ▼
git push
```

---

## Exemples

### Envoyer les derniers commits

```bash
git push
```

### Premier envoi d'une branche

```bash
git push -u origin main
```

---

## Bonnes pratiques

- Vérifier que le commit est terminé avant de faire un `git push`.
- Réaliser des commits cohérents avant de les envoyer.
- Effectuer un `git pull` si d'autres personnes travaillent sur le même dépôt.

---

## 💡 Comprendre

Un `git commit` enregistre votre travail **sur votre ordinateur**.

Un `git push` copie ensuite ces commits vers le dépôt distant (par exemple GitHub).

On peut voir Git comme deux historiques :

```text
Votre ordinateur                    GitHub

Commits locaux  ───────────────►  Commits distants
                    git push
```

Tant que vous n'avez pas exécuté `git push`, vos commits existent uniquement sur votre machine.

---

## À retenir

- `git push` envoie les commits vers le dépôt distant.
- Il n'envoie que les commits qui n'existent pas encore sur le serveur.
- Aucun nouveau commit n'est créé.
- Les fichiers doivent déjà avoir été commités avant d'être envoyés.

---

## Voir aussi

- git commit
- git pull
- git status
- git log