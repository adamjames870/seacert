package certificates

import (
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/domain/cert_types"
)

func TestCalculateExpiryDate(t *testing.T) {
	issueDate := time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC)

	t.Run("manual expiry set", func(t *testing.T) {
		manualExpiry := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
		cert := &Certificate{
			IssuedDate:   issueDate,
			ManualExpiry: manualExpiry,
		}
		cert.calculateExpiryDate()
		if !cert.ExpiryDate.Equal(manualExpiry) {
			t.Errorf("expected expiry %v, got %v", manualExpiry, cert.ExpiryDate)
		}
	})

	t.Run("validity months set", func(t *testing.T) {
		cert := &Certificate{
			IssuedDate: issueDate,
			CertType: cert_types.CertificateType{
				NormalValidityMonths: 12,
			},
		}
		expectedExpiry := getExpiryAfterValidity(issueDate, 12)
		cert.calculateExpiryDate()
		if !cert.ExpiryDate.Equal(expectedExpiry) {
			t.Errorf("expected expiry %v, got %v", expectedExpiry, cert.ExpiryDate)
		}
	})

	t.Run("nothing set", func(t *testing.T) {
		cert := &Certificate{
			IssuedDate: issueDate,
		}
		cert.calculateExpiryDate()
		if !cert.ExpiryDate.IsZero() {
			t.Errorf("expected zero expiry, got %v", cert.ExpiryDate)
		}
	})
}

func TestGetExpiryAfterValidity(t *testing.T) {
	tests := []struct {
		name           string
		issueDate      time.Time
		validityMonths int
		want           time.Time
	}{
		{
			name:           "normal case",
			issueDate:      time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
			validityMonths: 12,
			want:           time.Date(2024, 1, 14, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "leap year",
			issueDate:      time.Date(2024, 2, 29, 0, 0, 0, 0, time.UTC),
			validityMonths: 12,
			want:           time.Date(2025, 3, 28, 0, 0, 0, 0, time.UTC),
		},
		{
			name:           "end of month adjustment (31 Jan -> 30 Apr)",
			issueDate:      time.Date(2023, 1, 31, 0, 0, 0, 0, time.UTC),
			validityMonths: 3,
			want:           time.Date(2023, 5, 30, 0, 0, 0, 0, time.UTC),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := getExpiryAfterValidity(tt.issueDate, tt.validityMonths)
			if !got.Equal(tt.want) {
				t.Errorf("getExpiryAfterValidity() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestDaysInMonth(t *testing.T) {
	tests := []struct {
		year  int
		month time.Month
		want  int
	}{
		{2023, time.January, 31},
		{2023, time.February, 28},
		{2024, time.February, 29}, // Leap year
		{2023, time.April, 30},
	}

	for _, tt := range tests {
		got := daysInMonth(tt.year, tt.month)
		if got != tt.want {
			t.Errorf("daysInMonth(%d, %v) = %d, want %d", tt.year, tt.month, got, tt.want)
		}
	}
}
