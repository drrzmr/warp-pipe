package testdecoding

import "errors"

var (
	// ErrInvalidLogPosition returned for invalid log position access
	ErrInvalidLogPosition = errors.New("invalid log position")
)
