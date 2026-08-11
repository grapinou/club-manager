
---



## Objectif

Comprendre comment lire une commande Goose complète et à quoi servent les paramètres de connexion PostgreSQL.

Commande étudiée :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
up
```

---

# Structure générale

On peut découper la commande ainsi :

```text
goose
│
├── -dir migrations
│
├── postgres
│
├── "host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require"
│
└── up
```

Chaque élément répond à une question différente.

---

# 1. `goose`

```bash
goose
```

Lance le programme Goose.

C'est l'outil chargé de gérer les migrations.

---

# 2. `-dir migrations`

```bash
-dir migrations
```

Indique à Goose où se trouvent les fichiers de migration.

Dans Club Manager :

```text
club-manager/
└── migrations/
    └── 20260810164923_create_members.sql
```

Donc :

```text
-dir migrations
      ↓
chercher les migrations dans ./migrations
```

---

# 3. `postgres`

```bash
postgres
```

Indique à Goose quel moteur de base de données est utilisé.

Ici :

```text
postgres
   ↓
PostgreSQL
```

Il ne faut pas comprendre :

```text
exécuter PostgreSQL
```

mais plutôt :

```text
utiliser le driver / moteur PostgreSQL
```

Goose doit connaître le type de base auquel il parle.

---

# 4. La chaîne de connexion

```text
host=localhost
port=5432
user=club_manager
dbname=club_manager
sslmode=require
```

Cette partie indique à Goose comment se connecter au serveur PostgreSQL.

Conceptuellement :

```text
Goose
  │
  ├── où se trouve PostgreSQL ?
  ├── sur quel port ?
  ├── avec quel rôle ?
  ├── dans quelle base ?
  └── avec quel mode SSL ?
  ↓
PostgreSQL
```

---

# `host=localhost`

```text
host=localhost
```

Indique la machine sur laquelle PostgreSQL fonctionne.

Dans notre cas :

```text
localhost
```

signifie :

```text
la machine actuelle
```

Il correspond ici à :

```text
127.0.0.1
```

Nous l'avions vérifié avec :

```text
\conninfo
```

---

# `port=5432`

```text
port=5432
```

Indique le port réseau utilisé par PostgreSQL.

Nous l'avions observé avec :

```bash
pg_lsclusters
```

Résultat :

```text
Port : 5432
```

C'est le port par défaut de PostgreSQL.

---

# `user=club_manager`

```text
user=club_manager
```

Indique le rôle PostgreSQL utilisé par Goose.

Nous avons créé ce rôle avec :

```sql
CREATE ROLE club_manager WITH LOGIN;
```

Il s'agit du rôle destiné à l'application et aux migrations du projet.

On évite d'utiliser :

```text
postgres
```

car il s'agit du superutilisateur.

---

# `dbname=club_manager`

```text
dbname=club_manager
```

Indique la base de données dans laquelle Goose doit appliquer les migrations.

Nous avons créé cette base avec :

```sql
CREATE DATABASE club_manager OWNER club_manager;
```

Donc :

```text
serveur PostgreSQL
└── base club_manager
```

---

# `sslmode=require`

```text
sslmode=require
```

Demande que la connexion utilise SSL/TLS.

Lors de notre test avec `psql`, PostgreSQL indiquait :

```text
SSL connection
protocol: TLSv1.3
```

Donc notre connexion locale supporte bien SSL.

---

# 5. `up`

```bash
up
```

Indique l'action que Goose doit effectuer.

`up` signifie :

```text
appliquer les migrations qui ne l'ont pas encore été
```

Exemple :

```text
migration 1 → déjà appliquée
migration 2 → Pending
migration 3 → Pending

goose ... up
        ↓

migration 2 → appliquée
migration 3 → appliquée
```

---

# Lire la commande en français

La commande :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
up
```

peut se lire comme :

> Lance Goose, cherche les migrations dans le dossier `migrations`, utilise PostgreSQL, connecte-toi au serveur local sur le port 5432 avec le rôle `club_manager`, travaille dans la base `club_manager`, utilise SSL, puis applique les migrations en attente.

---

# Comparaison avec `psql`

Nous avions utilisé :

```bash
psql -h localhost -U club_manager -d club_manager
```

On retrouve pratiquement les mêmes informations.

```text
psql                     Goose
────────────────────────────────────────

-h localhost       →     host=localhost

-U club_manager    →     user=club_manager

-d club_manager    →     dbname=club_manager

                       + port=5432
                       + sslmode=require
```

Les deux programmes sont des clients PostgreSQL.

---

# Différence entre `psql` et Goose

## `psql`

```text
utilisateur
   ↓
psql
   ↓
PostgreSQL
```

Permet de :

- saisir du SQL manuellement ;
    
- observer les tables ;
    
- interroger les données ;
    
- administrer la base.
    

---

## Goose

```text
fichiers de migration
        ↓
      Goose
        ↓
   PostgreSQL
```

