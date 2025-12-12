package utils

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func getLogLevel(debug bool) zapcore.Level {
	if debug {
		return zapcore.DebugLevel
	}
	return zapcore.InfoLevel
}

func getTimeEncoder(development bool) zapcore.TimeEncoder {
	if development {
		return func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format("2006-01-02 15:04:05"))
		}
	}
	return zapcore.EpochMillisTimeEncoder
}

func InitLogger() *zap.Logger {
	development := os.Getenv("DEVELOPPEMENT") == "1"
	debug := os.Getenv("DEBUG") == "1"

	var cfg zap.Config
	if development {
		cfg = zap.NewDevelopmentConfig()
		cfg.Level = zap.NewAtomicLevelAt(getLogLevel(debug))
		cfg.Encoding = "console"
		cfg.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		cfg = zap.NewProductionConfig()
		cfg.Level = zap.NewAtomicLevelAt(getLogLevel(debug))
		cfg.EncoderConfig = zap.NewProductionEncoderConfig()
	}
	cfg.DisableStacktrace = true
	cfg.EncoderConfig.TimeKey = "ts"
	cfg.EncoderConfig.EncodeTime = getTimeEncoder(development)

	opts := []zap.Option{
		zap.AddStacktrace(zap.DPanicLevel),
	}
	if !development {
		opts = append(opts, zap.WithCaller(false))
	}

	logger, err := cfg.Build(opts...)

	if err != nil {
		panic(err)
	}
	return logger
}
