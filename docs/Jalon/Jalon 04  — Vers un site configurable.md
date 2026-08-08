
---


> [!info] Objectif du jalon  
> Cette fiche représente un **état figé de Club Manager**.
> 
> Elle ne doit pas évoluer avec le projet.
> 
> Son objectif est de montrer comment l'architecture s'est enrichie depuis le jalon précédent et pourquoi ces nouvelles responsabilités sont apparues.

---

## État du projet

Club Manager n'est plus seulement un serveur Go capable d'afficher plusieurs pages.

Le projet commence maintenant à devenir une **application configurable pour une association**.

L'idée de la première version est simple :

> Une petite association doit pouvoir disposer d'un site fonctionnel en modifiant principalement un fichier de configuration, sans modifier le code Go.

Les pages prévues pour cette première version sont :

- accueil ;
    
- présentation du club ;
    
- contact ;
    
- règlement intérieur ;
    
- lieu ;
    
- horaires.
    

---

# Architecture actuelle

```text
club-manager/
│
├── cmd/
│   └── server/
│       └── main.go
│
├── config/
│   └── config.json
│
├── internal/
│   ├── config/
│   │
│   ├── handlers/
│   │
│   ├── router/
│   │
│   └── views/
│       └── templates/
│           ├── layouts/
│           └── pages/
│
└── go.mod
```

Chaque dossier possède progressivement une responsabilité précise.

---

# Le chemin d'une requête

L'architecture HTTP peut être représentée ainsi :

```text
Navigateur
    │
    │ requête HTTP
    ↓
 Router
    │
    ↓
 Handler
    │
    ↓
 View
    │
    ↓
 Template
    │
    ↓
Réponse HTML
```

Mais une nouvelle chaîne est maintenant apparue : celle de la configuration.

```text
config.json
    │
    ↓
config.Load()
    │
    ↓
Config
    │
    ↓
main
    │
    ↓
router
    │
    ↓
handlers
```

Les deux chaînes se rejoignent donc dans les handlers :

```text
              config.json
                   │
                   ↓
                  Config
                   │
                   ↓
Navigateur → Router → Handler → View → Template
```

---

# `main` devient le point d'assemblage

Le rôle de `main` est maintenant plus clair.

Il :

1. charge la configuration ;
    
2. vérifie qu'elle peut être chargée ;
    
3. transmet la configuration au routeur ;
    
4. démarre le serveur HTTP.
    

Schématiquement :

```go
func main() {
    cfg, err := config.Load(configPath)

    // gestion de l'erreur

    mux := router.New(cfg)

    http.ListenAndServe(":8080", mux)
}
```

`main` ne connaît pas le contenu des pages.

Il assemble les différentes parties de l'application.

---

# Apparition de `Config`

Auparavant, les informations concernant l'association étaient directement présentes dans le code.

Exemple conceptuel :

```go
Description: "Nous sommes au 3 rue de Perlinpinpin"
```

Le problème est que cette information appartient à **l'association**, pas au fonctionnement de Club Manager.

Elle commence donc à être déplacée dans :

```text
config/config.json
```

Exemple :

```json
{
    "site_name": "Club Manager",

    "where": {
        "heading": "Où nous trouver ?",
        "description": "Nous sommes au 3 rue de Perlinpinpin"
    },

    "when": {
        "heading": "Quand nous trouver ?",
        "description": "..."
    }
}
```

---

# Structures imbriquées

Le fichier JSON possède une structure hiérarchique.

Elle est représentée en Go par des structures imbriquées :

```go
type Config struct {
    SiteName string
    Where    WhereConfig
    When     WhenConfig
}
```

avec par exemple :

```go
type WhereConfig struct {
    Heading     string
    Description string
}
```

Cela permet d'accéder naturellement aux données :

```go
cfg.SiteName

cfg.Where.Heading
cfg.Where.Description

cfg.When.Heading
cfg.When.Description
```

---

# Le handler change de responsabilité

Avant :

