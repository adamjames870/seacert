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
	return middleware.Cors(middleware.RequestID(middleware.Recovery(middleware.Logging(mux)))), nil
}

func createEndpoints(mux *http.ServeMux, state *internal.ApiState) error {

	authInfo := auth.Info{
		PublicKey:        os.Getenv("SUPABASE_PUBLIC_JWK"),
		ExpectedIssuer:   os.Getenv("SUPABASE_ISSUER"),
		ExpectedAudience: os.Getenv("SUPABASE_AUDIENCE"),
	}

	adapter := &userStoreAdapter{repo: state.Repo}

	authMw, errAuthMw := auth.NewAuthMiddleware(authInfo, adapter)
	if errAuthMw != nil {
		panic(errAuthMw)
	}

	adminMw := auth.RequireRole("admin")

	// ----------- ADMIN Handlers ----------------
	mux.Handle("GET /admin/healthz", admin.HandlerAdminHealthz())
	mux.Handle("POST /admin/reset", authMw(adminMw(admin.HandlerAdminReset(state))))
	mux.Handle("GET /admin/dbstats", authMw(adminMw(admin.HandlerAdminDbStats(state))))
	mux.Handle("PUT /admin/users", authMw(admin.HandlerAdminUpdateUser(state)))
	mux.Handle("GET /admin/users", authMw(admin.HandlerAdminGetUser(state)))
	mux.Handle("POST /admin/cert-types/resolve", authMw(adminMw(admin.HandlerAdminResolveCertType(state))))
	mux.Handle("POST /api/admin/ships/resolve", authMw(adminMw(admin.HandlerAdminResolveShip(state))))
	mux.Handle("POST /api/admin/ships/approve/{id}", authMw(adminMw(admin.HandlerAdminApproveShip(state))))

	// ----------- API Handlers ----------------
	mux.Handle("POST /api/certificates/upload-url", authMw(api.HandlerApiGetUploadURL(state)))
	mux.Handle("POST /api/certificates", authMw(api.HandlerApiAddCert(state)))
	mux.Handle("GET /api/certificates", authMw(api.HandlerApiGetCerts(state)))
	mux.Handle("GET /api/certificates/report", authMw(api.HandlerApiGetReport(state)))
	mux.Handle("PUT /api/certificates", authMw(api.HandlerApiUpdateCert(state)))
	mux.Handle("DELETE /api/certificates", authMw(api.HandlerApiDeleteCert(state)))

	mux.Handle("GET /api/cert-types", authMw(api.HandlerApiGetCertTypes(state)))
	mux.Handle("POST /api/cert-types", authMw(api.HandlerApiAddCertType(state)))
	mux.Handle("PUT /api/cert-types", authMw(adminMw(api.HandlerUpdateCertType(state))))

	mux.Handle("GET /api/issuers", authMw(api.HandlerApiGetIssuers(state)))
	mux.Handle("POST /api/issuers", authMw(api.HandlerApiAddIssuer(state)))
	mux.Handle("PUT /api/issuers", authMw(api.HandlerUpdateIssuer(state)))

	mux.Handle("GET /api/seatime/lookups", authMw(api.HandlerApiGetSeatimeLookups(state)))
	mux.Handle("POST /api/seatime", authMw(api.HandlerApiAddSeatime(state)))
	mux.Handle("GET /api/seatime", authMw(api.HandlerApiListSeatime(state)))
	mux.Handle("GET /api/ships", authMw(api.HandlerApiGetShips(state)))
	mux.Handle("POST /api/ships", authMw(api.HandlerApiAddShip(state)))
	mux.Handle("PATCH /api/ships", authMw(api.HandlerApiUpdateShip(state)))

	return nil
}
