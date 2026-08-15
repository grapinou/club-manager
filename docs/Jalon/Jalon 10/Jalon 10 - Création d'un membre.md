
---


## Objectif du jalon

Ce jalon marque la mise en place complète de la **création d'un membre** dans Club Manager.

Pour la première fois, une donnée saisie par un utilisateur dans une page HTML traverse toute l'application jusqu'à PostgreSQL.

Le chemin complet est maintenant :

```text
formulaire HTML
      ↓
requête HTTP POST
      ↓
routeur
      ↓
handler
      ↓
interface Queries
      ↓
sqlc
      ↓
pgx
      ↓
PostgreSQL
```

Puis, après l'insertion :

```text
PostgreSQL
      ↓
303 See Other
      ↓
GET /members/new
      ↓
formulaire
```

Ce jalon correspond donc au **Create** du CRUD :

```text
CRUD

C → Create   ← jalon actuel
R → Read
U → Update
D → Delete
```

---

# 1. Point de départ

La table `members` existe déjà dans PostgreSQL.

Elle contient notamment :

```text
id
first_name
last_name
birth_date
email
created_at
```

La création de cette table a été gérée par une migration Goose.

À partir de cette base, l'objectif était de permettre à l'application Go d'insérer réellement un membre.

---

# 2. La requête SQL

La création d'un membre commence par une requête SQL.

Conceptuellement :

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
RETURNING *;
```

Cette requête exprime directement l'opération métier :

> créer un nouveau membre dans PostgreSQL.

---

# 3. sqlc

Plutôt que d'écrire nous-mêmes tout le code permettant d'exécuter cette requête, nous utilisons `sqlc`.

`sqlc` analyse :

```text
requêtes SQL
+
schéma PostgreSQL
```

et génère du code Go typé.

On obtient notamment des éléments du type :

```go
CreateMemberParams
```

et une méthode :

```go
CreateMember(...)
```

La chaîne devient :

```text
SQL écrit par le développeur
        ↓
sqlc
        ↓
code Go généré
        ↓
pgx
        ↓
PostgreSQL
```

---

# 4. Pourquoi utiliser une interface devant sqlc ?

Le handler ne dépend pas directement du type concret généré par sqlc.

Il reçoit :

```go
database.Queries
```

Conceptuellement :

```go
func PostMemberHandler(queries database.Queries) http.HandlerFunc
```

Le handler sait uniquement que l'objet reçu possède la méthode nécessaire :

```go
CreateMember(...)
```

Il n'a pas besoin de savoir s'il s'agit :

- du vrai objet sqlc ;
- d'un faux objet utilisé dans un test ;
- d'une autre implémentation future.

Architecture :

```text
PostMemberHandler
       ↓
database.Queries
       ↓
   interface
    ↙      ↘
 sqlc      test
 réel      fake
```

C'est une application concrète des interfaces Go.

---

# 5. Injection de la dépendance

Dans `main`, le vrai objet permettant de travailler avec PostgreSQL est créé.

Conceptuellement :

```go
queries := dbsqlc.New(db)
```

Puis il est transmis au routeur :

```go
mux := router.New(cfg, queries)
```

Le routeur le transmet ensuite au handler :

```go
handlers.PostMemberHandler(queries)
```

La dépendance suit donc ce chemin :

```text
main
 ↓
router
 ↓
handler
```

Le handler ne crée pas lui-même sa connexion PostgreSQL.

Cette organisation permet de garder les responsabilités séparées.

---

# 6. La route POST

La route permettant de créer un membre est :

```text
POST /members
```

Dans le routeur :

```go
mux.HandleFunc(
    "POST /members",
    handlers.PostMemberHandler(queries),
)
```

Le routeur possède ici une responsabilité simple :

> lorsque je reçois un `POST /members`, j'appelle le handler correspondant.

---

# 7. Le handler de création

Le handler reçoit les données du formulaire grâce à :

```go
r.FormValue("FirstName")
r.FormValue("LastName")
r.FormValue("Birthdate")
r.FormValue("Email")
```

Il transforme ensuite ces données pour construire :

```go
dbsqlc.CreateMemberParams
```

Conceptuellement :

```text
requête HTTP
    ↓
r.FormValue(...)
    ↓
valeurs Go
    ↓
CreateMemberParams
    ↓
