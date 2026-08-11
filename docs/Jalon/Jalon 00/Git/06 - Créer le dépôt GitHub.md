MOC : [[00 - INITIALISATION PROJET]]

---


# Créer un dépôt GitHub

## Objectif

Créer un dépôt distant permettant :

- de sauvegarder le projet ;
- de partager le code ;
- de collaborer avec d'autres développeurs ;
- de publier son travail.

---

## Bonnes pratiques

Le dépôt GitHub doit être créé **après** :

- l'initialisation de Git (`git init`) ;
- la création du `README.md` ;
- la création du `.gitignore` ;
- le premier commit.

Ainsi, le dépôt distant reçoit immédiatement un historique propre.

---

## Paramètres recommandés

|Paramètre|Valeur|
|---|---|
|Repository name|`club-manager`|
|Visibility|Public|
|README|❌ Déjà présent localement|
|.gitignore|❌ Déjà présent localement|
|License|Aucune (ajout possible ultérieurement)|

---

## Pourquoi ne pas créer un README sur GitHub ?

Le projet possède déjà un `README.md`.

Créer un second README sur GitHub produirait un historique différent du dépôt local et nécessiterait une fusion (_merge_) dès le premier `push`.

---

## Dépôt local vs dépôt distant

```
          Dépôt local
      (ordinateur personnel)
               │
        git push / git pull
               │
               ▼
         Dépôt GitHub
```

- **Dépôt local** : environnement de travail sur l'ordinateur.
- **Dépôt distant** : copie hébergée sur GitHub, servant au partage, à la sauvegarde et à la collaboration.


# 1. `git remote add origin`

Commande :

```
git remote add origin https://github.com/grapinou/club-manager.git
```

Décomposons-la :

- `git` : on utilise Git.
- `remote` : on manipule les **dépôts distants**.
- `add` : on en ajoute un nouveau.
- `origin` : c'est le **nom** que l'on donne à ce dépôt distant.
- `https://github.com/grapinou/club-manager.git` : l'adresse du dépôt GitHub.

## Qu'est-ce qu'un "remote" ?

J'aime bien cette représentation :

```
             Dépôt distant
        (GitHub = origin)
               ▲
               │
     push      │      pull
               │
               ▼
         Dépôt local
      (ton ordinateur)
```

Le dépôt local ne connaît pas GitHub par magie.

La commande `git remote add` lui dit :

> "Quand je te parlerai de `origin`, tu sauras qu'il s'agit de ce dépôt GitHub."

Tu peux vérifier que cela a fonctionné avec :

```
git remote -v
```

Tu devrais obtenir quelque chose comme :

```
origin  https://github.com/grapinou/club-manager.git (fetch)
origin  https://github.com/grapinou/club-manager.git (push)
```

# 2. `git branch -M main`

Cette commande sert à **renommer** la branche courante en `main`.

Dans ton cas, je pense qu'elle est inutile.

Pourquoi ?

Parce que ton dépôt est déjà sur `main`.

Tu peux vérifier :

```
git branch
```

Si tu vois :

```
* main
```

➡️ **Ne fais pas cette commande.**

GitHub l'affiche parce qu'il ne sait pas si ton dépôt local est en `master` ou en `main`. C'est une commande générique.

# 3. `git push -u origin main`

C'est la commande la plus intéressante.

Décomposons-la :

```
git push -u origin main
```

- `push` : envoyer des commits vers un dépôt distant.
- `origin` : le dépôt distant.
- `main` : la branche à envoyer.
- `-u` : définir `origin/main` comme branche distante par défaut.

Pourquoi `-u` est-il pratique ?

La première fois, tu écris :

```
git push -u origin main
```

Ensuite, pour les prochains commits, un simple :

```
git push
```

suffira.

Git saura automatiquement :

- où envoyer les commits (`origin`) ;
- quelle branche utiliser (`main`).

# Authentification clée ssh github


## Je te recommande SSH

Pourquoi ?

Parce que :

- c'est la méthode utilisée par la majorité des développeurs ;
- tu n'as pas à ressaisir un mot de passe ou un token ;
- c'est plus sécurisé ;
- tu vas probablement utiliser GitHub régulièrement.

Une fois configuré, un simple :

```
git push
```

fonctionnera sans te demander d'identifiants.

## Première vérification

Peux-tu me donner le résultat de :

```
ls -la ~/.ssh
```

S'il y a déjà des fichiers comme :

```
id_ed25519
id_ed25519.pub
```

ou

```
id_rsa
id_rsa.pub
```

nous les utiliserons.

Sinon, nous en créerons une nouvelle.

Résultat : 

