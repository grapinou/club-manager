
---


# Go — Comprendre `context.Context`

## Objectif

Comprendre à quoi sert :

```go
context.Context
```

et pourquoi nous venons de rencontrer cette notion avec :

```go
pgxpool.New(ctx, ...)
db.Ping(ctx)
```

Dans Club Manager, nous avons écrit :

```go
ctx := context.Background()
```

puis :

```go
db, err := pgxpool.New(
    ctx,
    os.Getenv("DATABASE_URL"),
)
```

et :

```go
if err := db.Ping(ctx); err != nil {
    // ...
}
```

Pourquoi faut-il transmettre ce `ctx` ?

C'est précisément ce que cette fiche cherche à expliquer.

---

# 1. Première idée : un travail possède une durée de vie

Imaginons une opération :

```text
chercher un membre dans PostgreSQL
```

Elle commence :

```text
début
  │
  ▼
requête SQL
  │
  ▼
PostgreSQL travaille
  │
  ▼
résultat
```

Normalement, le travail va jusqu'au bout.

Mais plusieurs événements peuvent rendre ce travail inutile.

Par exemple :

```text
l'utilisateur ferme la page
```

ou :

```text
la requête HTTP est annulée
```

ou encore :

```text
nous avons décidé d'attendre maximum 2 secondes
```

Dans ces situations, il serait inutile de continuer à travailler pendant plusieurs secondes.

Nous aimerions pouvoir transmettre l'information :

> Ce travail n'est plus nécessaire.

C'est l'un des rôles principaux de `Context`.

La documentation officielle présente `Context` comme un moyen de faire circuler les délais et les signaux d'annulation entre les fonctions impliquées dans une même opération.

---

# 2. Une image mentale

On peut imaginer que chaque travail possède une sorte de fiche de mission :

```text
┌──────────────────────────────┐
│ Travail                      │
│                              │
│ encore valide ?              │
│ annulé ?                     │
│ date limite ?                │
│ informations liées ?         │
└──────────────────────────────┘
```

Cette fiche accompagne le travail lorsqu'il traverse plusieurs fonctions.

```text
Handler
   │
   │ ctx
   ▼
Repository
   │
   │ ctx
   ▼
pgx
   │
   │ ctx
   ▼
PostgreSQL
```

Ainsi, chaque étape sait dans quel contexte elle travaille.

---

# 3. Le `Context` ne contient pas le travail

C'est une distinction importante.

Le `Context` n'est pas :

```text
la requête SQL
```

ni :

```text
les données du membre
```

ni :

```text
la connexion PostgreSQL
```

Il accompagne le travail.

On peut représenter cela ainsi :

```text
travail réel
    │
    ├── requête SQL
    ├── calcul
    └── résultat

context
    │
    └── informations sur la durée de vie
        de ce travail
```

---

# 4. Un contexte traverse les fonctions

Imaginons :

```go
func A(ctx context.Context) {
    B(ctx)
}
```

puis :

```go
func B(ctx context.Context) {
    C(ctx)
}
```

et :

```go
func C(ctx context.Context) {
    // opération longue
}
```

Le même contexte peut suivre le travail :

```text
A
│
│ ctx
▼
B
│
│ ctx
▼
C
```

C'est pour cette raison que les fonctions Go utilisant un contexte ont généralement cette forme :

```go
func DoSomething(
    ctx context.Context,
    ...
) error
```

La convention officielle recommande de passer le `Context` explicitement comme premier paramètre, généralement nommé `ctx`, plutôt que de le stocker dans une structure.

---

# 5. Pourquoi `context.Background()` ?

Dans Club Manager, nous avons écrit :

```go
ctx := context.Background()
```

`Background()` crée un contexte racine.

On peut le représenter ainsi :

```text
context.Background()
        │
        ▼
      ctx
```

Il ne possède initialement :

```text
aucune annulation demandée
aucune date limite
aucune valeur particulière
```

Il constitue un point de départ à partir duquel d'autres contextes pourront être créés. La documentation officielle décrit `Background` comme un contexte non-nil, vide de deadline, de valeur et de signal d'annulation, généralement utilisé au niveau principal d'un programme ou pour l'initialisation.

---

# 6. Pourquoi dans `main` ?

Notre application démarre avec :

```go
func main()
```

Nous avons donc besoin d'un premier contexte.

```text
main
 │
 ▼
context.Background()
 │
 ▼
ctx
```

