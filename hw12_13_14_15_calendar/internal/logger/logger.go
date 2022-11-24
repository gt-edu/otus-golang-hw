package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger struct {
	zapLogger *zap.Logger
}

func New(level, preset string) (*Logger, error) {
	var config zap.Config
	var err error

	switch preset {
	case "development":
		config = zap.NewDevelopmentConfig()
	case "production":
		config = zap.NewProductionConfig()
	default:
		return nil, errors.New("unknown logging preset")
	}

	l := new(zapcore.Level)
	err = l.UnmarshalText([]byte(level))
	if err != nil {
		return nil, err
	}
	config.Level = zap.NewAtomicLevelAt(*l)

	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	log := &Logger{}
	log.zapLogger, err = config.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (l Logger) Sync() error {
	return l.zapLogger.Sync()
}

func (l Logger) Info(msg string) {
	l.zapLogger.Info(msg)
}

func (l Logger) Error(msg string) {
	l.zapLogger.Error(msg)
}

// TODO
