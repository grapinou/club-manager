
---


Cette fiche complète la première fiche sur les interfaces.

L'objectif n'est plus seulement de comprendre :

> « une struct peut satisfaire une interface »

mais de voir **pourquoi une interface devient utile dans une vraie architecture**.

---

# 1. Une interface décrit un contrat

Une interface Go décrit un ensemble de méthodes attendues.

Exemple :

```go
type Queries interface {
	CreateMember(
		ctx context.Context,
		arg dbsqlc.CreateMemberParams,
	) (dbsqlc.Member, error)
}
```

Cette interface dit :

> Tout type possédant cette méthode avec exactement cette signature peut être utilisé comme un `Queries`.

---

# 2. La signature complète compte

Pour satisfaire une interface, il ne suffit pas d'avoir une méthode portant le même nom.

Il faut que correspondent :

- le nom de la méthode ;
    
- les paramètres ;
    
- leur ordre ;
    
- leurs types ;
    
- les valeurs retournées ;
    
- leur ordre ;
    
- leurs types.
    

Ainsi :

```go
func (q *MyQueries) CreateMember(
	ctx context.Context,
	arg dbsqlc.CreateMemberParams,
) (dbsqlc.Member, error)
```

satisfait :

```go
type Queries interface {
	CreateMember(
		ctx context.Context,
		arg dbsqlc.CreateMemberParams,
	) (dbsqlc.Member, error)
}
```

En revanche :

```go
func (q *MyQueries) CreateMember(name string) error
```

ne satisfait pas cette interface.

Même si le nom est identique :

```text
CreateMember
```

la signature est différente.

---

# 3. Une interface définit un type par son comportement

Habituellement, avec une struct :

```go
type Dog struct {
	Name string
}
```

le type dépend de sa structure.

Avec une interface, le type dépend de son comportement.

```go
type Speaker interface {
	Speak()
}
```

Tout type possédant :

```go
Speak()
```

satisfait `Speaker`.

Exemple :

```go
type Dog struct{}

func (Dog) Speak() {
	fmt.Println("Wouf")
}
```

et :

```go
type Human struct{}

func (Human) Speak() {
	fmt.Println("Bonjour")
}
```

Les structures sont différentes :

```text
Dog
Human
```

mais les deux satisfont :

```text
Speaker
```

car les deux possèdent :

```go
Speak()
```

---

# 4. Une interface n'est pas implémentée explicitement

Dans d'autres langages, on pourrait rencontrer quelque chose comme :

```text
implements Speaker
```

Go ne fonctionne pas ainsi.

On écrit seulement :

```go
func (Dog) Speak() {}
```

À partir de ce moment :

```text
Dog possède Speak()
        │
        ▼
Dog satisfait Speaker
```

L'implémentation est implicite.

---

# 5. Application à Club Manager

sqlc génère notamment le type :

```go
*dbsqlc.Queries
```

et celui-ci possède une méthode :

```go
CreateMember(
	ctx context.Context,
	arg dbsqlc.CreateMemberParams,
) (dbsqlc.Member, error)
```

Nous pouvons définir notre propre interface :

```go
type Queries interface {
	CreateMember(
		ctx context.Context,
		arg dbsqlc.CreateMemberParams,
	) (dbsqlc.Member, error)
}
```

Comme `*dbsqlc.Queries` possède cette méthode avec exactement cette signature :

```text
*dbsqlc.Queries
      │
      │ possède CreateMember(...)
      ▼
   Queries
```

Il satisfait automatiquement l'interface.

---

# 6. Le type concret et le type attendu peuvent être différents

Dans `main` :

```go
queries := dbsqlc.New(db)
```

Le type réel de `queries` est :

```go
*dbsqlc.Queries
```

Mais le routeur peut demander :

```go
func New(
	cfg config.Config,
	queries Queries,
) *http.ServeMux
```

Le routeur ne demande donc plus précisément :

```text
*dbsqlc.Queries
```

Il demande :

```text
quelque chose qui satisfait Queries
```

---

# 7. Avant l'interface

Sans interface :

```go
func New(
	cfg config.Config,
	queries *dbsqlc.Queries,
) *http.ServeMux
```

Le routeur exige exactement :

```text
*dbsqlc.Queries
```

On ne peut donc lui fournir que ce type ou quelque chose compatible avec ce type précis.

Cela crée un lien fort :

```text
router
   │
   ▼
dbsqlc.Queries
   │
   ▼
PostgreSQL
```

---

# 8. Après l'interface

Avec :

