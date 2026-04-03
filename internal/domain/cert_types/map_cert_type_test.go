package cert_types

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func TestMapCertificateType(t *testing.T) {
	id := uuid.New()
	createdBy := uuid.New()
	now := time.Now().UTC()

	dbCertType := sqlc.CertificateType{
		ID:        id,
		CreatedAt: now,
		UpdatedAt: now,
		Name:      "Basic Safety Training",
		ShortName: "BST",
		StcwReference: sql.NullString{
			String: "VI/1",
			Valid:  true,
		},
		NormalValidityMonths: sql.NullInt32{
			Int32: 60,
			Valid: true,
		},
		Status: "approved",
		CreatedBy: uuid.NullUUID{
			UUID:  createdBy,
			Valid: true,
		},
	}

	domainCertType := CertificateType{
		Id:                   id,
		CreatedAt:            now,
		UpdatedAt:            now,
		Name:                 "Basic Safety Training",
		ShortName:            "BST",
		StcwReference:        "VI/1",
		NormalValidityMonths: 60,
		Status:               "approved",
		CreatedBy:            createdBy,
	}

	cbStr := createdBy.String()
	dtoCertType := dto.CertificateType{
		Id:                   id.String(),
		CreatedAt:            now,
		UpdatedAt:            now,
		Name:                 "Basic Safety Training",
		ShortName:            "BST",
		StcwRef:              "VI/1",
		NormalValidityMonths: 60,
		Status:               "approved",
		CreatedBy:            &cbStr,
	}

	t.Run("MapCertificateTypeDbToDomain", func(t *testing.T) {
		got := MapCertificateTypeDbToDomain(dbCertType)
		if got != domainCertType {
			t.Errorf("expected %+v, got %+v", domainCertType, got)
		}
	})

	t.Run("MapCertificateTypeDomainToDto", func(t *testing.T) {
		got := MapCertificateTypeDomainToDto(domainCertType)
		if got.Id != dtoCertType.Id ||
			got.Name != dtoCertType.Name ||
			got.ShortName != dtoCertType.ShortName ||
			got.StcwRef != dtoCertType.StcwRef ||
			got.NormalValidityMonths != dtoCertType.NormalValidityMonths ||
			got.Status != dtoCertType.Status ||
			*got.CreatedBy != *dtoCertType.CreatedBy {
			t.Errorf("expected %+v, got %+v", dtoCertType, got)
		}
	})

	t.Run("MapCertificateTypeDtoToDomain", func(t *testing.T) {
		got := MapCertificateTypeDtoToDomain(dtoCertType)
		if got != domainCertType {
			t.Errorf("expected %+v, got %+v", domainCertType, got)
		}
	})

	t.Run("MapCertificateTypeDomainToDb", func(t *testing.T) {
		got := MapCertificateTypeDomainToDb(domainCertType)
		if got.ID != dbCertType.ID ||
			got.Name != dbCertType.Name ||
			got.ShortName != dbCertType.ShortName ||
			got.StcwReference.String != dbCertType.StcwReference.String ||
			got.NormalValidityMonths.Int32 != dbCertType.NormalValidityMonths.Int32 ||
			got.Status != dbCertType.Status ||
			got.CreatedBy.UUID != dbCertType.CreatedBy.UUID {
			t.Errorf("expected %+v, got %+v", dbCertType, got)
		}
	})
}

func TestSuccessionReasonDbToDomain(t *testing.T) {
	tests := []struct {
		reason  sqlc.SuccessionReason
		want    SuccessionReason
		wantErr bool
	}{
		{sqlc.SuccessionReplaced, ReasonReplaced, false},
		{sqlc.SuccessionUpdated, ReasonUpdated, false},
		{"invalid", "", true},
	}

	for _, tt := range tests {
		t.Run(string(tt.reason), func(t *testing.T) {
			got, err := SuccessionReasonDbToDomain(tt.reason)
			if (err != nil) != tt.wantErr {
				t.Errorf("SuccessionReasonDbToDomain() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("SuccessionReasonDbToDomain() = %v, want %v", got, tt.want)
			}
		})
	}
}