Permet de :

- détecter les migrations ;
    
- les appliquer dans l'ordre ;
    
- mémoriser celles déjà exécutées ;
    
- revenir en arrière avec les migrations `Down`.
    

---

# Les paramètres ne sont pas propres à Goose

Les informations :

```text
host
port
user
password
database
```

décrivent la connexion PostgreSQL elle-même.

Elles seront utilisées par différents clients :

```text
psql
goose
application Go
```

Schéma :

```text
                 PostgreSQL
                     ▲
                     │
      ┌──────────────┼──────────────┐
      │              │              │
    psql           Goose        Club Manager
```

Tous doivent savoir comment atteindre la base.

---

# Le mot de passe

Le mot de passe n'était pas écrit dans la chaîne de connexion.

Nous avons utilisé :

```bash
read -s PGPASSWORD
export PGPASSWORD
```

Cela permet d'éviter d'écrire le mot de passe directement dans l'historique du terminal.

Puis Goose peut utiliser cette variable lors de la connexion.

---

# Les `\` dans la commande

Exemple :

```bash
goose -dir migrations postgres \
"host=localhost ..." \
up
```

Les `\` ne sont pas des options Goose.

Ils appartiennent au shell Bash.

Ils signifient :

```text
la commande continue à la ligne suivante
```

La commande suivante est donc strictement équivalente :

```bash
goose -dir migrations postgres "host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" up
```

Les retours à la ligne servent uniquement à améliorer la lisibilité.

---

# Vérifier l'état des migrations

Commande :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
status
```

Avant application :

```text
Pending -- 20260810164923_create_members.sql
```

Après application :

```text
migration appliquée
```

La commande `status` permet donc de savoir où en est la base.

---

# Première migration appliquée

Migration :

```text
20260810164923_create_members.sql
```

Commande :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
up
```

Résultat :

```text
OK 20260810164923_create_members.sql
```

Puis :

```text
successfully migrated database to version: 20260810164923
```

Cela signifie que la base est maintenant à cette version de migration.

---

# Une migration en erreur

Lors du premier essai, PostgreSQL a retourné :

```text
syntax error at or near ")"
```

La migration contenait notamment :

```sql
created_at TIMESTAMPZ NOT NULL DEFAULT NOW(),
);
```

Deux problèmes étaient présents :

```text
TIMESTAMPZ
    ↓
doit être TIMESTAMPTZ
```

et :

```text
NOW(),
     ↑
virgule inutile avant la fermeture de la table
```

Correction :

```sql
created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

# Que fait Goose lorsqu'une migration échoue ?

Une migration qui échoue n'est pas considérée comme appliquée.

Elle reste donc à exécuter.

Cela permet de :

```text
migration Pending
      ↓
erreur SQL
      ↓
corriger le fichier
      ↓
relancer goose up
```

Tant que la migration n'a pas été appliquée avec succès et diffusée sur d'autres environnements, on peut corriger directement son fichier.

---

# `Up` et `Down`

Une migration Goose contient généralement deux directions.

Exemple :

```sql
-- +goose Up

CREATE TABLE members (...);

-- +goose Down

DROP TABLE members;
```

## `Up`

```text
faire évoluer la base
```

## `Down`

```text
annuler cette évolution
```

Dans notre cas :

```text
Up
└── créer members

Down
└── supprimer members
```

---

# Comprendre et retenir

## Une commande Goose répond à quatre questions

```text
Où sont les migrations ?
→ -dir migrations

Quel moteur ?
→ postgres

Comment se connecter ?
→ host, port, user, dbname, sslmode

Que faire ?
→ up, status, down...
```

---

## Goose ne devine pas la base

Le fichier de migration indique :

```sql
CREATE TABLE members (...)
```

mais il ne dit pas :

```text
sur quel serveur ?
dans quelle base ?
avec quel utilisateur ?
```

Ces informations viennent de la commande de connexion.

---

## Migration et connexion sont deux choses différentes

```text
migration
→ décrit ce qui doit changer

connexion
→ indique où appliquer ce changement
```

Exemple :

```text
20260810164923_create_members.sql
               ↓
            Goose
               ↓
connexion PostgreSQL
               ↓
base club_manager
```

---

# Commande de référence

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" \
up
```

Lecture :

```text
goose
→ programme

-dir migrations
→ dossier des migrations

postgres
→ moteur PostgreSQL

host=localhost
→ serveur

port=5432
→ port

user=club_manager
→ rôle

dbname=club_manager
→ base

sslmode=require
→ connexion SSL obligatoire

up
→ appliquer les migrations
```

---

# État actuel de Club Manager

```text
Club Manager
│
├── migrations/
│   └── 20260810164923_create_members.sql
│
└── PostgreSQL
    └── club_manager
        └── public
            └── members
```

La première migration métier du projet est maintenant appliquée et suivie par Goose.