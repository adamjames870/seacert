package apiHttp

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/apiHttp/handlers"
)

func BuildRouter(state *internal.ApiState) (*http.ServeMux, error) {
	mux := http.NewServeMux()
	err := createEndpoints(mux, state)
	if err != nil {
		return nil, err
	}
	return mux, nil
}

func createEndpoints(mux *http.ServeMux, state *internal.ApiState) error {

	// ----------- ADMIN Handlers ----------------
	mux.Handle("GET /admin/healthz", handlers.HandlerAdminHealthz())
	mux.Handle("POST /admin/reset", handlers.HandlerApiReset(state))
	mux.Handle("GET /admin/dbstats", handlers.HandlerAdminDbStats(state))

	// ----------- API Handlers ----------------
	mux.Handle("POST /api/certificates", handlers.HandlerApiAddCert(state))
	mux.Handle("GET /api/certificates", handlers.HandlerApiGetCerts(state))

	mux.Handle("GET /api/cert-types", handlers.HandlerApiGetCertTypes(state))
	mux.Handle("POST /api/cert-types", handlers.HandlerApiAddCertType(state))

	mux.Handle("GET /api/issuers", handlers.HandlerApiGetIssuers(state))
	mux.Handle("POST /api/issuers", handlers.HandlerApiAddIssuer(state))

	return nil
}
