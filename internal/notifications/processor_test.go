package notifications

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/email/email_templates"
	"github.com/google/uuid"
)

func TestProcessor_BuildEmail_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			if id != userID {
				t.Fatalf("expected userID %s, got %s", userID, id)
			}
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Captain", Valid: true},
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	email, err := processor.BuildEmail(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(email.Recipient) != 1 || email.Recipient[0] != "sailor@example.com" {
		t.Errorf("expected recipient sailor@example.com, got %v", email.Recipient)
	}

	if email.Subject == "" {
		t.Errorf("expected non-empty subject")
	}

	if email.HtmlBody == "" {
		t.Errorf("expected non-empty HtmlBody")
	}

	if email.TextBody == "" {
		t.Errorf("expected non-empty TextBody")
	}

	if !strings.Contains(email.HtmlBody, "Captain") {
		t.Errorf("expected HtmlBody to contain user's name, got: %s", email.HtmlBody)
	}

	if !strings.Contains(email.TextBody, "Captain") {
		t.Errorf("expected TextBody to contain user's name, got: %s", email.TextBody)
	}
}

func TestProcessor_BuildEmail_NullUserID(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockRepository{}
	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{Valid: false},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, err := processor.BuildEmail(ctx, notification)
	if err == nil {
		t.Fatal("expected error for missing user_id, got nil")
	}
}

func TestProcessor_BuildEmail_UnsupportedNotificationType(t *testing.T) {
	ctx := context.Background()
	mockRepo := &mockRepository{}
	processor := NewProcessor(mockRepo)

	testCases := []string{
		string(TypeNoCertificates1Month),
		string(TypeCertificateExpirySummary),
		"unsupported_type",
		"",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			notification := sqlc.Notification{
				ID:               uuid.New(),
				UserID:           uuid.NullUUID{UUID: uuid.New(), Valid: true},
				NotificationType: tc,
				Status:           "pending",
			}

			_, err := processor.BuildEmail(ctx, notification)
			if err == nil {
				t.Fatalf("expected error for unsupported notification type %q, got nil", tc)
			}
		})
	}
}

func TestProcessor_BuildEmail_UserRepoError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{}, errors.New("database connection failure")
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, err := processor.BuildEmail(ctx, notification)
	if err == nil {
		t.Fatal("expected error when repo fails, got nil")
	}
}