```text
WhereHandler
│
├── connaît le titre
├── connaît l'adresse
├── construit les données
└── appelle la vue
```

Maintenant :

```text
WhereHandler
│
├── reçoit Config
├── récupère les données nécessaires
├── construit WhereData
└── appelle la vue
```

Le handler **ne décide plus du contenu de l'association**.

Il organise le traitement de la requête.

C'est une séparation des responsabilités supplémentaire.

---

# Exemple du trajet d'une donnée

Prenons l'adresse du club.

### 1. Configuration

```json
"where": {
    "description": "Nous sommes au 3 rue de Perlinpinpin"
}
```

### 2. Décodage JSON

Cette valeur devient :

```go
cfg.Where.Description
```

### 3. Handler

Le handler transmet cette valeur :

```go
Description: cfg.Where.Description,
```

### 4. Vue

La donnée est placée dans une structure destinée au template.

### 5. Template

```html
<p>
    {{ .Description }}
</p>
```

### 6. Navigateur

L'utilisateur voit finalement :

```text
Nous sommes au 3 rue de Perlinpinpin
```

---

# Évolution des responsabilités

## `config`

Responsabilité :

> Charger et représenter la configuration de l'application.

Il ne gère pas HTTP.

Il ne génère pas de HTML.

---

## `router`

Responsabilité :

> Associer une URL à un handler.

Exemple conceptuel :

```text
/        → HomeHandler
/club    → ClubHandler
/contact → ContactHandler
/rules   → RulesHandler
/where   → WhereHandler
/when    → WhenHandler
```

---

## `handlers`

Responsabilité :

> Traiter une requête HTTP et préparer les données nécessaires à la réponse.

Un handler fait progressivement le lien entre :

```text
Config
  +
requête HTTP
  ↓
données de vue
```

---

## `views`

Responsabilité :

> Transformer les données reçues en représentation HTML.

La vue connaît :

- les structures de données nécessaires au template ;
    
- le template à utiliser ;
    
- la manière de rendre ce template.
    

---

## `templates`

Responsabilité :

> Décrire la présentation HTML.

Ils ne doivent pas contenir la logique métier de l'application.

---

# Les tests évoluent également

Les tests ne vérifient plus seulement qu'une fonction peut être appelée.

Ils commencent notamment à vérifier :

- le chargement d'un fichier JSON ;
    
- les erreurs de chargement ;
    
- le décodage de structures imbriquées ;
    
- les réponses des handlers ;
    
- le contenu produit.
    

Le fichier :

```text
internal/config/testdata/config.json
```

permet également aux tests de disposer de leurs propres données.

Ainsi :

```text
configuration réelle
        ≠
configuration utilisée par les tests
```

---

# Ce que ce jalon apporte par rapport au précédent

Le projet précédent pouvait être résumé par :

```text
Route
 ↓
Handler
 ↓
View
 ↓
Template
```

Le projet actuel devient :

```text
             Configuration
                   ↓
                 Config
                   ↓
Route ─────────→ Handler
                   ↓
                  Data
                   ↓
                  View
                   ↓
                Template
```

Une nouvelle responsabilité est donc apparue :

> **La personnalisation de l'application est séparée de son code.**

---

# Une première frontière entre logiciel et données

On commence également à distinguer deux choses.

### Club Manager

Le logiciel :

```text
routes
handlers
views
templates
chargement de configuration
```

### L'association

Ses informations :

```text
nom
adresse
horaires
contact
présentation
règlement
```

L'objectif est progressivement d'obtenir :

```text
Club Manager
     +
config.json d'une association
     =
site de cette association
```

---

# Ce qui reste volontairement imparfait

Toutes les pages ne sont pas encore entièrement configurables.

Certaines informations sont encore écrites directement dans les handlers.

C'est notamment le cas de contenus tels que :

- contact ;
    
- règlement ;
    
- présentation du club ;
    
- accueil.
    

Ce n'est pas un problème.

