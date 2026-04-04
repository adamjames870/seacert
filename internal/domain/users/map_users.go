package users

import (
	"database/sql"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func MapUserDbToDomain(user sqlc.User) User {

	var consentTimestamp *time.Time
	if user.EmailConsentTimestamp.Valid {
		consentTimestamp = &user.EmailConsentTimestamp.Time
	}

	return User{
		Id:                    user.ID,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		Forename:              user.Forename.String,
		Surname:               user.Surname.String,
		Email:                 user.Email,
		Nationality:           user.Nationality.String,
		EmailConsent:          user.EmailConsent,
		EmailConsentTimestamp: consentTimestamp,
		EmailConsentVersion:   user.EmailConsentVersion.String,
		EmailConsentSource:    user.EmailConsentSource.String,
	}

}

func MapUserDomainToDto(user User) dto.User {

	return dto.User{
		Id:                    user.Id.String(),
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		Forename:              user.Forename,
		Surname:               user.Surname,
		Email:                 user.Email,
		Nationality:           user.Nationality,
		Role:                  user.Role,
		EmailConsent:          user.EmailConsent,
		EmailConsentTimestamp: user.EmailConsentTimestamp,
		EmailConsentVersion:   user.EmailConsentVersion,
		EmailConsentSource:    user.EmailConsentSource,
	}

}

func MapUserDtoToDomain(user dto.User) User {

	id, _ := uuid.Parse(user.Id)

	return User{
		Id:                    id,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		Forename:              user.Forename,
		Surname:               user.Surname,
		Email:                 user.Email,
		Nationality:           user.Nationality,
		Role:                  user.Role,
		EmailConsent:          user.EmailConsent,
		EmailConsentTimestamp: user.EmailConsentTimestamp,
		EmailConsentVersion:   user.EmailConsentVersion,
		EmailConsentSource:    user.EmailConsentSource,
	}

}

func MapUserDomainToDb(user User) sqlc.User {

	forename := domain.ToNullString(user.Forename)
	surname := domain.ToNullString(user.Surname)
	nationality := domain.ToNullString(user.Nationality)

	var consentTimestamp sql.NullTime
	if user.EmailConsentTimestamp != nil {
		consentTimestamp = sql.NullTime{Time: *user.EmailConsentTimestamp, Valid: true}
	}

	return sqlc.User{
		ID:                    user.Id,
		CreatedAt:             user.CreatedAt,
		UpdatedAt:             user.UpdatedAt,
		Forename:              forename,
		Surname:               surname,
		Email:                 user.Email,
		Nationality:           nationality,
		EmailConsent:          user.EmailConsent,
		EmailConsentTimestamp: consentTimestamp,
		EmailConsentVersion:   domain.ToNullString(user.EmailConsentVersion),
		EmailConsentSource:    domain.ToNullString(user.EmailConsentSource),
	}

}