```go
type Queries interface {
	CreateMember(
		ctx context.Context,
		arg dbsqlc.CreateMemberParams,
	) (dbsqlc.Member, error)
}
```

le routeur peut recevoir :

```go
func New(
	cfg config.Config,
	queries Queries,
) *http.ServeMux
```

Il ne dépend plus directement de l'implémentation concrète.

```text
                 Queries
                /       \
               /         \
              ▼           ▼
   *dbsqlc.Queries     FakeQueries
     production          tests
```

---

# 9. Le duo struct + méthodes

On peut créer un autre type :

```go
type FakeQueries struct{}
```

Puis lui donner la méthode demandée :

```go
func (FakeQueries) CreateMember(
	ctx context.Context,
	arg dbsqlc.CreateMemberParams,
) (dbsqlc.Member, error) {

	return dbsqlc.Member{}, nil
}
```

À partir de ce moment :

```text
FakeQueries
     │
     │ possède CreateMember(...)
     ▼
   Queries
```

`FakeQueries` satisfait donc également l'interface.

---

# 10. Pourquoi est-ce utile pour les tests ?

En production :

```go
queries := dbsqlc.New(db)

mux := router.New(cfg, queries)
```

On utilise :

```text
*dbsqlc.Queries
```

qui communique réellement avec PostgreSQL.

Pendant les tests :

```go
queries := FakeQueries{}

mux := New(cfg, queries)
```

On utilise :

```text
FakeQueries
```

sans PostgreSQL.

Les deux fonctionnent car le routeur ne demande pas un type concret.

Il demande seulement :

```text
Queries
```

---

# 11. Même interface, comportements différents

Les deux types peuvent posséder exactement la même méthode :

```text
CreateMember(...)
```

mais avoir des comportements totalement différents.

## Production

```go
func (*dbsqlc.Queries) CreateMember(...) (...) {
	// requête SQL réelle
}
```

Conceptuellement :

```text
CreateMember
    │
    ▼
PostgreSQL
```

## Test

```go
func (FakeQueries) CreateMember(...) (...) {
	return dbsqlc.Member{}, nil
}
```

Conceptuellement :

```text
CreateMember
    │
    ▼
valeur simulée
```

L'interface ne définit pas **comment** la méthode fonctionne.

Elle définit seulement :

> cette méthode doit exister avec cette signature.

---

# 12. Le contrat et l'implémentation

On peut distinguer deux choses.

## Le contrat

```go
type Queries interface {
	CreateMember(...)
}
```

Il décrit :

```text
ce que l'on peut faire
```

## L'implémentation

```go
*dbsqlc.Queries
```

ou :

```go
FakeQueries
```

décrit :

```text
comment cela est fait
```

Cela donne :

```text
Interface
   │
   │ définit
   ▼
contrat
   │
   ├───────────────┐
   ▼               ▼
production       tests
```

---

# 13. Une interface n'a pas besoin de connaître toute l'implémentation

`*dbsqlc.Queries` pourra posséder beaucoup de méthodes :

```text
CreateMember
GetMember
ListMembers
UpdateMember
DeleteMember
CreatePayment
ListPayments
...
```

Notre interface n'est pas obligée de toutes les reprendre.

Elle peut commencer par :

```go
type Queries interface {
	CreateMember(...)
}
```

Si plus tard Club Manager utilise :

```go
GetMember(...)
```

on pourra alors ajouter :

```go
type Queries interface {
	CreateMember(...)
	GetMember(...)
}
```

---

# 14. Ne pas créer tout le CRUD à l'avance

On pourrait être tenté d'écrire immédiatement :

```go
type Queries interface {
	CreateMember(...)
	GetMember(...)
	ListMembers(...)
	UpdateMember(...)
	DeleteMember(...)
}
```

Mais si seule `CreateMember` existe réellement aujourd'hui, cela ajoute des besoins artificiels.

Une règle utile est :

> Une interface doit décrire les comportements réellement nécessaires à son utilisateur.

On peut donc la faire évoluer progressivement.

Aujourd'hui :

```text
Queries
└── CreateMember
```

Plus tard :

```text
Queries
├── CreateMember
├── GetMember
├── ListMembers
├── UpdateMember
└── DeleteMember
```

si ces opérations deviennent réellement nécessaires.

---

# 15. L'interface permet de remplacer l'implémentation

C'est une des idées les plus importantes.

Sans interface :

```text
router
   │
   ▼
*dbsqlc.Queries
```

Avec interface :

```text
router
   │
   ▼
Queries
   │
   ├── *dbsqlc.Queries
   └── FakeQueries
```

