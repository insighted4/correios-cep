package postgres

import (
	"context"
	"time"

	"github.com/insighted4/correios-cep/pkg/errors"
	"github.com/insighted4/correios-cep/pkg/log"
	"github.com/insighted4/correios-cep/storage"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/sirupsen/logrus"
)

type Postgres struct {
	db     *pgxpool.Pool
	logger logrus.FieldLogger

	now func() time.Time
}

var _ storage.Storage = (*Postgres)(nil)

// Connect parses a database URL into options that can be used to connect to PostgreSQL.
func Connect(ctx context.Context, cfg *pgxpool.Config, now func() time.Time) (*Postgres, error) {
	const op errors.Op = "postgres.Connect"
	if now == nil {
		now = time.Now
	}

	if cfg == nil || cfg.ConnConfig == nil {
		return nil, errors.E(op, errors.KindUnexpected, "invalid database config")
	}

	logger := log.WithField("component", "postgres")
	logger.Infof("Connecting to postgresql://%s:%d/%s", cfg.ConnConfig.Host, cfg.ConnConfig.Port, cfg.ConnConfig.Database)

	db, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		return nil, errors.E(op, errors.KindUnexpected, err)
	}

	return &Postgres{db: db, logger: logger, now: now}, nil
}

func (p *Postgres) Close() error {
	p.logger.Info("Closing database")
	p.db.Close()
	return nil
}

func (p *Postgres) Check(ctx context.Context) error {
	const op errors.Op = "postgres.Check"

	if err := p.db.Ping(ctx); err != nil {
		return errors.E(op, errors.KindUnexpected, err)
	}

	return nil
}

func (p *Postgres) ExecTx(ctx context.Context, fn func(tx pgx.Tx) error, op errors.Op) error {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return errors.E(op, kind(err), err)
	}
	defer tx.Rollback(ctx)

	if _, err := tx.Exec(ctx, `SET TRANSACTION ISOLATION LEVEL SERIALIZABLE;`); err != nil {
		return errors.E(op, kind(err), err)
	}
	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func kind(err error) int {
	if err.Error() == pgx.ErrNoRows.Error() {
		return errors.KindNotFound
	}

	// TODO (danielnegri): Clean-up native Postgres error handling.
	//pqerr := &pq.Error{}
	//if errors.AsErr(err, &pqerr) {
	//	switch pqerr.Code {
	//	case pgErrUniqueViolation:
	//		return errors.KindAlreadyExists
	//	default:
	//		return errors.KindUnexpected
	//	}
	//}

	return errors.KindUnexpected
}
