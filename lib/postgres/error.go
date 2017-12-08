package postgres

import "errors"

var (
	// ErrInvalidDSN returned when DSN string are invalid
	ErrInvalidDSN = errors.New("invalid dsn")
)
