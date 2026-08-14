
---

## Objectif

`curl` est un outil en ligne de commande permettant d'envoyer des requêtes vers un serveur.

Il est particulièrement utile pour tester un backend sans avoir besoin :

- d'un navigateur ;
- d'un formulaire HTML ;
- d'une interface graphique ;
- d'un outil externe comme Postman.

Dans Club Manager, nous l'avons utilisé pour tester directement la route :

```text
POST /members
```

et vérifier que toute la chaîne fonctionnait jusqu'à PostgreSQL.

---

# 1. Principe général

On peut voir `curl` comme un petit client HTTP utilisable depuis le terminal.

Par exemple :

```bash
curl http://localhost:8080/
```

envoie une requête vers :

```text
http://localhost:8080/
```

Par défaut, `curl` utilise une requête HTTP `GET`.

Conceptuellement :

```text
curl
  ↓
requête HTTP
  ↓
serveur
  ↓
réponse HTTP
  ↓
terminal
```

---

# 2. Le POST utilisé dans Club Manager

Nous avons utilisé :

```bash
curl -i -X POST http://localhost:8080/members \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "FirstName=Robin" \
  -d "LastName=Des Bois" \
  -d "Birthdate=1990-05-12" \
  -d "Email=robin.desbois@example.com"
```

Cette commande reproduit ce que ferait un formulaire HTML envoyant les données d'un nouveau membre.

---

# 3. Décomposition de la commande

## `curl`

```bash
curl
```

lance le programme.

Il va créer et envoyer une requête HTTP.

---

## `-i`

```bash
-i
```

demande à `curl` d'afficher également les en-têtes de la réponse HTTP.

Sans `-i`, on aurait principalement vu :

```text
post membre
```

Avec `-i`, nous avons obtenu :

```text
HTTP/1.1 200 OK
Date: Fri, 14 Aug 2026 08:22:20 GMT
Content-Length: 12
Content-Type: text/plain; charset=utf-8

post membre
```

Cela permet d'observer non seulement le contenu de la réponse, mais également les informations HTTP associées.

---

# 4. `-X POST`

```bash
-X POST
```

indique la méthode HTTP utilisée.

Ici :

```text
POST
```

La requête envoyée est donc :

```text
POST /members
```

Dans notre test Go, cela correspondait à :

```go
http.MethodPost
```

---

# 5. L'adresse

```bash
http://localhost:8080/members
```

indique l'adresse à laquelle envoyer la requête.

Décomposition :

```text
http://
```

protocole utilisé.

```text
localhost
```

machine locale.

```text
8080
```

port utilisé par Club Manager.

```text
/members
```

route appelée.

Dans le routeur Go, nous avons :

```go
mux.HandleFunc(
    "POST /members",
    handlers.PostMemberHandler(queries),
)
```

---

# 6. `-H` : ajouter un header HTTP

Nous avons utilisé :

```bash
-H "Content-Type: application/x-www-form-urlencoded"
```

`-H` signifie que l'on ajoute un **header HTTP**.

Ici :

```text
Content-Type
```

indique la nature du contenu envoyé dans le corps de la requête.

La valeur :

```text
application/x-www-form-urlencoded
```

signifie :

> les données sont encodées comme celles d'un formulaire HTML classique.

---

# 7. L'équivalent dans notre test Go

Dans le test du handler, nous avions :

```go
request.Header.Set(
    "Content-Type",
    "application/x-www-form-urlencoded",
)
```

La commande `curl` :

```bash
-H "Content-Type: application/x-www-form-urlencoded"
```

fait donc exactement la même chose.

Correspondance :

```text
Go                                              curl

request.Header.Set(...)                         -H "Content-Type: ..."
```

---

# 8. `-d` : envoyer des données

Nous avons utilisé plusieurs fois :

```bash
-d
```

Par exemple :

```bash
-d "FirstName=Robin"
```

Cela signifie :

> ajouter cette donnée dans le corps de la requête.

Nous avons envoyé :

```bash
-d "FirstName=Robin"
-d "LastName=Des Bois"
-d "Birthdate=1990-05-12"
-d "Email=robin.desbois@example.com"
```

Ces données représentent notre formulaire.

---

# 9. Le corps de la requête

Les différentes valeurs envoyées avec `-d` sont transformées en un corps correspondant à un formulaire.

Conceptuellement :

```text
FirstName=Robin
LastName=Des Bois
Birthdate=1990-05-12
Email=robin.desbois@example.com
```

