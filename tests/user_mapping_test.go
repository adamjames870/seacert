package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/users"
	"github.com/google/uuid"
)

func TestUserMapping(t *testing.T) {
	id := uuid.New()
	now := time.Now()

	dbUser := sqlc.User{
		ID:          id,
		CreatedAt:   now,
		UpdatedAt:   now,
		Forename:    sql.NullString{String: "John", Valid: true},
		Surname:     sql.NullString{String: "Doe", Valid: true},
		Email:       "john@example.com",
		Nationality: sql.NullString{String: "British", Valid: true},
	}

	// DB to Domain
	domainUser := users.MapUserDbToDomain(dbUser)
	if domainUser.Id != id || domainUser.Forename != "John" || domainUser.Email != "john@example.com" {
		t.Errorf("DB to Domain mapping failed: %+v", domainUser)
	}

	// Domain to DTO
	dtoUser := users.MapUserDomainToDto(domainUser)
	if dtoUser.Id != id.String() || dtoUser.Forename != "John" || dtoUser.Email != "john@example.com" {
		t.Errorf("Domain to DTO mapping failed: %+v", dtoUser)
	}

	// DTO to Domain
	domainUser2 := users.MapUserDtoToDomain(dtoUser)
	if domainUser2.Id != id || domainUser2.Forename != "John" || domainUser2.Email != "john@example.com" {
		t.Errorf("DTO to Domain mapping failed: %+v", domainUser2)
	}

	// Domain to DB
	dbUser2 := users.MapUserDomainToDb(domainUser2)
	if dbUser2.ID != id || dbUser2.Forename.String != "John" || dbUser2.Email != "john@example.com" {
		t.Errorf("Domain to DB mapping failed: %+v", dbUser2)
	}
}
