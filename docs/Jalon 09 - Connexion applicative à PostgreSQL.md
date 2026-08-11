
---


## Objectif du jalon

Permettre pour la première fois à **Club Manager lui-même** de communiquer avec PostgreSQL.

Le jalon précédent avait permis de créer et faire évoluer la base :

```text
migrations
    │
    ▼
  Goose
    │
    ▼
PostgreSQL
```

Mais l'application Go ne possédait encore aucun accès à cette base.

L'objectif de ce jalon est donc de construire :

```text
Club Manager
     │
     ▼
PostgreSQL
```

---

# État avant le jalon

À la fin du Jalon 08 :

```text
Club Manager


PostgreSQL
    ▲
    │
  Goose
    ▲
    │
migrations/
```

La base existe.

La table :

```text
members
```

existe.

Goose peut modifier sa structure.

Mais Club Manager ne peut encore ni :

```text
SELECT
INSERT
UPDATE
DELETE
```

des données PostgreSQL.

---

# Nouvelle responsabilité

Une nouvelle couche apparaît :

```text
internal/database/
```

Son rôle est de gérer l'accès technique de l'application à PostgreSQL.

Organisation :

```text
internal/
├── config/
├── database/
├── handlers/
├── router/
└── views/
```

La base de données devient donc une responsabilité clairement identifiée dans l'architecture.

---

# pgx

Nous avons choisi :

```text
pgx
```

pour permettre à Go de communiquer avec PostgreSQL.

Plus précisément :

```text
pgxpool
```

est utilisé pour gérer un pool de connexions.

---

# Pourquoi un pool ?

Club Manager est un serveur HTTP.

Plusieurs utilisateurs pourront déclencher des traitements simultanément.

```text
Utilisateur A ──┐
Utilisateur B ──┼──► Club Manager
Utilisateur C ──┘
```

Ces traitements peuvent tous avoir besoin de PostgreSQL.

`pgxpool` fournit donc :

```text
                 pgxpool
               /    |    \
              ▼     ▼     ▼
         connexion connexion connexion
               \    |    /
                    ▼
               PostgreSQL
```

---

# Première connexion

Nous avons commencé par tester directement la connexion dans `main`.

Le principe :

```go
ctx := context.Background()

db, err := pgxpool.New(
    ctx,
    os.Getenv("DATABASE_URL"),
)
```

puis :

```go
db.Ping(ctx)
```

Le premier test a échoué avec :

```text
password authentication failed
```

Cette erreur a montré que :

```text
PostgreSQL trouvé             ✅
localhost:5432 accessible     ✅
utilisateur identifié         ✅
base identifiée               ✅

mot de passe                  ❌
```

Le mot de passe a ensuite été fourni temporairement dans la session avec :

```bash
export PGPASSWORD='...'
```

Après correction :

```text
Connexion à PostgreSQL établie
Serveur lancé sur http://localhost:8080
```

Club Manager communiquait pour la première fois réellement avec PostgreSQL.

---

# Variables d'environnement

La connexion utilise :

```text
DATABASE_URL
```

avec par exemple :

```text
host=localhost
port=5432
user=club_manager
dbname=club_manager
sslmode=require
```

et séparément :

```text
PGPASSWORD
```

pour le secret.

Nous distinguons ainsi :

```text
config.json
→ configuration fonctionnelle du club

variables d'environnement
→ configuration technique du serveur
```

Le mot de passe n'est donc pas placé dans Git.

---

# `context.Context`

La connexion avec pgx nous a également fait rencontrer :

```go
context.Context
```

Nous avons commencé avec :

```go
ctx := context.Background()
```

L'idée essentielle retenue est :

> Un contexte accompagne un travail et transporte sa durée de vie à travers les différentes couches.

Dans ce jalon :

```text
main
 │
 ▼
Context
 │
 ▼
pgx
 │
 ▼
PostgreSQL
```

Plus tard, pour les requêtes HTTP :

```text
HTTP
 │
 ▼
r.Context()
 │
 ▼
Handler
 │
 ▼
SQL
 │
 ▼
PostgreSQL
```

---

# Extraction de la responsabilité `database`

Une fois la connexion comprise et validée, sa gestion a été extraite de `main`.

Nous avons créé :

```text
internal/database/database.go
```

avec une fonction conceptuellement équivalente à :

