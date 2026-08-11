
---


# Go — PostgreSQL avec `pgx` et `pgxpool`

## Objectif

Comprendre comment une application Go peut communiquer avec PostgreSQL grâce à :

```text
pgx
```

et plus précisément :

```text
pgxpool
```

Dans Club Manager, cette étape marque un changement important.

Jusqu'ici :

```text
Application Go

PostgreSQL
    ▲
    │
  Goose
```

Goose savait communiquer avec PostgreSQL pour modifier la structure de la base.

Mais l'application Club Manager elle-même ne communiquait pas encore avec PostgreSQL.

Nous voulons maintenant obtenir :

```text
Application Go
      │
      ▼
   pgxpool
      │
      ▼
 PostgreSQL
```

---

# `pgx`

`pgx` est une bibliothèque Go permettant de communiquer avec PostgreSQL.

Le module utilisé est :

```go
github.com/jackc/pgx/v5
```

Il fournit notamment les outils nécessaires pour :

- ouvrir une connexion PostgreSQL ;
    
- envoyer des requêtes SQL ;
    
- récupérer des résultats ;
    
- gérer des transactions ;
    
- utiliser les types PostgreSQL ;
    
- gérer plusieurs connexions avec `pgxpool`.
    

---

# `pgx` et `pgxpool`

Il est important de distinguer deux notions.

## `pgx.Conn`

Une connexion simple peut être représentée par :

```text
Application Go
      │
      ▼
 connexion
      │
      ▼
 PostgreSQL
```

Conceptuellement :

```go
*pgx.Conn
```

représente une connexion unique.

Pour une petite commande ou un programme exécuté ponctuellement, cela peut être suffisant.

---

## `pgxpool.Pool`

Une application web doit pouvoir recevoir plusieurs requêtes.

Par exemple :

```text
Utilisateur A ──┐
Utilisateur B ──┼──► Club Manager
Utilisateur C ──┘
```

Ces requêtes peuvent avoir besoin de PostgreSQL au même moment.

Un pool permet de disposer de plusieurs connexions réutilisables :

```text
              Club Manager
                   │
                   ▼
                pgxpool
              /    |    \
             ▼     ▼     ▼
          connexion connexion connexion
              \    |    /
                   ▼
              PostgreSQL
```

`pgxpool` gère ce mécanisme.

---

# Pourquoi utiliser `pgxpool` dans Club Manager ?

Club Manager est un serveur HTTP.

Plusieurs requêtes peuvent être traitées simultanément.

Utiliser :

```go
pgxpool.Pool
```

est donc plus adapté qu'une connexion PostgreSQL unique.

Le pool peut :

- créer des connexions ;
    
- les réutiliser ;
    
- les distribuer aux traitements qui en ont besoin ;
    
- les récupérer après utilisation ;
    
- fermer les ressources lorsque l'application s'arrête.
    

---

# Installation

Nous avons ajouté `pgxpool` avec :

```bash
go get github.com/jackc/pgx/v5/pgxpool
```

Cette commande a modifié :

```text
go.mod
go.sum
```

---

# `go.mod`

Après l'installation, plusieurs dépendances sont apparues.

Par exemple :

```go
github.com/jackc/pgx/v5
github.com/jackc/puddle/v2
github.com/jackc/pgpassfile
github.com/jackc/pgservicefile
```

Nous n'avons pourtant demandé directement que :

```text
pgxpool
```

C'est normal.

Une bibliothèque peut elle-même dépendre d'autres bibliothèques.

On obtient donc une chaîne comme :

```text
Club Manager
     │
     ▼
   pgxpool
     │
     ├── pgx
     ├── puddle
     └── autres dépendances
```

Go gère automatiquement cet arbre de dépendances.

---

# `go.sum`

Un nouveau fichier est également apparu :

```text
go.sum
```

Il contient des empreintes permettant à Go de vérifier les modules téléchargés.

On peut simplifier ainsi :

```text
go.mod
│
├── quel module ?
└── quelle version ?

go.sum
│
└── ce module correspond-il bien
    au contenu attendu ?
```

`go.mod` et `go.sum` doivent être versionnés avec Git.

---

# Première connexion depuis Club Manager

Nous avons choisi de commencer simplement.

Avant de créer un package spécifique à la base de données, nous avons ajouté la connexion directement dans :

```text
cmd/server/main.go
```

L'objectif était de comprendre le mécanisme avant de créer une abstraction.

---

# Imports nécessaires

