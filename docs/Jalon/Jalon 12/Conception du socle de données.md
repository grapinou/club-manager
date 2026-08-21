
---


# Club Manager — Conception du socle de données

## Objectif

Nous sommes passés d'un premier modèle pédagogique très simple :

```text
members
```

à un véritable **socle de données relationnel** permettant de représenter les besoins d'une association :

- personnes ;
    
- adhésions et historique ;
    
- saisons ;
    
- activités ;
    
- cours d'essai ;
    
- comptes utilisateurs ;
    
- rôles et permissions futures.
    

L'ancienne table `members` avait surtout permis d'apprendre et de mettre en pratique :

```text
GET
→ formulaire
→ POST
→ INSERT PostgreSQL
```

Elle peut maintenant disparaître.

Ce premier travail reste utile pédagogiquement : il marque le passage entre **faire fonctionner une fonctionnalité** et **concevoir un modèle de données durable**.

---

# Principe général

Une personne n'est pas forcément un membre.

Une personne n'est pas forcément un utilisateur.

Une personne venant faire un essai n'est pas non plus une catégorie différente de personne.

Nous avons donc séparé :

```text
Person
│
├── peut devenir User
│
├── peut avoir des Memberships
│
└── peut avoir des TrialRegistrations
```

Cela évite de répéter :

```text
first_name
last_name
birth_date
email
...
```

dans plusieurs tables.

`persons` devient ainsi la source principale de l'identité.

---

# Vue générale du modèle

```text
persons
│
├── users
│    └── user_roles ───── roles
│
├── trial_registrations ───── activities
│
└── memberships
     ├── seasons
     ├── membership_types
     │
     └── membership_activities ───── activities
```

---

# 1. `persons`

## Responsabilité

> Qui est cette personne ?

```text
persons
- id
- first_name
- last_name
- birth_date
- phone_number
- email
- address
- created_at
```

Exemple :

```text
Alice Martin
12/05/1994
alice@example.com
...
```

Les informations liées à une adhésion, à un compte utilisateur ou à un essai **ne doivent pas être stockées ici**.

---

# 2. `seasons`

## Responsabilité

> À quelle saison administrative appartient une adhésion ?

```text
seasons
- id
- name
- starts_at
- ends_at
- is_active
- created_at
```

Exemple :

```text
2026-2027
01/09/2026
31/08/2027
active = true
```

Le nom de la saison est unique.

---

# 3. `membership_types`

## Responsabilité

> Quel type d'adhésion possède le membre ?

```text
membership_types
- id
- name
- is_active
- created_at
```

Exemples :

```text
Adulte
Enfant
Étudiant
Bienfaiteur
```

Ces valeurs sont **configurables par l'association**.

On évite donc de coder en dur dans Go :

```text
JJB adulte
JJB enfant
```

---

# Activité ≠ type d'adhésion

Une distinction importante est apparue pendant la conception.

```text
Type d'adhésion
────────────────
Adulte
Enfant
Étudiant
```

n'est pas la même chose que :

```text
Activité
────────
Jujitsu traditionnel
Jujitsu brésilien
```

Ces deux notions sont donc séparées.

---

# 4. `activities`

## Responsabilité

> Quelles activités l'association propose-t-elle ?

```text
activities
- id
- name
- is_active
- created_at
```

Exemples pour un club de jujitsu :

```text
Jujitsu traditionnel
Jujitsu brésilien
```

Mais Club Manager reste générique :

```text
Judo
Yoga
Escalade
Danse
Aïkido
...
```

---

# 5. `memberships`

## Responsabilité

> Quelle est l'adhésion de cette personne pour cette saison ?

```text
memberships
- id
- person_id
- season_id
- membership_type_id
- status
- joined_at
- ended_at
- created_at
- updated_at
```

Relations :

```text
person_id
→ persons.id

season_id
→ seasons.id

membership_type_id
→ membership_types.id
```

---

## Pourquoi `memberships` plutôt que `members` ?

Être membre n'est pas une propriété permanente de la personne.

Une personne peut avoir :

```text
Alice
├── adhésion 2025-2026
├── adhésion 2026-2027
└── adhésion 2028-2029
```

Chaque adhésion possède son propre historique.

On ne modifie donc pas une ancienne adhésion pour représenter une nouvelle saison.

On **crée une nouvelle `membership`**.

---

## `status`

Les statuts sont contrôlés par Club Manager :

```text
pending
active
ended
cancelled
```

Ils ne sont pas librement configurables par l'administrateur car le fonctionnement de l'application peut dépendre de leur signification.

Exemple :

```sql
status TEXT NOT NULL
    CHECK (
        status IN (
            'pending',
            'active',
            'ended',
            'cancelled'
        )
    )
```

---

## Contrainte d'unicité

```sql
UNIQUE (person_id, season_id)
```

