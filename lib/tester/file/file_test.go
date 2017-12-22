package file_test

import (
	"testing"

	"github.com/pagarme/warp-pipe/lib/tester/file"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {

	text := file.Load(t, "loremipson")
	file.MustMatch(t, "loremipson", text)
}

func TestLoadJson(t *testing.T) {
	sampleTable := []struct {
		LoadedName string `json:"name"`
		LoadedAge  int    `json:"age"`

		expectedName string
		expectedAge  int
	}{
		{expectedName: "jon doe", expectedAge: 33},
	}

	file.LoadJSON(t, "sample", &sampleTable)

	for _, tt := range sampleTable {
		require.Equal(t, tt.expectedName, tt.LoadedName)
		require.Equal(t, tt.expectedAge, tt.LoadedAge)
	}
}
