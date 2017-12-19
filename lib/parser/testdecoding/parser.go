package testdecoding

import (
	"github.com/looplab/fsm"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/parser/testdecoding/re"
)

// Parser parser struct
type Parser struct {
	fsm *fsm.FSM
	log []string
}

// Transaction the transactio itself
type Transaction struct {
	Operations []Operation
	ID         uint64
}

// Operation the operation unit
type Operation struct {
	Schema string
	Table  string
	Type   string
	Value  string
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

// Parse parses log lines
func (p *Parser) Parse() (t Transaction, err error) {
	t = Transaction{}

	for _, line := range p.log {
		switch p.fsm.Current() {
		case "idle":
			if re.Begin.MatchString(line) {
				err := p.fsm.Event("begin")
				if err != nil {
					return t, errors.Wrap(err, "(begin) transition error")
				}
			}
		case "parsing":
			if re.Commit.MatchString(line) {
				err := p.fsm.Event("commit")
				if err != nil {
					return t, errors.Wrap(err, "(commit) transition error")
				}
			}

			if re.Operation.MatchString(line) {
				matchList := re.Operation.FindStringSubmatch(line)
				t.Operations = append(t.Operations, Operation{
					Schema: matchList[0],
					Table:  matchList[1],
					Type:   matchList[2],
					Value:  matchList[3],
				})
			}
		}
	}

	return t, nil
}
