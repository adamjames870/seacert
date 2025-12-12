package handlers

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
)

func HandlerApiGetCertTypes(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/cert-types

		certTypes, errCertTypes := state.Queries.GetCertTypes(r.Context())
		if errCertTypes != nil {
			respondWithError(w, 500, "cannot load cert types: "+errCertTypes.Error())
			return
		}

		apiCertTypes := make([]cert_types.CertificateType, 0, len(certTypes))
		for _, cType := range certTypes {
			apiCertTypes = append(apiCertTypes, cert_types.MapCertificateType(cType))
		}

		respondWithJSON(w, 200, apiCertTypes)

	}
}
