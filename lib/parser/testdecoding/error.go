package testdecoding

import "errors"

var (
	// ErrInvalidLogPosition returned for invalid log position access
	ErrInvalidLogPosition = errors.New("invalid log position")
	// ErrInvalidMessage returned for invalid message
	ErrInvalidMessage = errors.New("invalid message")
	// ErrInconsistentTransaction returned for inconsistent transaction
	ErrInconsistentTransaction = errors.New("inconsistent transaction")
	// ErrInvalidFilteredMessage returned for invalid message
	ErrInvalidFilteredMessage = errors.New("invalid filtered message")
)
