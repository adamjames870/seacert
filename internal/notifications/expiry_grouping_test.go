package notifications

import (
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/google/uuid"
)

func makeCert(idStr string, certNumber string, expiryDate time.Time) sqlc.GetCertificatesForExpiryNotificationRow {
	var id uuid.UUID
	if idStr != "" {
		id = uuid.MustParse(idStr)
	} else {
		id = uuid.New()
	}
	return sqlc.GetCertificatesForExpiryNotificationRow{
		ID:                id,
		CertNumber:        certNumber,
		IssuedDate:        expiryDate.AddDate(-2, 0, 0),
		CertTypeID:        uuid.MustParse("10000000-0000-0000-0000-000000000001"),
		CertTypeName:      "Standard Qualification",
		CertTypeShortName: "SQ",
		ExpiryDate:        expiryDate,
	}
}

func TestGroupCertificatesByExpiry_Boundaries(t *testing.T) {
	// Stable reference date: 2026-08-28 14:30:00 UTC
	refTime := time.Date(2026, 8, 28, 14, 30, 0, 0, time.UTC)
	refDate := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)

	certExpiredYearsAgo := makeCert("10000000-0000-0000-0000-000000000001", "EXP-YEARS-AGO", refDate.AddDate(-3, 0, 0))
	certExpiredYesterday := makeCert("10000000-0000-0000-0000-000000000002", "EXP-YESTERDAY", refDate.AddDate(0, 0, -1))
	certExpiresToday := makeCert("10000000-0000-0000-0000-000000000003", "EXP-TODAY", refDate)
	certExpiresTomorrow := makeCert("10000000-0000-0000-0000-000000000004", "EXP-TOMORROW", refDate.AddDate(0, 0, 1))
	certExpiresExact1Month := makeCert("10000000-0000-0000-0000-000000000005", "EXP-EXACT-1M", refDate.AddDate(0, 1, 0))
	certExpires1MonthPlus1Day := makeCert("10000000-0000-0000-0000-000000000006", "EXP-1M-PLUS-1D", refDate.AddDate(0, 1, 1))
	certExpiresExact6Months := makeCert("10000000-0000-0000-0000-000000000007", "EXP-EXACT-6M", refDate.AddDate(0, 6, 0))
	certExpires6MonthsPlus1Day := makeCert("10000000-0000-0000-0000-000000000008", "EXP-6M-PLUS-1D", refDate.AddDate(0, 6, 1))
	certExpiresExact1Year := makeCert("10000000-0000-0000-0000-000000000009", "EXP-EXACT-1Y", refDate.AddDate(1, 0, 0))
	certExpires1YearPlus1Day := makeCert("10000000-0000-0000-0000-000000000010", "EXP-1Y-PLUS-1D", refDate.AddDate(1, 0, 1))
	certNoExpiry := makeCert("10000000-0000-0000-0000-000000000011", "EXP-ZERO-DATE", time.Time{})

	// Shuffle input order to test sorting
	input := []sqlc.GetCertificatesForExpiryNotificationRow{
		certExpires1YearPlus1Day,
		certExpires1MonthPlus1Day,
		certExpiredYesterday,
		certExpiresTomorrow,
		certNoExpiry,
		certExpiresExact1Year,
		certExpiredYearsAgo,
		certExpiresToday,
		certExpiresExact6Months,
		certExpiresExact1Month,
		certExpires6MonthsPlus1Day,
	}

	groups := GroupCertificatesByExpiry(input, refTime)

	// 1. Expired bucket
	if len(groups.Expired) != 2 {
		t.Fatalf("expected 2 in Expired, got %d", len(groups.Expired))
	}
	if groups.Expired[0].CertNumber != "EXP-YEARS-AGO" {
		t.Errorf("expected groups.Expired[0] to be EXP-YEARS-AGO, got %s", groups.Expired[0].CertNumber)
	}
	if groups.Expired[1].CertNumber != "EXP-YESTERDAY" {
		t.Errorf("expected groups.Expired[1] to be EXP-YESTERDAY, got %s", groups.Expired[1].CertNumber)
	}

	// 2. ExpiringOneMonth bucket (Today, Tomorrow, Exact 1 Month)
	if len(groups.ExpiringOneMonth) != 3 {
		t.Fatalf("expected 3 in ExpiringOneMonth, got %d", len(groups.ExpiringOneMonth))
	}
	if groups.ExpiringOneMonth[0].CertNumber != "EXP-TODAY" {
		t.Errorf("expected groups.ExpiringOneMonth[0] to be EXP-TODAY, got %s", groups.ExpiringOneMonth[0].CertNumber)
	}
	if groups.ExpiringOneMonth[1].CertNumber != "EXP-TOMORROW" {
		t.Errorf("expected groups.ExpiringOneMonth[1] to be EXP-TOMORROW, got %s", groups.ExpiringOneMonth[1].CertNumber)
	}
	if groups.ExpiringOneMonth[2].CertNumber != "EXP-EXACT-1M" {
		t.Errorf("expected groups.ExpiringOneMonth[2] to be EXP-EXACT-1M, got %s", groups.ExpiringOneMonth[2].CertNumber)
	}

	// 3. ExpiringSixMonths bucket (1 Month + 1 Day, Exact 6 Months)
	if len(groups.ExpiringSixMonths) != 2 {
		t.Fatalf("expected 2 in ExpiringSixMonths, got %d", len(groups.ExpiringSixMonths))
	}
	if groups.ExpiringSixMonths[0].CertNumber != "EXP-1M-PLUS-1D" {
		t.Errorf("expected groups.ExpiringSixMonths[0] to be EXP-1M-PLUS-1D, got %s", groups.ExpiringSixMonths[0].CertNumber)
	}
	if groups.ExpiringSixMonths[1].CertNumber != "EXP-EXACT-6M" {
		t.Errorf("expected groups.ExpiringSixMonths[1] to be EXP-EXACT-6M, got %s", groups.ExpiringSixMonths[1].CertNumber)
	}

	// 4. ExpiringOneYear bucket (6 Months + 1 Day, Exact 1 Year)
	if len(groups.ExpiringOneYear) != 2 {
		t.Fatalf("expected 2 in ExpiringOneYear, got %d", len(groups.ExpiringOneYear))
	}
	if groups.ExpiringOneYear[0].CertNumber != "EXP-6M-PLUS-1D" {
		t.Errorf("expected groups.ExpiringOneYear[0] to be EXP-6M-PLUS-1D, got %s", groups.ExpiringOneYear[0].CertNumber)
	}
	if groups.ExpiringOneYear[1].CertNumber != "EXP-EXACT-1Y" {
		t.Errorf("expected groups.ExpiringOneYear[1] to be EXP-EXACT-1Y, got %s", groups.ExpiringOneYear[1].CertNumber)
	}

	// 5. Excluded items verification (no certificate should appear in multiple buckets, total count = 9)
	seen := make(map[uuid.UUID]int)
	allGroups := [][]sqlc.GetCertificatesForExpiryNotificationRow{
		groups.Expired,
		groups.ExpiringOneMonth,
		groups.ExpiringSixMonths,
		groups.ExpiringOneYear,
	}
	for _, grp := range allGroups {
		for _, cert := range grp {
			seen[cert.ID]++
		}
	}

	for id, count := range seen {
		if count > 1 {
			t.Errorf("certificate %s appeared in %d groups", id, count)
		}
	}

	if _, exists := seen[certExpires1YearPlus1Day.ID]; exists {
		t.Errorf("certExpires1YearPlus1Day (>1 year) should have been excluded")
	}
	if _, exists := seen[certNoExpiry.ID]; exists {
		t.Errorf("certNoExpiry should have been excluded")
	}
}

