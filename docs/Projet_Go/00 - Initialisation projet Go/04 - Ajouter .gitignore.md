MOC : [[00 - INITIALISATION PROJET]]

---

# Comprendre le fichier `.gitignore`

On ignore les fichiers générés, temporaires, personnels ou contenant des secrets. On versionne le code source, la documentation, les fichiers de configuration nécessaires au projet et les dépendances déclaratives (go.mod, go.sum).

Particularité
	.gitignore :

	commence par un point ;

	n'a pas d'extension ;

	est un fichier caché sous Linux.

Le fichier `.gitignore` indique à Git quels fichiers ou dossiers ne doivent **pas** être suivis (versionnés).

L'objectif est de ne conserver dans le dépôt que les fichiers nécessaires au développement du projet.

|Élément|Pourquoi l'ignorer ?|Peut-il être régénéré ?|
|---|---|---|
|`*.exe`|Exécutable compilé sous Windows. Il est généré à partir du code source.|✅ Oui|
|`*.dll`|Bibliothèque dynamique Windows générée lors de la compilation.|✅ Oui|
|`*.so`|Bibliothèque partagée sous Linux générée lors de la compilation.|✅ Oui|
|`*.dylib`|Bibliothèque dynamique sous macOS.|✅ Oui|
|`*.test`|Exécutable généré par `go test -c`.|✅ Oui|
|`*.out`|Fichier de sortie généré par certains outils Go (tests, profils, couverture…).|✅ Oui|
|`coverage.out`|Rapport de couverture des tests (`go test -coverprofile`).|✅ Oui|
|`.env`|Contient des informations sensibles (mots de passe, clés API, secrets...). Chaque environnement possède généralement son propre fichier.|❌ Non|
|`.vscode/`|Paramètres personnels de Visual Studio Code (fenêtres, extensions, préférences...).|❌ Non|
|`.idea/`|Paramètres personnels des IDE JetBrains (GoLand, IntelliJ...).|❌ Non|

## À ne pas ignorer

Les fichiers suivants doivent être versionnés car ils sont indispensables au projet :

|Élément|Pourquoi le conserver ?|
|---|---|
|`go.mod`|Déclare le module Go et les dépendances directes du projet.|
|`go.sum`|Garantit l'intégrité et les versions exactes des dépendances téléchargées.|
|`README.md`|Documente le projet et son fonctionnement.|
|`LICENSE`|Indique la licence du projet.|
|`Makefile` (si présent)|Automatise les commandes courantes du projet.|

## Principe général

Avant d'ajouter un fichier au `.gitignore`, se poser la question :

> **« Si un autre développeur clone le dépôt, ce fichier est-il nécessaire pour développer ou exécuter le projet ? »**

- **Oui** → le versionner.
    
- **Non, il est généré automatiquement ou personnel** → l'ignorer.