package notifications

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/email"
	"github.com/adamjames870/seacert/internal/email/email_templates"
	"github.com/google/uuid"
)

type SenderFunc func(email email_templates.Email) (string, error)

type Processor struct {
	repo   domain.Repository
	sender SenderFunc
}

func NewProcessor(repo domain.Repository) *Processor {
	return &Processor{
		repo:   repo,
		sender: email.Send,
	}
}

func (p *Processor) BuildEmail(
	ctx context.Context,
	notification sqlc.Notification,
) (email_templates.Email, error) {
	if notification.NotificationType != string(TypeNoCertificates7Day) {
		return email_templates.Email{}, fmt.Errorf("unsupported notification type: %s", notification.NotificationType)
	}

	if !notification.UserID.Valid {
		return email_templates.Email{}, errors.New("notification user_id is missing or null")
	}

	user, err := p.repo.GetUserByID(ctx, notification.UserID.UUID)
	if err != nil {
		return email_templates.Email{}, fmt.Errorf("failed to get user: %w", err)
	}

	data := email_templates.NoCertsEmailData{
		FirstName: user.Forename.String,
		AppURL:    "https://www.seacert.app/certificates",
	}

	return email_templates.GetNoCertsEmail(data, []string{user.Email})
}

func (p *Processor) CreateDelivery(
	ctx context.Context,
	notification sqlc.Notification,
) (sqlc.EmailDelivery, email_templates.Email, error) {
	email, err := p.BuildEmail(ctx, notification)
	if err != nil {
		return sqlc.EmailDelivery{}, email_templates.Email{}, fmt.Errorf("failed to build email: %w", err)
	}

	if len(email.Recipient) != 1 {
		return sqlc.EmailDelivery{}, email_templates.Email{}, fmt.Errorf("expected exactly 1 recipient, got %d", len(email.Recipient))
	}

	deliveries, err := p.repo.GetEmailDeliveriesForNotification(ctx, notification.ID)
	if err != nil {
		return sqlc.EmailDelivery{}, email_templates.Email{}, fmt.Errorf("failed to get email deliveries: %w", err)
	}

	var maxAttempt int32 = 0
	for _, d := range deliveries {
		if d.Attempt > maxAttempt {
			maxAttempt = d.Attempt
		}
	}
	nextAttempt := maxAttempt + 1

	deliveryParams := sqlc.CreateEmailDeliveryParams{
		ID:             uuid.New(),
		NotificationID: notification.ID,
		Recipient:      email.Recipient[0],
		Provider:       ProviderResend,
		Attempt:        nextAttempt,
	}

	delivery, err := p.repo.CreateEmailDelivery(ctx, deliveryParams)
	if err != nil {
		return sqlc.EmailDelivery{}, email_templates.Email{}, fmt.Errorf("failed to create email delivery: %w", err)
	}

	return delivery, email, nil
}

func (p *Processor) SendDelivery(
	ctx context.Context,
	notification sqlc.Notification,
) (sqlc.EmailDelivery, error) {
	delivery, builtEmail, err := p.CreateDelivery(ctx, notification)
	if err != nil {
		return sqlc.EmailDelivery{}, fmt.Errorf("failed to create delivery: %w", err)
	}

	sender := p.sender
	if sender == nil {
		sender = email.Send
	}

	providerMessageID, sendErr := sender(builtEmail)
	if sendErr != nil {
		failedDelivery, markErr := p.repo.MarkEmailDeliveryFailed(ctx, sqlc.MarkEmailDeliveryFailedParams{
			ID: delivery.ID,
			ErrorMessage: sql.NullString{
				String: sendErr.Error(),
				Valid:  true,
			},
		})
		if markErr != nil {
			return delivery, fmt.Errorf("failed to send email: %w; additionally failed to mark delivery as failed: %v", sendErr, markErr)
		}
		return failedDelivery, fmt.Errorf("failed to send email: %w", sendErr)
	}

	sentDelivery, markErr := p.repo.MarkEmailDeliverySent(ctx, sqlc.MarkEmailDeliverySentParams{
		ID: delivery.ID,
		ProviderMessageID: sql.NullString{
			String: providerMessageID,
			Valid:  providerMessageID != "",
		},
	})
	if markErr != nil {
		return delivery, fmt.Errorf("email sent successfully with provider ID %q, but failed to record sent delivery state: %w", providerMessageID, markErr)
	}

	return sentDelivery, nil
}