Nous transmettons ensuite ce contexte à pgx :

```go
pgxpool.New(ctx, ...)
```

puis :

```go
db.Ping(ctx)
```

À ce stade du projet, notre contexte signifie essentiellement :

> Cette opération appartient au démarrage de l'application.

---

# 7. Le contexte devient surtout intéressant avec une requête HTTP

Imaginons bientôt un handler :

```go
func MemberHandler(
    w http.ResponseWriter,
    r *http.Request,
) {
    // ...
}
```

La requête HTTP possède déjà son propre contexte :

```go
r.Context()
```

On pourrait donc écrire :

```go
ctx := r.Context()
```

et transmettre ce contexte à PostgreSQL :

```go
member, err := queries.GetMember(
    ctx,
    id,
)
```

Le parcours devient :

```text
requête HTTP
     │
     ▼
 r.Context()
     │
     ▼
   handler
     │
     ▼
    sqlc
     │
     ▼
    pgx
     │
     ▼
PostgreSQL
```

Pour une requête HTTP entrante, Go annule le contexte de la requête notamment lorsque la connexion du client est fermée, lorsque la requête est annulée ou lorsque `ServeHTTP` se termine.

---

# 8. Pourquoi est-ce utile ?

Imaginons :

```text
Utilisateur
    │
    ▼
GET /members/42
```

Le serveur commence une opération PostgreSQL.

```text
Handler
   │
   ▼
PostgreSQL
   │
   └── opération longue...
```

Mais l'utilisateur ferme son navigateur.

Sans mécanisme d'annulation :

```text
utilisateur parti
      │
      ▼
serveur continue malgré tout
      │
      ▼
PostgreSQL continue
```

Avec le contexte de la requête :

```text
utilisateur parti
      │
      ▼
requête HTTP annulée
      │
      ▼
Context annulé
      │
      ▼
opération utilisant ce contexte
peut s'arrêter
```

Le contexte permet donc aux différentes couches impliquées dans une requête de recevoir le même signal d'annulation.

---

# 9. Attention : le Context n'arrête pas magiquement le code

C'est probablement l'une des idées les plus importantes.

Faire :

```go
cancel()
```

ne signifie pas :

```text
Go tue brutalement la fonction
```

Le contexte transmet un signal.

Le code qui effectue le travail doit être conçu pour observer ce signal ou appeler une API qui sait le faire.

Conceptuellement :

```text
Context
   │
   └── "annulation demandée"
              │
              ▼
       fonction concernée
              │
              ▼
      décide de s'arrêter
```

Les fonctions prenant un `Context`, comme de nombreuses opérations réseau ou base de données, peuvent utiliser ce signal pour arrêter le travail devenu inutile.

---

# 10. Créer un contexte annulable

On peut créer un nouveau contexte à partir d'un autre :

```go
ctx, cancel := context.WithCancel(parent)
```

Nous avons alors :

```text
parent
  │
  ▼
child ctx
```

et une fonction :

```go
cancel
```

capable de demander l'annulation.

Exemple :

```go
ctx, cancel := context.WithCancel(
    context.Background(),
)

defer cancel()
```

Puis quelque part :

```go
cancel()
```

Le contexte devient annulé.

```text
ctx
 │
 ▼
annulé
```

`WithCancel` crée un contexte enfant ; appeler la fonction d'annulation annule cet enfant et les contextes qui en dérivent. La documentation recommande d'appeler la fonction `cancel` afin de libérer les ressources associées.

---

# 11. Parent et enfant

Les contextes forment donc une sorte d'arbre.

```text
Background
    │
    ▼
 contexte A
    │
    ├─────────┐
    ▼         ▼
contexte B  contexte C
    │
    ▼
contexte D
```

Si :

```text
contexte A
```

est annulé, ses enfants le sont également :

```text
A annulé
│
├── B annulé
│   └── D annulé
│
└── C annulé
```

L'annulation se propage donc vers les contextes dérivés.

---

# 12. Pourquoi cette hiérarchie est intéressante ?

Prenons une requête HTTP :

```text
requête HTTP
     │
     ▼
   Context
     │
     ├── requête PostgreSQL
     │
     ├── appel API externe
     │
     └── autre traitement
```

Toutes ces opérations appartiennent à la même requête.

Si la requête principale disparaît :

