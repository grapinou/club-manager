MOC : [[00 - VIM COMMANDES]]

---


# Vim - Les commandes essentielles

> Le principe de Vim est simple : **on est toujours dans un mode**.

---

# Les modes

| Mode | Utilisation |
|-------|-------------|
| Normal | Naviguer et lancer des commandes |
| Insertion | Écrire du texte |
| Visuel | Sélectionner du texte |

## Passer en mode insertion

| Commande | Action |
|----------|--------|
| i | Insérer avant le curseur |
| a | Insérer après le curseur |
| o | Nouvelle ligne en dessous |
| O | Nouvelle ligne au-dessus |

## Revenir en mode normal

```
Esc
```

C'est probablement la touche la plus utilisée de Vim.

---

# Déplacements

## Caractères

| Commande | Action |
|----------|--------|
| h | Gauche |
| j | Bas |
| k | Haut |
| l | Droite |

---

## Mots

| Commande | Action |
|----------|--------|
| w | Mot suivant |
| b | Mot précédent |
| e | Fin du mot |

---

## Ligne

| Commande | Action |
|----------|--------|
| 0 | Début de ligne |
| ^ | Premier caractère non vide |
| $ | Fin de ligne |

---

## Document

| Commande | Action |
|----------|--------|
| gg | Début du fichier |
| G | Fin du fichier |
| Ctrl + d | Descendre d'une demi-page |
| Ctrl + u | Monter d'une demi-page |

---

# Copier / Couper / Coller

## Copier

| Commande | Action |
|----------|--------|
| yy | Copier la ligne |
| yw | Copier un mot |
| y$ | Copier jusqu'à la fin de ligne |

---

## Couper (Delete)

| Commande | Action |
|----------|--------|
| dd | Couper la ligne |
| dw | Couper un mot |
| d$ | Couper jusqu'à la fin de ligne |

---

## Coller

| Commande | Action |
|----------|--------|
| p | Coller après |
| P | Coller avant |

---

# Annuler

| Commande | Action |
|----------|--------|
| u | Annuler |
| Ctrl + r | Rétablir |

---

# Modifier rapidement

| Commande | Action |
|----------|--------|
| x | Supprimer un caractère |
| r | Remplacer un caractère |
| cw | Modifier un mot |
| cc | Modifier la ligne |
| C | Modifier jusqu'à la fin de ligne |

---

# Sélection (Mode Visuel)

| Commande | Action |
|----------|--------|
| v | Sélection caractère |
| V | Sélection ligne |
| Ctrl + v | Sélection bloc |

Ensuite :

- y → Copier
- d → Couper
- c → Modifier

---

# Recherche

| Commande | Action |
|----------|--------|
| /mot | Rechercher "mot" |
| n | Résultat suivant |
| N | Résultat précédent |

---

# Répéter

| Commande | Action |
|----------|--------|
| . | Refaire la dernière modification |

Exemple :

```
cwbonjour<Esc>
```

Puis :

```
.
```

refera exactement la même modification.

---

# Les nombres

On peut précéder presque toutes les commandes par un nombre.

Exemples :

```
5j
```

Descendre de 5 lignes.

```
3w
```

Avancer de 3 mots.

```
10dd
```

Supprimer 10 lignes.

```
4yy
```

Copier 4 lignes.

---

# Les commandes à connaître absolument

| Commande | Action |
|----------|--------|
| i | Insérer |
| Esc | Retour normal |
| h j k l | Déplacement |
| w b | Mot suivant/précédent |
| 0 $ | Début/fin de ligne |
| gg G | Début/fin du fichier |
| yy | Copier |
| dd | Couper |
| p | Coller |
| u | Annuler |
| Ctrl+r | Rétablir |
| cw | Modifier un mot |
| x | Supprimer un caractère |
| / | Rechercher |
| n | Résultat suivant |
| . | Répéter la dernière action |

---

# Astuce

Vim est basé sur des verbes.

Exemple :

```
d + w
```

→ supprimer un mot

```
d + $
```

→ supprimer jusqu'à la fin de ligne

```
y + y
```

→ copier une ligne

```
c + w
```

→ modifier un mot

On retient facilement :

- **d** = delete
- **y** = yank (copier)
- **c** = change