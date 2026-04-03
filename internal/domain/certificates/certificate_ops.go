package certificates

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log/slog"
	"time"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/adamjames870/seacert/internal/storage"
	"github.com/google/uuid"
)

func GetCertificates(ctx context.Context, repo domain.Repository, userId uuid.UUID) ([]Certificate, error) {
	certs, err := repo.GetCerts(ctx, userId)
	if err != nil {
		return nil, err
	}

	apiCerts := make([]Certificate, 0, len(certs))
	for _, cert := range certs {
		thisCert := MapCertificateViewDbToDomain(cert.ToCertView())
		if cert.HasPredecessors {
			predecessorIds, err := repo.GetPredecessors(ctx, thisCert.Id)
			if err != nil {
				return nil, err
			}
			predecessors, err := GetPredecessorsFromListOfIds(ctx, repo, predecessorIds, userId)
			if err != nil {
				return nil, err
			}
			thisCert.Predecessors = predecessors
		}
		apiCerts = append(apiCerts, thisCert)
	}

	return apiCerts, nil
}

func GetCertificateById(ctx context.Context, repo domain.Repository, certId uuid.UUID, userId uuid.UUID) (Certificate, error) {
	params := sqlc.GetCertFromIdParams{
		ID:     certId,
		UserID: userId,
	}

	dbCert, err := repo.GetCertFromId(ctx, params)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Certificate{}, domain.ErrNotFound
		}
		return Certificate{}, err
	}

	certView := dbCert.ToCertView()
	apiCert := MapCertificateViewDbToDomain(certView)

	if certView.HasPredecessors {
		predecessorIds, err := repo.GetPredecessors(ctx, apiCert.Id)
		if err != nil {
			return Certificate{}, err
		}
		predecessors, err := GetPredecessorsFromListOfIds(ctx, repo, predecessorIds, userId)
		if err != nil {
			return Certificate{}, err
		}
		apiCert.Predecessors = predecessors
	}

	return apiCert, nil
}

func GetPredecessorsFromListOfIds(ctx context.Context, repo domain.Repository, predecessors []sqlc.GetPredecessorsRow, userId uuid.UUID) ([]Predecesor, error) {
	certs := make([]Predecesor, 0, len(predecessors))
	for _, predecessor := range predecessors {
		oldCertUuid, _ := uuid.Parse(predecessor.OldCert.String())
		cert, err := GetCertificateById(ctx, repo, oldCertUuid, userId)
		if err != nil {
			return nil, err
		}
		replaceReason, err := cert_types.SuccessionReasonDbToDomain(predecessor.Reason)
		if err != nil {
			return nil, err
		}
		certs = append(certs, Predecesor{
			Cert:          cert,
			ReplaceReason: replaceReason,
		})
	}
	return certs, nil
}

