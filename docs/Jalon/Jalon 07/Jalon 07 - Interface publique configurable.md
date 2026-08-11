
---


## Objectif du jalon

Ce jalon représente l'arrivée d'une première véritable interface publique pour Club Manager.

L'application ne se contente plus :

- de recevoir des requêtes HTTP ;
    
- de sélectionner un handler ;
    
- de charger des données depuis la configuration ;
    
- de générer du HTML avec des templates.
    

Elle sait maintenant également :

- servir des fichiers statiques ;
    
- utiliser Bootstrap ;
    
- proposer une navigation responsive ;
    
- partager une mise en page commune ;
    
- afficher des images définies depuis la configuration.
    

Ce jalon marque donc le passage :

```text
application Go générant des pages
```

vers :

```text
petit site public configurable
```

---

# 1. Vue générale

L'architecture peut maintenant être représentée ainsi :

```text
                    config/config.json
                           │
                           ▼
                     config.Load()
                           │
                           ▼
                        Config
                           │
                           ▼
                         main
                           │
                           ▼
                     router.New(cfg)
                           │
             ┌─────────────┴─────────────┐
             │                           │
             ▼                           ▼
       routes des pages             /static/
             │                           │
             ▼                           ▼
     handlers publics               FileServer
             │                           │
             ▼                           ▼
       pageHandler()                 static/
             │                      ├── css/
             ▼                      └── images/
         PageData
             │
             ▼
         RenderPage()
             │
             ▼
     templates Go + Bootstrap
             │
             ▼
            HTML
             │
             ▼
         navigateur
```

---

# 2. Le fichier de configuration évolue

Les pages génériques reposent maintenant sur :

```go
type PageConfig struct {
    Title       string `json:"title"`
    Heading     string `json:"heading"`
    Description string `json:"description"`
    Image       string `json:"image"`
    ImageAlt    string `json:"image_alt"`
}
```

Une page peut donc contenir :

```json
"club": {
    "title": "Le Club",
    "heading": "Team Cat Ride",
    "description": "TCR est là pour t'accompagner et te faire progresser en vélo.",
    "image": "/static/images/club.jpeg",
    "image_alt": "L'ensemble des accompagnateurs"
}
```

Le JSON ne contient pas l'image elle-même.

Il contient :

```text
/static/images/club.jpeg
```

c'est-à-dire l'URL permettant au navigateur de demander cette image.

---

# 3. La configuration reste indépendante de la présentation

Le fichier JSON définit :

```text
quoi afficher
```

Le template décide :

```text
comment l'afficher
```

Par exemple :

```json
"image": "/static/images/club.jpeg"
```

ne contient aucune information concernant :

- la taille de l'image ;
    
- ses marges ;
    
- son comportement responsive ;
    
- son positionnement.
    

Ces décisions appartiennent au template HTML.

Cette séparation est importante :

```text
Config
    ↓
contenu

Template
    ↓
présentation
```

---

# 4. `PageConfig` vers `PageData`

Le handler générique réalise toujours la transformation entre la configuration et les données destinées à la vue.

```go
data := views.PageData{
    SiteName:    siteName,
    Title:       page.Title + " - " + siteName,
    Heading:     page.Heading,
    Description: page.Description,
    Image:       page.Image,
    ImageAlt:    page.ImageAlt,
}
```

Le parcours devient :

```text
PageConfig
    │
    ▼
pageHandler
    │
    ▼
PageData
```

Même si les structures se ressemblent, leurs responsabilités restent différentes.

### `PageConfig`

Représente les informations venant de la configuration.

### `PageData`

Représente les informations envoyées au template.

---

# 5. Une image facultative

Le template générique utilise maintenant :

```html
{{ if .Image }}

<img
    src="{{ .Image }}"
    alt="{{ .ImageAlt }}"
    class="img-fluid mb-4"
>

{{ end }}
```

On découvre ici une nouvelle fonctionnalité des templates Go :

```go
{{ if .Image }}
```

