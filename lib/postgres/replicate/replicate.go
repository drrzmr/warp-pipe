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
