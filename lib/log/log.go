package log

import (
	"go.uber.org/zap"
)

var development *zap.Logger

func init() {
	var err error
	development, err = zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
}

// Development return a child logger for given module
func Development(module string) *zap.Logger {

	return development.With(zap.String("module", module))
}
