package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"
)

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
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

func RespondWithError(w http.ResponseWriter, code int, msg string, err error) error {
	if code >= 500 {
		slog.Error("Internal server error", "code", code, "message", msg, "error", err)
	} else {
		slog.Warn("Client error", "code", code, "message", msg, "error", err)
	}
	return RespondWithJSON(w, code, map[string]string{"error": msg})
}
