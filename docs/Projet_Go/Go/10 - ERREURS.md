
---

## Définition

Une erreur indique qu'une opération n'a pas pu se dérouler comme prévu.

Exemples :

- un fichier n'existe pas ;
    
- une donnée est invalide ;
    
- un utilisateur n'est pas trouvé ;
    
- une connexion à la base de données échoue ;
    
- une valeur ne peut pas être convertie.
    

En Go, une erreur est généralement représentée par le type :

```go
error
```

---

## `error` est une interface

Le type `error` n'est pas une struct.

C'est une interface intégrée au langage Go :

```go
type error interface {
	Error() string
}
```

Cette interface demande une seule méthode :

```go
Error() string
```

Un type qui possède cette méthode satisfait automatiquement l'interface `error`.

```text
Un type possède Error() string
              ↓
Le type satisfait l'interface error
              ↓
Il peut être utilisé comme une erreur
```

---

## Une fonction qui peut échouer

En Go, une fonction qui peut échouer retourne souvent deux valeurs :

```go
résultat, erreur
```

Exemple :

```go
func FindMember(id int) (string, error) {
	// recherche du membre
}
```

Cette fonction retourne :

```go
string
```

pour le résultat, et :

```go
error
```

pour indiquer un éventuel problème.

Utilisation :

```go
member, err := FindMember(42)
```

Les deux variables reçoivent les valeurs retournées :

```text
member → résultat de la recherche
err    → erreur éventuelle
```

---

## Vérifier une erreur

La pratique courante consiste à vérifier immédiatement l'erreur :

```go
member, err := FindMember(42)

if err != nil {
	// gestion de l'erreur
}
```

`nil` représente ici l'absence d'erreur.

```text
err == nil
→ aucune erreur

err != nil
→ une erreur s'est produite
```

Exemple complet :

```go
member, err := FindMember(42)

if err != nil {
	fmt.Println("Impossible de trouver le membre :", err)
	return
}

fmt.Println(member)
```

---

## Pourquoi vérifier l'erreur immédiatement ?

Cela permet de séparer clairement deux chemins possibles :

```text
Appel de la fonction
        ↓
Une erreur existe ?
   ↙             ↘
 oui             non
  ↓               ↓
traiter         continuer
l'erreur        normalement
```

La structure suivante est donc très fréquente en Go :

```go
result, err := operation()

if err != nil {
	return err
}

// utilisation de result
```

---

## Créer une erreur avec `errors.New`

Le package `errors` permet de créer une erreur simple :

```go
import "errors"
```

Exemple :

```go
err := errors.New("membre introuvable")
```

La valeur créée peut être retournée par une fonction :

```go
func FindMember(id int) (string, error) {
	if id <= 0 {
		return "", errors.New("identifiant invalide")
	}

	return "Rémi", nil
}
```

Lorsque l'identifiant est invalide :

```go
return "", errors.New("identifiant invalide")
```

La fonction retourne :

- une chaîne vide comme résultat ;
    
- une erreur.
    

Lorsque tout fonctionne :

```go
return "Rémi", nil
```

La fonction retourne :

- le résultat ;
    
- aucune erreur.
    

---

## Créer une erreur avec `fmt.Errorf`

`fmt.Errorf` permet de construire une erreur contenant des informations variables.

```go
err := fmt.Errorf("membre %d introuvable", id)
```

Exemple :

```go
func FindMember(id int) (string, error) {
	if id != 42 {
		return "", fmt.Errorf(
			"membre avec l'identifiant %d introuvable",
			id,
		)
	}

	return "Rémi", nil
}
```

Pour l'identifiant `12`, le message produit sera :

```text
membre avec l'identifiant 12 introuvable
```

---

## `fmt.Errorf` et `t.Errorf`

Il faut distinguer :

```go
fmt.Errorf(...)
```

et :

```go
t.Errorf(...)
```

### `fmt.Errorf`

```go
err := fmt.Errorf("membre %d introuvable", id)
```

`fmt.Errorf` :

- construit une erreur ;
    
- retourne une valeur de type `error`.
    

Sa valeur peut être stockée :

```go
err := fmt.Errorf("une erreur est survenue")
```

ou retournée :

```go
return fmt.Errorf("une erreur est survenue")
```

### `t.Errorf`

```go
t.Errorf("statut obtenu : %d", response.Code)
```

`t.Errorf` :

- appartient au type `*testing.T` ;
    
- signale l'échec d'un test ;
    
- affiche un message formaté ;
    
- ne retourne pas une valeur de type `error`.
    

Ceci est donc incorrect :

```go
err := t.Errorf("le test a échoué")
```

`t.Errorf` ne retourne rien.

---

## Toutes les variables possèdent-elles une erreur ?

Non.

Une variable possède uniquement :

- les données définies par son type ;
    
- les méthodes définies pour son type.
    

Exemple :

```go
type Dog struct {
	Name string
}
```

Une variable de type `Dog` ne possède pas automatiquement une méthode `Errorf` :

```go
dog := Dog{Name: "Bob"}

dog.Errorf("erreur")
```

Ce code ne compile pas, car `Dog` ne définit aucune méthode `Errorf`.

Dans un test :

```go
t.Errorf(...)
```

