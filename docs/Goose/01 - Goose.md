
---

# Goose — Installation et configuration

## Objectif

Installer **Goose**, l'outil de migration que nous utiliserons avec PostgreSQL dans **Club Manager**, et comprendre :

- où Go installe les exécutables ;
    
- le rôle de `GOPATH` ;
    
- le rôle de `GOBIN` ;
    
- le rôle de `PATH` ;
    
- comment le shell trouve une commande.
    

---

# 1. Vérifier si Goose est installé

Commande :

```bash
goose --version
```

Version présente initialement :

```text
goose version: v3.24.3
```

Goose était donc déjà installé.

Pour comprendre son installation, nous avons recherché l'emplacement de l'exécutable.

---

# 2. Localiser une commande

Deux commandes permettent notamment de savoir quel exécutable sera utilisé :

```bash
which goose
```

ou :

```bash
command -v goose
```

Résultat initial :

```text
/usr/local/bin/goose
```

Le shell utilisait donc :

```text
/usr/local/bin/goose
```

---

# 3. Vérifier la configuration Go

Commande :

```bash
go env GOBIN
```

Résultat :

`GOBIN` n'était pas défini.

Puis :

```bash
go env GOPATH
```

Résultat :

```text
/home/sighto/go
```

---

# `GOPATH`

`GOPATH` représente un répertoire de travail utilisé par Go pour différents fichiers.

Dans notre configuration :

```text
GOPATH
   ↓
/home/sighto/go
```

Lorsqu'aucun `GOBIN` particulier n'est défini, `go install` place normalement les exécutables dans :

```text
$GOPATH/bin
```

Donc ici :

```text
/home/sighto/go/bin
```

---

# `GOBIN`

`GOBIN` permet de choisir explicitement le dossier dans lequel Go installe les exécutables.

Commande :

```bash
go env GOBIN
```

Dans notre cas, la valeur était vide.

Go utilise donc :

```text
$GOPATH/bin
```

comme emplacement par défaut.

---

# 4. Examiner l'ancienne installation

L'exécutable précédent se trouvait ici :

```text
/usr/local/bin/goose
```

Commande :

```bash
file /usr/local/bin/goose
```

Résultat :

```text
ELF 64-bit LSB executable, x86-64
```

Cela confirme qu'il s'agissait d'un véritable exécutable Linux compilé.

---

# 5. Vérifier le propriétaire du fichier

Commande :

```bash
ls -l /usr/local/bin/goose
```

Résultat :

```text
-rwxrwxr-x 1 sighto sighto ... /usr/local/bin/goose
```

Le fichier appartenait donc à :

```text
utilisateur : sighto
groupe      : sighto
```

---

# 6. Vérifier le propriétaire du dossier

Commande :

```bash
ls -ld /usr/local/bin
```

Résultat :

```text
drwxr-xr-x ... root root ... /usr/local/bin
```

Le répertoire appartient donc à :

```text
root
```

---

# Notion Linux : supprimer un fichier

Pour modifier le contenu d'un fichier, les permissions du fichier sont importantes.

Pour :

- créer ;
    
- supprimer ;
    
- renommer ;
    

un fichier, les permissions du **répertoire qui le contient** sont déterminantes.

Ainsi :

```text
/usr/local/bin/goose
        │
        └── appartient à sighto

/usr/local/bin
        │
        └── appartient à root
```

La suppression nécessite donc ici des privilèges administrateur.

---

# 7. Désinstaller Goose

`go install` ne possède pas de commande :

```text
go uninstall
```

Pour désinstaller un outil Go installé sous forme d'exécutable, on supprime simplement son binaire.

Commande utilisée :

```bash
sudo rm /usr/local/bin/goose
```

Puis vérification :

```bash
command -v goose
```

et :

```bash
goose --version
```

Goose n'était alors plus accessible.

---

# 8. Installer Goose avec Go

Installation :

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Décomposition conceptuelle :

```text
code source de Goose
        ↓
go install
        ↓
compilation
        ↓
création de l'exécutable
        ↓
$GOPATH/bin
```

---

# 9. Vérifier l'emplacement du nouvel exécutable

Commande :

```bash
ls -l ~/go/bin/goose
```

Résultat :

```text
/home/sighto/go/bin/goose
```

Goose était donc correctement installé.

Mais :

```bash
command -v goose
```

ne trouvait pas encore cette nouvelle installation.

---

# 10. Le rôle de `PATH`

Le problème venait de la variable :

```text
PATH
```

`PATH` contient la liste des répertoires dans lesquels le shell recherche les commandes.

Commande :

```bash
echo $PATH
```

Exemple obtenu :

```text
/usr/local/sbin
/usr/local/bin
/usr/sbin
/usr/bin
...
/usr/local/go/bin
/home/sighto/.local/bin
```

Le dossier suivant était absent :

```text
/home/sighto/go/bin
```

Donc :

```text
Goose installé ✅
```

mais :

```text
commande goose introuvable ❌
```

---

# Comprendre `PATH`

Quand on tape :

```bash
goose
```

Bash cherche un fichier exécutable nommé `goose` dans les dossiers du `PATH`.

Conceptuellement :

```text
goose
  ↓
PATH
  ├── /usr/local/bin
  ├── /usr/bin
  ├── ...
  └── /home/sighto/go/bin
```

Le premier exécutable correspondant est utilisé.

---

# 11. Ajouter temporairement Go au PATH

Commande :

