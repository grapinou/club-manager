
---

## project: Club Manager  
type: milestone  
milestone: 02  
git-tag: milestone-02-home-template  
git-commit: a36a014  
date: 2026-08-02  
status: frozen  
previous: milestone-01-handlers-routes

# Jalon 02 - Premier template

> [!info] État figé  
> Cette fiche représente Club Manager au tag Git :
> 
> `milestone-02-home-template`

---

## Objectif atteint

La page d'accueil n'est plus générée directement par son handler.

Club Manager possède maintenant une première séparation entre :

```text
traitement HTTP
```

et :

```text
présentation HTML
```

La page d'accueil sert d'expérimentation avant de généraliser cette architecture aux autres pages.

---

# Nouvelle architecture de l'accueil

```text
Navigateur
    ↓
 Router
    ↓
HomeHandler
    ↓
 HomeData
    ↓
RenderHome
    ↓
 Template
    ↓
   HTML
```

Une nouvelle couche apparaît donc :

```text
internal/views/
```

---

# Le handler change de rôle

Avant :

```go
fmt.Fprintln(w, "Bienvenue sur Club Manager")
```

Maintenant, le handler prépare les données :

```go
data := views.HomeData{
    Title:       "Club Manager",
    Heading:     "Bienvenue sur Club Manager",
    Description: "Une application destinée à faciliter la gestion d'une association.",
}
```

puis demande à la vue de produire la réponse :

```go
err := views.RenderHome(w, data)
```

Le handler ne contient donc plus directement le HTML.

---

# `HomeData`

Une structure intermédiaire apparaît :

```go
type HomeData struct {
    Title       string
    Heading     string
    Description string
}
```

Elle représente les données nécessaires à l'affichage de la page.

Le trajet devient :

```text
Handler
   ↓
HomeData
   ↓
Template
```

Cela constitue un premier exemple concret de passage de données entre différentes couches de l'application.

---

# La vue

La vue utilise plusieurs outils de Go :

```go
embed
html/template
io
```

Le template est intégré au programme avec :

```go
//go:embed templates/001_home.html
```

puis chargé avec :

```go
template.ParseFS(...)
```

et vérifié grâce à :

```go
template.Must(...)
```

La fonction :

```go
RenderHome(w io.Writer, data HomeData)
```

exécute finalement le template.

---

# Le template devient dynamique

Le HTML n'est plus entièrement écrit en dur.

Il utilise :

```html
<title>{{ .Title }}</title>

<h1>{{ .Heading }}</h1>

<p>
    {{ .Description }}
</p>
```

Les données provenant de Go sont injectées dans la page HTML.

---

# Une nouvelle séparation des responsabilités

Au jalon précédent :

```text
Handler
   ↓
texte
```

Maintenant :

```text
Handler
   │
   ├── prépare les données
   ↓
 View
   │
   ├── gère le rendu
   ↓
Template
   │
   └── décrit le HTML
```

C'est une évolution importante.

---

# Pourquoi seulement l'accueil ?

À ce stade, cette architecture n'est appliquée qu'à la page d'accueil.

Les autres handlers peuvent encore produire directement leur réponse.

C'est volontairement un état intermédiaire.

La page d'accueil joue le rôle de prototype permettant de comprendre :

- `embed`;
    
- `html/template`;
    
- `ParseFS`;
    
- `template.Must`;
    
- `io.Writer`;
    
- le passage de données ;
    
- la séparation handler / view.
    

Une fois ce modèle compris, il pourra être généralisé.

---

# Limite observée

Le template de la page d'accueil contient encore l'ensemble du document :

```text
<html>
├── head
└── body
    └── contenu de l'accueil
```

Si chaque page reproduit cette structure, on obtiendra rapidement :

```text
home.html
├── head
├── header
├── contenu
└── footer

club.html
├── head
├── header
├── contenu
└── footer
```

Une nouvelle duplication apparaît.

Elle prépare naturellement le jalon suivant :

> créer un layout partagé.

---

# Comprendre et retenir

> **Le handler prépare les données ; la vue produit la représentation.**

---

> **Une struct comme `HomeData` sert de contrat entre le code Go et le template.**

```text
Go
 ↓
HomeData
 ↓
HTML
```

---

> **Le template ne décide pas des données. Il décrit leur présentation.**

---

# État du jalon

```text
Serveur HTTP             ✅
Router                   ✅
Handlers séparés         ✅
Tests                    ✅

View pour Home           ✅
HomeData                 ✅
Template HTML            ✅
Passage de données       ✅
embed.FS                 ✅

Views pour toutes pages  ❌
Layout commun            ❌
Configuration            ❌
Base de données          ❌
```

**Ce jalon marque la première séparation entre le traitement HTTP et la génération de l'interface HTML.**