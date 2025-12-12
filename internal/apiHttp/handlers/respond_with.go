package handlers

import (
	"encoding/json"
	"net/http"
)

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) error {
	response, errMarshal := json.Marshal(payload)
	if errMarshal != nil {
		return errMarshal
	}
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(code)
	_, errWrite := w.Write(response)
	if errWrite != nil {
		return errWrite
	}
	return nil
}

func respondWithError(w http.ResponseWriter, code int, msg string) error {
	return respondWithJSON(w, code, map[string]string{"error": msg})
}
