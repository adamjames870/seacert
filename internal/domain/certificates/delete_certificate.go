package certificates

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/google/uuid"
)

func DeleteCertificate(state *internal.ApiState, ctx context.Context, certId string, userId uuid.UUID) error {

	certUuid, errId := uuid.Parse(certId)
	if errId != nil {
		return errId
	}

	params := sqlc.DeleteCertParams{
		ID:     certUuid,
		UserID: userId,
	}

	return state.Queries.DeleteCert(ctx, params)

}
