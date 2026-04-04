package admin

import (
	"net/http"

	"github.com/google/uuid"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/seatime"
	"github.com/adamjames870/seacert/internal/dto"
)

// Ship Types

func HandlerAdminAddShipType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsAddShipType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		st, err := seatime.CreateShipType(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Ship type created", "id", st.Id, "name", st.Name)
		handlers.RespondWithJSON(w, 201, st)
	}
}

func HandlerAdminUpdateShipType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsUpdateShipType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		st, err := seatime.UpdateShipType(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Ship type updated", "id", st.Id, "name", st.Name)
		handlers.RespondWithJSON(w, 200, st)
	}
}

func HandlerAdminDeleteShipType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			handlers.RespondWithError(w, r, 400, "Invalid ID", err)
			return
		}

		err = seatime.DeleteShipType(r.Context(), state.Repo, id)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Ship type deleted", "id", id)
		handlers.RespondWithJSON(w, 200, "Ship type deleted")
	}
}

// Voyage Types

func HandlerAdminAddVoyageType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsAddVoyageType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		vt, err := seatime.CreateVoyageType(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Voyage type created", "id", vt.Id, "name", vt.Name)
		handlers.RespondWithJSON(w, 201, vt)
	}
}

func HandlerAdminUpdateVoyageType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsUpdateVoyageType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		vt, err := seatime.UpdateVoyageType(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Voyage type updated", "id", vt.Id, "name", vt.Name)
		handlers.RespondWithJSON(w, 200, vt)
	}
}

func HandlerAdminDeleteVoyageType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			handlers.RespondWithError(w, r, 400, "Invalid ID", err)
			return
		}

		err = seatime.DeleteVoyageType(r.Context(), state.Repo, id)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Voyage type deleted", "id", id)
		handlers.RespondWithJSON(w, 200, "Voyage type deleted")
	}
}

// Period Types

func HandlerAdminAddPeriodType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsAddPeriodType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		pt, err := seatime.CreateSeatimePeriodType(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Period type created", "id", pt.Id, "name", pt.Name)
		handlers.RespondWithJSON(w, 201, pt)
	}
}

func HandlerAdminUpdatePeriodType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsUpdatePeriodType{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		pt, err := seatime.UpdateSeatimePeriodType(r.Context(), state.Repo, params)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Period type updated", "id", pt.Id, "name", pt.Name)
		handlers.RespondWithJSON(w, 200, pt)
	}
}

func HandlerAdminDeletePeriodType(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		idStr := r.PathValue("id")
		id, err := uuid.Parse(idStr)
		if err != nil {
			handlers.RespondWithError(w, r, 400, "Invalid ID", err)
			return
		}

		err = seatime.DeleteSeatimePeriodType(r.Context(), state.Repo, id)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Period type deleted", "id", id)
		handlers.RespondWithJSON(w, 200, "Period type deleted")
	}
}
