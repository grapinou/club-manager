
---


# Test d'une requête POST avec `httptest`

## Objectif

Lorsqu'un handler reçoit un formulaire en `POST`, son test doit reproduire ce que ferait réellement un navigateur :

1. créer les données du formulaire ;
2. encoder ces données ;
3. créer une requête HTTP `POST` ;
4. placer les données dans le corps de la requête ;
5. préciser le type de contenu envoyé ;
6. envoyer la requête au handler ;
7. vérifier la réponse ;
8. vérifier ce que le handler a transmis à la couche suivante.

Dans Club Manager, le test de création d'un membre permet également de voir un usage concret des **pointeurs**.

---

# 1. Création du formulaire

On commence par créer les données qui seraient normalement envoyées par un formulaire HTML.

```go
form := url.Values{}

firstName := "Robin"
lastName := "Des Bois"
birthdate := "1990-05-12"
email := "robin.desbois@example.com"

form.Set("FirstName", firstName)
form.Set("LastName", lastName)
form.Set("Birthdate", birthdate)
form.Set("Email", email)
```

`url.Values` représente des couples clé / valeur adaptés aux paramètres d'une URL ou d'un formulaire.

Conceptuellement :

```text
FirstName = Robin
LastName  = Des Bois
Birthdate = 1990-05-12
Email     = robin.desbois@example.com
```

---

# 2. Encoder le formulaire

Un formulaire HTTP n'est pas envoyé sous la forme d'une structure Go.

Il doit être transformé en texte.

C'est le rôle de :

```go
form.Encode()
```

On obtient quelque chose qui ressemble à :

```text
FirstName=Robin&LastName=Des+Bois&Birthdate=1990-05-12&Email=robin.desbois%40example.com
```

Les différentes valeurs sont séparées par :

```text
&
```

et chaque clé est associée à sa valeur avec :

```text
=
```

Certains caractères sont également encodés.

Par exemple :

```text
Des Bois
```

peut devenir :

```text
Des+Bois
```

---

# 3. Construire la requête POST

La ligne importante du test est :

```go
request := httptest.NewRequest(
    http.MethodPost,
    "/members",
    strings.NewReader(form.Encode()),
)
```

Elle peut être lue comme :

> Créer une requête HTTP `POST` vers `/members`, dont le corps contient le formulaire encodé.

La version sur une seule ligne est :

```go
request := httptest.NewRequest(http.MethodPost, "/members", strings.NewReader(form.Encode()))
```

---

# 4. Les trois arguments de `httptest.NewRequest`

La fonction peut être comprise approximativement ainsi :

```go
httptest.NewRequest(method, target, body)
```

Dans notre cas :

```go
httptest.NewRequest(
    http.MethodPost,
    "/members",
    strings.NewReader(form.Encode()),
)
```

## Premier argument : la méthode

```go
http.MethodPost
```

correspond à :

```text
POST
```

On teste donc bien une requête d'envoi de données.

---

## Deuxième argument : la destination

```go
"/members"
```

correspond à la route appelée.

La requête simulée est donc :

```text
POST /members
```

---

## Troisième argument : le corps de la requête

```go
strings.NewReader(form.Encode())
```

Cette partie mérite d'être décomposée.

D'abord :

```go
form.Encode()
```

produit une `string`.

Par exemple :

```text
FirstName=Robin&LastName=Des+Bois
```

Mais `httptest.NewRequest` attend quelque chose capable d'être **lu comme un flux de données**.

On utilise donc :

```go
strings.NewReader(...)
```

qui transforme la chaîne en lecteur.

Conceptuellement :

```text
formulaire
    ↓
form.Encode()
    ↓
string
    ↓
strings.NewReader(...)
    ↓
lecteur
    ↓
corps de la requête HTTP
```

---

# 5. Pourquoi utiliser un Reader ?

Le corps d'une requête HTTP peut contenir beaucoup de choses :

- du texte ;
- un formulaire ;
- du JSON ;
- un fichier ;
- une image ;
- etc.

