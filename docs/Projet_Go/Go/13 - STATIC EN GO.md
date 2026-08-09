
---

# Fichiers statiques avec Go

## Objectif

Permettre au serveur Go de fournir au navigateur des fichiers qui ne sont pas générés dynamiquement.

Exemples :

- fichiers CSS ;
    
- images ;
    
- icônes ;
    
- éventuellement fichiers JavaScript ;
    
- polices.
    

Dans Club Manager, ces fichiers sont placés dans :

```text
static/
```

---

# 1. Arborescence

Exemple :

```text
club-manager/
│
├── cmd/
├── config/
├── internal/
│
├── static/
│   ├── css/
│   │   └── style.css
│   │
│   └── images/
│
├── go.mod
└── README.md
```

Le fichier CSS se trouve donc physiquement ici :

```text
static/css/style.css
```

---

# 2. Demander le CSS depuis le HTML

Dans le template HTML :

```html
<link rel="stylesheet" href="/static/css/style.css">
```

Cette ligne ne demande pas directement à Go :

```text
ouvre static/css/style.css
```

Elle dit au **navigateur** :

> Demande au serveur la ressource située à l'URL `/static/css/style.css`.

Le navigateur effectue donc une nouvelle requête HTTP :

```http
GET /static/css/style.css
```

---

# 3. URL et chemin sur le disque

Il faut distinguer deux choses.

## URL publique

```text
/static/css/style.css
```

C'est ce que connaît le navigateur.

---

## Chemin réel du fichier

```text
static/css/style.css
```

C'est ce que connaît le serveur.

Ces deux chemins n'ont pas l'obligation d'être identiques.

Le serveur fait le lien entre les deux.

---

# 4. Créer un serveur de fichiers

Dans le routeur :

```go
staticFiles := http.FileServer(http.Dir("static"))
```

Décomposons.

---

## `http.Dir("static")`

```go
http.Dir("static")
```

définit :

```text
static/
```

comme racine des fichiers accessibles par ce serveur de fichiers.

On peut se représenter cela ainsi :

```text
racine FileServer
      │
      ▼
   static/
      │
      ├── css/
      └── images/
```

---

## `http.FileServer(...)`

```go
http.FileServer(...)
```

crée un handler HTTP capable de retourner des fichiers.

Exemple :

```go
http.FileServer(http.Dir("static"))
```

signifie :

> Lorsque je reçois un chemin de fichier, cherche-le à partir du dossier `static`.

---

# 5. Le problème du préfixe `/static/`

Le navigateur demande :

```text
/static/css/style.css
```

Mais `FileServer` possède déjà comme racine :

```text
static/
```

On voudrait donc lui transmettre seulement :

```text
css/style.css
```

et non :

```text
static/css/style.css
```

dans la requête.

C'est le rôle de :

```go
http.StripPrefix()
```

---

# 6. `http.StripPrefix`

On écrit :

```go
http.StripPrefix("/static/", staticFiles)
```

Le chemin HTTP :

```text
/static/css/style.css
```

devient :

```text
css/style.css
```

Important :

> `StripPrefix` modifie le chemin de la requête HTTP.

Il ne supprime pas un dossier réel sur le disque.

---

# 7. Déclarer la route

On peut alors écrire :

```go
mux.Handle(
	"GET /static/",
	http.StripPrefix("/static/", staticFiles),
)
```

Cette route signifie :

> Toutes les requêtes GET commençant par `/static/` sont envoyées vers le serveur de fichiers statiques.

Exemples :

```text
GET /static/css/style.css
GET /static/images/logo.png
GET /static/images/background.jpg
```

---

# 8. Parcours complet

Prenons :

```html
<link rel="stylesheet" href="/static/css/style.css">
```

Le parcours complet est :

```text
HTML
 │
 ▼
href="/static/css/style.css"
 │
 ▼
navigateur
 │
 ▼
GET /static/css/style.css
 │
 ▼
ServeMux
 │
 ▼
route "GET /static/"
 │
 ▼
StripPrefix("/static/")
 │
 ▼
css/style.css
 │
 ▼
FileServer
 │
 │ racine = static/
 ▼
static/css/style.css
 │
 ▼
fichier renvoyé au navigateur
```

---

# 9. Résumé des responsabilités

## Le navigateur

Il connaît :

```text
/static/css/style.css
```

---

## `ServeMux`

Il reconnaît :

```text
GET /static/
```

---

## `StripPrefix`

Il transforme :

```text
/static/css/style.css
```

en :

```text
css/style.css
```

---

## `FileServer`

Il cherche ce chemin dans :

```text
static/
```

Il obtient donc :

```text
static/css/style.css
```

---

# 10. Pourquoi ne pas écrire simplement `css/style.css` ?

On pourrait écrire :

```html
<link rel="stylesheet" href="css/style.css">
```

