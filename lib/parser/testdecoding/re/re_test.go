package re_test

import (
	"regexp"
	"strconv"
	"testing"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding/re"
	"github.com/pagarme/warp-pipe/lib/tester/file"
	"github.com/stretchr/testify/require"
)

func TestRegexp(t *testing.T) {

	controlMessage := map[string]*regexp.Regexp{
		"begin":  re.Begin,
		"commit": re.Commit,
	}

	operationMessage := map[string]*regexp.Regexp{
		"insert": re.Operation,
		"update": re.Operation,
		"delete": re.Operation,
	}

	var testTable []struct {
		Type    string `json:"type"`
		Message string `json:"message"`
		ID      uint64 `json:"id"`
		Schema  string `json:"schema"`
		Table   string `json:"table"`
		Columns []struct {
			Name  string `json:"name"`
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"columns"`
	}

	file.LoadJSON(t, "operations", &testTable)

	for _, tt := range testTable {
		if regexpControl := controlMessage[tt.Type]; regexpControl != nil {

			t.Run("Begin/Commit", func(t *testing.T) {
				matchList := regexpControl.FindStringSubmatch(tt.Message)
				require.NotNil(t, matchList)
				require.Len(t, matchList, 2)
				require.Equal(t, matchList[0], tt.Message)
				id, err := strconv.ParseUint(matchList[1], 10, 64)
				require.NoError(t, err)
				require.Equal(t, tt.ID, id)
			})
			continue
		}

		if regexpMessage := operationMessage[tt.Type]; regexpMessage != nil {

			t.Run("Operation", func(t *testing.T) {
				matchList := re.Operation.FindStringSubmatch(tt.Message)
				require.NotNil(t, matchList)
				require.Len(t, matchList, 5)
				require.Equal(t, tt.Message, matchList[0])
				require.Equal(t, tt.Schema, matchList[1])
				require.Equal(t, tt.Table, matchList[2])
				switch tt.Type {
				case "insert":
					require.Equal(t, "INSERT", matchList[3])
				case "update":
					require.Equal(t, "UPDATE", matchList[3])
				case "delete":
					require.Equal(t, "DELETE", matchList[3])
				default:
					require.FailNow(t, "invalid operation type", matchList[3])
				}

				t.Run("Columns", func(t *testing.T) {
					columnsMatch := re.Row.FindAllStringSubmatch(matchList[4], -1)
					require.NotNil(t, columnsMatch)
					require.Len(t, columnsMatch, len(tt.Columns))
					for i, col := range tt.Columns {

						t.Run(col.Name, func(t *testing.T) {
							require.Len(t, columnsMatch[i], 4)
							require.Equal(t, col.Name, columnsMatch[i][1])
							require.Equal(t, col.Type, columnsMatch[i][2])
							require.Equal(t, col.Value, columnsMatch[i][3])
						})
					}
				})
			})
			continue
		}

		require.FailNow(t, "invalid test case", tt)
	}
}
