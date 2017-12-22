package testdecoding_test

import (
	"strconv"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pagarme/warp-pipe/lib/tester/file"
)

func TestParser(t *testing.T) {

	var testTable []struct {
		Filename     string `json:"filename"`
		Transactions []struct {
			ID         uint64 `json:"id"`
			Operations []struct {
				Schema string `json:"schema"`
				Table  string `json:"table"`
				Type   string `json:"type"`
				Value  string `json:"value"`
			} `json:"operations"`
		} `json:"transactions"`
		ExpectedErrors []string `json:"expectedErrors"`
	}

	file.LoadJSON(t, "test-table", &testTable)

	// foreach file
	for _, tt := range testTable {
		t.Run(tt.Filename, func(t *testing.T) {

			var generatedTransactions []testdecoding.Transaction
			parser := testdecoding.NewParser(func(transaction testdecoding.Transaction) {
				generatedTransactions = append(generatedTransactions, transaction)
			})
			require.NotNil(t, parser)

			var errorList []error
			file.ForEachLine(t, tt.Filename, func(line string) {
				if err := parser.Parse(line); err != nil {
					errorList = append(errorList, err)
				}
			})
			require.Len(t, errorList, len(tt.ExpectedErrors))

			// foreach error
			for i, expectedError := range tt.ExpectedErrors {
				eid := strconv.FormatInt(int64(i), 10)

				t.Run(eid, func(t *testing.T) {
					require.Equal(t, expectedError, errors.Cause(errorList[i]).Error())
				})
			}

			require.Len(t, generatedTransactions, len(tt.Transactions))

			// foreach transaction
			for i, transaction := range tt.Transactions {
				tid := strconv.FormatUint(transaction.ID, 10)

				t.Run(tid, func(t *testing.T) {
					require.Equal(t, transaction.ID, generatedTransactions[i].ID)
					require.Len(t, generatedTransactions[i].Operations, len(transaction.Operations))
					// foreach operation
					for j, operation := range transaction.Operations {
						oid := strconv.FormatInt(int64(j), 10)

						t.Run(oid, func(t *testing.T) {
							require.Equal(t, operation.Schema, generatedTransactions[i].Operations[j].Schema)
							require.Equal(t, operation.Table, generatedTransactions[i].Operations[j].Table)
							require.Equal(t, operation.Type, generatedTransactions[i].Operations[j].Type)
							require.Equal(t, operation.Value, generatedTransactions[i].Operations[j].Value)
						})
					}
				})
			}

		})
	}
}