Si `Image` contient une valeur :

```text
image présente
    ↓
<img> générée
```

Si `Image` est vide :

```text
aucune image
    ↓
aucune balise <img>
```

Une page générique peut donc avoir ou non une image sans créer un nouveau template.

---

# 6. `ImageAlt`

Le champ :

```go
ImageAlt string
```

permet de produire :

```html
alt="{{ .ImageAlt }}"
```

Il fournit une description textuelle de l'image lorsqu'elle apporte une information.

Exemple :

```json
"image_alt": "L'ensemble des accompagnateurs"
```

Pour une image purement décorative, ce champ peut rester vide.

Le HTML généré possède alors :

```html
alt=""
```

La possibilité est donc présente sans rendre la configuration inutilement contraignante.

---

# 7. Les fichiers statiques

Une nouvelle partie apparaît dans l'arborescence du projet :

```text
static/
├── css/
│   └── style.css
│
└── images/
    ├── logo.jpeg
    └── club.jpeg
```

Contrairement aux templates, ces fichiers sont directement demandés par le navigateur.

---

# 8. Le routeur sert maintenant les fichiers statiques

Le routeur possède :

```go
staticFiles := http.FileServer(
    http.Dir("static"),
)

mux.Handle(
    "GET /static/",
    http.StripPrefix(
        "/static/",
        staticFiles,
    ),
)
```

Une requête :

```text
GET /static/images/club.jpeg
```

suit le parcours :

```text
GET /static/images/club.jpeg
            │
            ▼
        ServeMux
            │
            ▼
 StripPrefix("/static/")
            │
            ▼
    images/club.jpeg
            │
            ▼
 FileServer("static")
            │
            ▼
static/images/club.jpeg
```

Le serveur peut ainsi fournir :

- CSS ;
    
- images ;
    
- logos ;
    
- autres ressources statiques futures.
    

---

# 9. Une distinction importante

Le navigateur ne connaît pas :

```text
static/images/club.jpeg
```

comme chemin du projet.

Il connaît seulement :

```text
/static/images/club.jpeg
```

Il faut donc toujours distinguer :

```text
URL HTTP
```

et :

```text
chemin sur le disque
```

Le serveur fait le lien entre les deux.

---

# 10. Bootstrap

Bootstrap est maintenant chargé dans le layout commun.

Il fournit la majorité de la mise en forme générale :

```text
Bootstrap
│
├── responsive
├── navbar
├── containers
├── espacements
└── styles courants
```

Notre fichier :

```text
static/css/style.css
```

reste disponible pour les personnalisations spécifiques au club.

La philosophie retenue est :

```text
Bootstrap
    ↓
mise en forme standard

style.css
    ↓
personnalisation particulière
```

---

# 11. Pourquoi Bootstrap ?

Le projet Club Manager est avant tout un projet Go.

L'objectif de la V1 n'est pas de construire un système graphique complexe.

Bootstrap permet donc d'obtenir rapidement :

- une interface propre ;
    
- une navigation utilisable ;
    
- un comportement responsive ;
    
- des espacements cohérents ;
    
- une présentation standard.
    

Sans détourner le projet de son objectif principal.

---

# 12. Le layout commun

Le template `base.html` possède maintenant trois grandes parties :

```text
┌──────────────────────────────────────┐
│ Navbar                               │
├──────────────────────────────────────┤
│                                      │
│ Main                                 │
│                                      │
├──────────────────────────────────────┤
│ Footer                               │
└──────────────────────────────────────┘
```

Toutes les pages utilisant ce layout bénéficient automatiquement de la même structure.

---

# 13. La navbar

La navbar Bootstrap contient :

- le logo ;
    
- le nom du site ;
    
- les liens vers les différentes pages ;
    
- un bouton responsive pour les petits écrans.
    

Le nom du site reste fourni par :

```go
SiteName
```

et affiché avec :

```go
{{ .SiteName }}
```

Le logo est servi depuis :

