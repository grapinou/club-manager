
---

## project: Club Manager  
type: milestone  
milestone: 03  
git-tag: milestone-03-shared-layout  
git-commit: a56409d  
date: 2026-08-04  
status: frozen  
previous: milestone-02-home-template

# Jalon 03 - Layout commun

> [!info] État figé  
> Cette fiche représente Club Manager au tag Git :
> 
> `milestone-03-shared-layout`

---

## Objectif atteint

La première architecture basée sur les templates fonctionne.

Une nouvelle duplication est cependant apparue : chaque page pourrait avoir à répéter toute la structure HTML du site.

Club Manager introduit donc un **layout commun**.

---

# Nouvelle organisation des templates

Les templates sont maintenant séparés entre :

```text
templates/
├── layouts/
│   └── base.html
│
└── pages/
    ├── 001_home.html
    └── 002_club.html
```

On distingue ainsi :

```text
structure commune
       +
contenu spécifique
```

---

# `base.html`

Le layout définit la structure générale du site :

```html
{{ define "base" }}

<html>

<head>
    ...
    <title>{{ .Title }}</title>
</head>

<body>

    <header>
        ...
    </header>

    <main>
        {{ template "content" . }}
    </main>

    <footer>
        ...
    </footer>

</body>

</html>

{{ end }}
```

Il contient ce qui doit être partagé entre plusieurs pages :

- structure HTML ;
    
- `<head>` ;
    
- header ;
    
- zone principale ;
    
- footer.
    

---

# Les pages définissent leur contenu

Une page ne contient plus tout le document HTML.

Par exemple :

```html
{{ define "content" }}

<h1>{{ .Heading }}</h1>

<p>
    {{ .Description }}
</p>

{{ end }}
```

Le layout appelle ensuite :

```go
{{ template "content" . }}
```

On obtient donc :

```text
base.html
    │
    ├── header
    │
    ├── content ← page spécifique
    │
    └── footer
```

---

# Composition des templates

La vue charge maintenant plusieurs fichiers :

```go
template.ParseFS(
    homeFiles,
    "templates/layouts/base.html",
    "templates/pages/001_home.html",
)
```

Puis elle exécute :

```go
ExecuteTemplate(
    w,
    "base",
    data,
)
```

Le point d'entrée n'est donc plus le nom du fichier de page.

C'est le template défini par :

```go
{{ define "base" }}
```

---

# La page Club adopte l'architecture

Le principe n'est plus réservé à la page d'accueil.

`ClubHandler` construit désormais :

```go
data := views.ClubData{
    Title:       "Le club - Club Manager",
    Heading:     "Présentation du club",
    Description: "Découvrez l'association, son histoire et ses valeurs.",
}
```

puis appelle :

```go
views.RenderClub(w, data)
```

Une `ClubData`, une `RenderClub` et un template `002_club.html` apparaissent.

Le modèle commence donc à devenir reproductible.

---

# Architecture obtenue

```text
Navigateur
    ↓
 Router
    ↓
 Handler
    ↓
  Data
    ↓
  View
    ↓
┌───────────────────┐
│    base.html      │
│                   │
│ header            │
│                   │
│ ┌───────────────┐ │
│ │ page content  │ │
│ └───────────────┘ │
│                   │
│ footer            │
└───────────────────┘
    ↓
 navigateur
```

---

# Une architecture commence à émerger

Pour une page, on retrouve désormais régulièrement :

```text
route
 ↓
handler
 ↓
PageData
 ↓
RenderPage
 ↓
base + template de page
```

La page d'accueil et la page Club suivent maintenant ce modèle.

---

# Ce que le layout résout

Sans layout :

```text
home.html → structure HTML complète
club.html → structure HTML complète
```

Avec layout :

```text
base.html
├── structure commune
│
├── home content
└── club content
```

Cela réduit la duplication et permet de modifier la structure générale du site depuis un seul endroit.

---

# Limite encore présente

Le contenu reste écrit directement dans les handlers :

```go
Title: ...
Heading: ...
Description: ...
```

Club Manager sait donc maintenant **présenter proprement les données**, mais ces données sont encore étroitement liées au code.

Cette limite préparera une évolution future :

```text
contenu codé en dur
       ↓
configuration externe
```

---

# Comprendre et retenir

> **Un layout représente la structure commune des pages.**

---

> **Un template de page ne décrit que son contenu spécifique.**

---

> **La composition évite de répéter le même HTML dans chaque page.**

```text
layout commun
     +
contenu spécifique
     =
page complète
```

---

> **L'architecture handler → data → view → template commence maintenant à devenir un modèle réutilisable.**

---

# État du jalon

```text
Serveur HTTP               ✅
Router                     ✅
Handlers séparés           ✅
Tests                      ✅

Views                      ✅ Home / Club
Structures de données      ✅
Templates                  ✅
Layout commun              ✅
Composition des templates  ✅

Contact avec view          ❌
Rules avec view            ❌
Configuration              ❌
Base de données            ❌
```

**Ce jalon marque le passage de templates indépendants à une véritable structure de pages partagée et composable.**