fonctionne parce que `t` est de type :

```go
*testing.T
```

et que `testing.T` possède une méthode `Errorf`.

---

## Créer un type d'erreur personnalisé

Il est possible de créer une struct représentant une erreur :

```go
type MemberNotFoundError struct {
	ID int
}
```

Pour qu'elle satisfasse l'interface `error`, il faut lui ajouter :

```go
func (e MemberNotFoundError) Error() string {
	return fmt.Sprintf(
		"le membre %d est introuvable",
		e.ID,
	)
}
```

Le type possède maintenant la méthode demandée par l'interface :

```go
Error() string
```

Il peut donc être utilisé comme une erreur :

```go
func FindMember(id int) (string, error) {
	return "", MemberNotFoundError{ID: id}
}
```

---

## Retourner une erreur existante

Une fonction peut recevoir une erreur d'une autre fonction et la retourner :

```go
func LoadConfig() error {
	file, err := os.Open("config.json")

	if err != nil {
		return err
	}

	defer file.Close()

	return nil
}
```

Ici :

```go
os.Open("config.json")
```

peut retourner une erreur.

Si une erreur existe :

```go
if err != nil {
	return err
}
```

elle est transmise à l'appelant.

---

## Ajouter du contexte à une erreur

Retourner directement une erreur peut manquer de contexte :

```go
return err
```

Il est souvent préférable d'ajouter une information :

```go
return fmt.Errorf(
	"impossible d'ouvrir la configuration : %w",
	err,
)
```

Le symbole :

```go
%w
```

permet d'envelopper l'erreur originale.

Le message final pourra ressembler à :

```text
impossible d'ouvrir la configuration : fichier introuvable
```

L'erreur originale reste accessible.

---

## Afficher une erreur

Une erreur peut être affichée directement :

```go
fmt.Println(err)
```

Go appelle alors sa méthode :

```go
Error() string
```

Ces deux écritures produisent généralement le même texte :

```go
fmt.Println(err)
```

```go
fmt.Println(err.Error())
```

La première forme est habituellement suffisante.

---

## Erreur et arrêt du programme

Une erreur n'arrête pas automatiquement le programme.

Elle est simplement une valeur.

```go
err := errors.New("une erreur est survenue")
```

Le développeur doit décider quoi faire :

- retourner l'erreur ;
    
- afficher un message ;
    
- réessayer l'opération ;
    
- envoyer une réponse HTTP ;
    
- utiliser une valeur par défaut ;
    
- arrêter une partie du traitement.
    

Cette approche distingue Go des langages qui utilisent principalement des exceptions.

---

## Exemple avec un handler HTTP

```go
func MemberHandler(
	w http.ResponseWriter,
	r *http.Request,
) {
	member, err := FindMember(42)

	if err != nil {
		http.Error(
			w,
			"membre introuvable",
			http.StatusNotFound,
		)
		return
	}

	fmt.Fprintln(w, member)
}
```

Le déroulement est :

```text
Rechercher le membre
        ↓
Une erreur existe ?
   ↙             ↘
 oui             non
  ↓               ↓
réponse 404     afficher le membre
```

---

## Les erreurs dans les tests

Dans un test, nous vérifions parfois qu'une fonction retourne bien une erreur.

```go
func TestFindMemberWithInvalidID(t *testing.T) {
	_, err := FindMember(-1)

	if err == nil {
		t.Error(
			"une erreur était attendue",
		)
	}
}
```

Ici, le test échoue si aucune erreur n'est retournée.

```text
err == nil
→ la fonction n'a pas produit l'erreur attendue
→ le test échoue
```

Nous pouvons aussi vérifier qu'aucune erreur n'est produite :

```go
func TestFindMember(t *testing.T) {
	_, err := FindMember(42)

	if err != nil {
		t.Errorf(
			"aucune erreur attendue ; erreur obtenue : %v",
			err,
		)
	}
}
```

---

## Comprendre et retenir

En Go, une erreur est une valeur dont le type satisfait l'interface :

```go
type error interface {
	Error() string
}
```

Une fonction qui peut échouer retourne souvent :

```go
result, err
```

L'erreur est généralement vérifiée immédiatement :

```go
if err != nil {
	// traiter l'erreur
}
```

`nil` signifie qu'aucune erreur n'est présente.

Pour créer une erreur simple :

```go
errors.New("message")
```

Pour créer une erreur formatée :

```go
fmt.Errorf("message : %d", value)
```

Il faut distinguer :

```text
fmt.Errorf → crée et retourne une error

t.Errorf   → signale l'échec d'un test
```

Une erreur ne provoque pas automatiquement l'arrêt du programme.

Le développeur doit décider comment la traiter.

Remarque : 

## `%w` signifie-t-il « write » ?

Non. Il faut surtout le retenir comme :

```
w → wrap
```

c’est-à-dire :

```
envelopper une erreur
```

## À retenir

```
fmt.Errorf("contexte : %v", err)
```

affiche l’erreur comme une valeur.

```
fmt.Errorf("contexte : %w", err)
```

affiche **et enveloppe** l’erreur.

La bonne pratique, lorsqu’on retourne une erreur en ajoutant du contexte, est généralement :

```
return fmt.Errorf("description de l'opération : %w", err)
```