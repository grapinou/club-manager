package router

import (
	"context"

	"github.com/grapinou/club-manager/internal/database/dbsqlc"
)

type FakeQueries struct{}

func (f FakeQueries) CreatePerson(ctx context.Context, arg dbsqlc.CreatePersonParams) (dbsqlc.Person, error) {
	return dbsqlc.Person{}, nil
}
