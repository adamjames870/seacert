package users

import (
	"context"
	"database/sql"
	"errors"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func EnsureUserExists(state *internal.ApiState, ctx context.Context, id uuid.UUID, email string) (User, error) {

	user, errUser := state.Queries.GetUserByID(ctx, id)

	if errors.Is(errUser, sql.ErrNoRows) {

		userParams := dto.ParamsAddUser{
			Id:    id.String(),
			Email: email,
		}

		return WriteNewUser(state, ctx, userParams)

	}

	if errUser != nil {
		return User{}, errUser
	}

	return MapUserDbToDomain(user), nil

}

func GetUser(state *internal.ApiState, ctx context.Context, id uuid.UUID) (User, error) {

	user, errUser := state.Queries.GetUserByID(ctx, id)
	if errUser != nil {
		return User{}, errUser
	}

	return MapUserDbToDomain(user), nil

}
