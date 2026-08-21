

---


# Jalon — Refonte du modèle de données et validation de `Person`

## Situation du projet

Club Manager disposait initialement d'une table :

```text
members
```

Cette table avait été volontairement simple.

Elle nous a permis d'apprendre et de mettre en place :

- PostgreSQL ;
- Goose ;
- les migrations ;
- sqlc ;
- pgx ;
- les formulaires `POST` ;
- les handlers HTTP ;
- les interfaces ;
- les fakes pour les tests ;
- l'insertion de données depuis Go.

Cette première architecture a donc parfaitement rempli son rôle pédagogique.

Mais en faisant évoluer le projet, une limite est apparue :

> une personne et une adhésion ne représentent pas la même chose.

Une personne peut exister dans l'association sans être membre.

Elle peut par exemple :

- effectuer une séance d'essai ;
- avoir été membre une année précédente ;
- devenir membre plus tard ;
- posséder un compte utilisateur ;
- exercer un rôle dans l'association.

La table `members` mélangeait donc progressivement plusieurs concepts différents.

Ce jalon marque le passage d'un modèle simple centré sur `members` vers un véritable modèle relationnel centré sur `persons`.

---

# 1. Le changement conceptuel principal

L'ancien raisonnement était proche de :

```text
Member
=
une personne inscrite dans le club
```

Le nouveau raisonnement sépare :

```text
Person
=
identité d'une personne
```

et :

```text
Membership
=
relation entre une personne
et une saison d'adhésion
```

Ainsi :

```text
Person ≠ Membership
```

Une `Person` est relativement permanente.

Une `Membership` est temporelle.

---

# 2. Le concept métier de membre reste valable

La disparition de la table :

```text
members
```

ne signifie pas que le concept de membre disparaît.

Au contraire.

Un membre devient un **concept métier composé de plusieurs données**.

Par exemple :

```text
Person
    +
Membership active
    =
Member
```

On peut donc continuer à parler :

- d'un membre ;
- de la liste des membres ;
- de créer un membre ;

sans qu'une table PostgreSQL `members` soit nécessaire.

---

# 3. Nouveau modèle relationnel

Le modèle retenu est :

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
     └── membership_activities ───── activities
```

Ce modèle sépare plusieurs responsabilités.

---

# 4. Table `persons`

La table `persons` représente l'identité d'une personne.

Elle contient notamment :

```text
id
first_name
last_name
birth_date
phone_number
email
address
created_at
```

Les données suivantes sont obligatoires :

```text
first_name
last_name
birth_date
```

Les coordonnées sont optionnelles :

```text
phone_number
email
address
```

Une personne peut donc parfaitement exister sans adresse email.

C'est notamment important pour :

- les enfants ;
- les personnes partageant une adresse familiale ;
- les personnes ne souhaitant pas communiquer certains renseignements.

---

# 5. Pas d'unicité artificielle sur l'identité

Nous avons choisi de ne pas imposer :

```text
UNIQUE(first_name, last_name, birth_date)
```

car deux personnes différentes peuvent réellement posséder :

```text
même prénom
même nom
même date de naissance
```

PostgreSQL ne doit donc pas empêcher ce cas.

Une éventuelle détection des doublons probables pourra être réalisée plus tard dans la logique Go.

Principe :

> une contrainte SQL doit représenter une véritable impossibilité métier, pas simplement un cas improbable.

---

# 6. Table `seasons`

Une adhésion appartient à une saison.

La table `seasons` contient notamment :

```text
id
name
starts_at
ends_at
is_active
created_at
```

Exemple :

```text
Saison 2026-2027
```

La saison est séparée de l'adhésion car plusieurs adhésions utilisent la même saison.

---

# 7. Table `membership_types`

Le type d'adhésion est également séparé.

Exemples :

```text
Adult
Child
Student
Benefactor
```

La table contient notamment :

```text
id
name
is_active
created_at
```

Un type d'adhésion ne représente pas une activité sportive.

Il décrit la nature administrative ou tarifaire de l'adhésion.

---

# 8. Table `activities`

Les activités sont regroupées dans leur propre catalogue.

Exemples :

```text
Jujitsu traditionnel
Jujitsu brésilien
```

Structure :

```text
id
name
is_active
created_at
```

Une activité peut être utilisée aussi bien pour :

- une adhésion ;
- une séance d'essai ;
- d'autres fonctionnalités futures.

---

# 9. Table `memberships`

La table `memberships` représente réellement l'adhésion.

Elle relie :

```text
Person
+
Season
+
MembershipType
```

Elle contient notamment :

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

Les statuts actuellement prévus sont :

```text
pending
active
ended
cancelled
```

La contrainte :

```text
UNIQUE(person_id, season_id)
```

garantit actuellement qu'une personne ne possède qu'une adhésion par saison.

---

# 10. Relation entre adhésion et activités

Une adhésion peut concerner plusieurs activités.

Une activité peut concerner plusieurs adhésions.

Il s'agit donc d'une relation :

```text
many-to-many
```

représentée par :

```text
membership_activities
```

avec :

```text
membership_id
activity_id
```

et une clé primaire composite :

```text
PRIMARY KEY (
    membership_id,
    activity_id
)
```

---

# 11. Séances d'essai

Une personne peut participer à une séance d'essai sans devenir membre.

Il était donc important de ne pas rattacher directement les essais à `memberships`.

La table :

```text
trial_registrations
```

relie :

```text
Person
+
Activity
+
Date
```

avec notamment les statuts :

```text
registered
attended
cancelled
no_show
```

Cela permet par exemple :

```text
Person
    ↓