func TestProcessor_BuildEmail_UserWithNullForename(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "test@example.com",
				Forename: sql.NullString{Valid: false},
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	email, err := processor.BuildEmail(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if len(email.Recipient) != 1 || email.Recipient[0] != "test@example.com" {
		t.Errorf("expected recipient test@example.com, got %v", email.Recipient)
	}

	if email.Subject == "" || email.HtmlBody == "" || email.TextBody == "" {
		t.Errorf("expected subject, HtmlBody and TextBody to be non-empty")
	}
}

func TestProcessor_CreateDelivery_FirstAttempt(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()

	var createdParams sqlc.CreateEmailDeliveryParams
	var getDeliveriesCalled bool

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "first_attempt@example.com",
				Forename: sql.NullString{String: "First", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			if notifID != notificationID {
				t.Fatalf("expected notificationID %s, got %s", notificationID, notifID)
			}
			getDeliveriesCalled = true
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			createdParams = arg
			return sqlc.EmailDelivery{
				ID:             arg.ID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	delivery, email, err := processor.CreateDelivery(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !getDeliveriesCalled {
		t.Errorf("expected GetEmailDeliveriesForNotification to be called")
	}

	if createdParams.ID == uuid.Nil {
		t.Errorf("expected new non-nil UUID for delivery ID")
	}

	if createdParams.NotificationID != notificationID {
		t.Errorf("expected NotificationID %s, got %s", notificationID, createdParams.NotificationID)
	}

	if createdParams.Recipient != "first_attempt@example.com" {
		t.Errorf("expected recipient first_attempt@example.com, got %s", createdParams.Recipient)
	}

	if createdParams.Provider != ProviderResend {
		t.Errorf("expected provider %s, got %s", ProviderResend, createdParams.Provider)
	}

	if createdParams.Attempt != 1 {
		t.Errorf("expected attempt 1 for first delivery, got %d", createdParams.Attempt)
	}

	if delivery.Status != "pending" {
		t.Errorf("expected delivery status pending, got %s", delivery.Status)
	}

	if len(email.Recipient) != 1 || email.Recipient[0] != "first_attempt@example.com" {
		t.Errorf("expected email recipient first_attempt@example.com, got %v", email.Recipient)
	}
}

func TestProcessor_CreateDelivery_SubsequentAttempt(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()

	var createdParams sqlc.CreateEmailDeliveryParams

	// Existing deliveries unordered: attempt 2, attempt 1, attempt 4, attempt 3
	existingDeliveries := []sqlc.EmailDelivery{
		{ID: uuid.New(), NotificationID: notificationID, Recipient: "subsequent@example.com", Provider: "resend", Status: "failed", Attempt: 2},
		{ID: uuid.New(), NotificationID: notificationID, Recipient: "subsequent@example.com", Provider: "resend", Status: "failed", Attempt: 1},
		{ID: uuid.New(), NotificationID: notificationID, Recipient: "subsequent@example.com", Provider: "resend", Status: "failed", Attempt: 4},
		{ID: uuid.New(), NotificationID: notificationID, Recipient: "subsequent@example.com", Provider: "resend", Status: "failed", Attempt: 3},
	}

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "subsequent@example.com",
				Forename: sql.NullString{String: "Subsequent", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return existingDeliveries, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			createdParams = arg
			return sqlc.EmailDelivery{
				ID:             arg.ID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	delivery, _, err := processor.CreateDelivery(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if createdParams.Attempt != 5 {
		t.Errorf("expected attempt 5 (max 4 + 1), got %d", createdParams.Attempt)
	}

	if delivery.Attempt != 5 {
		t.Errorf("expected returned delivery attempt 5, got %d", delivery.Attempt)
	}
}

func TestProcessor_CreateDelivery_BuildEmailErrorPropagates(t *testing.T) {
	ctx := context.Background()

	mockRepo := &mockRepository{
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			t.Fatal("GetEmailDeliveriesForNotification should not be called when BuildEmail fails")
			return nil, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			t.Fatal("CreateEmailDelivery should not be called when BuildEmail fails")
			return sqlc.EmailDelivery{}, nil
		},
	}

	processor := NewProcessor(mockRepo)

	// Unsupported type
	notification := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{UUID: uuid.New(), Valid: true},
		NotificationType: "unsupported",
		Status:           "pending",
	}

	_, _, err := processor.CreateDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error for unsupported notification type, got nil")
	}

	// Missing user_id
	notificationNullUser := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{Valid: false},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, _, err = processor.CreateDelivery(ctx, notificationNullUser)
	if err == nil {
		t.Fatal("expected error for null user_id, got nil")
	}
}

func TestProcessor_CreateDelivery_GetDeliveriesError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "err@example.com",
				Forename: sql.NullString{String: "Err", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return nil, errors.New("failed to query delivery attempts")
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			t.Fatal("CreateEmailDelivery should not be called when GetEmailDeliveriesForNotification fails")
			return sqlc.EmailDelivery{}, nil
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, _, err := processor.CreateDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error when GetEmailDeliveriesForNotification fails, got nil")
	}
}

func TestProcessor_CreateDelivery_CreateDeliveryError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "err@example.com",
				Forename: sql.NullString{String: "Err", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{}, errors.New("insert failed")
		},
	}

	processor := NewProcessor(mockRepo)

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, _, err := processor.CreateDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error when CreateEmailDelivery fails, got nil")
	}
}

