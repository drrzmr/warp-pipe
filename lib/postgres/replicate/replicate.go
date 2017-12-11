package replicate

import (
	"github.com/jackc/pgx"

	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pkg/errors"
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

	err = r.conn.CreateReplicationSlot(r.config.Slot, r.config.Plugin)
	if err == nil {
		return nil
	}

	if postgres.IsError(err, postgres.DuplicateObject) && r.config.IgnoreDuplicateObjectError {
		return nil
	}

	return errors.Wrapf(err,
		"Something wrong with slot creation, slot: %s, plugin: %s",
		r.config.Slot,
		r.config.Plugin,
	)
}

func (r *Replicate) start() (err error) {

	err = r.conn.StartReplication(r.config.Slot, startLsn, timeLine, pluginArgs)
	return errors.Wrapf(err,
		"Something wrong with start replication, slot: %s, startLsn: %d, timeLine: %d, pluginArgs: %s",
		r.config.Slot,
		startLsn,
		timeLine,
		pluginArgs,
	)
}
