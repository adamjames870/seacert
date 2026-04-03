package handlers

import (
	"errors"
	"net/http"
	"testing"

	"github.com/adamjames870/seacert/internal/domain"
)

func TestMapDomainError(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		wantCode int
		wantMsg  string
	}{
		{
			name:     "NotFound",
			err:      domain.ErrNotFound,
			wantCode: http.StatusNotFound,
			wantMsg:  "Resource not found",
		},
		{
			name:     "Unauthorized",
			err:      domain.ErrUnauthorized,
			wantCode: http.StatusUnauthorized,
			wantMsg:  "Unauthorized",
		},
		{
			name:     "AlreadyExists",
			err:      domain.ErrAlreadyExists,
			wantCode: http.StatusConflict,
			wantMsg:  "Resource already exists",
		},
		{
			name:     "InvalidInput",
			err:      domain.ErrInvalidInput,
			wantCode: http.StatusBadRequest,
			wantMsg:  "Invalid input",
		},
		{
			name:     "InternalError",
			err:      errors.New("random error"),
			wantCode: http.StatusInternalServerError,
			wantMsg:  "Internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			code, msg := MapDomainError(tt.err)
			if code != tt.wantCode {
				t.Errorf("MapDomainError() code = %v, want %v", code, tt.wantCode)
			}
			if msg != tt.wantMsg {
				t.Errorf("MapDomainError() msg = %v, want %v", msg, tt.wantMsg)
			}
		})
	}
}
