
# Jalon 11 — Afficher la liste des membres

## Objectif du jalon

Afficher dans Club Manager l'ensemble des membres enregistrés dans PostgreSQL.

Le parcours complet devient :

```text
PostgreSQL
    ↓
sqlc
    ↓
[]dbsqlc.Member
    ↓
MembersHandler
    ↓
[]views.MemberData
    ↓
template HTML
    ↓
tableau des membres
```

Ce jalon complète le précédent :

```text
POST /members
→ création d'un membre

GET /members
→ consultation des membres
```

---

# 1. La requête SQL

Dans :

```text
internal/database/queries/members.sql
```

on ajoute une requête permettant de récupérer tous les membres :

```sql
-- name: ListMembers :many
SELECT
    id,
    first_name,
    last_name,
    birth_date,
    email,
    created_at
FROM members
ORDER BY last_name, first_name;
```

Le point important est :

```text
:many
```

Contrairement à une requête retournant un seul élément, `sqlc` génère ici une fonction retournant une slice :

```go
func (q *Queries) ListMembers(
    ctx context.Context,
) ([]Member, error)
```

Donc :

```text
un membre
→ dbsqlc.Member

plusieurs membres
→ []dbsqlc.Member
```

---

# 2. Code généré par sqlc

Après :

```bash
sqlc generate
```

sqlc génère notamment une boucle permettant de construire la slice :

```go
var items []Member

for rows.Next() {

    var i Member

    if err := rows.Scan(
        &i.ID,
        &i.FirstName,
        &i.LastName,
        &i.BirthDate,
        &i.Email,
        &i.CreatedAt,
    ); err != nil {
        return nil, err
    }

    items = append(items, i)
}
```

## Ce qu'il faut comprendre

```go
rows.Next()
```

avance jusqu'à la prochaine ligne du résultat SQL.

```go
var i Member
```

crée un membre temporaire.

```go
rows.Scan(&i.ID, ...)
```

écrit les données SQL directement dans les champs du membre.

Les `&` sont nécessaires car `Scan` doit connaître **l'adresse mémoire** des variables qu'il doit modifier.

Enfin :

```go
items = append(items, i)
```

ajoute le membre à la slice.

---

# 3. Faire évoluer l'interface `Queries`

Le handler ne dépend pas directement de :

```go
*dbsqlc.Queries
```

mais de notre interface :

```go
database.Queries
```

On y ajoute donc :

```go
ListMembers(
    ctx context.Context,
) ([]dbsqlc.Member, error)
```

L'interface contient maintenant les opérations nécessaires sur les membres :

```go
type Queries interface {

    CreateMember(
        ctx context.Context,
        arg dbsqlc.CreateMemberParams,
    ) (dbsqlc.Member, error)

    ListMembers(
        ctx context.Context,
    ) ([]dbsqlc.Member, error)
}
```

## Conséquence sur les tests

Toute structure utilisée comme faux `Queries` doit maintenant également posséder :

```go
ListMembers(...)
```

Même lorsqu'un test particulier ne l'utilise pas.

Cela rappelle une propriété importante des interfaces Go :

> Une structure satisfait une interface uniquement si elle possède toutes les méthodes demandées par cette interface.

---

# 4. Le handler récupère les membres

Le handler appelle :

```go
members, err := queries.ListMembers(r.Context())
```

`members` contient donc :

```go
[]dbsqlc.Member
```

Si PostgreSQL retourne une erreur :

```go
if err != nil {

    http.Error(
        w,
        "impossible d'afficher la liste des adhérents",
        http.StatusInternalServerError,
    )

    return
}
```

---

# 5. Ne pas envoyer directement les données SQL à la vue

Un `dbsqlc.Member` contient des types liés à PostgreSQL :

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

La vue n'a pas besoin de connaître :

```text
pgtype.Date
pgtype.Text
pgtype.Timestamptz
Valid
Time
```

On crée donc une structure destinée uniquement à l'affichage :

