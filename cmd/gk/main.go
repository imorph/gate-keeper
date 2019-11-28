package main

import (
	"fmt"
	"log"
	"os"
	"time"

	//"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"github.com/imorph/gate-keeper/pkg/server"
	"github.com/imorph/gate-keeper/pkg/signals"
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

	// ==================
	// Configuration
	var cfg struct {
		ListenHost     string
		LogLevel       string
		LoginThreshold int
		PassThreshold  int
		IPThreshold    int
		Logger         *zap.Logger
	}

	// command line flags
	pfs := pflag.NewFlagSet(version.GetAppName(), pflag.ContinueOnError)
	pfs.StringVar(&cfg.ListenHost, "listen-host", "127.0.0.1:10001", "ip:port server will listen on")
	pfs.StringVar(&cfg.LogLevel, "log-level", "info", "verbosity of logs, known levels are: debug, info, warn, error, fatal, panic")
	pfs.IntVar(&cfg.LoginThreshold, "max-login", 10, "Maximum alowed login attemps before ban")
	pfs.IntVar(&cfg.PassThreshold, "max-pass", 100, "Maximum alowed login attemps with same password before ban")
	pfs.IntVar(&cfg.IPThreshold, "max-ip", 1000, "Maximum alowed login attemps from single IP before ban")
	versionFlag := pfs.BoolP("version", "v", false, "get version number")

	// parse flags
	err := pfs.Parse(os.Args[1:])
	switch {
	case err == pflag.ErrHelp:
		os.Exit(0)
	case err != nil:
		pfs.PrintDefaults()
		return err
	case *versionFlag:
		fmt.Printf("%s-%s\n", version.GetVersion(), version.GetRevision())
		os.Exit(0)
	}

	// ==================
	// configure logging
	cfg.Logger, err = initZap(cfg.LogLevel)
	if err != nil {
		return err
	}
	defer func() {
		err := cfg.Logger.Sync()
		if err != nil {
			// may show "sync /dev/stderr: invalid argument" but this is ok?
			// https://github.com/uber-go/zap/issues/328
			cfg.Logger.Sugar().Warn("error syncing logger", err)
		}
	}()

	stdLog := zap.RedirectStdLog(cfg.Logger)
	defer stdLog()
	s := server.NewGateKeeperServer(cfg.ListenHost, cfg.Logger, cfg.IPThreshold, cfg.LoginThreshold, cfg.PassThreshold)
	go func(s *server.GateKeeperServer) {
		err = s.Start()
		if err != nil {
			cfg.Logger.Error("Error starting server", zap.Error(err))
			os.Exit(1)
		}
	}(s)

	cfg.Logger.Info("Application started",
		zap.Duration("startup_duration", time.Since(startTime)),
		zap.String("listen_address", cfg.ListenHost),
		zap.String("log_level", cfg.LogLevel),
		zap.String("app_name", version.GetAppName()),
		zap.String("version", version.GetVersion()),
		zap.String("revision", version.GetRevision()),
	)

	sig := signals.WaitForSigterm()

	//==========================
	//stop sequence
	sigtime := time.Now()
	s.Stop()
	cfg.Logger.Sugar().Info("OS signal received, stopping. ", "signal: ", sig)
	cfg.Logger.Info("Application stopped",
		zap.Duration("stop_duration", time.Since(sigtime)),
		zap.String("app_name", version.GetAppName()),
		zap.String("version", version.GetVersion()),
		zap.String("revision", version.GetRevision()),
	)

	return nil
}

//
// zap initialisation logic
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
