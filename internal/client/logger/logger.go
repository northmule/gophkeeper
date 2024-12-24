package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	*zap.SugaredLogger
}

func NewLogger(level string) (*Logger, error) {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return nil, err
	}
	cfg := zap.NewDevelopmentConfig()
	cfg.Level = lvl
	appLogger, err := cfg.Build()
	if err != nil {
		return nil, err
	}

	logger := &Logger{appLogger.Sugar()}

	return logger, nil
}
