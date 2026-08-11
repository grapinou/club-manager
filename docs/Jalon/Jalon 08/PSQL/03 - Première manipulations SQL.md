
---


## Objectif

Découvrir les premières manipulations dans la base `club_manager` avant de créer les vraies tables métier du projet.

À ce stade, nous avons vu :

- les schémas PostgreSQL ;
    
- la création d'une table ;
    
- la description d'une table ;
    
- l'insertion de données ;
    
- la lecture ;
    
- le filtrage ;
    
- la modification ;
    
- la suppression de données ;
    
- la suppression d'une table.
    

---

# 1. Lister les tables

Commande `psql` :

```text
\dt
```

Au départ, la base était vide :

```text
Did not find any relations.
```

Dans ce contexte, on peut retenir simplement :

```text
\dt
 ↓
liste des tables
```

---

# 2. Lister les schémas

Commande :

```text
\dn
```

Résultat :

```text
List of schemas

Name   | Owner
-------+------------------
public | pg_database_owner
```

PostgreSQL contient donc un schéma par défaut :

```text
public
```

---

# Notion : schéma

Un schéma est un espace de noms à l'intérieur d'une base PostgreSQL.

L'organisation peut être vue comme :

```text
Serveur PostgreSQL
└── base club_manager
    └── schéma public
        └── tables
```

Une table peut donc être désignée complètement par :

```text
public.test_items
```

Dans Club Manager, nous utiliserons pour l'instant simplement le schéma :

```text
public
```

---

# 3. Créer une première table

Nous avons créé une table temporaire :

```sql
CREATE TABLE test_items (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL
);
```

Résultat :

```text
CREATE TABLE
```

---

# Structure de la table

```text
test_items
├── id
│   ├── INTEGER
│   ├── IDENTITY
│   └── PRIMARY KEY
│
└── name
    ├── TEXT
    └── NOT NULL
```

---

# 4. Comprendre les éléments du `CREATE TABLE`

## `CREATE TABLE`

```sql
CREATE TABLE test_items
```

crée une nouvelle table appelée :

```text
test_items
```

---

## `INTEGER`

```sql
id INTEGER
```

La colonne contient un entier.

---

## `GENERATED ALWAYS AS IDENTITY`

```sql
GENERATED ALWAYS AS IDENTITY
```

demande à PostgreSQL de générer automatiquement la valeur de l'identifiant.

Nous n'avons donc pas besoin de fournir nous-mêmes :

```text
id = 1
id = 2
id = 3
```

PostgreSQL le fait pour nous.

---

## `PRIMARY KEY`

```sql
PRIMARY KEY
```

indique que cette colonne permet d'identifier chaque ligne de manière unique.

Ainsi :

```text
id = 1
```

ne peut correspondre qu'à une seule ligne.

---

## `TEXT`

```sql
name TEXT
```

indique que la colonne contient du texte.

---

## `NOT NULL`

```sql
NOT NULL
```

interdit l'absence de valeur.

La colonne doit donc obligatoirement contenir une donnée.

---

# 5. Vérifier la présence de la table

Commande :

```text
\dt
```

Résultat :

```text
Schema | Name       | Type  | Owner
-------+------------+-------+--------------
public | test_items | table | club_manager
```

Nous avons maintenant :

```text
club_manager
└── public
    └── test_items
```

---

# 6. Décrire une table

Commande :

```text
\d test_items
```

Résultat :

```text
Table "public.test_items"

Column | Type    | Nullable | Default
-------+---------+----------+------------------------------
id     | integer | not null | generated always as identity
name   | text    | not null |
```

Cette commande permet de voir rapidement :

- les colonnes ;
    
- les types ;
    
- les contraintes ;
    
- les valeurs par défaut ;
    
- les index.
    

---

# 7. Clé primaire et index

PostgreSQL affichait également :

```text
Indexes:
    "test_items_pkey" PRIMARY KEY, btree (id)
```

Lorsque nous avons déclaré :

```sql
PRIMARY KEY
```

PostgreSQL a créé automatiquement un index sur `id`.

Pour l'instant, il suffit de retenir :

```text
PRIMARY KEY
     ↓
identification unique
     ↓
index créé automatiquement
```

Nous approfondirons les index plus tard.

---

# 8. Insérer une donnée

Commande :

```sql
INSERT INTO test_items (name)
VALUES ('Premier test');
```

Résultat :

```text
INSERT 0 1
```

La partie importante est :

```text
1
```

Une ligne a été insérée.

---

# Comprendre `INSERT`

```sql
INSERT INTO test_items
```

signifie :

```text
insérer dans la table test_items
```

Puis :

```sql
(name)
```

indique la colonne renseignée.

Enfin :

```sql
VALUES ('Premier test');
```

indique la valeur à ajouter.

---

# 9. Pourquoi ne pas fournir `id` ?

Nous avons écrit :

```sql
INSERT INTO test_items (name)
VALUES ('Premier test');
```

et non :

```sql
INSERT INTO test_items (id, name)
VALUES (1, 'Premier test');
```

