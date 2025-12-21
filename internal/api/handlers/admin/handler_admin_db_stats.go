package admin

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerAdminDbStats(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		countCert, errCountCert := state.Queries.CountCertificates(r.Context())
		if errCountCert != nil {
			handlers.RespondWithError(w, 500, "Error counting certificates", errCountCert)
			return
		}

		countCertType, errCountCertType := state.Queries.CountCertTypes(r.Context())
		if errCountCertType != nil {
			handlers.RespondWithError(w, 500, "Error counting certificate types", errCountCertType)
			return
		}

		countIssuers, errCountIssuers := state.Queries.CountIssuers(r.Context())
		if errCountIssuers != nil {
			handlers.RespondWithError(w, 500, "Error counting issuers", errCountIssuers)
			return
		}

		countUsers, errCountUsers := state.Queries.CountUsers(r.Context())
		if errCountUsers != nil {
			handlers.RespondWithError(w, 500, "Error counting users", errCountUsers)
			return
		}

		user, ok := auth.UserFromContext(r.Context())
		if !ok {
			http.Error(w, "user not found in context", http.StatusUnauthorized)
			return
		}

		rv := dto.DbStats{
			CountCert:     int(countCert),
			CountCertType: int(countCertType),
			CountIssuer:   int(countIssuers),
			CountUsers:    int(countUsers),
			UserId:        user.Id,
			UserEmail:     user.Email,
		}

		handlers.RespondWithJSON(w, 200, rv)
	}

}
