package database

import (
	"context"

	"github.com/grapinou/club-manager/internal/database/dbsqlc"
)

type Queries interface {
	CreateMember(ctx context.Context, arg dbsqlc.CreateMemberParams) (dbsqlc.Member, error)

	ListMembers(ctx context.Context) ([]dbsqlc.Member, error)
}
