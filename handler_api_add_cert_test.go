package main

import (
	"testing"

	"github.com/adamjames870/seacert/models"
	"github.com/go-playground/validator/v10"
)

func TestParamsAddCertificate_ValidationSuccess(t *testing.T) {
	validate := validator.New()

	params := models.ParamsAddCertificate{
		Name:       "ISO 9001",
		CertNumber: "ABC123",
		Issuer:     "Lloyds",
		IssuedDate: "2020-01-01",
	}

	err := validate.Struct(params)
	if err != nil {
		t.Fatalf("expected validation to pass, got: %v", err)
	}
}

func TestParamsAddCertificate_ValidationFailure(t *testing.T) {
	validate := validator.New()

	// Missing required fields
	params := models.ParamsAddCertificate{}

	err := validate.Struct(params)
	if err == nil {
		t.Fatalf("expected validation to fail, got nil")
	}

	ves := err.(validator.ValidationErrors)
	if len(ves) != 4 {
		t.Fatalf("expected 4 validation errors, got %d", len(ves))
	}

	expectedFields := []string{"Name", "CertNumber", "Issuer", "IssuedDate"}
	for _, f := range expectedFields {
		found := false
		for _, ve := range ves {
			if ve.Field() == f && ve.Tag() == "required" {
				found = true
			}
		}
		if !found {
			t.Errorf("expected validation error for field %s", f)
		}
	}
}
