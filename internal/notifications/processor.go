package notifications

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/email"
	"github.com/adamjames870/seacert/internal/email/email_templates"
	"github.com/google/uuid"
)

var ErrNothingToSend = errors.New("notification has nothing to send")

type SenderFunc func(email email_templates.Email) (string, error)

type Processor struct {
	repo   domain.Repository
	sender SenderFunc
	now    func() time.Time
}

func NewProcessor(repo domain.Repository) *Processor {
	return &Processor{
		repo:   repo,
		sender: email.Send,
		now:    time.Now,
	}
}

func (p *Processor) BuildEmail(
	ctx context.Context,
	notification sqlc.Notification,
) (email_templates.Email, error) {
	if notification.NotificationType != string(TypeNoCertificates7Day) &&
		notification.NotificationType != string(TypeNoCertificates1Month) &&
		notification.NotificationType != string(TypeCertificateExpirySummary) {
		return email_templates.Email{}, fmt.Errorf("unsupported notification type: %s", notification.NotificationType)
	}

	if !notification.UserID.Valid {
		return email_templates.Email{}, errors.New("notification user_id is missing or null")
	}

	user, err := p.repo.GetUserByID(ctx, notification.UserID.UUID)
	if err != nil {
		return email_templates.Email{}, fmt.Errorf("failed to get user: %w", err)
	}

	switch notification.NotificationType {
	case string(TypeNoCertificates7Day), string(TypeNoCertificates1Month):
		data := email_templates.NoCertsEmailData{
			FirstName: user.Forename.String,
			AppURL:    "https://www.seacert.app/certificates",
		}
		return email_templates.GetNoCertsEmail(data, []string{user.Email})

	case string(TypeCertificateExpirySummary):
		certs, err := p.repo.GetCertificatesForExpiryNotification(ctx, notification.UserID.UUID)
		if err != nil {
			return email_templates.Email{}, fmt.Errorf("failed to get certificates for expiry notification: %w", err)
		}

		now := time.Now()
		if p.now != nil {
			now = p.now()
		}

		groups := GroupCertificatesByExpiry(certs, now)
		if len(groups.Expired) == 0 &&
			len(groups.ExpiringOneMonth) == 0 &&
			len(groups.ExpiringSixMonths) == 0 &&
			len(groups.ExpiringOneYear) == 0 {
			return email_templates.Email{}, ErrNothingToSend
		}

		mapItems := func(rows []sqlc.GetCertificatesForExpiryNotificationRow) []email_templates.ExpiringCertificateItem {
			items := make([]email_templates.ExpiringCertificateItem, 0, len(rows))
			for _, r := range rows {
				items = append(items, email_templates.ExpiringCertificateItem{
					Name:       r.CertTypeName,
					ExpiryDate: r.ExpiryDate.Format("02 Jan 2006"),
					Url:        "https://www.seacert.app/certificates",
				})
			}
			return items
		}

		data := email_templates.CertificateExpiryEmailData{
			FirstName:     user.Forename.String,
			AppURL:        "https://www.seacert.app/certificates",
			Expired:       mapItems(groups.Expired),
			NextMonth:     mapItems(groups.ExpiringOneMonth),
			NextSixMonths: mapItems(groups.ExpiringSixMonths),
			NextYear:      mapItems(groups.ExpiringOneYear),
		}

		return email_templates.GetCertificateExpiryEmail(data, []string{user.Email})

	default:
		return email_templates.Email{}, fmt.Errorf("unsupported notification type: %s", notification.NotificationType)
	}
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

func (p *Processor) ProcessNotification(
	ctx context.Context,
	notification sqlc.Notification,
) error {
	if notification.NotificationType != string(TypeNoCertificates7Day) &&
		notification.NotificationType != string(TypeNoCertificates1Month) &&
		notification.NotificationType != string(TypeCertificateExpirySummary) {
		return fmt.Errorf("unsupported notification type: %s", notification.NotificationType)
	}

	processingNotification, err := p.repo.MarkNotificationProcessing(ctx, notification.ID)
	if err != nil {
		return fmt.Errorf("failed to mark notification as processing: %w", err)
	}

	_, sendErr := p.SendDelivery(ctx, processingNotification)
	if sendErr != nil {
		if errors.Is(sendErr, ErrNothingToSend) {
			_, markErr := p.repo.MarkNotificationCompleted(ctx, notification.ID)
			if markErr != nil {
				return fmt.Errorf("nothing to send, but failed to record notification completion state: %w", markErr)
			}
			return nil
		}
		_, markErr := p.repo.MarkNotificationFailed(ctx, notification.ID)
		if markErr != nil {
			return fmt.Errorf("failed to send notification: %w; additionally failed to mark notification as failed: %v", sendErr, markErr)
		}
		return sendErr
	}

	_, err = p.repo.MarkNotificationCompleted(ctx, notification.ID)
	if err != nil {
		return fmt.Errorf("email sent successfully, but failed to record notification completion state: %w", err)
	}

	return nil
}

type BatchResult struct {
	Found     int `json:"found"`
	Completed int `json:"completed"`
	Failed    int `json:"failed"`
}

func (p *Processor) ProcessPendingNotifications(
	ctx context.Context,
	limit int,
) (BatchResult, error) {
	if limit <= 0 {
		return BatchResult{}, nil
	}

	pending, err := p.repo.GetPendingNotifications(ctx)
	if err != nil {
		return BatchResult{}, fmt.Errorf("failed to get pending notifications: %w", err)
	}

	if len(pending) > limit {
		pending = pending[:limit]
	}

	result := BatchResult{
		Found: len(pending),
	}

	for _, notification := range pending {
		if err := p.ProcessNotification(ctx, notification); err != nil {
			slog.ErrorContext(ctx, "failed to process notification in batch",
				"notification_id", notification.ID,
				"notification_type", notification.NotificationType,
				"error", err,
			)
			result.Failed++
		} else {
			result.Completed++
		}
	}

	return result, nil
}
