
---


# Club Manager — Recréer la base et appliquer les migrations Goose

## Objectif

Après avoir conçu le nouveau socle de données de Club Manager, nous avons recréé la base PostgreSQL afin de repartir d'un environnement propre.

L'objectif était de passer de :

```text
ancienne base de développement
└── ancien modèle `members`
```

à :

```text
nouvelle base vide
        ↓
migrations Goose
        ↓
nouveau schéma relationnel
```

Le nouveau schéma contient notamment :

```text
persons
seasons
membership_types
activities
memberships
trial_registrations
users
roles
user_roles
membership_activities
```

---

# 1. Supprimer l'ancienne base

Connexion à PostgreSQL en tant qu'administrateur :

```bash
sudo -u postgres psql
```

Puis :

```sql
DROP DATABASE club_manager;
```

Résultat :

```text
DROP DATABASE
```

Cela supprime :

```text
ancienne base club_manager
├── anciennes tables
├── anciennes données
└── ancien historique Goose
```

Mais cela ne supprime **pas** le rôle PostgreSQL :

```text
club_manager
```

---

# 2. Recréer la base

Toujours depuis `psql` :

```sql
CREATE DATABASE club_manager OWNER club_manager;
```

Résultat :

```text
CREATE DATABASE
```

La base appartient directement au rôle applicatif :

```text
database : club_manager
owner    : club_manager
```

On peut vérifier avec :

```sql
\l
```

Extrait obtenu :

```text
Name         | Owner
-------------+--------------
club_manager | club_manager
```

---

# 3. Vérifier la connexion

Quitter :

```sql
\q
```

Puis se connecter comme l'utilisateur applicatif :

```bash
psql -h localhost -U club_manager -d club_manager
```

La connexion doit arriver sur :

```text
club_manager=>
```

À ce stade, la base vient d'être créée et ne contient encore aucune table métier.

---

# 4. Valider les fichiers Goose

Depuis la racine du projet :

```bash
goose -dir migrations validate
```

Dans notre cas, la commande n'a affiché aucune sortie :

```text
sighto@sighto:~/club-manager$
```

Cela signifie que Goose n'a détecté aucun problème dans la structure de ses fichiers de migration.

## À retenir

```text
goose validate
        ↓
aucune erreur
        ↓
retour direct au terminal
        ↓
validation réussie
```

Une commande réussie n'affiche pas nécessairement :

```text
OK
```

---

# 5. Authentification PostgreSQL pour Goose

Lors du premier :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" status
```

PostgreSQL a retourné :

```text
password authentication failed for user "club_manager"
```

Cela ne signifiait pas que la base était mal créée.

Goose atteignait bien PostgreSQL, mais ne possédait simplement pas le mot de passe du rôle `club_manager`.

Nous avons donc défini :

```bash
export PGPASSWORD='mot_de_passe'
```

Puis relancé la commande.

## Important

```text
export PGPASSWORD=...
```

définit une variable d'environnement pour le shell courant.

Goose peut alors utiliser ce mot de passe lors de la connexion à PostgreSQL.

Il n'est pas nécessaire de lancer manuellement `psql` avant chaque commande Goose.

---

# 6. Vérifier l'état des migrations

Commande :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" status
```

Avant leur application, les migrations étaient toutes :

```text
Pending -- 0001_persons.sql
Pending -- 0002_seasons.sql
Pending -- 0003_membership_types.sql
Pending -- 0004_activities.sql
Pending -- 0005_memberships.sql
Pending -- 0006_trial_registrations.sql
Pending -- 0007_users.sql
Pending -- 0008_roles.sql
Pending -- 0009_user_roles.sql
Pending -- 0010_membership_activities.sql
```

`Pending` signifie :

> Goose connaît cette migration, mais elle n'a pas encore été exécutée sur cette base.

---

# 7. Appliquer les migrations

Commande :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" up
```

Résultat :

```text
OK 0001_persons.sql
OK 0002_seasons.sql
OK 0003_membership_types.sql
OK 0004_activities.sql
OK 0005_memberships.sql
OK 0006_trial_registrations.sql
OK 0007_users.sql
OK 0008_roles.sql
OK 0009_user_roles.sql
OK 0010_membership_activities.sql
```

Puis :

```text
goose: successfully migrated database to version: 10
```

Cela signifie que Goose a exécuté toutes les migrations dans l'ordre :

```text
0001
 ↓
