
---

# Bootstrap dans Club Manager

## Objectif

Utiliser Bootstrap pour obtenir rapidement une interface :

- propre ;
    
- lisible ;
    
- responsive ;
    
- standard ;
    
- facile à maintenir.
    

L'objectif n'est pas de faire un design complexe.

Dans Club Manager, Bootstrap sert surtout à :

- structurer la navigation ;
    
- aligner les contenus ;
    
- gérer les espacements ;
    
- rendre le site utilisable sur mobile ;
    
- éviter d'écrire beaucoup de CSS personnalisé.
    

---

# 1. Pourquoi Bootstrap ?

Sans Bootstrap, il faudrait écrire nous-mêmes une grande partie du CSS nécessaire pour :

- les marges ;
    
- les espacements ;
    
- les largeurs ;
    
- la navigation responsive ;
    
- les boutons ;
    
- les grilles ;
    
- la typographie.
    

Bootstrap fournit déjà ces éléments.

La philosophie retenue pour Club Manager est donc :

```text
Bootstrap
    ↓
mise en forme générale

style.css
    ↓
petites personnalisations spécifiques
```

L'objectif reste :

> consacrer l'essentiel du travail à Go et au backend.

---

# 2. Charger Bootstrap

Bootstrap est chargé dans le template commun :

```text
internal/views/templates/layouts/base.html
```

Dans le `<head>` :

```html
<link
    href="https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/css/bootstrap.min.css"
    rel="stylesheet"
>
```

Puis notre propre CSS est chargé ensuite :

```html
<link
    rel="stylesheet"
    href="/static/css/style.css"
>
```

L'ordre est important :

```text
Bootstrap
    ↓
style.css
```

Ainsi, notre CSS peut éventuellement modifier certains styles Bootstrap.

---

# 3. Le JavaScript Bootstrap

Avant la fermeture de `<body>` :

```html
<script
    src="https://cdn.jsdelivr.net/npm/bootstrap@5.3.8/dist/js/bootstrap.bundle.min.js"
></script>
```

Ce JavaScript n'est pas nécessaire pour toutes les fonctionnalités de Bootstrap.

Il devient cependant nécessaire pour certains composants interactifs.

Par exemple :

```text
navbar mobile
collapse
menu hamburger
```

---

# 4. La navbar

La navigation initiale ressemblait à :

```html
<nav>
    <a href="/">Accueil</a>
    <a href="/club">Le club</a>
    <a href="/where">Où ?</a>
    <a href="/when">Quand ?</a>
    <a href="/contact">Contact</a>
    <a href="/rules">Règlement</a>
</nav>
```

Elle fonctionne, mais ne possède aucune mise en forme ni comportement responsive.

Bootstrap permet de la transformer en véritable navbar.

---

# 5. Structure de la navbar Bootstrap

```html
<nav class="navbar navbar-expand-lg bg-body-tertiary">

    <div class="container">

        <a
            class="navbar-brand"
            href="/"
        >
            {{ .SiteName }}
        </a>

        <button
            class="navbar-toggler"
            type="button"
            data-bs-toggle="collapse"
            data-bs-target="#mainNavbar"
            aria-controls="mainNavbar"
            aria-expanded="false"
            aria-label="Afficher la navigation"
        >
            <span class="navbar-toggler-icon"></span>
        </button>

        <div
            class="collapse navbar-collapse"
            id="mainNavbar"
        >

            <div class="navbar-nav">

                <a class="nav-link" href="/">
                    Accueil
                </a>

                <a class="nav-link" href="/club">
                    Le club
                </a>

                <a class="nav-link" href="/where">
                    Où ?
                </a>

                <a class="nav-link" href="/when">
                    Quand ?
                </a>

                <a class="nav-link" href="/contact">
                    Contact
                </a>

                <a class="nav-link" href="/rules">
                    Règlement
                </a>

            </div>

        </div>

    </div>

</nav>
```

---

# 6. `navbar`

La classe :

```html
navbar
```

indique à Bootstrap :

> cet élément est une barre de navigation.

Elle fournit notamment :

- une structure adaptée ;
    
- des espacements ;
    
- des comportements propres aux composants Bootstrap.
    

---

# 7. `navbar-expand-lg`

```html
navbar-expand-lg
```

permet à la navbar d'être responsive.

Sur un écran suffisamment large :

```text
Club Manager   Accueil   Club   Où ?   Quand ?   Contact
```

Sur un petit écran :

```text
Club Manager                          ☰
```

Les liens peuvent alors être affichés ou masqués grâce au bouton hamburger.

