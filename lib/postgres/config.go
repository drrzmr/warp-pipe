package postgres

import "time"

const (
	// MinConnectTimeout min timeout necessary for connect
	MinConnectTimeout = 1 * time.Second
)

// Config object
type Config struct {
	Driver         string
	ConnectTimeout time.Duration
}

// DSN return dsn (data source name)
func (c *Config) DSN() (dsn string) {
	return ""
}
