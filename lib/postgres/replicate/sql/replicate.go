package sql

import (
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/retry"
)

/*
 * https://www.postgresql.org/docs/9.4/static/logicaldecoding-example.html
 * https://www.postgresql.org/docs/9.4/static/functions-admin.html
 */

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

// Connect public method
func (r *Replicate) Connect() (err error) {

	if r.isConnected() {
		return nil
	}

	if err = r.connect(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

// CreateSlot public method
func (r *Replicate) CreateSlot() (created bool, err error) {

	created, err = createSlot(r.db, r.config.Replicate.Slot, r.config.Replicate.Plugin)
	return created, errors.WithStack(err)
}

// ListSlots public method
func (r *Replicate) ListSlots() (result []SlotInfo, err error) {

	result, err = listSlots(r.db)
	return result, errors.WithStack(err)
}

// GetAllChanges public method
func (r *Replicate) GetAllChanges() (result []ReplicationEvent, err error) {

	result, err = getAllChanges(r.db, r.config.Replicate.Slot)
	return result, errors.WithStack(err)
}

// PeekAllChanges public method
func (r *Replicate) PeekAllChanges() (result []ReplicationEvent, err error) {

	result, err = peekAllChanges(r.db, r.config.Replicate.Slot)
	return result, errors.WithStack(err)
}

// DropSlot public method
func (r *Replicate) DropSlot() (err error) {

	err = dropSlot(r.db, r.config.Replicate.Slot)
	return errors.WithStack(err)
}

// Close public method
func (r *Replicate) Close() (err error) {

	err = r.db.Close()
	if err == nil {
		r.db = nil
	}
	return errors.Wrap(err, "Could not close connection")
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
