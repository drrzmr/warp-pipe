package postgres

import (
	"errors"
	"strconv"

	"github.com/jackc/pgx"
)

// https://www.postgresql.org/docs/9.4/static/errcodes-appendix.html

// ErrorCode postgres error code
type ErrorCode uint64

const (
	// DuplicateObject postgres error code
	DuplicateObject ErrorCode = 42710
)

var (
	// ErrInvalidDSN returned when DSN string are invalid
	ErrInvalidDSN = errors.New("invalid dsn")
)

func (code ErrorCode) str() string {

	return strconv.FormatUint(uint64(code), 10)

}

// IsError return true if err match with given error code
func IsError(err error, code ErrorCode) bool {

	pgErr, ok := err.(pgx.PgError)
	return ok && pgErr.Code == code.str()
}