```go
func New(
    ctx context.Context,
    databaseURL string,
) (*pgxpool.Pool, error)
```

Son rôle :

```text
database.New
     │
     ├── crée pgxpool
     ├── teste PostgreSQL
     │
     └── retourne le pool
```

---

# `main` reste le point d'assemblage

Le rôle de `main` devient :

```text
main
│
├── config.Load()
│
├── context.Background()
│
├── database.New()
│
├── router.New()
└── ListenAndServe()
```

Il ne connaît plus les détails de création du pool.

Il assemble simplement les différentes responsabilités de l'application.

---

# `db`

La fonction :

```go
database.New(...)
```

retourne :

```go
*pgxpool.Pool
```

stocké dans :

```go
db
```

Conceptuellement :

```text
db
 │
 ▼
pgxpool
 │
 ▼
PostgreSQL
```

`db` n'est donc pas PostgreSQL lui-même.

C'est l'objet par lequel Club Manager peut obtenir des connexions vers PostgreSQL.

---

# sqlc

Une fois la connexion disponible, une nouvelle question est apparue :

> Comment utiliser nos requêtes SQL depuis Go sans écrire manuellement tout le code répétitif ?

Nous avons introduit :

```text
sqlc
```

---

# Installation de sqlc

sqlc a été réinstallé via Go :

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

L'exécutable utilisé est :

```text
/home/sighto/go/bin/sqlc
```

Version lors du jalon :

```text
v1.31.1
```

---

# Configuration de sqlc

Un fichier :

```text
sqlc.yaml
```

est ajouté à la racine.

Il relie :

```text
migrations/
```

à :

```text
internal/database/queries/
```

et génère du code vers :

```text
internal/database/dbsqlc/
```

La chaîne est :

```text
migrations
     │
     ▼
    sqlc
     ▲
     │
queries/*.sql
     │
     ▼
dbsqlc/*.go
```

---

# Première requête SQL applicative

Nous avons créé :

```text
internal/database/queries/members.sql
```

avec une première requête :

```sql
-- name: CreateMember :one
INSERT INTO members (
    first_name,
    last_name,
    birth_date,
    email
)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING
    id,
    first_name,
    last_name,
    birth_date,
    email,
    created_at;
```

Cette requête décrit explicitement comment créer un membre dans PostgreSQL.

---

# Génération du code

Commande :

```bash
sqlc generate
```

sqlc produit :

```text
internal/database/dbsqlc/
├── db.go
├── members.sql.go
└── models.go
```

Ces fichiers sont automatiquement générés.

---

# Modèle `Member`

sqlc utilise notre migration pour produire une représentation Go de :

```text
members
```

Le modèle contient conceptuellement :

```text
Member
│
├── ID
├── FirstName
├── LastName
├── BirthDate
├── Email
└── CreatedAt
```

La structure PostgreSQL devient donc directement exploitable depuis Go.

---

# `CreateMember`

Notre annotation :

```sql
-- name: CreateMember :one
```

a permis de générer une méthode Go :

```text
CreateMember(...)
```

Le chemin devient :

```text
SQL écrit par nous
       │
       ▼
      sqlc
       │
       ▼
méthode Go typée
```

---

# Première interface réelle : `DBTX`

sqlc génère une interface :

```text
DBTX
```

qui demande notamment :

```text
Exec
Query
QueryRow
```

Notre :

```go
*pgxpool.Pool
```

possède ces méthodes.

Il satisfait donc automatiquement cette interface.

```text
*pgxpool.Pool
      │
      ▼
     DBTX
```

Cela constitue un cas concret d'utilisation des interfaces implicites de Go.

---

# `dbsqlc.New(db)`

Le code généré par sqlc possède également :

```go
dbsqlc.New(db)
```

Conceptuellement :

```text
db
│
│ *pgxpool.Pool
▼
dbsqlc.New(db)
│
▼
queries
│
│ *dbsqlc.Queries
▼
méthodes SQL générées
```

Nous obtenons donc deux objets importants.

## `db`

```text
db
→ accès technique à PostgreSQL
```

## `queries`

```text
queries
→ accès aux requêtes SQL générées
```

---

# Architecture à la fin du Jalon 09

La chaîne technique est maintenant :

```text
                  Club Manager
                       │
                       ▼
                    dbsqlc
                       │
                       ▼
                    pgxpool
                       │
                       ▼
                  PostgreSQL
```

Pendant le développement :