func TestProcessor_SendDelivery_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	var senderCalled bool
	var markSentParams sqlc.MarkEmailDeliverySentParams

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Popeye", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliverySentFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
			markSentParams = arg
			return sqlc.EmailDelivery{
				ID:                arg.ID,
				NotificationID:    notificationID,
				Recipient:         "sailor@example.com",
				Provider:          ProviderResend,
				ProviderMessageID: arg.ProviderMessageID,
				Status:            "sent",
				Attempt:           1,
				SentAt:            sql.NullTime{Time: time.Now(), Valid: true},
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		senderCalled = true
		if len(email.Recipient) != 1 || email.Recipient[0] != "sailor@example.com" {
			t.Errorf("unexpected email recipient: %v", email.Recipient)
		}
		return "msg_12345", nil
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	delivery, err := processor.SendDelivery(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !senderCalled {
		t.Errorf("expected email sender to be called")
	}

	if markSentParams.ID != deliveryID {
		t.Errorf("expected markSent delivery ID %s, got %s", deliveryID, markSentParams.ID)
	}

	if !markSentParams.ProviderMessageID.Valid || markSentParams.ProviderMessageID.String != "msg_12345" {
		t.Errorf("expected provider message ID msg_12345, got %v", markSentParams.ProviderMessageID)
	}

	if delivery.Status != "sent" {
		t.Errorf("expected delivery status sent, got %s", delivery.Status)
	}

	if !delivery.SentAt.Valid {
		t.Errorf("expected SentAt to be valid")
	}

	if delivery.ProviderMessageID.String != "msg_12345" {
		t.Errorf("expected delivery ProviderMessageID msg_12345, got %s", delivery.ProviderMessageID.String)
	}
}

func TestProcessor_SendDelivery_SendFailure(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	var markFailedParams sqlc.MarkEmailDeliveryFailedParams
	createDeliveryCalls := 0

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Popeye", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			createDeliveryCalls++
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliveryFailedFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error) {
			markFailedParams = arg
			return sqlc.EmailDelivery{
				ID:             arg.ID,
				NotificationID: notificationID,
				Recipient:      "sailor@example.com",
				Provider:       ProviderResend,
				Status:         "failed",
				Attempt:        1,
				ErrorMessage:   arg.ErrorMessage,
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	resendErr := errors.New("resend API rate limit exceeded")
	processor.sender = func(email email_templates.Email) (string, error) {
		return "", resendErr
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	delivery, err := processor.SendDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error from SendDelivery on send failure, got nil")
	}

	if !errors.Is(err, resendErr) {
		t.Errorf("expected wrapped resendErr, got %v", err)
	}

	if createDeliveryCalls != 1 {
		t.Errorf("expected exactly 1 delivery creation (no retry attempt), got %d", createDeliveryCalls)
	}

	if markFailedParams.ID != deliveryID {
		t.Errorf("expected markFailed delivery ID %s, got %s", deliveryID, markFailedParams.ID)
	}

	if !markFailedParams.ErrorMessage.Valid || markFailedParams.ErrorMessage.String != resendErr.Error() {
		t.Errorf("expected error message %q, got %v", resendErr.Error(), markFailedParams.ErrorMessage)
	}

	if delivery.Status != "failed" {
		t.Errorf("expected delivery status failed, got %s", delivery.Status)
	}
}

func TestProcessor_SendDelivery_CreateDeliveryError(t *testing.T) {
	ctx := context.Background()
	var senderCalled bool

	mockRepo := &mockRepository{
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return nil, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		senderCalled = true
		return "msg_123", nil
	}

	// Invalid notification (unsupported type)
	notification := sqlc.Notification{
		ID:               uuid.New(),
		UserID:           uuid.NullUUID{UUID: uuid.New(), Valid: true},
		NotificationType: "invalid_type",
		Status:           "pending",
	}

	_, err := processor.SendDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error on CreateDelivery failure, got nil")
	}

	if senderCalled {
		t.Errorf("sender should not be called when CreateDelivery fails")
	}
}

func TestProcessor_SendDelivery_SendSuccess_MarkSentError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	dbErr := errors.New("db connection lost")

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Popeye", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliverySentFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{}, dbErr
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		return "resend_msg_abc", nil
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, err := processor.SendDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error when MarkEmailDeliverySent fails, got nil")
	}

	if !errors.Is(err, dbErr) {
		t.Errorf("expected wrapped dbErr, got %v", err)
	}
}