```bash
export PATH="$HOME/go/bin:$PATH"
```

Décomposition :

```text
PATH="$HOME/go/bin:$PATH"
```

signifie :

```text
nouveau PATH
│
├── $HOME/go/bin
└── ancien PATH
```

Ici :

```text
$HOME
```

vaut :

```text
/home/sighto
```

Donc :

```text
$HOME/go/bin
```

correspond à :

```text
/home/sighto/go/bin
```

---

# 12. Vérifier la résolution de la commande

Commande :

```bash
command -v goose
```

Résultat :

```text
/home/sighto/go/bin/goose
```

Le shell trouve maintenant le bon exécutable.

---

# 13. Vérifier la version

Commande :

```bash
goose --version
```

Résultat :

```text
goose version: v3.27.3
```

Goose est donc correctement installé.

---

# 14. Rendre le PATH permanent

La commande :

```bash
export PATH="$HOME/go/bin:$PATH"
```

ne modifie que le terminal courant.

Pour rendre cette configuration permanente, ajouter dans :

```text
~/.profile
```

la ligne :

```bash
export PATH="$HOME/go/bin:$PATH"
```

Puis recharger le fichier :

```bash
source ~/.profile
```

---

# Vérification finale

Commande :

```bash
command -v goose
```

Résultat :

```text
/home/sighto/go/bin/goose
```

Puis :

```bash
goose --version
```

Résultat :

```text
goose version: v3.27.3
```

Installation validée.

---

# Architecture obtenue

```text
/home/sighto
└── go
    └── bin
        └── goose
```

Avec :

```text
GOPATH
   ↓
/home/sighto/go
```

et :

```text
PATH
   ↓
contient /home/sighto/go/bin
```

Le shell peut donc trouver :

```text
/home/sighto/go/bin/goose
```

lorsque l'on tape simplement :

```bash
goose
```

---

# `GOPATH`, `GOBIN` et `PATH`

## `GOPATH`

Répertoire utilisé par Go.

Dans notre cas :

```text
/home/sighto/go
```

---

## `GOBIN`

Répertoire explicite dans lequel Go peut installer les exécutables.

Dans notre configuration :

```text
GOBIN = vide
```

Go utilise donc :

```text
$GOPATH/bin
```

---

## `PATH`

Variable du shell.

Elle indique où Linux doit chercher les commandes.

```text
Go
└── installe goose dans $GOPATH/bin

Bash
└── cherche goose dans $PATH
```

Ce sont deux mécanismes différents.

---

# Comprendre et retenir

## Installer un programme ne suffit pas toujours

Un exécutable peut parfaitement exister :

```text
/home/sighto/go/bin/goose
```

tout en donnant :

```text
goose: commande introuvable
```

si son répertoire n'est pas présent dans :

```text
PATH
```

---

## `go install` compile un exécutable

La commande :

```bash
go install github.com/pressly/goose/v3/cmd/goose@latest
```

ne fait pas seulement télécharger Goose.

Conceptuellement :

```text
module Go
   ↓
téléchargement
   ↓
compilation
   ↓
exécutable
   ↓
$GOPATH/bin/goose
```

---

## Goose n'est pas encore une dépendance de Club Manager

Nous avons installé :

```text
goose
```

comme **outil de développement**.

Il fonctionne depuis le terminal :

```text
développeur
    ↓
goose
    ↓
PostgreSQL
```

Cela ne signifie pas encore que le programme Go Club Manager importe Goose.

C'est comparable à :

```text
git
psql
goose
```

qui sont des outils utilisés autour du projet.

---

## `command -v` est très utile

Pour savoir quel exécutable sera réellement utilisé :

```bash
command -v goose
```

Même principe avec :

```bash
command -v go
command -v git
command -v psql
```

Cela permet de détecter rapidement :

- plusieurs installations ;
    
- un mauvais `PATH` ;
    
- un ancien exécutable encore présent.
    

---

# Commandes à retenir

|Commande|Fonction|
|---|---|
|`goose --version`|afficher la version de Goose|
|`command -v goose`|afficher l'exécutable utilisé|
|`which goose`|localiser une commande|
|`go env GOPATH`|afficher `GOPATH`|
|`go env GOBIN`|afficher `GOBIN`|
|`echo $PATH`|afficher les dossiers recherchés par le shell|
|`file fichier`|identifier le type d'un fichier|
|`ls -l fichier`|afficher propriétaire et permissions|
|`ls -ld dossier`|afficher les permissions d'un dossier|
|`go install ...@latest`|compiler et installer un outil Go|
|`export PATH="...:$PATH"`|modifier temporairement le PATH|
|`source ~/.profile`|recharger la configuration du profil|

---

# État final

```text
Goose              : installé ✅
Version             : v3.27.3
Exécutable          : /home/sighto/go/bin/goose
GOPATH              : /home/sighto/go
GOBIN               : non défini
$GOPATH/bin         : présent dans PATH ✅
Commande `goose`    : fonctionnelle ✅
```

---

# Prochaine étape

Goose est maintenant prêt à être utilisé avec Club Manager.

La prochaine étape sera de comprendre le principe d'une migration puis de créer notre première migration versionnée dans le dépôt :

```text
Club Manager
└── migrations
    └── première migration SQL
            ↓
          Goose
            ↓
       PostgreSQL
```

À partir de là, les modifications de structure de la base ne seront plus faites manuellement : elles pourront être décrites dans des fichiers SQL et suivies avec Git.