package retry

import (
	"errors"
	"time"
)

const (
	// DefaultSleep default time to sleep between retries
	DefaultSleep = 100 * time.Millisecond
	// DefaultAttempts default number of attempts
	DefaultAttempts = uint64(100)
)

// CallbackFunc callback func for retry.Do
type CallbackFunc func() (err error)

type config struct {
	Sleep       time.Duration
	MaxAttempts uint64
}

var (
	// Config store retry configuration
	Config = config{
		Sleep:       DefaultSleep,
		MaxAttempts: DefaultAttempts,
	}
	// ErrTimeout returned when retry.Do finish by timeout
	ErrTimeout = errors.New("retry timeout")
	// ErrMaxAttempts returned when retry.Do finish by max attempts
	ErrMaxAttempts = errors.New("max attempts")
	// ErrContinue can be used by CallbackFunc to make retry.Do not finish
	ErrContinue = errors.New("continue")
)

// Do exec retry calling given CallbackFunc
func Do(timeout time.Duration, fn CallbackFunc) (err, innerErr error) {

	attempts := uint64(0)
	innerErr = ErrTimeout

	doneAt := time.Now().Add(timeout)
	for time.Now().Before(doneAt) {

		innerErr = fn()
		attempts++

		if innerErr == nil {
			return nil, nil
		}

		if attempts > Config.MaxAttempts {
			return ErrMaxAttempts, innerErr
		}

		time.Sleep(Config.Sleep)
	}

	return ErrTimeout, innerErr
}