Nous avons ajouté :

```go
import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/grapinou/club-manager/internal/config"
	"github.com/grapinou/club-manager/internal/router"
)
```

Les nouveautés sont :

```go
"context"
"os"
"github.com/jackc/pgx/v5/pgxpool"
```

---

# Le `context`

Nous créons :

```go
ctx := context.Background()
```

Pour l'instant, on peut considérer :

```text
context.Background()
```

comme le contexte racine de l'application.

Il est transmis aux opérations qui peuvent notamment :

- être annulées ;
    
- posséder une durée maximale ;
    
- dépendre de la durée de vie d'une requête.
    

Nous approfondirons `context.Context` dans une fiche dédiée.

Pour cette première connexion :

```go
ctx := context.Background()
```

nous fournit simplement le contexte nécessaire à `pgxpool`.

---

# Créer le pool

Nous avons ensuite écrit :

```go
db, err := pgxpool.New(
	ctx,
	os.Getenv("DATABASE_URL"),
)
```

Décomposons cette instruction.

---

## `pgxpool.New`

```go
pgxpool.New(...)
```

crée un nouveau pool PostgreSQL.

Le résultat est stocké dans :

```go
db
```

`db` représente donc notre pool de connexions.

Conceptuellement :

```text
db
│
└── pool PostgreSQL
```

---

# `DATABASE_URL`

Le deuxième argument est :

```go
os.Getenv("DATABASE_URL")
```

`os.Getenv` permet de lire une variable d'environnement.

Nous avons défini :

```bash
export DATABASE_URL="host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require"
```

La variable contient les informations nécessaires pour trouver PostgreSQL.

---

# Décomposition de la connexion

```text
host=localhost
```

PostgreSQL se trouve sur la machine locale.

---

```text
port=5432
```

Le serveur PostgreSQL écoute sur le port `5432`.

---

```text
user=club_manager
```

L'application se connecte avec l'utilisateur PostgreSQL :

```text
club_manager
```

---

```text
dbname=club_manager
```

La base utilisée est :

```text
club_manager
```

---

```text
sslmode=require
```

La connexion demande l'utilisation de SSL.

---

# Séparer la connexion du mot de passe

Nous n'avons pas placé le mot de passe directement dans :

```text
DATABASE_URL
```

Nous avons utilisé une autre variable d'environnement :

```bash
export PGPASSWORD='mot_de_passe'
```

La séparation devient :

```text
DATABASE_URL
├── serveur
├── port
├── utilisateur
├── base
└── paramètres

PGPASSWORD
└── secret
```

C'est particulièrement intéressant pour éviter d'écrire un mot de passe dans :

```text
main.go
config.json
Git
```

---

# Pourquoi ne pas utiliser `config.json` ?

Le fichier :

```text
config/config.json
```

contient les informations fonctionnelles propres à l'association.

Par exemple :

```text
nom du site
titre des pages
description
adresse
horaires
images
```

La connexion PostgreSQL appartient plutôt à l'environnement dans lequel l'application est exécutée.

On distingue donc :

```text
config.json
      │
      ▼
configuration fonctionnelle
de Club Manager
```

et :

```text
DATABASE_URL
PGPASSWORD
      │
      ▼
configuration technique
de l'environnement
```

Cette séparation évite également de placer des secrets dans le dépôt Git.

---

# Les variables sont temporaires

Lorsque nous faisons :

```bash
export PGPASSWORD='...'
```

la variable existe dans la session du shell courant.

Si le terminal est fermé, la variable disparaît.

Même principe pour :

```bash
export DATABASE_URL="..."
```

C'est pourquoi, lors de notre premier test, PostgreSQL a refusé l'authentification : le mot de passe n'avait pas été redéfini dans la session.

---

# Première erreur obtenue

Lors du premier lancement :

```bash
go run ./cmd/server
```

nous avons obtenu :

```text
failed to connect to `user=club_manager database=club_manager`:
127.0.0.1:5432 (localhost):
failed SASL auth:
FATAL: password authentication failed for user "club_manager"
(SQLSTATE 28P01)
```

Cette erreur était intéressante.

Elle indiquait que plusieurs choses fonctionnaient déjà.

```text
PostgreSQL trouvé             ✅
localhost:5432 accessible     ✅
utilisateur identifié         ✅
base identifiée               ✅

authentification              ❌
```

Le problème était donc très ciblé.

---

# Correction

Nous avons redéfini le mot de passe dans la session :

