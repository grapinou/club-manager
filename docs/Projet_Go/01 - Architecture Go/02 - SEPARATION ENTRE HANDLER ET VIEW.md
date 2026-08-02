

---

# Séparation entre handler et view

## Contexte

Dans la première architecture de Club Manager, chaque handler écrit directement le contenu de la réponse HTTP.

Exemple :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Bienvenue sur Club Manager")
}
```

Le handler remplit alors deux responsabilités :

1. recevoir et traiter la requête HTTP ;
    
2. produire le contenu affiché à l'utilisateur.
    

Cette solution est suffisante pour une première architecture minimale.

Cependant, elle devient rapidement limitée lorsque la page contient un véritable document HTML.

---

## Nouvelle séparation des responsabilités

Nous introduisons maintenant une nouvelle couche :

```text
Routeur
   ↓
Handler
   ↓
View
   ↓
Template HTML
```

Chaque élément possède une responsabilité différente.

### Le routeur

Le routeur associe une méthode HTTP et une URL à un handler.

```go
mux.HandleFunc("GET /{$}", handlers.HomeHandler)
```

Il répond à la question :

> Quelle fonction doit recevoir cette requête ?

Il ne construit pas la page HTML.

---

### Le handler

Le handler reçoit la requête HTTP et coordonne la réponse.

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := views.RenderHome(w)

	if err != nil {
		http.Error(
			w,
			"Erreur interne du serveur",
			http.StatusInternalServerError,
		)
	}
}
```

Il répond notamment aux questions suivantes :

- quelle vue faut-il appeler ?
    
- quelles données faut-il lui transmettre ?
    
- que faut-il faire si le rendu échoue ?
    
- quel statut HTTP faut-il retourner ?
    

Le handler conserve donc la responsabilité liée au protocole HTTP.

---

### La view

La view produit le contenu destiné à être envoyé au navigateur.

```go
func RenderHome(w io.Writer) error {
	return homeTemplate.ExecuteTemplate(
		w,
		"home.html",
		nil,
	)
}
```

Elle répond à la question :

> Comment transformer un template et des données en contenu HTML ?

La view ne choisit pas la route et ne démarre pas le serveur.

---

### Le template

Le template décrit la structure de la page.

```html
<!DOCTYPE html>
<html lang="fr">
<head>
    <meta charset="UTF-8">
    <title>Club Manager</title>
</head>

<body>
    <h1>Bienvenue sur Club Manager</h1>
</body>
</html>
```

Il contient principalement :

- la structure HTML ;
    
- le texte affiché ;
    
- les emplacements réservés aux futures données dynamiques ;
    
- plus tard, les appels aux éléments communs comme la navigation.
    

Le package `html/template` est prévu pour produire du HTML et applique un échappement adapté aux données insérées dans les pages.

---

## Parcours d'une requête

Lorsqu'un navigateur demande la page d'accueil :

```text
GET /
```

le parcours devient :

```text
1. Le navigateur envoie GET /
                 ↓
2. Le routeur trouve HomeHandler
                 ↓
3. HomeHandler appelle RenderHome
                 ↓
4. RenderHome exécute home.html
                 ↓
5. Le HTML est écrit dans la réponse
                 ↓
6. Le navigateur reçoit et affiche la page
```

Dans le code :

```text
"GET /{$}"
      ↓
HomeHandler
      ↓
views.RenderHome
      ↓
homeTemplate.ExecuteTemplate
      ↓
home.html
```

---

## Qui associe la route au fichier HTML ?

Go ne déduit pas automatiquement que :

```text
GET /contact
```

doit utiliser :

```text
contact.html
```

L'association est écrite explicitement dans notre code.

Première association :

```go
mux.HandleFunc("GET /contact", handlers.ContactHandler)
```

Elle relie la route au handler.

Deuxième association :

