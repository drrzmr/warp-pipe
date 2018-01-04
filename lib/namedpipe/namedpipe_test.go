package namedpipe_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/namedpipe"
)

func TestNamedPipe(t *testing.T) {

	var (
		pipe     = namedpipe.New("testpipe")
		lineList []string
		line     string
		err      error

		ctx, cancel   = context.WithCancel(context.Background())
		lineCh, errCh = pipe.Loop(ctx)
	)

	time.AfterFunc(1*time.Second, func() {
		cancel()
		file, err := os.OpenFile("testpipe", os.O_WRONLY, os.ModeNamedPipe)
		require.NoError(t, err)
		defer file.Close()

		_, err = file.WriteString("testdata\n")
		require.NoError(t, err)
	})

	for ok := true; ok; {
		select {
		case line, ok = <-lineCh:
			if ok {
				lineList = append(lineList, line)
			}
		case err, ok = <-errCh:
			require.NoError(t, err)
		}
	}

	require.Equal(t, []string{"testdata"}, lineList)
	for _, line := range lineList {
		fmt.Println("from pipe:", line)
	}
}