Car la colonne `id` utilise :

```sql
GENERATED ALWAYS AS IDENTITY
```

PostgreSQL génère donc automatiquement l'identifiant.

---

# 10. Lire les données

Commande :

```sql
SELECT * FROM test_items;
```

Résultat :

```text
id | name
---+--------------
1  | Premier test
```

---

# Comprendre `SELECT`

```sql
SELECT
```

signifie sélectionner des données.

Le symbole :

```text
*
```

signifie :

```text
toutes les colonnes
```

Puis :

```sql
FROM test_items
```

indique dans quelle table chercher.

Donc :

```sql
SELECT * FROM test_items;
```

peut être lu comme :

```text
sélectionner toutes les colonnes
dans la table test_items
```

---

# 11. Ajouter une deuxième ligne

Commande :

```sql
INSERT INTO test_items (name)
VALUES ('Deuxième test');
```

PostgreSQL génère automatiquement :

```text
id = 2
```

Nous obtenons :

```text
id | name
---+----------------
1  | Premier test
2  | Deuxième test
```

L'identifiant sert à distinguer clairement les lignes.

---

# 12. Filtrer avec `WHERE`

Pour sélectionner uniquement la ligne dont l'identifiant vaut `2` :

```sql
SELECT * FROM test_items
WHERE id = 2;
```

Résultat :

```text
id | name
---+----------------
2  | Deuxième test
```

---

# Comprendre `WHERE`

`WHERE` ajoute une condition.

```text
SELECT
  ↓
quoi lire ?

FROM
  ↓
dans quelle table ?

WHERE
  ↓
sous quelle condition ?
```

Exemple :

```sql
WHERE id = 2
```

signifie :

```text
uniquement les lignes pour lesquelles id vaut 2
```

---

# 13. Modifier une donnée avec `UPDATE`

Commande :

```sql
UPDATE test_items
SET name = 'Test modifié'
WHERE id = 2;
```

---

# Comprendre `UPDATE`

```sql
UPDATE test_items
```

indique la table à modifier.

```sql
SET name = 'Test modifié'
```

indique la nouvelle valeur.

```sql
WHERE id = 2
```

indique la ligne concernée.

---

# Attention à `UPDATE`

Sans `WHERE` :

```sql
UPDATE test_items
SET name = 'Test modifié';
```

toutes les lignes seraient modifiées.

Il faut donc être particulièrement vigilant avec :

```sql
UPDATE
```

et vérifier la clause :

```sql
WHERE
```

---

# 14. Supprimer une ligne avec `DELETE`

Commande :

```sql
DELETE FROM test_items
WHERE id = 2;
```

Cela supprime uniquement la ligne dont :

```text
id = 2
```

---

# Attention à `DELETE`

Sans `WHERE` :

```sql
DELETE FROM test_items;
```

toutes les lignes seraient supprimées.

La table continuerait cependant d'exister.

---

# 15. Supprimer la table avec `DROP TABLE`

Commande :

```sql
DROP TABLE test_items;
```

Cette fois, ce n'est pas seulement le contenu qui est supprimé.

La table elle-même disparaît.

---

# `DELETE` vs `DROP TABLE`

## `DELETE`

```sql
DELETE FROM test_items;
```

Résultat :

```text
table conservée
données supprimées
```

Schéma :

```text
test_items
├── structure ✅
└── données   ❌
```

---

## `DROP TABLE`

```sql
DROP TABLE test_items;
```

Résultat :

```text
table supprimée
structure supprimée
données supprimées
```

Schéma :

```text
test_items ❌
```

---

# 16. Vérifier la suppression

Après :

```sql
DROP TABLE test_items;
```

on peut vérifier avec :

```text
\dt
```

La base devrait de nouveau répondre :

```text
Did not find any relations.
```

Nous revenons ainsi à une base propre.

---

# CRUD

Les opérations vues correspondent au principe général appelé **CRUD**.

## Create

Créer une donnée :

```sql
INSERT
```

Et créer une structure :

```sql
CREATE TABLE
```

---

## Read

Lire :

```sql
SELECT
```

---

## Update

Modifier :

```sql
UPDATE
```

---

## Delete

Supprimer :

```sql
DELETE
```

On peut résumer :

```text
CRUD
│
├── Create → INSERT
├── Read   → SELECT
├── Update → UPDATE
└── Delete → DELETE
```

`CREATE TABLE` et `DROP TABLE` concernent plutôt la structure de la base elle-même.

---

# Commandes `psql` vues

|Commande|Fonction|
|---|---|
|`\dt`|lister les tables|
|`\dn`|lister les schémas|
|`\d table`|décrire une table|

---

# SQL vu

## Créer une table

```sql
CREATE TABLE test_items (
    id INTEGER GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name TEXT NOT NULL
);
```

## Ajouter une ligne

```sql
INSERT INTO test_items (name)
VALUES ('Premier test');
```

## Lire toutes les lignes

```sql
SELECT * FROM test_items;
```

## Filtrer

```sql
SELECT * FROM test_items
WHERE id = 2;
```

