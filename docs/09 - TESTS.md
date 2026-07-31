
---

- [[09.01 - Anatomie d'un premier test HTTP en Go]]
# Pourquoi tester son code ?

## Définition

Un test est un morceau de code qui vérifie automatiquement qu'une partie de l'application se comporte comme prévu.

Un test peut par exemple vérifier que :

- une page existe ;
    
- une requête HTTP renvoie le bon statut ;
    
- une fonction retourne la bonne valeur ;
    
- une erreur est correctement détectée ;
    
- une modification n'a pas cassé un comportement existant.
    

Un test compare donc :

```text
le résultat obtenu
        avec
le résultat attendu
```

---

## Pourquoi écrire des tests ?

### Vérifier que le code fonctionne

Lorsqu'on écrit une fonction, on peut la tester manuellement.

Par exemple, pour vérifier une page web, on peut :

1. lancer le serveur ;
    
2. ouvrir le navigateur ;
    
3. saisir l'adresse ;
    
4. regarder le résultat.
    

Cette méthode fonctionne, mais elle devient rapidement longue et répétitive.

Un test permet d'effectuer cette vérification automatiquement.

```text
Commande de test
       ↓
Exécution du code
       ↓
Comparaison avec le résultat attendu
       ↓
Succès ou échec
```

---

## Éviter les régressions

Une régression est un comportement qui fonctionnait auparavant, mais qui ne fonctionne plus après une modification.

Exemple :

```text
La page /contact fonctionne.
        ↓
Nous modifions le routeur.
        ↓
La page /contact ne fonctionne plus.
```

Sans test, le problème peut passer inaperçu.

Avec un test, l'erreur est immédiatement signalée :

```text
FAIL
```

Les tests permettent donc de protéger les fonctionnalités déjà développées.

---

## Modifier le code avec plus de confiance

Lorsqu'un projet grandit, une modification peut avoir des conséquences ailleurs dans l'application.

Les tests permettent de vérifier rapidement que les anciens comportements sont toujours valides.

Ils rendent notamment plus sereines les opérations suivantes :

- déplacer du code ;
    
- renommer une fonction ;
    
- modifier une architecture ;
    
- remplacer une implémentation ;
    
- simplifier une fonction ;
    
- ajouter une nouvelle fonctionnalité.
    

Les tests ne garantissent pas que le code ne contient aucun problème.

Ils apportent cependant une protection contre les erreurs déjà envisagées et décrites dans les tests.

---

## Les tests décrivent le comportement attendu

Un test ne sert pas uniquement à trouver des erreurs.

Il permet aussi de décrire ce que l'application doit faire.

Exemple :

```go
func TestHomeHandler(t *testing.T)
```

Ce nom indique qu'il existe un comportement attendu pour `HomeHandler`.

Un autre développeur peut lire le test pour comprendre :

- comment utiliser la fonction ;
    
- quelles données lui transmettre ;
    
- quelle réponse attendre ;
    
- quels cas sont importants.
    

Un test peut donc également servir de documentation.

---

## Les tests obligent à préciser les attentes

Avant d'écrire un test, il faut répondre à une question :

> Quel comportement voulons-nous obtenir ?

Pour une page web, nous pouvons par exemple définir que :

```text
GET /
```

doit :

- répondre avec un statut `200 OK` ;
    
- afficher la page d'accueil.
    

Pour une page inconnue :

```text
GET /page-inconnue
```

nous pouvons décider qu'elle doit :

- répondre avec un statut `404 Not Found`.
    

Le test oblige donc à transformer une intention vague en comportement précis.

---

## Test manuel et test automatisé

### Test manuel

```text
Lancer le serveur
Ouvrir le navigateur
Visiter la page
Observer le résultat
```

Avantages :

- simple à comprendre ;
    
- utile pour vérifier l'apparence générale ;
    
- permet d'observer l'application comme un utilisateur.
    

Inconvénients :

- lent ;
    
- répétitif ;
    
- facile à oublier ;
    
- difficile à reproduire exactement.
    

### Test automatisé

```text
go test ./...
```

Avantages :

- rapide ;
    
- répétable ;
    
- précis ;
    
- exécutable après chaque modification ;
    
- utile pour détecter les régressions.
    

Inconvénients :

- demande un apprentissage ;
    
- nécessite de définir les comportements attendus ;
    
- représente du code supplémentaire à maintenir.
    

Les tests manuels et automatisés sont complémentaires.

---

## Que faut-il tester ?

Il faut principalement tester les comportements qui appartiennent à notre application.

Dans Club Manager, nous pourrons vérifier que :

```text
GET /          → page d'accueil
GET /club      → présentation du club
GET /contact   → page de contact
GET /rules     → règlement intérieur
GET /inconnue  → erreur 404
```

Nous ne cherchons pas à tester le fonctionnement interne de la bibliothèque standard de Go.

Par exemple, nous n'avons pas besoin de vérifier que `http.NewServeMux()` fonctionne correctement.

Ce comportement est déjà testé par les développeurs de Go.

Nous devons vérifier que nous avons correctement configuré et utilisé le routeur.

---

## Faut-il tout tester ?

Non.

Tester chaque ligne de code n'est pas nécessairement utile.

Il faut privilégier les comportements importants :

- fonctionnalités principales ;
    
- règles métier ;
    
- traitements pouvant produire des erreurs ;
    
- cas limites ;
    
- éléments pouvant facilement être cassés par une modification.
    

Pour le premier jalon de Club Manager, quelques tests simples suffisent.

L'objectif est d'apprendre progressivement et de construire de bonnes habitudes.

---

## Quand écrire les tests ?

Les tests peuvent être écrits :

- avant le code ;
    
- pendant le développement ;
    
- après l'implémentation ;
    
- lorsqu'un bug est découvert.
    

Dans notre progression, nous allons commencer par ajouter des tests après avoir créé une première architecture fonctionnelle.

Cela permet de comprendre le fonctionnement des tests sur du code déjà connu.

Plus tard, nous pourrons éventuellement écrire certains tests avant l'implémentation.

---

## Organisation des fichiers de test en Go

Un fichier de test Go se termine par :

```text
_test.go
```

Exemple :

```text
home.go
home_test.go
```

Le fichier `home.go` contient le code de l'application.

Le fichier `home_test.go` contient les tests associés.

Une organisation possible dans Club Manager serait :

```text
internal/
└── handlers/
    ├── home.go
    └── home_test.go
```

---

## Exécuter les tests

Pour tester tous les packages du projet :

```bash
go test ./...
```

Pour tester le package courant :

```bash
go test
```

Lorsque les tests réussissent :

```text
ok
```

Lorsqu'un test échoue :

```text
FAIL
```

Le message d'erreur doit aider à comprendre :

- ce qui était attendu ;
    
- ce qui a réellement été obtenu.
    

---

## Les tests dans la progression de Club Manager

Notre progression devient :

```text
Architecture minimale
        ↓
Séparation des responsabilités
        ↓
Handlers
        ↓
Routeur
        ↓
Premiers tests automatisés
```

L'ajout des tests à ce stade est intéressant parce que :

- le projet est encore petit ;
    
- le code est facile à comprendre ;
    
- les comportements sont simples ;
    
- les erreurs sont plus faciles à identifier ;
    
- nous pouvons apprendre sans logique métier complexe.
    

---

## Comprendre et retenir

Un test vérifie automatiquement qu'un comportement correspond à ce qui était attendu.

```text
Résultat obtenu == résultat attendu
```

Les tests permettent principalement de :

- vérifier le fonctionnement du code ;
    
- détecter les régressions ;
    
- modifier le projet avec plus de confiance ;
    
- documenter les comportements attendus ;
    
- préciser les règles de l'application.
    

Un test ne prouve pas que le programme est parfait.

Il vérifie uniquement les situations que nous avons choisi de tester.

Dans Club Manager, nous commencerons par tester les pages et le routeur avant d'aborder des comportements plus complexes.