Mais il faut comprendre que ceci est une **URL relative**.

Le navigateur ne connaît pas l'organisation du projet Go.

Il interprète seulement une URL.

Par exemple, depuis :

```text
http://localhost:8080/
```

cela pourrait produire :

```text
http://localhost:8080/css/style.css
```

donc :

```http
GET /css/style.css
```

Il faudrait alors créer une route correspondant à :

```text
/css/
```

---

# 11. Pourquoi utiliser `/static/` ?

Utiliser :

```text
/static/
```

n'est pas obligatoire.

C'est surtout une convention permettant de distinguer clairement les ressources statiques.

Exemple :

```text
/                  page dynamique
/club              page dynamique
/contact           page dynamique
/rules             page dynamique

/static/css/...    fichier statique
/static/images/... fichier statique
```

On réserve ainsi une partie claire de l'espace des URL.

---

# 12. Chemin absolu côté site

Avec :

```html
href="/static/css/style.css"
```

le `/` initial signifie :

> repartir depuis la racine du site.

Ainsi, quelle que soit la page courante :

```text
/
/club
/contact
/rules
```

le navigateur demandera toujours :

```text
/static/css/style.css
```

---

# 13. Différence avec un chemin relatif

## Chemin relatif

```html
href="css/style.css"
```

Le chemin dépend de l'URL actuelle.

---

## Chemin depuis la racine

```html
href="/static/css/style.css"
```

Le navigateur repart toujours de :

```text
/
```

C'est généralement plus adapté pour les ressources communes au site.

---

# 14. Code complet minimal

Dans le routeur :

```go
func New(cfg config.Config) *http.ServeMux {

	mux := http.NewServeMux()

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

	mux.HandleFunc(
		"GET /{$}",
		handlers.HomeHandler(cfg),
	)

	return mux
}
```

Dans `base.html` :

```html
<head>

    <meta charset="UTF-8">

    <meta
        name="viewport"
        content="width=device-width, initial-scale=1.0"
    >

    <title>{{ .Title }}</title>

    <link
        rel="stylesheet"
        href="/static/css/style.css"
    >

</head>
```

---

# 15. Exemple de CSS de test

```css
body {
    font-family: sans-serif;
}
```

Ce fichier est enregistré dans :

```text
static/css/style.css
```

---

# 16. Tester directement le serveur statique

Avec le serveur lancé :

```bash
go run ./cmd/server
```

on peut accéder directement à :

```text
http://localhost:8080/static/css/style.css
```

Si le navigateur affiche le contenu du fichier CSS, cela signifie que :

```text
ServeMux
+
StripPrefix
+
FileServer
```

fonctionnent correctement.

---

# 17. Ce que `FileServer` évite

Sans `FileServer`, il faudrait écrire manuellement un handler capable de :

- ouvrir chaque fichier ;
    
- déterminer son type ;
    
- lire son contenu ;
    
- l'envoyer au navigateur ;
    
- gérer les erreurs.
    

`http.FileServer` fournit déjà cette mécanique.

---

# 18. À retenir sur `StripPrefix`

Sans transformation, la requête serait :

```text
/static/css/style.css
```

alors que notre racine de fichiers est déjà :

```text
static/
```

On veut donc transmettre à `FileServer` :

```text
css/style.css
```

Le rôle de :

```go
http.StripPrefix("/static/", ...)
```

est précisément de réaliser cette transformation.

---

# Comprendre et retenir

## Que fait le `<link>` HTML ?

```html
<link rel="stylesheet" href="/static/css/style.css">
```

Il provoque une requête HTTP du navigateur :

```http
GET /static/css/style.css
```

---

## Que fait `http.Dir("static")` ?

Il définit le dossier :

```text
static/
```

comme racine du système de fichiers servi.

---

## Que fait `http.FileServer` ?

Il fournit des fichiers au navigateur à partir d'un système de fichiers.

---

## Que fait `http.StripPrefix` ?

Il enlève un préfixe du **chemin HTTP** avant de transmettre la requête au handler suivant.

---

## Pourquoi utiliser `/static/` ?

Pour avoir une zone d'URL clairement réservée aux ressources statiques.

---

## Le navigateur connaît-il le dossier `static/` du projet ?

Non.

Le navigateur ne connaît que des URL.

C'est le serveur Go qui associe :

```text
URL
```

à :

```text
fichier sur le disque
```

---

# À retenir

```text
/static/css/style.css
          │
          ▼
GET /static/css/style.css
          │
          ▼
StripPrefix("/static/")
          │
          ▼
css/style.css
          │
          ▼
FileServer(http.Dir("static"))
          │
          ▼
static/css/style.css
```

La notion essentielle est donc :

> Une URL HTTP et un chemin sur le disque sont deux choses différentes.

Le serveur fait le lien entre les deux.