Le routeur n'a plus besoin de connaître le détail de l'objet reçu.

Il connaît seulement son contrat.

---

# 16. Exemple complet simplifié

Interface :

```go
type Queries interface {
	CreateMember(
		ctx context.Context,
		arg dbsqlc.CreateMemberParams,
	) (dbsqlc.Member, error)
}
```

Production :

```go
queries := dbsqlc.New(db)

mux := router.New(cfg, queries)
```

Test :

```go
type FakeQueries struct{}

func (FakeQueries) CreateMember(
	ctx context.Context,
	arg dbsqlc.CreateMemberParams,
) (dbsqlc.Member, error) {

	return dbsqlc.Member{}, nil
}
```

Puis :

```go
queries := FakeQueries{}

mux := New(cfg, queries)
```

Même fonction :

```go
router.New(...)
```

Deux implémentations différentes.

---

# 17. Comment raisonner face à une interface ?

Lorsqu'une fonction attend :

```go
func Use(q Queries)
```

il ne faut pas se demander :

> Est-ce que mon objet est un `Queries` ?

Il faut plutôt se demander :

> Est-ce que mon type possède toutes les méthodes demandées par `Queries`, avec exactement les bonnes signatures ?

Par exemple :

```text
Queries demande :

CreateMember(context.Context, CreateMemberParams)
    → (Member, error)
```

On vérifie alors :

```text
*dbsqlc.Queries
    └── CreateMember(...) ✅

FakeQueries
    └── CreateMember(...) ✅
```

Les deux sont acceptés.

---

# 18. Une interface est un contrat comportemental

On peut résumer ainsi :

```text
struct
→ décrit principalement des données

méthodes
→ décrivent ce qu'un type sait faire

interface
→ décrit ce qu'un utilisateur attend qu'un type sache faire
```

Ainsi :

```text
struct + méthodes
       │
       ▼
comportement disponible
       │
       ▼
correspond au contrat ?
       │
      oui
       ▼
interface satisfaite
```

---

# 19. Application directe à Club Manager

Notre architecture peut devenir :

```text
main
 │
 ├── database.New()
 │       │
 │       ▼
 │     pgxpool
 │
 ├── dbsqlc.New(db)
 │       │
 │       ▼
 │  *dbsqlc.Queries
 │       │
 │       ▼
 │     Queries
 │
 ▼
router
 │
 ▼
handlers
 │
 ▼
CreateMember(...)
```

Pour les tests :

```text
test
 │
 ▼
FakeQueries
 │
 ▼
Queries
 │
 ▼
router
 │
 ▼
handlers
```

PostgreSQL n'est plus nécessaire pour tester le comportement HTTP qui ne dépend pas réellement de la base.

---

# Comprendre et retenir

> **Une interface définit des méthodes attendues, pas des données.**

```go
type Queries interface {
	CreateMember(...)
}
```

---

> **La signature complète de la méthode doit correspondre.**

Le nom seul ne suffit pas.

```text
nom
+ paramètres
+ types
+ retours
= signature compatible
```

---

> **Un type satisfait implicitement une interface.**

```text
*dbsqlc.Queries
      │
      │ possède les bonnes méthodes
      ▼
    Queries
```

Aucun `implements` n'est nécessaire.

---

> **Plusieurs types peuvent satisfaire la même interface.**

```text
              Queries
             /       \
            ▼         ▼
*dbsqlc.Queries   FakeQueries
  production        tests
```

---

> **L'interface sépare le contrat de l'implémentation.**

```text
interface
→ ce que l'on attend

struct + méthodes
→ comment c'est réalisé
```

---

> **Une interface ne doit pas nécessairement contenir tout le CRUD dès le départ.**

On commence par les besoins réels :

```text
Queries
└── CreateMember
```

et on la fait évoluer lorsque de nouveaux besoins apparaissent.

---

# À retenir pour Club Manager

```text
queries := dbsqlc.New(db)
```

retourne :

```text
*dbsqlc.Queries
```

Ce type possède :

```text
CreateMember(...)
```

Si notre interface `Queries` demande exactement cette méthode :

```text
*dbsqlc.Queries satisfait Queries
```

Le routeur peut alors demander :

```go
func New(cfg config.Config, queries Queries) *http.ServeMux
```

et recevoir aussi bien :

```text
*dbsqlc.Queries
```

que :

```text
FakeQueries
```

C'est ce qui nous permet de conserver une vraie base en production tout en utilisant une implémentation légère dans les tests.