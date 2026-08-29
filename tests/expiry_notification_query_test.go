package tests

import (
	"context"
	"database/sql"
	"os"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/notifications"
	"github.com/adamjames870/seacert/internal/repository/postgres"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

func getTestDB(t *testing.T) *sql.DB {
	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		dbURL = "user=adam dbname=seacert sslmode=disable host=/var/run/postgresql"
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		t.Skipf("skipping test; could not open db: %v", err)
	}

	if err := db.Ping(); err != nil {
		t.Skipf("skipping test; could not ping db: %v", err)
	}

	return db
}

func TestGetCertificatesForExpiryNotification_User009DummyData(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()
	repo := postgres.NewRepository(db)

	// User 009 from dummy records
	user009ID := uuid.MustParse("10000000-0000-0000-0000-000000000009")

	certs, err := repo.GetCertificatesForExpiryNotification(ctx, user009ID)
	if err != nil {
		t.Fatalf("failed to get certificates for expiry notification: %v", err)
	}

	// In dummy records, User 009 has 5 certificates:
	// 1. EXP-EXPIRED-001: issued 25 months ago (expired ~1 month ago) -> MUST be returned
	// 2. EXP-1MONTH-001: issued 23.5 months ago (expires in ~15 days) -> MUST be returned
	// 3. EXP-6MONTH-001: issued 20 months ago (expires in ~4 months) -> MUST be returned
	// 4. EXP-1YEAR-001: issued 16 months ago (expires in ~8 months) -> MUST be returned
	// 5. EXP-NOT-DUE-001: issued 2 months ago (expires in ~22 months) -> MUST be EXCLUDED (> 1 year)

	if len(certs) != 4 {
		t.Fatalf("expected 4 certificates for user 009, got %d", len(certs))
	}

	expectedCertNumbers := []string{
		"EXP-EXPIRED-001",
		"EXP-1MONTH-001",
		"EXP-6MONTH-001",
		"EXP-1YEAR-001",
	}

	for i, expectedNumber := range expectedCertNumbers {
		if certs[i].CertNumber != expectedNumber {
			t.Errorf("at index %d: expected cert number %s, got %s", i, expectedNumber, certs[i].CertNumber)
		}
		if certs[i].CertTypeName != "Seafarer Medical - UK - ENG1" {
			t.Errorf("at index %d: expected cert type name 'Seafarer Medical - UK - ENG1', got %s", i, certs[i].CertTypeName)
		}
		if certs[i].CertTypeShortName != "ENG1" {
			t.Errorf("at index %d: expected short name 'ENG1', got %s", i, certs[i].CertTypeShortName)
		}
		if certs[i].ExpiryDate.IsZero() {
			t.Errorf("at index %d: expected non-zero expiry date", i)
		}
	}

	// Verify ordering is strictly ascending by expiry date
	for i := 1; i < len(certs); i++ {
		if certs[i].ExpiryDate.Before(certs[i-1].ExpiryDate) {
			t.Errorf("certs not ordered by expiry_date asc: index %d (%v) is before index %d (%v)",
				i, certs[i].ExpiryDate, i-1, certs[i-1].ExpiryDate)
		}
	}

	// Verify certificates belong ONLY to user 009 (user 003 has other certs)
	user003ID := uuid.MustParse("10000000-0000-0000-0000-000000000003")
	user003Certs, err := repo.GetCertificatesForExpiryNotification(ctx, user003ID)
	if err != nil {
		t.Fatalf("failed to query for user 003: %v", err)
	}

	// User 003 has a cert issued NOW with 60 months validity, so expiry is 5 years away (> 1 year) -> should return 0
	for _, c := range user003Certs {
		if c.CertNumber == "EXP-EXPIRED-001" || c.CertNumber == "EXP-1MONTH-001" {
			t.Errorf("user 003 unexpectedly received user 009's certificate: %s", c.CertNumber)
		}
	}
}

