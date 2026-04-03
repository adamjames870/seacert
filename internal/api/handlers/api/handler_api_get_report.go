package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
)

func HandlerApiGetReport(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		certs, err := getAllCertificates(state, r.Context(), userId)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		pdf, err := GenerateCertificatesReport(certs)
		if err != nil {
			handlers.RespondWithError(w, r, 500, "Error generating report", err)
			return
		}

		w.Header().Set("Content-Type", "application/pdf")
		w.Header().Set("Content-Disposition", "attachment; filename=\"seacert-report.pdf\"")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write(pdf.GetBytes())
	}
}