func TestProcessor_SendDelivery_SendFailure_MarkFailedError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	sendErr := errors.New("provider timeout")
	dbErr := errors.New("db write failure")

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Popeye", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliveryFailedFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{}, dbErr
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		return "", sendErr
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	_, err := processor.SendDelivery(ctx, notification)
	if err == nil {
		t.Fatal("expected error when both send and MarkEmailDeliveryFailed fail, got nil")
	}

	if !errors.Is(err, sendErr) {
		t.Errorf("expected wrapped sendErr, got %v", err)
	}
}

func TestProcessor_ProcessNotification_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	var markProcessingCalled bool
	var markCompletedCalled bool
	var markFailedCalled bool
	var senderCalled bool

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "captain@example.com",
				Forename: sql.NullString{String: "Ahab", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliverySentFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:                arg.ID,
				NotificationID:    notificationID,
				Recipient:         "captain@example.com",
				Provider:          ProviderResend,
				ProviderMessageID: arg.ProviderMessageID,
				Status:            "sent",
				Attempt:           1,
			}, nil
		},
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			if id != notificationID {
				t.Fatalf("expected notificationID %s, got %s", notificationID, id)
			}
			markProcessingCalled = true
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "processing",
			}, nil
		},
		MarkNotificationCompletedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			if id != notificationID {
				t.Fatalf("expected notificationID %s, got %s", notificationID, id)
			}
			markCompletedCalled = true
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "completed",
			}, nil
		},
		MarkNotificationFailedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markFailedCalled = true
			return sqlc.Notification{}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		senderCalled = true
		return "msg_resend_999", nil
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	err := processor.ProcessNotification(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	if !markProcessingCalled {
		t.Errorf("expected MarkNotificationProcessing to be called")
	}

	if !senderCalled {
		t.Errorf("expected sender to be called")
	}

	if !markCompletedCalled {
		t.Errorf("expected MarkNotificationCompleted to be called")
	}

	if markFailedCalled {
		t.Errorf("MarkNotificationFailed should not be called on success")
	}
}

func TestProcessor_ProcessNotification_OrderOfTransitions(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	var eventOrder []string

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "captain@example.com",
				Forename: sql.NullString{String: "Ahab", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliverySentFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:     arg.ID,
				Status: "sent",
			}, nil
		},
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			eventOrder = append(eventOrder, "mark_processing")
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "processing",
			}, nil
		},
		MarkNotificationCompletedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			eventOrder = append(eventOrder, "mark_completed")
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "completed",
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		eventOrder = append(eventOrder, "send_email")
		return "msg_id_123", nil
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	err := processor.ProcessNotification(ctx, notification)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	expectedOrder := []string{"mark_processing", "send_email", "mark_completed"}
	if len(eventOrder) != len(expectedOrder) {
		t.Fatalf("expected event order %v, got %v", expectedOrder, eventOrder)
	}
	for i, ev := range expectedOrder {
		if eventOrder[i] != ev {
			t.Errorf("at step %d: expected %s, got %s", i, ev, eventOrder[i])
		}
	}
}

