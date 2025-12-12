package handlers

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/google/uuid"
)

func HandlerApiGetCerts(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/certificates

		// load CertTypes to dictionary

		certTypes, errCertTypes := state.Queries.GetCertTypes(r.Context())
		if errCertTypes != nil {
			respondWithError(w, 500, "cannot load cert types: "+errCertTypes.Error())
			return
		}

		certTypeMap := make(map[uuid.UUID]sqlc.CertificateType)
		for _, cType := range certTypes {
			certTypeMap[cType.ID] = cType
		}

		// load certs

		certs, errCerts := state.Queries.GetCerts(r.Context())
		if errCerts != nil {
			respondWithError(w, 500, "cannot load certs: "+errCerts.Error())
			return
		}

		// TODO - convert to lazy loading

		apiCerts := make([]certificates.Certificate, 0, len(certs))
		for _, dbCert := range certs {
			apiCerts = append(apiCerts, certificates.MapCertificate(dbCert, certTypeMap[dbCert.CertTypeID]))
		}

		respondWithJSON(w, 200, apiCerts)

	}
}
func HandlerApiGetCertFromId(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/certificates/{certID}

		certId, errId := uuid.Parse(r.PathValue("certId"))
		if errId != nil {
			respondWithError(w, 400, "cannot parse cert id to uuid: "+errId.Error())
			return
		}

		dbCert, errCert := state.Queries.GetCertFromId(r.Context(), certId)
		if errCert != nil {
			respondWithError(w, 404, "cannot load cert: "+errCert.Error())
			return
		}

		dbCertType, errCertType := state.Queries.GetCertTypeFromId(r.Context(), dbCert.CertTypeID)
		if errCertType != nil {
			respondWithError(w, 500, "cannot load cert type: "+errCertType.Error())
			return
		}

		respondWithJSON(w, 200, certificates.MapCertificate(dbCert, dbCertType))

	}
}
