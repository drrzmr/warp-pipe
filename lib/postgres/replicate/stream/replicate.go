package stream

import (
	"context"

	"github.com/jackc/pgx"
	"github.com/pkg/errors"

	"github.com/pagarme/warp-pipe/lib/log"
	"github.com/pagarme/warp-pipe/lib/postgres"
	"go.uber.org/zap"
)

const (
	pluginArgs = `("include-xids" '1', "include-timestamp" '1', "skip-empty-xacts" '0', "only-local" '0')`
	startLsn   = uint64(0)
	timeLine   = int64(-1)
)

var logger = log.Development("replicate")

// Replicate object
type Replicate struct {
	config   postgres.Config
	conn     *pgx.ReplicationConn
	listener EventListener
}

// New create a Replicate object
func New(config postgres.Config) *Replicate {

	return &Replicate{
		config:   config,
		conn:     nil,
		listener: nil,
	}
}

// Config return the address of Database config object
func (r *Replicate) Config() *postgres.Config {
	return &r.config
}

// Connect to postgres
func (r *Replicate) Connect() (err error) {

	logger.Debug("--> Connect()")
	defer logger.Debug("<-- Connect()")

	if r.isConnected() {
		logger.Info("already connected")
		return nil
	}

	if err = r.connect(); err != nil {
		logger.Error("connect error")
		return errors.WithStack(err)
	}

	if err = r.createSlot(); err != nil {
		logger.Error("create slot error")
		return errors.WithStack(err)
	}

	return nil
}

// Start replication
func (r *Replicate) Start(ctx context.Context, listener EventListener) (started bool, err error) {

	logger.Debug("--> Start()")
	defer logger.Debug("<-- Start()")

	if !r.isConnected() {
		logger.Error("not connected")
		return false, nil
	}

	if r.isStarted() {
		logger.Error("already started")
		return false, nil
	}

	if err = r.start(); err != nil {
		logger.Error("start error")
		return false, errors.WithStack(err)
	}

	r.listener = listener

	if err = r.listener.Run(ctx); err != nil {

		// filter context canceled
		canceled := errors.Cause(err) == context.Canceled
		logger.DebugIf(canceled, "context canceled")
		logger.ErrorIf(!canceled, "listener run", zap.Error(err))
		return false, errors.WithStack(err)
	}

	return true, nil
}

// SendStandByStatus method
func (r *Replicate) SendStandByStatus(position uint64) (err error) {

	var status *pgx.StandbyStatus

	if status, err = pgx.NewStandbyStatus(position); err != nil {
		return errors.Wrapf(err, "create new standby status object failed, position: %d", position)
	}

	err = r.conn.SendStandbyStatus(status)
	logger.DebugIf(err == nil, "sent standby status", zap.Uint64("position", position))

	return errors.Wrapf(err, "send stand by status failed, position: %d", position)
}

func (r *Replicate) isConnected() (connected bool) {

	return r.conn != nil
}

func (r *Replicate) isStarted() (connected bool) {

	return r.listener != nil
}

func (r *Replicate) connect() (err error) {

	config := pgx.ConnConfig{
		Host:     r.config.Host,
		Port:     r.config.Port,
		User:     r.config.User,
		Password: r.config.Password,
		Database: r.config.Database,
	}

	r.conn, err = pgx.ReplicationConnect(config)

	return errors.Wrapf(err,
		"Could not connect to postgres, host: %s, user: %s, database: %s",
		r.config.Host,
		r.config.User,
		r.config.Database)
}

func (r *Replicate) createSlot() (err error) {

	err = r.conn.CreateReplicationSlot(r.config.Replicate.Slot, r.config.Replicate.Plugin)
	if err == nil {
		return nil
	}

	if postgres.IsError(err, postgres.DuplicateObject) && r.config.Replicate.IgnoreDuplicateObjectError {
		return nil
	}

	return errors.Wrapf(err,
		"Something wrong with slot creation, slot: %s, plugin: %s",
		r.config.Replicate.Slot,
		r.config.Replicate.Plugin,
	)
}

func (r *Replicate) start() (err error) {

	err = r.conn.StartReplication(r.config.Replicate.Slot, startLsn, timeLine, pluginArgs)
	return errors.Wrapf(err,
		"Something wrong with start replication, slot: %s, startLsn: %d, timeLine: %d, pluginArgs: %s",
		r.config.Replicate.Slot,
		startLsn,
		timeLine,
		pluginArgs,
	)
}
