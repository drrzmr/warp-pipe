package testdecoding

import (
	"github.com/looplab/fsm"
	"github.com/pkg/errors"
)

// Parser parser struct
type Parser struct {
	fsm *fsm.FSM
	log []string
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

// Append append log lines
func (p *Parser) Append(str string) {
	p.log = append(p.log, str)
}

// Log returns the log lines to be processed
func (p *Parser) Log(position uint) (line string, err error) {

	maxPos := uint(len(p.log) - 1)
	if position > maxPos {
		return "", errors.Wrapf(ErrInvalidLogPosition, "position: %d, max position: %d", position, maxPos)
	}

	return p.log[position], nil
}
