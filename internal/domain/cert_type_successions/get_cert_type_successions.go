package cert_type_successions

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/domain/cert_types"
)

func GetAllReplaceableByThisId(state *internal.ApiState, ctx context.Context, certType cert_types.CertificateType) (CertTypeSuccessions, error) {

	dbSuccessions, errSuccessions := state.Queries.GetAllReplaceableByMe(ctx, certType.Id)
	if errSuccessions != nil {
		return CertTypeSuccessions{}, errSuccessions
	}

	apiSuccessions := CertTypeSuccessions{
		CertType:   certType,
		CanReplace: make([]Succession, len(dbSuccessions)),
	}

	for _, succession := range dbSuccessions {
		reason, errReason := cert_types.SuccessionReasonDbToDomain(succession.ReplaceReason)
		if errReason != nil {
			return CertTypeSuccessions{}, errReason
		}

		replaceableCertType, errReplaceable := state.Queries.GetCertTypeFromId(ctx, succession.ID)
		if errReplaceable != nil {
			return CertTypeSuccessions{}, errReplaceable
		}

		apiSuccessions.CanReplace = append(apiSuccessions.CanReplace, Succession{
			CertType: cert_types.MapCertificateTypeDbToDomain(replaceableCertType),
			Reason:   reason,
		})
	}

	return apiSuccessions, nil

}

func GetAllThatCanReplaceThisId(state *internal.ApiState, ctx context.Context, certType cert_types.CertificateType) (CertTypeSuccessions, error) {

	dbSuccessions, errSuccessions := state.Queries.GetAllThatCanReplaceMe(ctx, certType.Id)
	if errSuccessions != nil {
		return CertTypeSuccessions{}, errSuccessions
	}

	apiSuccessions := CertTypeSuccessions{
		CertType:      certType,
		ReplaceableBy: make([]Succession, len(dbSuccessions)),
	}

	for _, succession := range dbSuccessions {
		reason, errReason := cert_types.SuccessionReasonDbToDomain(succession.ReplaceReason)
		if errReason != nil {
			return CertTypeSuccessions{}, errReason
		}

		replacingCertType, errReplacing := state.Queries.GetCertTypeFromId(ctx, succession.ID)
		if errReplacing != nil {
			return CertTypeSuccessions{}, errReplacing
		}

		apiSuccessions.ReplaceableBy = append(apiSuccessions.ReplaceableBy, Succession{
			CertType: cert_types.MapCertificateTypeDbToDomain(replacingCertType),
			Reason:   reason,
		})

	}

	return apiSuccessions, nil

}