func TestGroupCertificatesByExpiry_TimeOfDayIndependence(t *testing.T) {
	// Reference time late at night: 23:59:59
	lateRefTime := time.Date(2026, 8, 28, 23, 59, 59, 0, time.UTC)
	// Expiry date set at 00:00:00 on the same day
	certExpiresSameDayEarly := makeCert("10000000-0000-0000-0000-000000000001", "EXP-SAME-DAY-EARLY", time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC))

	groups := GroupCertificatesByExpiry([]sqlc.GetCertificatesForExpiryNotificationRow{certExpiresSameDayEarly}, lateRefTime)

	if len(groups.Expired) != 0 {
		t.Errorf("same calendar day should not be marked Expired even if time-of-day of reference is later")
	}
	if len(groups.ExpiringOneMonth) != 1 {
		t.Errorf("same calendar day must be in ExpiringOneMonth, got count %d", len(groups.ExpiringOneMonth))
	}
}

func TestGroupCertificatesByExpiry_CalendarArithmeticMonthEnd(t *testing.T) {
	// Test month-end handling: January 31 in leap year 2024
	refTime := time.Date(2024, 1, 31, 10, 0, 0, 0, time.UTC)
	refDate := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	// In Go: time.Date(2024, 1, 31, ...).AddDate(0, 1, 0) normalizes to 2024-03-02 (Feb 2024 has 29 days)
	boundary1Month := refDate.AddDate(0, 1, 0)
	boundary6Months := refDate.AddDate(0, 6, 0)
	boundary1Year := refDate.AddDate(1, 0, 0)

	certExact1M := makeCert("10000000-0000-0000-0000-000000000001", "EXP-1M", boundary1Month)
	cert1MPlus1D := makeCert("10000000-0000-0000-0000-000000000002", "EXP-1M-1D", boundary1Month.AddDate(0, 0, 1))
	certExact6M := makeCert("10000000-0000-0000-0000-000000000003", "EXP-6M", boundary6Months)
	cert6MPlus1D := makeCert("10000000-0000-0000-0000-000000000004", "EXP-6M-1D", boundary6Months.AddDate(0, 0, 1))
	certExact1Y := makeCert("10000000-0000-0000-0000-000000000005", "EXP-1Y", boundary1Year)
	cert1YPlus1D := makeCert("10000000-0000-0000-0000-000000000006", "EXP-1Y-1D", boundary1Year.AddDate(0, 0, 1))

	groups := GroupCertificatesByExpiry([]sqlc.GetCertificatesForExpiryNotificationRow{
		certExact1M,
		cert1MPlus1D,
		certExact6M,
		cert6MPlus1D,
		certExact1Y,
		cert1YPlus1D,
	}, refTime)

	if len(groups.ExpiringOneMonth) != 1 || groups.ExpiringOneMonth[0].CertNumber != "EXP-1M" {
		t.Errorf("expected EXP-1M in ExpiringOneMonth, got %v", groups.ExpiringOneMonth)
	}
	if len(groups.ExpiringSixMonths) != 2 || groups.ExpiringSixMonths[0].CertNumber != "EXP-1M-1D" || groups.ExpiringSixMonths[1].CertNumber != "EXP-6M" {
		t.Errorf("expected EXP-1M-1D and EXP-6M in ExpiringSixMonths, got %v", groups.ExpiringSixMonths)
	}
	if len(groups.ExpiringOneYear) != 2 || groups.ExpiringOneYear[0].CertNumber != "EXP-6M-1D" || groups.ExpiringOneYear[1].CertNumber != "EXP-1Y" {
		t.Errorf("expected EXP-6M-1D and EXP-1Y in ExpiringOneYear, got %v", groups.ExpiringOneYear)
	}
}

