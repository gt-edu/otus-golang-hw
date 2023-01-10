package logger

import (
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Sync() error
	Debug(msg string)
	Info(msg string)
	Warn(msg string)
	Error(msg string)
}

type ZapLogger struct {
	zapLogger *zap.Logger
}

func New(level, preset string) (*ZapLogger, error) {
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

	log := &ZapLogger{}
	log.zapLogger, err = config.Build()
	if err != nil {
		return nil, err
	}

	return log, nil
}

func (l ZapLogger) Sync() error {
	return l.zapLogger.Sync()
}

func (l ZapLogger) Debug(msg string) {
	l.zapLogger.Debug(msg)
}

func (l ZapLogger) Info(msg string) {
	l.zapLogger.Info(msg)
}

func (l ZapLogger) Warn(msg string) {
	l.zapLogger.Warn(msg)
}

func (l ZapLogger) Error(msg string) {
	l.zapLogger.Error(msg)
}

// TODO
