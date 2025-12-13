package handlers

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiGetCerts(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/certificates

		// load CertTypes to dictionary

		certs, errCerts := certificates.GetCertificates(state, r.Context())
		if errCerts != nil {
			respondWithError(w, 500, "cannot load certs: "+errCerts.Error())
			return
		}

		// TODO - convert to lazy loading

		rv := make([]dto.Certificate, 0, len(certs))
		for _, cert := range certs {
			rv = append(rv, certificates.MapCertificateDomainToDto(cert))
		}

		respondWithJSON(w, 200, rv)

	}
}
func HandlerApiGetCertFromId(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/certificates/{certID}

		certId := r.PathValue("certId")

		cert, errCert := certificates.GetCertificateFromId(state, r.Context(), certId)
		if errCert != nil {
			respondWithError(w, 404, "cannot load cert: "+errCert.Error())
			return
		}

		rv := certificates.MapCertificateDomainToDto(cert)

		respondWithJSON(w, 200, rv)

	}
}