```text
requête annulée
       │
       ▼
 Context annulé
       │
       ├── PostgreSQL peut s'arrêter
       ├── API externe peut s'arrêter
       └── traitement peut s'arrêter
```

On peut donc gérer la durée de vie d'un ensemble de travaux liés.

---

# 13. `WithTimeout`

Un autre cas fréquent consiste à dire :

> Cette opération peut durer au maximum deux secondes.

On peut écrire :

```go
ctx, cancel := context.WithTimeout(
    parent,
    2*time.Second,
)
defer cancel()
```

Le contexte possède alors une limite :

```text
début
  │
  ├──────────── 2 secondes ────────────┐
  │                                    │
  ▼                                    ▼
travail                          contexte annulé
```

Si le travail finit avant :

```text
travail terminé
      │
      ▼
cancel()
```

Sinon, le timeout expire et le contexte est automatiquement annulé.

`WithTimeout` crée un contexte dérivé avec une échéance relative ; il est équivalent à une deadline calculée à partir de l'heure courante. La fonction `cancel` doit tout de même être appelée pour libérer rapidement les ressources si le travail finit avant le timeout.

---

# 14. Exemple PostgreSQL

Plus tard, on pourrait écrire :

```go
ctx, cancel := context.WithTimeout(
    r.Context(),
    2*time.Second,
)
defer cancel()

member, err := queries.GetMember(
    ctx,
    memberID,
)
```

Le parcours serait :

```text
r.Context()
    │
    ▼
WithTimeout(2s)
    │
    ▼
    ctx
    │
    ▼
GetMember
    │
    ▼
pgx
    │
    ▼
PostgreSQL
```

Deux événements pourraient alors arrêter le travail :

```text
client déconnecté
```

ou :

```text
2 secondes écoulées
```

---

# 15. `WithDeadline`

Il existe également :

```go
context.WithDeadline(...)
```

La différence est surtout la manière d'exprimer la limite.

Avec :

```go
WithTimeout
```

on dit :

```text
maximum 2 secondes
```

Avec :

```go
WithDeadline
```

on dit :

```text
arrête au plus tard à telle heure
```

Exemple conceptuel :

```text
Timeout
→ durée

Deadline
→ instant précis
```

La documentation officielle définit `WithDeadline` à partir d'un instant absolu, tandis que `WithTimeout` construit cette échéance à partir d'une durée.

---

# 16. `WithValue`

Un contexte peut aussi transporter certaines informations :

```go
context.WithValue(...)
```

Par exemple, dans une application web :

```text
identifiant d'une requête
```

ou certaines informations de traçage.

Mais il faut être prudent.

Un `Context` ne doit pas devenir :

```text
un sac contenant toutes les variables
```

La documentation officielle recommande de réserver les valeurs de contexte aux données liées à la requête qui doivent traverser plusieurs API, et non de les utiliser comme paramètres optionnels ordinaires.

---

# 17. Ce qu'il ne faut pas mettre dans un Context

Par exemple :

```text
configuration de l'application
connexion PostgreSQL
prix d'une cotisation
structure Member
adresse du club
```

ne devraient normalement pas être transportés avec :

```go
context.WithValue
```

Ces données ont d'autres moyens naturels de circuler.

Le contexte doit rester associé à la vie du travail en cours.

---

# 18. Le Context et une requête HTTP

Cette image est particulièrement utile pour Club Manager.

```text
Utilisateur
    │
    ▼
Requête HTTP
    │
    ▼
r.Context()
    │
    ▼
Handler
    │
    ▼
requête SQL
    │
    ▼
pgxpool
    │
    ▼
PostgreSQL
```

Le contexte relie donc :

```text
durée de vie de la requête HTTP
```

à :

```text
durée de vie du travail PostgreSQL
```

---

# 19. `context.Background()` ou `r.Context()` ?

C'est une distinction importante.

## Au démarrage de l'application

Nous avons :

```go
ctx := context.Background()
```

puis :

```go
db.Ping(ctx)
```

Cela concerne une opération liée au démarrage du programme.

```text
main
 │
 ▼
Background
```

---

## Pendant une requête HTTP

Nous préférerons généralement utiliser :

```go
ctx := r.Context()
```

car ce contexte est lié à la requête en cours.

```text
requête utilisateur
       │
       ▼
   r.Context()
```

Ainsi :

```text
utilisateur abandonne la requête
              │
              ▼
        Context annulé
```