```go
func ContactHandler(w http.ResponseWriter, r *http.Request) {
	err := views.RenderContact(w)
	// Gestion de l'erreur...
}
```

Elle relie le handler à la fonction de rendu.

Troisième association :

```go
func RenderContact(w io.Writer) error {
	return pages.ExecuteTemplate(
		w,
		"contact.html",
		nil,
	)
}
```

Elle relie la fonction de rendu au template HTML.

Le parcours complet est donc explicite :

```text
GET /contact
      ↓
ContactHandler
      ↓
RenderContact
      ↓
contact.html
```

Il n'existe pas de magie basée sur le nom des fichiers.

C'est notre code qui construit chaque association.

---

## Pourquoi cette séparation est-elle utile ?

### Le handler reste lisible

Sans séparation, un handler pourrait finir par contenir une longue chaîne HTML :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "<!DOCTYPE html>")
	fmt.Fprintln(w, "<html>")
	fmt.Fprintln(w, "<body>")
	fmt.Fprintln(w, "<h1>Bienvenue</h1>")
	// ...
}
```

Le code Go et la présentation seraient mélangés.

Avec une view :

```go
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	err := views.RenderHome(w)
	// ...
}
```

Le handler exprime seulement son intention :

```text
Rendre la page d'accueil
```

---

### Le HTML reste dans un fichier HTML

Le template peut être lu et modifié comme un document HTML classique.

```text
internal/views/templates/home.html
```

Cela facilite :

- la lecture de la structure de la page ;
    
- l'ajout futur de CSS ;
    
- l'utilisation de Bootstrap ;
    
- la création d'une navigation ;
    
- le partage d'éléments communs entre plusieurs pages.
    

---

### Les responsabilités deviennent testables séparément

Nous pouvons distinguer plusieurs comportements :

```text
Le routeur dirige-t-il / vers HomeHandler ?
Le handler répond-il correctement ?
La view produit-elle le HTML attendu ?
Le template contient-il les informations nécessaires ?
```

Nous ne serons pas obligés de tout tester de la même manière.

---

## Est-ce déjà une architecture MVC ?

Cette organisation commence à séparer la présentation du traitement HTTP.

```text
Handler → coordination HTTP
View    → présentation
```

Elle se rapproche de certains principes utilisés dans une architecture MVC, mais Club Manager ne possède pas encore toutes les couches d'une telle architecture.

Nous n'avons notamment pas encore :

- de véritable modèle métier ;
    
- de base de données ;
    
- de repositories ;
    
- de services applicatifs ;
    
- de données complexes transmises aux vues.
    

Il est donc plus précis de parler pour l'instant de :

> séparation entre la gestion HTTP et le rendu de la présentation.

Nous pourrons comparer cette architecture à MVC lorsque les modèles et la logique métier apparaîtront.

---

## Architecture obtenue

```text
club-manager/
├── cmd/
│   └── server/
│       └── main.go
│
└── internal/
    ├── handlers/
    │   └── 001_home.go
    │
    ├── router/
    │   └── 001_router.go
    │
    └── views/
        ├── 001_home.go
        │
        └── templates/
            └── home.html
```

Les responsabilités deviennent :

```text
main.go
└── lance l'application

router
└── associe les routes aux handlers

handlers
└── gère les requêtes et les réponses HTTP

views
└── prépare et exécute les templates

templates
└── décrit les pages HTML
```

---

## Comprendre et retenir

Le routeur choisit le handler :

```text
route → handler
```

Le handler choisit la view :

```text
handler → fonction de rendu
```

La fonction de rendu choisit le template :

```text
fonction de rendu → fichier HTML
```

Go ne choisit pas automatiquement `home.html` ou `contact.html` selon l'URL.

L'association est explicitement décrite par notre code :

```text
GET /contact
      ↓
ContactHandler
      ↓
RenderContact
      ↓
contact.html
```

Cette séparation évite de mélanger :

```text
gestion HTTP
et
construction du HTML
```

