package utils

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type ZapLoggerWrapper struct {
	zapLogger *zap.Logger
}

var (
	Logger *ZapLoggerWrapper
	once   sync.Once
)

func (l *ZapLoggerWrapper) GetLogger() *zap.Logger {
	return l.zapLogger
}

func InitializeLogger() {
	var lvl zapcore.Level
	switch strings.ToUpper(strings.TrimSpace(os.Getenv("LOG_LEVEL"))) {
	case "ERR", "ERROR":
		lvl = zapcore.ErrorLevel
	case "WARN", "WARNING":
		lvl = zapcore.WarnLevel
	case "INFO":
		lvl = zapcore.InfoLevel
	case "DEBUG":
		lvl = zapcore.DebugLevel
	case "FATAL":
		lvl = zapcore.FatalLevel
	default:
		lvl = zapcore.InfoLevel
	}

	once.Do(func() {
		globalLevel := lvl

		highPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= zapcore.ErrorLevel
		})
		lowPriority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
			return lvl >= globalLevel && lvl < zapcore.ErrorLevel
		})
		consoleInfos := zapcore.Lock(os.Stdout)
		consoleErrors := zapcore.Lock(os.Stderr)

		zapConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",
			LevelKey:       "level",
			CallerKey:      "caller",
			MessageKey:     "message",
			StacktraceKey:  "stacktrace",
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.LowercaseLevelEncoder,
			EncodeTime:     zapcore.ISO8601TimeEncoder,
			EncodeDuration: zapcore.SecondsDurationEncoder,
			EncodeCaller:   zapcore.ShortCallerEncoder,
		}
		encoder := zapcore.NewJSONEncoder(zapConfig)

		core := zapcore.NewTee(
			zapcore.NewCore(encoder, consoleErrors, highPriority),
			zapcore.NewCore(encoder, consoleInfos, lowPriority),
		)

		zapLogger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
		zap.RedirectStdLog(zapLogger)

		Logger = &ZapLoggerWrapper{zapLogger: zapLogger}
	})
}

func (l *ZapLoggerWrapper) Debug(msg string, fields ...zap.Field) {
	l.zapLogger.Debug(msg, fields...)
}

func (l *ZapLoggerWrapper) Info(msg string, fields ...zap.Field) {
	l.zapLogger.Info(msg, fields...)
}

func (l *ZapLoggerWrapper) Warn(msg string, fields ...zap.Field) {
	l.zapLogger.Warn(msg, fields...)
}

func (l *ZapLoggerWrapper) Error(msg interface{}, fields ...zap.Field) {
	if msg == nil {
		return
	}

	switch v := msg.(type) {
	case string:
		l.zapLogger.Error(v, fields...)
	case error:
		l.zapLogger.Error(v.Error(), fields...)
	case fmt.Stringer:
		l.zapLogger.Error(v.String(), fields...)
	default:
		l.zapLogger.Error(fmt.Sprintf("%v", v), fields...)
	}
}

func (l *ZapLoggerWrapper) Fatal(msg interface{}, fields ...zap.Field) {
	if msg == nil {
		return
	}

	switch v := msg.(type) {
	case string:
		l.zapLogger.Fatal(v, fields...)
	case error:
		l.zapLogger.Fatal(v.Error(), fields...)
	case fmt.Stringer:
		l.zapLogger.Fatal(v.String(), fields...)
	default:
		l.zapLogger.Fatal(fmt.Sprintf("%v", msg), fields...)
	}
}
