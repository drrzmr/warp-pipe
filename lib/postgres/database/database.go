package database

import "github.com/pagarme/warp-pipe/lib/postgres"

// Database object
type Database struct {
	config postgres.Config
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
