
---

## Objectif

Cette fiche constitue un **jalon figé** de l'architecture de Club Manager.

Elle fait suite à :

- [[00 - Architecture minimale en Go]]
    
- [[01 - SEPARATION DES RESPONSABILITES]]
    

L'objectif n'est pas de maintenir cette fiche à jour, mais de conserver une photographie du projet à ce stade de son développement.

Ce jalon montre notamment :

- l'introduction d'une configuration externe en JSON ;
    
- la transmission explicite des dépendances ;
    
- l'introduction des vues et templates HTML ;
    
- la création d'un layout commun ;
    
- la factorisation des pages ayant la même structure ;
    
- la différence entre une page générique et une page spécifique.
    

---

# 1. Architecture actuelle

À ce stade, le chemin principal d'une requête est :

```text
config/config.json
        │
        ▼
   config.Load()
        │
        ▼
   config.Config
        │
        ▼
      main
        │
        ▼
  router.New(cfg)
        │
        ▼
     handlers
        │
        ▼
       views
        │
        ▼
    templates HTML
        │
        ▼
      navigateur
```

L'application possède maintenant plusieurs responsabilités clairement séparées.

---

# 2. Le rôle de `main`

Le fichier :

```text
cmd/server/main.go
```

reste le point d'entrée de l'application.

Son rôle est maintenant principalement de :

1. charger la configuration ;
    
2. arrêter l'application si la configuration ne peut pas être chargée ;
    
3. transmettre la configuration au routeur ;
    
4. démarrer le serveur HTTP.
    

Exemple :

```go
cfg, err := config.Load(configPath)

mux := router.New(cfg)

http.ListenAndServe(":8080", mux)
```

`main` ne connaît pas le détail des pages.

Il assemble les différentes parties de l'application.

---

# 3. La configuration

La configuration est stockée dans :

```text
config/config.json
```

Elle est chargée par le package :

```text
internal/config
```

La structure principale est :

```go
type Config struct {
	SiteName string
	Home     PageConfig
	Club     PageConfig
	Contact  ContactConfig
	Rules    PageConfig
	Where    PageConfig
	When     PageConfig
}
```

Les pages simples utilisent une structure commune :

```go
type PageConfig struct {
	Title       string
	Heading     string
	Description string
}
```

Cette structure est utilisée par :

```text
Home
Club
Rules
Where
When
```

La page `Contact` possède des besoins supplémentaires.

Elle conserve donc sa propre structure :

```go
type ContactConfig struct {
	Title        string
	Heading      string
	Description  string
	EmailAddress string
	PhoneNumber  string
}
```

---

# 4. Une première factorisation : `PageConfig`

Au départ, chaque page possédait sa propre structure :

```text
HomeConfig
ClubConfig
RulesConfig
WhereConfig
WhenConfig
```

Elles contenaient pourtant exactement les mêmes champs.

Cette répétition a conduit à la création de :

```go
PageConfig
```

On ne factorise donc pas par anticipation.

On factorise parce qu'une **répétition réelle est apparue dans le code**.

```text
HomeConfig  ─┐
ClubConfig  ─┤
RulesConfig ─┼──► PageConfig
WhereConfig ─┤
WhenConfig  ─┘
```

---

# 5. Le routeur

Le package :

```text
internal/router
```

associe les routes HTTP aux handlers.

Exemple :

```go
mux.HandleFunc("GET /{$}", handlers.HomeHandler(cfg))
mux.HandleFunc("GET /club", handlers.ClubHandler(cfg))
mux.HandleFunc("GET /contact", handlers.ContactHandler(cfg))
mux.HandleFunc("GET /where", handlers.WhereHandler(cfg))
mux.HandleFunc("GET /when", handlers.WhenHandler(cfg))
mux.HandleFunc("GET /rules", handlers.RulesHandler(cfg))
```

Le routeur reste très lisible.

Il indique directement :

```text
route → handler
```

Il ne s'occupe :

- ni de construire les données ;
    
- ni de générer le HTML ;
    
- ni de lire le fichier JSON.
    

---

# 6. Injection de la configuration

La configuration est chargée une seule fois dans `main`.

