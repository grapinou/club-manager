
# Club Manager — Index des fiches

Cet espace regroupe les fiches utilisées pendant le développement du projet **Club Manager**.

Les notes servent à :

- documenter la progression du projet ;
    
- comprendre les notions rencontrées ;
    
- conserver les commandes importantes ;
    
- expliquer les choix techniques ;
    
- construire progressivement une documentation de référence.
    

---

## Organisation des fiches

Les fichiers sont numérotés afin de conserver un ordre logique.

- `00` désigne généralement la fiche principale d’un thème ;
    
- `01`, `02`, `03`... représentent les étapes ou notions principales ;
    
- `01.01`, `01.02`... approfondissent une fiche principale.
    

Exemple :

```text
04 - INTERFACE
├── 04.01 - Interface composée
└── 04.02 - Interface utilisation
```

---

# Projet Club Manager

## Feuille de route

- [[Projet_Go/Feuille de route projet go|Feuille de route du projet Go]]
    

---

## 00 — Initialisation du projet

Fiches consacrées à la préparation de l’environnement, à Git, à GitHub et à l’initialisation du module Go.

- [[Projet_Go/00 - Initialisation projet Go/00 - INITIALISATION PROJET|Initialisation du projet]]
    
- [[Projet_Go/00 - Initialisation projet Go/01 - Installer la dernière version de Go|Installer la dernière version de Go]]
    
    - [[Projet_Go/00 - Initialisation projet Go/01.01 - wget|wget]]
        
    - [[Projet_Go/00 - Initialisation projet Go/01.02 - tar|tar]]
        
    - [[Projet_Go/00 - Initialisation projet Go/01.03 - absolu|Chemins absolus et relatifs]]
        
- [[Projet_Go/00 - Initialisation projet Go/02 - Vérifier si Git est bien installé|Vérifier si Git est installé]]
    
- [[Projet_Go/00 - Initialisation projet Go/03 - git init|Initialiser un dépôt Git]]
    
- [[Projet_Go/00 - Initialisation projet Go/04 - Ajouter .gitignore|Ajouter un fichier .gitignore]]
    
- [[Projet_Go/00 - Initialisation projet Go/05 - Premier commit|Effectuer le premier commit]]
    
- [[Projet_Go/00 - Initialisation projet Go/06 - Créer le dépôt GitHub|Créer le dépôt GitHub]]
    
    - [[Projet_Go/00 - Initialisation projet Go/06.01 - clé ssh|Clé SSH]]
        
- [[Projet_Go/00 - Initialisation projet Go/07 - go mod init|Initialiser le module Go]]
    

---

## 01 — Architecture du projet Go

Fiches consacrées à l’organisation des dossiers et à la mise en place de l’architecture minimale de Club Manager.

- [[Projet_Go/01 - Architecture Go/00 - ARCHITECTURE MINIMALE EN GO|Architecture minimale en Go]]
    
    - [[Projet_Go/01 - Architecture Go/00.01 - mkdir -p|La commande mkdir -p]]
        

---

## Notions Go

### 00 — Serveur HTTP

- [[Projet_Go/Go/00 - SERVEUR GO|Serveur Go]]
    

### 01 — Fonctions

- [[Projet_Go/Go/01 - FONCTION|Fonctions]]
    
    - [[Projet_Go/Go/01.01 - Fonction exportation|Exportation d’une fonction]]
        
    - [[Projet_Go/Go/01.02 - Fonction New|Convention New]]
        

### 02 — Packages

- [[Projet_Go/Go/02 - PACKAGE|Packages]]
    
    - [[Projet_Go/Go/02.01 - Package import|Imports de packages]]
        

### 03 — Structs et méthodes

- [[Projet_Go/Go/03 - STRUCT ET METHODE|Structs et méthodes]]
    

### 04 — Interfaces

- [[Projet_Go/Go/04 - INTERFACE|Interfaces]]
    
    - [[Projet_Go/Go/04.01 - Interface composée|Interfaces composées]]
        
    - [[Projet_Go/Go/04.02 - Interface utilisation|Utilisation des interfaces]]
        

### 05 — Pointeurs

- [[Projet_Go/Go/05 - POINTEUR|Pointeurs]]
    

### 06 — Affichage et écriture

- [[Projet_Go/Go/06 - PRINT|Fonctions d’affichage]]
    
    - [[Projet_Go/Go/06.01 - fmt.Fprintln|fmt.Fprintln]]
        
    - [[Projet_Go/Go/06.02 - io.Writer|io.Writer]]
        

### 07 — Handlers HTTP

- [[Projet_Go/Go/07 - HANDLER|Handlers HTTP]]
    

### 08 — Routeur HTTP

- [[Projet_Go/Go/08 - ROUTEUR|Routeur]]
    
    - [[Projet_Go/Go/08.01 - Routeur ServeMux|ServeMux]]
        

---

# Git

Ces fiches présentent le cycle classique de travail avec Git et les principales commandes utilisées pendant le développement.

- [[00 - Base Git/00 - BASE GIT|Base Git]]
    
- [[00 - Base Git/01 - git status|git status]]
    
- [[00 - Base Git/02 - git add|git add]]
    
- [[00 - Base Git/03 - git commit|git commit]]
    
- [[00 - Base Git/04 - git push|git push]]
    
- [[00 - Base Git/05 - git diff|git diff]]
    
- [[00 - Base Git/06 - git log|git log]]
    
- [[00 - Base Git/07 - BRANCHE|Branches]]
    

---

# Vim

Ces fiches regroupent les commandes Vim utiles pour écrire et modifier du code efficacement.

- [[02 - Vim commandes/00 - VIM COMMANDES|Commandes Vim]]
    
- [[02 - Vim commandes/01 - L'essentiel|L’essentiel]]
    
- [[02 - Vim commandes/02 - Les combinaisons|Les combinaisons utiles]]
    

---

# Parcours conseillé

Pour suivre la progression du projet dans l’ordre :

1. [[Projet_Go/00 - Initialisation projet Go/00 - INITIALISATION PROJET|Initialiser le projet]]
    
2. [[Projet_Go/01 - Architecture Go/00 - ARCHITECTURE MINIMALE EN GO|Créer l’architecture minimale]]
    
3. [[Projet_Go/Go/00 - SERVEUR GO|Créer le premier serveur Go]]
    
4. [[Projet_Go/Go/01 - FONCTION|Comprendre les fonctions]]
    
5. [[Projet_Go/Go/02 - PACKAGE|Comprendre les packages]]
    
6. [[Projet_Go/Go/03 - STRUCT ET METHODE|Comprendre les structs et les méthodes]]
    
7. [[Projet_Go/Go/04 - INTERFACE|Comprendre les interfaces]]
    
8. [[Projet_Go/Go/05 - POINTEUR|Comprendre les pointeurs]]
    
9. [[Projet_Go/Go/06 - PRINT|Comprendre l’écriture des réponses]]
    
10. [[Projet_Go/Go/07 - HANDLER|Comprendre les handlers]]
    
11. [[Projet_Go/Go/08 - ROUTEUR|Comprendre le routeur]]
    

---

# Mise à jour de l’index

Lorsqu’une nouvelle fiche est créée :

1. lui donner un nom explicite ;
    
2. lui attribuer un numéro cohérent ;
    
3. la relier à sa fiche principale avec une propriété `MOC` ;
    
4. ajouter son lien dans cet index ;
    
5. enregistrer la modification avec Git.
    

L’index doit rester simple : il sert à retrouver les fiches et à visualiser la progression, sans recopier leur contenu.