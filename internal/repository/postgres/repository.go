package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/google/uuid"
)

type repository struct {
	db *sql.DB
	*sqlc.Queries
}

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (r *repository) ResetAll(ctx context.Context) error {
	return r.WithTx(ctx, func(txRepo domain.Repository) error {
		_ = r.Queries.ResetEmailDeliveries(ctx)
		_ = r.Queries.ResetNotifications(ctx)
		_ = r.Queries.ResetSeatimePeriods(ctx)
		_ = r.Queries.ResetSeatime(ctx)
		_ = r.Queries.ResetShips(ctx)
		_ = r.Queries.ResetSuccessions(ctx)
		_ = r.Queries.ResetCerts(ctx)
		_ = r.Queries.ResetCertTypes(ctx)
		_ = r.Queries.ResetIssuers(ctx)
		_ = r.Queries.ResetUsers(ctx)
		_ = r.Queries.DeleteExpiredPromptCaches(ctx)
		return nil
	})
}

func (r *repository) CreateNotification(ctx context.Context, arg sqlc.CreateNotificationParams) (sqlc.Notification, error) {
	return r.Queries.CreateNotification(ctx, arg)
}

func (r *repository) GetUsersEligibleForNoCertificates7Day(ctx context.Context) ([]uuid.UUID, error) {
	return r.Queries.GetUsersEligibleForNoCertificates7Day(ctx)
}

func (r *repository) GetUsersEligibleForNoCertificates1Month(ctx context.Context) ([]uuid.UUID, error) {
	return r.Queries.GetUsersEligibleForNoCertificates1Month(ctx)
}

func (r *repository) GetCandidateUsersForExpiryNotification(ctx context.Context) ([]uuid.UUID, error) {
	return r.Queries.GetCandidateUsersForExpiryNotification(ctx)
}

func (r *repository) GetNotificationByKey(ctx context.Context, notificationKey string) (sqlc.Notification, error) {
	return r.Queries.GetNotificationByKey(ctx, notificationKey)
}

func (r *repository) GetPendingNotifications(ctx context.Context) ([]sqlc.Notification, error) {
	return r.Queries.GetPendingNotifications(ctx)
}

func (r *repository) MarkNotificationProcessing(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
	return r.Queries.MarkNotificationProcessing(ctx, id)
}

func (r *repository) MarkNotificationCompleted(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
	return r.Queries.MarkNotificationCompleted(ctx, id)
}

func (r *repository) MarkNotificationFailed(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
	return r.Queries.MarkNotificationFailed(ctx, id)
}

func (r *repository) CreateEmailDelivery(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
	return r.Queries.CreateEmailDelivery(ctx, arg)
}

func (r *repository) MarkEmailDeliverySent(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
	return r.Queries.MarkEmailDeliverySent(ctx, arg)
}

func (r *repository) MarkEmailDeliveryFailed(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error) {
	return r.Queries.MarkEmailDeliveryFailed(ctx, arg)
}

func (r *repository) GetEmailDeliveriesForNotification(ctx context.Context, notificationID uuid.UUID) ([]sqlc.EmailDelivery, error) {
	return r.Queries.GetEmailDeliveriesForNotification(ctx, notificationID)
}

func (r *repository) WithTx(ctx context.Context, fn func(domain.Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := r.Queries.WithTx(tx)
	err = fn(&repository{
		db:      r.db,
		Queries: qtx,
	})

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