---

# 8. Le bouton hamburger

Le bouton possède :

```html
data-bs-toggle="collapse"
```

et :

```html
data-bs-target="#mainNavbar"
```

Il contrôle l'élément :

```html
<div
    class="collapse navbar-collapse"
    id="mainNavbar"
>
```

Le lien entre les deux est donc :

```text
data-bs-target="#mainNavbar"
               │
               ▼
        id="mainNavbar"
```

Lorsque l'utilisateur clique :

```text
bouton hamburger
        ↓
Bootstrap JavaScript
        ↓
cherche #mainNavbar
        ↓
affiche ou masque la navigation
```

---

# 9. `navbar-brand`

Le nom du site utilise :

```html
<a
    class="navbar-brand"
    href="/"
>
    {{ .SiteName }}
</a>
```

`navbar-brand` représente généralement :

- le nom du site ;
    
- le logo ;
    
- ou les deux.
    

Dans Club Manager, `SiteName` vient déjà des données transmises par Go.

```text
Config
   ↓
PageData
   ↓
{{ .SiteName }}
   ↓
navbar-brand
```

---

# 10. Ajouter un logo

Si le logo est stocké dans :

```text
static/images/logo.png
```

il peut être utilisé dans la navbar :

```html
<a
    class="navbar-brand d-flex align-items-center gap-2"
    href="/"
>

    <img
        src="/static/images/logo.png"
        alt=""
        width="32"
        height="32"
    >

    {{ .SiteName }}

</a>
```

Les classes utilisées sont :

```text
d-flex
```

pour placer les éléments avec Flexbox,

```text
align-items-center
```

pour les centrer verticalement,

et :

```text
gap-2
```

pour ajouter un espace entre le logo et le texte.

---

# 11. Le `container`

Bootstrap fournit la classe :

```html
container
```

Elle sert principalement à :

- limiter la largeur du contenu ;
    
- créer des marges adaptées à la taille de l'écran ;
    
- aligner les différentes parties du site.
    

Exemple :

```html
<div class="container">
```

---

# 12. Aligner navbar, contenu et footer

On utilise le même `container` dans les trois parties principales.

## Navbar

```html
<nav class="navbar navbar-expand-lg bg-body-tertiary">
    <div class="container">
        ...
    </div>
</nav>
```

## Contenu

```html
<main class="container py-4">
    {{ template "content" . }}
</main>
```

## Footer

```html
<footer class="border-top py-3">
    <div class="container">
        ...
    </div>
</footer>
```

Ainsi :

```text
Navbar  ─┐
Main     ├── même alignement horizontal
Footer  ─┘
```

Visuellement :

```text
┌──────────────────────────────────────────────┐
│                                             │
│     Club Manager   Accueil   Club ...       │
│     │                                       │
│     │ Titre de la page                      │
│     │ Description                           │
│     │                                       │
│     │ Club Manager                          │
│                                             │
└──────────────────────────────────────────────┘
      ↑
      même limite horizontale
```

---

# 13. `container` et `container-fluid`

Il existe notamment :

```html
container
```

et :

```html
container-fluid
```

## `container`

La largeur maximale évolue selon la taille de l'écran.

Le contenu reste centré.

---

## `container-fluid`

Occupe pratiquement toute la largeur disponible.

Dans Club Manager, nous préférons :

```html
container
```

afin d'avoir un alignement commun entre :

- navbar ;
    
- contenu ;
    
- footer.
    

---

# 14. Les classes d'espacement

Bootstrap fournit des classes permettant d'éviter d'écrire du CSS pour les espacements courants.

Exemple :

```html
py-4
```

Décomposition :

```text
p = padding

y = axe vertical

4 = niveau d'espacement
```

Donc :

```text
py-4
```

ajoute :

```text
padding-top
+
padding-bottom
```

---

# 15. Pourquoi `py-4` sur le `<main>` ?

Nous utilisons :

```html
<main class="container py-4">
```

Le `container` gère l'alignement horizontal.

Le `py-4` ajoute de l'espace vertical.

Cela évite que le contenu soit collé à la navbar ou au footer.

```text
Navbar
────────────────────────

      espace

Titre

Description

      espace

────────────────────────
Footer
```

---

# 16. Pourquoi ne pas mettre `py-4` sur le container de la navbar ?

La navbar possède déjà ses propres espacements Bootstrap.

Ajouter :

```html
<div class="container py-4">
```

augmenterait fortement sa hauteur.

Le rôle du `container` de la navbar est principalement :

> aligner son contenu horizontalement avec le reste du site.

