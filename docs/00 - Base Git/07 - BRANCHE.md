MOC : [[00 - BASE GIT]]

---

## Rôle

Une branche permet de travailler sur une évolution du projet sans modifier directement la branche principale (`main`).

Elle crée une copie de travail indépendante à partir d'un point donné de l'historique.

---

## Quand l'utiliser ?

Lorsqu'une modification peut être développée séparément du code principal :

- Nouvelle fonctionnalité
- Correction importante
- Expérimentation
- Test d'une idée

Exemple :

```text
main
 |
 |------ Nouvelle fonctionnalité
          |
          branche "licences"
```

---

## Procédure classique

### 1. Créer une branche

```bash
git switch -c nom-de-la-branche
```

Exemple :

```bash
git switch -c gestion-licences
```

Cette commande :
- crée la nouvelle branche ;
- bascule automatiquement dessus.

---

### 2. Travailler sur la branche

Modifier les fichiers puis créer des commits :

```bash
git add .
git commit -m "Ajout de la gestion des licences"
```

Les commits sont enregistrés uniquement dans cette branche.

---

### 3. Revenir sur main

```bash
git switch main
```

La branche principale redevient active.

---

### 4. Fusionner la branche

Lors d'un `git merge`, **il faut être positionné sur la branche qui doit recevoir les modifications**.


```bash
git switch main
git merge nom-de-la-branche
```

Exemple :

```bash
git merge gestion-licences
```

Les modifications de la branche sont intégrées dans `main`.

Ici :

- tu es sur `main` ✅
- tu demandes à Git : "prends les changements de `nom-de-la-branche` et ajoute-les dans `main`".

---

### 5. Supprimer la branche

Une fois la fusion terminée :

```bash
git branch -d nom-de-la-branche
```

Exemple :

```bash
git branch -d gestion-licences
```

---

## 💡 Comprendre

Une branche est simplement un pointeur vers un commit.

Elle permet de créer un chemin parallèle dans l'historique.

Exemple :

```text
A --- B --- C -------- F   main
           \
            D --- E        gestion-licences
```

La branche `gestion-licences` contient les commits D et E.

Après une fusion :

```text
A --- B --- C -------- F   main
           \          /
            D --- E --
```

Les modifications rejoignent l'historique principal.

---

## À retenir

- Une branche permet de développer sans toucher à `main`.
- `git switch -c` crée une branche et se positionne dessus.
- Les commits réalisés appartiennent à la branche active.
- `git merge` intègre une branche dans une autre.
- Une branche terminée peut être supprimée.
- La branche active est la branche qui reçoit la fusion.

Remarque : 

### Que se passerait-il si tu étais sur la mauvaise branche ?

Imaginons :

```
git switch gestion-licences
git merge main
```

Git va intégrer `main` **dans** `gestion-licences`.

Ce n'est pas une erreur technique : Git fera exactement ce que tu lui demandes.

Mais le résultat serait différent :

```
gestion-licences

A --- B --- C --- D --- E --- (fusion de main)
```

Tu aurais ramené les modifications de `main` dans ta branche de fonctionnalité.

---

## Voir aussi

- git commit
- git log
