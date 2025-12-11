package main

import (
	"net/http"

	"github.com/adamjames870/seacert/internal/database"
	"github.com/adamjames870/seacert/mapper"
	"github.com/adamjames870/seacert/models"
	"github.com/google/uuid"
)

func (state *apiState) handlerApiGetCerts(w http.ResponseWriter, r *http.Request) {

	// GET api/certificates

	// load CertTypes to dictionary

	certTypes, errCertTypes := state.db.GetCertTypes(r.Context())
	if errCertTypes != nil {
		respondWithError(w, 500, "cannot load cert types: "+errCertTypes.Error())
		return
	}

	certTypeMap := make(map[uuid.UUID]database.CertificateType)
	for _, cType := range certTypes {
		certTypeMap[cType.ID] = cType
	}

	// load certs

	certs, errCerts := state.db.GetCerts(r.Context())
	if errCerts != nil {
		respondWithError(w, 500, "cannot load certs: "+errCerts.Error())
		return
	}

	// TODO - convert to lazy loading

	apiCerts := make([]models.Certificate, 0, len(certs))
	for _, dbCert := range certs {
		apiCerts = append(apiCerts, mapper.MapCertificate(dbCert, certTypeMap[dbCert.CertTypeID]))
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

	dbCertType, errCertType := state.db.GetCertTypeFromId(r.Context(), dbCert.CertTypeID)
	if errCertType != nil {
		respondWithError(w, 500, "cannot load cert type: "+errCertType.Error())
		return
	}

	respondWithJSON(w, 200, mapper.MapCertificate(dbCert, dbCertType))

}
