package users

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func TestMapUser(t *testing.T) {
	id := uuid.New()
	now := time.Now().UTC()

	dbUser := sqlc.User{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Forename: sql.NullString{
			String: "John",
			Valid:  true,
		},
		Surname: sql.NullString{
			String: "Doe",
			Valid:  true,
		},
		Email: "john.doe@example.com",
		Nationality: sql.NullString{
			String: "British",
			Valid:  true,
		},
	}

	domainUser := User{
		Id:          id,
		CreatedAt:   now,
		UpdatedAt:   now,
		Forename:    "John",
		Surname:     "Doe",
		Email:       "john.doe@example.com",
		Nationality: "British",
		Role:        "user", // Default in MapUserDomainToDto if not set? No, MapUserDomainToDto just copies it.
	}
	domainUser.Role = "user"

	dtoUser := dto.User{
		Id:          id.String(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Forename:    "John",
		Surname:     "Doe",
		Email:       "john.doe@example.com",
		Nationality: "British",
		Role:        "user",
	}

	t.Run("MapUserDbToDomain", func(t *testing.T) {
		got := MapUserDbToDomain(dbUser)
		// Role is not in DB User struct (sqlc.User might not have it if it's in a separate table or just missing from mapping)
		// Let's check sqlc.User structure from internal/database/sqlc/models.go
		if got.Id != domainUser.Id ||
			got.Forename != domainUser.Forename ||
			got.Surname != domainUser.Surname ||
			got.Email != domainUser.Email ||
			got.Nationality != domainUser.Nationality {
			t.Errorf("expected %+v, got %+v", domainUser, got)
		}
	})

	t.Run("MapUserDomainToDto", func(t *testing.T) {
		got := MapUserDomainToDto(domainUser)
		if got != dtoUser {
			t.Errorf("expected %+v, got %+v", dtoUser, got)
		}
	})

	t.Run("MapUserDtoToDomain", func(t *testing.T) {
		got := MapUserDtoToDomain(dtoUser)
		if got != domainUser {
			t.Errorf("expected %+v, got %+v", domainUser, got)
		}
	})

	t.Run("MapUserDomainToDb", func(t *testing.T) {
		got := MapUserDomainToDb(domainUser)
		if got.ID != dbUser.ID ||
			got.Forename.String != dbUser.Forename.String ||
			got.Surname.String != dbUser.Surname.String ||
			got.Email != dbUser.Email ||
			got.Nationality.String != dbUser.Nationality.String {
			t.Errorf("expected %+v, got %+v", dbUser, got)
		}
	})
}
