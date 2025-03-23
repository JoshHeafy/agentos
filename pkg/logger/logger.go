package logger

import (
	"context"
	"os"

	"github.com/rs/zerolog"
)

type Adapter struct {
	log zerolog.Logger
}

func NewZeroLog(serviceName string) Adapter {
	logger := zerolog.New(os.Stdout).
		With().
		Timestamp().
		Str("service", serviceName).
		Caller().
		Logger()

	return Adapter{
		log: logger,
	}
}

func (a Adapter) Debug(ctx context.Context, msg string, keysAndValues ...interface{}) {
	a.log.Debug().Ctx(ctx).Fields(keysAndValues).Caller(1).Msg(msg)
}

func (a Adapter) Info(ctx context.Context, msg string, keysAndValues ...interface{}) {
	a.log.Info().Ctx(ctx).Fields(keysAndValues).Caller(1).Msg(msg)
}

func (a Adapter) Warn(ctx context.Context, msg string, keysAndValues ...interface{}) {
	a.log.Warn().Ctx(ctx).Fields(keysAndValues).Caller(1).Msg(msg)
}

func (a Adapter) Error(ctx context.Context, msg string, keysAndValues ...interface{}) {
	a.log.Error().Ctx(ctx).Fields(keysAndValues).Caller(1).Msg(msg)
}

func (a Adapter) Fatal(ctx context.Context, msg string, keysAndValues ...interface{}) {
	a.log.Fatal().Ctx(ctx).Fields(keysAndValues).Caller(1).Msg(msg)
}