Il serait donc trop restrictif de demander simplement une `string`.

Go utilise ici le concept de :

```go
io.Reader
```

Un `Reader` signifie essentiellement :

> quelque chose dans lequel on peut lire des données.

`strings.NewReader(...)` permet donc de transformer une chaîne en objet satisfaisant `io.Reader`.

C'est un nouvel exemple de la puissance des **interfaces** en Go.

`httptest.NewRequest` n'a pas besoin de connaître précisément le type utilisé.

Il lui suffit que celui-ci sache être lu.

---

# 6. Le `Content-Type`

Après avoir créé la requête :

```go
request.Header.Set(
    "Content-Type",
    "application/x-www-form-urlencoded",
)
```

ou sur une ligne :

```go
request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
```

Cette ligne indique au serveur :

> Les données présentes dans le corps de cette requête sont celles d'un formulaire HTML classique encodé.

---

# 7. Pourquoi le `Content-Type` est important ?

Le corps contient seulement des octets.

Par lui-même, le serveur ne sait pas nécessairement comment les interpréter.

Il pourrait recevoir :

```text
FirstName=Robin&LastName=Des+Bois
```

Mais doit-il considérer cela comme :

- du texte ?
- du JSON ?
- un formulaire ?
- autre chose ?

Le header :

```http
Content-Type: application/x-www-form-urlencoded
```

donne cette information.

Il indique :

> Interprète le corps comme un formulaire encodé.

C'est particulièrement important ici puisque le handler utilise :

```go
r.FormValue("FirstName")
```

```go
r.FormValue("LastName")
```

etc.

Pour que Go interprète correctement le corps POST comme un formulaire, il doit connaître son type.

---

# 8. La requête simulée complète

Notre test construit donc quelque chose qui ressemble à une vraie requête HTTP :

```http
POST /members
Content-Type: application/x-www-form-urlencoded

FirstName=Robin&LastName=Des+Bois&Birthdate=1990-05-12&Email=...
```

`httptest` nous permet de fabriquer cette requête sans :

- démarrer réellement le serveur ;
- ouvrir un port ;
- utiliser un navigateur ;
- effectuer une véritable connexion réseau.

---

# 9. `httptest.NewRequest` retourne un pointeur

Lorsque l'on écrit :

```go
request := httptest.NewRequest(...)
```

le type réel de `request` est :

```go
*http.Request
```

La fonction retourne donc un **pointeur vers une `http.Request`**.

On pourrait imaginer son type ainsi :

```go
var request *http.Request
```

---

# 10. Que signifie `*http.Request` ?

Si :

```go
http.Request
```

signifie :

> une valeur de type `Request`

alors :

```go
*http.Request
```

signifie :

> un pointeur vers une valeur de type `Request`.

Schématiquement :

```text
request
   │
   │ contient une adresse
   ▼
┌──────────────────────┐
│ http.Request         │
│ Method               │
│ Header               │
│ Body                 │
│ URL                  │
│ ...                  │
└──────────────────────┘
```

La variable `request` ne contient pas directement toute la structure.

Elle permet d'accéder à la structure.

---

# 11. Pourquoi peut-on écrire `request.Header` ?

Puisque `request` est un pointeur :

```go
*http.Request
```

on pourrait imaginer devoir écrire :

```go
(*request).Header
```

Mais Go simplifie automatiquement cette écriture.

On peut donc écrire :

```go
request.Header
```

à la place de :

```go
(*request).Header
```

C'est ce que l'on fait ici :

```go
request.Header.Set(...)
```

Go comprend automatiquement :

> accéder à la structure `http.Request` pointée par `request`, puis à son champ `Header`.

---

# 12. Pourquoi une requête est-elle passée par pointeur ?

Les handlers HTTP utilisent :

```go
func(w http.ResponseWriter, r *http.Request)
```

et non :

```go
func(w http.ResponseWriter, r http.Request)
```

Il y a plusieurs intérêts.

### Éviter une copie inutile

Une `http.Request` contient beaucoup d'informations.