func TestGroupCertificatesByExpiry_StableSorting(t *testing.T) {
	refTime := time.Date(2026, 5, 1, 12, 0, 0, 0, time.UTC)
	sameExpiry := time.Date(2026, 6, 1, 0, 0, 0, 0, time.UTC)

	certB := makeCert("10000000-0000-0000-0000-000000000002", "CERT-B", sameExpiry)
	certA := makeCert("10000000-0000-0000-0000-000000000001", "CERT-A", sameExpiry)

	groups := GroupCertificatesByExpiry([]sqlc.GetCertificatesForExpiryNotificationRow{certB, certA}, refTime)

	if len(groups.ExpiringOneMonth) != 2 {
		t.Fatalf("expected 2 certs, got %d", len(groups.ExpiringOneMonth))
	}
	if groups.ExpiringOneMonth[0].CertNumber != "CERT-A" || groups.ExpiringOneMonth[1].CertNumber != "CERT-B" {
		t.Errorf("expected CERT-A before CERT-B due to ID secondary sort, got [%s, %s]",
			groups.ExpiringOneMonth[0].CertNumber, groups.ExpiringOneMonth[1].CertNumber)
	}
}

func TestGroupCertificatesByExpiry_EmptyInput(t *testing.T) {
	refTime := time.Date(2026, 8, 28, 0, 0, 0, 0, time.UTC)
	groups := GroupCertificatesByExpiry(nil, refTime)

	if groups.Expired == nil || groups.ExpiringOneMonth == nil || groups.ExpiringSixMonths == nil || groups.ExpiringOneYear == nil {
		t.Errorf("expected non-nil initialized slices in groups")
	}
	if len(groups.Expired) != 0 || len(groups.ExpiringOneMonth) != 0 || len(groups.ExpiringSixMonths) != 0 || len(groups.ExpiringOneYear) != 0 {
		t.Errorf("expected empty slices for all groups")
	}
}
