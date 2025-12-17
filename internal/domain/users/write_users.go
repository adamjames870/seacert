package users

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewUser(state *internal.ApiState, ctx context.Context, params dto.ParamsAddUser) (User, error) {

	id, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return User{}, errParse
	}

	newUser := sqlc.CreateUserParams{
		ID:    id,
		Email: params.Email,
	}

	dbUser, errWriteNewUser := state.Queries.CreateUser(ctx, newUser)
	if errWriteNewUser != nil {
		return User{}, errWriteNewUser
	}

	apiUser := MapUserDbToDomain(dbUser)
	return apiUser, nil

}