On évite de recopier toute la structure lors de son passage à une fonction.

### Travailler sur la même requête

Le pointeur permet de manipuler la même requête tout au long de son traitement.

On retrouve donc une idée importante :

> Le pointeur est particulièrement utile lorsque plusieurs parties du programme doivent travailler sur le même objet.

---

# 13. Le deuxième pointeur important du test : `recordingQueries`

Dans le test, on crée une structure permettant d'enregistrer les données reçues par `CreateMember`.

```go
type recordingQueries struct {
    CreateMemberParams dbsqlc.CreateMemberParams
}
```

Puis :

```go
func (q *recordingQueries) CreateMember(
    ctx context.Context,
    arg dbsqlc.CreateMemberParams,
) (dbsqlc.Member, error) {

    q.CreateMemberParams = arg

    return dbsqlc.Member{}, nil
}
```

Le point important est :

```go
(q *recordingQueries)
```

Il s'agit d'un **receiver pointeur**.

---

# 14. Pourquoi utiliser un receiver pointeur ?

La méthode doit modifier la structure :

```go
q.CreateMemberParams = arg
```

Le but du faux `Queries` est justement de retenir ce que le handler lui a envoyé.

Après l'appel du handler, le test veut pouvoir faire :

```go
queries.CreateMemberParams.FirstName
```

et retrouver :

```text
Robin
```

Il faut donc modifier **la structure originale**.

---

# 15. Ce qui se passerait avec un receiver valeur

Imaginons :

```go
func (q recordingQueries) CreateMember(...) {
    q.CreateMemberParams = arg
}
```

Ici, `q` serait une copie.

Schématiquement :

```text
recordingQueries original
        │
        │ copie
        ▼
       q
```

La méthode modifierait :

```text
q
```

mais pas la structure originale.

À la fin de la méthode, la copie disparaîtrait.

Le test ne pourrait donc pas retrouver correctement les paramètres enregistrés.

---

# 16. Avec un pointeur

Avec :

```go
func (q *recordingQueries) CreateMember(...)
```

on obtient :

```text
q
│
│ adresse
▼
recordingQueries original
```

Lorsque l'on fait :

```go
q.CreateMemberParams = arg
```

on modifie directement :

```text
recordingQueries original
```

Le test peut ensuite observer cette modification.

---

# 17. L'opérateur `&` : récupérer une adresse

Dans le test :

```go
queries := &recordingQueries{}
```

Le symbole :

```text
&
```

signifie ici :

> prendre l'adresse de cette valeur.

Sans `&` :

```go
queries := recordingQueries{}
```

le type serait :

```go
recordingQueries
```

Avec :

```go
queries := &recordingQueries{}
```

le type devient :

```go
*recordingQueries
```

On obtient donc un pointeur.

---

# 18. Relation entre `&` et `*`

Les deux symboles vont souvent ensemble.

## `&` : obtenir une adresse

```go
queries := &recordingQueries{}
```

signifie :

```text
donne-moi l'adresse de cette structure
```

Le résultat est :

```go
*recordingQueries
```

---

## `*` dans un type : pointeur vers

```go
*recordingQueries
```

signifie :

> pointeur vers un `recordingQueries`.

Même principe pour :

```go
*http.Request
```

qui signifie :

> pointeur vers une `http.Request`.

---

# 19. Petit exemple indépendant

```go
type Member struct {
    Name string
}

member := Member{
    Name: "Robin",
}

pointer := &member
```

On peut représenter cela ainsi :

```text
member
┌───────────────┐
│ Name: Robin   │
└───────────────┘
       ▲
       │
       │
    pointer
```

Le type de :

```go
member
```

est :

```go
Member
```

Le type de :

```go
pointer
```

est :

```go
*Member
```

---

# 20. Modifier la valeur via le pointeur

On peut écrire :

```go
pointer.Name = "Arthur"
```

Go comprend automatiquement :

```go
(*pointer).Name = "Arthur"
```

La structure originale est modifiée.

Après cela :