TrialRegistration
    ↓
éventuellement
    ↓
Membership
```

---

# 12. Comptes utilisateurs

Une personne et un utilisateur de l'application sont également deux concepts différents.

La table :

```text
users
```

est reliée à :

```text
persons
```

Une personne peut donc exister sans disposer d'un compte applicatif.

Le compte contient notamment :

```text
person_id
login_email
password_hash
is_active
```

Une distinction importante apparaît alors :

```text
Person.Email
```

représente une coordonnée de contact.

Alors que :

```text
User.LoginEmail
```

représente une identité d'authentification.

Ces deux données n'ont pas forcément la même responsabilité.

---

# 13. Rôles

Les rôles ne sont pas placés directement dans `persons`.

Ils sont liés aux comptes utilisateurs.

Structure :

```text
users
    ↓
user_roles
    ↓
roles
```

Cela permettra plus tard d'avoir par exemple :

```text
Admin
Président
Secrétaire
Trésorier
```

La réflexion détaillée sur les permissions est volontairement reportée.

---

# 14. Ce qui est volontairement reporté

Certaines fonctionnalités identifiées ne sont pas encore implémentées.

Notamment :

```text
payments
licences
documents
guardian relationships
emergency contacts
permissions
role_permissions
modules configurables
tables configurables
```

Principe retenu :

> construire uniquement ce dont nous avons besoin maintenant, mais éviter une architecture qui empêcherait les évolutions futures.

---

# 15. Répartition des responsabilités entre Go et PostgreSQL

Un choix important a été fait sur la responsabilité des validations.

## Go

Go prend en charge notamment :

```text
TrimSpace
validation des formulaires
conversion des dates
conversion des champs optionnels
normalisation éventuelle des emails
messages d'erreur utilisateur
détection future de doublons probables
```

## PostgreSQL

PostgreSQL garantit :

```text
PRIMARY KEY
FOREIGN KEY
NOT NULL
UNIQUE
CHECK
cohérence relationnelle
```

Résumé :

> Go décide à quoi doivent ressembler les données.

> PostgreSQL garantit que le modèle relationnel reste cohérent.

---

# 16. Ne pas dupliquer inutilement les validations

Par exemple, nous pourrions écrire en SQL :

```sql
CHECK (TRIM(first_name) <> '')
```

Mais le handler Go vérifie déjà ce cas.

Nous avons donc choisi de ne pas nécessairement dupliquer chaque règle dans les deux couches.

Il faut distinguer :

```text
validation utilisateur
```

et :

```text
intégrité structurelle
```

---

# 17. Reconstruction des migrations

Comme le projet est encore en développement et qu'aucune donnée importante n'avait besoin d'être conservée, il a été décidé de repartir proprement.

La base de développement a été recréée.

Les nouvelles migrations sont :

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

Cette reconstruction évite de maintenir artificiellement un ancien modèle devenu obsolète.

Git conserve de toute façon l'historique pédagogique du projet.

---

# 18. Validation des migrations avec Goose

Les migrations sont d'abord vérifiées avec :

```bash
goose -dir migrations validate
```

Une validation réussie peut ne produire aucun message.

Puis elles sont appliquées :

```bash
goose -dir migrations postgres \
"host=localhost port=5432 user=club_manager dbname=club_manager sslmode=require" up
```

La base est alors migrée jusqu'à :

```text
version 10
```

---

# 19. Vérification PostgreSQL

Les tables ont ensuite été inspectées directement avec `psql`.

Connexion :

```bash
psql -h localhost -U club_manager -d club_manager
```

Liste des tables :

```text
\dt
```

Inspection d'une table :

```text
\d memberships
```

ou :

```text
\d persons
```

Cela permet de vérifier que les migrations produisent réellement le schéma attendu.

---

# 20. Refonte des requêtes sqlc

L'ancienne requête :

```text
CreateMember
```

ne correspondait plus au nouveau modèle.

Elle n'a donc pas été simplement renommée.

Créer une `Person` n'est pas équivalent à créer un `Member`.

L'ancienne requête a été supprimée et remplacée par :

```sql
-- name: CreatePerson :one
INSERT INTO persons (
    first_name,
    last_name,
    birth_date,
    phone_number,
    email,
    address
)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;
```

Puis :

```bash
sqlc generate
```

---

# 21. Types générés par sqlc

sqlc génère notamment :

```go
type Person struct {
    ID          int32
    FirstName   string
    LastName    string
    BirthDate   pgtype.Date
    PhoneNumber pgtype.Text
    Email       pgtype.Text
    Address     pgtype.Text
    CreatedAt   pgtype.Timestamptz
}
```

et :

```go
type CreatePersonParams struct {
    FirstName   string
    LastName    string
    BirthDate   pgtype.Date
    PhoneNumber pgtype.Text
    Email       pgtype.Text
    Address     pgtype.Text
}
```

Ainsi que :

```go
func (q *Queries) CreatePerson(
    ctx context.Context,
    arg CreatePersonParams,
) (Person, error)
```

---

# 22. Nettoyage de l'ancien modèle `Member`

Après la régénération sqlc, l'ancien modèle `Member` existait encore dans plusieurs parties du code.

Nous avons volontairement utilisé :

```bash
go test ./...
```

pour laisser le compilateur nous montrer progressivement les dépendances restantes.

Cela a permis de supprimer ou adapter :

```text
interfaces
handlers
tests
routes
fakes
ancien code généré
```

Principe intéressant :

> le compilateur peut servir de guide pendant une refactorisation.

---

# 23. Évolution de l'interface `Queries`

Le routeur et les handlers ne dépendent pas directement de :

```go
*dbsqlc.Queries
```

Ils dépendent d'une interface définie par l'application.

Par exemple :

```go
type Queries interface {
    CreatePerson(
        ctx context.Context,
        arg dbsqlc.CreatePersonParams,
    ) (dbsqlc.Person, error)
}
```

Ainsi :

```text
handler
   ↓