## Modifier

```sql
UPDATE test_items
SET name = 'Test modifié'
WHERE id = 2;
```

## Supprimer une ligne

```sql
DELETE FROM test_items
WHERE id = 2;
```

## Supprimer une table

```sql
DROP TABLE test_items;
```

---

# Comprendre et retenir

## Une base contient des schémas

```text
base
 ↓
schéma
 ↓
table
 ↓
colonnes
 ↓
lignes
```

Dans notre cas :

```text
club_manager
└── public
    └── table
```

---

## Une table possède une structure

Par exemple :

```text
test_items
├── id
└── name
```

Chaque colonne possède :

- un nom ;
    
- un type ;
    
- éventuellement des contraintes.
    

---

## Une ligne représente une donnée

Exemple :

```text
id | name
---+--------------
1  | Premier test
```

Une table peut donc être vue comme :

```text
structure
   +
ensemble de lignes
```

---

## L'identifiant permet de cibler une ligne

```text
id = 1
id = 2
id = 3
```

Chaque identifiant correspond à une ligne précise.

Il peut ensuite être utilisé avec :

```sql
WHERE id = ...
```

---

## `WHERE` est particulièrement important

Avec :

```sql
SELECT
```

une erreur de `WHERE` peut simplement afficher trop de données.

Avec :

```sql
UPDATE
```

ou :

```sql
DELETE
```

une erreur peut modifier ou supprimer beaucoup de lignes.

Bonne habitude :

> Toujours vérifier attentivement le `WHERE` avant un `UPDATE` ou un `DELETE`.

---

# État à la fin de l'exercice

Nous avons créé temporairement :

```text
club_manager
└── public
    └── test_items
```

puis nous l'avons supprimée avec :

```sql
DROP TABLE test_items;
```

La base peut donc revenir à :

```text
club_manager
└── public
    └── aucune table
```

Nous avons ainsi découvert le SQL de base sans polluer la future structure métier de Club Manager.

---

# Prochaine étape

Nous pouvons maintenant commencer à réfléchir à la première véritable table du projet.

À partir de maintenant, une nouvelle question apparaît :

```text
Quelles données doit réellement stocker Club Manager ?
```

Ce sera le début du travail de **modélisation de la base de données** avant de connecter PostgreSQL au code Go.


### Remarque sur les schémas : 

Un schéma ne sert pas à repérer une base de données. Il sert surtout à **organiser les objets à l’intérieur d’une même base**.

Dans PostgreSQL, la hiérarchie est plutôt :

```
Serveur PostgreSQL
└── Base de données
    └── Schéma
        ├── tables
        ├── vues
        ├── séquences
        ├── fonctions
        └── autres objets
```

Dans notre cas :

```
PostgreSQL
└── club_manager
    └── public
        └── ...
```

Le schéma `public` est donc un espace de noms à l’intérieur de `club_manager`.

L’un des intérêts est de pouvoir avoir **deux objets portant le même nom**, tant qu’ils sont dans des schémas différents. Par exemple :

```
club_manager
├── public
│   └── members
│
└── archive
    └── members
```

Les deux tables peuvent exister simultanément :

```
public.members
archive.members
```

Le schéma sert aussi à **organiser fonctionnellement** une grosse application. On pourrait imaginer :

```
club_manager
├── members
│   ├── members
│   └── memberships
│
├── accounting
│   ├── invoices
│   └── payments
│
└── inventory
    ├── equipment
    └── stock_movements
```

Cela dit, pour Club Manager au début, je ne ferais surtout pas ça. `public` suffit largement. Multiplier les schémas trop tôt ajouterait de la complexité sans bénéfice réel.

Autre intérêt important : **les permissions**. PostgreSQL permet de donner des droits différents selon les schémas. On peut par exemple autoriser un rôle à utiliser :

```
accounting
```

mais pas :

```
administration
```

Les schémas servent donc aussi de frontière logique pour les droits.

Enfin, il existe la notion de `search_path`. Lorsque tu écris :

```
SELECT * FROM members;
```

PostgreSQL cherche la table dans les schémas configurés dans ce chemin de recherche. Avec la configuration habituelle, il trouve par exemple :

```
public.members
```

C’est pour cela que nous n’avons pas besoin d’écrire systématiquement :

```
SELECT * FROM public.members;
```

La façon la plus simple de retenir la notion serait :

```
Base de données
= un ensemble global de données

Schéma
= un espace de rangement / namespace dans cette base

Table
= une structure contenant réellement les données
```

Une analogie assez correcte serait celle des packages Go :

```
base PostgreSQL
└── schémas
    └── objets
```

comme :

```
projet Go
└── packages
    └── fonctions/types
```

Ce n’est pas exactement la même chose, mais l’idée de **namespace pour organiser et éviter les collisions de noms** est très proche.

Pour Club Manager, je garderais donc :

```
club_manager
└── public
    ├── members
    ├── ...
    └── ...
```

jusqu’au jour où un vrai besoin justifiera plusieurs schémas.
