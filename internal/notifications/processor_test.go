package notifications

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"testing"

	"github.com/adamjames870/seacert/internal/database/sqlc"
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
