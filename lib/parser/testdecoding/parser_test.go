package testdecoding_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pagarme/warp-pipe/lib/tester/file"
)

func TestParser(t *testing.T) {

	t.Run("NewParser", func(t *testing.T) {
		parser := testdecoding.NewParser(nil)
		require.NotNil(t, parser)
	})

	t.Run("Parse", func(t *testing.T) {

		operationTable := make([]struct {
			Table string
			Type  string
			Value string
		}, 0)

		file.LoadJSON(t, "operations", &operationTable)

		parser := testdecoding.NewParser(func(transaction testdecoding.Transaction) {

			require.Equal(t, uint64(627), transaction.ID)
			for i, op := range operationTable {
				require.Equal(t, op.Table, transaction.Operations[i].Table)
				require.Equal(t, op.Type, transaction.Operations[i].Type)
				require.Equal(t, op.Value, transaction.Operations[i].Value)
			}

		})

		file.ForEachLine(t, "transaction", func(line string) {
			err := parser.Parse(line)
			require.NoError(t, err)
		})

		file.ForEachLine(t, "transaction", func(line string) {
			err := parser.Parse(line)
			require.NoError(t, err)
		})

		file.ForEachLine(t, "transaction", func(line string) {
			err := parser.Parse(line)
			require.NoError(t, err)
		})
	})
}
