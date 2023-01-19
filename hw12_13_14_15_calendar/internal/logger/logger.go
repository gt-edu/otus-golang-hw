package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	Sync() error
	Debug(msg string, fields ...zapcore.Field)
	Info(msg string, fields ...zapcore.Field)
	Warn(msg string, fields ...zapcore.Field)
	Error(msg string, fields ...zapcore.Field)
}

func New(level, preset string) (Logger, error) {
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

	return config.Build()
}

func SafeLoggerSync(logg Logger) {
	func(logg Logger) {
		err := logg.Sync()
		if err != nil && !strings.Contains(err.Error(), "sync /dev/stderr: invalid argument") {
			_, _ = fmt.Fprintf(os.Stderr, "Error during logger syncing: %v", err)
		}
	}(logg)
}