```text
               migrations/
                /       \
               ▼         ▼
            Goose       sqlc
               │         ▲
               ▼         │
          PostgreSQL  queries/
                         │
                         ▼
                    Go généré
```

En réunissant les deux :

```text
queries/*.sql
      │
      ▼
     sqlc
      │
      ▼
   dbsqlc
      │
      ▼
     db
      │
      ▼
  pgxpool
      │
      ▼
PostgreSQL
```

---

# Évolution architecturale

Avant le Jalon 08 :

```text
Club Manager
→ application sans base persistante
```

Jalon 08 :

```text
migrations
    │
    ▼
PostgreSQL
```

Jalon 09 :

```text
Club Manager
     │
     ▼
PostgreSQL
```

La base de données n'est donc plus seulement une infrastructure extérieure au projet.

Elle devient une dépendance utilisée par l'application.

---

# Ce que sait faire le projet à ce stade

```text
Serveur HTTP                        ✅
Router                              ✅
Handlers                            ✅
Views                               ✅
Templates                           ✅
Configuration                       ✅
Interface publique                  ✅
PostgreSQL                          ✅
Migrations Goose                    ✅
Connexion Go → PostgreSQL           ✅
Pool pgx                            ✅
Package database                    ✅
sqlc                                ✅
Requête CreateMember définie        ✅
Code Go SQL généré                  ✅

Création via requête HTTP           ❌
Handler membre                      ❌
POST                                ❌
Validation des données utilisateur  ❌
```

Cette frontière est volontaire.

---

# Pourquoi arrêter le jalon ici ?

Le problème résolu dans ce jalon était :

> Comment Club Manager peut-il communiquer proprement avec PostgreSQL ?

La réponse est maintenant :

```text
database.New
      │
      ▼
   pgxpool
      │
      ▼
 PostgreSQL
```

avec :

```text
SQL
 │
 ▼
sqlc
 │
 ▼
dbsqlc
```

La prochaine question est différente :

> Comment une requête HTTP déclenche-t-elle une opération métier dans PostgreSQL ?

Elle introduira :

```text
HTTP POST
    │
    ▼
Handler
    │
    ▼
validation
    │
    ▼
r.Context()
    │
    ▼
CreateMember
    │
    ▼
PostgreSQL
```

Cette nouvelle chaîne mérite donc son propre jalon.

---

# Prochain jalon

## Jalon 10 — Premier flux HTTP vers PostgreSQL

Objectif envisagé :

```text
Navigateur
    │
    │ POST
    ▼
Router
    │
    ▼
Handler
    │
    ▼
CreateMember
    │
    ▼
PostgreSQL
```

Ce jalon introduira probablement :

- une route `POST` ;
    
- un handler dédié ;
    
- la réception de données utilisateur ;
    
- la validation ;
    
- `r.Context()` ;
    
- l'utilisation réelle de `queries.CreateMember()` ;
    
- la gestion des erreurs PostgreSQL ;
    
- les tests correspondants.
    

---

# Comprendre et retenir

> **Le Jalon 08 a créé la base et son système de migrations.**

```text
Goose → PostgreSQL
```

---

> **Le Jalon 09 connecte l'application à cette base.**

```text
Go → pgxpool → PostgreSQL
```

---

> **sqlc relie notre SQL à notre code Go.**

```text
SQL → sqlc → Go typé
```

---

> **`database` gère la connexion.**

```text
database.New()
→ *pgxpool.Pool
```

---

> **`dbsqlc` contient le code généré à partir de nos requêtes.**

```text
dbsqlc.New(db)
→ *dbsqlc.Queries
```

---

> **L'application possède maintenant tout le chemin technique nécessaire pour atteindre PostgreSQL.**

Mais elle ne possède pas encore de fonctionnalité HTTP utilisant ce chemin.

C'est précisément le rôle du prochain jalon.

---

# Résumé du Jalon 09

```text
         migrations
          /      \
         ▼        ▼
      Goose      sqlc
         │        ▲
         │        │
         │     queries
         │        │
         │        ▼
         │     dbsqlc
         │        │
         ▼        ▼
        PostgreSQL
             ▲
             │
          pgxpool
             ▲
             │
             db
             ▲
             │
       Club Manager
```

**Club Manager est désormais techniquement connecté à sa base PostgreSQL et dispose de requêtes SQL utilisables depuis Go.**