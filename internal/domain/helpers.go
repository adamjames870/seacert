package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

type SuccessionReason string

const (
	ReasonUpdated  SuccessionReason = "updated"
	ReasonReplaced SuccessionReason = "replaced"
)

func SuccessionReasonFromString(reason string) SuccessionReason {
	switch reason {
	case "updated":
		return ReasonUpdated
	case "replaced":
		return ReasonReplaced
	default:
		return ""
	}
}

func ToNullStringFromPointer(s *string) sql.NullString {
	if s == nil {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: *s, Valid: true}
}

func ToNullInt32FromPointer(i *int32) sql.NullInt32 {
	if i == nil {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: *i, Valid: true}
}

func ToNullBoolFromPointer(b *bool) sql.NullBool {
	if b == nil {
		return sql.NullBool{Valid: false}
	}
	return sql.NullBool{Bool: *b, Valid: true}
}

func ToNullString(s string) sql.NullString {
	return sql.NullString{String: s, Valid: true}
}

func ToNullInt32OrNil(v int32) sql.NullInt32 {
	if v == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: v, Valid: true}
}

func ToNullUUIDFromStringPointer(id *string) uuid.NullUUID {
	if id == nil || *id == "" {
		return uuid.NullUUID{Valid: false}
	}

	parsed, err := uuid.Parse(*id)
	if err != nil {
		return uuid.NullUUID{Valid: false}
	}

	return uuid.NullUUID{UUID: parsed, Valid: true}
}

func ToNullTimeFromStringPointer(t *string) sql.NullTime {
	if t == nil || *t == "" {
		return sql.NullTime{Valid: false}
	}

	parsed, err := time.Parse(time.RFC3339, *t)
	if err != nil {
		return sql.NullTime{Valid: false}
	}

	return sql.NullTime{Time: parsed, Valid: true}
}
