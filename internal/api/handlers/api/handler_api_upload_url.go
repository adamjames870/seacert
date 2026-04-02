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

		// 1. Get a clean version of the original filename without its extension
		originalBase := filepath.Base(params.Filename)
		ext := filepath.Ext(originalBase)
		filenameWithoutExt := strings.TrimSuffix(originalBase, ext)

		// 2. Default extension if missing
		if ext == "" {
			if strings.Contains(params.ContentType, "pdf") {
				ext = ".pdf"
			} else if strings.Contains(params.ContentType, "jpeg") || strings.Contains(params.ContentType, "jpg") {
				ext = ".jpg"
			}
		}

		// 3. Generate a compact timestamp (YYYYMMDD-HHMMSS)
		timestamp := time.Now().Format("20060102-150405")

		// 4. Combine them into the final key: users/<user-id>/certs/<filename>-<timestamp><ext>
		fileKey := fmt.Sprintf("users/%s/certs/%s-%s%s", userId.String(), filenameWithoutExt, timestamp, ext)

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