Pour les requêtes entrantes du serveur HTTP, Go fournit ce contexte directement via `Request.Context()`.

---

# 20. Exemple futur dans Club Manager

Imaginons :

```go
func MemberHandler(
    w http.ResponseWriter,
    r *http.Request,
) {
    ctx := r.Context()

    member, err := queries.GetMember(
        ctx,
        memberID,
    )

    // ...
}
```

Le flux devient :

```text
GET /members/42
       │
       ▼
MemberHandler
       │
       ▼
r.Context()
       │
       ▼
GetMember(ctx, 42)
       │
       ▼
pgx
       │
       ▼
PostgreSQL
```

Le contexte suit l'opération de bout en bout.

---

# 21. Pourquoi ne pas créer un nouveau Background partout ?

On pourrait être tenté d'écrire :

```go
context.Background()
```

dans chaque fonction.

Mais cela casserait la chaîne.

Exemple :

```text
requête HTTP
     │
     ▼
r.Context()
     │
     ▼
 Handler
     │
     X
context.Background()
     │
     ▼
 PostgreSQL
```

Si la requête HTTP est annulée :

```text
r.Context()
   │
   ▼
annulé
```

le nouveau :

```go
context.Background()
```

ne dépend pas de lui.

PostgreSQL pourrait donc continuer son travail.

Il est préférable de **propager le contexte existant** à travers la chaîne d'appels. C'est précisément la convention recommandée par le package `context`.

---

# 22. Le contexte comme chaîne

Une bonne manière de le retenir est :

```text
événement initial
      │
      ▼
   Context
      │
      ▼
 fonction A
      │
      ▼
 fonction B
      │
      ▼
 fonction C
```

Le contexte traverse toute la chaîne.

Si l'événement initial disparaît :

```text
annulation
    │
    ▼
 Context
    │
    ▼
toute la chaîne peut être prévenue
```

---

# 23. Une analogie

Imaginons un restaurant.

Un client commande :

```text
1 pizza
```

Le serveur transmet la commande :

```text
client
  │
  ▼
serveur
  │
  ▼
cuisine
```

La commande joue le rôle du travail.

Mais imaginons que le client parte avant que la pizza soit commencée.

Il faut aussi transmettre :

```text
commande annulée
```

Sinon :

```text
cuisine prépare la pizza
      │
      ▼
personne ne la veut
```

Le contexte ressemble à cette information qui accompagne le travail :

```text
Commande
+
état de la commande
```

Dans un serveur :

```text
requête
+
Context
```

---

# 24. Ce que Context résout

Sans contexte :

```text
A appelle B
B appelle C
C lance un travail

A ne veut plus le résultat

mais C ne le sait pas
```

Avec contexte :

```text
A
│
├── ctx
▼
B
│
├── ctx
▼
C

A annule ctx
     │
     ▼
C peut être informé
```

---

# 25. Context et goroutines

Les contextes sont particulièrement utiles lorsque plusieurs traitements liés fonctionnent en parallèle.

Par exemple :

```text
                 requête
                    │
                    ▼
                 Context
                /   |   \
               ▼    ▼    ▼
             DB    API   calcul
```

Tous peuvent partager le même contexte.

Si la requête disparaît :

```text
Context annulé
      │
      ├── DB
      ├── API
      └── calcul
```

Le même `Context` peut être utilisé simultanément par plusieurs goroutines.

---

# 26. Une interface

`context.Context` est une interface.

Elle définit notamment les méthodes permettant de consulter :

```text
deadline
annulation
erreur d'annulation
valeurs
```

Nous n'avons généralement pas besoin de créer nous-mêmes une structure implémentant cette interface.

Nous utilisons les fonctions du package :

```go
context.Background()
context.WithCancel(...)
context.WithTimeout(...)
context.WithDeadline(...)
context.WithValue(...)
```

La définition officielle de `Context` fournit notamment `Deadline`, `Done`, `Err` et `Value`.

---

# 27. `ctx.Done()`

Un contexte expose notamment :

```go
ctx.Done()
```

Cela permet à du code de savoir qu'une annulation a été demandée.

Conceptuellement :

```text
travail
  │
  ├── continue normalement
  │
  └── surveille ctx.Done()
             │
             ▼
          annulé
             │
             ▼
           stop
```

On retrouvera surtout ce mécanisme lorsque nous travaillerons avec les goroutines.

---

# 28. `ctx.Err()`

