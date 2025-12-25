package sqlc

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type CertView struct {
	ID                    uuid.UUID
	CreatedAt             time.Time
	UpdatedAt             time.Time
	CertNumber            string
	IssuedDate            time.Time
	AlternativeName       sql.NullString
	Remarks               sql.NullString
	ManualExpiry          sql.NullTime
	CertTypeID            uuid.UUID
	CertTypeCreatedAt     time.Time
	CertTypeUpdatedAt     time.Time
	CertTypeName          string
	CertTypeShortName     string
	CertTypeStcwReference sql.NullString
	NormalValidityMonths  sql.NullInt32
	IssuerID              uuid.UUID
	IssuerCreatedAt       time.Time
	IssuerUpdatedAt       time.Time
	IssuerName            string
	IssuerCountry         sql.NullString
	IssuerWebsite         sql.NullString
	Deleted               bool
	HasSuccessor          bool
	HasPredecessors       bool
}

func (r GetCertsRow) ToCertView() CertView {
	return CertView{
		ID:                    r.ID,
		CreatedAt:             r.CreatedAt,
		UpdatedAt:             r.UpdatedAt,
		CertNumber:            r.CertNumber,
		IssuedDate:            r.IssuedDate,
		AlternativeName:       r.AlternativeName,
		Remarks:               r.Remarks,
		ManualExpiry:          r.ManualExpiry,
		CertTypeID:            r.CertTypeID,
		CertTypeCreatedAt:     r.CertTypeCreatedAt,
		CertTypeUpdatedAt:     r.CertTypeUpdatedAt,
		CertTypeName:          r.CertTypeName,
		CertTypeShortName:     r.CertTypeShortName,
		CertTypeStcwReference: r.CertTypeStcwReference,
		NormalValidityMonths:  r.NormalValidityMonths,
		IssuerID:              r.IssuerID,
		IssuerCreatedAt:       r.IssuerCreatedAt,
		IssuerUpdatedAt:       r.IssuerUpdatedAt,
		IssuerName:            r.IssuerName,
		IssuerCountry:         r.IssuerCountry,
		IssuerWebsite:         r.IssuerWebsite,
		Deleted:               r.Deleted,
		HasSuccessor:          r.HasSuccessor,
		HasPredecessors:       r.HasPredecessors,
	}
}

func (r GetCertFromIdRow) ToCertView() CertView {
	return CertView{
		ID:                    r.ID,
		CreatedAt:             r.CreatedAt,
		UpdatedAt:             r.UpdatedAt,
		CertNumber:            r.CertNumber,
		IssuedDate:            r.IssuedDate,
		AlternativeName:       r.AlternativeName,
		Remarks:               r.Remarks,
		ManualExpiry:          r.ManualExpiry,
		CertTypeID:            r.CertTypeID,
		CertTypeCreatedAt:     r.CertTypeCreatedAt,
		CertTypeUpdatedAt:     r.CertTypeUpdatedAt,
		CertTypeName:          r.CertTypeName,
		CertTypeShortName:     r.CertTypeShortName,
		CertTypeStcwReference: r.CertTypeStcwReference,
		NormalValidityMonths:  r.NormalValidityMonths,
		IssuerID:              r.IssuerID,
		IssuerCreatedAt:       r.IssuerCreatedAt,
		IssuerUpdatedAt:       r.IssuerUpdatedAt,
		IssuerName:            r.IssuerName,
		IssuerCountry:         r.IssuerCountry,
		IssuerWebsite:         r.IssuerWebsite,
		Deleted:               r.Deleted,
		HasSuccessor:          r.HasSuccessor,
		HasPredecessors:       r.HasPredecessors,
	}
}
