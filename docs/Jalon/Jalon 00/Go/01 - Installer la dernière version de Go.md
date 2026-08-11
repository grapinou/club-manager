MOC : [[01 - Installer la dernière version de Go]]

---

- Voir la version installé

	Dans le terminal : 
		`go version`
			`go version go1.25.0 linux/amd64`
	
		which go
			/usr/local/go/bin/go

- Désinstaller Go 

	`sudo rm -rf /usr/local/go`

- Aller dans le tmp pour télécharger l'archive Go
	- sera supprimer au prochain démarrage

	`cd /tmp`

- Avec [[01.01 - wget]], télécharger l'archive

	`wget https://go.dev/dl/go1.26.5.linux-amd64.tar.gz`

- Vérifier que l'archive est correcte avec [[01.02 - tar]]

	tar -tzf /tmp/go1.26.5.linux-amd64.tar.gz | head

- Installer Go à partir de l'archive

`sudo tar -C /usr/local -xzf /tmp/go1.26.5.linux-amd64.tar.gz`

remarque : chemin en [[01.03 - absolu]]

- Vérifier le Go installer

	`go version`
		`go version go1.26.5 linux/amd64`

- Mettre Go au PATH

```
export PATH=$PATH:/usr/local/go/bin

export PATH="$PATH:$(go env GOPATH)/bin"
```

- La première ligne

ajoute le dossier contenant l'exécutable `go`.

Avant :

```
PATH=/usr/local/bin:/usr/bin:/bin
```

Après :

```
PATH=/usr/local/bin:/usr/bin:/bin:/usr/local/go/bin
```

Cela permet au shell de trouver la commande :

```
go
```

au lieu d'avoir à taper :

```
/usr/local/go/bin/go
```

Permet de trouver le compilateur et les outils Go.

- La deuxième ligne :

```
export PATH="$PATH:$(go env GOPATH)/bin"
```

est différente.

Elle ajoute le dossier :

```
/home/sighto/go/bin
```

dans ton cas.

Pourquoi ?

Avec Go, quand tu installes un outil en ligne de commande :

```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

Go place l'exécutable ici :

```
$(go env GOPATH)/bin
```

Donc chez toi :

```
/home/sighto/go/bin
```

Tu obtiendras par exemple :

```
/home/sighto/go/bin/goose
/home/sighto/go/bin/sqlc
```

Pour pouvoir ensuite simplement taper :

```
goose version
```

ou :

```
sqlc version
```

il faut que ce dossier soit dans le PATH.

Les commandes installées avec :

```
go install paquet
```

sont placées dans :

```
$(go env GOPATH)/bin
```