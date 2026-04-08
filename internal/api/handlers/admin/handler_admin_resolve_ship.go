package admin

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/seatime"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerAdminResolveShip(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsResolveShip{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		err := seatime.ResolveShip(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Ship resolved", "provisional_id", params.ProvisionalId, "replacement_id", params.ReplacementId)

		handlers.RespondWithJSON(w, 200, "Ship resolved and provisional ship deleted")
	}
}

func HandlerAdminApproveShip(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			handlers.RespondWithError(w, r, 400, "Invalid ship ID", err)
			return
		}

		err = seatime.ApproveShip(r.Context(), state.Repo, id)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Ship approved", "ship_id", id)

		handlers.RespondWithJSON(w, 200, "Ship approved")
	}
}
