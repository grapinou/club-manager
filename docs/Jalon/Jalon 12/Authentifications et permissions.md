

# Club Manager — Lignes directrices d’architecture

## Objectif

Club Manager doit rester une application :

- simple à comprendre ;
    
- modulaire ;
    
- sécurisée ;
    
- adaptable aux besoins de différentes associations ;
    
- maintenable dans le temps.
    

L’idée générale est de distinguer clairement :

1. le **cœur de l’application** ;
    
2. les **modules métier** ;
    
3. les **données personnalisées** créées par un administrateur.
    

---

## 1. Le cœur de l’application

Le **Core** regroupe les éléments indispensables au fonctionnement de Club Manager.

Exemples :

- configuration ;
    
- utilisateurs ;
    
- authentification ;
    
- rôles ;
    
- permissions ;
    
- routage ;
    
- gestion des modules activés.
    

Le Core ne doit pas contenir toute la logique métier de l’application.

Il fournit plutôt les briques communes utilisées par les différents modules.

---

## 2. Les modules métier

Un module correspond à une fonctionnalité prévue et développée dans Club Manager.

Exemples :

- Members ;
    
- Accounting ;
    
- Memberships ;
    
- Equipment ;
    
- Events ;
    
- Shop.
    

Un module peut contenir :

- ses routes ;
    
- ses handlers ;
    
- ses vues ;
    
- ses requêtes SQL ;
    
- ses tables ;
    
- ses permissions ;
    
- ses règles métier.
    

Exemple pour le module `Members` :

```
Members
├── members.read
├── members.create
├── members.update
└── members.delete
```

Un module peut être activé ou désactivé selon les besoins de l’association.

Si un module est désactivé, ses routes peuvent simplement ne pas être enregistrées.

```
module activé ?
      ↓ oui
enregistrement des routes
      ↓
contrôle des permissions
      ↓
handler
```

---

## 3. Permissions et sécurité

Les permissions doivent être liées aux actions, et non directement aux rôles.

On préfère :

```
members.read
members.create
members.update
members.delete
```

plutôt que de disperser dans le code des conditions comme :

```
if user.Role == "admin" || user.Role == "secretary" {
    // ...
}
```

Un rôle devient alors simplement un regroupement de permissions.

Exemple :

```
Secrétaire
├── members.read
├── members.create
└── members.update
```

Le routeur associe une route à un middleware d’autorisation.

```
requête
   ↓
routeur
   ↓
middleware de permission
   ↓
handler
   ↓
base de données
```

Ainsi, les handlers peuvent rester concentrés sur leur responsabilité métier.

---

## 4. Les tables ne gèrent pas les permissions applicatives

Une table doit avant tout représenter et stocker des données cohérentes.

Exemple :

```
members
├── id
├── first_name
├── last_name
├── birth_date
├── email
└── created_at
```

On évite d’ajouter dans une table métier des colonnes du type :

```
can_be_seen_by_secretary
can_be_modified_by_admin
```

Ces règles appartiennent au système d’autorisation.

La base de données conserve néanmoins ses propres responsabilités :

- contraintes `NOT NULL` ;
    
- clés primaires ;
    
- clés étrangères ;
    
- contraintes d’unicité ;
    
- cohérence des données ;
    
- droits PostgreSQL de l’utilisateur utilisé par l’application.
    

---

## 5. Développement avec un administrateur

Pendant le développement, il est possible d’utiliser un utilisateur administrateur fictif.

Cela permet de développer les fonctionnalités sans construire immédiatement tout le système de connexion.

Exemple conceptuel :

```
devUser := User{
    ID:   1,
    Role: "admin",
}
```

Il faudra néanmoins tester rapidement plusieurs cas :

```
Anonyme
Membre
Membre du bureau
Administrateur
```

L’administrateur ne doit pas être le seul cas testé, car il possède généralement toutes les permissions et peut masquer des erreurs d’autorisation.

---

## 6. Ne pas permettre à l’administrateur de créer directement des tables SQL

Techniquement, il serait possible de permettre à un administrateur d’exécuter des créations de tables.

Cependant, ce choix compliquerait fortement l’architecture.

Aujourd’hui, le schéma suit une chaîne maîtrisée :

```
Goose
  ↓
PostgreSQL
  ↓
sqlc
  ↓
code Go typé
```

Si un administrateur pouvait créer librement une table :

```
CREATE TABLE materiel (...);
```

le code Go compilé ne connaîtrait pas cette table.

`sqlc` ne pourrait pas avoir généré les types et méthodes correspondants.

On perdrait alors plusieurs avantages :

- schéma connu ;
    
- migrations maîtrisées ;
    
- historique clair ;
    
- génération de code typé ;
    
- tests reproductibles ;
    
- maintenance simple.
    

### Ligne directrice

> L’administrateur ne modifie jamais directement le schéma PostgreSQL.

