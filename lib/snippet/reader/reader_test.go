package reader

import (
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRun(t *testing.T) {

	// sender
	sr, sw := io.Pipe()
	require.NotNil(t, sr)
	require.NotNil(t, sw)

	// receiver
	rr, rw := io.Pipe()
	require.NotNil(t, rr)
	require.NotNil(t, rw)

	text := "testing...\n"

	// buf -> sw -> sr
	go func() {
		buf := []byte(text)
		n, err := sw.Write(buf)
		require.Equal(t, len(buf), n)
		require.NoError(t, err)

		err = sw.Close()
		require.NoError(t, err)
	}()

	// sr -> rw -> rr
	go func() {
		err := Run(sr, rw)
		require.NoError(t, err)
	}()

	// rr -> receivedBuf
	receivedBuf, err := ioutil.ReadAll(rr)
	receivedText := string(receivedBuf)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("-> %s", text), receivedText)
}
