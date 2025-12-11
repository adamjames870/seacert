package main

import (
	"net/http"

	"github.com/adamjames870/seacert/mapper"
	"github.com/adamjames870/seacert/models"
)

func (state *apiState) handlerApiGetCertTypes(w http.ResponseWriter, r *http.Request) {

	// GET api/cert-types

	certTypes, errCertTypes := state.db.GetCertTypes(r.Context())
	if errCertTypes != nil {
		respondWithError(w, 500, "cannot load cert types: "+errCertTypes.Error())
		return
	}

	apiCertTypes := make([]models.CertificateType, 0, len(certTypes))
	for _, cType := range certTypes {
		apiCertTypes = append(apiCertTypes, mapper.MapCertificateType(cType))
	}

	respondWithJSON(w, 200, apiCertTypes)

}
