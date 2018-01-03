package reader_test

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"testing"

	cmd "github.com/pagarme/warp-pipe/cmd/reader"
	"github.com/pagarme/warp-pipe/lib/snippet/reader"

	"github.com/stretchr/testify/require"
)

func Test_CmdReader(t *testing.T) {

	var err error

	// sender
	sr, sw := io.Pipe()
	require.NotNil(t, sr)
	require.NotNil(t, sw)

	// receiver
	rr, rw := io.Pipe()
	require.NotNil(t, rr)
	require.NotNil(t, rw)

	text := "testing...\n"

	confReader := &reader.Config{
		InputStream:  sr,
		OutputStream: rw,
	}

	buf := &bytes.Buffer{}
	readerCmd := cmd.New(confReader)
	readerCmd.SetOutput(buf)
	out := readerCmd.OutOrStdout()
	require.Equal(t, buf, out)

	// buf -> sw -> sr
	go func() {
		var errSender error
		buf := []byte(text)

		n, errSender := sw.Write(buf)
		require.Equal(t, len(buf), n)
		require.NoError(t, errSender)

		errSender = sw.Close()
		require.NoError(t, errSender)
	}()

	// sr -> rw -> rr
	go func() {
		errCmd := readerCmd.Execute()
		require.NoError(t, errCmd)
	}()

	// rr -> receivedBuf
	receivedBuf, err := ioutil.ReadAll(rr)
	receivedText := string(receivedBuf)
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("-> %s", text), receivedText)
}
