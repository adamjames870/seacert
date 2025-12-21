package tests

import (
	"database/sql"
	"testing"
	"time"

	"github.com/adamjames870/seacert/internal/domain"
	"github.com/google/uuid"
)

func TestToNullStringFromPointer(t *testing.T) {
	s := "test"
	tests := []struct {
		name  string
		input *string
		want  sql.NullString
	}{
		{"valid string", &s, sql.NullString{String: "test", Valid: true}},
		{"nil string", nil, sql.NullString{Valid: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.ToNullStringFromPointer(tt.input)
			if got != tt.want {
				t.Errorf("ToNullStringFromPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullString(t *testing.T) {
	got := domain.ToNullString("test")
	want := sql.NullString{String: "test", Valid: true}
	if got != want {
		t.Errorf("ToNullString() = %v, want %v", got, want)
	}
}

func TestToNullInt32OrNil(t *testing.T) {
	tests := []struct {
		name  string
		input int32
		want  sql.NullInt32
	}{
		{"valid int", 123, sql.NullInt32{Int32: 123, Valid: true}},
		{"zero int", 0, sql.NullInt32{Valid: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.ToNullInt32OrNil(tt.input)
			if got != tt.want {
				t.Errorf("ToNullInt32OrNil() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullUUIDFromStringPointer(t *testing.T) {
	validUUIDStr := "550e8400-e29b-41d4-a716-446655440000"
	validUUID := uuid.MustParse(validUUIDStr)
	invalidUUIDStr := "not-a-uuid"
	emptyStr := ""

	tests := []struct {
		name  string
		input *string
		want  uuid.NullUUID
	}{
		{"valid uuid", &validUUIDStr, uuid.NullUUID{UUID: validUUID, Valid: true}},
		{"invalid uuid", &invalidUUIDStr, uuid.NullUUID{Valid: false}},
		{"empty string", &emptyStr, uuid.NullUUID{Valid: false}},
		{"nil pointer", nil, uuid.NullUUID{Valid: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.ToNullUUIDFromStringPointer(tt.input)
			if got != tt.want {
				t.Errorf("ToNullUUIDFromStringPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToNullTimeFromStringPointer(t *testing.T) {
	validTimeStr := "2023-10-27T10:00:00Z"
	validTime, _ := time.Parse(time.RFC3339, validTimeStr)
	invalidTimeStr := "not-a-time"
	emptyStr := ""

	tests := []struct {
		name  string
		input *string
		want  sql.NullTime
	}{
		{"valid time", &validTimeStr, sql.NullTime{Time: validTime, Valid: true}},
		{"invalid time", &invalidTimeStr, sql.NullTime{Valid: false}},
		{"empty string", &emptyStr, sql.NullTime{Valid: false}},
		{"nil pointer", nil, sql.NullTime{Valid: false}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := domain.ToNullTimeFromStringPointer(tt.input)
			if got != tt.want {
				t.Errorf("ToNullTimeFromStringPointer() = %v, want %v", got, tt.want)
			}
		})
	}
}
