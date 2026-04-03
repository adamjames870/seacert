package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/adamjames870/seacert/internal/api/middleware"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload any) error {
	response, errMarshal := json.Marshal(payload)
	if errMarshal != nil {
		slog.Error("Failed to marshal JSON response", "error", errMarshal)
		return errMarshal
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_, errWrite := w.Write(response)
	if errWrite != nil {
		slog.Error("Failed to write response", "error", errWrite)
		return errWrite
	}
	return nil
}

func RespondWithError(w http.ResponseWriter, r *http.Request, code int, msg string, err error) error {
	requestID := middleware.GetRequestID(r.Context())
	if code >= 500 {
		slog.Error("Internal server error", "code", code, "message", msg, "error", err, "request_id", requestID)
	} else {
		slog.Warn("Client error", "code", code, "message", msg, "error", err, "request_id", requestID)
	}
	return RespondWithJSON(w, code, map[string]string{"error": msg, "request_id": requestID})
}
