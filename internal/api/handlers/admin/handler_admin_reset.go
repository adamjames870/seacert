package admin

import (
	"fmt"
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
)

func HandlerAdminReset(state *internal.ApiState) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		if !state.IsDev {
			handlers.RespondWithError(w, 403, "")
		}

		errResetCerts := state.Queries.ResetCerts(r.Context())
		if errResetCerts != nil {
			fmt.Println(errResetCerts)
			handlers.RespondWithError(w, 500, errResetCerts.Error())
		}

		errResetCertTypes := state.Queries.ResetCertTypes(r.Context())
		if errResetCertTypes != nil {
			fmt.Println(errResetCertTypes)
			handlers.RespondWithError(w, 500, errResetCertTypes.Error())
		}

		errResetIssuers := state.Queries.ResetIssuers(r.Context())
		if errResetIssuers != nil {
			fmt.Println(errResetIssuers)
			handlers.RespondWithError(w, 500, errResetIssuers.Error())
		}

		errResetUsers := state.Queries.ResetUsers(r.Context())
		if errResetUsers != nil {
			fmt.Println(errResetUsers)
			handlers.RespondWithError(w, 500, errResetUsers.Error())
		}

		handlers.RespondWithJSON(w, 200, map[string]string{"message": "db reset"})

	}
}
