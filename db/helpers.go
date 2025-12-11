package db

import (
	"database/sql"

	"github.com/google/uuid"
)

func toNullString(s string) sql.NullString {
	if s == "" {
		return sql.NullString{Valid: false}
	}
	return sql.NullString{String: s, Valid: true}
}

func toNullInt32OrNil(v int32) sql.NullInt32 {
	if v == 0 {
		return sql.NullInt32{Valid: false}
	}
	return sql.NullInt32{Int32: v, Valid: true}
}

func toNullUUID(id *uuid.UUID) uuid.NullUUID {
	if id == nil {
		return uuid.NullUUID{
			Valid: false,
		}
	}
	return uuid.NullUUID{
		UUID:  *id,
		Valid: true,
	}
}