```text
/static/images/logo.jpeg
```

---

# 14. Responsive

La classe :

```html
navbar-expand-lg
```

permet à Bootstrap d'adapter la navigation à la largeur de l'écran.

Écran large :

```text
[logo] TCR Club Manager   Accueil   Le club   Où ?   Quand ? ...
```

Petit écran :

```text
[logo] TCR Club Manager                         ☰
```

Le JavaScript Bootstrap assure l'ouverture et la fermeture du menu.

---

# 15. Le `container`

La navbar, le contenu principal et le footer utilisent le même principe de `container`.

Cela permet d'obtenir un alignement commun :

```text
       même limite
           ↓

Navbar     │ ...
Main       │ ...
Footer     │ ...
```

Le contenu reste ainsi centré et lisible même sur de grands écrans.

---

# 16. Le contenu principal

Le contenu utilise :

```html
<main class="container py-4">
```

`container` gère principalement :

```text
largeur
+
alignement horizontal
```

`py-4` ajoute :

```text
padding vertical
```

donc de l'espace :

```text
au-dessus
+
en dessous
```

du contenu.

---

# 17. Le footer

Le footer reste volontairement simple.

Il utilise notamment :

```text
border-top
py-3
container
text-body-secondary
```

L'objectif n'est pas encore d'ajouter :

- réseaux sociaux ;
    
- nombreux liens ;
    
- widgets ;
    
- informations complexes.
    

Pour la V1, un footer sobre suffit.

---

# 18. Les pages génériques restent génériques

L'ajout d'une image n'a pas nécessité :

```text
ClubTemplate
HomeTemplate
RulesTemplate
WhereTemplate
WhenTemplate
```

Nous conservons :

```text
PageConfig
     ↓
pageHandler
     ↓
PageData
     ↓
000_page.html
```

L'image est simplement devenue une capacité supplémentaire de la page générique.

C'est important :

> nous avons enrichi l'abstraction existante sans créer une nouvelle abstraction.

---

# 19. Contact reste spécifique

La page Contact reste différente des pages génériques.

Elle possède notamment :

```text
EmailAddress
PhoneNumber
```

Cela justifie toujours :

```text
ContactConfig
ContactData
ContactHandler
template Contact
```

Toutes les pages ne doivent pas nécessairement entrer dans le même modèle.

---

# 20. Éditer le site sans modifier Go

L'architecture commence maintenant à permettre une distinction intéressante entre :

```text
développeur
```

et :

```text
personne qui édite le contenu
```

Pour modifier par exemple :

- le nom du site ;
    
- un titre ;
    
- une description ;
    
- une image ;
    

il n'est plus nécessaire de modifier les handlers Go.

Le contenu peut être changé depuis :

```text
config/config.json
```

et les images peuvent être déposées dans :

```text
static/images/
```

---

# 21. Limite actuelle de la configuration

Le JSON reste pratique pour :

- textes courts ;
    
- coordonnées ;
    
- chemins d'images ;
    
- valeurs structurées.
    

Il devient moins agréable pour rédiger :

- plusieurs paragraphes ;
    
- des listes ;
    
- de longs règlements ;
    
- des articles ;
    
- des contenus avec mise en forme.
    

Une évolution possible plus tard serait donc :

```text
JSON
    ↓
données structurées

Markdown
    ↓
contenus rédactionnels
```

Mais cette évolution n'est pas nécessaire à la V1.

Elle répondra à un besoin seulement lorsqu'il apparaîtra réellement.

---

# 22. Ce qui n'est volontairement pas présent

À ce jalon, le projet ne possède toujours pas :

- PostgreSQL ;
    
- migrations ;
    
- sqlc ;
    
- authentification ;
    
- utilisateurs ;
    
- rôles ;
    
- administration ;
    
- éditeur Markdown ;
    
- système de téléchargement d'images ;
    
- galerie ;
    
- contenu dynamique complexe.
    

Ce n'est pas un manque.