func CreateCertificate(ctx context.Context, repo domain.Repository, params dto.ParamsAddCertificate, userId uuid.UUID) (Certificate, error) {
	issuedDateNull := domain.ToNullTimeFromStringPointer(&params.IssuedDate)
	if !issuedDateNull.Valid {
		return Certificate{}, domain.ErrInvalidInput
	}
	issuedDate := issuedDateNull.Time

	certTypeUuid, _ := uuid.Parse(params.CertTypeId)
	issuerUuid, _ := uuid.Parse(params.IssuerId)

	newCert := sqlc.CreateCertParams{
		ID:              uuid.New(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
		UserID:          userId,
		CertNumber:      params.CertNumber,
		IssuedDate:      issuedDate,
		CertTypeID:      certTypeUuid,
		IssuerID:        issuerUuid,
		AlternativeName: domain.ToNullStringFromPointer(params.AlternativeName),
		Remarks:         domain.ToNullStringFromPointer(params.Remarks),
		ManualExpiry:    domain.ToNullTimeFromStringPointer(params.ManualExpiry),
		DocumentPath:    domain.ToNullStringFromPointer(params.DocumentPath),
	}

	var dbCert sqlc.Certificate
	err := repo.WithTx(ctx, func(txRepo domain.Repository) error {
		var err error
		dbCert, err = txRepo.CreateCert(ctx, newCert)
		if err != nil {
			return err
		}

		if params.Supersedes != nil && params.SupersedeReason != nil {
			oldCertUuid, _ := uuid.Parse(*params.Supersedes)
			_, err := txRepo.CreateSuccession(ctx, sqlc.CreateSuccessionParams{
				ID:      uuid.New(),
				OldCert: oldCertUuid,
				NewCert: dbCert.ID,
				Reason:  sqlc.SuccessionReason(*params.SupersedeReason),
			})
			if err != nil {
				return err
			}
		}
		return nil
	})

	if err != nil {
		return Certificate{}, err
	}

	return GetCertificateById(ctx, repo, dbCert.ID, userId)
}

func UpdateCertificate(ctx context.Context, repo domain.Repository, store storage.Storage, logger *slog.Logger, params dto.ParamsUpdateCertificate, userId uuid.UUID) (Certificate, error) {
	certUuid, errParse := uuid.Parse(params.Id)
	if errParse != nil {
		return Certificate{}, domain.ErrInvalidInput
	}

	oldCert, err := GetCertificateById(ctx, repo, certUuid, userId)
	if err != nil {
		return Certificate{}, err
	}

	var docPath *string
	if params.DocumentPath != nil {
		if string(params.DocumentPath) == "null" {
			docPath = new(string)
			*docPath = ""
		} else {
			var pathStr string
			if err := json.Unmarshal(params.DocumentPath, &pathStr); err != nil {
				return Certificate{}, err
			}
			docPath = &pathStr
		}

		if oldCert.DocumentPath != "" && *docPath != "" && oldCert.DocumentPath != *docPath && store != nil {
			errDelete := store.DeleteObject(ctx, oldCert.DocumentPath)
			if errDelete != nil {
				logger.Error("Failed to delete old certificate from storage", "path", oldCert.DocumentPath, "error", errDelete)
			}
		}
	}

	updateParams := sqlc.UpdateCertificateParams{
		ID:              certUuid,
		CertNumber:      domain.ToNullStringFromPointer(params.CertNumber),
		IssuedDate:      domain.ToNullTimeFromStringPointer(params.IssuedDate),
		CertTypeID:      domain.ToNullUUIDFromStringPointer(params.CertTypeId),
		IssuerID:        domain.ToNullUUIDFromStringPointer(params.IssuerId),
		AlternativeName: domain.ToNullStringFromPointer(params.AlternativeName),
		Remarks:         domain.ToNullStringFromPointer(params.Remarks),
		DocumentPath:    domain.ToNullStringFromPointer(docPath),
		Deleted:         domain.ToNullBoolFromPointer(params.Deleted),
	}

	if params.ManualExpiry != nil {
		updateParams.ManualExpiryDoUpdate = true
		if string(params.ManualExpiry) == "null" {
			updateParams.ManualExpiry = sql.NullTime{Valid: false}
		} else {
			var expiryStr string
			if err := json.Unmarshal(params.ManualExpiry, &expiryStr); err != nil {
				return Certificate{}, err
			}
			updateParams.ManualExpiry = domain.ToNullTimeFromStringPointer(&expiryStr)
		}
	}

	dbCert, err := repo.UpdateCertificate(ctx, updateParams)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Certificate{}, domain.ErrNotFound
		}
		return Certificate{}, err
	}

	return GetCertificateById(ctx, repo, dbCert.ID, userId)
}

func DeleteCertificate(ctx context.Context, repo domain.Repository, store storage.Storage, logger *slog.Logger, certId uuid.UUID, userId uuid.UUID) error {
	cert, err := GetCertificateById(ctx, repo, certId, userId)
	if err != nil {
		return err
	}

	if cert.DocumentPath != "" && store != nil {
		err := store.DeleteObject(ctx, cert.DocumentPath)
		if err != nil {
			logger.Error("Failed to delete certificate document", "path", cert.DocumentPath, "error", err)
		} else {
			logger.Info("Deleted certificate document", "path", cert.DocumentPath)
		}
	}

	return repo.DeleteCert(ctx, sqlc.DeleteCertParams{
		ID:     certId,
		UserID: userId,
	})
}

func (cert *Certificate) calculateExpiryDate() {
	if !cert.ManualExpiry.IsZero() {
		cert.ExpiryDate = cert.ManualExpiry
	} else if cert.CertType.NormalValidityMonths != 0 {
		cert.ExpiryDate = getExpiryAfterValidity(cert.IssuedDate, int(cert.CertType.NormalValidityMonths))
	} else {
		cert.ExpiryDate = time.Time{}
	}
}

func getExpiryAfterValidity(issueDate time.Time, validityMonths int) time.Time {
	target := issueDate.AddDate(0, validityMonths, 0)
	expectedMonth := (int(issueDate.Month()) + validityMonths - 1) % 12
	if expectedMonth < 0 {
		expectedMonth += 12
	}
	expectedMonth += 1
	if int(target.Month()) != expectedMonth {
		target = time.Date(target.Year(), target.Month(), 1, 0, 0, 0, 0, target.Location()).AddDate(0, 0, -1)
	}
	return time.Date(target.Year(), target.Month(), target.Day(), 0, 0, 0, 0, target.Location()).AddDate(0, 0, -1)
}

func daysInMonth(year int, month time.Month) int {
	// 1st of next month, minus one day
	firstOfNext := time.Date(year, month+1, 1, 0, 0, 0, 0, time.UTC)
	lastOfThis := firstOfNext.AddDate(0, 0, -1)
	return lastOfThis.Day()
}