0002
 ↓
...
 ↓
0010
```

et que la base est maintenant à la version :

```text
10
```

---

# 8. Vérifier l'historique Goose

Nous avons relancé :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" status
```

Les migrations ne sont alors plus indiquées comme `Pending`.

Elles possèdent toutes une date d'application.

Conceptuellement :

```text
avant `up`

0001 → Pending
0002 → Pending
...
0010 → Pending
```

devient :

```text
après `up`

0001 → appliquée
0002 → appliquée
...
0010 → appliquée
```

---

# 9. Vérifier directement dans PostgreSQL

Connexion :

```bash
psql -h localhost -U club_manager -d club_manager
```

Puis :

```sql
\dt
```

Résultat :

```text
activities
goose_db_version
membership_activities
membership_types
memberships
persons
roles
seasons
trial_registrations
user_roles
users
```

Nous retrouvons donc bien :

```text
10 tables métier
+
1 table interne Goose
```

---

# `goose_db_version`

La table :

```text
goose_db_version
```

n'a pas été créée directement par nos migrations métier.

Elle est utilisée par Goose pour savoir :

```text
quelles migrations ont été exécutées
```

et :

```text
à quelle version se trouve la base
```

Elle appartient donc à l'outil de migration, pas au modèle métier de Club Manager.

---

# 10. Inspecter une table avec `\d`

La commande :

```sql
\d memberships
```

permet d'inspecter précisément une table.

Elle affiche notamment :

```text
colonnes
types
valeurs par défaut
NOT NULL
index
clés primaires
UNIQUE
CHECK
clés étrangères
tables qui la référencent
```

---

# Exemple : `memberships`

PostgreSQL nous a montré :

```text
id
person_id
season_id
membership_type_id
status
joined_at
ended_at
created_at
updated_at
```

---

## Clé primaire

```text
memberships_pkey
PRIMARY KEY (id)
```

Cela correspond à :

```text
id
→ identifiant unique d'une adhésion
```

---

## Contrainte `UNIQUE`

```text
UNIQUE (person_id, season_id)
```

PostgreSQL l'affiche sous la forme :

```text
memberships_person_id_season_id_key
```

Cette contrainte garantit :

> une personne ne peut pas posséder deux adhésions pour une même saison dans le modèle actuel.

---

# `CHECK` sur `status`

La migration contenait conceptuellement :

```sql
CHECK (
    status IN (
        'pending',
        'active',
        'ended',
        'cancelled'
    )
)
```

PostgreSQL l'affiche sous une forme interne comme :

```text
status = ANY (
    ARRAY[
        'pending',
        'active',
        'ended',
        'cancelled'
    ]
)
```

Les deux expressions représentent la même règle.

PostgreSQL peut réécrire les contraintes sous une forme interne différente de celle saisie dans la migration.

---

# Clés étrangères de `memberships`

Nous retrouvons :

```text
membership_type_id
→ membership_types.id

person_id
→ persons.id

season_id
→ seasons.id
```

Cela confirme que les relations créées dans les migrations sont réellement présentes dans PostgreSQL.

---

# Relation inverse

PostgreSQL affiche également :

```text
Referenced by:
membership_activities
```

Cela signifie :

```text
membership_activities.membership_id
              │
              ▼
       memberships.id
```

Une `membership_activity` ne peut donc pas exister sans `membership` correspondante.

---

# 11. Vérifier `user_roles`

Commande :

```sql
\d user_roles
```

La table possède :

```text
user_id
role_id
```

et aucune colonne `id`.

Sa clé primaire est :

```text
PRIMARY KEY (user_id, role_id)
```

---

## Pourquoi une clé primaire composée ?

Le couple :

```text
user_id + role_id
```

identifie déjà parfaitement la relation.

Exemple :

```text
Alice + president
```

Il n'est donc pas nécessaire d'ajouter :

```text
id
```

PostgreSQL empêchera également :

```text
Alice → president
Alice → president
```

d'être enregistré deux fois.

---

## Clés étrangères

```text
user_id
→ users.id

role_id
→ roles.id
```

Cela matérialise :

```text
users
  │
  ▼
user_roles
  │
  ▼
roles
```

Un utilisateur peut donc avoir plusieurs rôles.

---

# 12. Vérifier `membership_activities`

Commande :

```sql
\d membership_activities
```

