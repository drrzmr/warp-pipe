package testdecoding_test

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding"
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
}
