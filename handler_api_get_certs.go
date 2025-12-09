package main

import (
	"net/http"

	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

func (state *apiState) handlerApiGetCertFromId(w http.ResponseWriter, r *http.Request) {

	// GET api/certificates/{certID}

	certId, errId := uuid.Parse(r.PathValue("certId"))
	if errId != nil {
		respondWithError(w, 400, "cannot parse cert id to uuid: "+errId.Error())
		return
	}

	dbCert, errCert := state.db.GetCertFromId(r.Context(), certId)
	if errCert != nil {
		respondWithError(w, 404, "cannot load cert: "+errCert.Error())
		return
	}

	rv := models.Certificate{
		ID:         dbCert.ID,
		CreatedAt:  dbCert.CreatedAt,
		UpdatedAt:  dbCert.UpdatedAt,
		Name:       dbCert.Name,
		CertNumber: dbCert.CertNumber,
		Issuer:     dbCert.Issuer,
		IssuedDate: dbCert.IssuedDate,
	}

	respondWithJSON(w, 200, rv)

}
