package ping

import (
	"bufio"
	"bytes"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

type config struct {
	CommandPath string
	Count       uint64
	Prefix      string
	Suffix      string
}

// Config object
var Config = config{
	CommandPath: "/bin/ping",
	Count:       1,
	Prefix:      "--- ",
	Suffix:      " ---",
}

// Run execute ping command using Config
func Run(t *testing.T, host string) (transmitted, received uint64) {
	t.Helper()

	var (
		count = fmt.Sprintf("-c%d", Config.Count)
		ping  = Config.CommandPath
	)

	cmd := exec.Command(ping, host, count)
	require.NotNil(t, cmd)

	buf, err := cmd.Output()
	require.NoError(t, err)
	require.NotNil(t, buf)

	line := findLine(t, buf)
	require.NotEmpty(t, line)

	return parseLine(t, line)
}

func tryMatch(str, prefix, suffix string) (match bool) {

	return strings.HasPrefix(str, prefix) && strings.HasSuffix(str, suffix)
}

func findLine(t *testing.T, buf []byte) (line string) {
	t.Helper()

	reader := bytes.NewReader(buf)
	require.NotNil(t, reader)

	scanner := bufio.NewScanner(reader)
	require.NotNil(t, scanner)

	for scanner.Scan() {
		str := string(scanner.Bytes())
		match := tryMatch(str, Config.Prefix, Config.Suffix)
		if !match {
			continue
		}

		require.True(t, scanner.Scan())
		return string(scanner.Bytes())
	}

	return ""
}

func parseLine(t *testing.T, line string) (transmitted, received uint64) {
	t.Helper()

	words := strings.Split(strings.TrimSpace(line), " ")
	require.Len(t, words, 10)

	var err error

	transmitted, err = strconv.ParseUint(words[0], 10, 64)
	require.NoError(t, err)

	received, err = strconv.ParseUint(words[3], 10, 64)
	require.NoError(t, err)

	return transmitted, received
}
