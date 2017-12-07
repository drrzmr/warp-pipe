package testdecoding

import (
	"github.com/looplab/fsm"
)

// Parser parser struct
type Parser struct {
	fsm *fsm.FSM
}

// NewParser return a new parser struct
func NewParser() *Parser {

	parserStateMachine := fsm.NewFSM(
		"idle",
		fsm.Events{
			{
				Name: "begin",
				Src:  []string{"idle"},
				Dst:  "parsing",
			}, {
				Name: "parse",
				Src:  []string{"parsing"},
				Dst:  "parsing",
			}, {
				Name: "commit",
				Src:  []string{"parsing"},
				Dst:  "idle",
			},
		},
		fsm.Callbacks{
			"begin":  func(e *fsm.Event) {},
			"parse":  func(e *fsm.Event) {},
			"commit": func(e *fsm.Event) {},
		},
	)

	return &Parser{
		fsm: parserStateMachine,
	}
}

// FSM returns parser FSM
func (p *Parser) FSM() (fsm *fsm.FSM) {
	return p.fsm
}