queries.CreateMember(...)
```

---

# 8. Conversion de la date

Une donnée provenant d'un formulaire HTTP est initialement du texte.

Par exemple :

```text
1990-05-12
```

Pour PostgreSQL, la date doit être convertie.

Nous utilisons :

```go
time.Parse(
    "2006-01-02",
    r.FormValue("Birthdate"),
)
```

Le format :

```text
2006-01-02
```

est le format de référence utilisé par Go pour représenter :

```text
année-mois-jour
```

Puis la valeur est placée dans un type PostgreSQL compatible :

```go
pgtype.Date{
    Time:  birthDate,
    Valid: true,
}
```

---

# 9. Les types PostgreSQL de pgx

Certaines données sont représentées avec des types spécifiques à `pgx`.

Par exemple :

```go
pgtype.Date
```

et :

```go
pgtype.Text
```

Ces structures permettent notamment de représenter une valeur et son état de validité.

Exemple :

```go
pgtype.Text{
    String: email,
    Valid:  true,
}
```

Cela permet également de gérer plus tard les valeurs SQL :

```text
NULL
```

de manière explicite.

---

# 10. Le test du handler

Avant de brancher réellement PostgreSQL, le comportement du handler a été testé.

Nous avons créé :

```go
recordingQueries
```

Cette structure joue le rôle d'un faux objet `Queries`.

Elle possède une méthode :

```go
CreateMember(...)
```

compatible avec l'interface attendue par le handler.

Son rôle n'est pas d'écrire dans PostgreSQL.

Son rôle est d'enregistrer ce que le handler essaie d'envoyer à la base.

Conceptuellement :

```text
PostMemberHandler
       ↓
recordingQueries
       ↓
enregistre les paramètres reçus
       ↓
test
```

Cela permet de vérifier :

```text
FirstName
LastName
BirthDate
Email
```

sans utiliser une véritable base de données.

---

# 11. Les pointeurs dans le test

Le faux objet est créé avec :

```go
queries := &recordingQueries{}
```

`&` signifie :

> récupérer l'adresse de cette valeur.

Le type de `queries` devient donc :

```go
*recordingQueries
```

La méthode utilise également un receiver pointeur :

```go
func (q *recordingQueries) CreateMember(...)
```

Cela permet de modifier la structure originale :

```go
q.CreateMemberParams = arg
```

Après l'appel du handler, le test peut donc consulter :

```go
queries.CreateMemberParams
```

et vérifier ce que le handler a réellement envoyé.

C'est un exemple concret de l'intérêt des pointeurs :

> plusieurs parties du programme travaillent sur le même objet plutôt que sur des copies indépendantes.

---

# 12. Simulation du formulaire avec `httptest`

Le formulaire envoyé au handler est construit avec :

```go
form := url.Values{}
```

Puis :

```go
form.Set("FirstName", firstName)
form.Set("LastName", lastName)
form.Set("Birthdate", birthdate)
form.Set("Email", email)
```

Le formulaire est ensuite encodé :

```go
form.Encode()
```

Ce qui produit quelque chose ressemblant à :

```text
FirstName=Robin&LastName=Des+Bois&Birthdate=1990-05-12&Email=...
```

---

# 13. Création de la requête POST de test

La requête est créée avec :

```go
request := httptest.NewRequest(
    http.MethodPost,
    "/members",
    strings.NewReader(form.Encode()),
)
```

Décomposition :

```text
http.MethodPost
        ↓
méthode POST

"/members"
        ↓
route appelée

form.Encode()
        ↓
formulaire transformé en texte

strings.NewReader(...)
        ↓
transforme la chaîne en lecteur
```

Le corps de la requête peut ainsi être lu comme un flux de données.

---

# 14. Le Content-Type

Le test ajoute :

```go
request.Header.Set(
    "Content-Type",
    "application/x-www-form-urlencoded",
)
```

Ce header signifie :

> le corps de cette requête contient les données d'un formulaire HTML classique.

Sans cette information, le serveur ne saurait pas nécessairement comment interpréter les données reçues.

---

# 15. Validation avec curl

Après le test du handler, nous avons testé le vrai serveur avec `curl`.

Commande utilisée :

```bash
curl -i -X POST http://localhost:8080/members \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "FirstName=Robin" \
  -d "LastName=Des Bois" \
  -d "Birthdate=1990-05-12" \
  -d "Email=robin.desbois@example.com"
