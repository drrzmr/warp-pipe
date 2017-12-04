package waitfor

import (
	"fmt"
	"net"
	"time"

	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/retry"
)

// TCPPort block until port is achievable
func TCPPort(timeout time.Duration, host string, port uint16) (err error) {

	var (
		conn    net.Conn
		address = fmt.Sprintf("%s:%d", host, port)
	)

	err, innerErr := retry.Do(timeout, func() (err error) {

		conn, err = net.Dial("tcp", address)
		if err == nil {
			conn.Close()
		}

		return err
	})

	if err != nil {
		return errors.Wrapf(innerErr, "error: %s, address: %s, maxAttempts: %d, timeout: %s",
			err.Error(),
			address,
			retry.Config.MaxAttempts,
			timeout.String(),
		)
	}

	return nil
}
