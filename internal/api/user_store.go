package api

import (
	"context"

	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
)

type userStoreAdapter struct {
	repo domain.Repository
}

func (a *userStoreAdapter) EnsureUserExists(ctx context.Context, id uuid.UUID, email string) (users.User, error) {
	return users.EnsureUserExists(ctx, a.repo, id, email)
}
