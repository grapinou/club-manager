package database

import (
	"context"

	"github.com/grapinou/club-manager/internal/database/dbsqlc"
)

type PersonQueries interface {
	CreatePerson(ctx context.Context, arg dbsqlc.CreatePersonParams) (dbsqlc.Person, error)
	ListPersons(ctx context.Context) ([]dbsqlc.Person, error)
	GetPersonByID(ctx context.Context, id int32) (dbsqlc.Person, error)
	UpdatePerson(ctx context.Context, arg dbsqlc.UpdatePersonParams) (dbsqlc.Person, error)
}

type Queries interface {
	PersonQueries
}