```bash
export PGPASSWORD='...'
```

Puis relancé :

```bash
go run ./cmd/server
```

Cette fois :

```text
Connexion à PostgreSQL établie

Serveur lancé sur http://localhost:8080
```

La connexion entre Club Manager et PostgreSQL fonctionne donc réellement.

---

# `pgxpool.New` ne suffit pas

Un point important est que :

```go
pgxpool.New(...)
```

crée le pool, mais ne constitue pas à lui seul notre preuve que PostgreSQL est accessible.

On peut représenter cela ainsi :

```text
pgxpool.New()
      │
      ▼
configuration du pool
```

Ce n'est pas encore nécessairement :

```text
PostgreSQL répond réellement
```

Pour vérifier la connexion, nous avons utilisé :

```go
db.Ping(ctx)
```

---

# Tester la connexion avec `Ping`

Le code est :

```go
if err := db.Ping(ctx); err != nil {
	log.Fatalf(
		"impossible de se connecter à PostgreSQL : %v",
		err,
	)
}
```

`Ping` demande réellement au pool d'effectuer une opération de connexion.

Le chemin devient :

```text
db.Ping(ctx)
     │
     ▼
pgxpool
     │
     ▼
tentative de connexion
     │
     ▼
PostgreSQL
```

Si PostgreSQL répond correctement :

```text
Ping réussi
```

Sinon :

```text
erreur
```

---

# Pourquoi notre erreur de mot de passe était utile ?

L'erreur :

```text
password authentication failed
```

nous a montré que :

```go
db.Ping(ctx)
```

effectuait bien un vrai test.

Sans `Ping`, le programme aurait pu créer son pool sans révéler immédiatement le problème d'authentification.

---

# Fermer le pool

Après la création du pool, nous avons ajouté :

```go
defer db.Close()
```

---

## `Close`

```go
db.Close()
```

ferme les connexions détenues par le pool.

Cela permet de libérer proprement les ressources utilisées par PostgreSQL.

---

# Pourquoi `defer` ?

Avec :

```go
defer db.Close()
```

nous demandons :

> Lorsque `main` se terminera, appelle `db.Close()`.

Le fonctionnement conceptuel est :

```text
création du pool
      │
      ▼
defer db.Close()
      │
      ▼
application fonctionne
      │
      ▼
main se termine
      │
      ▼
db.Close()
```

Le pool reste donc disponible pendant toute la durée de vie du serveur.

---

# Code obtenu

Notre première connexion ressemble à :

```go
ctx := context.Background()

db, err := pgxpool.New(
	ctx,
	os.Getenv("DATABASE_URL"),
)
if err != nil {
	log.Fatalf(
		"impossible de créer le pool PostgreSQL : %v",
		err,
	)
}
defer db.Close()

if err := db.Ping(ctx); err != nil {
	log.Fatalf(
		"impossible de se connecter à PostgreSQL : %v",
		err,
	)
}

log.Println("Connexion à PostgreSQL établie")
```

---

# Décomposition complète

```text
context.Background()
        │
        ▼
        ctx
        │
        ▼
pgxpool.New(
    ctx,
    DATABASE_URL
)
        │
        ▼
       db
        │
        ├── defer db.Close()
        │
        ▼
   db.Ping(ctx)
        │
        ▼
   PostgreSQL
        │
        ▼
Connexion valide
```

---

# Place dans `main`

Le rôle actuel de `main` devient :

```text
main
│
├── charger config.json
│
├── créer le contexte
│
├── créer pgxpool
│
├── vérifier PostgreSQL
│
├── construire le routeur
│
└── démarrer le serveur HTTP
```

Conceptuellement :

```text
                main
                 │
        ┌────────┼────────┐
        │        │        │
        ▼        ▼        ▼
     Config   PostgreSQL  Router
                 │
                 ▼
               pgxpool
```

---

# Pourquoi commencer dans `main` ?

Nous aurions pu créer directement :

```text
internal/database/
```

Mais nous avons volontairement commencé avec le code dans `main`.

L'objectif était d'abord de comprendre :

- ce qu'est `pgxpool` ;
    
- comment créer un pool ;
    
- comment transmettre une configuration ;
    
- comment tester la connexion ;
    
- comment fermer le pool ;
    
- comment gérer les erreurs.
    

Nous avons maintenant une vraie duplication potentielle de responsabilité dans `main`.

Cela nous donne une raison concrète pour créer plus tard :

```text
internal/database/
```