Dans le modèle actuel :

> Une personne ne peut avoir qu'une seule adhésion pour une même saison.

---

# 6. `membership_activities`

Une adhésion peut donner accès à plusieurs activités.

Exemple :

```text
Adhésion de Lucas
│
├── Jujitsu traditionnel
└── Jujitsu brésilien
```

On utilise donc une **table de liaison** :

```text
membership_activities
- membership_id
- activity_id
```

Relations :

```text
membership_id
→ memberships.id

activity_id
→ activities.id
```

Sa clé primaire est directement :

```sql
PRIMARY KEY (
    membership_id,
    activity_id
)
```

Il n'est pas nécessaire d'ajouter un `id`.

Cela empêche également :

```text
membership 12 → JJB
membership 12 → JJB
```

d'être enregistré deux fois.

---

# 7. `trial_registrations`

## Responsabilité

> Quels essais une personne a-t-elle réservés ou effectués ?

```text
trial_registrations
- id
- person_id
- activity_id
- trial_date
- status
- created_at
```

Relations :

```text
person_id
→ persons.id

activity_id
→ activities.id
```

---

## Statuts d'un essai

```text
registered
attended
cancelled
no_show
```

Ils répondent à :

> Que s'est-il passé avec cette inscription ?

Ils ne décrivent pas l'activité.

Exemple :

```text
Lucas
│
├── activité : Jujitsu brésilien
├── date : 18/09/2026
└── status : attended
```

---

## Nombre d'essais

On ne stocke pas :

```text
number_of_trials
```

car cette information existe déjà indirectement dans la base.

Elle est calculée avec les inscriptions :

```text
COUNT(trial_registrations)
```

Cela évite une redondance pouvant devenir incohérente.

La limite d'essais sera une **règle métier**.

Par exemple :

```text
maximum = 2 essais
```

ou éventuellement :

```text
2 essais maximum
par personne et par activité
```

---

# 8. `users`

## Responsabilité

> Cette personne possède-t-elle un compte permettant d'utiliser Club Manager ?

```text
users
- id
- person_id
- login_email
- password_hash
- is_active
- created_at
```

Relation :

```text
person_id
→ persons.id
```

---

## Email de contact ≠ email de connexion

On distingue :

```text
persons.email
→ moyen de contacter la personne
```

et :

```text
users.login_email
→ identifiant d'authentification
```

C'est notamment utile pour les enfants pouvant partager l'adresse électronique d'un parent.

`login_email` doit être unique.

---

# 9. `roles`

## Responsabilité

> Quels rôles existent dans Club Manager ?

```text
roles
- id
- name
```

Exemples :

```text
admin
president
secretary
treasurer
```

Le nom du rôle est unique.

---

# 10. `user_roles`

Un utilisateur peut posséder plusieurs rôles.

Exemple :

```text
Alice
├── president
└── treasurer
```

On utilise donc une nouvelle table de liaison :

```text
user_roles
- user_id
- role_id
```

Relations :

```text
user_id
→ users.id

role_id
→ roles.id
```

La clé primaire est :

```sql
PRIMARY KEY (
    user_id,
    role_id
)
```

Là encore, aucun `id` supplémentaire n'est nécessaire.

---

# Clés étrangères

Une clé étrangère relie une ligne à une autre table.

Exemple :

```sql
person_id INTEGER NOT NULL
    REFERENCES persons(id)
```

Cela signifie :

```text
memberships.person_id
        │
        ▼
persons.id
```

PostgreSQL garantit alors qu'une adhésion ne peut pas être créée pour une personne inexistante.

---

# Tables de liaison

Deux tables de liaison sont apparues :

```text
user_roles
```

et :

```text
membership_activities
```

Elles servent à représenter des relations **plusieurs-à-plusieurs**.

Exemple :

```text
User
 │
 ├── admin
 └── treasurer
```

et :

```text
Membership
 │
 ├── traditionnel
 └── brésilien
```

Le couple des deux clés étrangères peut directement devenir la clé primaire.

---

# Ordre des migrations Goose

Une table référencée doit exister **avant** la table qui possède la clé étrangère.

L'ordre retenu est donc :

```text
0001_persons.sql

0002_seasons.sql

0003_membership_types.sql

0004_activities.sql

0005_memberships.sql

0006_trial_registrations.sql

0007_users.sql

0008_roles.sql

0009_user_roles.sql

0010_membership_activities.sql
```

Les tables de liaison arrivent après les tables qu'elles relient.

---

# Pourquoi repartir de migrations propres ?

La première migration `create_members` avait été réalisée dans un but pédagogique.

Elle a permis d'apprendre :

```text
PostgreSQL
Goose
sqlc
GET
POST
CreateMember
```

Le modèle ayant changé profondément, elle n'a plus de raison d'appartenir au schéma final.

