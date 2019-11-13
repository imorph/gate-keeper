package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ardanlabs/conf"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/imorph/gate-keeper/pkg/version"
)

func main() {
	if err := run(); err != nil {
		log.Println("error :", err)
		os.Exit(1)
	}
}

func run() error {
	startTime := time.Now()

	var cfg struct {
		ListenHost string `conf:"default:127.0.0.1:10001"`
		LogLevel   string `conf:"default:info"`
	}

	if err := conf.Parse(os.Args[1:], "GK", &cfg); err != nil {
		if err == conf.ErrHelpWanted {
			usage, err := conf.Usage("GK", &cfg)
			if err != nil {
				return errors.Wrap(err, "generating config usage")
			}
			fmt.Println(usage)
			return nil
		}
		return errors.Wrap(err, "parsing config")
	}

	// configure logging
	logger, _ := initZap(cfg.LogLevel)
	defer func() {
		err := logger.Sync()
		if err != nil {
			// may show "sync /dev/stderr: invalid argument" but this is ok?
			// https://github.com/uber-go/zap/issues/328
			logger.Sugar().Info("error syncing logger", err)
		}
	}()

	stdLog := zap.RedirectStdLog(logger)
	defer stdLog()

	logger.Info("Application started",
		zap.Duration("startup_duration", time.Since(startTime)),
		zap.String("listen_address", cfg.ListenHost),
		zap.String("log_level", cfg.LogLevel),
		zap.String("app_name", version.GetAppName()),
		zap.String("version", version.GetVersion()),
		zap.String("revision", version.GetRevision()),
	)

	return nil
}

func initZap(logLevel string) (*zap.Logger, error) {
	level := zap.NewAtomicLevelAt(zapcore.InfoLevel)
	switch logLevel {
	case "debug":
		level = zap.NewAtomicLevelAt(zapcore.DebugLevel)
	case "info":
		level = zap.NewAtomicLevelAt(zapcore.InfoLevel)
	case "warn":
		level = zap.NewAtomicLevelAt(zapcore.WarnLevel)
	case "error":
		level = zap.NewAtomicLevelAt(zapcore.ErrorLevel)
	case "fatal":
		level = zap.NewAtomicLevelAt(zapcore.FatalLevel)
	case "panic":
		level = zap.NewAtomicLevelAt(zapcore.PanicLevel)
	default:
		log.Fatal("Dont know this log level:", logLevel, "known levels are: debug, info, warn, error, fatal, panic")
	}

	zapEncoderConfig := zapcore.EncoderConfig{
		TimeKey:        "ts",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	zapConfig := zap.Config{
		Level:       level,
		Development: false,
		Sampling: &zap.SamplingConfig{
			Initial:    100,
			Thereafter: 100,
		},
		Encoding:         "json",
		EncoderConfig:    zapEncoderConfig,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
	}

	return zapConfig.Build()
}
