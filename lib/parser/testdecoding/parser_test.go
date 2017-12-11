package testdecoding_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
	"github.com/pagarme/warp-pipe/lib/tester/file"
)

func TestParser(t *testing.T) {

	t.Run("NewParser", func(t *testing.T) {
		parser := testdecoding.NewParser()
		require.NotNil(t, parser)
		fsm := parser.FSM()
		require.NotNil(t, fsm)
	})

	t.Run("FSMTransitions", func(t *testing.T) {
		parser := testdecoding.NewParser()
		fsm := parser.FSM()

		require.Equal(t, fsm.Current(), "idle")
		fsm.Event("begin")
		require.Equal(t, fsm.Current(), "parsing")
		fsm.Event("parse")
		require.Equal(t, fsm.Current(), "parsing")
		fsm.Event("parse")
		require.Equal(t, fsm.Current(), "parsing")
		fsm.Event("commit")
		require.Equal(t, fsm.Current(), "idle")
	})

	t.Run("Append", func(t *testing.T) {
		parser := testdecoding.NewParser()

		testTable := []struct {
			line string
		}{
			{"BEGIN 627"},
			{"COMMIT 627"},
		}
		var (
			readLine string
			err      error
		)

		for i, testCase := range testTable {
			parser.Append(testCase.line)
			readLine, err = parser.Log(uint(i))
			require.NoError(t, err)
			require.Equal(t, readLine, testCase.line)
		}
	})

	t.Run("Parse", func(t *testing.T) {
		parser := testdecoding.NewParser()

		file.ForEachLine(t, "transaction", func(line string) {
			parser.Append(line)
		})

		transaction, err := parser.Parse()

		require.NoError(t, err)
		require.NotNil(t, transaction)
		require.NotEmpty(t, transaction)
	})
}
