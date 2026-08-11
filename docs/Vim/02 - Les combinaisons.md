MOC : [[00 - VIM COMMANDES]]

---

# Vim - Les combinaisons indispensables pour le développement

> Les vraies forces de Vim viennent de ses **combinaisons**.  
> On décrit **quoi faire** puis **sur quoi le faire**.

---

# Les objets de texte

| Objet | Signification |
|--------|---------------|
| w | mot |
| s | phrase |
| p | paragraphe |
| " | entre guillemets |
| ' | entre apostrophes |
| ( | entre parenthèses |
| [ | entre crochets |
| { | entre accolades |
| < | entre chevrons |
| t | balise HTML/XML |

---

# Les commandes

| Commande | Action |
|----------|--------|
| d | Delete (supprimer) |
| c | Change (modifier) |
| y | Yank (copier) |
| v | Visual (sélectionner) |

---

# "Inner" et "Around"

Deux lettres sont extrêmement importantes.

| Lettre | Signification |
|---------|---------------|
| i | Inner (à l'intérieur) |
| a | Around (avec les délimiteurs) |

Exemple :

```
("Bonjour")
```

- `ci(` → remplace uniquement **Bonjour**
- `ca(` → remplace **(Bonjour)**

---

# Les commandes les plus utiles

## Modifier rapidement

| Commande | Action |
|----------|--------|
| ciw | Modifier le mot |
| cib | Modifier le contenu entre parenthèses |
| ci( | Idem |
| ci[ | Modifier le contenu entre crochets |
| ci{ | Modifier le contenu entre accolades |
| ci" | Modifier le contenu entre guillemets |
| ci' | Modifier le contenu entre apostrophes |
| cit | Modifier le contenu d'une balise HTML |

---

## Supprimer rapidement

| Commande | Action |
|----------|--------|
| diw | Supprimer un mot |
| di( | Supprimer le contenu des parenthèses |
| di{ | Supprimer le contenu des accolades |
| di" | Supprimer le contenu entre guillemets |
| da" | Supprimer les guillemets compris |
| da( | Supprimer les parenthèses comprises |

---

## Sélection rapide

| Commande | Action |
|----------|--------|
| viw | Sélectionner un mot |
| vi( | Sélectionner le contenu des parenthèses |
| vi{ | Sélectionner le contenu des accolades |
| vi" | Sélectionner le contenu des guillemets |
| va( | Sélectionner les parenthèses comprises |

---

# Parenthèses et accolades

## Aller à la paire correspondante

```
%
```

Très utile pour :

```
if (...) {

}
```

ou

```
func test() {

}
```

Le curseur saute directement sur la parenthèse ou l'accolade correspondante.

---

# Recherche intelligente

## Rechercher le mot sous le curseur

```
*
```

Recherche le mot entier.

Puis :

```
n
```

Suivant.

```
N
```

Précédent.

---

# Répéter

```
.
```

La commande la plus puissante de Vim.

Exemple :

```
ciwUtilisateur<Esc>
```

Puis sur chaque mot suivant :

```
.
```

Le remplacement est reproduit instantanément.

---

# Indentation

## Décaler une ligne

```
>>
```

Décale à droite.

```
<<
```

Décale à gauche.

---

## Réindenter

```
==
```

Réindente la ligne.

```
gg=G
```

Réindente tout le fichier.

Très pratique après un gros copier/coller.

---

# Déplacer des lignes

```
dd
p
```

Déplace une ligne vers le bas.

```
dd
P
```

Déplace une ligne au-dessus.

---

# Copier plusieurs lignes

```
5yy
```

Copie 5 lignes.

```
5dd
```

Supprime 5 lignes.

```
5p
```

Colle 5 lignes copiées.

---

# Les macros

Commencer l'enregistrement :

```
qa
```

Arrêter :

```
q
```

Rejouer :

```
@a
```

Rejouer plusieurs fois :

```
10@a
```

Très utile pour automatiser une tâche répétitive.

---

# Les commandes de développeur

## Aller à la définition (selon l'IDE)

```
gd
```

## Retour

```
Ctrl + o
```

## Avancer

```
Ctrl + i
```

*(Ces commandes dépendent du plugin Vim de ton IDE.)*

---

# Les indispensables

| Commande | Action |
|----------|--------|
| ciw | Modifier un mot |
| ci" | Modifier le texte entre guillemets |
| ci( | Modifier le contenu des parenthèses |
| ci{ | Modifier le contenu des accolades |
| diw | Supprimer un mot |
| vi( | Sélectionner des parenthèses |
| % | Aller à la parenthèse correspondante |
| * | Rechercher le mot sous le curseur |
| . | Répéter la dernière action |
| >> | Indenter |
| << | Désindenter |
| == | Réindenter la ligne |
| gg=G | Réindenter tout le fichier |
| qa | Enregistrer une macro |
| @a | Rejouer une macro |

---

# À retenir

Les commandes suivent presque toujours le même schéma :

```
<action><portée><objet>
```

Exemples :

```
ciw
```

→ Change Inner Word

```
di"
```

→ Delete Inner Quotes

```
va(
```

→ Visual Around Parentheses

```
yi{
```

→ Yank Inner Braces

Une fois cette logique comprise, il devient facile de "deviner" de nombreuses commandes sans les avoir apprises par cœur.