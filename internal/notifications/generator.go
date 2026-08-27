package notifications

import (
	"context"
	"encoding/json"
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
