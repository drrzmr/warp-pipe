package testdecoding_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pagarme/warp-pipe/lib/tester/file"
)

func TestParser(t *testing.T) {

	/*t.Run("NewParser", func(t *testing.T) {
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
	})*/

	t.Run("Transaction", func(t *testing.T) {

		table := []struct {
			filename    string
			expectError bool
			expectNil   bool
			expectLen   bool
		}{
			{filename: "transaction-without-begin", expectError: true, expectNil: true, expectLen: false},
			{filename: "transaction-without-commit", expectError: false, expectNil: true, expectLen: false},
			{filename: "transaction-invalid", expectError: true, expectNil: true, expectLen: false},
			{filename: "transaction-without-begin-id", expectError: true, expectNil: true, expectLen: false},
			{filename: "transaction-without-commit-id", expectError: true, expectNil: true, expectLen: false},
			{filename: "transaction-without-operations", expectError: false, expectNil: false, expectLen: false},
			{filename: "transaction-multiples", expectError: false, expectNil: false, expectLen: true},
			{filename: "transaction-messy", expectError: true, expectNil: true, expectLen: false},
			{filename: "transaction-messy-but-valid", expectError: true, expectNil: false, expectLen: true},
		}

		for _, tt := range table {
			t.Run(tt.filename, func(t *testing.T) {
				parser := testdecoding.NewParser(func(transaction testdecoding.Transaction) {
					if tt.expectLen {
						require.True(t, len(transaction.Operations) > 0)
					} else {
						require.Len(t, transaction.Operations, 0)
					}

					if tt.expectNil {
						require.Nil(t, transaction)
					} else {
						require.NotNil(t, transaction)
					}
				})

				var errList []error
				file.ForEachLine(t, tt.filename, func(line string) {
					err := parser.Parse(line)
					if err != nil {
						errList = append(errList, err)
					}
				})

				if tt.expectError {
					require.True(t, len(errList) > 0)
				} else {
					require.True(t, len(errList) == 0)
				}
			})
		}
	})
}
