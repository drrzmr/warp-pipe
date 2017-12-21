package testdecoding

import (
	"strconv"

	"github.com/looplab/fsm"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/parser/testdecoding/event"
	"github.com/pagarme/warp-pipe/lib/parser/testdecoding/re"
	"github.com/pagarme/warp-pipe/lib/parser/testdecoding/state"
)

// TransactionFunc callback called for new transaction
type TransactionFunc func(transaction Transaction)

// Parser parser struct
type Parser struct {
	transaction  *Transaction
	stateMachine *fsm.FSM
	callback     TransactionFunc
}

// Transaction the transaction itself
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

var logger = log.Development("parser")

// NewParser return a new parser struct
func NewParser(transactionFunc TransactionFunc) *Parser {

	return &Parser{
		transaction: nil,
		callback:    transactionFunc,
		stateMachine: fsm.NewFSM(
			state.Idle,
			fsm.Events{
				{
					Name: event.BeginIn,
					Dst:  state.Start,
					Src:  []string{state.Idle},
				},
				{
					Name: event.BeginOut,
					Dst:  state.Started,
					Src:  []string{state.Start},
				},
				{
					Name: event.OperationIn,
					Dst:  state.Store,
					Src:  []string{state.Started, state.Stored},
				},
				{
					Name: event.OperationOut,
					Dst:  state.Stored,
					Src:  []string{state.Store},
				},
				{
					Name: event.CommitIn,
					Dst:  state.Publish,
					Src:  []string{state.Stored, state.Started},
				},
				{
					Name: event.CommitOut,
					Dst:  state.Idle,
					Src:  []string{state.Publish},
				},
			},
			fsm.Callbacks{
				"enter_state": func(e *fsm.Event) {
					logger.Debug("enter state",
						zap.String("event", e.Event),
						zap.String("src", e.Src),
						zap.String("dst", e.Dst),
					)
				},
			},
		),
	}
}

// Parse parses log lines
func (p *Parser) Parse(msg string) (err error) {

	if list := re.Operation.FindStringSubmatch(msg); list != nil {
		err = p.handleOperation(list[1:])
	} else if list := re.Begin.FindStringSubmatch(msg); list != nil {
		err = p.handleBegin(list[1:])
	} else if list := re.Commit.FindStringSubmatch(msg); list != nil {
		err = p.handleCommit(list[1:])
	} else {
		err = errors.Wrapf(ErrInvalidMessage, "message: %s", msg)
	}

	return errors.WithStack(err)
}

func (p *Parser) handleBegin(list []string) (err error) {
	if len(list) != 1 {
		return errors.Wrapf(ErrInvalidFilteredMessage, "filteredMessage: %v", list)
	}

	if err = p.stateMachine.Event(event.BeginIn); err != nil {
		return errors.Wrapf(err,
			"[handleBegin] state: %s, list: %v",
			p.stateMachine.Current(),
			list)
	}

	id, err := strconv.ParseUint(list[0], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "could not parse begin id: %s", list[0])
	}

	p.transaction = &Transaction{
		ID:         id,
		Operations: []Operation{},
	}

	if err = p.stateMachine.Event(event.BeginOut); err != nil {
		return errors.Wrapf(err,
			"[handleBegin] state: %s, list: %v",
			p.stateMachine.Current(),
			list)
	}

	return nil
}

func (p *Parser) handleCommit(list []string) (err error) {
	if len(list) != 1 {
		return errors.Wrapf(ErrInvalidFilteredMessage, "filteredMessage: %v", list)
	}

	if err = p.stateMachine.Event(event.CommitIn); err != nil {
		return errors.Wrapf(err,
			"[handleBegin] state: %s, list: %v",
			p.stateMachine.Current(),
			list)
	}

	id, err := strconv.ParseUint(list[0], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "could not parse commit id: %s", list[0])
	}

	if id != p.transaction.ID {
		return errors.Wrapf(ErrInconsistentTransaction,
			"begin id: %d, commit id: %d",
			p.transaction.ID,
			id,
		)
	}

	p.callback(*p.transaction)
	p.transaction = nil

	if err = p.stateMachine.Event(event.CommitOut); err != nil {
		return errors.Wrapf(err,
			"[handleBegin] state: %s, list: %v",
			p.stateMachine.Current(),
			list)
	}

	return nil
}

func (p *Parser) handleOperation(list []string) (err error) {
	if len(list) != 4 {
		return errors.Wrapf(ErrInvalidFilteredMessage, "filteredMessage: %v", list)
	}

	if err = p.stateMachine.Event(event.OperationIn); err != nil {
		return errors.Wrapf(err,
			"[handleBegin] state: %s, list: %v",
			p.stateMachine.Current(),
			list)
	}

	p.transaction.Operations = append(p.transaction.Operations, Operation{
		Schema: list[0],
		Table:  list[1],
		Type:   list[2],
		Value:  list[3],
	})

	if err = p.stateMachine.Event(event.OperationOut); err != nil {
		return errors.Wrapf(err,
			"[handleBegin] state: %s, list: %v",
			p.stateMachine.Current(),
			list)
	}

	return nil
}
