package main

func (state *apiState) CreateEndpoints() error {

	// ----------- ADMIN Handlers ----------------
	state.mux.HandleFunc("GET /admin/healthz", healthzHandler)
	state.mux.HandleFunc("POST /admin/reset", state.handlerApiReset)

	// ----------- API Handlers ----------------
	state.mux.HandleFunc("POST /api/certificates", state.handlerApiAddCert)
	state.mux.HandleFunc("GET /api/certificates", state.handlerApiGetCerts)
	state.mux.HandleFunc("GET /api/certificates/{certId}", state.handlerApiGetCertFromId)
	state.mux.HandleFunc("GET /api/cert-types", state.handlerApiGetCertTypes)
	state.mux.HandleFunc("POST /api/cert-types", state.handlerApiAddCertType)

	return nil
}
