
---


## Objectif

Finaliser l'installation de PostgreSQL pour **Club Manager** et vérifier qu'une application pourra se connecter à la base avec un rôle dédié.

À la fin de cette étape, nous devons avoir :

```text
PostgreSQL installé ✅
Serveur actif ✅
Base club_manager créée ✅
Rôle club_manager créé ✅
Mot de passe configuré ✅
Connexion applicative testée ✅
```

---

# 1. État du serveur PostgreSQL

Le cluster PostgreSQL est actif :

```bash
pg_lsclusters
```

Résultat :

```text
Ver Cluster Port Status Owner    Data directory
16  main    5432 online postgres /var/lib/postgresql/16/main
```

Notre serveur PostgreSQL fonctionne donc sur :

```text
port : 5432
état : online
```

---

# 2. Base et rôle créés

Nous avons créé un rôle PostgreSQL destiné à l'application :

```sql
CREATE ROLE club_manager WITH LOGIN;
```

Puis une base :

```sql
CREATE DATABASE club_manager OWNER club_manager;
```

Nous obtenons donc :

```text
PostgreSQL
│
├── rôle postgres
│   └── administration
│
├── rôle club_manager
│   └── application
│
└── base club_manager
    └── propriétaire : club_manager
```

---

# 3. Pourquoi ne pas utiliser `postgres` ?

Le rôle :

```text
postgres
```

est un superutilisateur.

Il possède notamment des permissions permettant de :

- créer des bases ;
    
- créer des rôles ;
    
- administrer PostgreSQL ;
    
- contourner certaines protections.
    

L'application Club Manager ne doit donc pas utiliser ce compte.

Nous utilisons :

```text
club_manager
```

qui possède uniquement les droits nécessaires à son fonctionnement.

---

# 4. Ajouter un mot de passe au rôle

Depuis `psql`, connecté comme administrateur :

```text
\password club_manager
```

PostgreSQL demande :

```text
Enter new password for user "club_manager":
Enter it again:
```

Le mot de passe ne s'affiche pas pendant la saisie.

C'est normal.

---

# Pourquoi utiliser `\password` ?

Il serait possible d'écrire :

```sql
ALTER ROLE club_manager PASSWORD 'mot_de_passe';
```

Mais cela place le mot de passe directement dans une commande SQL.

La commande :

```text
\password
```

est préférable lors d'une manipulation interactive.

---

# 5. Utilisateur Linux et rôle PostgreSQL

Une confusion est apparue avec :

```bash
sudo -u club_manager psql
```

Résultat :

```text
sudo: utilisateur club_manager inconnu
```

Cette erreur est normale.

## Pourquoi ?

`sudo` travaille avec les **utilisateurs Linux**.

Or :

```text
club_manager
```

est un **rôle PostgreSQL**, pas un utilisateur Linux.

Il faut distinguer les deux systèmes.

```text
Linux
│
├── sighto
└── postgres
```

et :

```text
PostgreSQL
│
├── postgres
└── club_manager
```

Les deux noms `postgres` existent dans les deux systèmes, ce qui peut prêter à confusion.

---

# Comprendre : utilisateur Linux vs rôle PostgreSQL

## Utilisateur Linux

Il appartient au système d'exploitation.

Exemple :

```text
sighto
postgres
```

Il peut être utilisé avec :

```bash
sudo -u utilisateur ...
```

---

## Rôle PostgreSQL

Il appartient au serveur PostgreSQL.

Exemple :

```text
postgres
club_manager
```

Il est utilisé pour :

- se connecter à PostgreSQL ;
    
- posséder des bases ou tables ;
    
- recevoir des permissions.
    

---

# 6. Tester une vraie connexion applicative

La connexion a été testée depuis le terminal avec :

```bash
psql -h localhost -U club_manager -d club_manager
```

Décomposition :

```text
psql
│
├── -h localhost
│   └── serveur PostgreSQL
│
├── -U club_manager
│   └── rôle PostgreSQL
│
└── -d club_manager
    └── base utilisée
```

PostgreSQL demande ensuite :

```text
Password for user club_manager:
```

Après saisie du mot de passe, la connexion fonctionne.

---

# 7. Prompt obtenu

Le prompt devient :

```text
club_manager=>
```

C'est différent du prompt administrateur :

```text
postgres=#
```

On peut retenir simplement :

```text
=#   → rôle avec privilèges élevés

=>   → rôle ordinaire
```

Le rôle `club_manager` n'est donc pas superutilisateur.

C'est exactement ce que nous recherchons.

---

# 8. Vérifier les informations de connexion

Commande :

```text
\conninfo
```

Résultat obtenu :

```text
You are connected to database "club_manager"
as user "club_manager"
on host "localhost"
(address "127.0.0.1")
at port "5432".
```

La connexion est donc :

```text
host     : localhost
adresse  : 127.0.0.1
port     : 5432
database : club_manager
user     : club_manager
```

---

# 9. Socket Unix vs connexion réseau

