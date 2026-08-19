
# Club Manager — Index des fiches

Cet espace regroupe la documentation construite pendant le développement de **Club Manager**.

L'organisation suit principalement les **jalons du projet**.

L'idée est simple :

> Une fiche est rangée dans le jalon où la notion a été rencontrée et utilisée pour la première fois.

Ainsi, chaque jalon constitue une photographie de l'état du projet et des connaissances nécessaires pour le comprendre à ce moment-là.

---

# Organisation

```text
docs/
│
├── Jalon/
│   ├── Jalon 00/
│   ├── Jalon 01/
│   ├── Jalon 02/
│   ├── ...
│   └── Jalon 08/
│
├── Git général/
│
└── Vim/
```

Les dossiers `Git général` et `Vim` contiennent des références transversales qui ne correspondent pas à une étape particulière de Club Manager.

---

# Progression de Club Manager

```text
Jalon 00
Initialisation du projet
        │
        ▼
Jalon 01
Handlers et routes
        │
        ▼
Jalon 02
Premier template
        │
        ▼
Jalon 03
Layout commun
        │
        ▼
Jalon 04
Vers un site configurable
        │
        ▼
Jalon 05
Site associatif configurable
        │
        ▼
Jalon 06
Configuration et pages génériques
        │
        ▼
Jalon 07
Interface publique configurable
        │
        ▼
Jalon 08
PostgreSQL, psql et Goose
```

---

# Jalon 00 — Initialisation du projet

## Objectif

Créer le projet Club Manager et disposer d'un premier serveur web Go minimal.

Ce jalon pose les fondations :

```text
environnement
    +
Git
    +
module Go
    +
architecture minimale
    +
serveur HTTP
```

## Fiche principale

- [[00 - INITIALISATION PROJET|Initialisation du projet]]
    

## Mise en place de l'environnement

- [[01 - Installer la dernière version de Go|Installer Go]]
    
    - [[01.01 - wget|wget]]
        
    - [[01.02 - tar|tar]]
        
    - [[01.03 - absolu|Chemins absolus et relatifs]]
        

## Git et GitHub

