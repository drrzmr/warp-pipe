package postgres

import "time"

// Config object
type Config struct {
	Driver         string
	ConnectTimeout time.Duration
}

// DSN return dsn (data source name)
func (c *Config) DSN() (dsn string) {
	return ""
}
