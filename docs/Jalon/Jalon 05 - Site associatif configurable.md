
---

## project: Club Manager  
type: milestone  
milestone: 05  
git-tag: milestone-05-configurable-association-website  
git-commit: 7372490076bc7bb8bb58e324a7d7aa1545f40dc6  
date: 2026-08-08  
status: frozen  
previous: milestone-04-configurable-site

# Jalon 05 - Site associatif configurable

> [!info] État figé  
> Cette fiche représente Club Manager au commit :
> 
> `7372490076bc7bb8bb58e324a7d7aa1545f40dc6`
> 
> Tag associé :
> 
> `milestone-05-configurable-association-website`
> 
> Cette fiche constitue une photographie du projet et ne doit pas évoluer avec lui.

---

## Objectif atteint

Au jalon précédent, Club Manager avait commencé à sortir certaines données du code Go pour les placer dans un fichier de configuration.

Le mécanisme est maintenant généralisé.

Les principales informations propres à l'association sont définies dans :

```text
config/config.json
```

Le code Go décrit davantage **le fonctionnement de l'application**, tandis que le fichier JSON décrit davantage **le site associatif que cette application doit afficher**.

On obtient une première séparation nette entre :

```text
logiciel
   +
configuration
```

---

# Architecture générale

```text
                    config/config.json
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
          ┌────────────────┼────────────────┐
          ↓                ↓                ↓
     HomeHandler      ClubHandler      ContactHandler
          │                │                │
          └────────────────┼────────────────┘
                           │
                    autres handlers
                           │
                           ↓
                       ViewData
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

La configuration suit donc maintenant un trajet complet à travers l'application.

---

# Le fichier `config.json`

La configuration contient désormais l'identité générale du site :

```json
"site_name": "TCR Club Manager"
```

ainsi que les données de chaque page :

```text
home
club
contact
rules
where
when
```

Exemple :

```json
"home": {
    "title": "Accueil",
    "heading": "Bienvenue !",
    "description": "La Team Cat Ride (TCR) vous accueille sur son site."
}
```

Le fichier peut donc être modifié pour adapter le contenu du site sans changer les handlers.

---

# Une configuration structurée

La configuration JSON possède une représentation directe en Go :

```go
type Config struct {
    SiteName string
    Home     HomeConfig
    Club     ClubConfig
    Contact  ContactConfig
    Rules    RulesConfig
    Where    WhereConfig
    When     WhenConfig
}
```

Chaque partie de la configuration possède sa propre structure.

Par exemple :

```go
type HomeConfig struct {
    Title       string
    Heading     string
    Description string
}
```

Et la page de contact possède des informations supplémentaires :

```go
type ContactConfig struct {
    Title        string
    Heading      string
    Description  string
    EmailAddress string
    PhoneNumber  string
}
```

Le modèle Go reflète ainsi la structure du JSON.

---

# Correspondance JSON ↔ Go

Grâce aux tags :

```go
json:"home"
```

ou :

```go
json:"email_address"
```

Go sait faire correspondre les données du JSON aux champs des structures.

On peut représenter cela ainsi :

```text
config.json

"contact"
    │
    ↓
ContactConfig
    │
    ├── Title
    ├── Heading
    ├── Description
    ├── EmailAddress
    └── PhoneNumber
```

Le fichier JSON devient une représentation sérialisée de la configuration utilisée par l'application.

---

# Chargement au démarrage

`main` connaît le chemin du fichier :

```go
const configPath = "config/config.json"
```

Puis charge la configuration :

```go
cfg, err := config.Load(configPath)
```

Si le chargement échoue :

```go
log.Fatalf(...)
```

arrête immédiatement le programme.

C'est cohérent avec le rôle de la configuration :

> l'application ne doit pas démarrer si elle ne peut pas récupérer les informations nécessaires à son fonctionnement.

---

# `main` devient le point d'assemblage

Une fois chargée, la configuration est transmise au routeur :

```go
mux := router.New(cfg)
```

`main` ne connaît toujours pas le fonctionnement détaillé des pages.

Son rôle est plutôt d'assembler les grandes parties :

```text
configuration
      +
routeur
      +
