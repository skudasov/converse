package log

import (
	"go.uber.org/zap"
)

var L *zap.SugaredLogger

func init() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	L = logger.Named("main").Sugar()
}