```go
type MemberData struct {
    ID        int32
    FirstName string
    LastName  string
    BirthDate string
    Email     string
    CreatedAt string
}
```

La vue reçoit alors uniquement des types simples.

---

# 6. Transformer les données dans le handler

Le handler joue ici le rôle d'adaptateur :

```text
données techniques
      ↓
   handler
      ↓
données d'affichage
```

On crée :

```go
var memberData []views.MemberData
```

puis on parcourt les membres :

```go
for _, member := range members {
```

Pour chaque membre, on construit un :

```go
views.MemberData
```

---

# 7. Représenter correctement un email absent

Dans PostgreSQL :

```text
email existant
→ TEXT

email absent
→ NULL
```

Il est important de conserver cette distinction.

On ne stocke donc pas :

```text
"Aucun mail"
```

dans PostgreSQL.

`"Aucun mail"` est uniquement une **représentation destinée à l'utilisateur**.

## Lors de la création

Dans le POST :

```go
email := r.FormValue("Email")
```

puis :

```go
Email: pgtype.Text{
    String: email,
    Valid:  email != "",
},
```

Si :

```go
email == ""
```

alors :

```go
Valid == false
```

et PostgreSQL reçoit :

```sql
NULL
```

---

# 8. Transformer `NULL` pour la vue

Dans la liste :

```go
for _, member := range members {

    email := "Aucun mail"

    if member.Email.Valid {
        email = member.Email.String
    }

    // ...
}
```

Ainsi :

```text
PostgreSQL NULL
    ↓
pgtype.Text{Valid: false}
    ↓
handler
    ↓
"Aucun mail"
    ↓
HTML
```

## Pourquoi `email` est dans la boucle ?

Il faut écrire :

```go
for _, member := range members {

    email := "Aucun mail"
```

et non :

```go
email := "Aucun mail"

for _, member := range members {
```

Sinon la valeur modifiée pour un membre pourrait être conservée pour le membre suivant.

Chaque itération doit repartir de :

```text
"Aucun mail"
```

---

# 9. Formater les dates

PostgreSQL/sqlc fournit notamment :

```go
member.BirthDate.Time
```

et :

```go
member.CreatedAt.Time
```

## Date de naissance

```go
member.BirthDate.Time.Format("02/01/2006")
```

Exemple :

```text
12/05/1990
```

## Date de création

```go
member.CreatedAt.Time.Format("02/01/2006 15:04")
```

Exemple :

```text
15/08/2026 10:30
```

La date de naissance est une simple date.

La création correspond à un instant précis : afficher l'heure est donc pertinent.

---

# 10. Construction des données de la vue

La transformation complète ressemble à :

```go
for _, member := range members {

    email := "Aucun mail"

    if member.Email.Valid {
        email = member.Email.String
    }

    memberData = append(
        memberData,
        views.MemberData{
            ID:        member.ID,
            FirstName: member.FirstName,
            LastName:  member.LastName,
            BirthDate: member.BirthDate.Time.Format("02/01/2006"),
            Email:     email,
            CreatedAt: member.CreatedAt.Time.Format("02/01/2006 15:04"),
        },
    )
}
```

Puis :

```go
data := views.MembersData{
    SiteName: cfg.SiteName,
    Title:    "Membres - " + cfg.SiteName,
    Members:  memberData,
}
```

---

# 11. Afficher les membres dans le template

La page reçoit :

```go
type MembersData struct {
    SiteName string
    Title    string
    Members  []MemberData
}
```

Dans le template :

```html
{{ range .Members }}
```

est l'équivalent de :

```go
for _, member := range members {
```

---

## Le changement de contexte de `.`

Avant le `range` :

```text
.
→ MembersData
```

Donc :

```html
{{ .Members }}
```

accède à la slice.

À l'intérieur :

```html
{{ range .Members }}
```

le `.` devient le `MemberData` courant.

On peut donc écrire :

```html
<td>{{ .FirstName }}</td>
<td>{{ .LastName }}</td>
<td>{{ .BirthDate }}</td>
<td>{{ .Email }}</td>
<td>{{ .CreatedAt }}</td>
```

