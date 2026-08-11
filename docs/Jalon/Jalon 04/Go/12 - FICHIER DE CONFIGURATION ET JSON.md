

---

- [[12.01 - Charger et tester une configuration JSON en Go]]
- [[12.02 - Transmettre une configuration de main jusqu'aux handlers]]
# JSON et fichiers de configuration

## Objectif

JSON est un format texte permettant de représenter des données structurées.

Le nom JSON signifie :

```text
JavaScript Object Notation
```

Même si son origine est liée à JavaScript, JSON est aujourd’hui utilisé par de nombreux langages :

- Go ;
    
- Python ;
    
- Java ;
    
- JavaScript ;
    
- Rust ;
    
- PHP ;
    
- C#.
    

Dans Club Manager, JSON pourra servir à enregistrer des paramètres généraux dans un fichier extérieur au code.

Exemple :

```json
{
  "site_name": "Club Manager"
}
```

---

# Pourquoi utiliser un fichier JSON ?

Sans fichier de configuration, une information peut être écrite directement dans le code :

```go
siteName := "Club Manager"
```

Cette valeur est alors compilée avec le programme.

Pour la modifier, il faut généralement :

1. ouvrir le code source ;
    
2. modifier la valeur ;
    
3. recompiler l’application ;
    
4. relancer l’application.
    

Avec un fichier JSON :

```json
{
  "site_name": "Club Manager"
}
```

la donnée se trouve en dehors du code Go.

```text
Code Go
└── fonctionnement de l’application

Fichier JSON
└── paramètres de l’application
```

Cela permet de modifier certains paramètres sans modifier directement le code source.

---

# Un format d’échange de données

JSON est souvent utilisé pour transmettre des données entre plusieurs systèmes.

Exemple :

```text
Application A
     ↓
données JSON
     ↓
Application B
```

On le retrouve notamment dans :

- les fichiers de configuration ;
    
- les API web ;
    
- les échanges entre un frontend et un backend ;
    
- les fichiers exportés ;
    
- certaines bases de données ;
    
- les outils de développement.
    

JSON est populaire parce qu’il est :

- relativement simple à lire ;
    
- textuel ;
    
- indépendant d’un langage particulier ;
    
- facile à produire et à analyser automatiquement.
    

---

# Un objet JSON

La structure suivante est un objet JSON :

```json
{
  "site_name": "Club Manager"
}
```

Un objet JSON est entouré par des accolades :

```json
{
}
```

Il contient des associations entre des clés et des valeurs.

```text
clé → valeur
```

Dans notre exemple :

```text
site_name → Club Manager
```

---

## La clé

La clé est :

```json
"site_name"
```

Une clé JSON est toujours écrite entre guillemets doubles.

Elle permet d’identifier une donnée.

```text
site_name
└── nom utilisé pour retrouver la valeur
```

---

## La valeur

La valeur associée est :

```json
"Club Manager"
```

Cette valeur est une chaîne de caractères.

L’association complète est :

```json
"site_name": "Club Manager"
```

Le caractère `:` sépare la clé de sa valeur.

```text
clé : valeur
```

---

# Objet, dictionnaire et map

Selon le langage utilisé, une structure associant des clés à des valeurs peut porter différents noms.

|Contexte|Nom courant|
|---|---|
|JSON|Objet|
|Python|Dictionnaire|
|Go|`map` ou `struct`|
|JavaScript|Objet|
|Java|`Map`|
|Rust|`HashMap` ou structure|

Nous pouvons donc dire qu’un objet JSON ressemble à un dictionnaire.

Cependant, dans une documentation précise sur JSON, le terme correct est :

```text
objet JSON
```

---

# Plusieurs propriétés

Un objet JSON peut contenir plusieurs associations :

```json
{
  "site_name": "Club Manager",
  "contact_email": "contact@example.com",
  "logo": "logo.png"
}
```

Chaque association est appelée une propriété.

```text
Objet JSON
├── propriété site_name
├── propriété contact_email
└── propriété logo
```

Les propriétés sont séparées par des virgules :

```json
{
  "site_name": "Club Manager",
  "contact_email": "contact@example.com"
}
```

La dernière propriété ne doit pas être suivie d’une virgule.

Incorrect :

```json
{
  "site_name": "Club Manager",
}
```

Correct :

```json
{
  "site_name": "Club Manager"
}
```

---

# Les types de valeurs JSON

JSON possède un petit nombre de types.

## Chaîne de caractères

```json
{
  "site_name": "Club Manager"
}
```

Une chaîne est placée entre guillemets doubles.

---

## Nombre

```json
{
  "membership_price": 150
}
```

Un nombre n’utilise pas de guillemets.

```json
150
```

n’est pas la même chose que :

```json
"150"
```

Le premier est un nombre.

Le second est une chaîne de caractères.

---

## Booléen

```json
{
  "registration_open": true
}
```

Les deux valeurs booléennes sont :

```json
true
false
```

Elles sont écrites sans majuscule et sans guillemets.

---

## Valeur nulle

```json
{
  "logo": null
}
```

`null` signifie qu’aucune valeur n’est présente.

---

## Tableau

```json
{
  "colors": [
    "blue",
    "white",
    "black"
  ]
}
```

Un tableau est entouré de crochets :

```json
[
]
```

Ses valeurs sont séparées par des virgules.

---

## Objet imbriqué

Un objet peut contenir un autre objet :

```json
{
  "site": {
    "name": "Club Manager",
    "logo": "logo.png"
  }
}
```

La structure devient :

```text
objet principal
└── site
    ├── name
    └── logo
```

---

# Exemple plus complet

Un futur fichier de configuration pourrait ressembler à ceci :

```json
{
  "site_name": "Club Manager",
  "contact_email": "contact@example.com",
  "registration_open": true,
  "theme": {
    "primary_color": "blue",
    "secondary_color": "white"
  }
}
```

Cependant, il est préférable de commencer avec une configuration minimale :

```json
{
  "site_name": "Club Manager"
}
```

Nous ajouterons de nouvelles propriétés uniquement lorsqu’un besoin réel apparaîtra.

---

# JSON et structure Go

Le fichier JSON :

```json
{
  "site_name": "Club Manager"
}
```

peut correspondre à la structure Go suivante :

```go
type Config struct {
	SiteName string `json:"site_name"`
}
```

La structure Go possède un champ :

```go
SiteName string
```

Le fichier JSON possède une propriété :

```json
"site_name"
```

L’annotation suivante crée le lien entre les deux noms :

```go
`json:"site_name"`
```

La correspondance devient :

```text
JSON                    Go

"site_name"    →        SiteName
"Club Manager" →        valeur du champ
```

Après le décodage, nous pourrons obtenir :

```go
cfg.SiteName
```

avec la valeur :

```text
Club Manager
```

---

# Pourquoi les noms sont-ils différents ?

En Go, les champs exportés commencent par une majuscule :

```go
SiteName
```

Dans un fichier JSON, une convention fréquente consiste à écrire les noms en minuscules avec des underscores :

```json
"site_name"
```

Cette convention est appelée :

```text
snake_case
```

Exemples :

```text
site_name
contact_email
primary_color
registration_open
```

Go utilise généralement le style :

```text
PascalCase
```

Exemples :

```text
SiteName
ContactEmail
PrimaryColor
RegistrationOpen
```

L’annotation JSON permet donc de faire correspondre les deux conventions.

```go
SiteName string `json:"site_name"`
```

---

# Les règles importantes de syntaxe

## Utiliser des guillemets doubles

Correct :

```json
{
  "site_name": "Club Manager"
}
```

Incorrect :

```json
{
  'site_name': 'Club Manager'
}
```

JSON utilise des guillemets doubles, et non des apostrophes.

---

## Séparer la clé et la valeur avec `:`

```json
"site_name": "Club Manager"
```

---

## Séparer les propriétés avec une virgule

```json
{
  "site_name": "Club Manager",
  "logo": "logo.png"
}
```

---

## Ne pas ajouter de virgule finale

Incorrect :

```json
{
  "site_name": "Club Manager",
}
```

---

## Respecter les accolades et les crochets

Objet :

```json
{
}
```

Tableau :

```json
[
]
```

---

## JSON ne permet normalement pas les commentaires

Ceci n’est pas du JSON valide :

```json
{
  // Nom affiché dans l'application
  "site_name": "Club Manager"
}
```

Le format JSON standard ne possède pas de syntaxe pour écrire des commentaires.

Les explications doivent donc être placées dans une documentation séparée.

---

# JSON n’est pas un langage de programmation

JSON permet de décrire des données.

Il ne permet pas d’écrire :

- des fonctions ;
    
- des conditions ;
    
- des boucles ;
    
- des calculs ;
    
- des méthodes.
    

Un fichier JSON ne contient pas de comportement.

```text
Go
└── comportement et logique

JSON
└── données
```

Par exemple :

```json
{
  "site_name": "Club Manager"
}
```

décrit une valeur, mais ne précise pas comment l’application doit l’utiliser.

C’est le programme Go qui décide :

- quand lire le fichier ;
    
- comment vérifier les données ;
    
- où afficher le nom ;
    
- quoi faire en cas d’erreur.
    

---

# JSON comme fichier de configuration

Dans Club Manager, la progression pourrait être :

```text
config/config.json
        ↓
config.Load()
        ↓
Config
        ↓
main.go
        ↓
router
        ↓
handlers
        ↓
views
        ↓
templates
```

Le fichier contient les paramètres :

```json
{
  "site_name": "Club Manager"
}
```

Le package `config` lit le fichier.

Le programme transforme ensuite le JSON en structure Go :

```go
type Config struct {
	SiteName string `json:"site_name"`
}
```

La valeur peut enfin être utilisée dans l’application :

```go
cfg.SiteName
```

---

# Configuration générale et contenu des pages

Toutes les données n’ont pas besoin d’être placées dans le fichier JSON.

## Configuration générale

```text
site_name
logo
contact_email
couleurs du site
```

Ces valeurs peuvent être partagées par plusieurs pages.

## Contenu propre à une page

```text
titre de la page
en-tête
description
texte de présentation
```

Ces valeurs peuvent rester dans le code tant qu’elles ne doivent pas être personnalisables.

Exemple :

```go
data := views.HomeData{
	SiteName:    cfg.SiteName,
	Title:       "Accueil",
	Heading:     "Bienvenue sur Club Manager",
	Description: "Une application destinée à faciliter la gestion d'une association.",
}
```

La séparation devient :

```text
cfg.SiteName
└── configuration extérieure

Title, Heading, Description
└── contenu propre à la page
```

---

# Avantages du JSON

JSON présente plusieurs avantages :

- il est lisible dans un éditeur de texte ;
    
- il est pris en charge par de nombreux langages ;
    
- il représente des données structurées ;
    
- il permet de séparer certaines données du code ;
    
- il peut être analysé avec la bibliothèque standard de Go ;
    
- il est très répandu dans les outils web.
    

---

# Limites du JSON

JSON possède également quelques limites :

- il n’accepte pas les commentaires ;
    
- une virgule oubliée peut rendre le fichier invalide ;
    
- les guillemets doivent être correctement placés ;
    
- un utilisateur peut casser le fichier en modifiant sa structure ;
    
- il faut relancer l’application si le fichier est lu uniquement au démarrage ;
    
- il devient moins pratique pour des utilisateurs totalement non techniques.
    

À long terme, Club Manager pourra proposer une page d’administration.

```text
Formulaire d’administration
          ↓
Base de données
          ↓
Paramètres personnalisés
```

Le fichier JSON reste néanmoins utile pour découvrir le principe de configuration externe et pour conserver certains paramètres techniques.

---

# Comprendre et retenir

JSON est un format texte permettant de représenter des données structurées.

```json
{
  "site_name": "Club Manager"
}
```

Cette structure est appelée un objet JSON.

Elle associe une clé à une valeur :

```text
site_name → Club Manager
```

Dans d’autres langages, une structure similaire peut être appelée dictionnaire, map ou objet.

Les principaux symboles sont :

```text
{ }  → objet
[ ]  → tableau
:    → sépare une clé et une valeur
,    → sépare plusieurs propriétés
" "  → délimite une chaîne de caractères
```

JSON décrit des données, mais ne contient pas de logique.

```text
JSON → données
Go   → comportement
```

Dans Club Manager, le JSON permettra de déplacer certains paramètres généraux hors du code :

```text
config.json
    ↓
Config
    ↓
application
```

Il faut cependant conserver une distinction entre :

```text
configuration générale
et
contenu propre aux pages
```