```
sighto@sighto:~/club-manager$ ls -la ~/.ssh
total 8
drwx------  2 sighto sighto 4096 juil. 22  2025 .
drwxr-x--- 26 sighto sighto 4096 juil. 24 17:04 ..
-rw-------  1 sighto sighto    0 juil. 22  2025 authorized_keys
```


Parfait. Cette sortie nous apprend plusieurs choses.

```
~/.ssh
└── authorized_keys
```

Il n'y a **pas de clé SSH personnelle**.

Le fichier `authorized_keys` sert à **autoriser d'autres machines à se connecter à celle-ci** (par exemple si tu utilises SSH pour te connecter à ton PC). Il ne sert pas à t'authentifier auprès de GitHub.

Nous allons donc créer ta première paire de clés SSH. C'est une étape importante que tu réutiliseras pour GitHub, GitLab, certains serveurs Linux, etc.


# Étape 1 : Générer une clé SSH

Je te conseille d'utiliser l'algorithme **Ed25519**, qui est aujourd'hui la recommandation de GitHub.

Exécute :

```
ssh-keygen -t ed25519 -C "remigrapinou@gmail.com"
```

Décomposition de la commande :

- `ssh-keygen` : génère une paire de clés SSH.
- `-t ed25519` : choisit l'algorithme Ed25519.
- `-C` : ajoute un commentaire (ici ton adresse e-mail) pour identifier la clé.

---

## Tu devrais voir quelque chose comme :

```
Generating public/private ed25519 key pair.
Enter file in which to save the key (/home/sighto/.ssh/id_ed25519):
```

👉 **Appuie simplement sur Entrée** pour accepter l'emplacement par défaut.

Ensuite :

```
Enter passphrase (empty for no passphrase):
```

C'est une étape importante.

Tu as deux possibilités :

- **laisser vide** (Entrée) : plus pratique, mais si quelqu'un vole ta clé privée, il peut l'utiliser.
- **mettre une phrase de passe** : plus sécurisé, mais il faudra la saisir (ou utiliser un agent SSH).

### Ce que je te conseille

Comme tu es en train de construire un environnement de développement sérieux, je te recommande **de mettre une phrase de passe**. Tu pourras ensuite utiliser un agent SSH pour éviter de la retaper à chaque fois.

Une bonne phrase de passe est différente d'un mot de passe : elle peut être une petite phrase facile à retenir mais difficile à deviner.

## Ensuite

Une fois la clé créée, exécute :

```
ls -la ~/.ssh
```

Tu devrais voir apparaître :

```
id_ed25519
id_ed25519.pub
```

- `id_ed25519` : **clé privée** → ne jamais partager.
- `id_ed25519.pub` : **clé publique** → c'est celle que nous ajouterons à GitHub.

---

💡 **Petit point sécurité pour ta fiche Obsidian :**

Une paire de clés SSH fonctionne comme une serrure :

- 🔒 **Clé privée** (`id_ed25519`) : reste sur ton ordinateur, ne se partage jamais.
- 🔑 **Clé publique** (`id_ed25519.pub`) : peut être envoyée à GitHub ou à un serveur pour t'autoriser à te connecter.

Tu as maintenant :

```
~/.ssh/
├── authorized_keys
├── id_ed25519
└── id_ed25519.pub
```

Regardons le rôle de chacun :

|Fichier|Rôle|À faire|
|---|---|---|
|`id_ed25519`|🔒 Clé privée|**Ne jamais la partager**|
|`id_ed25519.pub`|🔑 Clé publique|À transmettre à GitHub|
|`authorized_keys`|Autorisations SSH entrantes|Ne concerne pas GitHub|

Les permissions sont également bonnes :

```
-rw------- id_ed25519
```

La clé privée n'est lisible que par ton utilisateur `sighto`.

C'est exactement ce que l'on veut.

---

# Étape suivante : ajouter la clé publique à GitHub

Nous allons récupérer le contenu de la clé publique.

Commande :

```
cat ~/.ssh/id_ed25519.pub
```

Tu vas obtenir une ligne ressemblant à :

```
ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAI.... remigrapinou@gmail.com
```

⚠️ Celle-ci peut être partagée.  
C'est justement celle que GitHub doit connaître.


## Sur GitHub

Va dans :

```
GitHub
 → ta photo de profil (en haut à droite)
 → Settings
 → SSH and GPG keys
 → New SSH key
```

Puis :

**Title**

Je mettrais quelque chose d'explicite :

```
Ubuntu PC - Club Manager
```

ou :

```
sighto Ubuntu desktop
```

Cela sera utile si un jour tu ajoutes ton portable.

**Key type :**

Laisse :

```
Authentication Key
```

**Key :**

Colle le résultat de :

```
cat ~/.ssh/id_ed25519.pub
```

Puis :

```
Add SSH key
```


# Ensuite : tester la connexion

Une fois la clé ajoutée :

```
ssh -T git@github.com
```