func TestProcessor_ProcessNotification_SendFailure_MarksFailed(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	var markProcessingCalled bool
	var markFailedCalled bool
	var markCompletedCalled bool
	var markFailedID uuid.UUID

	resendErr := errors.New("resend down")

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "captain@example.com",
				Forename: sql.NullString{String: "Ahab", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliveryFailedFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:           arg.ID,
				Status:       "failed",
				ErrorMessage: arg.ErrorMessage,
			}, nil
		},
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markProcessingCalled = true
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "processing",
			}, nil
		},
		MarkNotificationCompletedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markCompletedCalled = true
			return sqlc.Notification{}, nil
		},
		MarkNotificationFailedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markFailedCalled = true
			markFailedID = id
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "failed",
			}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		return "", resendErr
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	err := processor.ProcessNotification(ctx, notification)
	if err == nil {
		t.Fatal("expected error from ProcessNotification on send failure, got nil")
	}

	if !errors.Is(err, resendErr) {
		t.Errorf("expected wrapped resendErr, got %v", err)
	}

	if !markProcessingCalled {
		t.Errorf("expected MarkNotificationProcessing to be called")
	}

	if !markFailedCalled {
		t.Errorf("expected MarkNotificationFailed to be called")
	}

	if markFailedID != notificationID {
		t.Errorf("expected MarkNotificationFailed with ID %s, got %s", notificationID, markFailedID)
	}

	if markCompletedCalled {
		t.Errorf("MarkNotificationCompleted should not be called when sending fails")
	}
}

func TestProcessor_ProcessNotification_UnsupportedType(t *testing.T) {
	ctx := context.Background()

	var markProcessingCalled bool
	var markCompletedCalled bool
	var markFailedCalled bool
	var senderCalled bool

	mockRepo := &mockRepository{
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markProcessingCalled = true
			return sqlc.Notification{}, nil
		},
		MarkNotificationCompletedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markCompletedCalled = true
			return sqlc.Notification{}, nil
		},
		MarkNotificationFailedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markFailedCalled = true
			return sqlc.Notification{}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		senderCalled = true
		return "msg", nil
	}

	testCases := []string{
		string(TypeNoCertificates1Month),
		string(TypeCertificateExpirySummary),
		"unsupported_type",
		"",
	}

	for _, tc := range testCases {
		t.Run(tc, func(t *testing.T) {
			markProcessingCalled = false
			markCompletedCalled = false
			markFailedCalled = false
			senderCalled = false

			notification := sqlc.Notification{
				ID:               uuid.New(),
				UserID:           uuid.NullUUID{UUID: uuid.New(), Valid: true},
				NotificationType: tc,
				Status:           "pending",
			}

			err := processor.ProcessNotification(ctx, notification)
			if err == nil {
				t.Fatalf("expected error for unsupported notification type %q, got nil", tc)
			}

			if markProcessingCalled {
				t.Errorf("MarkNotificationProcessing should not be called for unsupported type %q", tc)
			}
			if senderCalled {
				t.Errorf("sender should not be called for unsupported type %q", tc)
			}
			if markCompletedCalled {
				t.Errorf("MarkNotificationCompleted should not be called for unsupported type %q", tc)
			}
			if markFailedCalled {
				t.Errorf("MarkNotificationFailed should not be called for unsupported type %q", tc)
			}
		})
	}
}

func TestProcessor_ProcessNotification_MarkProcessingError(t *testing.T) {
	ctx := context.Background()
	notificationID := uuid.New()
	userID := uuid.New()

	dbErr := errors.New("db error updating to processing")
	var senderCalled bool
	var markCompletedCalled bool
	var markFailedCalled bool

	mockRepo := &mockRepository{
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			return sqlc.Notification{}, dbErr
		},
		MarkNotificationCompletedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markCompletedCalled = true
			return sqlc.Notification{}, nil
		},
		MarkNotificationFailedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markFailedCalled = true
			return sqlc.Notification{}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		senderCalled = true
		return "msg", nil
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	err := processor.ProcessNotification(ctx, notification)
	if err == nil {
		t.Fatal("expected error when MarkNotificationProcessing fails, got nil")
	}

	if !errors.Is(err, dbErr) {
		t.Errorf("expected wrapped dbErr, got %v", err)
	}

	if senderCalled {
		t.Errorf("sender should not be called when MarkNotificationProcessing fails")
	}

	if markCompletedCalled {
		t.Errorf("MarkNotificationCompleted should not be called")
	}

	if markFailedCalled {
		t.Errorf("MarkNotificationFailed should not be called")
	}
}