database.Queries
```

et non :

```text
handler
   ↓
implémentation PostgreSQL concrète
```

Cela permet notamment d'utiliser :

```text
vrai dbsqlc.Queries
```

en production et :

```text
FakeQueries
```

dans les tests.

---

# 24. Helpers pour `pgtype`

sqlc utilise les types pgx pour représenter certains champs PostgreSQL.

Par exemple :

```go
pgtype.Date
pgtype.Text
```

Des helpers ont donc été créés pour centraliser les conversions.

---

# 25. Helper `pgTypeDate`

Un champ HTML `date` retourne une chaîne telle que :

```text
2020-03-15
```

Go utilise le layout :

```go
"2006-01-02"
```

Le helper transforme cette chaîne en :

```go
pgtype.Date
```

Principe :

```go
func pgTypeDate(value string) (pgtype.Date, error) {
    t, err := time.Parse("2006-01-02", value)

    if err != nil {
        return pgtype.Date{}, fmt.Errorf(
            "conversion de la date %q impossible : %w",
            value,
            err,
        )
    }

    return pgtype.Date{
        Time:  t,
        Valid: true,
    }, nil
}
```

---

# 26. Helper `pgTypeText`

Les champs optionnels doivent pouvoir devenir :

```sql
NULL
```

et pas nécessairement :

```text
""
```

Le helper réalise :

```text
chaîne non vide
→ pgtype.Text{Valid:true}
```

et :

```text
chaîne vide
→ pgtype.Text{Valid:false}
→ NULL PostgreSQL
```

Exemple :

```go
func pgTypeText(value string) pgtype.Text {
    value = strings.TrimSpace(value)

    if value == "" {
        return pgtype.Text{
            Valid: false,
        }
    }

    return pgtype.Text{
        String: value,
        Valid:  true,
    }
}
```

---

# 27. Tests des helpers

Les helpers ont été testés séparément.

Pour `pgTypeDate` :

```text
date valide
→ conversion correcte

date invalide
→ erreur
```

Pour `pgTypeText` :

```text
texte normal
→ Valid = true

espaces autour
→ TrimSpace

texte vide
→ Valid = false
→ NULL
```

Cette étape illustre le principe retenu :

```text
faire une petite fonction
        ↓
la tester
        ↓
la valider
        ↓
