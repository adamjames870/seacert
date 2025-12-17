package users

import (
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func MapUserDbToDomain(user sqlc.User) User {

	return User{
		Id:          user.ID,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Forename:    user.Forename.String,
		Surname:     user.Surname.String,
		Email:       user.Email,
		Nationality: user.Nationality.String,
	}

}

func MapUserDomainToDto(user User) dto.User {

	return dto.User{
		Id:          user.Id.String(),
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Forename:    user.Forename,
		Surname:     user.Surname,
		Email:       user.Email,
		Nationality: user.Nationality,
	}

}

func MapUserDtoToDomain(user dto.User) User {

	id, _ := uuid.Parse(user.Id)

	return User{
		Id:          id,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Forename:    user.Forename,
		Surname:     user.Surname,
		Email:       user.Email,
		Nationality: user.Nationality,
	}

}

func MapUserDomainToDb(user User) sqlc.User {

	forename := domain.ToNullString(user.Forename)
	surname := domain.ToNullString(user.Surname)
	nationality := domain.ToNullString(user.Nationality)

	return sqlc.User{
		ID:          user.Id,
		CreatedAt:   user.CreatedAt,
		UpdatedAt:   user.UpdatedAt,
		Forename:    forename,
		Surname:     surname,
		Email:       user.Email,
		Nationality: nationality,
	}

}
