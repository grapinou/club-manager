
---


## Objectif

Mettre en place PostgreSQL pour le projet **Club Manager** et comprendre les premières différences avec SQLite.

À ce stade, nous ne connectons pas encore l'application Go à PostgreSQL.

L'objectif est d'abord de comprendre :

- le serveur PostgreSQL ;
    
- le client `psql` ;
    
- les clusters ;
    
- les bases de données ;
    
- les rôles ;
    
- la connexion à une base.
    

---

# 1. Vérifier si PostgreSQL est installé

Commande :

```bash
psql --version
```

Résultat initial :

```text
La commande « psql » n'a pas été trouvée
```

PostgreSQL n'était donc pas installé.

---

# 2. Installer PostgreSQL

Mise à jour de la liste des paquets :

```bash
sudo apt update
```

Installation :

```bash
sudo apt install postgresql
```

Le paquet `postgresql` installe notamment :

- le serveur PostgreSQL ;
    
- le client `psql` ;
    
- les outils nécessaires à son fonctionnement.
    

---

# 3. Vérifier l'installation

Commande :

```bash
psql --version
```

Résultat :

```text
psql (PostgreSQL) 16.14 (Ubuntu 16.14-0ubuntu0.24.04.1)
```

PostgreSQL 16 est donc installé.

---

# 4. Vérifier le service PostgreSQL

Commande :

```bash
sudo systemctl status postgresql
```

Résultat important :

```text
Active: active (exited)
```

Cela ne signifie pas que PostgreSQL est arrêté.

Sous Ubuntu, `postgresql.service` est un service général qui gère les différentes instances PostgreSQL.

Il faut donc regarder les clusters pour connaître l'état réel du serveur.

---

# 5. Vérifier le cluster PostgreSQL

Commande :

```bash
pg_lsclusters
```

Résultat :

```text
Ver Cluster Port Status Owner    Data directory
16  main    5432 online postgres /var/lib/postgresql/16/main
```

Le serveur est donc bien actif.

## Signification

|Élément|Valeur|Rôle|
|---|---|---|
|Version|`16`|version majeure de PostgreSQL|
|Cluster|`main`|nom de l'instance|
|Port|`5432`|port réseau utilisé|
|Status|`online`|serveur actif|
|Owner|`postgres`|utilisateur système propriétaire|
|Data directory|`/var/lib/postgresql/16/main`|emplacement des données|

---

# Notion : cluster PostgreSQL

Dans ce contexte, un **cluster** n'est pas forcément un ensemble de plusieurs machines.

Un cluster PostgreSQL représente essentiellement une instance PostgreSQL possédant :

- ses bases ;
    
- sa configuration ;
    
- son répertoire de données ;
    
- son port réseau.
    

Dans notre cas :

```text
PostgreSQL 16
└── cluster main
    ├── port 5432
    ├── owner postgres
    └── status online
```

---

# 6. Première connexion avec `psql`

Commande :

```bash
sudo -u postgres psql
```

Décomposition :

```text
sudo
└── exécuter une commande avec une autre identité

-u postgres
└── utiliser l'utilisateur système postgres

psql
└── lancer le client PostgreSQL
```

Résultat :

```text
postgres=#
```

Nous sommes maintenant dans l'interface `psql`.

---

# 7. Terminal Linux vs psql

Avant :

```text
sighto@sighto:~$
```

Nous sommes dans le shell Linux.

Après :

```text
postgres=#
```

Nous sommes dans `psql`.

À partir de là, les commandes saisies sont :

- du SQL ;
    
- ou des commandes propres à `psql`.
    

---

# 8. SQL et commandes `psql`

## SQL

Exemple :

```sql
SELECT ...;
CREATE TABLE ...;
```

Une instruction SQL se termine généralement par :

```text
;
```

## Commandes psql

Les commandes propres à `psql` commencent par :

```text
\
```

Exemples :

```text
\l
\du
\c
\conninfo
\q
```

Elles ne nécessitent pas de `;`.

---

# 9. Lister les bases de données

Commande :

```text
\l
```

Bases présentes après l'installation :

```text
postgres
template0
template1
```

## `postgres`

Base fournie par défaut.

Elle peut notamment servir de base de connexion pour les opérations administratives.

## `template1`

Modèle utilisé normalement lors de la création d'une nouvelle base.

Schématiquement :

```text
template1
    ↓
CREATE DATABASE
    ↓
nouvelle base
```

## `template0`

Modèle de référence que l'on évite généralement de modifier.

---

# Différence importante avec SQLite

Avec SQLite :

```text
Application
    ↓
fichier .db
```

Une base est essentiellement un fichier.

Avec PostgreSQL :

```text
Application
    ↓
serveur PostgreSQL
    ├── base 1
    ├── base 2
    ├── base 3
    └── ...
```

PostgreSQL est un véritable serveur capable de gérer plusieurs bases.

---

# 10. Lister les rôles

Commande :

```text
\du
```

Résultat initial :

```text
Role name | Attributes
postgres  | Superuser, Create role, Create DB, Replication, Bypass RLS
```

Le rôle `postgres` possède des privilèges très élevés.

Il s'agit du compte administrateur PostgreSQL.

---

# Notion : rôle PostgreSQL

PostgreSQL utilise principalement la notion de **rôle**.

Un rôle peut :

- posséder des objets ;
    
- recevoir des permissions ;
    
- éventuellement se connecter à PostgreSQL.
    

Un rôle capable de se connecter possède l'attribut :

```text
LOGIN
```

Ainsi :

```sql
CREATE ROLE club_manager;
```

crée un rôle sans connexion directe.

Alors que :

```sql
CREATE ROLE club_manager WITH LOGIN;
```

crée un rôle capable de se connecter.

