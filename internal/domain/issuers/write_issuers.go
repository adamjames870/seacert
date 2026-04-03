package issuers

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

func GetIssuers(ctx context.Context, repo domain.Repository) ([]Issuer, error) {
	issuers, err := repo.GetIssuers(ctx)
	if err != nil {
		return nil, err
	}

	apiIssuers := make([]Issuer, 0, len(issuers))
	for _, issuer := range issuers {
		apiIssuers = append(apiIssuers, MapIssuerDbToDomain(issuer))
	}

	return apiIssuers, nil
}

func GetIssuerById(ctx context.Context, repo domain.Repository, id uuid.UUID) (Issuer, error) {
	issuer, err := repo.GetIssuerById(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issuer{}, domain.ErrNotFound
		}
		return Issuer{}, err
	}

	return MapIssuerDbToDomain(issuer), nil
}

func GetIssuerByName(ctx context.Context, repo domain.Repository, name string) (Issuer, error) {
	issuer, err := repo.GetIssuerByName(ctx, name)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issuer{}, domain.ErrNotFound
		}
		return Issuer{}, err
	}

	return MapIssuerDbToDomain(issuer), nil
}

func CreateIssuer(ctx context.Context, repo domain.Repository, params dto.ParamsAddIssuer) (Issuer, error) {
	newIssuer := sqlc.CreateIssuerParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      params.Name,
		Country:   domain.ToNullStringFromPointer(params.Country),
		Website:   domain.ToNullStringFromPointer(params.Website),
	}

	dbIssuer, err := repo.CreateIssuer(ctx, newIssuer)
	if err != nil {
		return Issuer{}, err
	}

	return MapIssuerDbToDomain(dbIssuer), nil
}

func UpdateIssuer(ctx context.Context, repo domain.Repository, params dto.ParamsUpdateIssuer) (Issuer, error) {
	uuidId, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return Issuer{}, domain.ErrInvalidInput
	}

	updateIssuer := sqlc.UpdateIssuerParams{
		ID:      uuidId,
		Name:    domain.ToNullStringFromPointer(params.Name),
		Country: domain.ToNullStringFromPointer(params.Country),
		Website: domain.ToNullStringFromPointer(params.Website),
	}

	dbIssuer, err := repo.UpdateIssuer(ctx, updateIssuer)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Issuer{}, domain.ErrNotFound
		}
		return Issuer{}, err
	}

	return MapIssuerDbToDomain(dbIssuer), nil
}