Le projet étant encore en développement, nous pouvons repartir d'une base vierge avec :

```text
0001_persons
...
0010_membership_activities
```

Git conserve malgré tout l'ancien travail et permet de voir cette évolution.

---

# Responsabilité de Go et PostgreSQL

Une décision importante a été prise concernant la validation et la normalisation.

Nous ne voulons pas que Go et PostgreSQL effectuent tous les deux les mêmes transformations.

Sinon :

```text
Go transforme
+
PostgreSQL transforme
```

et lorsqu'un problème apparaît, il devient plus difficile de savoir quelle couche en est responsable.

---

## Go : mise en forme et validation

Go s'occupera notamment de :

```text
suppression des espaces inutiles
normalisation des emails
validation du formulaire
validation des dates
détection des doublons probables
messages d'erreur
```

Exemple :

```go
strings.TrimSpace(value)
```

ou :

```go
strings.ToLower(
    strings.TrimSpace(email),
)
```

---

## PostgreSQL : intégrité structurelle

PostgreSQL garantit le modèle.

Exemples :

```text
PRIMARY KEY
FOREIGN KEY
NOT NULL
UNIQUE
CHECK sur les statuts
```

La séparation peut être résumée ainsi :

> **Go décide à quoi doivent ressembler les données.**
> 
> **PostgreSQL garantit que le modèle relationnel reste cohérent.**

---

# Doublons de personnes

On ne met pas :

```sql
UNIQUE (
    first_name,
    last_name,
    birth_date
)
```

car deux personnes différentes peuvent théoriquement avoir le même :

```text
nom
prénom
date de naissance
```

On utilisera plutôt ces informations pour **détecter un doublon probable** avant de créer une personne.

```text
CreatePerson
      │
      ▼
recherche
nom + prénom + naissance
      │
      ├── trouvé
      │      ↓
      │   vérifier / réutiliser
      │
      └── absent
             ↓
         création
```

---

# Principe de conception retenu

Chaque table doit répondre à **une question précise**.

```text
persons
→ Qui est cette personne ?

memberships
→ Quelle est son adhésion ?

seasons
→ Pour quelle saison ?

membership_types
→ Quel type d'adhésion ?

activities
→ Quelle activité ?

trial_registrations
→ Quel essai a-t-elle effectué ?

users
→ Peut-elle se connecter ?

roles
→ Quels rôles possède-t-elle ?
```

Cette séparation permet d'éviter une énorme table contenant des dizaines de colonnes sans rapport direct les unes avec les autres.

---

# Évolution future possible

Le modèle peut ensuite accueillir naturellement :

```text
memberships
├── payments
├── licences
└── documents

persons
└── person_relationships
     ├── responsable légal
     └── contact d'urgence

roles
└── role_permissions
     └── permissions
```

Ces notions ne sont pas nécessaires pour construire le socle actuel.

---

# Comprendre et retenir

## Une personne n'est pas un membre

```text
Person
→ identité

Membership
→ adhésion pendant une période
```

---

## Une personne n'est pas un utilisateur

```text
Person
→ personne réelle

User
→ compte informatique
```

---

## Un essai n'est pas une personne

```text
Person
→ Lucas

TrialRegistration
→ Lucas vient essayer le JJB le 18 septembre
```

---

## Une activité n'est pas un type d'adhésion

```text
Activity
→ Jujitsu brésilien

MembershipType
→ Enfant
```

---

## Une table de liaison représente une relation

```text
membership_activities
→ quelles activités appartiennent à une adhésion ?

user_roles
→ quels rôles possède un utilisateur ?
```

---

## Une clé étrangère protège les relations

```text
person_id REFERENCES persons(id)
```

signifie :

> Cette ligne doit obligatoirement correspondre à une personne existante.

---

## Le nombre d'essais n'est pas stocké

Il est calculé à partir des données existantes :

```text
COUNT(trial_registrations)
```

---

## Go et PostgreSQL n'ont pas la même responsabilité

```text
Go
→ nettoyer
→ normaliser
→ valider
→ expliquer les erreurs

PostgreSQL
→ garantir les relations
→ garantir les contraintes fortes
→ protéger l'intégrité du modèle
```

---

# Résumé du jalon

Nous sommes passés de :

```text
members
├── prénom
├── nom
├── naissance
└── email
```

à :

```text
                  persons
                     │
       ┌─────────────┼──────────────┐
       │             │              │
       ▼             ▼              ▼
     users      memberships   trial_registrations
       │             │              │
       │        ┌────┼────┐         │
       ▼        ▼    ▼    ▼         ▼
 user_roles  season type activities activities
       │
       ▼
     roles
```

Ce jalon marque le passage :

```text
apprendre à enregistrer une ligne
                ↓
concevoir un véritable modèle relationnel
                ↓
construire le socle de Club Manager
```