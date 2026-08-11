
---


## Objectif du jalon

À ce stade, **Club Manager possède maintenant une vraie base de données PostgreSQL versionnée par des migrations**.

Le projet ne se contente donc plus uniquement de servir des pages HTTP et de charger sa configuration.

Il commence à disposer d'une couche de persistance.

Ce jalon correspond à l'état du projet après :

- l'installation et la configuration de PostgreSQL ;
    
- la création de la base `club_manager` ;
    
- la création de l'utilisateur PostgreSQL `club_manager` ;
    
- la prise en main de `psql` ;
    
- la mise en place de Goose ;
    
- la création de la première migration ;
    
- l'application de cette migration dans PostgreSQL.
    

---

# 1. Évolution générale de Club Manager

Avant ce jalon, le projet ressemblait essentiellement à :

```text
Navigateur
    │
    ▼
Router
    │
    ▼
Handlers
    │
    ▼
Views / Templates

Configuration
    │
    └──► application
```

Il possède maintenant une nouvelle brique :

```text
Navigateur
    │
    ▼
Router
    │
    ▼
Handlers
    │
    ▼
Views / Templates


Configuration


PostgreSQL
    ▲
    │
Migrations
    ▲
    │
  Goose
```

Pour l'instant, l'application Go n'utilise pas encore directement PostgreSQL.

La base de données existe cependant et sa structure peut maintenant évoluer proprement grâce aux migrations.

---

# 2. PostgreSQL

PostgreSQL est le **système de gestion de base de données** utilisé par Club Manager.

Son rôle sera de stocker les données persistantes de l'application.

Par exemple :

```text
members
roles
memberships
payments
equipment
...
```

Contrairement au fichier de configuration, ces données peuvent être :

- nombreuses ;
    
- modifiées fréquemment ;
    
- recherchées ;
    
- filtrées ;
    
- reliées entre elles.
    

---

# 3. `psql`

`psql` est le client en ligne de commande de PostgreSQL.

Il permet de communiquer directement avec le serveur PostgreSQL.

On peut le voir comme :

```text
Utilisateur
    │
    ▼
  psql
    │
    ▼
PostgreSQL
```

`psql` n'est donc **pas la base de données**.

C'est un outil permettant de lui envoyer des commandes SQL et de l'administrer.

---

## Exemples d'utilisation

Connexion à PostgreSQL :

```bash
psql -U club_manager -d club_manager
```

Une fois connecté, il est possible d'exécuter du SQL :

```sql
SELECT * FROM members;
```

ou des commandes propres à `psql`.

Par exemple, afficher les tables :

```text
\dt
```

Afficher la structure d'une table :

```text
\d members
```

Quitter :

```text
\q
```

---

# 4. Pourquoi avoir utilisé `psql` ?

`psql` nous permet notamment de vérifier ce que fait réellement l'application ou Goose.

Par exemple :

```text
Goose
  │
  │ exécute une migration
  ▼
PostgreSQL
```

Puis :

```text
psql
  │
  │ inspection
  ▼
PostgreSQL
```

On peut ainsi vérifier que la table créée existe réellement.

C'est particulièrement utile pour apprendre car cela permet de voir directement l'état de la base.

---

# 5. Goose

Goose est un outil de **migration de base de données**.

Une migration représente une modification contrôlée de la structure de la base.

Par exemple :

```text
Créer la table members
```

puis plus tard :

```text
Ajouter une colonne email
```

puis :

```text
Créer la table memberships
```

Chaque évolution possède ainsi une trace dans le projet.

---

# 6. Pourquoi utiliser des migrations ?

Sans migrations, on pourrait modifier directement PostgreSQL :

```sql
CREATE TABLE members (...);
```

Le problème est que cette modification existerait uniquement dans cette base.

Un autre développeur ou un nouveau serveur ne saurait pas forcément comment reconstruire la même structure.

Avec Goose :

```text
Code source
│
├── application Go
│
└── migrations
      │
      ├── migration 1
      ├── migration 2
      └── migration 3
```

La structure de la base devient donc elle aussi **versionnée et reproductible**.

---

# 7. Première migration de Club Manager

La première migration concerne les membres.

Le fichier ressemble à :

```text
migrations/
└── 20260810164923_create_members.sql
```

Le nombre placé devant le nom représente la version de la migration.

Par exemple :

```text
20260810164923
```

permet à Goose de savoir dans quel ordre appliquer les migrations.

---

# 8. Structure d'une migration Goose

Une migration SQL Goose contient généralement deux parties :

```sql
-- +goose Up
```

et :

```sql
-- +goose Down
```

## `Up`

La partie `Up` décrit comment appliquer la migration.

Par exemple :

```sql
-- +goose Up

CREATE TABLE members (
    ...
);
```

---

## `Down`

La partie `Down` décrit comment revenir en arrière.

Par exemple :

```sql
-- +goose Down

DROP TABLE members;
```

On obtient donc :

```text
Up
 │
 ▼
état précédent ─────────► nouvel état

Down
 │
 ▼
nouvel état ────────────► état précédent
```

---

# 9. Exécution de la migration

La migration a été appliquée avec :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
up
```

Résultat obtenu :

```text
OK   20260810164923_create_members.sql
goose: successfully migrated database to version: 20260810164923
```

La base est donc maintenant à la version :

```text
20260810164923
```

---

# 10. Comprendre la commande Goose

La commande :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
up
```

peut être décomposée ainsi.

### `goose`

