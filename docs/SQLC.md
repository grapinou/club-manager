
---


## Objectif

Comprendre le rôle de **sqlc** dans Club Manager et comment il s'intègre avec PostgreSQL, Goose et `pgxpool`.

À ce stade du projet, nous avons déjà :

```text
Goose
→ fait évoluer la structure PostgreSQL

pgxpool
→ permet à l'application Go d'accéder à PostgreSQL
```

Il manque encore un moyen pratique d'écrire nos requêtes SQL et de les utiliser depuis Go.

C'est le rôle de **sqlc**.

---

# Qu'est-ce que sqlc ?

`sqlc` est un générateur de code.

Nous écrivons nous-mêmes du SQL :

```sql
SELECT ...
INSERT ...
UPDATE ...
DELETE ...
```

Puis sqlc analyse :

- la structure de la base ;
    
- nos requêtes SQL ;
    

et génère du code Go typé.

Le principe général est :

```text
schéma PostgreSQL
        │
        │
        ▼
      sqlc ◄──── requêtes SQL
        │
        ▼
   code Go généré
```

---

# sqlc n'est pas un ORM

Avec sqlc, nous continuons à écrire explicitement le SQL.

Par exemple :

```sql
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

sqlc ne construit pas cette requête pour nous.

Il construit **le code Go permettant de l'utiliser proprement**.

On peut donc résumer :

```text
Nous
│
└── écrivons le SQL

sqlc
│
└── génère le Go correspondant
```

---

# Place de sqlc dans Club Manager

Nous avons maintenant trois outils complémentaires.

```text
                    PostgreSQL
                    ▲         ▲
                    │         │
                 pgxpool    Goose
                    ▲         ▲
                    │         │
              Club Manager  migrations/
                    ▲
                    │
                  sqlc
                    ▲
               ┌────┴────┐
               │         │
          migrations   queries
```

Leur responsabilité est différente.

## Goose

```text
migrations
    │
    ▼
  Goose
    │
    ▼
PostgreSQL
```

Goose **exécute** les migrations.

Il fait évoluer la structure réelle de PostgreSQL.

---

## pgxpool

```text
Club Manager
     │
     ▼
  pgxpool
     │
     ▼
