package replicate

import (
	"github.com/jackc/pgx"

	"github.com/pagarme/warp-pipe/lib/postgres"
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

	return nil
}

func (r *Replicate) isStarted() (started bool) {

	return r.conn != nil
}
