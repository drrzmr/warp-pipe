package log

import (
	"sync"

	"go.uber.org/zap"
)

type register struct {
	name   string
	logger **zap.Logger
}

var (
	mainLogger   *zap.Logger
	registerList []register
	once         sync.Once
	zapConfig    = zap.NewDevelopmentConfig()
)

// Setup configure and create main logger
func Setup(config Config) {

	once.Do(func() {
		zapConfig.OutputPaths = []string{config.Stdout}
		zapConfig.ErrorOutputPaths = []string{config.Stderr}

		var err error
		if mainLogger, err = zapConfig.Build(); err != nil {
			panic(err)
		}

		for _, r := range registerList {
			*r.logger = mainLogger.Named(r.name)
		}
	})
}

// Register a logger with given name
func Register(logger **zap.Logger, name string) {

	registerList = append(registerList, register{
		name:   name,
		logger: logger,
	})

	if isSetuped() {
		*logger = mainLogger.Named(name)
	}
}

func isSetuped() (setuped bool) { return mainLogger != nil }
