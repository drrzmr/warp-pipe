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
		operationTable := []struct {
			Table string
			Type  string
			Value string
		}{
			{
				Table: "table_with_pk",
				Type:  "INSERT",
				Value: "a[integer]:1 " +
					"b[character varying]:'Backup and Restore' " +
					"c[timestamp without time zone]:'2017-11-30 17:59:33.825033'",
			}, {
				Table: "table_with_pk",
				Type:  "INSERT",
				Value: "a[integer]:2 " +
					"b[character varying]:'Tuning' " +
					"c[timestamp without time zone]:'2017-11-30 17:59:33.825033'",
			}, {
				Table: "table_with_pk",
				Type:  "INSERT",
				Value: "a[integer]:3 " +
					"b[character varying]:'Replication' " +
					"c[timestamp without time zone]:'2017-11-30 17:59:33.825033'",
			}, {
				Table: "table_with_pk",
				Type:  "DELETE",
				Value: "a[integer]:1 " +
					"c[timestamp without time zone]:'2017-11-30 17:59:33.825033'",
			}, {
				Table: "table_with_pk",
				Type:  "DELETE",
				Value: "a[integer]:2 " +
					"c[timestamp without time zone]:'2017-11-30 17:59:33.825033'",
			}, {
				Table: "table_without_pk",
				Type:  "INSERT",
				Value: "a[integer]:1 " +
					"b[numeric]:2.34 " +
					"c[text]:'Tapir'",
			}, {
				Table: "table_without_pk",
				Type:  "UPDATE",
				Value: "a[integer]:1 " +
					"b[numeric]:2.34 " +
					"c[text]:'Anta'",
			},
		}

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
