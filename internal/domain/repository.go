package domain

import (
	"context"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/google/uuid"
)

// Repository defines the interface for database operations needed by the domain.
type Repository interface {
	// Certificates
	GetCerts(ctx context.Context, userID uuid.UUID) ([]sqlc.GetCertsRow, error)
	GetCertFromId(ctx context.Context, arg sqlc.GetCertFromIdParams) (sqlc.GetCertFromIdRow, error)
	GetPredecessors(ctx context.Context, newCert uuid.UUID) ([]sqlc.GetPredecessorsRow, error)
	CreateCert(ctx context.Context, arg sqlc.CreateCertParams) (sqlc.Certificate, error)
	UpdateCertificate(ctx context.Context, arg sqlc.UpdateCertificateParams) (sqlc.Certificate, error)
	DeleteCert(ctx context.Context, arg sqlc.DeleteCertParams) error
	CreateSuccession(ctx context.Context, arg sqlc.CreateSuccessionParams) (sqlc.Succession, error)

	// Certificate Types
	GetCertTypes(ctx context.Context) ([]sqlc.CertificateType, error)
	GetCertTypesForUser(ctx context.Context, createdBy uuid.NullUUID) ([]sqlc.CertificateType, error)
	GetCertTypeFromId(ctx context.Context, id uuid.UUID) (sqlc.CertificateType, error)
	GetCertTypeFromName(ctx context.Context, name string) (sqlc.CertificateType, error)
	CreateCertType(ctx context.Context, arg sqlc.CreateCertTypeParams) (sqlc.CertificateType, error)
	UpdateCertType(ctx context.Context, arg sqlc.UpdateCertTypeParams) (sqlc.CertificateType, error)
	UpdateCertTypeReferences(ctx context.Context, arg sqlc.UpdateCertTypeReferencesParams) error
	DeleteCertType(ctx context.Context, id uuid.UUID) error

	// Issuers
	GetIssuers(ctx context.Context) ([]sqlc.Issuer, error)
	GetIssuerById(ctx context.Context, id uuid.UUID) (sqlc.Issuer, error)
	GetIssuerByName(ctx context.Context, name string) (sqlc.Issuer, error)
	CreateIssuer(ctx context.Context, arg sqlc.CreateIssuerParams) (sqlc.Issuer, error)
	UpdateIssuer(ctx context.Context, arg sqlc.UpdateIssuerParams) (sqlc.Issuer, error)

	// Users
	GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	GetUserByEmail(ctx context.Context, email string) (sqlc.User, error)
	CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error)
	UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error)

	ResetAll(ctx context.Context) error

	WithTx(ctx context.Context, fn func(Repository) error) error
}