La première fois, tu risques d'avoir :

```
The authenticity of host 'github.com' can't be established.
```

Réponds :

```
yes
```

pour déverrouiller utiliser la phrase clée.

Si tout fonctionne, tu devrais voir :

```
Hi grapinou! You've successfully authenticated, but GitHub does not provide shell access.
```

C'est un message normal.

Cela veut dire :

✅ GitHub reconnaît ta clé  
✅ Ton ordinateur est authentifié  
✅ SSH fonctionne

# Ensuite seulement

## Ce que nous voulons faire

Actuellement ton dépôt local pointe probablement encore vers :

```
https://github.com/grapinou/club-manager.git
```

On veut remplacer cette adresse HTTPS par l'adresse SSH :

```
git@github.com:grapinou/club-manager.git
```

La commande correcte est :

```
git remote set-url origin git@github.com:grapinou/club-manager.git
```

Décomposons :

```
git remote
```

→ on manipule les dépôts distants

```
set-url
```

→ on modifie l'adresse d'un dépôt distant existant

```
origin
```

→ le nom donné à ton dépôt GitHub

```
git@github.com:grapinou/club-manager.git
```

→ la nouvelle adresse SSH

---

## Vérifions avant de modifier

Tu peux faire :

```
git remote -v
```

Tu devrais actuellement avoir :

```
origin  https://github.com/grapinou/club-manager.git (fetch)
origin  https://github.com/grapinou/club-manager.git (push)
```

Après :

```
git remote set-url origin git@github.com:grapinou/club-manager.git
```

refais :

```
git remote -v
```

Tu dois obtenir :

```
origin  git@github.com:grapinou/club-manager.git (fetch)
origin  git@github.com:grapinou/club-manager.git (push)
```

## Ensuite seulement

On fera :

```
git push -u origin main
```

Et normalement GitHub acceptera la connexion sans demander de mot de passe.

---

Petite note pour ta future fiche Obsidian : cette confusion est très formatrice.

Une adresse comme :

```
git@github.com:grapinou/club-manager.git
```

n'est pas une commande. C'est une **URL de transport Git**.

Elle peut être utilisée avec :

```
git clone
git remote add
git remote set-url
```

mais elle n'a aucun sens toute seule dans un terminal.

Fais d'abord :

```
git remote -v
```

et montre-moi le résultat. On vérifiera que tout est bien aligné avant de changer quoi que ce soit. 🙂

```
sighto@sighto:~/club-manager$ git push -u origin main
Énumération des objets: 4, fait.
Décompte des objets: 100% (4/4), fait.
Compression par delta en utilisant jusqu'à 12 fils d'exécution
Compression des objets: 100% (4/4), fait.
Écriture des objets: 100% (4/4), 1.42 Kio | 1.42 Mio/s, fait.
Total 4 (delta 0), réutilisés 0 (delta 0), réutilisés du pack 0
To github.com:grapinou/club-manager.git
 * [new branch]      main -> main
la branche 'main' est paramétrée pour suivre 'origin/main'.

```


## Ce que confirme ton message

### 1) Les objets ont été envoyés

```
Énumération des objets: 4, fait.
Décompte des objets: 100% (4/4), fait.
```

Git a trouvé les éléments à envoyer.

Dans ton cas, il s'agit probablement de :

```
.gitignore
README.md
historique du commit
métadonnées Git
```

---

### 2) Création de la branche distante

La ligne importante :

```
* [new branch]      main -> main
```

Signifie :

```
Ordinateur                      GitHub

main
 |
 |  git push
 |
 v
origin/main
```

Ta branche locale `main` a été créée également sur GitHub.

---

### 3) Le suivi de branche est configuré

La dernière ligne est très importante :

```
la branche 'main' est paramétrée pour suivre 'origin/main'.
```

Cela correspond au `-u` de :

```
git push -u origin main
```

Le `-u` signifie :

> "Associe ma branche locale main à la branche distante origin/main."

Maintenant Git sait que :

```
main  <---->  origin/main
```

Donc à partir de maintenant :

```
git push
```

suffira.

Tu n'auras plus besoin d'écrire :

```
git push origin main
```

# Vérification depuis GitHub

Tu peux maintenant aller sur :

```
https://github.com/grapinou/club-manager
```

Tu devrais voir :

- ton README ;
- ton `.gitignore` ;
- ton premier commit.

Le dépôt est officiellement créé. 🚀

---

# Petite étape de propreté maintenant

Je te propose de vérifier ton état local :

```
git status
```

Tu devrais obtenir :

```
On branch main
Your branch is up to date with 'origin/main'.

nothing to commit, working tree clean
```

Cela signifie :

- tu es sur la bonne branche ;
- ton dépôt local est synchronisé ;
- aucune modification en attente.