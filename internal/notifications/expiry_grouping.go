package notifications

import (
	"sort"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
)

// CertificateExpiryGroups holds certificates divided into the four expiry buckets
// required by the monthly certificate-expiry summary email.
type CertificateExpiryGroups struct {
	Expired           []sqlc.GetCertificatesForExpiryNotificationRow
	ExpiringOneMonth  []sqlc.GetCertificatesForExpiryNotificationRow
	ExpiringSixMonths []sqlc.GetCertificatesForExpiryNotificationRow
	ExpiringOneYear   []sqlc.GetCertificatesForExpiryNotificationRow
}

// GroupCertificatesByExpiry is a pure function that divides certificates into four
// mutually exclusive expiry categories based on a reference time `now`:
//   - Expired: expiry < today
//   - Expiring within next month: expiry >= today AND expiry <= today + 1 month
//   - Expiring within 6 months: expiry > today + 1 month AND expiry <= today + 6 months
//   - Expiring within 1 year: expiry > today + 6 months AND expiry <= today + 1 year
//
// Certificates expiring more than 1 year away, or those without an expiry date, are excluded.
// Within each group, certificates are sorted deterministically by expiry date ascending,
// with certificate ID as a secondary tie-breaker.
func GroupCertificatesByExpiry(
	certificates []sqlc.GetCertificatesForExpiryNotificationRow,
	now time.Time,
) CertificateExpiryGroups {
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)
	oneMonth := today.AddDate(0, 1, 0)
	sixMonths := today.AddDate(0, 6, 0)
	oneYear := today.AddDate(1, 0, 0)

	groups := CertificateExpiryGroups{
		Expired:           make([]sqlc.GetCertificatesForExpiryNotificationRow, 0),
		ExpiringOneMonth:  make([]sqlc.GetCertificatesForExpiryNotificationRow, 0),
		ExpiringSixMonths: make([]sqlc.GetCertificatesForExpiryNotificationRow, 0),
		ExpiringOneYear:   make([]sqlc.GetCertificatesForExpiryNotificationRow, 0),
	}

	for _, cert := range certificates {
		if cert.ExpiryDate.IsZero() {
			continue
		}

		expiry := time.Date(cert.ExpiryDate.Year(), cert.ExpiryDate.Month(), cert.ExpiryDate.Day(), 0, 0, 0, 0, time.UTC)

		if expiry.Before(today) {
			groups.Expired = append(groups.Expired, cert)
		} else if !expiry.After(oneMonth) {
			groups.ExpiringOneMonth = append(groups.ExpiringOneMonth, cert)
		} else if !expiry.After(sixMonths) {
			groups.ExpiringSixMonths = append(groups.ExpiringSixMonths, cert)
		} else if !expiry.After(oneYear) {
			groups.ExpiringOneYear = append(groups.ExpiringOneYear, cert)
		}
		// Certificates expiring > 1 year are excluded
	}

	sortGroup := func(group []sqlc.GetCertificatesForExpiryNotificationRow) {
		sort.SliceStable(group, func(i, j int) bool {
			if !group[i].ExpiryDate.Equal(group[j].ExpiryDate) {
				return group[i].ExpiryDate.Before(group[j].ExpiryDate)
			}
			return group[i].ID.String() < group[j].ID.String()
		})
	}

	sortGroup(groups.Expired)
	sortGroup(groups.ExpiringOneMonth)
	sortGroup(groups.ExpiringSixMonths)
	sortGroup(groups.ExpiringOneYear)

	return groups
}