```

`curl` a joué le rôle d'un client HTTP sans utiliser de navigateur.

---

# 16. Correspondance entre httptest et curl

Nous avons pu observer que les deux opérations étaient presque identiques :

```text
TEST GO                                  CURL

http.MethodPost                          -X POST

"/members"                               /members

form.Set("FirstName", "Robin")           -d "FirstName=Robin"

form.Encode()                            données envoyées avec -d

Content-Type                             -H "Content-Type: ..."
```

Cela permet de comprendre que `httptest` ne simule pas quelque chose d'abstrait.

Il construit réellement les différents éléments d'une requête HTTP.

---

# 17. Premier POST réel

Le serveur a répondu :

```text
HTTP/1.1 200 OK

post membre
```

Puis PostgreSQL a été interrogé avec :

```sql
SELECT * FROM members;
```

Résultat :

```text
 id | first_name | last_name | birth_date |           email           |          created_at
----+------------+-----------+------------+---------------------------+-------------------------------
  1 | Robin      | Des Bois  | 1990-05-12 | robin.desbois@example.com | 2026-08-14 10:22:20.847458+02
```

À cet instant, toute la chaîne backend était validée :

```text
curl
 ↓
HTTP
 ↓
router
 ↓
handler
 ↓
sqlc
 ↓
pgx
 ↓
PostgreSQL
```

---

# 18. Création du vrai formulaire HTML

Une page permettant de créer un membre a ensuite été ajoutée.

Le formulaire contient notamment :

```html
<form action="/members" method="post">
```

Cette ligne signifie :

> lorsque l'utilisateur valide le formulaire, envoyer une requête `POST` vers `/members`.

C'est l'équivalent HTML de :

```bash
curl -X POST http://localhost:8080/members
```

---

# 19. Le rôle de `name`

Un champ possède par exemple :

```html
<input
    type="text"
    id="firstName"
    name="FirstName"
>
```

Le point important pour le serveur est :

```html
name="FirstName"
```

Le handler cherche :

```go
r.FormValue("FirstName")
```

Il faut donc une correspondance exacte :

```text
FORMULAIRE                     HANDLER

name="FirstName"       →       FormValue("FirstName")

name="LastName"        →       FormValue("LastName")

name="Birthdate"       →       FormValue("Birthdate")

name="Email"           →       FormValue("Email")
```

Le champ `id` sert notamment au HTML et au lien avec le `label`.

Le champ `name` définit la donnée réellement envoyée au serveur.

---

# 20. Le navigateur remplace curl

Avec le formulaire HTML, le navigateur construit maintenant automatiquement la requête que nous avions créée avec `curl`.

```text
curl                            navigateur

-X POST                         method="post"

URL /members                    action="/members"

-d "FirstName=..."              name="FirstName"

-d "LastName=..."               name="LastName"

Content-Type                    automatique
```

Le navigateur ne fait donc pas quelque chose de magique.

Il automatise simplement la construction de la requête HTTP.

---

# 21. La vue du formulaire

Une vue dédiée au formulaire a été créée.

Organisation :

```text
internal/
└── views/
    ├── 008_member_form.go
    └── templates/
        └── pages/
            └── 008_member_form.html
```

La vue Go utilise le même principe que les autres pages :

```text
embed
 ↓
template.ParseFS
 ↓
RenderMemberForm
```

Elle utilise également le layout commun :

```text
base.html
```

---

# 22. Les données de la vue

La structure nécessaire au formulaire reste volontairement simple :

```go
type MemberFormData struct {
    SiteName string
    Title    string
}
```

Principe retenu :

> une vue reçoit uniquement les données dont elle a besoin.

---

# 23. Le handler GET du formulaire

Une route séparée permet d'afficher le formulaire :

```text
GET /members/new
```

Le handler :

1. reçoit la configuration ;
2. construit `MemberFormData` ;
3. appelle `RenderMemberForm`.

Architecture :

```text
GET /members/new
       ↓
MemberFormHandler
       ↓
MemberFormData
       ↓
RenderMemberForm
       ↓
HTML
```

---

# 24. Séparation GET / POST

Nous avons maintenant deux routes avec deux responsabilités différentes :

```text
GET /members/new
```

sert à :

> afficher le formulaire.

Alors que :

```text
POST /members
```

sert à :

> traiter le formulaire et créer le membre.

Cette séparation rend le fonctionnement particulièrement lisible.

```text
GET
 ↓
