package logger

import (
	"go.uber.org/zap"
)

var Logger *zap.SugaredLogger

func Init(level string) {
	var config zap.Config

	if level == "development" {
		config = zap.NewDevelopmentConfig()
	} else {
		config = zap.NewProductionConfig()
	}

	// Set log level
	switch level {
	case "debug":
		config.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
	case "info":
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	case "warn":
		config.Level = zap.NewAtomicLevelAt(zap.WarnLevel)
	case "error":
		config.Level = zap.NewAtomicLevelAt(zap.ErrorLevel)
	default:
		config.Level = zap.NewAtomicLevelAt(zap.InfoLevel)
	}

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}

	Logger = logger.Sugar()
}

func Debug(args ...interface{}) {
	Logger.Debug(args...)
}

func Info(args ...interface{}) {
	Logger.Info(args...)
}

func Warn(args ...interface{}) {
	Logger.Warn(args...)
}

func Error(args ...interface{}) {
	Logger.Error(args...)
}

func Fatal(args ...interface{}) {
	Logger.Fatal(args...)
}

func Debugf(template string, args ...interface{}) {
	Logger.Debugf(template, args...)
}

func Infof(template string, args ...interface{}) {
	Logger.Infof(template, args...)
}

func Warnf(template string, args ...interface{}) {
	Logger.Warnf(template, args...)
}

func Errorf(template string, args ...interface{}) {
	Logger.Errorf(template, args...)
}

func Fatalf(template string, args ...interface{}) {
	Logger.Fatalf(template, args...)
}
