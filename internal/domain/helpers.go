package domain

import (
	"database/sql"
	"time"

	"github.com/google/uuid"
)

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

	// Try RFC3339 first
	parsed, err := time.Parse(time.RFC3339, *t)
	if err == nil {
		return sql.NullTime{Time: parsed, Valid: true}
	}

	// Try ISO Date (YYYY-MM-DD)
	parsed, err = time.Parse("2006-01-02", *t)
	if err == nil {
		// Set to start of day in UTC
		parsed = time.Date(parsed.Year(), parsed.Month(), parsed.Day(), 0, 0, 0, 0, time.UTC)
		return sql.NullTime{Time: parsed, Valid: true}
	}

	return sql.NullTime{Valid: false}
}
