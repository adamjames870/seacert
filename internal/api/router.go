package api

import (
	"net/http"
	"os"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers/admin"
	"github.com/adamjames870/seacert/internal/api/handlers/api"
	"github.com/adamjames870/seacert/internal/api/middleware"
)

func BuildRouter(state *internal.ApiState) (http.Handler, error) {
	mux := http.NewServeMux()
	err := createEndpoints(mux, state)
	if err != nil {
		return nil, err
	}
	return middleware.Cors(middleware.Logging(mux)), nil
}

func createEndpoints(mux *http.ServeMux, state *internal.ApiState) error {

	authInfo := auth.Info{
		PublicKey:        os.Getenv("SUPABASE_PUBLIC_JWK"),
		ExpectedIssuer:   os.Getenv("SUPABASE_ISSUER"),
		ExpectedAudience: os.Getenv("SUPABASE_AUDIENCE"),
	}

	adapter := &userStoreAdapter{state: state}

	authMw, errAuthMw := auth.NewAuthMiddleware(authInfo, adapter)
	if errAuthMw != nil {
		panic(errAuthMw)
	}

	// ----------- ADMIN Handlers ----------------
	mux.Handle("GET /admin/healthz", admin.HandlerAdminHealthz())
	mux.Handle("POST /admin/reset", authMw(admin.HandlerAdminReset(state)))
	mux.Handle("GET /admin/dbstats", authMw(admin.HandlerAdminDbStats(state)))
	mux.Handle("PUT /admin/users", authMw(admin.HandlerAdminUpdateUser(state)))
	mux.Handle("GET /admin/users", authMw(admin.HandlerAdminGetUser(state)))

	// ----------- API Handlers ----------------
	mux.Handle("POST /api/certificates", authMw(api.HandlerApiAddCert(state)))
	mux.Handle("GET /api/certificates", authMw(api.HandlerApiGetCerts(state)))
	mux.Handle("PUT /api/certificates", authMw(api.HandlerApiUpdateCert(state)))

	mux.Handle("GET /api/cert-types", authMw(api.HandlerApiGetCertTypes(state)))
	mux.Handle("POST /api/cert-types", authMw(api.HandlerApiAddCertType(state)))

	mux.Handle("GET /api/issuers", authMw(api.HandlerApiGetIssuers(state)))
	mux.Handle("POST /api/issuers", authMw(api.HandlerApiAddIssuer(state)))

	return nil
}