Colonnes :

```text
membership_id
activity_id
```

Clé primaire :

```text
PRIMARY KEY (
    membership_id,
    activity_id
)
```

Clés étrangères :

```text
membership_id
→ memberships.id

activity_id
→ activities.id
```

Cette table matérialise :

```text
memberships
      │
      ▼
membership_activities
      │
      ▼
activities
```

Une adhésion peut avoir plusieurs activités.

Une même activité peut appartenir à plusieurs adhésions.

Il s'agit donc d'une relation :

```text
plusieurs-à-plusieurs
```

---

# 13. Deux tables de liaison validées

Nous avons maintenant deux exemples très similaires :

```text
user_roles
```

et :

```text
membership_activities
```

## `user_roles`

```text
users
   │
   └── roles
```

avec :

```text
PRIMARY KEY (user_id, role_id)
```

## `membership_activities`

```text
memberships
   │
   └── activities
```

avec :

```text
PRIMARY KEY (
    membership_id,
    activity_id
)
```

---

# 14. Cycle complet utilisé

La procédure suivie a été :

```text
suppression ancienne DB
        ↓
création nouvelle DB
        ↓
goose validate
        ↓
goose status
        ↓
migrations Pending
        ↓
goose up
        ↓
goose status
        ↓
migrations appliquées
        ↓
psql
        ↓
\dt
        ↓
\d memberships
\d user_roles
\d membership_activities
```

Cela permet de vérifier progressivement chaque niveau.

---

# Pourquoi ne pas lancer directement `goose up` ?

Nous avons préféré avancer étape par étape :

```text
valider les fichiers
        ↓
vérifier que Goose les voit
        ↓
vérifier qu'ils sont Pending
        ↓
les appliquer
        ↓
inspecter le résultat
```

Cette approche est particulièrement utile lorsqu'on modifie la structure d'une base.

Elle permet de savoir précisément à quelle étape apparaît une erreur.

---

# Comprendre et retenir

## `goose validate`

```text
vérifie les fichiers de migration
```

Il n'applique rien à PostgreSQL.

---

## `goose status`

```text
compare :
migrations présentes
        avec
migrations déjà exécutées
```

---

## `goose up`

```text
exécute les migrations Pending
dans l'ordre
```

---

## `\dt`

```text
liste les tables
```

---

## `\d nom_table`

```text
inspecte la structure réelle
d'une table PostgreSQL
```

---

## `Pending`

```text
migration connue
mais pas encore exécutée
```

---

## `PRIMARY KEY (a, b)`

Une clé primaire peut être composée de plusieurs colonnes.

```text
(a, b)
```

devient alors l'identité de la ligne.

C'est particulièrement adapté aux tables de liaison.

---

## `FOREIGN KEY`

Une clé étrangère exprime :

```text
cette valeur doit correspondre
à une ligne existante
dans une autre table
```

Exemple :

```text
membership_id
→ memberships.id
```

---

# Point important : environnement de développement

Nous avons volontairement supprimé l'ancienne base car :

```text
ancien modèle pédagogique
        ↓
nouveau modèle relationnel
```

Nous étions encore dans un environnement de développement sans données importantes à conserver.

Dans une application en production, on ne supprimerait évidemment pas la base de cette manière.

On utiliserait de nouvelles migrations pour transformer progressivement les données existantes.

---

# Résultat du jalon

Le nouveau socle n'existe maintenant plus seulement sous forme de fichiers SQL.

Il existe réellement dans PostgreSQL :

```text
                 persons
                    │
       ┌────────────┼────────────┐
       │            │            │
       ▼            ▼            ▼
     users      memberships   trial_registrations
       │            │               │
       ▼            │               ▼
  user_roles        │           activities
       │            │
       ▼            ▼
     roles   membership_activities
                    │
                    ▼
                activities
```

Les migrations sont toutes appliquées :

```text
version Goose = 10
```

Les clés étrangères sont présentes.

Les contraintes sont présentes.

Les clés primaires composées sont présentes.

La base est donc prête pour la prochaine étape :

```text
PostgreSQL
    ↓
sqlc
    ↓
requêtes autour de persons
    ↓
code Go
```

# À retenir en une phrase

> Goose construit et versionne le schéma ; PostgreSQL matérialise et garantit les relations ; `psql` permet de vérifier que ce qui a été conçu existe réellement dans la base.