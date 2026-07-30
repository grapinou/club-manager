MOC : [[00 - INITIALISATION PROJET]]

---

## Étape Git : initialiser le dépôt

Objectif : créer un dépôt Git local.

Pour l'instant, rien n'est envoyé sur GitHub.

Schéma :
```
Ordinateur
    |
    |
    v
Dépôt Git local
    |
    |
    v
Dépôt GitHub (plus tard)
```

git init crée simplement l'historique du projet sur ta machine.

Il crée un dossier caché `.git` contenant :

- l'historique ;
- les commits ;
- la configuration du dépôt.
  
## Commande

Dans le dossier du projet :

```
git init
```

Tu devrais obtenir :

```
Initialized empty Git repository in /home/sighto/club-manager/.git/
```

---

## Ce qui vient de se passer ?

Git a créé un dossier caché :

```
club-manager
│
└── .git
```

Ce dossier contient :

- l'historique des versions ;
- les informations de configuration du dépôt ;
- les branches ;
- les commits.

⚠️ On ne touche jamais directement au dossier `.git`.

---

# Vérifier que Git fonctionne

Commande :

```
git status
```

Résultat probable :

```
On branch master

No commits yet

nothing to commit
```

ou :

```
On branch main

No commits yet

nothing to commit
```

Cela signifie :

- Git fonctionne ;
- le dépôt existe ;
- aucun fichier n'est encore suivi.

# Changer le nom master en main

Actuellement, ton `git init` peut encore créer `master` selon ta configuration.

Tu peux vérifier :

```
git config --global init.defaultBranch
```

Si rien ne s'affiche, tu peux définir :

```
git config --global init.defaultBranch main
```

Ensuite, les prochains :

```
git init
```

créeront directement une branche `main`.

## Si ton dépôt actuel est déjà en `master`

Ce n'est pas grave.

Tu peux voir la branche actuelle :

```
git branch
```

Si tu obtiens :

```
* master
```

Remarque : si aucun commit n'a été fait, aucun retour n'est donné ! il n'y a pas de branche matérialisée

Tu peux la renommer :

```
git branch -m master main
```

fonctionne même si aucun commit n'a été fait

Explication :

- `branch` : manipuler les branches ;
- `-m` : move/rename (renommer).

## Pour vérifier quelle branche Git prévoit d'utiliser

Tu peux faire :

```
git symbolic-ref --short HEAD
```

Si tu as configuré `main`, cela devrait répondre :

```
main
```