afficher

POST
 ↓
modifier les données
```

---

# 25. Validation avec le vrai navigateur

Le formulaire a ensuite été rempli directement dans le navigateur.

Le navigateur a envoyé :

```text
POST /members
```

Le membre est apparu dans PostgreSQL.

La chaîne complète de l'application a donc été validée :

```text
utilisateur
    ↓
formulaire HTML
    ↓
navigateur
    ↓
POST /members
    ↓
router
    ↓
PostMemberHandler
    ↓
database.Queries
    ↓
sqlc
    ↓
pgx
    ↓
PostgreSQL
```

---

# 26. Problème de la réponse directe après POST

Initialement, le handler terminait par :

```go
fmt.Fprintln(w, "post membre")
```

Le navigateur affichait donc simplement :

```text
post membre
```

Cela prouvait que la création fonctionnait, mais ce n'est pas un comportement adapté à une vraie application.

Un autre problème existe :

> si l'utilisateur rafraîchit une page résultant directement d'un POST, le navigateur peut proposer de renvoyer les données.

Cela pourrait entraîner une nouvelle insertion.

---

# 27. POST → Redirect → GET

La réponse directe a donc été remplacée par :

```go
http.Redirect(
    w,
    r,
    "/members/new",
    http.StatusSeeOther,
)
```

`http.StatusSeeOther` correspond au statut :

```text
303 See Other
```

Le serveur dit au navigateur :

> Le traitement du POST est terminé. Effectue maintenant un GET vers cette autre page.

Le fonctionnement devient :

```text
POST /members
     ↓
création du membre
     ↓
303 See Other
     ↓
GET /members/new
     ↓
affichage du formulaire
```

---

# 28. Pourquoi la redirection est importante ?

Sans redirection :

```text
POST
 ↓
page affichée directement
 ↓
rafraîchissement
 ↓
risque de renvoyer le POST
```

Avec redirection :

```text
POST
 ↓
modification des données
 ↓
redirect
 ↓
GET
 ↓
affichage
```

Le navigateur se retrouve finalement sur une requête `GET`.

Un rafraîchissement recharge donc simplement la page.

Ce principe est appelé :

```text
POST / Redirect / GET
```

ou :

```text
PRG
```

---

# 29. Responsabilités obtenues à la fin du jalon

## PostgreSQL

Responsabilité :

> stocker les membres.

---

## Goose

Responsabilité :

> faire évoluer le schéma de la base de données.

---

## SQL

Responsabilité :

> décrire les opérations effectuées sur les données.

---

## sqlc

Responsabilité :

> générer du code Go typé à partir des requêtes SQL.

---

## pgx

Responsabilité :

> communiquer avec PostgreSQL depuis Go.

---

## Interface `database.Queries`

Responsabilité :

> définir ce dont les handlers ont besoin pour accéder aux données.

---

## Handler POST

Responsabilité :

> recevoir et convertir les données du formulaire puis demander la création du membre.

---

## Handler GET

Responsabilité :

> préparer les données nécessaires à l'affichage du formulaire.

---

## Vue Go

Responsabilité :

> charger et exécuter les templates.

---

## Template HTML

Responsabilité :

> afficher le formulaire.

---

## Routeur

Responsabilité :

> relier une méthode et une URL au bon handler.

---

# 30. Architecture finale du jalon

```text
                         ┌─────────────────┐
                         │   Navigateur    │
                         └────────┬────────┘
                                  │
                     GET /members/new
                                  │
                                  ▼
                         ┌─────────────────┐
                         │     Router      │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │ Handler GET     │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │      View       │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │ Formulaire HTML │
                         └────────┬────────┘
                                  │
                            POST /members
                                  │
                                  ▼
                         ┌─────────────────┐
                         │     Router      │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │ Handler POST    │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │ database.Queries│
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │      sqlc       │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │       pgx       │
                         └────────┬────────┘
                                  │
                                  ▼
                         ┌─────────────────┐
                         │   PostgreSQL    │
                         └────────┬────────┘
                                  │
                            303 See Other
                                  │
                                  ▼
                          GET /members/new
```

---

# 31. Ce que ce jalon apporte par rapport aux précédents

Jusqu'ici, Club Manager savait principalement :

```text
recevoir une requête GET
        ↓
préparer des données
        ↓
afficher une page
```

À partir de ce jalon, l'application sait aussi :

```text
recevoir des données utilisateur
        ↓