serveur HTTP
```

On commence à voir apparaître une notion importante d'architecture :

> **le point d'entrée assemble les dépendances nécessaires à l'application.**

---

# Le routeur transmet la configuration

Le constructeur du routeur reçoit maintenant :

```go
func New(cfg config.Config) *http.ServeMux
```

Puis chaque handler reçoit cette configuration :

```go
mux.HandleFunc(
    "GET /{$}",
    handlers.HomeHandler(cfg),
)
```

Même principe pour :

```text
/club
/contact
/where
/when
/rules
```

Le routeur ne lit pas la configuration.

Il ne décide pas non plus de son contenu.

Il ne fait que la transmettre aux handlers qui en ont besoin.

---

# Pourquoi les handlers retournent-ils une fonction ?

Avant la configuration, un handler pouvait simplement être :

```go
func HomeHandler(
    w http.ResponseWriter,
    r *http.Request,
)
```

Maintenant, il faut lui fournir `cfg`.

On obtient donc :

```go
func HomeHandler(cfg config.Config) http.HandlerFunc
```

Cette fonction reçoit la configuration puis retourne le véritable handler HTTP.

Schématiquement :

```text
Config
  │
  ↓
HomeHandler(cfg)
  │
  ↓
http.HandlerFunc
  │
  ↓
requêtes HTTP futures
```

La closure permet au handler retourné de continuer à utiliser `cfg`.

---

# Le handler ne définit plus le contenu

Prenons la page Contact.

Le handler construit les données de vue à partir de la configuration :

```go
data := views.ContactData{
    SiteName:     cfg.SiteName,
    Title:        cfg.Contact.Title + " - " + cfg.SiteName,
    Heading:      cfg.Contact.Heading,
    Description:  cfg.Contact.Description,
    EmailAddress: cfg.Contact.EmailAddress,
    PhoneNumber:  cfg.Contact.PhoneNumber,
}
```

Le handler conserve donc une responsabilité importante :

```text
Config
   ↓
adaptation
   ↓
ViewData
```

Il fait le lien entre les données de l'application et les données attendues par la vue.

---

# Une donnée de bout en bout

Prenons par exemple :

```json
"email_address": "teamcatride@miaoumail.com"
```

Son trajet est :

```text
config.json
    │
    ↓
json.Decoder
    │
    ↓
Config.Contact.EmailAddress
    │
    ↓
ContactHandler
    │
    ↓
ContactData.EmailAddress
    │
    ↓
template
    │
    ↓
HTML envoyé au navigateur
```

C'est une chaîne complète de circulation de données.

---

# Le rôle des différentes couches

## `config`

Responsabilité :

```text
charger
+
représenter
```

la configuration de l'application.

Il ne connaît pas HTTP.

Il ne connaît pas les templates.

---

## `main`

Responsabilité :

```text
charger les dépendances
+
assembler l'application
+
démarrer le serveur
```

---

## `router`

Responsabilité :

```text
requête HTTP
      ↓
bon handler
```

Il transmet également la configuration reçue lors de sa construction.

---

## `handlers`

Responsabilité :

```text
configuration
      +
requête HTTP
      ↓
données nécessaires à la vue
```

Ils organisent le traitement d'une page.

---

## `views`

Responsabilité :

```text
ViewData
   ↓
templates
   ↓
HTML
```

---

## `templates`

Responsabilité :

```text
présentation
```

Ils ne chargent ni la configuration ni les données eux-mêmes.

---

# Évolution entre les jalons 04 et 05

## Jalon 04

Le principe de configuration existe.

```text
Where    ✅ configurable
When     ✅ configurable

autres pages
         encore partiellement liées au code
```

Le projet va **vers** un site configurable.

---

## Jalon 05

Le principe est généralisé :

```text
SiteName   ✅
Home       ✅
Club       ✅
Contact    ✅
Rules      ✅
Where      ✅
When       ✅
```

On peut donc résumer la progression :

```text
Jalon 04
expérimentation / transition
        ↓
