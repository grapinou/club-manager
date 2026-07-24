# Club Manager

Club Manager est une application de gestion d'associations développée en Go.

L'objectif du projet est de concevoir une plateforme robuste, modulaire et facilement maintenable permettant de répondre aux besoins courants des associations, tout en mettant en œuvre les bonnes pratiques de développement logiciel.

Ce projet est également un support d'apprentissage et de démonstration de compétences. Chaque choix d'architecture est réfléchi et documenté afin de construire une application pérenne et évolutive.

## Fonctionnalités prévues

* Gestion des adhérents
* Gestion des rôles et des permissions
* Authentification des utilisateurs
* Gestion des cotisations
* Gestion des cours et des événements
* Gestion documentaire
* Comptabilité
* Modules optionnels (boutique, statistiques, etc.)

## Architecture

Le projet repose sur une architecture orientée **backend**.

Le serveur Go est responsable de l'ensemble de la logique métier :

* gestion des requêtes HTTP ;
* authentification et autorisation ;
* validation des données ;
* accès à la base de données ;
* génération des pages HTML.

L'interface utilisateur est réalisée avec les templates HTML de Go et enrichie avec HTMX afin d'offrir une navigation fluide sans dépendre d'un framework JavaScript.

## Technologies

* Go
* PostgreSQL
* Goose
* sqlc
* HTML Templates
* HTMX
* Git

## Objectifs du projet

* Concevoir une architecture logicielle claire et maintenable.
* Mettre en pratique les bonnes pratiques du développement backend.
* Développer une application facilement extensible grâce à une organisation modulaire.
* Documenter les choix techniques et les étapes de développement.
* Constituer un projet de référence pouvant être présenté lors d'un entretien ou d'un concours.

## État du projet

🚧 Le projet est actuellement en cours de développement.


## Principes de développement

Ce projet est développé selon les principes suivants :

* privilégier la simplicité avant la complexité ;
* écrire un code lisible et facilement maintenable ;
* séparer clairement la logique métier de l'accès aux données et de la présentation ;
* documenter les décisions d'architecture ;
* faire évoluer l'application de manière incrémentale en conservant un historique Git clair et cohérent.