PostgreSQL
```

`pgxpool` fournit à l'application un moyen de communiquer avec PostgreSQL.

---

## sqlc

```text
migrations
     │
     ▼
    sqlc ◄──── queries/*.sql
     │
     ▼
code Go généré
```

sqlc travaille pendant le développement.

Il ne remplace ni Goose ni pgx.

---

# Installation

sqlc a été installé avec Go :

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

L'exécutable se trouve dans :

```text
/home/sighto/go/bin/sqlc
```

Vérification :

```bash
which sqlc
```

puis :

```bash
sqlc version
```

Version utilisée lors de ce jalon :

```text
v1.31.1
```

---

# Le fichier `sqlc.yaml`

La configuration de sqlc se trouve à la racine du projet :

```text
sqlc.yaml
```

Configuration utilisée :

```yaml
version: "2"

sql:
  - engine: "postgresql"
    schema: "migrations"
    queries: "internal/database/queries"
    gen:
      go:
        package: "dbsqlc"
        out: "internal/database/dbsqlc"
        sql_package: "pgx/v5"
```

---

# `version`

```yaml
version: "2"
```

Il s'agit de la version du format du fichier de configuration sqlc.

Ce n'est ni :

- la version de PostgreSQL ;
    
- la version de Go ;
    
- la version de sqlc.
    

---

# `engine`

```yaml
engine: "postgresql"
```

On indique à sqlc que notre SQL utilise PostgreSQL.

---

# `schema`

```yaml
schema: "migrations"
```

sqlc doit connaître la structure des tables pour déterminer les types des colonnes.

Nous utilisons directement :

```text
migrations/
```

où Goose possède déjà nos migrations.

Par exemple :

```sql
CREATE TABLE members (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    first_name TEXT NOT NULL,
    last_name TEXT NOT NULL,
    birth_date DATE NOT NULL,
    email TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

La même migration possède donc deux utilisations :

```text
             migrations/
              /       \
             ▼         ▼
          Goose       sqlc
             │         │
             ▼         ▼
        modifie DB   comprend
                     le schéma
```

Il n'est pas nécessaire de recopier le schéma dans un autre fichier.

---

# `queries`

```yaml
queries: "internal/database/queries"
```

Ce dossier contient **le SQL que nous écrivons nous-mêmes**.

Organisation actuelle :

```text
internal/database/
├── database.go
├── queries/
│   └── members.sql
└── dbsqlc/
```

---

# Première requête : `CreateMember`

Dans :

```text
internal/database/queries/members.sql
```

nous avons écrit :

```sql
-- name: CreateMember :one
INSERT INTO members (
    first_name,
    last_name,
    birth_date,
    email
) VALUES (
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

---

# Annotation sqlc

Cette ligne :

```sql
-- name: CreateMember :one
```

est interprétée par sqlc.

Elle contient deux informations importantes.

## `CreateMember`

```text
CreateMember
```

sera utilisé pour créer une méthode Go :

```go
CreateMember(...)
```

---

## `:one`

```text
:one
```

indique que la requête retourne exactement une ligne.

Ici :

```sql
INSERT ...
RETURNING ...
```

crée un membre puis retourne le membre créé.

---

# Les paramètres `$1`, `$2`, `$3`, `$4`

PostgreSQL utilise :

```sql
$1
$2
$3
$4
```

comme paramètres de la requête.

Ici :

```text
$1 → first_name
$2 → last_name
$3 → birth_date
$4 → email
```

Nous ne fournissons pas :

```text
id
created_at
```

car PostgreSQL les génère lui-même.

```text
id
→ IDENTITY

created_at
→ DEFAULT NOW()
```

---

# `sqlc generate`

Une fois le schéma, la configuration et les requêtes prêts :

```bash
sqlc generate
```

sqlc analyse :

```text
migrations/
      +
queries/
      +
sqlc.yaml
```

puis génère du Go.

---

# Code généré

Dans notre projet :

```text
internal/database/dbsqlc/
├── db.go
├── members.sql.go
└── models.go
```

Ces fichiers sont générés automatiquement.

> Ils ne doivent pas être modifiés manuellement.

Une modification doit être faite dans :

```text
migrations/
```

ou :

```text
queries/
```

puis :

```bash
sqlc generate
```

---

# `models.go`

sqlc a transformé notre table PostgreSQL en structure Go.

Conceptuellement :

```go
type Member struct {
    ID        int32
    FirstName string
    LastName  string
    BirthDate pgtype.Date
    Email     pgtype.Text
    CreatedAt pgtype.Timestamptz
}
```

On voit ici la correspondance entre PostgreSQL et Go :

```text
PostgreSQL             Go

INTEGER          →     int32
TEXT NOT NULL    →     string
DATE             →     pgtype.Date
TEXT nullable    →     pgtype.Text
TIMESTAMPTZ      →     pgtype.Timestamptz
```

---

# Pourquoi `Email` n'est-il pas un `string` ?

Dans notre migration :

```sql
email TEXT
```

contrairement à :

```sql
first_name TEXT NOT NULL
```

`email` accepte :

```text
une chaîne
OU
NULL
```

Un simple :

```go
string
```

ne permet pas de représenter cette distinction.

sqlc génère donc :

```go
pgtype.Text
```

qui sait représenter une valeur PostgreSQL nullable.

---

# `CreateMemberParams`

sqlc a également généré une structure pour les arguments nécessaires à la requête :

```go
type CreateMemberParams struct {
    FirstName string
    LastName  string
    BirthDate pgtype.Date
    Email     pgtype.Text
}
```

Cela correspond directement à :

```text
$1
$2
$3
$4
```

de notre requête SQL.

---

# La méthode `CreateMember`

sqlc génère ensuite une méthode ressemblant à :

```go
func (q *Queries) CreateMember(
    ctx context.Context,
    arg CreateMemberParams,
) (Member, error)
```

On pourra donc écrire :

```go
member, err := queries.CreateMember(
    ctx,
    params,
)
```

au lieu d'écrire manuellement :

```text
QueryRow
+
paramètres
+
Scan
+
conversion des types
```

---

# Ce que sqlc génère à notre place

Sans sqlc, nous aurions dû écrire quelque chose ressemblant à :

```go
row := db.QueryRow(
    ctx,
    query,
    firstName,
    lastName,
    birthDate,
    email,
)

var member Member

err := row.Scan(
    &member.ID,
    &member.FirstName,
    &member.LastName,
    &member.BirthDate,
    &member.Email,
    &member.CreatedAt,
)
```

sqlc produit ce code automatiquement.

Nous pouvons donc nous concentrer sur :

```text
le SQL
+
la logique de l'application
```

---

# `db.go`

Ce fichier contient notamment une interface :

```go
type DBTX interface {
    Exec(...)
    Query(...)
    QueryRow(...)
}
```

L'idée est très importante.

sqlc ne demande pas :

```text
obligatoirement un *pgxpool.Pool
```

Il demande :

> Donne-moi quelque chose qui sait exécuter `Exec`, `Query` et `QueryRow`.

---

# L'interface `DBTX`

Nous avions déjà étudié les interfaces Go.

Ici, nous rencontrons leur utilisation dans un cas réel.

Notre :

```go
*pgxpool.Pool
```

possède les méthodes nécessaires à `DBTX`.

Il satisfait donc automatiquement cette interface.

```text
*pgxpool.Pool
      │
      ├── Exec
      ├── Query
      └── QueryRow
      │
      ▼
     DBTX
```

Il n'y a aucun :

```go
implements DBTX
```

à écrire.

C'est le principe des interfaces implicites de Go.

---

# `dbsqlc.New`

sqlc génère également :

```go
func New(db DBTX) *Queries {
    return &Queries{
        db: db,
    }
}
```

Notre pool peut donc être passé directement :

```go
queries := dbsqlc.New(db)
```

Le chemin devient :

```text
database.New(...)
      │
      ▼
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
CreateMember(...)
```

---

# `db` et `queries`

Il faut bien distinguer les deux.

## `db`

```go
db
```

est notre moyen technique d'accéder à PostgreSQL.

Son type réel est :

```go
*pgxpool.Pool
```

---

## `queries`

```go
queries
```

est l'objet généré par sqlc permettant d'utiliser nos requêtes SQL.

Son type est :

```go
*dbsqlc.Queries
```

Conceptuellement :

```text
queries
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

# Architecture obtenue

```text
                    Club Manager
                         │
                         ▼
                      queries
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

Et pendant le développement :

```text
migrations ──────┐
                 │
                 ▼
                sqlc
                 ▲
                 │
queries/*.sql ───┘
                 │
                 ▼
          code Go généré
```

---

# Workflow sqlc

Lorsqu'une nouvelle requête est nécessaire :

```text
1. écrire le SQL
       │
       ▼
2. ajouter l'annotation sqlc
       │
       ▼
3. sqlc generate
       │
       ▼
4. utiliser la méthode Go générée
```

Exemple :

```sql
-- name: GetMember :one
SELECT ...
```

puis :

```bash
sqlc generate
```

pour obtenir conceptuellement :

```go
queries.GetMember(...)
```

---

# Comprendre et retenir

> **sqlc transforme du SQL écrit par nous en code Go typé.**

---

> **sqlc n'est pas un ORM.**

Nous gardons la maîtrise du SQL.

---

> **Les migrations décrivent le schéma.**

```text
migrations/
   │
   ├── Goose → exécute
   │
   └── sqlc  → analyse
```

---

> **Les requêtes SQL sont écrites dans `queries/`.**

```text
queries/
└── members.sql
```

---

> **Le code généré vit dans `dbsqlc/`.**

```text
dbsqlc/
├── db.go
├── models.go
└── members.sql.go
```

---

> **On ne modifie pas directement les fichiers générés.**

On modifie :

```text
SQL
```

puis :

```bash
sqlc generate
```

---

> **`db` et `queries` ont deux rôles différents.**

```text
db
→ accès technique à PostgreSQL

queries
→ requêtes SQL typées
```

---

> **`*pgxpool.Pool` satisfait l'interface `DBTX`.**

C'est un cas concret d'interface implicite en Go.

---

# Phrase à retenir

> **Nous écrivons le SQL, sqlc écrit le code Go répétitif qui permet de l'utiliser.**

La chaîne complète devient :

```text
SQL
 │
 ▼
sqlc
 │
 ▼
Go typé
 │
 ▼
pgxpool
 │
 ▼
PostgreSQL
```