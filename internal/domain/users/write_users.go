package users

import (
	"context"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewUser(state *internal.ApiState, ctx context.Context, params dto.ParamsAddUser) (User, error) {

	id, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return User{}, errParse
	}

	newUser := sqlc.CreateUserParams{
		ID:        id,
		Email:     params.Email,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	dbUser, errWriteNewUser := state.Queries.CreateUser(ctx, newUser)
	if errWriteNewUser != nil {
		return User{}, errWriteNewUser
	}

	apiUser := MapUserDbToDomain(dbUser)
	return apiUser, nil

}
func UpdateUser(state *internal.ApiState, ctx context.Context, params dto.ParamsUpdateUser) (User, error) {

	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return User{}, errParse
	}

	foreName := domain.ToNullStringFromPointer(params.Forename)
	surname := domain.ToNullStringFromPointer(params.Surname)
	nationality := domain.ToNullStringFromPointer(params.Nationality)

	updatedUser := sqlc.UpdateUserParams{
		ID:          uuidId,
		Forename:    foreName,
		Surname:     surname,
		Nationality: nationality,
	}

	dbUser, errUpdateUser := state.Queries.UpdateUser(ctx, updatedUser)
	if errUpdateUser != nil {
		return User{}, errUpdateUser
	}

	return MapUserDbToDomain(dbUser), nil

}
