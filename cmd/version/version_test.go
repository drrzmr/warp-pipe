package version

import (
	"bytes"
	"fmt"
	"testing"

	"github.com/pagarme/warp-pipe/config"
	"github.com/stretchr/testify/require"
)

func Test_CmdVersion(t *testing.T) {

	buf := &bytes.Buffer{}
	versionCmd := New()
	versionCmd.SetOutput(buf)
	out := versionCmd.OutOrStdout()
	require.Equal(t, buf, out)
	err := versionCmd.Execute()
	require.NoError(t, err)
	require.Equal(t, fmt.Sprintf("version: %s\n", config.AppVersion), buf.String())
}
