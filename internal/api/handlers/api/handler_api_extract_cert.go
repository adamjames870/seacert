package api

import (
	"io"
	"net/http"
	"strings"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/certificates"
)

func HandlerApiExtractCert(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Limit request size to 10MB
		r.Body = http.MaxBytesReader(w, r.Body, 10<<20)
		if err := r.ParseMultipartForm(10 << 20); err != nil {
			handlers.RespondWithError(w, r, 400, "File too large or invalid multipart form", err)
			return
		}

		file, header, err := r.FormFile("certificate")
		if err != nil {
			handlers.RespondWithError(w, r, 400, "Missing 'certificate' file in form-data", err)
			return
		}
		defer file.Close()

		mimeType := header.Header.Get("Content-Type")
		if !isValidMimeType(mimeType) {
			handlers.RespondWithError(w, r, 400, "Invalid file type. Only JPEG, PNG, WEBP, and PDF are supported.", nil)
			return
		}

		fileBytes, err := io.ReadAll(file)
		if err != nil {
			handlers.RespondWithError(w, r, 500, "Failed to read uploaded file", err)
			return
		}

		// Call Gemini service
		extracted, err := certificates.ExtractCertificateData(r.Context(), state.Gemini, fileBytes, mimeType)
		if err != nil {
			state.Logger.Error("Gemini extraction failed", "error", err)
			handlers.RespondWithError(w, r, 500, "Failed to extract certificate data", err)
			return
		}

		handlers.RespondWithJSON(w, 200, extracted)
	}
}

func isValidMimeType(mimeType string) bool {
	validTypes := []string{
		"image/jpeg",
		"image/png",
		"image/webp",
		"application/pdf",
	}
	for _, t := range validTypes {
		if strings.HasPrefix(mimeType, t) {
			return true
		}
	}
	return false
}
