package logger

import "go.uber.org/zap"

var logger *zap.Logger
var sugar *zap.SugaredLogger

func init() {
	var err error
	logger, err = zap.NewProduction()
	neverFail(err)
	sugar = logger.Sugar()
}

func Get() *zap.SugaredLogger {
	return sugar
}