```go
fmt.Println(member.Name)
```

affiche :

```text
Arthur
```

C'est exactement l'idée utilisée avec `recordingQueries`.

---

# 21. Pourquoi `&recordingQueries{}` est particulièrement important ici ?

On écrit :

```go
queries := &recordingQueries{}
```

puis :

```go
PostMemberHandler(queries)(response, request)
```

Le handler reçoit donc un objet satisfaisant l'interface :

```go
database.Queries
```

Il appelle :

```go
queries.CreateMember(...)
```

Notre faux `Queries` exécute alors :

```go
q.CreateMemberParams = arg
```

Comme `q` est un pointeur, les données restent enregistrées dans la structure créée par le test.

On peut ensuite vérifier :

```go
queries.CreateMemberParams.FirstName
```

---

# 22. Le chemin complet

Le test peut être représenté ainsi :

```text
url.Values
    │
    │ form.Encode()
    ▼
chaîne encodée
    │
    │ strings.NewReader()
    ▼
io.Reader
    │
    │ httptest.NewRequest()
    ▼
*http.Request
    │
    │ Content-Type
    ▼
PostMemberHandler
    │
    │ r.FormValue()
    ▼
CreateMemberParams
    │
    │ queries.CreateMember()
    ▼
*recordingQueries
    │
    │ enregistrement des paramètres
    ▼
vérifications du test
```

---

# 23. Le rôle du faux `Queries`

Notre `recordingQueries` n'insère rien dans PostgreSQL.

Il sert uniquement à répondre à la question :

> Qu'est-ce que le handler aurait envoyé à la base de données ?

C'est pour cela que sa méthode :

```go
CreateMember(...)
```

enregistre simplement :

```go
q.CreateMemberParams = arg
```

puis retourne :

```go
return dbsqlc.Member{}, nil
```

Ainsi le test porte sur le **handler**, pas sur PostgreSQL.

---

# 24. Ce que vérifie réellement le test

Le test vérifie plusieurs responsabilités.

### La requête HTTP fonctionne

```go
response.Code == http.StatusOK
```

### Le handler produit la réponse attendue

```go
strings.Contains(body, "post membre")
```

### Le prénom a été transmis

```go
queries.CreateMemberParams.FirstName == firstName
```

### Le nom a été transmis

```go
queries.CreateMemberParams.LastName == lastName
```

### La date a été convertie correctement

```go
queries.CreateMemberParams.BirthDate
```

### L'email a été transmis

```go
queries.CreateMemberParams.Email
```

Le test vérifie donc principalement le chemin :

```text
HTTP
 ↓
Handler
 ↓
conversion des données
 ↓
interface Queries
```

sans avoir besoin d'une véritable base de données.

---

# Comprendre et retenir

## Pour créer une requête POST de formulaire

```go
form := url.Values{}

form.Set("FirstName", "Robin")

request := httptest.NewRequest(
    http.MethodPost,
    "/members",
    strings.NewReader(form.Encode()),
)

request.Header.Set(
    "Content-Type",
    "application/x-www-form-urlencoded",
)
```

À retenir :

```text
url.Values
    ↓ Encode()
string
    ↓ NewReader()
Reader
    ↓ NewRequest()
requête HTTP
```

---

## Pour les pointeurs

```go
T
```

signifie :

> une valeur de type `T`.

```go
*T
```

signifie :

> un pointeur vers une valeur de type `T`.

```go
&value
```

signifie :

> récupérer l'adresse de `value`.

Dans notre test :

```go
queries := &recordingQueries{}
```

donne :

```go
*recordingQueries
```

et :

```go
func (q *recordingQueries) CreateMember(...)
```

permet à la méthode de modifier la structure originale.

---

# À retenir en une phrase

Le formulaire est **encodé puis placé dans le corps d'une requête POST**, le `Content-Type` indique au handler comment interpréter ce corps, et le pointeur `*recordingQueries` permet au faux `Queries` de **conserver les données reçues afin que le test puisse les vérifier après l'appel du handler**.