Elle est ensuite transmise :

```text
main
 │
 │ cfg
 ▼
router.New(cfg)
 │
 │ cfg
 ▼
HomeHandler(cfg)
```

Un handler ne va donc pas chercher lui-même une configuration globale.

Sa dépendance lui est fournie explicitement.

Exemple :

```go
func HomeHandler(cfg config.Config) http.HandlerFunc {
	return pageHandler(cfg.SiteName, cfg.Home)
}
```

Cette approche rend les dépendances visibles et facilite également les tests.

---

# 7. Une deuxième factorisation : `PageData`

Les pages simples ont également besoin des mêmes données pour leur affichage :

```text
SiteName
Title
Heading
Description
```

Une structure commune a donc été créée :

```go
type PageData struct {
	SiteName    string
	Title       string
	Heading     string
	Description string
}
```

On distingue ainsi deux représentations différentes :

```text
PageConfig
    │
    │ données venant de la configuration
    ▼
PageData
    │
    │ données destinées à l'affichage
    ▼
template HTML
```

Même si elles se ressemblent actuellement, elles n'ont pas la même responsabilité.

Par exemple :

```go
Title: page.Title + " - " + siteName
```

montre qu'une donnée de configuration peut être transformée avant d'être envoyée à la vue.

---

# 8. Une troisième factorisation : `pageHandler`

Les handlers des pages simples réalisaient tous le même travail :

```text
récupérer PageConfig
        ↓
construire PageData
        ↓
appeler RenderPage
        ↓
gérer une éventuelle erreur
```

Cette mécanique a été regroupée dans :

```go
func pageHandler(
	siteName string,
	page config.PageConfig,
) http.HandlerFunc
```

`pageHandler` :

1. reçoit uniquement les données dont il a besoin ;
    
2. construit `views.PageData` ;
    
3. appelle `views.RenderPage` ;
    
4. gère l'erreur de rendu.
    

---

# 9. Les handlers publics deviennent très simples

Les handlers conservent leur nom métier.

Par exemple :

```go
func HomeHandler(cfg config.Config) http.HandlerFunc {
	return pageHandler(cfg.SiteName, cfg.Home)
}
```

et :

```go
func ClubHandler(cfg config.Config) http.HandlerFunc {
	return pageHandler(cfg.SiteName, cfg.Club)
}
```

Ils indiquent simplement **quelle configuration de page utiliser**.

On obtient :

```text
HomeHandler ─────► cfg.Home  ───┐
ClubHandler ─────► cfg.Club  ───┤
RulesHandler ────► cfg.Rules ───┤
WhereHandler ────► cfg.Where ───┼──► pageHandler
WhenHandler ─────► cfg.When  ───┘
```

Cette solution permet de conserver un routeur explicite :

```go
handlers.HomeHandler(cfg)
handlers.ClubHandler(cfg)
```

tout en évitant de répéter la mécanique interne.

---

# 10. Pourquoi `pageHandler` n'est pas exportée

La fonction commence par une minuscule :

```go
pageHandler
```

Elle n'est donc accessible qu'à l'intérieur du package `handlers`.

C'est volontaire.

Le reste de l'application doit connaître :

```text
HomeHandler
ClubHandler
RulesHandler
WhereHandler
WhenHandler
ContactHandler
```

mais n'a aucune raison de connaître la mécanique utilisée pour construire une page générique.

`pageHandler` est un **détail d'implémentation du package handlers**.

---

# 11. Les vues

Le package :

```text
internal/views
```

est responsable du rendu HTML.

Pour les pages simples, on possède maintenant une seule fonction :

```go
RenderPage(w, data)
```

et une seule structure :

```go
PageData
```

Le template générique est :

```text
templates/pages/000_page.html
```

Il utilise notamment :

```html
<h1>{{ .Heading }}</h1>

<p>
    {{ .Description }}
</p>
```

---

# 12. Le layout commun

Toutes les pages utilisent :

```text
templates/layouts/base.html
```

Le layout contient les éléments communs :

```text
base.html
│
├── titre HTML
├── navigation
├── header
├── main
│    └── {{ template "content" . }}
└── footer
```