func TestProcessor_ProcessNotification_SendSuccess_MarkCompletedError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	dbErr := errors.New("db failed on mark completed")
	senderCalls := 0
	var markFailedCalled bool

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Popeye", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliverySentFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliverySentParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:     arg.ID,
				Status: "sent",
			}, nil
		},
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "processing",
			}, nil
		},
		MarkNotificationCompletedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			return sqlc.Notification{}, dbErr
		},
		MarkNotificationFailedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			markFailedCalled = true
			return sqlc.Notification{}, nil
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		senderCalls++
		return "msg_sent_123", nil
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	err := processor.ProcessNotification(ctx, notification)
	if err == nil {
		t.Fatal("expected error when MarkNotificationCompleted fails, got nil")
	}

	if !errors.Is(err, dbErr) {
		t.Errorf("expected wrapped dbErr, got %v", err)
	}

	if !strings.Contains(err.Error(), "email sent successfully") {
		t.Errorf("expected error to clearly state email was sent successfully, got %v", err)
	}

	if senderCalls != 1 {
		t.Errorf("expected exactly 1 sender call (do not resend), got %d", senderCalls)
	}

	if markFailedCalled {
		t.Errorf("MarkNotificationFailed should not be called when send was successful")
	}
}

func TestProcessor_ProcessNotification_SendFailure_MarkFailedError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()
	notificationID := uuid.New()
	deliveryID := uuid.New()

	sendErr := errors.New("resend connection timeout")
	dbErr := errors.New("db disk full")

	mockRepo := &mockRepository{
		GetUserByIDFunc: func(ctx context.Context, id uuid.UUID) (sqlc.User, error) {
			return sqlc.User{
				ID:       userID,
				Email:    "sailor@example.com",
				Forename: sql.NullString{String: "Popeye", Valid: true},
			}, nil
		},
		GetEmailDeliveriesForNotificationFunc: func(ctx context.Context, notifID uuid.UUID) ([]sqlc.EmailDelivery, error) {
			return []sqlc.EmailDelivery{}, nil
		},
		CreateEmailDeliveryFunc: func(ctx context.Context, arg sqlc.CreateEmailDeliveryParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:             deliveryID,
				NotificationID: arg.NotificationID,
				Recipient:      arg.Recipient,
				Provider:       arg.Provider,
				Status:         "pending",
				Attempt:        arg.Attempt,
			}, nil
		},
		MarkEmailDeliveryFailedFunc: func(ctx context.Context, arg sqlc.MarkEmailDeliveryFailedParams) (sqlc.EmailDelivery, error) {
			return sqlc.EmailDelivery{
				ID:           arg.ID,
				Status:       "failed",
				ErrorMessage: arg.ErrorMessage,
			}, nil
		},
		MarkNotificationProcessingFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			return sqlc.Notification{
				ID:               notificationID,
				UserID:           uuid.NullUUID{UUID: userID, Valid: true},
				NotificationType: string(TypeNoCertificates7Day),
				Status:           "processing",
			}, nil
		},
		MarkNotificationFailedFunc: func(ctx context.Context, id uuid.UUID) (sqlc.Notification, error) {
			return sqlc.Notification{}, dbErr
		},
	}

	processor := NewProcessor(mockRepo)
	processor.sender = func(email email_templates.Email) (string, error) {
		return "", sendErr
	}

	notification := sqlc.Notification{
		ID:               notificationID,
		UserID:           uuid.NullUUID{UUID: userID, Valid: true},
		NotificationType: string(TypeNoCertificates7Day),
		Status:           "pending",
	}

	err := processor.ProcessNotification(ctx, notification)
	if err == nil {
		t.Fatal("expected error when both sending and MarkNotificationFailed fail, got nil")
	}

	if !errors.Is(err, sendErr) {
		t.Errorf("expected wrapped sendErr, got %v", err)
	}

	if !strings.Contains(err.Error(), dbErr.Error()) {
		t.Errorf("expected error to preserve MarkNotificationFailed context %q, got %v", dbErr.Error(), err)
	}
}
