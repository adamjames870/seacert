package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

var validate = validator.New()

func DecodeAndValidate(r *http.Request, payload any) error {
	if err := json.NewDecoder(r.Body).Decode(payload); err != nil {
		return fmt.Errorf("invalid request payload: %w", err)
	}

	if err := validate.Struct(payload); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errMsgs []string
			for _, ve := range validationErrors {
				errMsgs = append(errMsgs, fmt.Sprintf("field %s failed validation: %s", ve.Field(), ve.Tag()))
			}
			return fmt.Errorf("validation failed: %s", strings.Join(errMsgs, ", "))
		}
		return err
	}

	return nil
}
