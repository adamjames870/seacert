package issuers

import (
	"database/sql"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func MapIssuerDbToDomain(issuer sqlc.Issuer) Issuer {

	return Issuer{
		Id:        issuer.ID,
		CreatedAt: issuer.CreatedAt,
		UpdatedAt: issuer.UpdatedAt,
		Name:      issuer.Name,
		Country:   issuer.Country.String,
		Website:   issuer.Website.String,
	}

}

func MapIssuerDomainToDto(issuer Issuer) dto.Issuer {

	return dto.Issuer{
		Id:        issuer.Id.String(),
		CreatedAt: issuer.CreatedAt,
		UpdatedAt: issuer.UpdatedAt,
		Name:      issuer.Name,
		Country:   issuer.Country,
		Website:   issuer.Website,
	}

}

func MapIssuerDtoToDomain(issuer dto.Issuer) Issuer {

	id, _ := uuid.Parse(issuer.Id)

	return Issuer{
		Id:        id,
		CreatedAt: issuer.CreatedAt,
		UpdatedAt: issuer.UpdatedAt,
		Name:      issuer.Name,
		Country:   issuer.Country,
		Website:   issuer.Website,
	}

}

func MapIssuerDomainToDb(issuer Issuer) sqlc.Issuer {

	country := sql.NullString{String: issuer.Country, Valid: true}
	website := sql.NullString{String: issuer.Website, Valid: true}

	return sqlc.Issuer{
		ID:        issuer.Id,
		CreatedAt: issuer.CreatedAt,
		UpdatedAt: issuer.UpdatedAt,
		Name:      issuer.Name,
		Country:   country,
		Website:   website,
	}

}