Jalon 05
généralisation
```

---

# Changer d'association

C'est probablement l'une des conséquences les plus importantes de ce jalon.

Le code utilise actuellement une association fictive :

```text
Team Cat Ride
```

Mais le contenu peut être remplacé dans `config.json`.

Par exemple :

```json
{
    "site_name": "Mon Association",

    "home": {
        "title": "Accueil",
        "heading": "Bienvenue",
        "description": "Bienvenue sur le site de notre association."
    }
}
```

Le code Go n'a pas besoin d'être réécrit simplement parce que l'identité de l'association change.

C'est le début du caractère **générique** de Club Manager.

---

# Ce qui appartient maintenant à la configuration

À ce stade :

```text
nom du site
titres des pages
titres principaux
descriptions
email
téléphone
lieu
horaires sous forme descriptive
```

sont des données de configuration.

---

# Ce qui n'appartient pas à la configuration

Cette étape permet aussi de commencer à tracer une frontière.

Des données comme :

```text
membres
cotisations
paiements
comptabilité
horaires variables
inscriptions
rôles
```

ne doivent pas nécessairement devenir de grosses sections dans `config.json`.

Ces données :

- peuvent changer régulièrement ;
    
- peuvent être nombreuses ;
    
- peuvent devoir être recherchées ;
    
- peuvent avoir un historique ;
    
- peuvent être reliées entre elles.
    

Elles conduiront plutôt vers une base de données.

---

# Configuration ≠ base de données

Cette distinction devient importante pour la suite du projet.

## Configuration

Adaptée à des données :

```text
peu nombreuses
rarement modifiées
nécessaires au démarrage
liées au paramétrage du site
```

Exemples :

```text
nom du site
contenu institutionnel
coordonnées
textes simples
```

---

## Base de données

Adaptée à des données :

```text
nombreuses
dynamiques
recherchables
reliées
historiques
```

Exemples futurs :

```text
Member
Membership
Payment
Schedule
Role
```

Le fichier de configuration et la future base de données répondent donc à des besoins différents.

---

# Les tests évoluent avec l'architecture

Le chargement de la configuration est testé avec un fichier indépendant :

```text
internal/config/testdata/config.json
```

Le test construit ensuite une configuration attendue :

```go
expected := Config{
    ...
}
```

et compare :

```go
if cfg != expected
```

Cela permet de vérifier le contrat complet :

```text
JSON
 ↓
Load()
 ↓
Config
```

et non plus seulement quelques champs isolés.

---

# Des données de test distinctives

Les tests utilisent également des valeurs comme :

```text
Accueil test
Le Club test
Contact test
Règlement test
Où test
Quand test
```

Cette différence est utile.

Un test ne vérifie plus seulement qu'un mot générique comme `Contact` apparaît dans la réponse.

Il permet davantage de vérifier que :

```text
valeur construite par le test
          ↓
        Config
          ↓
       Handler
          ↓
       réponse
```

fonctionne réellement.

---

# Une première injection de dépendance

Sans introduire de framework ni d'abstraction supplémentaire, le projet utilise maintenant un principe très général.

Le handler ne crée pas lui-même sa configuration.

Elle lui est fournie :

```go
handlers.HomeHandler(cfg)
```

On peut appeler cela une **injection de dépendance**.

Très simplement :

> une fonction reçoit ce dont elle a besoin au lieu d'aller le chercher elle-même.

Cela facilite notamment :

- la compréhension des dépendances ;
    
- les tests ;
    
- le remplacement des données ;
    
- l'évolution future du projet.
    

---

# Une architecture plus testable

Un handler peut être construit avec une configuration spécifique au test :

```go
cfg := config.Config{
    SiteName: "Site de test",
}
```

Puis :

```go
handler := HomeHandler(cfg)
```

Le test n'a donc pas besoin de modifier :

```text
config/config.json
```

pour tester le handler.

Cela limite les dépendances entre tests.

---

# Ce que nous n'avons volontairement pas abstrait

Les structures :

```text
HomeConfig
ClubConfig
RulesConfig
WhereConfig
WhenConfig
```

se ressemblent beaucoup.

De même, plusieurs handlers suivent presque exactement le même schéma.

Cette duplication est actuellement visible.

Elle n'est pas nécessairement un problème à corriger immédiatement.

Au contraire, elle permet de comprendre concrètement :

```text
ce qui est identique
```

et :

```text
ce qui est réellement différent
```

avant d'introduire une abstraction.

C'est une bonne base pour une future réflexion sur :

> **quand la duplication justifie-t-elle un refactoring ?**

---

# Limites actuelles

Le système de configuration fonctionne, mais reste volontairement simple.

Par exemple :

```text
champ JSON absent
        ↓
