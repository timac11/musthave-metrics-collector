package logger

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger = zap.NewNop().Sugar()

// Initialize init private singleton logger in Sugar mode with production config
func Initialize(level string) error {
	lvl, err := zap.ParseAtomicLevel(level)
	if err != nil {
		return err
	}
	cfg := zap.NewProductionConfig()
	cfg.Level = lvl

	zl, err := cfg.Build()
	if err != nil {
		return err
	}
	logger = zl.Sugar()
	return nil
}

func Debug(msg string, params any) {
	logger.Debug(msg, params)
}

func Info(msg string, params any) {
	logger.Info(msg, params)
}

func Error(msg string, params any) {
	logger.Error(msg, params)
}
