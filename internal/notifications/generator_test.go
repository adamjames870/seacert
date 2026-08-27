package notifications

import (
	"context"
	"testing"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/google/uuid"
)

type mockRepository struct {
	CreateNotificationFunc                    func(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error)
	GetUsersEligibleForNoCertificates7DayFunc func(ctx context.Context) ([]uuid.UUID, error)
	GetUserByIDFunc                           func(ctx context.Context, id uuid.UUID) (sqlc.User, error)
	CreateEmailDeliveryFunc                   func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error)
	GetEmailDeliveriesForNotificationFunc     func(ctx context.Context, notificationID uuid.UUID) ([]sqlc.EmailDelivery, error)
	MarkEmailDeliverySentFunc                 func(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error)
	MarkEmailDeliveryFailedFunc               func(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error)
}

func (m *mockRepository) CreateNotification(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error) {
	return m.CreateNotificationFunc(ctx, arg)
}

func (m *mockRepository) GetUsersEligibleForNoCertificates7Day(ctx context.Context) ([]uuid.UUID, error) {
	return m.GetUsersEligibleForNoCertificates7DayFunc(ctx)
}

// Panic for unimplemented methods
func (m *mockRepository) GetPendingNotifications(ctx context.Context) ([]sqlc.Notification, error) {
	panic("not implemented")
}
func (m *mockRepository) MarkNotificationProcessing(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
	panic("not implemented")
}
func (m *mockRepository) MarkNotificationCompleted(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
	panic("not implemented")
}
func (m *mockRepository) MarkNotificationFailed(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
	panic("not implemented")
}
func (m *mockRepository) GetCerts(ctx context.Context, userID uuid.UUID) ([]sqlc.GetCertsRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetCertFromId(ctx context.Context, arg sqlc.GetCertFromIdParams) (sqlc.GetCertFromIdRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetPredecessors(ctx context.Context, newCert uuid.UUID) ([]sqlc.GetPredecessorsRow, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateCert(ctx context.Context, arg sqlc.CreateCertParams) (sqlc.Certificate, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateCertificate(ctx context.Context, arg sqlc.UpdateCertificateParams) (sqlc.Certificate, error) {
	panic("not implemented")
}
func (m *mockRepository) DeleteCert(ctx context.Context, arg sqlc.DeleteCertParams) error {
	panic("not implemented")
}
func (m *mockRepository) CreateSuccession(ctx context.Context, arg sqlc.CreateSuccessionParams) (sqlc.Succession, error) {
	panic("not implemented")
}
func (m *mockRepository) GetCertTypes(ctx context.Context) ([]sqlc.CertificateType, error) {
	panic("not implemented")
}
func (m *mockRepository) GetCertTypesForUser(ctx context.Context, createdBy uuid.NullUUID) ([]sqlc.CertificateType, error) {
	panic("not implemented")
}
func (m *mockRepository) GetCertTypeFromId(ctx context.Context, id uuid.UUID) (sqlc.CertificateType, error) {
	panic("not implemented")
}
func (m *mockRepository) GetCertTypeFromName(ctx context.Context, name string) (sqlc.CertificateType, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateCertType(ctx context.Context, arg sqlc.CreateCertTypeParams) (sqlc.CertificateType, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateCertType(ctx context.Context, arg sqlc.UpdateCertTypeParams) (sqlc.CertificateType, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateCertTypeReferences(ctx context.Context, arg sqlc.UpdateCertTypeReferencesParams) error {
	panic("not implemented")
}
func (m *mockRepository) DeleteCertType(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
func (m *mockRepository) GetIssuers(ctx context.Context) ([]sqlc.Issuer, error) {
	panic("not implemented")
}
func (m *mockRepository) GetIssuerById(ctx context.Context, id uuid.UUID) (sqlc.Issuer, error) {
	panic("not implemented")
}
func (m *mockRepository) GetIssuerByName(ctx context.Context, name string) (sqlc.Issuer, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateIssuer(ctx context.Context, arg sqlc.CreateIssuerParams) (sqlc.Issuer, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateIssuer(ctx context.Context, arg sqlc.UpdateIssuerParams) (sqlc.Issuer, error) {
	panic("not implemented")
}
func (m *mockRepository) GetUserByID(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
	if m.GetUserByIDFunc != nil {
		return m.GetUserByIDFunc(ctx, id)
	}
	panic("not implemented")
}
func (m *mockRepository) GetUserByEmail(ctx context.Context, email string) (sqlc.User, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateUser(ctx context.Context, arg sqlc.CreateUserParams) (sqlc.User, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateUser(ctx context.Context, arg sqlc.UpdateUserParams) (sqlc.User, error) {
	panic("not implemented")
}
func (m *mockRepository) GetShipTypes(ctx context.Context) ([]sqlc.ShipType, error) {
	panic("not implemented")
}
func (m *mockRepository) GetShipTypeById(ctx context.Context, id uuid.UUID) (sqlc.ShipType, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateShipType(ctx context.Context, arg sqlc.CreateShipTypeParams) (sqlc.ShipType, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateShipType(ctx context.Context, arg sqlc.UpdateShipTypeParams) (sqlc.ShipType, error) {
	panic("not implemented")
}
func (m *mockRepository) DeleteShipType(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
func (m *mockRepository) GetVoyageTypes(ctx context.Context) ([]sqlc.VoyageType, error) {
	panic("not implemented")
}
func (m *mockRepository) GetVoyageTypeById(ctx context.Context, id uuid.UUID) (sqlc.VoyageType, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateVoyageType(ctx context.Context, arg sqlc.CreateVoyageTypeParams) (sqlc.VoyageType, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateVoyageType(ctx context.Context, arg sqlc.UpdateVoyageTypeParams) (sqlc.VoyageType, error) {
	panic("not implemented")
}
func (m *mockRepository) DeleteVoyageType(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
func (m *mockRepository) GetSeatimePeriodTypes(ctx context.Context) ([]sqlc.SeatimePeriodType, error) {
	panic("not implemented")
}
func (m *mockRepository) GetSeatimePeriodTypeById(ctx context.Context, id uuid.UUID) (sqlc.SeatimePeriodType, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateSeatimePeriodType(ctx context.Context, arg sqlc.CreateSeatimePeriodTypeParams) (sqlc.SeatimePeriodType, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateSeatimePeriodType(ctx context.Context, arg sqlc.UpdateSeatimePeriodTypeParams) (sqlc.SeatimePeriodType, error) {
	panic("not implemented")
}
func (m *mockRepository) DeleteSeatimePeriodType(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
func (m *mockRepository) GetShipByImo(ctx context.Context, imoNumber string) (sqlc.GetShipByImoRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetShipById(ctx context.Context, id uuid.UUID) (sqlc.GetShipByIdRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetShips(ctx context.Context) ([]sqlc.GetShipsRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetShipsForUser(ctx context.Context, createdBy uuid.NullUUID) ([]sqlc.GetShipsForUserRow, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateShip(ctx context.Context, arg sqlc.CreateShipParams) (sqlc.Ship, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateShip(ctx context.Context, arg sqlc.UpdateShipParams) (sqlc.Ship, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateShipStatus(ctx context.Context, arg sqlc.UpdateShipStatusParams) (sqlc.Ship, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateShipReferences(ctx context.Context, arg sqlc.UpdateShipReferencesParams) error {
	panic("not implemented")
}
func (m *mockRepository) DeleteShip(ctx context.Context, id uuid.UUID) error {
	panic("not implemented")
}
func (m *mockRepository) CreateSeatime(ctx context.Context, arg sqlc.CreateSeatimeParams) (sqlc.Seatime, error) {
	panic("not implemented")
}
func (m *mockRepository) CreateSeatimePeriod(ctx context.Context, arg sqlc.CreateSeatimePeriodParams) (sqlc.SeatimePeriod, error) {
	panic("not implemented")
}
func (m *mockRepository) UpdateSeatime(ctx context.Context, arg sqlc.UpdateSeatimeParams) (sqlc.Seatime, error) {
	panic("not implemented")
}
func (m *mockRepository) DeleteSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) error {
	panic("not implemented")
}
func (m *mockRepository) GetSeatimeByUserId(ctx context.Context, userID uuid.UUID) ([]sqlc.GetSeatimeByUserIdRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetSeatimePeriods(ctx context.Context, seatimeID uuid.UUID) ([]sqlc.GetSeatimePeriodsRow, error) {
	panic("not implemented")
}
func (m *mockRepository) GetOverlappingSeatime(ctx context.Context, arg sqlc.GetOverlappingSeatimeParams) ([]sqlc.Seatime, error) {
	panic("not implemented")
}
func (m *mockRepository) DeleteSeatime(ctx context.Context, arg sqlc.DeleteSeatimeParams) error {
	panic("not implemented")
}
func (m *mockRepository) ResetAll(ctx context.Context) error { panic("not implemented") }
func (m *mockRepository) GetPromptCache(ctx context.Context, cacheKey string) (sqlc.PromptCach, error) {
	panic("not implemented")
}
func (m *mockRepository) UpsertPromptCache(ctx context.Context, arg sqlc.UpsertPromptCacheParams) error {
	panic("not implemented")
}
func (m *mockRepository) DeleteExpiredPromptCaches(ctx context.Context) error {
	panic("not implemented")
}
func (m *mockRepository) CreateEmailDelivery(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
	if m.CreateEmailDeliveryFunc != nil {
		return m.CreateEmailDeliveryFunc(ctx, arg)
	}
	panic("not implemented")
}
func (m *mockRepository) MarkEmailDeliverySent(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
	if m.MarkEmailDeliverySentFunc != nil {
		return m.MarkEmailDeliverySentFunc(ctx, arg)
	}
	panic("not implemented")
}
func (m *mockRepository) MarkEmailDeliveryFailed(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error) {
	if m.MarkEmailDeliveryFailedFunc != nil {
		return m.MarkEmailDeliveryFailedFunc(ctx, arg)
	}
	panic("not implemented")
}
func (m *mockRepository) GetEmailDeliveriesForNotification(ctx context.Context, notificationID uuid.UUID) ([]sqlc.EmailDelivery, error) {
	if m.GetEmailDeliveriesForNotificationFunc != nil {
		return m.GetEmailDeliveriesForNotificationFunc(ctx, notificationID)
	}
	panic("not implemented")
}
func (m *mockRepository) WithTx(ctx context.Context, fn func(domain.Repository) error) error {
	panic("not implemented")
}

func TestGenerateNoCertificates7Day(t *testing.T) {
	repo := new(mockRepository)
	generator := NewGenerator(repo)

	user1 := uuid.New()
	user2 := uuid.New()
	eligibleUsers := []uuid.UUID{user1, user2}

	repo.GetUsersEligibleForNoCertificates7DayFunc = func(ctx context.Context) ([]uuid.UUID, error) {
		return eligibleUsers, nil
	}

	repo.CreateNotificationFunc = func(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error) {
		return sqlc.Notification{}, nil
	}

	count, err := generator.GenerateNoCertificates7Day(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if count != 2 {
		t.Errorf("expected 2 notifications, got %d", count)
	}
}
