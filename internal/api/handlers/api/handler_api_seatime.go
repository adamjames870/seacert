package api

import (
	"net/http"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/api/auth"
	"github.com/adamjames870/seacert/internal/api/handlers"
	"github.com/adamjames870/seacert/internal/domain/seatime"
	"github.com/adamjames870/seacert/internal/dto"
)

func HandlerApiAddSeatime(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsAddSeatime{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		st, err := seatime.CreateSeatime(r.Context(), state.Repo, params, userId, auth.IsAdmin(r.Context()))
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		state.Logger.Info("Seatime created", "user_id", userId, "seatime_id", st.Id)
		rv := seatime.MapSeatimeDomainToDto(st)

		handlers.RespondWithJSON(w, 201, rv)
	}
}

func HandlerApiListSeatime(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		sts, err := seatime.GetSeatime(r.Context(), state.Repo, userId)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		var rv []dto.Seatime
		for _, st := range sts {
			rv = append(rv, seatime.MapSeatimeDomainToDto(st))
		}

		handlers.RespondWithJSON(w, 200, rv)
	}
}

func HandlerApiGetSeatimeLookups(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		lookups, err := seatime.GetSeatimeLookups(r.Context(), state.Repo)
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		handlers.RespondWithJSON(w, 200, lookups)
	}
}

func HandlerApiGetShips(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		ships, err := seatime.GetShips(r.Context(), state.Repo, &userId, auth.IsAdmin(r.Context()))
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		var rv []dto.Ship
		for _, s := range ships {
			rv = append(rv, seatime.MapShipToDto(s))
		}

		handlers.RespondWithJSON(w, 200, rv)
	}
}

func HandlerApiAddShip(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsAddShip{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		s, err := seatime.CreateShipStandalone(r.Context(), state.Repo, params, userId, auth.IsAdmin(r.Context()))
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		handlers.RespondWithJSON(w, 201, seatime.MapShipToDto(s))
	}
}

func HandlerApiUpdateShip(state *internal.ApiState) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		params := dto.ParamsUpdateShip{}
		if err := handlers.DecodeAndValidate(r, &params); err != nil {
			handlers.RespondWithError(w, r, 400, err.Error(), err)
			return
		}

		userId, errId := auth.UserIdFromContext(r.Context())
		if errId != nil {
			handlers.RespondWithError(w, r, 401, "Unauthorized", errId)
			return
		}

		s, err := seatime.UpdateShip(r.Context(), state.Repo, params, userId, auth.IsAdmin(r.Context()))
		if err != nil {
			code, msg := handlers.MapDomainError(err)
			handlers.RespondWithError(w, r, code, msg, err)
			return
		}

		handlers.RespondWithJSON(w, 200, seatime.MapShipToDto(s))
	}
}