C'est une décision d'architecture.

La V1 doit rester simple.

---

# 23. Évolution depuis les premiers jalons

L'application minimale pouvait être représentée par :

```text
requête
   ↓
handler
   ↓
texte
```

L'architecture actuelle ressemble maintenant davantage à :

```text
                  configuration
                       │
                       ▼
requête → router → handler
                       │
                       ▼
                    données
                       │
                       ▼
                    template
                       │
                ┌──────┴──────┐
                ▼             ▼
            Bootstrap       static
                              │
                         ┌────┴────┐
                         ▼         ▼
                        CSS      images
                │
                ▼
               HTML
                │
                ▼
            navigateur
```

Chaque responsabilité commence à avoir un emplacement identifiable.

---

# 24. Responsabilités actuelles

## `main`

Démarre l'application et assemble les dépendances.

## `config`

Charge et représente la configuration.

## `router`

Associe les requêtes HTTP aux handlers et sert les fichiers statiques.

## `handlers`

Préparent les données nécessaires aux pages.

## `views`

Gèrent les données destinées à l'affichage et exécutent les templates.

## `templates`

Définissent la structure HTML.

## `static`

Contient les ressources directement servies au navigateur.

## Bootstrap

Fournit la majorité de la présentation standard.

---

# 25. Principe architectural observé

L'évolution jusqu'à ce jalon suit toujours le même principe :

```text
Écrire simplement
        ↓
observer un besoin
        ↓
comprendre la responsabilité
        ↓
ajouter le minimum nécessaire
```

Exemples :

```text
handlers répétés
      ↓
pageHandler
```

```text
pages similaires
      ↓
PageConfig + PageData
```

```text
besoin de CSS et d'images
      ↓
static/
```

```text
besoin d'une interface propre rapidement
      ↓
Bootstrap
```

```text
besoin d'illustrer une page
      ↓
Image + ImageAlt
```

---

# Comprendre et retenir

## Où est le contenu ?

Principalement dans :

```text
config/config.json
```

---

## Où sont les images ?

Dans :

```text
static/images/
```

---

## Que contient le JSON pour une image ?

Pas l'image.

Seulement son chemin public :

```text
/static/images/club.jpeg
```

---

## Qui sert l'image ?

```text
http.FileServer
```

via la route :

```text
/static/
```

---

## Qui décide si l'image doit être affichée ?

Le template :

```go
{{ if .Image }}
```

---

## Qui décide de son apparence ?

Le HTML et Bootstrap.

---

## À quoi sert Bootstrap ?

À fournir rapidement une interface standard, responsive et propre sans faire du frontend le sujet principal du projet.

---

## Pourquoi ne pas ajouter Markdown maintenant ?

Parce que la V1 n'en a pas encore besoin.

Le JSON suffit actuellement pour les informations simples et structurées.

---

# État du projet au jalon 7

```text
Club Manager
│
├── serveur HTTP Go
│
├── configuration JSON
│
├── injection de la configuration
│
├── routeur
│
│   ├── pages
│
│   └── fichiers statiques
│
├── handlers
│   ├── handlers publics
│   ├── pageHandler générique
│   └── Contact spécifique
│
├── views
│   ├── PageData
│   └── ContactData
│
├── templates
│   ├── layout commun
│   ├── pages génériques
│   └── page Contact
│
├── Bootstrap
│   ├── navbar responsive
│   ├── containers
│   └── footer
│
└── static
    ├── css
    └── images
```

---

# À retenir

Ce jalon marque une évolution importante :

> Club Manager possède maintenant le socle d'un site public simple, configurable et présentable.

La configuration fournit le contenu.

Go organise et transforme les données.

Les templates construisent le HTML.

Bootstrap fournit une présentation standard.

Le serveur de fichiers statiques fournit les ressources comme les images et le CSS.

Le point essentiel reste cependant la simplicité :

> La V1 doit démontrer une architecture claire et fonctionnelle, pas accumuler des fonctionnalités.