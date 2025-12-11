package main

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/db"
	"github.com/adamjames870/seacert/mapper"
	"github.com/adamjames870/seacert/models"
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

	cert, certErr := db.WriteNewCert(r.Context(), *state.db, params)
	if certErr != nil {
		respondWithError(w, 500, "error writing cert: "+certErr.Error())
		return
	}

	certType, certTypeError := state.db.GetCertTypeFromId(r.Context(), cert.CertTypeID)
	if certTypeError != nil {
		respondWithError(w, 500, "error retrieving cert type: "+certTypeError.Error())
	}

	rv := mapper.MapCertificate(cert, certType)

	respondWithJSON(w, 201, rv)

}
