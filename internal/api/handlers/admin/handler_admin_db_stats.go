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
		if !state.IsDev {
			handlers.RespondWithError(w, 403, "")
		}

		countCert, errCountCert := state.Queries.CountCertificates(r.Context())
		if errCountCert != nil {
			handlers.RespondWithError(w, 500, errCountCert.Error())
		}

		countCertType, errCountCertType := state.Queries.CountCertTypes(r.Context())
		if errCountCertType != nil {
			handlers.RespondWithError(w, 500, errCountCertType.Error())
		}

		countIssuers, errCountIssuers := state.Queries.CountIssuers(r.Context())
		if errCountIssuers != nil {
			handlers.RespondWithError(w, 500, errCountIssuers.Error())
		}

		countUsers, errCountUsers := state.Queries.CountUsers(r.Context())
		if errCountUsers != nil {
			handlers.RespondWithError(w, 500, errCountUsers.Error())
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
