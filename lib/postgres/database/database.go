package database

import (
	"database/sql"
	"fmt"
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

	if err = d.createDatabaseIfNotExist(); err != nil {
		return errors.WithStack(err)
	}

	err = d.connect()
	return errors.WithStack(err)
}

// DB return database object
func (d *Database) DB() *sql.DB {
	return d.db
}

// Disconnect from database
func (d *Database) Disconnect() (err error) {
	if !d.isConnected() {
		return nil
	}

	if err = d.db.Close(); err == nil {
		d.db = nil
		return nil
	}

	dsn, _ := d.config.DSN(true, false)
	return errors.Wrapf(err, "Could not disconnect from database, dsn: %s", dsn)
}

func (d *Database) isConnected() (connected bool) {

	return d.db != nil
}

func (d *Database) connect() (err error) {

	var (
		driver  = d.config.Driver
		dsn, m  = d.config.DSN(true, true)
		timeout = d.config.ConnectTimeout
	)

	if len(m) != 0 {
		return errors.Wrapf(postgres.ErrInvalidDSN, "Missing list should be empty, and not: %v", m)
	}

	d.db, err = connect(driver, dsn, timeout)

	return errors.WithStack(err)
}

func (d *Database) createDatabaseIfNotExist() (err error) {

	if !d.config.CreateDatabaseIfNotExist {
		return nil
	}

	var (
		db       *sql.DB
		driver   = d.config.Driver
		timeout  = d.config.ConnectTimeout
		database = d.config.Database
		dsn, _   = d.config.DSN(false, true)
		query    = fmt.Sprintf(`CREATE DATABASE "%s"`, database)
	)

	if db, err = connect(driver, dsn, timeout); err != nil {
		return errors.WithStack(err)
	}
	defer db.Close()

	_, err = db.Exec(query)
	return errors.Wrapf(err, "Could not create, query: %s, database: %s", query, database)
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