devient quelque chose de proche de :

```text
FirstName=Robin&LastName=Des+Bois&Birthdate=1990-05-12&Email=robin.desbois%40example.com
```

C'est le même principe que :

```go
form.Encode()
```

dans notre test Go.

---

# 10. Correspondance avec `url.Values`

Dans notre test, nous avions :

```go
form := url.Values{}

form.Set("FirstName", firstName)
form.Set("LastName", lastName)
form.Set("Birthdate", birthdate)
form.Set("Email", email)
```

Puis :

```go
form.Encode()
```

Avec `curl`, on écrit directement :

```bash
-d "FirstName=Robin"
-d "LastName=Des Bois"
-d "Birthdate=1990-05-12"
-d "Email=robin.desbois@example.com"
```

On peut donc faire la correspondance suivante :

```text
Go                                      curl

form.Set("FirstName", "Robin")          -d "FirstName=Robin"

form.Set("LastName", "Des Bois")        -d "LastName=Des Bois"

form.Encode()                           données assemblées par curl
```

---

# 11. Correspondance complète entre le test Go et curl

Dans notre test :

```go
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

Avec `curl` :

```bash
curl -i -X POST http://localhost:8080/members \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "FirstName=Robin" \
  -d "LastName=Des Bois" \
  -d "Birthdate=1990-05-12" \
  -d "Email=robin.desbois@example.com"
```

Correspondance :

```text
httptest.NewRequest(...)            curl

http.MethodPost                     -X POST

"/members"                          http://localhost:8080/members

form.Set(...)                       -d "..."

form.Encode()                       données du formulaire

Content-Type                        -H "Content-Type: ..."
```

---

# 12. Différence entre `httptest` et `curl`

Avec `httptest`, nous testions directement le handler en Go.

```text
httptest
    ↓
PostMemberHandler
    ↓
recordingQueries
```

Nous utilisions volontairement un faux objet :

```go
recordingQueries
```

La base PostgreSQL n'était donc pas utilisée.

Le test permettait de vérifier :

> Est-ce que le handler transmet les bonnes données à la couche base de données ?

---

Avec `curl`, nous avons utilisé le vrai serveur.

```text
curl
  ↓
HTTP
  ↓
router
  ↓
PostMemberHandler
  ↓
sqlc
  ↓
pgx
  ↓
PostgreSQL
```

Cette fois, toute l'application était réellement utilisée.

---

# 13. Le résultat obtenu

Après avoir exécuté :

```bash
curl -i -X POST http://localhost:8080/members \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "FirstName=Robin" \
  -d "LastName=Des Bois" \
  -d "Birthdate=1990-05-12" \
  -d "Email=robin.desbois@example.com"
```

le serveur a répondu :

```text
HTTP/1.1 200 OK
Date: Fri, 14 Aug 2026 08:22:20 GMT
Content-Length: 12
Content-Type: text/plain; charset=utf-8

post membre
```

---

# 14. Que signifie `200 OK` ?

```text
HTTP/1.1 200 OK
```

signifie que la requête a été traitée avec succès.

Dans notre handler, l'insertion se fait avant :

```go
fmt.Fprintln(w, "post membre")
```

Si `CreateMember` échoue :

```go
_, err = queries.CreateMember(r.Context(), datas)

if err != nil {
    http.Error(
        w,
        "db insert non fait",
        http.StatusInternalServerError,
    )
    return
}
```

on obtient une erreur HTTP.

Le fait d'obtenir :

```text
200 OK
```

puis :

```text
post membre
```

montre donc que le handler est arrivé jusqu'au bout.

---

# 15. `Content-Type` de la réponse

Nous avons reçu :

```text
Content-Type: text/plain; charset=utf-8
```

Cette fois, il ne s'agit pas du `Content-Type` envoyé avec la requête.

Il s'agit du type de contenu **renvoyé par le serveur**.

La réponse étant :

```text
post membre
```

Go indique qu'il s'agit de texte :

```text
text/plain
```

avec un encodage :

```text
utf-8
```

Il faut donc distinguer :

```text
REQUÊTE
Content-Type: application/x-www-form-urlencoded
```

et :

```text
RÉPONSE
Content-Type: text/plain; charset=utf-8
```

---

# 16. Vérification dans PostgreSQL

Après le `curl`, nous nous sommes connectés à PostgreSQL :

```bash
psql -h localhost -U club_manager -d club_manager
```

Puis :

```sql
SELECT * FROM members;
```

Le résultat obtenu était :

```text
 id | first_name | last_name | birth_date |           email           |          created_at
