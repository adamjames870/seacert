package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/adamjames870/seacert/internal/database/sqlc"
	"github.com/adamjames870/seacert/internal/domain"
)

type repository struct {
	db *sql.DB
	*sqlc.Queries
}

func NewRepository(db *sql.DB) domain.Repository {
	return &repository{
		db:      db,
		Queries: sqlc.New(db),
	}
}

func (r *repository) ResetAll(ctx context.Context) error {
	return r.WithTx(ctx, func(txRepo domain.Repository) error {
		_ = r.Queries.ResetSeatimePeriods(ctx)
		_ = r.Queries.ResetSeatime(ctx)
		_ = r.Queries.ResetShips(ctx)
		_ = r.Queries.ResetSuccessions(ctx)
		_ = r.Queries.ResetCerts(ctx)
		_ = r.Queries.ResetCertTypes(ctx)
		_ = r.Queries.ResetIssuers(ctx)
		_ = r.Queries.ResetUsers(ctx)
		return nil
	})
}

func (r *repository) WithTx(ctx context.Context, fn func(domain.Repository) error) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	qtx := r.Queries.WithTx(tx)
	err = fn(&repository{
		db:      r.db,
		Queries: qtx,
	})

	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err
	}

	return tx.Commit()
}
