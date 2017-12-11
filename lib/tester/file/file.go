package file

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

var update = flag.Bool("update", false, "update output")

// Load filename from testdata directory
func Load(t *testing.T, filename string) (text string) {
	t.Helper()

	fn := fmt.Sprintf("%s.%s", filename, Config.InputExtension)
	path := filepath.Join("testdata", fn)
	buf, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.NotNil(t, buf)

	text = string(buf)

	if *update {
		Store(t, filename, text)
	}

	return text
}

// Store text into filename using Config.OutputExtension
func Store(t *testing.T, filename, text string) {
	t.Helper()

	fn := fmt.Sprintf("%s.%s", filename, Config.OutputExtension)
	path := filepath.Join("testdata", fn)
	err := ioutil.WriteFile(path, []byte(text), 0644)
	require.NoError(t, err)
}

// MustMatch verify if text match with filename
func MustMatch(t *testing.T, filename, text string) {
	t.Helper()

	fn := fmt.Sprintf("%s.%s", filename, Config.OutputExtension)
	path := filepath.Join("testdata", fn)
	buf, err := ioutil.ReadFile(path)
	require.NoError(t, err)
	require.NotNil(t, buf)

	golden := string(buf)
	require.Equal(t, golden, text)
}

// ForEachLineFunc callback argument to ForEachLine func
type ForEachLineFunc func(line string)

// ForEachLine call ForEachLineFunc for each line in filename
func ForEachLine(t *testing.T, filename string, cb ForEachLineFunc) {
	t.Helper()

	fn := fmt.Sprintf("%s.%s", filename, Config.InputExtension)
	path := filepath.Join("testdata", fn)
	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		cb(scanner.Text())
	}

	require.NoError(t, scanner.Err())
}