Les vraies tables restent sous le contrôle du développement et des migrations Goose.

---

## 7. Permettre malgré tout des données personnalisées

L’administrateur peut avoir besoin de stocker des informations que Club Manager n’avait pas prévues.

Exemples :

- clés du dojo ;
    
- matériel particulier ;
    
- véhicules ;
    
- costumes ;
    
- prêts ;
    
- équipements spécifiques ;
    
- registres internes.
    

Plutôt que de créer une vraie table PostgreSQL, l’administrateur pourrait créer un **registre personnalisé**.

Exemple :

```
Registre : Matériel

Champs :
- Nom
- Quantité
- État
- Date d’achat
```

L’utilisateur a l’impression de créer une nouvelle structure de données, mais PostgreSQL conserve un schéma stable.

---

## 8. Structure possible pour les données personnalisées

Une solution consiste à utiliser des tables génériques.

Exemple :

```
custom_collections
├── id
└── name
```

```
custom_fields
├── id
├── collection_id
├── name
└── type
```

```
custom_records
├── id
├── collection_id
└── data
```

La colonne `data` pourrait éventuellement utiliser le type PostgreSQL `JSONB`.

Exemple :

```
{
    "nom": "Tatami",
    "quantite": 40,
    "etat": "bon"
}
```

Un autre registre pourrait utiliser la même infrastructure avec une structure différente :

```
{
    "nom": "Caméra",
    "numero_serie": "ABC123",
    "date_achat": "2026-05-12"
}
```

Le choix de `JSONB` devra être étudié avant implémentation.

L’objectif est surtout de conserver cette possibilité architecturale.

---

## 9. Module métier et registre personnalisé ne sont pas équivalents

Un module possède une vraie logique métier.

Exemple :

```
Members
   ↓
validation
   ↓
création
   ↓
cotisation
   ↓
licence
   ↓
autres règles métier
```

Un registre personnalisé est volontairement plus simple :

```
Créer
Lire
Modifier
Supprimer
```

Il sert essentiellement à stocker des données structurées que Club Manager ne connaît pas à l’avance.

### Règle

> Tout ne doit pas devenir dynamique.

Les fonctionnalités importantes et complexes restent de vrais modules développés et testés.

Les registres personnalisés servent aux besoins simples et spécifiques.

---

## 10. Permissions des registres personnalisés

Les registres personnalisés pourront eux aussi s’intégrer au système de permissions.

Exemple :

```
custom.material.read
custom.material.create
custom.material.update
custom.material.delete
```

L’administrateur pourra alors attribuer ces permissions à différents rôles.

Exemple :

```
Responsable matériel
├── lecture
├── création
├── modification
└── suppression

Secrétaire
└── lecture

Membre
└── aucune permission
```

---

## Architecture générale envisagée

```
                        CLUB MANAGER
                             │
             ┌───────────────┴────────────────┐
             │                                │
            Core                       Fonctionnalités
             │                                │
    ┌────────┼────────┐              ┌────────┴─────────┐
    │        │        │              │                  │
  Users    Roles  Permissions      Modules        Custom Data
```

Ou, en termes de responsabilités :

```
CORE
│
├── utilisateurs
├── authentification
├── permissions
├── configuration
└── gestion des modules

MODULES
│
├── Members
├── Accounting
├── Equipment
└── ...

CUSTOM DATA
│
├── registres personnalisés
├── champs configurables
└── données définies par l’association
```

---

## Principes à retenir

### 1. Le schéma SQL reste maîtrisé

Les tables physiques sont créées et modifiées par des migrations Goose.

```
Goose → PostgreSQL → sqlc → Go
```

### 2. Les permissions restent hors des tables métier

Les tables stockent les données.

Le système d’autorisation décide qui peut agir dessus.

### 3. Les rôles regroupent des permissions

```
Rôle
  ↓
ensemble de permissions
```

### 4. Les modules portent la logique métier

Une fonctionnalité complexe devient un vrai module.

### 5. Les besoins spécifiques simples utilisent des registres personnalisés

L’administrateur peut définir ses propres structures sans modifier directement PostgreSQL.

### 6. Ne pas abstraire trop tôt

La modularité doit apparaître à partir de besoins réels.

On commence par `Members`, puis on observe ce qui est réutilisable lorsqu’un deuxième module apparaît.

---

## Comprendre et retenir

> **Le Core fournit les briques communes.**

> **Un module représente une vraie fonctionnalité métier.**

> **Un registre personnalisé permet de stocker des données imprévues sans modifier le schéma SQL.**

> **Les permissions contrôlent les actions, pas les tables.**

> **Les rôles regroupent des permissions.**

> **L’administrateur configure l’application, mais ne modifie pas directement PostgreSQL.**

> **Goose reste la source de vérité pour l’évolution du schéma.**

Cette séparation doit permettre à Club Manager de rester à la fois simple, sûr et adaptable.