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

	// Seatime
	GetShipTypes(ctx context.Context) ([]sqlc.ShipType, error)
	GetShipTypeById(ctx context.Context, id uuid.UUID) (sqlc.ShipType, error)
	CreateShipType(ctx context.Context, arg sqlc.CreateShipTypeParams) (sqlc.ShipType, error)
	UpdateShipType(ctx context.Context, arg sqlc.UpdateShipTypeParams) (sqlc.ShipType, error)
	DeleteShipType(ctx context.Context, id uuid.UUID) error

	GetVoyageTypes(ctx context.Context) ([]sqlc.VoyageType, error)
	GetVoyageTypeById(ctx context.Context, id uuid.UUID) (sqlc.VoyageType, error)
	CreateVoyageType(ctx context.Context, arg sqlc.CreateVoyageTypeParams) (sqlc.VoyageType, error)
	UpdateVoyageType(ctx context.Context, arg sqlc.UpdateVoyageTypeParams) (sqlc.VoyageType, error)
	DeleteVoyageType(ctx context.Context, id uuid.UUID) error

	GetSeatimePeriodTypes(ctx context.Context) ([]sqlc.SeatimePeriodType, error)
	GetSeatimePeriodTypeById(ctx context.Context, id uuid.UUID) (sqlc.SeatimePeriodType, error)
	CreateSeatimePeriodType(ctx context.Context, arg sqlc.CreateSeatimePeriodTypeParams) (sqlc.SeatimePeriodType, error)
	UpdateSeatimePeriodType(ctx context.Context, arg sqlc.UpdateSeatimePeriodTypeParams) (sqlc.SeatimePeriodType, error)
	DeleteSeatimePeriodType(ctx context.Context, id uuid.UUID) error
	GetShipByImo(ctx context.Context, imoNumber string) (sqlc.GetShipByImoRow, error)
	GetShipById(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error)
	GetShips(ctx context.Context) ([]sqlc.GetShipsRow, error)
	GetShipsForUser(ctx context.Context, createdBy uuid.NullUUID) ([]sqlc.GetShipsForUserRow, error)
	CreateShip(ctx context.Context, arg sqlc.CreateShipParams) (sqlc.Ship, error)
	UpdateShip(ctx context.Context, arg sqlc.UpdateShipParams) (sqlc.Ship, error)
	UpdateShipStatus(ctx context.Context, arg sqlc.UpdateShipStatusParams) (sqlc.Ship, error)
	UpdateShipReferences(ctx context.Context, arg sqlc.UpdateShipReferencesParams) error
	DeleteShip(ctx context.Context, id uuid.UUID) error
	CreateSeatime(ctx context.Context, arg sqlc.CreateSeatimeParams) (sqlc.Seatime, error)
	CreateSeatimePeriod(ctx context.Context, arg sqlc.CreateSeatimePeriodParams) (sqlc.SeatimePeriod, error)
	UpdateSeatime(ctx context.Context, arg sqlc.UpdateSeatimeParams) (sqlc.Seatime, error)
	DeleteSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) error
	GetSeatimeByUserId(ctx context.Context, userID uuid.UUID) ([]sqlc.GetSeatimeByUserIdRow, error)
	GetSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) ([]sqlc.GetSeatimePeriodsRow, error)
	DeleteSeatime(ctx context.Context, arg sqlc.DeleteSeatimeParams) error

	ResetAll(ctx context.Context) error

	WithTx(ctx context.Context, fn func(Repository) error) error
}