passer à la suivante
```

---

# 28. Nouveau `PostPersonHandler`

Le handler de création d'une personne remplace progressivement l'ancien handler `Member`.

Structure générale :

```go
func PostPersonHandler(
    queries database.Queries,
) http.HandlerFunc {

    return func(
        w http.ResponseWriter,
        r *http.Request,
    ) {
        // ...
    }
}
```

Le handler :

1. lit les valeurs du formulaire ;
2. nettoie prénom et nom ;
3. vérifie les champs obligatoires ;
4. convertit la date ;
5. convertit les champs optionnels ;
6. construit `CreatePersonParams` ;
7. appelle `CreatePerson`.

---

# 29. Validation du prénom et du nom

PostgreSQL possède :

```sql
NOT NULL
```

mais :

```text
""
```

n'est pas `NULL`.

Le handler doit donc vérifier :

```go
firstName := strings.TrimSpace(
    r.FormValue("FirstName"),
)

lastName := strings.TrimSpace(
    r.FormValue("LastName"),
)

if firstName == "" || lastName == "" {
    http.Error(
        w,
        "Nom et prénom obligatoires",
        http.StatusBadRequest,
    )

    return
}
```

---

# 30. Conversion de la date

La date du formulaire est convertie avant l'appel sqlc :

```go
birthDate, err := pgTypeDate(
    r.FormValue("Birthdate"),
)

if err != nil {
    http.Error(
        w,
        "Date de naissance invalide",
        http.StatusBadRequest,
    )

    return
}
```

Une date invalide est une erreur utilisateur :

```text
HTTP 400
```

et non une panne du serveur.

---

# 31. Construction de `CreatePersonParams`

Les données sont ensuite regroupées :

```go
datas := dbsqlc.CreatePersonParams{
    FirstName: firstName,
    LastName:  lastName,

    BirthDate: birthDate,

    PhoneNumber: pgTypeText(
        r.FormValue("PhoneNumber"),
    ),

    Email: pgTypeText(
        r.FormValue("Email"),
    ),

    Address: pgTypeText(
        r.FormValue("Address"),
    ),
}
```

Puis :

```go
_, err = queries.CreatePerson(
    r.Context(),
    datas,
)
```

---

# 32. Pourquoi utiliser `r.Context()`

Le contexte de la requête HTTP est transmis jusqu'à PostgreSQL.

```text
HTTP Request
    ↓
r.Context()
    ↓
Handler
    ↓
sqlc
    ↓
pgx
    ↓
PostgreSQL
```

Si la requête HTTP est annulée, les couches inférieures peuvent également recevoir cette information.

---

# 33. Tests du `PostPersonHandler`

Un fake spécifique enregistre les appels reçus.

Exemple conceptuel :

```go
type recordingQueries struct {
    CreatePersonParams dbsqlc.CreatePersonParams
    CreatePersonCalled bool
    CreatePersonError  error
}
```

La méthode :

```go
func (q *recordingQueries) CreatePerson(
    ctx context.Context,
    arg dbsqlc.CreatePersonParams,
) (dbsqlc.Person, error) {

    q.CreatePersonParams = arg
    q.CreatePersonCalled = true

    return dbsqlc.Person{},
        q.CreatePersonError
}
```

---

# 34. Pourquoi enregistrer les paramètres dans le fake ?

Cela permet de vérifier que le handler transmet réellement :

```text
les bonnes valeurs
```

et pas seulement qu'une méthode a été appelée.

Exemple :

```text
formulaire
FirstName = "   Robin "

handler
TrimSpace

fake reçoit
FirstName = "Robin"
```

Le test vérifie donc aussi la transformation effectuée par le handler.

---

# 35. Scénarios testés sur le handler

Quatre chemins principaux sont maintenant couverts.

## Cas nominal

Formulaire valide :

```text
→ CreatePerson appelé
→ paramètres corrects
→ HTTP 200 actuellement
```

---

## Nom ou prénom absent

```text
→ HTTP 400
→ CreatePerson NON appelé
```

C'est important.

Un test ne doit pas seulement vérifier le résultat visible.

Il peut également vérifier :

> ce qui ne doit surtout pas se produire.

---

## Date invalide

```text
→ HTTP 400
→ CreatePerson NON appelé
```

---

## Erreur de base simulée

Le fake peut être configuré :

```go
CreatePersonError:
    errors.New("database error")
```

Cela ne provoque pas réellement une erreur PostgreSQL.

Le fake se comporte simplement comme si la base avait échoué.

Le handler doit alors répondre :

```text
HTTP 500
```

et :

```text
CreatePersonCalled = true
```

---

# 36. Deuxième niveau de tests : intégration PostgreSQL

Les tests du handler utilisent un fake.

Ils vérifient donc :

```text
HTTP
 ↓
