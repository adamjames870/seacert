package api

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
)

type userStoreAdapter struct {
	state *internal.ApiState
}

func (a *userStoreAdapter) EnsureUserExists(ctx context.Context, id uuid.UUID, email string) (users.User, error) {
	return users.EnsureUserExists(a.state, ctx, id, email)
}