- [[02 - Vérifier si Git est bien installé|Vérifier l'installation de Git]]
    
- [[03 - git init|Initialiser un dépôt Git]]
    
- [[04 - Ajouter .gitignore|Créer le .gitignore]]
    
- [[05 - Premier commit|Premier commit]]
    
- [[06 - Créer le dépôt GitHub|Créer le dépôt GitHub]]
    
    - [[06.01 - clé ssh|Clé SSH]]
        

### Workflow Git de base

- [[00 - BASE GIT|Base Git]]
    
- [[01 - git status|git status]]
    
- [[02 - git add|git add]]
    
- [[03 - git commit|git commit]]
    
- [[04 - git push|git push]]
    
- [[05 - git diff|git diff]]
    
- [[06 - git log|git log]]
    

## Module Go

- [[07 - go mod init|Initialiser le module Go]]
    

## Architecture minimale

- [[00 - ARCHITECTURE MINIMALE EN GO|Architecture minimale en Go]]
    
    - [[00.01 - mkdir -p|mkdir -p]]
        
- [[00 - SERVEUR GO|Serveur HTTP Go]]
    

## Notions Go rencontrées

### Fonctions

- [[01 - FONCTION|Fonctions]]
    
    - [[01.01 - Fonction exportation|Fonctions exportées]]
        

### Packages

- [[02 - PACKAGE|Packages]]
    
    - [[02.01 - Package import|Imports]]
        

### Structs et méthodes

- [[03 - STRUCT ET METHODE|Structs et méthodes]]
    

### Interfaces

- [[04 - INTERFACE|Interfaces]]
    
    - [[04.01 - Interface composée|Interfaces composées]]
        
    - [[04.02 - Interface utilisation|Utiliser une interface]]
        

### Pointeurs

- [[05 - POINTEUR|Pointeurs]]
    

### Écriture et affichage

- [[06 - PRINT|Fonctions d'affichage]]
    
    - [[06.01 - fmt.Fprintln|fmt.Fprintln]]
        
    - [[06.02 - io.Writer|io.Writer]]
        

## État obtenu

```text
Navigateur
    │
    ▼
Serveur HTTP Go
    │
    ▼
Réponse texte
```

---

# Jalon 01 — Handlers et routes

## Objectif

Séparer les différentes responsabilités du premier serveur HTTP.

Le projet passe progressivement de :

```text
main.go
└── fait presque tout
```

à :

```text
main
 │
 ▼
router
 │
 ▼
handlers
```

## Fiche de jalon

- [[Jalon 01 - Handlers et routes]]
    

## Architecture

- [[01 - SEPARATION DES RESPONSABILITES|Séparation des responsabilités]]
    

## Handlers

- [[07 - HANDLER|Handlers HTTP]]
    

## Routeur

- [[08 - ROUTEUR|Routeur HTTP]]
    
    - [[08.01 - Routeur ServeMux|ServeMux]]
        

## Constructeur `New`

- [[01.02 - Fonction New|Convention New]]
    

## Tests

- [[09 - TESTS|Tests en Go]]
    
    - [[09.01 - Anatomie d'un premier test HTTP en Go|Anatomie d'un test HTTP]]
        
    - [[09.02 t.Helper()|t.Helper()]]
        

## Gestion des erreurs

- [[10 - ERREURS|Erreurs en Go]]
    

## État obtenu

```text
Navigateur
    │
    ▼
Router
    │
    ▼
Handler
    │
    ▼
Réponse
```

Les responsabilités du serveur commencent à être séparées.

---

# Jalon 02 — Premier template

## Objectif

Séparer le traitement HTTP de la génération du HTML.

## Fiche de jalon

- [[Jalon 02 - Premier template]]
    

## Architecture

- [[02 - SEPARATION ENTRE HANDLER ET VIEW|Séparation entre handler et view]]
    

## Templates Go

- [[11 - TEMPLATES|Templates]]
    
    - [[11.01 - embed|embed]]
        
    - [[11.02 - ParseFS|ParseFS]]
        
    - [[11.03 - Passage de données|Passage de données]]
        

## Architecture obtenue

```text
Navigateur
    │
    ▼
Router
    │
    ▼
Handler
    │
    ▼
HomeData
    │
    ▼
View
    │
    ▼
Template
    │
    ▼
HTML
```

Une nouvelle couche apparaît :

```text
internal/views/
```

---

# Jalon 03 — Layout commun

## Objectif

Éviter de répéter la structure HTML complète dans chaque page.

## Fiche de jalon

- [[Jalon 03 - Layout commun]]
    

## Composition des templates

- [[11.04 - Layout commun et composition des templates|Layout commun et composition]]
    

## Architecture obtenue

```text
             base.html
                │
        ┌───────┴───────┐
        │               │
structure commune    contenu
                        │
                    page.html
```

Une page devient :

```text
layout commun
      +
contenu spécifique
      =
page complète
```

---

# Jalon 04 — Vers un site configurable

## Objectif

Séparer les informations propres à une association du code de Club Manager.

Une nouvelle chaîne apparaît :

```text
config.json
    │
    ▼
config.Load()
    │
    ▼
Config
    │
    ▼
main
    │
    ▼
router
    │
    ▼
handlers
```

## Fiche de jalon

- [[Jalon 04 — Vers un site configurable]]
    

## Configuration et JSON

- [[12 - FICHIER DE CONFIGURATION ET JSON|Fichier de configuration et JSON]]
    
    - [[12.01 - Charger et tester une configuration JSON en Go|Charger et tester une configuration]]
        
    - [[12.02 - Transmettre une configuration de main jusqu'aux handlers|Transmettre Config jusqu'aux handlers]]
        

## Closures

- [[01.03 - Comprendre les closures|Closures]]
    

Les closures sont notamment utilisées pour construire des handlers conservant l'accès à la configuration.

## Logs et formatage

- [[06.03 - fmt ou log|fmt ou log]]
    
- [[06.04 - Les verbes de formatage avec fmt|Verbes de formatage avec fmt]]
    

Notions rencontrées notamment :

```text
%v
%q
%w
```

## Git

- [[08 - SEPARER DES COMMITS|Séparer plusieurs changements en commits]]
    

## Architecture obtenue

```text
              config.json
                   │
                   ▼
                 Config
                   │
                   ▼
Navigateur → Router → Handler → View → Template
```

Club Manager commence à séparer :

```text
logiciel générique
```

de :

```text
informations propres à l'association
```

---

# Jalon 05 — Site associatif configurable

## Objectif

Généraliser progressivement l'utilisation de la configuration aux différentes pages du site.

## Fiche de jalon

- [[Jalon 05 - Site associatif configurable]]
    

## Évolution principale

Le principe :

```text
contenu écrit dans les handlers
```

évolue vers :

```text
config.json
    │
    ▼
Config
    │
    ▼
handlers
    │
    ▼
views
```

Le projet devient progressivement réutilisable pour plusieurs associations.

---

# Jalon 06 — Configuration et pages génériques

## Objectif

Observer les duplications apparues entre les différentes pages et introduire une abstraction commune lorsqu'elle devient réellement utile.

## Fiche de jalon

- [[Jalon 06 - Configuration et pages génériques]]
    

## Évolution principale

Plusieurs pages peuvent maintenant suivre le même parcours :

```text
PageConfig
    │
    ▼
pageHandler
    │
    ▼
PageData
    │
    ▼
RenderPage
    │
    ▼
template générique
```

Cette étape marque une évolution importante :

> L'abstraction apparaît pour supprimer une duplication réellement observée dans le projet.

---

# Jalon 07 — Interface publique configurable

## Objectif

Transformer l'application Go fonctionnelle en première véritable interface publique utilisable.

## Fiche de jalon

- [[Jalon 07 - Interface publique configurable]]
    

## Fichiers statiques

- [[13 - STATIC EN GO|Fichiers statiques avec Go]]
    

Notions utilisées :

```text
http.FileServer
http.Dir
http.StripPrefix
```

Le serveur sait maintenant fournir :

```text
/static/css/...
/static/images/...
```

## Bootstrap

- [[14 - BOOTSTRAP|Bootstrap]]
    

Bootstrap apporte notamment :

- une mise en page responsive ;
    
- une navbar ;
    
- des containers ;
    
- des espacements cohérents ;
    
- une présentation générale simple.
    

## Architecture obtenue

```text
                    config.json
                         │
                         ▼
                       Config
                         │
                         ▼
Navigateur ──► Router ──┬──► handlers ──► views ──► templates
                        │
                        └──► /static/ ──► FileServer
```

Club Manager dispose maintenant d'une première interface publique configurable.

---

# Jalon 08 — PostgreSQL, psql et Goose

## Objectif

Introduire la persistance des données et commencer à versionner la structure de la base.

## Fiche de jalon

- [[Jalon 08 - PostgreSQL, psql et Goose]]
    

## PostgreSQL et `psql`

- [[01 - Installation et première prise en main de PostgresSQL|Installation et première prise en main de PostgreSQL]]
    
- [[02 - Finaliser linstallation et tester la connexion applicative|Finaliser l'installation et tester la connexion]]
    
- [[03 - Première manipulations SQL|Premières manipulations SQL]]
    

## Goose

- [[01 - Goose|Goose]]
    
- [[02 - Goose Comprendre une commande de migration PostgreSQL|Comprendre une commande Goose]]
    

## Architecture obtenue

```text
Application Go


PostgreSQL
    ▲
    │
  Goose
    ▲
    │
migrations/
```

La première migration introduit la table :

```text
members
```

Goose permet maintenant de faire évoluer la structure de PostgreSQL de manière :

- versionnée ;
    
- reproductible ;
    
- contrôlée.
    

## État actuel

```text
Serveur HTTP                ✅
Router                      ✅
Handlers                    ✅
Views                       ✅
Templates                   ✅
Configuration               ✅
Pages génériques            ✅
Interface publique          ✅
Tests                       ✅
PostgreSQL                  ✅
Migrations Goose            ✅

Connexion Go → PostgreSQL   ⏳
sqlc                        ⏳
Création d'un membre        ⏳
```

---

# Références générales

Certaines fiches ne correspondent pas à un jalon particulier.

Elles peuvent être utilisées pendant toute la durée du projet.

---

## Git général

- [[BRANCHE|Branches Git]]
    

Cette partie pourra accueillir les notions Git qui ne correspondent pas directement à une étape précise de Club Manager.

Par exemple, à terme :

```text
rebase
cherry-pick
stash
tags
worktree
```

uniquement lorsqu'elles seront réellement utiles.

---

## Vim

- [[00 - VIM COMMANDES|Commandes Vim]]
    
- [[01 - L'essentiel|L'essentiel]]
    
- [[02 - Les combinaisons|Combinaisons utiles]]
    

Ces fiches constituent une référence pratique pour l'édition du code et ne dépendent pas de l'architecture de Club Manager.

---

# Comment utiliser cette documentation ?

Il existe deux manières principales de naviguer dans les fiches.

## Revoir l'évolution du projet

Suivre les jalons dans l'ordre :

```text
00
↓
01
↓
02
↓
03
↓
04
↓
05
↓
06
↓
07
↓
08
```

Cela permet de comprendre :

> Pourquoi l'architecture actuelle existe-t-elle ?

Chaque abstraction peut ainsi être reliée au problème qui a provoqué son apparition.

---

## Rechercher une notion précise

Les liens entre les fiches permettent également de naviguer par concept.

Exemple :

```text
Interface
   │
   ├── interface composée
   │
   └── io.Writer
```

ou :

```text
Templates
   │
   ├── embed
   ├── ParseFS
   ├── passage de données
   └── layout commun
```

Une notion reste rangée dans son jalon d'origine même si elle est réutilisée plus tard.

---

# Règle de classement

Lorsqu'une nouvelle fiche apparaît, se poser d'abord la question :

> Dans quel jalon avons-nous eu besoin de comprendre cette notion pour la première fois ?

Si la réponse est claire :

```text
fiche
  ↓
jalon correspondant
```

Exemple :

```text
pgxpool
   ↓
prochain jalon consacré
à la connexion Go ↔ PostgreSQL
```

Si la notion n'est liée à aucun état particulier du projet et sert de référence générale :

```text
fiche
  ↓
référence générale
```

Exemple :

```text
commandes Vim
branches Git
```

---

# Philosophie de la documentation

Les fiches ne cherchent pas uniquement à documenter **comment** fonctionne Club Manager.

Elles cherchent aussi à conserver :

```text
Pourquoi ce choix ?
Pourquoi maintenant ?
Quel problème cherchait-on à résoudre ?
Quelle notion avons-nous apprise ?
```

L'ensemble de la documentation doit donc permettre de reconstruire le raisonnement :

```text
problème
   ↓
notion
   ↓
solution
   ↓
nouvelle architecture
```

---

# Comprendre et retenir

> **Les jalons racontent l'histoire du projet.**

---

> **Les fiches techniques expliquent les notions rencontrées pendant cette histoire.**

---

> **Une fiche appartient au jalon où la notion devient nécessaire pour la première fois.**

---

> **Une notion réutilisée plus tard ne change pas de jalon.**

On crée simplement des liens entre les fiches.

---

> **L'architecture de Club Manager n'est pas construite à l'avance.**

Elle évolue lorsque le projet rencontre un problème concret :

```text
besoin
  ↓
problème observé
  ↓
solution
  ↓
nouveau jalon
```

C'est cette évolution que cette documentation cherche à conserver.




