package users

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func CreateUser(ctx context.Context, repo domain.Repository, params dto.ParamsAddUser) (User, error) {
	id, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return User{}, domain.ErrInvalidInput
	}

	newUser := sqlc.CreateUserParams{
		ID:                    id,
		Email:                 params.Email,
		CreatedAt:             time.Now(),
		UpdatedAt:             time.Now(),
		EmailConsent:          false,
		EmailConsentTimestamp: sql.NullTime{},
		EmailConsentVersion:   sql.NullString{},
		EmailConsentSource:    sql.NullString{},
	}

	dbUser, err := repo.CreateUser(ctx, newUser)
	if err != nil {
		return User{}, err
	}

	return MapUserDbToDomain(dbUser), nil
}

func UpdateUser(ctx context.Context, repo domain.Repository, params dto.ParamsUpdateUser) (User, error) {
	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return User{}, domain.ErrInvalidInput
	}

	var consentTimestamp sql.NullTime
	if params.EmailConsent != nil {
		consentTimestamp = sql.NullTime{Time: time.Now(), Valid: true}
	}

	updatedUser := sqlc.UpdateUserParams{
		ID:                    uuidId,
		Forename:              domain.ToNullStringFromPointer(params.Forename),
		Surname:               domain.ToNullStringFromPointer(params.Surname),
		Nationality:           domain.ToNullStringFromPointer(params.Nationality),
		EmailConsent:          domain.ToNullBoolFromPointer(params.EmailConsent),
		EmailConsentTimestamp: consentTimestamp,
		EmailConsentVersion:   domain.ToNullStringFromPointer(params.EmailConsentVersion),
		EmailConsentSource:    domain.ToNullStringFromPointer(params.EmailConsentSource),
	}

	dbUser, err := repo.UpdateUser(ctx, updatedUser)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, domain.ErrNotFound
		}
		return User{}, err
	}

	return MapUserDbToDomain(dbUser), nil
}

func GetUser(ctx context.Context, repo domain.Repository, id uuid.UUID) (User, error) {
	user, err := repo.GetUserByID(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return User{}, domain.ErrNotFound
		}
		return User{}, err
	}

	return MapUserDbToDomain(user), nil
}

func EnsureUserExists(ctx context.Context, repo domain.Repository, id uuid.UUID, email string) (User, error) {
	user, err := repo.GetUserByID(ctx, id)
	if err == nil {
		return MapUserDbToDomain(user), nil
	}

	if errors.Is(err, sql.ErrNoRows) {
		userParams := dto.ParamsAddUser{
			Id:    id.String(),
			Email: email,
		}
		return CreateUser(ctx, repo, userParams)
	}

	return User{}, err
}
