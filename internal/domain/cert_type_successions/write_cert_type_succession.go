package cert_type_successions

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
	"github.com/adamjames870/seacert/internal/dto"
	"github.com/google/uuid"
)

func WriteNewCertTypeSuccession(state *internal.ApiState, ctx context.Context, params dto.ParamsAddCertTypeSuccession) (CertTypeSuccession, error) {

	idReplacingType, errReplacingType := uuid.Parse(params.ReplacingType)
	if errReplacingType != nil {
		return CertTypeSuccession{}, errReplacingType
	}

	idReplaceableType, errReplaceableType := uuid.Parse(params.ReplaceableType)
	if errReplaceableType != nil {
		return CertTypeSuccession{}, errReplaceableType
	}

	reason, errReason := sqlc.SuccessionReasonFromString(params.ReplaceReason)
	if errReason != nil {
		return CertTypeSuccession{}, errReason
	}

	newSuccession := sqlc.CreateTypeSuccessionParams{
		ID:                  uuid.New(),
		ReplacingCertType:   idReplacingType,
		ReplaceableCertType: idReplaceableType,
		ReplaceReason:       reason,
	}

	dbSuccession, errWriteNewSuccession := state.Queries.CreateTypeSuccession(ctx, newSuccession)
	if errWriteNewSuccession != nil {
		return CertTypeSuccession{}, errWriteNewSuccession
	}

	replacing, errReplacing := cert_types.GetCertTypeFromId(state, ctx, idReplacingType.String())
	if errReplacing != nil {
		return CertTypeSuccession{}, errReplacing
	}

	replaceable, errReplaceable := cert_types.GetCertTypeFromId(state, ctx, idReplaceableType.String())
	if errReplaceable != nil {
		return CertTypeSuccession{}, errReplaceable
	}

	apiSuccession := MapSuccessionDbToDomain(dbSuccession, replacing, replaceable)

	return apiSuccession, nil

}
