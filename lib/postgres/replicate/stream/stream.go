package stream

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream/handler"
	"github.com/pagarme/warp-pipe/lib/postgres/replicate/stream/listener"
	"go.uber.org/zap"
)

const (
	pluginArgs = `("include-xids" '1', "include-timestamp" '1', "skip-empty-xacts" '0', "only-local" '0')`
	startLsn   = uint64(0)
	timeLine   = int64(-1)
)

var logger = log.Development("stream")

// Stream object
type Stream struct {
	config   postgres.Config
	conn     *pgx.ReplicationConn
	listener listener.EventListener
}

// New create a Replicate object
func New(config postgres.Config) *Stream {

	return &Stream{
		config:   config,
		conn:     nil,
		listener: nil,
	}
}

// Config return the address of Database config object
func (s *Stream) Config() *postgres.Config {
	return &s.config
}

// NewDefaultEventListener return new default listener
func (s *Stream) NewDefaultEventListener(h handler.EventHandler) listener.EventListener {

	return listener.NewDefaultEventListener(s.conn, s.config.Streaming.WaitMessageTimeout, h)
}

// Connect to postgres
func (s *Stream) Connect() (err error) {

	logger.Debug("--> Connect()")
	defer logger.Debug("<-- Connect()")

	if s.isConnected() {
		logger.Info("already connected")
		return nil
	}

	if err = s.connect(); err != nil {
		logger.Error("connect error")
		return errors.WithStack(err)
	}

	if err = s.createSlot(); err != nil {
		logger.Error("create slot error")
		return errors.WithStack(err)
	}

	return nil
}

// Start replication
func (s *Stream) Start(ctx context.Context, listener listener.EventListener) (started bool, err error) {

	logger.Debug("--> Start()")
	defer logger.Debug("<-- Start()")

	if !s.isConnected() {
		logger.Error("not connected")
		return false, nil
	}

	if s.isStarted() {
		logger.Error("already started")
		return false, nil
	}

	if err = s.start(); err != nil {
		logger.Error("start error")
		return false, errors.WithStack(err)
	}

	s.listener = listener

	if err = s.listener.Listen(ctx); err != nil {

		// filter context canceled
		canceled := errors.Cause(err) == context.Canceled
		logger.DebugIf(canceled, "context canceled")
		logger.ErrorIf(!canceled, "listener run", zap.Error(err))
		return false, errors.WithStack(err)
	}

	return true, nil
}

// SendStandByStatus method
func (s *Stream) SendStandByStatus(position uint64) (err error) {

	var status *pgx.StandbyStatus

	if status, err = pgx.NewStandbyStatus(position); err != nil {
		return errors.Wrapf(err, "create new standby status object failed, position: %d", position)
	}

	err = s.conn.SendStandbyStatus(status)
	logger.DebugIf(err == nil, "sent standby status", zap.Uint64("position", position))

	return errors.Wrapf(err, "send stand by status failed, position: %d", position)
}

func (s *Stream) isConnected() (connected bool) { return s.conn != nil }
func (s *Stream) isStarted() (started bool)     { return s.listener != nil }

func (s *Stream) connect() (err error) {

	config := pgx.ConnConfig{
		Host:     s.config.Host,
		Port:     s.config.Port,
		User:     s.config.User,
		Password: s.config.Password,
		Database: s.config.Database,
	}

	s.conn, err = pgx.ReplicationConnect(config)

	return errors.Wrapf(err,
		"Could not connect to postgres, host: %s, user: %s, database: %s",
		s.config.Host,
		s.config.User,
		s.config.Database)
}

func (s *Stream) createSlot() (err error) {

	err = s.conn.CreateReplicationSlot(s.config.Replicate.Slot, s.config.Replicate.Plugin)
	if err == nil {
		return nil
	}

	if postgres.IsError(err, postgres.DuplicateObject) && s.config.Replicate.IgnoreDuplicateObjectError {
		return nil
	}

	return errors.Wrapf(err,
		"Something wrong with slot creation, slot: %s, plugin: %s",
		s.config.Replicate.Slot,
		s.config.Replicate.Plugin,
	)
}

func (s *Stream) start() (err error) {

	err = s.conn.StartReplication(s.config.Replicate.Slot, startLsn, timeLine, pluginArgs)
	return errors.Wrapf(err,
		"Something wrong with start replication, slot: %s, startLsn: %d, timeLine: %d, pluginArgs: %s",
		s.config.Replicate.Slot,
		startLsn,
		timeLine,
		pluginArgs,
	)
}