Handler
 ↓
FakeQueries
```

Ils ne vérifient pas :

```text
sqlc
pgx
PostgreSQL
```

Un test d'intégration a donc été ajouté.

---

# 37. Nouvelle requête `GetPersonByID`

Pour vérifier réellement ce qui a été écrit dans PostgreSQL, une requête de lecture a été ajoutée :

```sql
-- name: GetPersonByID :one
SELECT *
FROM persons
WHERE id = $1;
```

Puis :

```bash
sqlc generate
```

---

# 38. Convention de nommage des requêtes

Le nom :

```text
GetPersonByID
```

a été préféré à :

```text
GetPerson
```

car il indique explicitement le critère de recherche.

Cela prépare une convention future :

```text
GetPersonByID
GetPersonByEmail

ListPersonsByLastName
ListPersonsByBirthDate
```

Principe :

```text
Get
→ un élément

List
→ plusieurs éléments possibles

By
→ critère de recherche
```

---

# 39. Pourquoi ne pas utiliser `GetLastPerson`

Une première idée aurait pu être :

```sql
SELECT *
FROM persons
ORDER BY id DESC
LIMIT 1;
```

Mais cela revient à demander :

```text
donne-moi la dernière personne
```

et non :

```text
donne-moi exactement
celle que je viens de créer
```

`CreatePerson` retourne déjà la ligne créée et donc :

```go
createdPerson.ID
```

Le test peut utiliser précisément :

```go
GetPersonByID(
    ctx,
    createdPerson.ID,
)
```

Cette méthode est beaucoup plus robuste.

---

# 40. Emplacement du test d'intégration

Le test est placé dans :

```text
internal/database/
    persons_integration_test.go
```

et non dans :

```text
internal/database/dbsqlc/
```

car `dbsqlc` contient du code généré.

Organisation :

```text
internal/database/dbsqlc
→ code généré

internal/database
→ notre code
→ interfaces
→ connexion
→ tests d'intégration
```

---

# 41. Différence avec un test de handler

Test de handler :

```text
HTTP
 ↓
Handler
 ↓
FakeQueries
```

Test d'intégration :

```text
Test
 ↓
sqlc
 ↓
pgx
 ↓
PostgreSQL réel
```

Les deux types de tests ne remplacent pas l'un l'autre.

Ils répondent à des questions différentes.

---

# 42. Connexion PostgreSQL depuis le test

Le test utilise la vraie fonction :

```go
database.New(...)
```

avec :

```text
DATABASE_URL
```

Exemple :

```go
ctx := context.Background()

databaseURL :=
    os.Getenv("DATABASE_URL")

if databaseURL == "" {
    t.Fatal(
        "DATABASE_URL doit être définie pour le test d'intégration",
    )
}

db, err := New(
    ctx,
    databaseURL,
)

if err != nil {
    t.Fatalf(
        "connexion à PostgreSQL impossible : %v",
        err,
    )
}

defer db.Close()
```

---

# 43. Attention à l'environnement de l'IDE

Le premier lancement du test depuis l'IDE a échoué avec :

```text
DATABASE_URL doit être définie
pour le test d'intégration
```

Alors que le test fonctionnait depuis le terminal.

La raison est qu'une variable définie par :

```bash
export DATABASE_URL='...'
```

est transmise aux processus lancés depuis ce terminal.

Un IDE déjà démarré ne possède pas nécessairement cette variable.

Exemple :

```text
terminal
│
├── export DATABASE_URL
│
└── go test
      ✅
```

mais :

```text
terminal
└── export DATABASE_URL

IDE déjà lancé
└── Run test
      DATABASE_URL absente
```

Pour l'instant, les tests d'intégration sont donc lancés simplement depuis le terminal.

---

# 44. Transaction de test

Le test démarre une transaction :

```go
tx, err := db.Begin(ctx)

if err != nil {
    t.Fatalf(
        "impossible de démarrer la transaction : %v",
        err,
    )
}
```

Puis prépare son annulation :

```go
defer func() {
    _ = tx.Rollback(ctx)
}()
```

---

# 45. Comprendre une transaction

Une transaction permet de regrouper plusieurs opérations SQL :

```text
BEGIN
  ↓
INSERT
  ↓
SELECT
  ↓
assertions
  ↓
ROLLBACK
```

Pendant la transaction, les modifications existent réellement.

Le test peut donc relire la personne qu'il vient de créer.

---

# 46. Utilisation de `WithTx`

sqlc génère :

```go
func (
    q *Queries,
) WithTx(
    tx pgx.Tx,
) *Queries
```

On peut donc écrire :

```go
queries := dbsqlc.New(db)

