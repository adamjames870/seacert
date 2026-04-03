package handlers

import (
	"errors"
	"net/http"

	"github.com/adamjames870/seacert/internal/domain"
)

func MapDomainError(err error) (int, string) {
	if errors.Is(err, domain.ErrNotFound) {
		return http.StatusNotFound, "Resource not found"
	}
	if errors.Is(err, domain.ErrUnauthorized) {
		return http.StatusUnauthorized, "Unauthorized"
	}
	if errors.Is(err, domain.ErrAlreadyExists) {
		return http.StatusConflict, "Resource already exists"
	}
	if errors.Is(err, domain.ErrInvalidInput) {
		return http.StatusBadRequest, "Invalid input"
	}
	return http.StatusInternalServerError, "Internal server error"
}