---

# 12. Le tableau HTML

La page peut maintenant afficher :

```html
<table class="table">

    <thead>
        <tr>
            <th>ID</th>
            <th>Prénom</th>
            <th>Nom</th>
            <th>Date de naissance</th>
            <th>Email</th>
            <th>Créé le</th>
        </tr>
    </thead>

    <tbody>

        {{ range .Members }}

        <tr>
            <td>{{ .ID }}</td>
            <td>{{ .FirstName }}</td>
            <td>{{ .LastName }}</td>
            <td>{{ .BirthDate }}</td>
            <td>{{ .Email }}</td>
            <td>{{ .CreatedAt }}</td>
        </tr>

        {{ end }}

    </tbody>

</table>
```

La vue ne fait pratiquement aucune logique.

Elle reçoit des données déjà prêtes à afficher.

---

# 13. Tester `MembersHandler`

Le test utilise un faux `Queries`.

On lui donne directement les membres qu'il devra retourner :

```go
type recordingQueries struct {
    Members []dbsqlc.Member
}
```

Puis :

```go
func (q *recordingQueries) ListMembers(
    ctx context.Context,
) ([]dbsqlc.Member, error) {

    return q.Members, nil
}
```

---

# 14. Le faux doit reproduire les données de la base

Il ne faut pas donner directement :

```go
Email: "robin@example.com"
```

car `dbsqlc.Member` utilise :

```go
pgtype.Text
```

On simule donc réellement ce que PostgreSQL/sqlc retournerait :

```go
Email: pgtype.Text{
    String: "robin.desbois@example.com",
    Valid:  true,
},
```

Pour un email absent :

```go
Email: pgtype.Text{
    Valid: false,
},
```

Cela simule :

```sql
NULL
```

---

# 15. Construire des dates dans le test

Pour la date de naissance :

```go
BirthDate: pgtype.Date{
    Time: time.Date(
        1990,
        5,
        12,
        0,
        0,
        0,
        0,
        time.UTC,
    ),
    Valid: true,
},
```

Pour `CreatedAt` :

```go
CreatedAt: pgtype.Timestamptz{
    Time: time.Date(
        2009,
        10,
        10,
        15,
        4,
        0,
        0,
        time.UTC,
    ),
    Valid: true,
},
```

---

# 16. Tester la réponse HTTP

On crée une requête :

```go
request := httptest.NewRequest(
    http.MethodGet,
    "/members",
    nil,
)
```

et un faux `ResponseWriter` :

```go
response := httptest.NewRecorder()
```

Puis :

```go
MembersHandler(
    cfg,
    queries,
)(response, request)
```

On vérifie le statut :

```go
if response.Code != http.StatusOK {

    t.Errorf(
        "statut obtenu : %d, statut attendu : %d",
        response.Code,
        http.StatusOK,
    )
}
```

---

# 17. Tester le HTML généré

On récupère :

```go
body := response.Body.String()
```

Puis on peut vérifier différentes transformations :

```go
strings.Contains(body, "Robin")
```

vérifie qu'un membre est présent.

```go
strings.Contains(body, "12/05/1990")
```

vérifie le formatage de la date de naissance.

```go
strings.Contains(
    body,
    "robin.desbois@example.com",
)
```

vérifie un email existant.

```go
strings.Contains(body, "Aucun mail")
```

vérifie la transformation :

```text
NULL
→ Aucun mail
```

Enfin :

```go
strings.Contains(
    body,
    "10/10/2009 15:04",
)
```

vérifie le formatage de `CreatedAt`.

---

# 18. Le test a réellement trouvé une erreur

Lors du développement, le handler utilisait :

```go
member.CreatedAt.Time.Format("02/01/2006")
```

alors que le comportement attendu était :

```go
member.CreatedAt.Time.Format(
    "02/01/2006 15:04",
)
```

Le site fonctionnait visuellement, mais le test a signalé :

```text
la réponse ne contient pas la date de création formatée
```