queries =
    queries.WithTx(tx)
```

À partir de là :

```go
queries.CreatePerson(...)
queries.GetPersonByID(...)
```

sont exécutés dans la transaction.

---

# 47. Test d'intégration

Le principe est :

```go
createdPerson, err :=
    queries.CreatePerson(
        ctx,
        params,
    )

if err != nil {
    t.Fatalf(...)
}

savedPerson, err :=
    queries.GetPersonByID(
        ctx,
        createdPerson.ID,
    )

if err != nil {
    t.Fatalf(...)
}
```

Puis les données sont comparées :

```go
if savedPerson.FirstName !=
    params.FirstName {

    t.Errorf(...)
}
```

et :

```go
if savedPerson.LastName !=
    params.LastName {

    t.Errorf(...)
}
```

---

# 48. Pourquoi relire après insertion ?

Un :

```text
INSERT réussi
```

prouve que PostgreSQL a accepté la commande.

Mais le test veut vérifier davantage :

```text
ce qui a réellement été enregistré
```

La stratégie devient :

```text
INSERT
  ↓
récupérer ID
  ↓
SELECT par ID
  ↓
comparer
```

Cela vérifie ensemble :

```text
CreatePerson
sqlc
pgx
PostgreSQL
GetPersonByID
```

---

# 49. `ROLLBACK`

À la fin du test :

```go
tx.Rollback(ctx)
```

annule les modifications de la transaction.

Exemple :

```text
avant le test

persons
→ vide
```

Pendant le test :

```text
persons
→ Robin Des Bois
```

Puis :

```text
ROLLBACK
```

Après :

```text
persons
→ vide
```

---

# 50. Vérification manuelle du rollback

Après le test, connexion à PostgreSQL :

```bash
psql -h localhost \
    -U club_manager \
    -d club_manager
```

Puis :

```sql
SELECT * FROM persons;
```

Résultat obtenu :

```text
 id | first_name | last_name | ...
----+------------+-----------+-----
(0 rows)
```

Le test a donc réellement :

```text
inséré
↓
relu
↓
vérifié
↓
annulé
```

La base n'est pas polluée par les tests.

---

# 51. Pourquoi `defer` pour le rollback ?

Le rollback est placé dans :

```go
defer func() {
    _ = tx.Rollback(ctx)
}()
```

Cela garantit son exécution lorsque la fonction de test se termine.

Même si le test rencontre :

```go
t.Fatalf(...)
```

les fonctions différées avec `defer` sont exécutées avant la sortie.

C'est une sécurité importante pour le nettoyage.

---

# 52. Identifiants et rollback

Un rollback annule l'insertion.

Il ne garantit pas pour autant que les identifiants générés soient réutilisés.

Exemple :

```text
ID 1
ID 2

test :
ID 3
ROLLBACK

prochaine insertion :
ID 4
```

Il peut donc exister des trous :

```text
1
2
4
5
```

Ce n'est pas un problème.

Une clé primaire doit être :

```text
unique
```

pas :

```text
continue
```

---

# 53. Résultat de `go test ./...`

L'ensemble des tests passe :

```text
cmd/server
[no test files]

internal/config
ok

internal/database
ok

internal/database/dbsqlc
[no test files]

internal/handlers
ok

internal/router
ok

internal/views
[no test files]
```

Le projet est donc à nouveau entièrement vert après la refonte.

---

# 54. Les niveaux de tests actuels

Nous disposons maintenant de plusieurs niveaux complémentaires.

## Niveau 1 — fonctions isolées

Exemples :

```text
pgTypeDate
pgTypeText
```

Chaîne :

```text
fonction
 ↓
assertions
```

---

## Niveau 2 — handlers

Chaîne :

```text
HTTP
 ↓
Handler
 ↓
FakeQueries
```

Permet de tester :

```text
validation
codes HTTP
paramètres
appels
erreurs
```

sans PostgreSQL.

---

## Niveau 3 — intégration PostgreSQL

Chaîne :

```text
sqlc
 ↓
pgx
 ↓
PostgreSQL
```

Permet de vérifier que les vraies couches techniques fonctionnent ensemble.

---

# 55. Architecture obtenue

L'architecture de la création d'une personne est désormais :

```text
HTML Form
    ↓
HTTP POST
    ↓
Router
    ↓
PostPersonHandler
    ↓
database.Queries
    ↓
dbsqlc.Queries
    ↓
pgx
    ↓
PostgreSQL
```

Chaque niveau possède une responsabilité identifiable.

---

# 56. Intérêt de l'interface `Queries`

L'interface permet :

```text
PostPersonHandler
       ↓
