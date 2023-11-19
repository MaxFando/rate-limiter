package logger

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/utils"

	"go.uber.org/zap"
	gormLogger "gorm.io/gorm/logger"
)

const (
	errorLogTemplate = "[PgSQL ERROR] [%.3fms] [rows:%v] [file:%s] %s"
	warnLogTemplate  = "[PgSQL WARN] [%.3fms] [rows:%v] %s %s"
	infoLogTemplate  = "[PgSQL INFO] [%.3fms] [rows:%v] %s"
)

type GormZapLoggerWrapper struct {
	logger                    *zap.SugaredLogger
	logLevel                  gormLogger.LogLevel
	IgnoreRecordNotFoundError bool
	SlowThreshold             time.Duration
}

func NewGormZapLoggerWrapper(logger *zap.SugaredLogger) *GormZapLoggerWrapper {
	return &GormZapLoggerWrapper{logger: logger}
}

func (l *GormZapLoggerWrapper) LogMode(level gormLogger.LogLevel) gormLogger.Interface {
	l.logLevel = level
	return l
}

func (l *GormZapLoggerWrapper) Info(ctx context.Context, msg string, data ...interface{}) {
	if strings.Contains(msg, "replacing callback") {
		return
	}

	if l.logLevel >= gormLogger.Info {
		l.logger.Infof(msg, data...)
	}
}

func (l *GormZapLoggerWrapper) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Warn {
		l.logger.Warnf(msg, data...)
	}
}

func (l *GormZapLoggerWrapper) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.logLevel >= gormLogger.Error {
		l.logger.Errorf(msg, data...)
	}
}

func (l *GormZapLoggerWrapper) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if l.logLevel <= gormLogger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	switch {
	case err != nil && l.logLevel >= gormLogger.Error && (!errors.Is(err, gormLogger.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.Error(ctx, errorLogTemplate, float64(elapsed.Nanoseconds())/1e6, rows, utils.FileWithLineNum(), sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.logLevel >= gormLogger.Warn:
		slowLog := fmt.Sprintf("SLOW SQL >= %v", l.SlowThreshold)
		l.Warn(ctx, warnLogTemplate, float64(elapsed.Nanoseconds())/1e6, rows, slowLog, sql)
	case l.logLevel == gormLogger.Info:
		l.Info(ctx, infoLogTemplate, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}