Les pages ne redéfinissent donc plus la structure générale du document HTML.

Elles fournissent uniquement leur contenu.

---

# 13. La navigation

Une navigation commune permet maintenant d'accéder aux différentes routes :

```text
Accueil
Le club
Où ?
Quand ?
Contact
Règlement
```

Comme elle appartient au layout commun, une modification de cette navigation est immédiatement visible sur toutes les pages.

Cela montre l'intérêt concret d'un layout partagé.

---

# 14. Le cas particulier de Contact

`Contact` n'est volontairement pas passé dans `pageHandler`.

Cette page possède des informations supplémentaires :

```text
EmailAddress
PhoneNumber
```

Elle conserve donc :

```text
ContactConfig
        ↓
ContactHandler
        ↓
ContactData
        ↓
RenderContact
        ↓
003_contact.html
```

C'est un point important :

> Factoriser ne signifie pas forcer tous les éléments à utiliser la même structure.

On factorise ce qui est réellement identique et on laisse les cas spécifiques rester spécifiques.

---

# 15. Parcours complet d'une page générique

Pour `/club` :

```text
GET /club
    │
    ▼
router
    │
    ▼
ClubHandler(cfg)
    │
    │ sélectionne cfg.Club
    ▼
pageHandler(
    cfg.SiteName,
    cfg.Club,
)
    │
    ▼
PageData
    │
    ▼
RenderPage
    │
    ├── base.html
    └── 000_page.html
    │
    ▼
HTML
    │
    ▼
navigateur
```

---

# 16. Parcours de la configuration

Au démarrage :

```text
config/config.json
        │
        ▼
     os.Open
        │
        ▼
json.Decoder
        │
        ▼
  config.Config
        │
        ▼
       main
        │
        ▼
      router
        │
        ▼
     handlers
```

La configuration est donc chargée une fois puis distribuée aux composants qui en ont besoin.

---

# 17. Évolution depuis le jalon 01

## Jalon 01

L'architecture mettait principalement en évidence :

```text
main
 ↓
router
 ↓
handlers
```

L'objectif était de sortir progressivement les responsabilités de `main.go`.

---

## Jalon 02

L'architecture actuelle est devenue :

```text
                config.json
                    │
                    ▼
                  Config
                    │
                    ▼
main ───────────► router
                    │
                    ▼
                 handlers
                    │
              ┌─────┴─────┐
              │           │
         pageHandler   ContactHandler
              │           │
              ▼           ▼
          PageData    ContactData
              │           │
              ▼           ▼
        RenderPage   RenderContact
              │           │
              └─────┬─────┘
                    ▼
                templates
                    │
                    ▼
                 HTML
```

La différence essentielle est que le projet commence maintenant à posséder une véritable **chaîne de transformation des données**.

---

# 18. Comparaison

|Architecture minimale|Jalon 01|Jalon 02|
|---|---|---|
|`main` gère presque tout|`main` démarre l'application|`main` assemble l'application|
|routes proches du démarrage|routeur dédié|routeur + injection de Config|
|handlers simples|handlers séparés|handlers spécialisés et factorisés|
|réponse HTTP directe|séparation route / traitement|données de vue dédiées|
|pas de configuration|préparation de l'architecture|configuration JSON chargée|
|pas de vue structurée|préparation des templates|package `views`|
|pas de layout|—|`base.html` commun|
|duplication acceptable|responsabilités séparées|répétitions factorisées|

---

# 19. Principe important observé

L'évolution du projet suit maintenant un principe intéressant :

```text
Écrire simplement
      ↓
Observer la répétition
      ↓
Comprendre ce qui est réellement commun
      ↓
Factoriser
```

On n'a pas commencé le projet directement avec :

```text
PageConfig
PageData
pageHandler
RenderPage
```

Ces abstractions sont apparues parce que plusieurs implémentations concrètes étaient devenues identiques.

Cela permet d'éviter une architecture inutilement complexe.

---

# 20. Responsabilités actuelles

## `main`

> Démarrer et assembler l'application.

---

## `config`

> Charger et représenter la configuration externe.

---

## `router`

> Associer une requête HTTP à un handler.

