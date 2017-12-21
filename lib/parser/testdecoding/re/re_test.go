package re_test

import (
	"strconv"
	"testing"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding/re"
	"github.com/pagarme/warp-pipe/lib/tester/file"
	"github.com/stretchr/testify/require"
)

func TestMessageRegexp(t *testing.T) {

	testTable := make([]struct {
		Begin     bool
		Commit    bool
		Operation bool
		Insert    bool
		Update    bool
		Delete    bool
		Message   string
		ID        uint64

		Schema string
		Table  string
		Values string
	}, 0)

	file.LoadJSON(t, "operations", &testTable)

	t.Run("Begin/Commit", func(t *testing.T) {
		for _, tt := range testTable {
			var matchList []string

			if tt.Begin {
				matchList = re.Begin.FindStringSubmatch(tt.Message)
			}

			if tt.Commit {
				matchList = re.Commit.FindStringSubmatch(tt.Message)
			}

			if !tt.Begin && !tt.Commit {
				require.Nil(t, matchList)
				continue
			}
			require.NotNil(t, matchList)
			require.Len(t, matchList, 2)
			require.Equal(t, matchList[0], tt.Message)
			id, err := strconv.ParseUint(matchList[1], 10, 64)
			require.NoError(t, err)
			require.Equal(t, tt.ID, id)
		}
	})
}

func TestOperationRegexp(t *testing.T) {

	testTable := make([]struct {
		Begin     bool
		Commit    bool
		Operation bool
		Insert    bool
		Update    bool
		Delete    bool
		Message   string
		ID        uint64

		Schema string
		Table  string
		Values string
	}, 0)

	file.LoadJSON(t, "operations", &testTable)

	t.Run("Operation", func(t *testing.T) {
		for _, tt := range testTable {
			matchList := re.Operation.FindStringSubmatch(tt.Message)
			if !tt.Operation {
				require.Nil(t, matchList)
				continue
			}
			require.NotNil(t, matchList)
			require.Len(t, matchList, 5)
			require.Equal(t, tt.Message, matchList[0])
			require.Equal(t, tt.Schema, matchList[1])
			require.Equal(t, tt.Table, matchList[2])
			if tt.Insert {
				require.Equal(t, "INSERT", matchList[3])
			} else if tt.Update {
				require.Equal(t, "UPDATE", matchList[3])
			} else if tt.Delete {
				require.Equal(t, "DELETE", matchList[3])
			} else {
				require.FailNow(t, "invalid operation type", matchList[3])
			}
			require.Equal(t, tt.Values, matchList[4])
		}
	})
}

func TestRowRegexp(t *testing.T) {

	testTable := make([]struct {
		Begin     bool
		Commit    bool
		Operation bool
		Insert    bool
		Update    bool
		Delete    bool
		Message   string
		ID        uint64

		Schema string
		Table  string
		Values string
	}, 0)

	file.LoadJSON(t, "operations", &testTable)

	t.Run("Row", func(t *testing.T) {
		tt := testTable[1]
		matchList := re.Operation.FindStringSubmatch(tt.Message)
		require.NotNil(t, matchList)
		require.Len(t, matchList, 5)
		require.Equal(t, tt.Message, matchList[0])
		require.Equal(t, tt.Schema, matchList[1])
		require.Equal(t, tt.Table, matchList[2])
		require.Equal(t, tt.Values, matchList[4])

		rows := re.Row.FindAllStringSubmatch(matchList[4], -1)

		for _, row := range rows {
			require.Len(t, row, 4)
		}
	})
}