database.Queries
       ↓
       ├── dbsqlc.Queries
       │      production
       │
       └── recordingQueries
              tests
```

Le handler ne connaît pas l'implémentation concrète.

Il connaît seulement le contrat :

```text
CreatePerson(...)
```

Cela améliore :

```text
testabilité
découplage
maintenabilité
lisibilité
```

---

# 57. La testabilité influence l'architecture

Une observation importante de ce jalon est que les tests ne servent pas uniquement à détecter des bugs.

Le besoin de tester facilement pousse naturellement vers :

```text
dépendances explicites
interfaces
petites fonctions
responsabilités séparées
```

On obtient donc :

```text
testabilité
    ↓
architecture plus claire
    ↓
code plus facilement maintenable
```

---

# 58. Méthode de développement qui émerge

Une manière de travailler s'est naturellement imposée :

```text
faire la fonction
        ↓
la tester
        ↓
corriger
        ↓
valider
        ↓
faire la suivante
```

Ce fonctionnement est proche d'un développement incrémental orienté tests.

Ce n'est pas nécessairement du TDD strict :

```text
test avant code
```

mais plutôt :

```text
petite évolution
+
validation immédiate
```

Cette approche évite :

```text
5 fonctionnalités
↓
50 erreurs
↓
où est le problème ?
```

et favorise :

```text
1 petite évolution
↓
test
↓
vert
↓
suite
```

---

# 59. Les tests comme documentation

Les tests décrivent également le comportement attendu.

Par exemple :

```text
prénom manquant
→ 400
→ aucune écriture DB
```

ou :

```text
date invalide
→ 400
→ aucune écriture DB
```

ou :

```text
erreur DB
→ 500
```

Un lecteur peut donc utiliser les tests pour comprendre le contrat d'un handler.

---

# 60. Tests et bonnes pratiques

Le nombre de tests n'est pas, à lui seul, une preuve de qualité.

Ce qui est intéressant ici est que les tests vérifient des comportements précis.

Par exemple :

```text
entrée valide
→ appel DB
```

```text
entrée invalide
→ pas d'appel DB
```

```text
erreur DB
→ réponse adaptée
```

```text
insertion réelle
→ lecture réelle
→ valeurs correctes
```

Ils ne vérifient donc pas uniquement :

```text
ça compile
```

mais bien :

```text
le comportement attendu
```

---

# 61. Ancien et nouveau modèle

Le chemin parcouru peut être résumé ainsi :

```text
members
   ↓
premières migrations
   ↓
CreateMember
   ↓
POST membre
   ↓
tests
   ↓
compréhension des limites du modèle
   ↓
refonte métier
   ↓
persons
   ↓
memberships
   ↓
activities
   ↓
users
   ↓
roles
   ↓
CreatePerson
   ↓
tests unitaires
   ↓
test d'intégration PostgreSQL
```

L'ancien modèle n'était donc pas une erreur.

Il constituait une étape d'apprentissage ayant permis de comprendre pourquoi le nouveau modèle était nécessaire.

---

# 62. Principe important retenu

> Ne pas chercher à concevoir immédiatement l'architecture parfaite.

Nous avons d'abord construit quelque chose de simple :

```text
Member
```

Puis son utilisation réelle a révélé les concepts cachés :

```text
Person
Membership
Season
Activity
User
Role
```

La nouvelle architecture résulte donc d'un besoin compris grâce à l'ancienne.

---

# 63. Autre principe important

> Une refactorisation n'est pas simplement un renommage.

Nous n'avons pas remplacé mécaniquement :

```text
Member
```

par :

```text
Person
```

car :

```text
CreateMember
≠
CreatePerson
```

Créer une personne ne crée pas encore une adhésion.

Plus tard :

```text
Créer un membre
```

pourra signifier quelque chose comme :

```text
trouver/créer Person
        ↓
créer Membership
        ↓
associer Activities
```

---

# 64. État actuel du concept `Person`

La création d'une personne est maintenant validée à plusieurs niveaux :

```text
pgTypeDate
→ testé
✅

pgTypeText
→ testé
✅

PostPersonHandler
→ testé avec fake
✅

CreatePerson
→ PostgreSQL réel
✅

GetPersonByID
→ PostgreSQL réel
✅

transaction
→ testée
✅

rollback
→ vérifié manuellement
✅