Après l'annulation, on peut également obtenir la raison générale :

```go
ctx.Err()
```

Par exemple :

```text
context canceled
```

ou :

```text
context deadline exceeded
```

Le package `context` expose ces états afin que le code puisse distinguer une annulation d'une échéance dépassée.

---

# 29. Une règle importante

Un contexte doit suivre le travail auquel il appartient.

```text
travail A
   │
   └── Context A

travail B
   │
   └── Context B
```

Il ne faut pas considérer `Context` comme une variable globale de l'application.

La documentation Go recommande d'ailleurs de ne pas stocker les contextes dans les structures mais de les passer explicitement aux fonctions qui en ont besoin.

---

# 30. Pourquoi `ctx` est souvent le premier argument ?

On rencontre souvent :

```go
func GetMember(
    ctx context.Context,
    id int,
) (...)
```

et non :

```go
func GetMember(
    id int,
    ctx context.Context,
) (...)
```

La convention Go consiste à placer le contexte en premier :

```text
fonction(
    contexte du travail,
    données du travail,
)
```

La documentation officielle recommande explicitement cette forme.

---

# 31. Notre utilisation actuelle

Dans Club Manager :

```go
ctx := context.Background()
```

puis :

```go
db, err := pgxpool.New(
    ctx,
    os.Getenv("DATABASE_URL"),
)
```

et :

```go
db.Ping(ctx)
```

Le parcours est :

```text
main
 │
 ▼
context.Background()
 │
 ▼
ctx
 │
 ├────► pgxpool.New
 │
 └────► db.Ping
             │
             ▼
         PostgreSQL
```

Nous avons donc rencontré `Context` pour la première fois parce que pgx permet à ses opérations de participer à cette mécanique de durée de vie et d'annulation.

---

# 32. Plus tard avec sqlc

Lorsque sqlc générera nos fonctions, nous rencontrerons probablement des signatures ressemblant à :

```go
func (q *Queries) GetMember(
    ctx context.Context,
    id int32,
) (...)
```

Nous pourrons alors transmettre directement :

```go
r.Context()
```

depuis le handler.

```text
HTTP
 │
 ▼
r.Context()
 │
 ▼
handler
 │
 ▼
sqlc
 │
 ▼
pgx
 │
 ▼
PostgreSQL
```

Le contexte deviendra ainsi beaucoup plus concret.

---

# Comprendre et retenir

> **Un Context accompagne un travail.**

Il contient surtout des informations concernant la durée de vie de ce travail.

---

> **Le Context permet de transmettre un signal d'annulation.**

```text
travail devenu inutile
        │
        ▼
Context annulé
        │
        ▼
les opérations concernées
peuvent s'arrêter
```

---

> **Le contexte traverse les fonctions.**

```text
Handler
  │ ctx
  ▼
Service
  │ ctx
  ▼
Database
```

On ne recrée pas un `Background()` à chaque étage.

---

> **`context.Background()` est un point de départ.**

Dans notre cas :

```text
main
 ↓
Background
 ↓
connexion PostgreSQL
```

---

> **Pour une requête HTTP, `r.Context()` est lié à cette requête.**

```text
HTTP request
     │
     ▼
r.Context()
```

Si la requête disparaît, son contexte peut signaler cette disparition aux opérations qui en dépendent.

---

> **`WithTimeout` ajoute une durée maximale.**

```text
Context parent
      │
      ▼
WithTimeout
      │
      ▼
Context enfant
avec limite de temps
```

---

> **Le contexte n'arrête pas magiquement une fonction.**

Il signale :

```text
"ce travail doit s'arrêter"
```

Les opérations doivent être conçues pour respecter ce signal.

---

> **Le Context n'est pas un sac à variables.**

Il ne remplace pas :

```text
Config
struct
arguments de fonction
base de données
état global
```

Les valeurs de contexte sont réservées à des données liées à la requête qui doivent traverser plusieurs couches.

---

# La phrase à retenir

> **`context.Context` permet de faire voyager la durée de vie d'un travail à travers le programme.**

Dans Club Manager :

```text
requête utilisateur
        │
        ▼
     Context
        │
        ▼
     Handler
        │
        ▼
      sqlc
        │
        ▼
       pgx
        │
        ▼
   PostgreSQL
```

Le contexte permet à toutes ces couches de comprendre :

```text
le travail est-il encore utile ?
```

