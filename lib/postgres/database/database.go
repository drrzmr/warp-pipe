package database

import (
	"database/sql"
	"time"

	// connect pgx to database/sql
	_ "github.com/jackc/pgx/stdlib"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/retry"
)

// Database object
type Database struct {
	config postgres.Config
	db     *sql.DB
}

// New returns a new Database object with given configuration
func New(config postgres.Config) *Database {
	return &Database{
		config: config,
	}
}

// Config return the address of Database config object
func (d *Database) Config() *postgres.Config {
	return &d.config
}

// Connect to database using configuration given on New()
func (d *Database) Connect() (err error) {

	if d.isConnected() {
		return nil
	}

	err = d.connect()
	return errors.WithStack(err)
}

func (d *Database) isConnected() (connected bool) {

	return d.db != nil
}

func (d *Database) connect() (err error) {

	var (
		driver  = d.config.Driver
		dsn     = d.config.DSN()
		timeout = d.config.ConnectTimeout
	)

	d.db, err = connect(driver, dsn, timeout)

	return errors.WithStack(err)
}

func connect(driver, dsn string, timeout time.Duration) (db *sql.DB, err error) {

	if db, err = sql.Open(driver, dsn); err != nil {
		return nil, errors.Wrapf(err,
			"Could not open, driver: %s, dsn: %s",
			driver,
			dsn,
		)
	}

	err, innerErr := retry.Do(timeout, func() error {
		return db.Ping()
	})

	if err == nil {
		return db, nil /* Success \o/ */
	}

	return nil, errors.Wrapf(innerErr,
		"Could not connect, retry error: %s, driver: %s, dsn: %s, maxAttempts: %d, timeout: %s",
		err.Error(),
		driver,
		dsn,
		retry.Config.MaxAttempts,
		timeout.String(),
	)
}
