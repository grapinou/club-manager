MOC : [[00 - INITIALISATION PROJET]]

---


## Étape 1 : Vérifier que Git est installé

Dans un terminal, tape :

```
git --version
```

Si Git est installé, tu obtiendras quelque chose comme :

```
git version 2.43.0
```

(peu importe le numéro exact).

Si tu obtiens un message du type :

```
git: command not found
```

alors Git n'est pas installé.

### 2. Emplacement de Git

```
which git
```

Résultat :

```
/usr/bin/git
```

C'est exactement ce à quoi on s'attend sur Ubuntu.

---

### 3. Pourquoi `type -a git` affiche deux chemins ?

Tu as obtenu :

```
git est /usr/bin/git
git est /bin/git
```

C'est normal.

Sur les versions récentes d'Ubuntu, `/bin` est un **lien symbolique** vers `/usr/bin` (ou inversement selon l'organisation du système).

Tu peux le vérifier avec :

```
ls -ld /bin
```

Tu devrais obtenir quelque chose comme :

```
lrwxrwxrwx ... /bin -> usr/bin
```

Autrement dit :

- `/bin/git`
- `/usr/bin/git`

désignent en réalité **le même exécutable**.

Il n'y a donc pas deux installations de Git.

---

### 4. Configuration Git

Tu as :

```
user.name=grapinou
user.email=remigrapinou@gmail.com
```

Git sait donc qui tu es lorsque tu feras un commit.