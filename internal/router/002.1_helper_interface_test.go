package router

import (
	"context"

	"github.com/grapinou/club-manager/internal/database/dbsqlc"
)

type FakeQueries struct{}

func (f FakeQueries) CreateMember(ctx context.Context, arg dbsqlc.CreateMemberParams) (dbsqlc.Member, error) {
	return dbsqlc.Member{}, nil
}

func (f FakeQueries) ListMembers(ctx context.Context) ([]dbsqlc.Member, error) {
	return nil, nil
}
