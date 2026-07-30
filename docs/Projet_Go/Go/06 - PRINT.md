

---


Concernant la dénomination :

|Fonction|Lettre|Destination|
|---|---|---|
|`Println()`|_(aucune)_|La sortie standard (le terminal)|
|`Fprintln()`|**F = File**|Un objet qui implémente `io.Writer`|
|`Sprintln()`|**S = String**|Une chaîne de caractères|


## `Println()`

```go
fmt.Println("Bonjour")
```

Écrit dans la sortie standard (`stdout`).

Pour un programme en console, cela correspond généralement au terminal.

```
Programme
    │
    ▼
Terminal
```

---

## [[06.01 - fmt.Fprintln]]

Le **F** signifie historiquement **File**.

À l'origine, cette fonction servait à écrire dans un fichier.

Exemple :

```go
fmt.Fprintln(monFichier, "Bonjour")
```

Un fichier (`*os.File`) implémente l'interface `io.Writer`, donc cela fonctionne.

Remarque : [[06.02 - io.Writer]]

Mais en Go, la fonction est plus générale.

Aujourd'hui, on peut écrire dans **n'importe quel objet** qui implémente `io.Writer` :

- un fichier (`*os.File`) ✅
- une connexion réseau (`net.Conn`) ✅
- un `http.ResponseWriter` ✅
- un tampon mémoire (`bytes.Buffer`) ✅

Autrement dit, même si le **F** vient de **File**, il faut plutôt retenir aujourd'hui :

> **F = écrire dans un Writer fourni en argument.**

Par exemple :

```go
fmt.Fprintln(w, "Bienvenue")
```

`w` n'est pas un fichier.

C'est un `http.ResponseWriter`.

---

## `Sprintln()`

Le **S** signifie **String**.

Au lieu d'écrire quelque part, la fonction construit une chaîne de caractères et la renvoie.

```go
message := fmt.Sprintln("Bonjour")
```

Après cette ligne :

```
message
```

contient :

```
Bonjour\n
```

Aucun affichage n'a lieu.

Tu peux ensuite :

```go
fmt.Println(message)
```

ou

```
fmt.Fprintln(w, message)
```

---

## Une logique très élégante

Observe ces trois fonctions :

```go
fmt.Println("Bonjour")
```

↓

Écrit dans le terminal.

---

```go
fmt.Fprintln(w, "Bonjour")
```

↓

Écrit dans un `Writer`.

---

```go
message := fmt.Sprintln("Bonjour")
```

↓

Écrit dans une `String`.

La seule chose qui change est la **destination**.

Le comportement reste le même :

- formater les données ;
- ajouter un retour à la ligne (`ln`).

---

## Et si on enlève le `ln` ?

Tu trouveras aussi :

```go
fmt.Print()
fmt.Fprint()
fmt.Sprint()
```

Cette fois :

- même destination ;
- **mais sans retour à la ligne**.

Le suffixe :

```
ln
```

signifie simplement :

> **Line**

Il ajoute automatiquement :

```
\n
```

à la fin.

---

## Une astuce pour retenir

Tu peux retenir le package `fmt` comme ceci :

|Fonction|Signification|Écrit où ?|
|---|---|---|
|`Print`|Print|Terminal (`stdout`)|
|`Fprint`|**File** (ou plus généralement `io.Writer`)|Objet fourni|
|`Sprint`|String|Variable de type `string`|

Puis les variantes :

- `Print` → sans retour à la ligne.
- `Println` → avec retour à la ligne.
- `Printf` → avec une chaîne de format (`%s`, `%d`, etc.).