valeur zéro Go
```

Un champ texte absent donnera généralement :

```go
""
```

Le chargement vérifie actuellement que le fichier existe et que le JSON peut être décodé.

Il ne valide pas encore toutes les règles métier possibles de la configuration.

Cela pourra faire l'objet d'une amélioration ultérieure.

---

# Concepts Go mobilisés

Ce jalon rassemble déjà de nombreux concepts :

```text
struct
struct imbriquées
tags JSON
encoding/json
os.File
gestion des erreurs
%w
fonctions
valeurs de retour
closures
http.HandlerFunc
ServeMux
constructeur New()
embed.FS
html/template
io.Writer
tests
httptest
testdata
```

Ils ne sont plus utilisés isolément.

Ils coopèrent maintenant dans une architecture complète.

---

# Vision du projet à ce jalon

Club Manager peut être représenté ainsi :

```text
                ┌──────────────────┐
                │   config.json    │
                └────────┬─────────┘
                         ↓
                ┌──────────────────┐
                │      Config      │
                └────────┬─────────┘
                         ↓
┌────────┐      ┌──────────────────┐
│ Client │ ───→ │      Router      │
└────────┘      └────────┬─────────┘
                         ↓
                ┌──────────────────┐
                │     Handler      │
                └────────┬─────────┘
                         ↓
                ┌──────────────────┐
                │     ViewData     │
                └────────┬─────────┘
                         ↓
                ┌──────────────────┐
                │       View       │
                └────────┬─────────┘
                         ↓
                ┌──────────────────┐
                │     Template     │
                └────────┬─────────┘
                         ↓
                      HTML
```

---

# Progression des jalons

```text
01 - Handlers et routes
        │
        ↓
séparation du routage et du traitement HTTP

02 - Premier template
        │
        ↓
séparation du traitement et de la présentation

03 - Layout commun
        │
        ↓
séparation de la structure commune
et du contenu des pages

04 - Vers un site configurable
        │
        ↓
apparition d'une configuration externe

05 - Site associatif configurable
        │
        ↓
généralisation de la configuration
à l'ensemble du socle
```

On voit apparaître une progression générale :

```text
séparer les responsabilités
          ↓
faire circuler les données
          ↓
rendre les composants indépendants
          ↓
rendre l'application configurable
```

---

# Comprendre et retenir

> **Le code Go définit maintenant principalement le comportement de Club Manager ; `config.json` définit une partie importante de l'identité et du contenu de l'association.**

---

> **La configuration est chargée une seule fois au démarrage puis transmise aux composants qui en ont besoin.**

```text
config.json
    ↓
Config
    ↓
main
    ↓
router
    ↓
handlers
```

---

> **Un composant ne devrait pas forcément aller chercher lui-même ses dépendances. Elles peuvent lui être fournies.**

```go
HomeHandler(cfg)
```

est un premier exemple simple de ce principe.

---

> **Une configuration et une base de données ne répondent pas au même besoin.**

```text
configuration
→ données peu nombreuses et relativement stables

base de données
→ données dynamiques, nombreuses et structurées
```

---

> **Le jalon 05 marque le moment où Club Manager commence réellement à devenir générique.**

Changer l'association ne signifie plus nécessairement modifier le code Go.

---

# État du jalon

```text
Serveur HTTP                         ✅
Routeur                              ✅
Routes limitées à GET                ✅

Handlers séparés                     ✅
Views séparées                       ✅
Templates                            ✅
Layout commun                        ✅

Chargement JSON                      ✅
Config structurée                    ✅
Config transmise depuis main         ✅
Config transmise aux handlers        ✅

SiteName configurable                ✅
Home configurable                    ✅
Club configurable                    ✅
Contact configurable                 ✅
Rules configurable                   ✅
Where configurable                   ✅
When configurable                    ✅

Tests de chargement                  ✅
Testdata indépendant                 ✅
Handlers testables avec Config       ✅

Validation avancée de Config         ❌
Base de données                      ❌
Membres                              ❌
Cotisations                          ❌
Authentification                     ❌
Rôles / permissions                  ❌
```

---

# Conclusion du jalon

Le jalon 04 avait introduit une question :

> Peut-on sortir les informations propres à l'association du code Go ?

Le jalon 05 apporte une première réponse complète :

```text
oui
```

pour le contenu relativement statique du site.

Club Manager n'est désormais plus seulement :

```text
un site écrit pour une association
```

Il commence à devenir :

```text
un logiciel capable de représenter
différentes associations
à partir d'une configuration
```

**Ce jalon marque le passage d'un site contenant des données configurables à un véritable socle de site associatif configurable.**