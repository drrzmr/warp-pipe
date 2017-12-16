package stream

import (
	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/postgres"
)

const (
	pluginArgs = `("include-xids" '1', "include-timestamp" '1', "skip-empty-xacts" '0', "only-local" '0')`
	startLsn   = uint64(0)
	timeLine   = int64(-1)
)

// Replicate object
type Replicate struct {
	config postgres.Config
	conn   *pgx.ReplicationConn
}

// New create a Replicate object
func New(config postgres.Config) *Replicate {

	return &Replicate{
		config: config,
		conn:   nil,
	}
}

// Config return the address of Database config object
func (r *Replicate) Config() *postgres.Config {
	return &r.config
}

// Start replication
func (r *Replicate) Start() (err error) {

	if r.isStarted() {
		return nil
	}

	if err = r.connect(); err != nil {
		return errors.WithStack(err)
	}

	if err = r.createSlot(); err != nil {
		return errors.WithStack(err)
	}

	if err = r.start(); err != nil {
		return errors.WithStack(err)
	}

	return nil
}

func (r *Replicate) isStarted() (started bool) {

	return r.conn != nil
}

func (r *Replicate) connect() (err error) {

	config := pgx.ConnConfig{
		Host:     r.config.Host,
		Port:     r.config.Port,
		User:     r.config.User,
		Password: r.config.Password,
		Database: r.config.Database,
	}

	r.conn, err = pgx.ReplicationConnect(config)

	return errors.Wrapf(err,
		"Could not connect to postgres, host: %s, user: %s, database: %s",
		r.config.Host,
		r.config.User,
		r.config.Database)
}

func (r *Replicate) createSlot() (err error) {

	err = r.conn.CreateReplicationSlot(r.config.Replicate.Slot, r.config.Replicate.Plugin)
	if err == nil {
		return nil
	}

	if postgres.IsError(err, postgres.DuplicateObject) && r.config.Replicate.IgnoreDuplicateObjectError {
		return nil
	}

	return errors.Wrapf(err,
		"Something wrong with slot creation, slot: %s, plugin: %s",
		r.config.Replicate.Slot,
		r.config.Replicate.Plugin,
	)
}

func (r *Replicate) start() (err error) {

	err = r.conn.StartReplication(r.config.Replicate.Slot, startLsn, timeLine, pluginArgs)
	return errors.Wrapf(err,
		"Something wrong with start replication, slot: %s, startLsn: %d, timeLine: %d, pluginArgs: %s",
		r.config.Replicate.Slot,
		startLsn,
		timeLine,
		pluginArgs,
	)
}
