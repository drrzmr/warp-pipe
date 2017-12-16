package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/retry"
)

// Replicate object
type Replicate struct {
	db     *sqlx.DB
	config postgres.Config
}

// New return a sql replicate pointer
func New(config postgres.Config) (r *Replicate) {
	return &Replicate{
		config: config,
		db:     nil,
	}
}

func (r *Replicate) connect() (err error) {

	var (
		dsn, m     = r.config.DSN(true, true)
		safeDsn, _ = r.config.DSN(true, false)
		driver     = r.config.SQL.Driver
		timeout    = r.config.SQL.ConnectTimeout
	)

	if len(m) != 0 {
		return errors.Wrapf(postgres.ErrInvalidDSN, "Missing list should be empty, and not: %v", m)
	}

	err, innerErr := retry.Do(timeout, func() error {
		r.db, err = sqlx.Connect(driver, dsn)
		return err
	})

	if err == nil {
		return nil
	}

	return errors.Wrapf(innerErr,
		"Could not connect, retry error: %s, driver: %s, dsn: %s, maxAttempts: %d, timeout: %s",
		err.Error(),
		driver,
		safeDsn,
		retry.Config.MaxAttempts,
		timeout.String(),
	)
}

func (r *Replicate) isConnected() (connected bool) {
	return r.db != nil
}
