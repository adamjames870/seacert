package handlers

import (
	"context"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/certificates"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiGetCerts(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		// GET api/certificates
		// GET api/certificates?id=<uuid>

		// TODO - convert to lazy loading

		idParam := r.URL.Query().Get("id")

		if idParam == "" {
			rv, err := getAllCertificates(state, r.Context())
			if err != nil {
				respondWithError(w, 500, err.Error())
				return
			}
			respondWithJSON(w, 200, rv)
			return
		}

		if idParam != "" {
			rv, err := certificates.GetCertificateFromId(state, r.Context(), idParam)
			if err != nil {
				respondWithError(w, 404, err.Error())
			}
			respondWithJSON(w, 200, certificates.MapCertificateDomainToDto(rv))
			return
		}

	}
}

func getAllCertificates(state *internal.ApiState, ctx context.Context) ([]dto.Certificate, error) {

	certs, errCerts := certificates.GetCertificates(state, ctx)
	if errCerts != nil {
		return nil, errCerts
	}

	rv := make([]dto.Certificate, 0, len(certs))
	for _, cert := range certs {
		rv = append(rv, certificates.MapCertificateDomainToDto(cert))
	}

	return rv, nil

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
