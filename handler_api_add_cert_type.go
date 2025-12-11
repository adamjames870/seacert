package main

import (
	"encoding/json"
	"net/http"

	"github.com/adamjames870/seacert/db"
	"github.com/adamjames870/seacert/mapper"
	"github.com/adamjames870/seacert/models"
)

func (state *apiState) handlerApiAddCertType(w http.ResponseWriter, r *http.Request) {

	// POST api/cert-types

	decoder := json.NewDecoder(r.Body)
	params := models.ParamsAddCertificateType{}
	errDecode := decoder.Decode(&params)
	if errDecode != nil {
		respondWithError(w, 400, "unable to decode json: "+errDecode.Error())
		return
	}

	certType, certTypeErr := db.WriteNewCertType(r.Context(), *state.db, params)
	if certTypeErr != nil {
		respondWithError(w, 500, certTypeErr.Error())
		return
	}

	rv := mapper.MapCertificateType(certType)

	respondWithJSON(w, 201, rv)

}