---

## handlers publics

> Choisir le comportement correspondant à une route.

Exemple :

```go
HomeHandler
```

sélectionne :

```go
cfg.Home
```

---

## `pageHandler`

> Implémenter le comportement commun des pages simples.

---

## `views`

> Préparer et exécuter les templates HTML.

---

## templates

> Décrire la présentation HTML.

---

# 21. Ce qui est volontairement absent

À ce stade, Club Manager ne possède toujours pas :

- PostgreSQL ;
    
- Goose ;
    
- sqlc ;
    
- authentification ;
    
- utilisateurs ;
    
- rôles et permissions ;
    
- gestion des adhérents ;
    
- logique métier complexe.
    

Ce n'est pas un manque dans ce jalon.

L'objectif actuel est de construire une base web claire avant d'introduire la persistance et la logique métier.

---

# 22. Petites limites de l'état actuel

Ce jalon représente également les imperfections actuelles.

Dans `internal/views/000_page.go`, certains commentaires proviennent encore de l'ancienne vue Home et parlent de :

```text
homeFiles
homeTemplate
RenderHome
```

alors que le code est désormais générique.

Ces commentaires devront être nettoyés.

Le fichier `base.html` place actuellement également la navigation dans `<head>`.

Pour une structure HTML correcte, la navigation devra être déplacée dans `<body>`, probablement dans le `<header>`.

Ces éléments constituent de petites corrections de finition et non une remise en cause de l'architecture.

---

# 23. Ce qui a été validé

À ce jalon :

- la configuration JSON est chargée ;
    
- les routes publiques fonctionnent ;
    
- la configuration est transmise au routeur puis aux handlers ;
    
- les pages génériques fonctionnent avec `pageHandler` ;
    
- la page Contact conserve son rendu spécifique ;
    
- le layout est partagé ;
    
- la navigation fonctionne sur l'ensemble des pages ;
    
- les tests existants passent.
    

---

# 24. Suite logique

L'architecture fonctionnelle des pages publiques est maintenant suffisamment stable pour commencer à travailler sur leur présentation.

La progression naturelle vers la V1 devient :

```text
Architecture des pages
        ✅
        │
        ▼
Nettoyage HTML / commentaires
        │
        ▼
Présentation / CSS / Bootstrap
        │
        ▼
Pages publiques propres
        │
        ▼
Tests et finition
        │
        ▼
       V1
```

La base de données viendra ensuite lorsque l'application commencera à gérer des données réellement dynamiques.

---

# Comprendre et retenir

### Pourquoi existe-t-il maintenant `PageConfig` ?

Parce que plusieurs configurations de pages avaient exactement la même structure.

---

### Pourquoi existe-t-il `PageData` en plus de `PageConfig` ?

Parce que les données stockées dans la configuration et les données nécessaires à l'affichage n'ont pas la même responsabilité.

---

### Pourquoi existe-t-il `pageHandler` ?

Parce que plusieurs handlers exécutaient exactement le même algorithme.

---

### Pourquoi garder `HomeHandler`, `ClubHandler`, etc. ?

Parce qu'ils donnent un sens métier aux routes et rendent le routeur immédiatement compréhensible.

---

### Pourquoi `pageHandler` commence-t-elle par une minuscule ?

Parce qu'elle est un détail interne au package `handlers`.

---

### Pourquoi Contact reste-t-il séparé ?

Parce qu'il possède des données et un affichage spécifiques.

---

### Quel principe architectural ressort de ce jalon ?

> Séparer les responsabilités, puis factoriser uniquement lorsqu'une répétition réelle apparaît.

---

# À retenir

```text
config.json
    ↓
Config
    ↓
main
    ↓
router
    ↓
handler public
    ↓
pageHandler
    ↓
PageData
    ↓
RenderPage
    ↓
template
    ↓
HTML
```

Club Manager n'est plus seulement une architecture HTTP découpée en plusieurs packages.

Le projet possède désormais une chaîne claire allant de la **configuration externe jusqu'à la génération de la page HTML**, avec des responsabilités distinctes et des abstractions apparues progressivement en réponse aux besoins réels du code.