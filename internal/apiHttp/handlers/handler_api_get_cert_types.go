package handlers

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiGetCertTypes(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/cert-types

		certTypes, errCertTypes := cert_types.GetCertTypes(state, r.Context())
		if errCertTypes != nil {
			respondWithError(w, 500, "Unable to get cert types: "+errCertTypes.Error())
		}

		rv := make([]dto.CertificateType, 0, len(certTypes))
		for _, certType := range certTypes {
			rv = append(rv, cert_types.MapCertificateTypeDomainToDto(certType))
		}

		respondWithJSON(w, 200, rv)

	}
}