L'abstraction arrivera donc parce qu'un besoin est apparu.

---

# Goose et pgx ne font pas la même chose

Il est important de ne pas confondre les deux.

## Goose

Goose sert principalement à modifier la structure de la base.

```text
migrations SQL
      │
      ▼
    Goose
      │
      ▼
 PostgreSQL
```

Exemple :

```sql
CREATE TABLE members (...);
```

---

## pgx

`pgx` sert à l'application pendant son fonctionnement.

```text
Club Manager
     │
     ▼
  pgxpool
     │
     ▼
PostgreSQL
```

Plus tard, Club Manager pourra par exemple demander :

```sql
SELECT ...
INSERT ...
UPDATE ...
DELETE ...
```

---

# Goose et pgx travaillent donc en parallèle

```text
                    PostgreSQL
                    ▲        ▲
                    │        │
                  pgx      Goose
                    ▲
                    │
              Club Manager
```

Ils ne communiquent pas directement entre eux.

Tous les deux communiquent avec PostgreSQL.

---

# Et sqlc ?

Une autre étape viendra plus tard :

```text
sqlc
```

Son rôle sera encore différent.

Nous écrirons du SQL :

```sql
INSERT INTO members ...
```

Puis sqlc pourra générer du code Go typé.

La chaîne future ressemblera à :

```text
SQL
 │
 ▼
sqlc
 │
 ▼
code Go généré
 │
 ▼
pgxpool
 │
 ▼
PostgreSQL
```

Mais pour l'instant :

```text
pgxpool ✅

sqlc    plus tard
```

---

# Architecture à la fin de cette étape

Avant :

```text
Club Manager


PostgreSQL
    ▲
    │
  Goose
```

Maintenant :

```text
             Club Manager
                  │
                  ▼
               pgxpool
                  │
                  ▼
              PostgreSQL
                  ▲
                  │
                Goose
                  ▲
                  │
             migrations/
```

Club Manager possède maintenant deux chemins différents vers PostgreSQL :

```text
développement / évolution du schéma
           ↓
         Goose
```

et :

```text
fonctionnement de l'application
           ↓
        pgxpool
```

---

# Prochaine évolution logique

Le code fonctionne.

Mais `main` possède maintenant une nouvelle responsabilité :

```text
connexion à PostgreSQL
```

On pourra donc envisager :

```text
internal/database/
```

avec quelque chose comme :

```text
main
 │
 ├── config.Load()
 │
 ├── database.New()
 │
 ├── router.New()
 │
 └── ListenAndServe()
```

Le prochain refactoring pourra ainsi simplifier `main`.

---

# Comprendre et retenir

> **`pgx` permet au code Go de communiquer avec PostgreSQL.**

---

> **`pgxpool` gère un ensemble de connexions PostgreSQL réutilisables.**

Pour un serveur web :

```text
plusieurs requêtes
       ↓
     pgxpool
       ↓
plusieurs connexions
       ↓
   PostgreSQL
```

---

> **`pgxpool.New` crée le pool.**

```go
db, err := pgxpool.New(...)
```

---

> **`db.Ping(ctx)` vérifie réellement que PostgreSQL répond.**

```text
New
 ↓
créer le pool

Ping
 ↓
tester PostgreSQL
```

---

> **`defer db.Close()` garantit que le pool sera fermé lorsque l'application se termine.**

---

> **Les informations techniques de connexion sont placées dans l'environnement.**

```text
DATABASE_URL
PGPASSWORD
```

et non dans :

```text
config.json
```

---

> **`config.json` et les variables d'environnement n'ont pas la même responsabilité.**

```text
config.json
→ données fonctionnelles du club

variables d'environnement
→ configuration technique du serveur
```

---

> **Goose et pgx ont des responsabilités différentes.**

```text
Goose
→ structure de PostgreSQL

pgx
→ communication de l'application avec PostgreSQL
```

---

# Résumé

```text
DATABASE_URL
PGPASSWORD
      │
      ▼
  pgxpool.New
      │
      ▼
    Pool
      │
      ├── Ping
      │     │
      │     ▼
      │ PostgreSQL
      │
      └── Close
```

Et dans Club Manager :

```text
Navigateur
    │
    ▼
Club Manager
    │
    ├── HTTP
    │
    └── pgxpool
          │
          ▼
      PostgreSQL
```

**Cette étape marque la première connexion réelle entre l'application Go Club Manager et sa base PostgreSQL.**

