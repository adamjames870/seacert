package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strings"
	"time"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/google/uuid"
)

type ParamsUploadURL struct {
	Filename    string `json:"filename" validate:"required"`
	ContentType string `json:"content-type" validate:"required"`
}

type ResponseUploadURL struct {
	UploadURL string `json:"upload-url"`
	FileKey   string `json:"file-key"`
}

func HandlerApiGetUploadURL(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if state.Storage == nil {
			handlers.RespondWithError(w, 501, "Storage not configured", nil)
			return
		}

		decoder := json.NewDecoder(r.Body)
		params := ParamsUploadURL{}
		errDecode := decoder.Decode(&params)
		if errDecode != nil {
			handlers.RespondWithError(w, 400, "Invalid request payload", errDecode)
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, 401, "Unauthorized", errId)
			return
		}

		// Generate a unique file key: users/<user_id>/certs/<uuid><ext>
		ext := filepath.Ext(params.Filename)
		if ext == "" {
			// Default extensions if missing
			if strings.Contains(params.ContentType, "pdf") {
				ext = ".pdf"
			} else if strings.Contains(params.ContentType, "jpeg") || strings.Contains(params.ContentType, "jpg") {
				ext = ".jpg"
			}
		}

		fileKey := fmt.Sprintf("users/%s/certs/%s%s", userId.String(), uuid.New().String(), ext)

		uploadURL, err := state.Storage.GetPresignedUploadURL(r.Context(), fileKey, params.ContentType, 15*time.Minute)
		if err != nil {
			handlers.RespondWithError(w, 500, "Error generating upload URL", err)
			return
		}

		handlers.RespondWithJSON(w, 200, ResponseUploadURL{
			UploadURL: uploadURL,
			FileKey:   fileKey,
		})
	}
}
