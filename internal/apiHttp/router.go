package apiHttp

import (
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/apiHttp/auth"
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

	authInfo := auth.Info{
		JwksUrl:          os.Getenv("SUPABASE_JWKS_URL"),
		ApiKey:           os.Getenv("SUPABASE_PUBLIC_JWK"),
		ExpectedIssuer:   os.Getenv("SUPABASE_ISSUER"),
		ExpectedAudience: os.Getenv("SUPABASE_AUDIENCE"),
	}

	authMw, errAuthMw := auth.NewAuthMiddleware(authInfo, state)
	if errAuthMw != nil {
		panic(errAuthMw)
	}

	// ----------- ADMIN Handlers ----------------
	mux.Handle("GET /admin/healthz", handlers.HandlerAdminHealthz())
	mux.Handle("POST /admin/reset", authMw(handlers.HandlerApiReset(state)))
	mux.Handle("GET /admin/dbstats", authMw(handlers.HandlerAdminDbStats(state)))

	// ----------- API Handlers ----------------
	mux.Handle("POST /api/certificates", authMw(handlers.HandlerApiAddCert(state)))
	mux.Handle("GET /api/certificates", authMw(handlers.HandlerApiGetCerts(state)))

	mux.Handle("GET /api/cert-types", authMw(handlers.HandlerApiGetCertTypes(state)))
	mux.Handle("POST /api/cert-types", authMw(handlers.HandlerApiAddCertType(state)))

	mux.Handle("GET /api/issuers", authMw(handlers.HandlerApiGetIssuers(state)))
	mux.Handle("POST /api/issuers", authMw(handlers.HandlerApiAddIssuer(state)))

	return nil
}