Programme que l'on souhaite exécuter.

```text
goose
```

---

### `-dir migrations`

Indique où se trouvent les fichiers de migration.

```text
-dir migrations
```

Donc :

```text
Club Manager
│
└── migrations/
```

---

### `postgres`

Indique à Goose le type de base utilisé.

```text
postgres
```

Goose sait ainsi quel pilote SQL utiliser.

---

### Chaîne de connexion

```text
host=localhost
port=5432
user=club_manager
dbname=club_manager
sslmode=require
```

Ces informations expliquent à Goose **où se trouve la base et comment s'y connecter**.

On peut lire cette chaîne comme :

```text
Serveur     : localhost
Port        : 5432
Utilisateur : club_manager
Base        : club_manager
SSL         : requis
```

---

### `up`

Enfin :

```text
up
```

demande à Goose d'appliquer toutes les migrations qui n'ont pas encore été exécutées.

---

# 11. Goose mémorise les migrations appliquées

Goose doit savoir quelles migrations ont déjà été exécutées.

Il ajoute pour cela sa propre table dans PostgreSQL.

On peut donc avoir quelque chose conceptuellement proche de :

```text
PostgreSQL
│
├── members
│
└── goose_db_version
```

La table `goose_db_version` permet à Goose de suivre l'état des migrations.

Ainsi, si :

```text
migration 1 → appliquée
migration 2 → appliquée
migration 3 → non appliquée
```

un :

```bash
goose ... up
```

n'appliquera que :

```text
migration 3
```

---

# 12. Séparation des responsabilités

Ce jalon introduit une nouvelle séparation importante.

## PostgreSQL

Stocke les données.

```text
PostgreSQL
→ persistance
```

## `psql`

Permet à l'humain d'interroger et d'administrer PostgreSQL.

```text
psql
→ interaction manuelle
```

## Goose

Gère l'évolution de la structure de PostgreSQL.

```text
Goose
→ migrations
```

## Go

Contiendra la logique applicative.

```text
Go
→ application
```

À terme :

```text
Utilisateur
    │
    ▼
Application Go
    │
    ▼
Requêtes SQL
    │
    ▼
PostgreSQL
```

Pendant que :

```text
Goose
    │
    ▼
structure de PostgreSQL
```

---

# 13. Ce qui existe maintenant

À ce jalon, Club Manager possède donc :

```text
club-manager/
│
├── cmd/
│   └── server/
│
├── internal/
│   ├── handlers/
│   ├── router/
│   └── ...
│
├── migrations/
│   └── 20260810164923_create_members.sql
│
├── config/
│
├── go.mod
│
└── ...
```

Et côté infrastructure :

```text
Club Manager
│
├── serveur HTTP Go
│
├── configuration
│
├── PostgreSQL
│
└── Goose
```

---

# 14. Ce qui manque encore

La base existe.

La table `members` existe.

Mais l'application Go ne sait pas encore l'utiliser.

Nous sommes actuellement ici :

```text
Application Go

        ✕ connexion absente

PostgreSQL
```

La prochaine étape consistera à créer le lien :

```text
Application Go
        │
        ▼
PostgreSQL
```

Puis à permettre au programme d'exécuter des requêtes.

---

# 15. Suite logique

Une progression possible est maintenant :

```text
1. PostgreSQL
      ✅

2. migrations Goose
      ✅

3. connexion Go → PostgreSQL
      ↓

4. requêtes SQL
      ↓

5. sqlc
      ↓

6. repository / accès aux données
      ↓

7. handler POST
      ↓

8. création d'un membre
```

L'un des premiers parcours complets pourra donc être :

```text
Formulaire
    │
    ▼
POST /members
    │
    ▼
Handler
    │
    ▼
Requête SQL
    │
    ▼
PostgreSQL
    │
    ▼
table members
```

---

# 16. Pourquoi ce jalon est important ?

Jusqu'à présent, Club Manager permettait principalement de travailler sur :

- HTTP ;
    
- les routes ;
    
- les handlers ;
    
- les vues ;
    
- les structs ;
    
- la configuration ;
    
- la séparation des responsabilités.
    

Ce jalon introduit une nouvelle grande notion :

> **la persistance des données**

Les données peuvent maintenant survivre à l'arrêt du serveur Go.

```text
Serveur Go arrêté
       │
       ✕
       │
PostgreSQL
       │
       └── les données restent présentes
```

---

# Comprendre et retenir

### PostgreSQL

> PostgreSQL est le système qui stocke durablement les données de Club Manager.

### `psql`

> `psql` est un client en ligne de commande permettant à un humain de communiquer avec PostgreSQL.

### Goose

> Goose ne stocke pas les données de l'application : il gère les modifications de la structure de la base.

### Migration

> Une migration décrit une évolution versionnée de la structure de la base de données.

### `goose up`

> Applique les migrations qui n'ont pas encore été exécutées.

### `goose down`

> Permet de revenir en arrière sur une migration lorsque celle-ci définit une opération `Down`.

### Chaîne de connexion

> Elle indique à un programme où se trouve PostgreSQL et comment s'y connecter.

### Idée essentielle du jalon

Avant :

```text
Club Manager → pages web
```

Maintenant :

```text
Club Manager
├── pages web
└── base PostgreSQL versionnée
```

La prochaine grande étape consiste à réunir ces deux mondes :

```text
HTTP
 │
 ▼
Go
 │
 ▼
SQL
 │
 ▼
PostgreSQL
```