---

# 17. Le footer

Le footer peut rester extrêmement simple :

```html
<footer class="border-top py-3">

    <div class="container">

        <p class="mb-0 text-body-secondary">
            {{ .SiteName }}
        </p>

    </div>

</footer>
```

---

# 18. Classes utilisées dans le footer

## `border-top`

```html
border-top
```

ajoute une séparation légère au-dessus du footer.

---

## `py-3`

```html
py-3
```

ajoute un espace vertical.

---

## `mb-0`

```html
mb-0
```

signifie :

```text
margin-bottom: 0
```

Cela retire la marge basse habituelle du paragraphe.

---

## `text-body-secondary`

```html
text-body-secondary
```

rend le texte plus discret que le texte principal.

Cela convient bien à un footer.

---

# 19. Structure actuelle du layout

Le layout se rapproche maintenant de :

```html
<body>

    <nav class="navbar navbar-expand-lg bg-body-tertiary">

        <div class="container">

            ...

        </div>

    </nav>

    <main class="container py-4">

        {{ template "content" . }}

    </main>

    <footer class="border-top py-3">

        <div class="container">

            <p class="mb-0 text-body-secondary">
                {{ .SiteName }}
            </p>

        </div>

    </footer>

</body>
```

On obtient donc trois grandes zones :

```text
┌───────────────────────────┐
│ Navbar                    │
├───────────────────────────┤
│                           │
│ Main                      │
│                           │
├───────────────────────────┤
│ Footer                    │
└───────────────────────────┘
```

---

# 20. Bootstrap ne remplace pas le HTML

Bootstrap repose toujours sur du HTML classique.

Par exemple :

```html
<nav>
```

reste une balise HTML.

Bootstrap ajoute simplement :

```html
class="navbar"
```

Le principe est donc :

```text
HTML
   ↓
structure

Bootstrap
   ↓
présentation et comportement
```

---

# 21. Bootstrap ne remplace pas notre CSS

Nous conservons :

```text
static/css/style.css
```

Mais son rôle doit rester limité.

Bootstrap gère :

```text
mise en page générale
responsive
espacements
composants standards
```

Notre CSS pourra gérer :

```text
couleurs du club
petites adaptations
éléments réellement spécifiques
```

Cela évite de réécrire ce que Bootstrap sait déjà faire.

---

# 22. Principe retenu pour la V1

La V1 doit rester simple.

Il n'est donc pas nécessaire d'utiliser tous les composants Bootstrap disponibles.

Nous utilisons uniquement ceux qui répondent à un besoin réel.

Pour l'instant :

```text
Bootstrap
│
├── navbar
├── container
├── espacements
├── footer
└── responsive
```

Pas besoin pour la V1 de multiplier :

- carrousels ;
    
- animations ;
    
- modales ;
    
- composants complexes ;
    
- effets visuels ;
    
- JavaScript personnalisé.
    

---

# 23. Philosophie générale

La logique reste la même que pour l'architecture Go :

```text
Faire simple
    ↓
observer les besoins
    ↓
ajouter seulement ce qui est utile
```

Avec Bootstrap :

```text
HTML simple
    +
quelques classes Bootstrap
    +
très peu de CSS personnalisé
```

est suffisant pour obtenir rapidement une interface propre.

---

# Comprendre et retenir

## À quoi sert Bootstrap ?

Bootstrap fournit des composants et des classes CSS permettant de construire rapidement une interface standard et responsive.

---

## Que fait `container` ?

Il limite et centre la largeur du contenu.

Il permet aussi d'aligner :

```text
navbar
main
footer
```

---

## Que fait `py-4` ?

Il ajoute du padding vertical :

```text
haut
+
bas
```

---

## Que fait `navbar-expand-lg` ?

Il permet à la navbar d'être développée sur grand écran et repliable sur les écrans plus petits.

---

## Pourquoi charger le JavaScript Bootstrap ?

Parce que certains composants interactifs en ont besoin.

Dans notre cas :

```text
bouton hamburger
    ↓
collapse de la navbar
```

---

## Pourquoi conserver `style.css` ?

Pour les personnalisations propres à Club Manager.

Mais Bootstrap doit prendre en charge la majorité du style générique.

---

# À retenir

Bootstrap n'est pas le cœur du projet.

Il est un outil permettant de construire rapidement une interface correcte pendant que l'essentiel du travail reste concentré sur :

```text
Go
architecture
HTTP
templates
configuration
puis base de données
```

Pour la V1 :

> quelques composants Bootstrap bien choisis valent mieux qu'une interface complexe.

