package database

import (
	"context"

	"github.com/grapinou/club-manager/internal/database/dbsqlc"
)

type Queries interface {
	CreatePerson(ctx context.Context, arg dbsqlc.CreatePersonParams) (dbsqlc.Person, error)
}