Ce jalon représente précisément **la période de transition** durant laquelle le principe de configuration vient d'être introduit et commence à être généralisé.

---

# Une duplication commence à apparaître

Les structures :

```go
WhereConfig
WhenConfig
```

se ressemblent fortement.

Les handlers associés possèdent également une structure très proche.

On pourrait être tenté de tout généraliser immédiatement.

Pour l'instant, cette duplication est volontairement conservée.

Elle permettra plus tard de poser une vraie question de conception :

> À partir de quel moment une duplication justifie-t-elle la création d'une abstraction ?

Le futur refactoring sera ainsi motivé par un problème réellement observé dans le projet.

---

# Concepts Go mis en pratique

À ce stade, Club Manager permet déjà de manipuler concrètement :

- packages ;
    
- fonctions ;
    
- méthodes de la bibliothèque standard ;
    
- structs ;
    
- structs imbriquées ;
    
- pointeurs ;
    
- closures ;
    
- gestion des erreurs ;
    
- JSON ;
    
- fichiers ;
    
- `io.Writer` ;
    
- HTTP ;
    
- handlers ;
    
- `http.HandlerFunc` ;
    
- `ServeMux` ;
    
- templates HTML ;
    
- `embed.FS` ;
    
- tests ;
    
- logs.
    

Ces notions ne sont plus seulement étudiées indépendamment.

Elles participent à une même application.

---

# Lien avec la préparation aux concours

Club Manager commence à jouer le rôle de support pratique.

Par exemple :

```text
Structure de données
→ Config et structs imbriquées

Entrées / sorties
→ lecture de config.json

Encodage de données
→ JSON

Réseau
→ HTTP

Architecture logicielle
→ séparation des responsabilités

Tests
→ validation du comportement

Gestion des erreurs
→ chargement de configuration et rendu

Git
→ évolution progressive et jalons
```

Une notion théorique peut ainsi être associée à une réalisation concrète.

---

# Prochaine direction

La prochaine phase consiste à terminer le principe qui vient d'être introduit :

```text
where    ✅
when     ✅

contact  → à rendre configurable
rules    → à rendre configurable
club     → à rendre configurable
home     → à rendre configurable
```

On pourra alors obtenir :

```text
       config.json
           │
           ↓
     Club Manager
           │
           ↓
Site associatif personnalisé
```

À ce moment-là, le socle fonctionnel de la **V1 configurable** sera presque terminé.

Bootstrap pourra ensuite intervenir pour apporter une présentation commune et responsive.

---

# Comprendre et retenir

> **Un handler ne devrait pas contenir les informations propres à une association.**
> 
> Ces informations appartiennent à la configuration.

---

> **`main` devient le point d'assemblage de l'application.**
> 
> Il charge la configuration, construit le routeur et démarre le serveur.

---

> **Une struct Go peut représenter naturellement la hiérarchie d'un objet JSON.**

```text
JSON imbriqué
     ↓
structs imbriquées
```

---

> **La séparation des responsabilités s'améliore au fur et à mesure que le projet rencontre de nouveaux besoins.**
> 
> L'architecture n'est pas créée pour elle-même : elle apparaît pour résoudre des problèmes concrets.

---

> **Club Manager commence à séparer le logiciel de l'association qui l'utilise.**

```text
logiciel générique
      +
configuration spécifique
      =
site personnalisé
```

---

# État du jalon

```text
Serveur HTTP                 ✅
Router                       ✅
Handlers séparés             ✅
Views séparées               ✅
Templates                    ✅
Layout commun                ✅
Tests                        ✅
Chargement JSON              ✅
Config transmise au routeur  ✅
Config transmise aux handlers✅
Where configurable           ✅
When configurable            ✅

Contact configurable         ⏳
Rules configurable           ⏳
Club configurable            ⏳
Home configurable            ⏳
Bootstrap                    ⏳
Base de données              — plus tard
```

**Ce jalon marque le passage d'un site Go structuré à un véritable socle de site associatif configurable.**