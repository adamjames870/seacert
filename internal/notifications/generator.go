package notifications

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/google/uuid"
)

type Generator struct {
	repo domain.Repository
}

func NewGenerator(repo domain.Repository) *Generator {
	return &Generator{repo: repo}
}

func (g *Generator) GenerateNoCertificates7Day(ctx context.Context) (int, error) {
	eligibleUserIDs, err := g.repo.GetUsersEligibleForNoCertificates7Day(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get eligible users: %w", err)
	}

	count := 0
	for _, userID := range eligibleUserIDs {
		notificationKey := fmt.Sprintf("no-certificates:%s:7d", userID.String())

		_, err := g.repo.CreateNotification(ctx, sqlc.CreateNotificationParams{
			ID:               uuid.New(),
			UserID:           uuid.NullUUID{UUID: userID, Valid: true},
			NotificationType: string(TypeNoCertificates7Day),
			NotificationKey:  notificationKey,
			Payload:          json.RawMessage("{}"),
			ScheduledAt:      time.Now(),
		})
		if err != nil {
			return count, fmt.Errorf("failed to create notification for user %s: %w", userID, err)
		}
		count++
	}

	return count, nil
}

func (g *Generator) GenerateNoCertificates1Month(ctx context.Context) (int, error) {
	eligibleUserIDs, err := g.repo.GetUsersEligibleForNoCertificates1Month(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get eligible users: %w", err)
	}

	count := 0
	for _, userID := range eligibleUserIDs {
		notificationKey := fmt.Sprintf("no-certificates:%s:1m", userID.String())

		_, err := g.repo.CreateNotification(ctx, sqlc.CreateNotificationParams{
			ID:               uuid.New(),
			UserID:           uuid.NullUUID{UUID: userID, Valid: true},
			NotificationType: string(TypeNoCertificates1Month),
			NotificationKey:  notificationKey,
			Payload:          json.RawMessage("{}"),
			ScheduledAt:      time.Now(),
		})
		if err != nil {
			return count, fmt.Errorf("failed to create notification for user %s: %w", userID, err)
		}
		count++
	}

	return count, nil
}

func (g *Generator) GenerateCertificateExpirySummary(ctx context.Context, now time.Time) (int, error) {
	candidateUserIDs, err := g.repo.GetCandidateUsersForExpiryNotification(ctx)
	if err != nil {
		return 0, fmt.Errorf("failed to get candidate users: %w", err)
	}

	monthStr := now.Format("2006-01")
	count := 0

	for _, userID := range candidateUserIDs {
		notificationKey := fmt.Sprintf("certificate-expiry-summary:%s:%s", userID.String(), monthStr)

		// Check if notification for this key already exists
		_, err := g.repo.GetNotificationByKey(ctx, notificationKey)
		if err == nil {
			// Notification already exists (pending, processing, completed, or failed), skip
			continue
		}
		if !errors.Is(err, sql.ErrNoRows) {
			return count, fmt.Errorf("failed to check existing notification for user %s: %w", userID, err)
		}

		// Retrieve certificates for user
		certs, err := g.repo.GetCertificatesForExpiryNotification(ctx, userID)
		if err != nil {
			return count, fmt.Errorf("failed to get certificates for expiry notification for user %s: %w", userID, err)
		}

		groups := GroupCertificatesByExpiry(certs, now)
		if len(groups.Expired) == 0 &&
			len(groups.ExpiringOneMonth) == 0 &&
			len(groups.ExpiringSixMonths) == 0 &&
			len(groups.ExpiringOneYear) == 0 {
			continue
		}

		_, err = g.repo.CreateNotification(ctx, sqlc.CreateNotificationParams{
			ID:               uuid.New(),
			UserID:           uuid.NullUUID{UUID: userID, Valid: true},
			NotificationType: string(TypeCertificateExpirySummary),
			NotificationKey:  notificationKey,
			Payload:          json.RawMessage("{}"),
			ScheduledAt:      now,
		})
		if err != nil {
			return count, fmt.Errorf("failed to create notification for user %s: %w", userID, err)
		}
		count++
	}

	return count, nil
}
