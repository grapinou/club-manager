MOC : [[00 - INITIALISATION PROJET]]

---

Je résumerais les commandes ainsi :

|Commande|Rôle|Vérification|
|---|---|---|
|`git status`|Affiche l'état du dépôt et des fichiers.|Vérifier les fichiers suivis, modifiés ou en attente.|
|`git add .`|Ajoute tous les fichiers du dossier courant à la zone de préparation.|`git status` affiche **Changes to be committed**.|
|`git commit -m "Message"`|Enregistre un instantané des fichiers préparés dans l'historique Git.|`git log` affiche le nouveau commit.|

---

## Une remarque sur les messages de commit

Comme tu souhaites que ce projet soit une vitrine de ton travail, je te proposerai de suivre une règle simple :

> **Chaque message de commit doit expliquer "ce qui change", pas "ce que tu as fait".**

Par exemple :

- ✅ `Initial project structure`
- ✅ `Add README describing project goals`
- ✅ `Initialize Go module`
- ✅ `Configure PostgreSQL connection`

Ces messages sont courts, mais ils racontent l'évolution du projet. Si quelqu'un lit l'historique un an plus tard, il comprendra immédiatement la progression.

## Git log

Tu viens de faire quelque chose de très important : tu as créé le premier point de référence de ton projet. À partir de maintenant, Git peut revenir à cet état à tout moment.

Ton git log montre :

```
commit 302790060835f27912c23a0482187a07661c0d84 (HEAD -> main)
Author: grapinou <remigrapinou@gmail.com>
Date:   Fri Jul 24 16:00:24 2026 +0200

    Add README and gitignore
```

Quelques explications.

### Le hash du commit

	302790060835f27912c23a0482187a07661c0d84

C'est l'identifiant unique du commit.

En pratique, Git accepte généralement les 7 ou 8 premiers caractères :

	3027900

C'est suffisant pour identifier ce commit dans ton dépôt.

### HEAD

Tu vois :

	HEAD -> main
C'est une notion fondamentale.
HEAD signifie :

	"L'endroit où je me trouve actuellement."

Pour l'instant :

```
main
 │
 ▼
3027900
 ↑
HEAD
```

Quand tu feras un nouveau commit, HEAD avancera automatiquement.

Le message du commit
Tu as choisi :
	Add README and gitignore
Je le trouve très bien.

Il est :
- court ;
- explicite ;
- rédigé en anglais (c'est la convention la plus répandue).

C'est exactement le type de message que je mettrais sur un projet GitHub.

Une petite astuce pour git log
À partir de maintenant, je te conseille plutôt d'utiliser :

Bash
```
git log --oneline --graph --decorate
```

Pour l'instant, tu verras :

	 3027900 (HEAD -> main) Add README and gitignore
	 
Mais dans quelques mois, quand tu auras 200 commits, cela donnera quelque chose comme :

```
5c7d2e1 (HEAD -> main) Add accounting module
8fd1203 Add role permissions
42bc781 Add authentication
3027900 Add README and gitignore
```

Je trouve cette vue beaucoup plus agréable.