---

# 11. Créer le rôle Club Manager

Commande :

```sql
CREATE ROLE club_manager WITH LOGIN;
```

Résultat :

```text
CREATE ROLE
```

Vérification :

```text
\du
```

Résultat :

```text
club_manager |
postgres     | Superuser, Create role, Create DB, ...
```

Le rôle `club_manager` :

- peut se connecter ;
    
- n'est pas superutilisateur ;
    
- ne peut pas créer librement d'autres rôles ;
    
- ne peut pas créer librement d'autres bases.
    

C'est volontaire.

---

# Principe de sécurité

L'application Club Manager ne doit pas utiliser le compte :

```text
postgres
```

car il dispose de beaucoup trop de privilèges.

On préfère :

```text
PostgreSQL
│
├── postgres
│   └── administration
│
└── club_manager
    └── application
```

Principe général :

> Une application ne doit disposer que des permissions dont elle a réellement besoin.

---

# 12. Créer la base `club_manager`

Commande :

```sql
CREATE DATABASE club_manager OWNER club_manager;
```

Résultat :

```text
CREATE DATABASE
```

Cette commande est exécutée par `postgres`.

Elle crée une base nommée :

```text
club_manager
```

dont le propriétaire est :

```text
club_manager
```

---

# 13. Architecture obtenue

Nous avons maintenant :

```text
Serveur PostgreSQL
│
├── rôles
│   ├── postgres
│   │   └── administrateur
│   │
│   └── club_manager
│       └── rôle destiné à l'application
│
└── bases
    ├── postgres
    ├── template0
    ├── template1
    └── club_manager
        └── owner : club_manager
```

---

# 14. Se connecter à la base `club_manager`

Depuis `psql` :

```text
\c club_manager
```

Résultat :

```text
You are now connected to database "club_manager" as user "postgres".
```

Le prompt devient :

```text
club_manager=#
```

Attention : cela indique principalement que la **base actuelle** est `club_manager`.

Nous sommes encore connectés avec le rôle :

```text
postgres
```

Donc :

```text
Base : club_manager
Utilisateur PostgreSQL : postgres
```

---

# Base et utilisateur sont deux notions différentes

Il faut bien distinguer :

```text
Rôle
└── qui suis-je ?

Base
└── où suis-je connecté ?
```

Dans notre situation actuelle :

```text
postgres
   │
   │ se connecte à
   ▼
club_manager
```

L'objectif futur sera :

```text
Application Go
     │
     ▼
rôle club_manager
     │
     ▼
base club_manager
```

---

# Commandes vues

|Commande|Fonction|
|---|---|
|`psql --version`|afficher la version du client PostgreSQL|
|`sudo systemctl status postgresql`|consulter l'état du service PostgreSQL|
|`pg_lsclusters`|afficher les clusters PostgreSQL|
|`sudo -u postgres psql`|ouvrir `psql` avec l'utilisateur système `postgres`|
|`\l`|lister les bases|
|`\du`|lister les rôles|
|`\c nom_base`|changer de base|
|`\conninfo`|afficher les informations de connexion|
|`\q`|quitter `psql`|

---

# SQL vu

Créer un rôle pouvant se connecter :

```sql
CREATE ROLE club_manager WITH LOGIN;
```

Créer une base avec un propriétaire donné :

```sql
CREATE DATABASE club_manager OWNER club_manager;
```

---

# SQLite vs PostgreSQL

## SQLite

```text
Programme
   │
   ▼
Bibliothèque SQLite
   │
   ▼
fichier database.db
```

Peu de configuration nécessaire.

---

## PostgreSQL

```text
Programme
   │
   │ connexion
   ▼
Serveur PostgreSQL
   │
   ├── rôles
   ├── bases
   ├── permissions
   └── données
```

Cela introduit de nouvelles notions :

- serveur ;
    
- client ;
    
- port réseau ;
    
- authentification ;
    
- rôles ;
    
- permissions ;
    
- plusieurs bases ;
    
- connexions.
    

---

# Comprendre et retenir

## PostgreSQL est un serveur

Contrairement à SQLite, PostgreSQL fonctionne indépendamment du programme Go.

```text
Club Manager ≠ PostgreSQL
```

Club Manager deviendra un **client** de PostgreSQL.

---

## `psql` n'est pas PostgreSQL

`psql` est un client permettant de communiquer avec le serveur.

```text
psql
  ↓
PostgreSQL
```

De la même manière, plus tard :

```text
Club Manager
    ↓
PostgreSQL
```

---

## Une base et un utilisateur sont différents

On peut être :

```text
utilisateur postgres
```

tout en étant connecté à :

```text
base club_manager
```

---

## Ne pas utiliser le superutilisateur dans l'application

Le rôle :

```text
postgres
```

sert à l'administration.

Le rôle :

```text
club_manager
```

sera utilisé par l'application.

Cela limite les dégâts possibles en cas d'erreur ou de problème de sécurité.

---

# Où en sommes-nous ?

Installation :

```text
PostgreSQL 16.14 ✅
```

Serveur :

```text
cluster main
port 5432
online ✅
```

Rôle applicatif :

```text
club_manager ✅
```

Base :

```text
club_manager ✅
```

Propriétaire :

```text
club_manager ✅
```

Connexion actuelle :

```text
base : club_manager
rôle : postgres
```

---

# Prochaine étape

La prochaine étape sera de faire fonctionner une véritable connexion :

```text
rôle club_manager
        ↓
base club_manager
```

Il faudra donc aborder :

- l'authentification ;
    
- le mot de passe PostgreSQL ;
    
- les informations de connexion.
    

Nous pourrons ensuite tester cette connexion indépendamment de Go avant de commencer l'intégration dans Club Manager.