----+------------+-----------+------------+---------------------------+-------------------------------
  1 | Robin      | Des Bois  | 1990-05-12 | robin.desbois@example.com | 2026-08-14 10:22:20.847458+02
```

Cela valide définitivement que les données ont été enregistrées.

---

# 17. Ce que nous avons réellement testé

Le test avec `curl` a validé :

```text
curl
  ↓
requête HTTP POST
  ↓
route POST /members
  ↓
PostMemberHandler
  ↓
lecture du formulaire avec r.FormValue()
  ↓
conversion des données
  ↓
CreateMemberParams
  ↓
queries.CreateMember()
  ↓
sqlc
  ↓
pgx
  ↓
PostgreSQL
  ↓
ligne dans members
```

Contrairement au test unitaire du handler, il s'agit donc d'un test manuel beaucoup plus proche du fonctionnement réel de l'application.

---

# 18. Pourquoi curl est utile pour développer un backend ?

Lorsqu'on développe un backend, l'interface graphique n'est pas nécessaire pour vérifier que le serveur fonctionne.

On peut développer :

```text
route
 ↓
handler
 ↓
base de données
```

et tester directement cette partie avec `curl`.

L'interface HTML peut arriver ensuite.

Cela permet de séparer les problèmes.

Par exemple :

```text
curl fonctionne
mais formulaire HTML ne fonctionne pas
```

signifie probablement :

> le backend fonctionne, il faut chercher le problème du côté du formulaire.

À l'inverse :

```text
curl ne fonctionne pas
```

indique que le problème se trouve probablement déjà dans le backend.

---

# 19. Commandes curl utiles

## Faire un GET simple

```bash
curl http://localhost:8080/
```

---

## Voir les headers et le contenu

```bash
curl -i http://localhost:8080/
```

---

## Faire un POST

```bash
curl -X POST http://localhost:8080/members
```

---

## Ajouter un header

```bash
curl \
  -H "Content-Type: application/x-www-form-urlencoded" \
  http://localhost:8080/members
```

---

## Envoyer une donnée

```bash
curl \
  -d "FirstName=Robin" \
  http://localhost:8080/members
```

---

## Envoyer plusieurs valeurs

```bash
curl \
  -d "FirstName=Robin" \
  -d "LastName=Des Bois" \
  http://localhost:8080/members
```

---

# 20. Une remarque sur `-X POST`

Dans notre commande, nous avons utilisé :

```bash
-X POST
```

C'est très explicite et donc intéressant pour apprendre.

Cependant, lorsque `curl` reçoit des données avec :

```bash
-d
```

il utilise normalement automatiquement une requête `POST`.

Cette commande :

```bash
curl \
  -d "FirstName=Robin" \
  http://localhost:8080/members
```

enverra donc déjà un `POST`.

Pour apprendre et lire facilement la commande, conserver :

```bash
-X POST
```

reste néanmoins très clair.

---

# 21. `curl` comme outil de diagnostic

Une commande `curl` permet de tester indépendamment chaque route.

Par exemple :

```bash
curl -i http://localhost:8080/
```

permet de tester l'accueil.

Puis :

```bash
curl -i http://localhost:8080/contact
```

la page contact.

Et :

```bash
curl -i -X POST http://localhost:8080/members ...
```

la création d'un membre.

On peut donc considérer `curl` comme un outil simple pour dialoguer directement avec le serveur HTTP.

---

# Comprendre et retenir

`curl` est un **client HTTP en ligne de commande**.

La commande :

```bash
curl -i -X POST http://localhost:8080/members \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "FirstName=Robin"
```

peut être lue ainsi :

```text
curl
│
├── -i
│      afficher les headers de la réponse
│
├── -X POST
│      utiliser la méthode POST
│
├── URL
│      destination de la requête
│
├── -H
│      ajouter un header HTTP
│
└── -d
       envoyer une donnée dans le corps
```

Le point essentiel est que `curl` permet de tester directement :

```text
HTTP
 ↓
backend
 ↓
base de données
```

sans avoir besoin d'une interface HTML.

---

# À retenir en une phrase

`curl` permet d'envoyer directement des requêtes HTTP depuis le terminal ; dans Club Manager, il nous a permis de vérifier que `POST /members` fonctionnait réellement depuis la réception du formulaire jusqu'à l'insertion du membre dans PostgreSQL.