go test ./...
→ vert
✅
```

---

# 65. Ce que ce jalon garantit

À partir de maintenant, nous savons que :

```text
une Person peut être représentée en PostgreSQL
```

```text
les migrations créent correctement son modèle
```

```text
sqlc génère les bons types
```

```text
Go peut construire les paramètres
```

```text
le handler sait valider un formulaire
```

```text
CreatePerson fonctionne
```

```text
GetPersonByID fonctionne
```

```text
les erreurs sont testées
```

```text
PostgreSQL peut être testé sans polluer la base
```

---

# 66. Ce que ce jalon ne fait pas encore

Il ne crée pas encore un membre au sens métier.

Il crée uniquement :

```text
Person
```

Il reste ensuite à construire progressivement :

```text
Membership
```

puis les relations nécessaires.

C'est volontaire.

Nous validons d'abord :

```text
Person
```

avant d'empiler :

```text
Membership
Activities
Season
MembershipType
```

---

# 67. Prochaine direction

La base étant maintenant solide, le travail peut continuer avec la même méthode :

```text
faire une petite brique
        ↓
tester
        ↓
valider
        ↓
continuer
```

La prochaine grande étape logique sera autour de :

```text
Membership
```

mais sans remettre en cause `Person`.

Le modèle commencera alors réellement à permettre :

```text
Person
   ↓
Membership
   ↓
Member
```

---

# 68. Comprendre et retenir

> Une personne est une identité ; une adhésion est une relation temporelle.

> `Person` et `Member` ne sont donc pas synonymes dans le modèle de données.

> Le concept métier de membre peut exister sans table `members`.

> Les tables doivent représenter des concepts cohérents et séparés.

> PostgreSQL protège l'intégrité relationnelle.

> Go gère principalement la validation et la préparation des données.

> sqlc transforme les requêtes SQL en méthodes Go typées.

> Une interface permet au code métier de dépendre d'un contrat plutôt que d'une implémentation concrète.

> Les fakes permettent de tester les handlers sans PostgreSQL.

> Un test d'intégration vérifie plusieurs composants réels ensemble.

> Pour tester une insertion, il est pertinent d'insérer puis de relire exactement la ligne grâce à son ID.

> Une transaction permet de tester réellement PostgreSQL.

> `ROLLBACK` permet ensuite d'annuler les modifications du test.

> `defer` garantit le nettoyage lorsque le test se termine.

> Les trous dans les identifiants après un rollback sont normaux.

> Un identifiant doit être unique, pas continu.

> Le compilateur et les tests sont de très bons guides pendant une refactorisation.

> Une ancienne architecture simple n'est pas nécessairement mauvaise : elle peut permettre de comprendre pourquoi une architecture plus riche devient nécessaire.

> Une bonne évolution consiste à construire une petite brique, la tester, la valider, puis seulement construire la suivante.

---

# 69. Jalon atteint

Avant :

```text
Member
 ↓
table unique
 ↓
responsabilités progressivement mélangées
```

Maintenant :

```text
Person
│
├── future Membership
├── TrialRegistration
└── User

Membership
├── Season
├── MembershipType
└── Activities
```

Et pour `Person` :

```text
HTML
 ↓
HTTP
 ↓
Handler
 ↓
Interface Queries
 ↓
sqlc
 ↓
pgx
 ↓
PostgreSQL
```

avec :

```text
tests helpers
      +
tests handler
      +
fake database
      +
test d'intégration
      +
transaction
      +
rollback
```

Le modèle `Person` constitue maintenant une fondation validée sur laquelle les prochaines fonctionnalités peuvent être construites.

---

# 70. Git — commit du jalon

Après vérification :

```bash
go test ./...
```

puis :

```bash
git status
```

le jalon peut être enregistré avec :

```bash
git add .
git commit -m "complete person data model and integration tests"
```

Puis :

```bash
git push
```

---

# 71. Git — tag du jalon

Ce jalon représente plus qu'un simple ajout de test.

Il marque :

```text
la refonte du modèle relationnel
+
le remplacement de Member par Person comme fondation
+
la validation complète de Person
```

Un tag peut donc être créé.

Exemple :

```bash
git tag -a v0.2.0-person-model \
    -m "Person data model and PostgreSQL integration validated"
```

Puis :

```bash
git push origin v0.2.0-person-model
```

Le tag permet de revenir précisément à cet état avant de commencer la construction des relations métier suivantes.

---

# État final du jalon

```text
Refonte du modèle de données
✅

Person séparée de Membership
✅

nouveau schéma PostgreSQL
✅

10 migrations Goose
✅

sqlc régénéré
✅

ancien modèle Member nettoyé
✅

interface Queries adaptée
✅

helpers pgtype
✅

tests helpers
✅

PostPersonHandler
✅

tests handler
✅

GetPersonByID
✅

test d'intégration PostgreSQL
✅

transaction
✅

rollback
✅

go test ./...
✅
```

## Le projet peut maintenant passer au jalon suivant.

```text
Person
  ✅
  ↓
Membership
  prochain objectif
```