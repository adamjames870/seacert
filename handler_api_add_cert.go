package main

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

func (state *apiState) handlerApiAddCert(w http.ResponseWriter, r *http.Request) {
	// POST api/certificates
	decoder := json.NewDecoder(r.Body)
	params := models.ParamsAddCertificate{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	cert, certErr := writeNewCert(r.Context(), *state.db, params)
	if certErr != nil {
		respondWithError(w, 500, certErr.Error())
		return
	}

	rv := models.Certificate{
		ID:         cert.ID,
		CreatedAt:  cert.CreatedAt,
		UpdatedAt:  cert.UpdatedAt,
		Name:       cert.Name,
		CertNumber: cert.CertNumber,
		Issuer:     cert.Issuer,
		IssuedDate: cert.IssuedDate,
	}

	respondWithJSON(w, 201, rv)

}

func writeNewCert(context context.Context, db database.Queries, params models.ParamsAddCertificate) (database.Certificate, error) {

	issuedDate, errParse := time.Parse("2006-01-02", params.IssuedDate)
	if errParse != nil {
		return database.Certificate{}, errParse
	}

	newCert := database.CreateCertParams{
		ID:         uuid.New(),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
		Name:       params.Name,
		CertNumber: params.CertNumber,
		Issuer:     params.Issuer,
		IssuedDate: issuedDate,
	}

	return db.CreateCert(context, newCert)
}
