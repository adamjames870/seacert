package certificates

import (
	"context"

	"github.com/adamjames870/seacert/internal"
	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/google/uuid"
)

func DeleteCertificate(state *internal.ApiState, ctx context.Context, certId string, userId uuid.UUID) error {

	cert, err := GetCertificateFromId(state, ctx, certId, userId)
	if err != nil {
		return err
	}

	if cert.DocumentPath != "" && state.Storage != nil {
		errDelete := state.Storage.DeleteObject(ctx, cert.DocumentPath)
		if errDelete != nil {
			state.Logger.Error("Failed to delete certificate document from R2 on certificate deletion", "path", cert.DocumentPath, "error", errDelete)
		} else {
			state.Logger.Info("Deleted certificate document from R2 on certificate deletion", "path", cert.DocumentPath)
		}
	}

	certUuid, _ := uuid.Parse(certId)

	params := sqlc.DeleteCertParams{
		ID:     certUuid,
		UserID: userId,
	}

	return state.Queries.DeleteCert(ctx, params)

}
