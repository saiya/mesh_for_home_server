package logger

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger *zap.Logger
var sugar *zap.SugaredLogger
var loglevel *zap.AtomicLevel

type loggerContextKeyType struct{}

var loggerContextKey = loggerContextKeyType{}

func init() {
	cfg := zap.NewProductionConfig()
	cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	loglevel = &cfg.Level
	if os.Getenv("DEBUG") != "" {
		EnableDebugLog()
	}

	var err error
	logger, err = cfg.Build()
	neverFail(err)
	sugar = logger.Sugar()
}

func Wrap(ctx context.Context, keyAndValues ...interface{}) context.Context {
	return context.WithValue(
		ctx,
		loggerContextKey,
		GetFrom(ctx).With(keyAndValues...),
	)
}

func GetFrom(ctx context.Context) *zap.SugaredLogger {
	logger, ok := ctx.Value(loggerContextKey).(*zap.SugaredLogger)
	if ok {
		return logger
	}
	return Get()
}

func Get() *zap.SugaredLogger {
	return sugar
}

func EnableDebugLog() {
	loglevel.SetLevel(zap.DebugLevel)
}
