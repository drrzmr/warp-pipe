package log

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger object
type Logger struct {
	zap.Logger
}

var development *zap.Logger

func init() {
	var err error
	development, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

// Development return a child logger for given module
func Development(module string) *Logger {

	zapLogger := development.With(zap.String("module", module))

	return &Logger{
		*zapLogger,
	}
}

// DebugIf helper
func (l *Logger) DebugIf(cond bool, msg string, fields ...zapcore.Field) {

	if !cond {
		return
	}

	l.Debug(msg, fields...)
}

// ErrorIf helper
func (l *Logger) ErrorIf(cond bool, msg string, fields ...zapcore.Field) {

	if !cond {
		return
	}

	l.Error(msg, fields...)
}