func TestGetCertificatesForExpiryNotification_EdgeCasesInTx(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Use a transaction so we can rollback test fixtures
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}
	defer tx.Rollback()

	qtx := sqlc.New(tx)

	testUserID := uuid.New()
	issuerID := uuid.MustParse("10000000-0000-0000-0000-000000000001")
	expiringTypeID := uuid.MustParse("10000000-0000-0000-0000-000000000004") // ENG1: 24 months

	// Create a non-expiring cert type (normal_validity_months = NULL)
	nonExpiringTypeID := uuid.New()
	_, err = qtx.CreateCertType(ctx, sqlc.CreateCertTypeParams{
		ID:                   nonExpiringTypeID,
		Name:                 "Non-Expiring Test Qualification",
		ShortName:            "NETQ",
		StcwReference:        sql.NullString{String: "VI/1", Valid: true},
		NormalValidityMonths: sql.NullInt32{Valid: false},
		Status:               "approved",
		CreatedBy:            uuid.NullUUID{Valid: false},
	})
	if err != nil {
		t.Fatalf("failed to create non-expiring cert type: %v", err)
	}

	// Create test user
	_, err = qtx.CreateUser(ctx, sqlc.CreateUserParams{
		ID:          testUserID,
		Forename:    sql.NullString{String: "Expiry", Valid: true},
		Surname:     sql.NullString{String: "Tester", Valid: true},
		Email:       "expiry-tester-" + testUserID.String() + "@example.com",
		Nationality: sql.NullString{String: "GB", Valid: true},
	})
	if err != nil {
		t.Fatalf("failed to create test user: %v", err)
	}

	now := time.Now().UTC().Truncate(time.Second)

	// 1. Non-expiring certificate (should be EXCLUDED)
	_, err = qtx.CreateCert(ctx, sqlc.CreateCertParams{
		ID:         uuid.New(),
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     testUserID,
		CertTypeID: nonExpiringTypeID,
		CertNumber: "NON-EXPIRING-001",
		IssuerID:   issuerID,
		IssuedDate: now.AddDate(-2, 0, 0),
	})
	if err != nil {
		t.Fatalf("failed to create cert 1: %v", err)
	}

	// 2. Non-expiring certificate WITH manual_expiry within 6 months (should be INCLUDED)
	manualExpiryDate := now.AddDate(0, 3, 0)
	_, err = qtx.CreateCert(ctx, sqlc.CreateCertParams{
		ID:           uuid.New(),
		CreatedAt:    now,
		UpdatedAt:    now,
		UserID:       testUserID,
		CertTypeID:   nonExpiringTypeID,
		CertNumber:   "MANUAL-EXP-001",
		IssuerID:     issuerID,
		IssuedDate:   now.AddDate(-1, 0, 0),
		ManualExpiry: sql.NullTime{Time: manualExpiryDate, Valid: true},
	})
	if err != nil {
		t.Fatalf("failed to create cert 2: %v", err)
	}

	// 3. Expiring cert type with manual_expiry override (manual expiry > 1 year -> should be EXCLUDED)
	_, err = qtx.CreateCert(ctx, sqlc.CreateCertParams{
		ID:           uuid.New(),
		CreatedAt:    now,
		UpdatedAt:    now,
		UserID:       testUserID,
		CertTypeID:   expiringTypeID, // ENG1 24 months
		CertNumber:   "MANUAL-OVERRIDE-FAR",
		IssuerID:     issuerID,
		IssuedDate:   now.AddDate(-2, 0, 0), // would have expired now by normal validity, but manual is 2 years away
		ManualExpiry: sql.NullTime{Time: now.AddDate(2, 0, 0), Valid: true},
	})
	if err != nil {
		t.Fatalf("failed to create cert 3: %v", err)
	}

	// 4. Expiring cert type deleted/retired (should be EXCLUDED)
	deletedCertID := uuid.New()
	_, err = qtx.CreateCert(ctx, sqlc.CreateCertParams{
		ID:         deletedCertID,
		CreatedAt:  now,
		UpdatedAt:  now,
		UserID:     testUserID,
		CertTypeID: expiringTypeID,
		CertNumber: "DELETED-CERT-001",
		IssuerID:   issuerID,
		IssuedDate: now.AddDate(-1, 0, 0),
	})
	if err != nil {
		t.Fatalf("failed to create cert 4: %v", err)
	}
	_, err = qtx.UpdateCertificate(ctx, sqlc.UpdateCertificateParams{
		ID:      deletedCertID,
		Deleted: sql.NullBool{Bool: true, Valid: true},
	})
	if err != nil {
		t.Fatalf("failed to mark cert 4 as deleted: %v", err)
	}

	// Query for test user
	results, err := qtx.GetCertificatesForExpiryNotification(ctx, testUserID)
	if err != nil {
		t.Fatalf("GetCertificatesForExpiryNotification failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected exactly 1 result (the manual-expiry one), got %d: %+v", len(results), results)
	}

	if results[0].CertNumber != "MANUAL-EXP-001" {
		t.Errorf("expected MANUAL-EXP-001, got %s", results[0].CertNumber)
	}

	if !results[0].ExpiryDate.Equal(manualExpiryDate) {
		t.Errorf("expected expiry date %v, got %v", manualExpiryDate, results[0].ExpiryDate)
	}
}

func TestCandidateUsers_And_GenerateExpirySummary_Integration(t *testing.T) {
	db := getTestDB(t)
	defer db.Close()

	ctx := context.Background()

	// Run within transaction and roll back at end so dummy data remains clean
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to start tx: %v", err)
	}
	defer tx.Rollback()

	queries := sqlc.New(tx)

	candidateUserIDs, err := queries.GetCandidateUsersForExpiryNotification(ctx)
	if err != nil {
		t.Fatalf("failed to get candidate users: %v", err)
	}

	user009ID := uuid.MustParse("10000000-0000-0000-0000-000000000009")
	found009 := false
	for _, id := range candidateUserIDs {
		if id == user009ID {
			found009 = true
			break
		}
	}
	if !found009 {
		t.Errorf("expected user 009 to be in candidate users list, got %v", candidateUserIDs)
	}

	// Test notification key query
	testKey := "test-notification-key-" + uuid.New().String()
	notifID := uuid.New()
	_, err = queries.CreateNotification(ctx, sqlc.CreateNotificationParams{
		ID:               notifID,
		UserID:           uuid.NullUUID{UUID: user009ID, Valid: true},
		NotificationType: string(notifications.TypeCertificateExpirySummary),
		NotificationKey:  testKey,
		Payload:          []byte("{}"),
		ScheduledAt:      time.Now(),
	})
	if err != nil {
		t.Fatalf("failed to create notification in tx: %v", err)
	}

	fetched, err := queries.GetNotificationByKey(ctx, testKey)
	if err != nil {
		t.Fatalf("failed to fetch notification by key: %v", err)
	}
	if fetched.ID != notifID {
		t.Errorf("expected notification ID %s, got %s", notifID, fetched.ID)
	}
	if fetched.NotificationType != string(notifications.TypeCertificateExpirySummary) {
		t.Errorf("expected notification type %s, got %s", notifications.TypeCertificateExpirySummary, fetched.NotificationType)
	}
}
