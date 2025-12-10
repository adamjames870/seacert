package main

import (
	"net/http"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

func (state *apiState) handlerApiGetCerts(w http.ResponseWriter, r *http.Request) {

	// GET api/certificates

	certs, errCerts := state.db.GetCerts(r.Context())
	if errCerts != nil {
		respondWithError(w, 500, "cannot load certs: "+errCerts.Error())
		return
	}

	// TODO - convert to lazy loading

	apiCerts := make([]models.Certificate, 0, len(certs))
	for _, dbCert := range certs {
		apiCerts = append(apiCerts, convertCertStruct(dbCert))
	}

	respondWithJSON(w, 200, apiCerts)

}

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

	respondWithJSON(w, 200, convertCertStruct(dbCert))

}

func convertCertStruct(dbCert database.Certificate) models.Certificate {
	return models.Certificate{
		ID:         dbCert.ID,
		CreatedAt:  dbCert.CreatedAt,
		UpdatedAt:  dbCert.UpdatedAt,
		Name:       dbCert.Name,
		CertNumber: dbCert.CertNumber,
		Issuer:     dbCert.Issuer,
		IssuedDate: dbCert.IssuedDate,
	}
}