Lors de la première connexion administrateur :

```bash
sudo -u postgres psql
```

PostgreSQL utilisait un socket Unix :

```text
/var/run/postgresql
```

Schéma :

```text
psql
 │
 ▼
socket Unix
 │
 ▼
PostgreSQL
```

---

Lorsque nous faisons :

```bash
psql -h localhost -U club_manager -d club_manager
```

nous précisons un hôte.

La connexion devient :

```text
psql
 │
 ▼
localhost:5432
 │
 ▼
PostgreSQL
```

C'est une connexion réseau locale.

---

# 10. Pourquoi cette connexion ressemble à celle de Go ?

L'application Club Manager devra disposer d'informations similaires :

```text
host     = localhost
port     = 5432
user     = club_manager
password = ********
database = club_manager
```

Elle pourra alors faire :

```text
Club Manager
     │
     │ connexion PostgreSQL
     ▼
localhost:5432
     │
     ▼
rôle club_manager
     │
     ▼
base club_manager
```

La commande `psql` nous permet donc de tester manuellement ce que l'application devra faire automatiquement plus tard.

---

# 11. Connexion SSL

Lors de la connexion, PostgreSQL indique :

```text
SSL connection
(protocol: TLSv1.3,
cipher: TLS_AES_256_GCM_SHA384,
compression: off)
```

La connexion TCP locale utilise donc ici SSL/TLS.

Ce n'est pas une erreur.

---

# Architecture obtenue

À la fin de cette installation :

```text
Ubuntu
│
└── PostgreSQL 16
    │
    └── cluster main
        │
        ├── port 5432
        │
        ├── rôle postgres
        │   └── administration
        │
        ├── rôle club_manager
        │   └── connexion de l'application
        │
        └── base club_manager
            └── propriétaire : club_manager
```

---

# Connexion actuelle

La connexion testée avec succès est :

```text
Client     : psql

Serveur    : localhost
Adresse    : 127.0.0.1
Port       : 5432

Utilisateur PostgreSQL :
club_manager

Base :
club_manager
```

---

# Commandes à retenir

## Entrer comme administrateur local

```bash
sudo -u postgres psql
```

---

## Lister les rôles

```text
\du
```

---

## Lister les bases

```text
\l
```

---

## Changer de base

```text
\c nom_base
```

---

## Modifier un mot de passe

```text
\password nom_role
```

---

## Afficher les informations de connexion

```text
\conninfo
```

---

## Se connecter comme l'application

```bash
psql -h localhost -U club_manager -d club_manager
```

---

## Quitter psql

```text
\q
```

---

# Comprendre et retenir

## PostgreSQL possède ses propres utilisateurs

Les rôles PostgreSQL ne sont pas les utilisateurs Linux.

```text
Linux ≠ PostgreSQL
```

Ainsi :

```text
rôle PostgreSQL club_manager
```

ne signifie pas qu'il existe :

```text
utilisateur Linux club_manager
```

---

## Une application utilise un compte dédié

Club Manager ne doit pas se connecter comme :

```text
postgres
```

mais comme :

```text
club_manager
```

Cela applique le principe du moindre privilège.

---

## Une connexion PostgreSQL nécessite plusieurs informations

Contrairement à SQLite, il ne suffit pas de connaître le chemin d'un fichier.

Il faut notamment connaître :

```text
serveur
port
base
utilisateur
mot de passe
```

---

## `localhost` désigne notre propre machine

Dans notre configuration actuelle :

```text
localhost
```

correspond à :

```text
127.0.0.1
```

PostgreSQL et Club Manager fonctionneront donc pour l'instant sur la même machine.

---

# SQLite vs PostgreSQL

## SQLite

```text
Application
     │
     ▼
fichier .db
```

La base est directement manipulée par l'application.

---

## PostgreSQL

```text
Application
     │
     │ utilisateur
     │ mot de passe
     │ réseau
     ▼
Serveur PostgreSQL
     │
     ▼
Base de données
```

PostgreSQL est un service indépendant.

---

# État final

|Élément|État|
|---|---|
|PostgreSQL 16.14|✅|
|cluster `main`|✅|
|port `5432`|✅|
|serveur `online`|✅|
|rôle `club_manager`|✅|
|LOGIN|✅|
|mot de passe|✅|
|base `club_manager`|✅|
|propriétaire `club_manager`|✅|
|connexion TCP locale|✅|
|connexion avec rôle applicatif|✅|

---

# Étape suivante

L'installation de PostgreSQL est maintenant terminée.

Avant de connecter Go à PostgreSQL, nous pouvons commencer à explorer la base :

```text
club_manager=>
```

et notamment vérifier les tables présentes avec :

```text
\dt
```

À ce stade, la base devrait être vide.

Nous pourrons ensuite introduire progressivement :

```text
tables
    ↓
SQL
    ↓
migrations
    ↓
Go
    ↓
Club Manager
```

L'objectif reste de comprendre chaque couche avant d'ajouter la suivante.