Le test a donc rempli son rôle :

> vérifier automatiquement le comportement attendu et détecter une régression ou une erreur difficile à remarquer visuellement.

---

# Architecture obtenue

```text
                    GET /members
                         │
                         ▼
                  MembersHandler
                         │
                         │ ListMembers()
                         ▼
                 database.Queries
                         │
                         ▼
                       sqlc
                         │
                         ▼
                    PostgreSQL
                         │
                         ▼
                 []dbsqlc.Member
                         │
                         │ transformation
                         ▼
                []views.MemberData
                         │
                         ▼
                 MembersData
                         │
                         ▼
                    template
                         │
                         ▼
                 tableau HTML
```

---

# Séparation des responsabilités

## PostgreSQL

Conserve la donnée réelle :

```text
email absent
→ NULL
```

---

## sqlc

Transforme les lignes SQL en structures Go :

```text
PostgreSQL
→ dbsqlc.Member
```

---

## Handler

Récupère et prépare les données destinées à l'affichage :

```text
pgtype.Text
→ string

NULL
→ "Aucun mail"

time.Time
→ "12/05/1990"
```

---

## View

Décrit les données nécessaires à la page :

```go
MemberData
MembersData
```

---

## Template

Se contente essentiellement d'afficher :

```html
{{ .FirstName }}
{{ .LastName }}
{{ .Email }}
```

---

## Test

Vérifie que toutes ces transformations produisent le résultat attendu.

---

# Ce que ce jalon introduit

Ce jalon permet de travailler plusieurs concepts importants ensemble :

- requête SQL retournant plusieurs lignes ;
    
- `sqlc` avec `:many` ;
    
- slices ;
    
- boucle `range` en Go ;
    
- `rows.Next()` ;
    
- `rows.Scan()` et adresses mémoire ;
    
- extension d'une interface ;
    
- faux objet pour les tests ;
    
- séparation modèle SQL / modèle de vue ;
    
- gestion de `NULL` avec `pgtype.Text.Valid` ;
    
- formatage des dates avec `time.Format()` ;
    
- boucle `range` dans `html/template` ;
    
- changement du contexte `.` dans un template ;
    
- test HTTP avec `httptest`;
    
- vérification du HTML généré.
    

---

# Comprendre et retenir

## 1. La base conserve la donnée, pas son affichage

```text
NULL
```

signifie :

> aucune valeur enregistrée.

```text
"Aucun mail"
```

signifie :

> représentation choisie pour l'utilisateur.

Les deux ne doivent pas être confondus.

---

## 2. Le handler fait le lien entre deux mondes

```text
modèle base de données
        ↓
     handler
        ↓
modèle de vue
```

Le handler peut donc transformer :

```text
pgtype.Text
→ string
```

ou :

```text
time.Time
→ date lisible
```

---

## 3. Une vue doit recevoir des données simples

Il est préférable que le template connaisse :

```go
Email string
```

plutôt que :

```go
Email pgtype.Text
```

La vue ne doit pas avoir besoin de comprendre PostgreSQL.

---

## 4. `range` parcourt une collection

En Go :

```go
for _, member := range members {
```

Dans un template :

```html
{{ range .Members }}
```

Le principe est le même.

---

## 5. Un test doit simuler l'entrée réelle

Le handler reçoit normalement :

```go
dbsqlc.Member
```

Le faux doit donc fournir des `dbsqlc.Member`, avec leurs vrais types :

```go
pgtype.Date
pgtype.Text
pgtype.Timestamptz
```

Puis le test vérifie la transformation effectuée par le handler.

---

# Résultat du jalon

Club Manager sait maintenant :

```text
créer un membre
        +
récupérer les membres
        +
transformer leurs données
        +
afficher la liste
        +
tester cet affichage
```

Le cycle commence donc à ressembler à un véritable CRUD :

```text
Create ✅
Read   ✅

Update ⏳
Delete ⏳
```

Le socle permettant de gérer réellement les membres est maintenant en place.