les valider / convertir
        ↓
les transmettre à la couche SQL
        ↓
modifier durablement PostgreSQL
```

Il s'agit donc d'une évolution importante de l'architecture.

L'application n'est plus seulement capable **d'afficher des informations**.

Elle devient capable **de modifier son état persistant**.

---

# 32. Méthode utilisée pour construire ce jalon

Le développement a été réalisé progressivement.

### Étape 1

Construire la couche PostgreSQL.

### Étape 2

Créer la requête SQL.

### Étape 3

Générer la couche Go avec sqlc.

### Étape 4

Créer l'interface utilisée par le handler.

### Étape 5

Créer le handler POST.

### Étape 6

Tester le handler avec un faux `Queries`.

### Étape 7

Tester le vrai backend avec `curl`.

### Étape 8

Vérifier directement PostgreSQL.

### Étape 9

Créer le formulaire HTML.

### Étape 10

Créer la vue et le handler GET.

### Étape 11

Tester depuis le navigateur.

### Étape 12

Ajouter le pattern POST → Redirect → GET.

Cette progression a permis de tester chaque couche séparément avant de les assembler.

---

# 33. Intérêt de cette approche

Si le formulaire final ne fonctionnait pas, nous savions déjà que :

```text
PostgreSQL fonctionne
sqlc fonctionne
pgx fonctionne
le handler fonctionne
POST /members fonctionne
```

car cela avait été validé précédemment avec :

```text
httptest
```

puis :

```text
curl
```

On aurait donc pu concentrer la recherche du problème sur :

```text
formulaire
ou
navigateur
```

C'est un principe important du développement :

> construire et valider progressivement les différentes couches réduit fortement la zone dans laquelle rechercher une erreur.

---

# 34. Comprendre et retenir

La création d'un membre n'est pas une seule fonction.

C'est une chaîne de responsabilités :

```text
HTML
 ↓
HTTP
 ↓
router
 ↓
handler
 ↓
interface
 ↓
sqlc
 ↓
pgx
 ↓
PostgreSQL
```

Chaque couche possède un rôle limité.

Le handler ne connaît pas les détails de PostgreSQL.

Le template ne connaît pas sqlc.

Le routeur ne transforme pas les données.

La base de données ne connaît pas HTTP.

Cette séparation permet de garder une architecture compréhensible même lorsque l'application grandit.

---

# 35. Les notions importantes rencontrées

Ce jalon a permis d'utiliser concrètement :

- requête HTTP `POST` ;
- formulaire HTML ;
- `Content-Type` ;
- `application/x-www-form-urlencoded` ;
- `url.Values` ;
- `form.Encode()` ;
- `strings.NewReader()` ;
- `io.Reader` ;
- `httptest.NewRequest()` ;
- headers HTTP ;
- `curl` ;
- interfaces Go ;
- injection de dépendances ;
- faux objet de test ;
- pointeurs ;
- `&` ;
- receiver `*T` ;
- `time.Parse()` ;
- types `pgtype` ;
- sqlc ;
- pgx ;
- PostgreSQL ;
- templates Go ;
- `GET /members/new` ;
- `POST /members` ;
- `303 See Other` ;
- `http.Redirect()` ;
- pattern POST / Redirect / GET.

---

# 36. À retenir en une phrase

> La création d'un membre traverse désormais toutes les couches de Club Manager : le navigateur envoie un formulaire en `POST`, le handler convertit les données, l'interface permet de les transmettre à sqlc puis pgx jusqu'à PostgreSQL, avant qu'une redirection `303` ramène proprement l'utilisateur sur une requête `GET`.

---

# 37. État du projet après ce jalon

Le **Create** du CRUD des membres est maintenant fonctionnel.

```text
CREATE ✅

READ   → prochain jalon

UPDATE → plus tard

DELETE → plus tard
```

Le prochain jalon peut donc commencer sur une nouvelle responsabilité :

> récupérer les membres depuis PostgreSQL et les afficher dans une page.

Ce futur chemin sera approximativement :

```text
GET /members
      ↓
handler
      ↓
sqlc
      ↓
SELECT
      ↓
PostgreSQL
      ↓
[]Member
      ↓
vue
      ↓
{{ range }}
      ↓
liste HTML des membres
```

Ce travail constitue un nouveau jalon indépendant : **la lecture et l'affichage des membres**.

