package main

func (state *apiState) LoadState() error {
	return nil
}

func (state *apiState) CreateEndpoints() error {

	// ----------- API Handlers ----------------
	state.mux.HandleFunc("GET /api/healthz", healthzHandler)

	return nil
}
