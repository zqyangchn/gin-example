package gormzap

import (
	"context"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// LogLevel
type LogLevel int

const (
	Silent LogLevel = iota + 1
	Error
	Warn
	Info
)

type Config struct {
	SlowThreshold time.Duration
	LogLevel      LogLevel
}

func New(zap *zap.Logger, config Config) logger.Interface {
	return &Logger{
		Logger: *zap,
		Config: config,
	}
}

type Logger struct {
	zap.Logger
	Config
}

// LogMode log mode
func (g *Logger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *g
	newLogger.LogLevel = LogLevel(level)
	return &newLogger
}

// Info print info
func (g Logger) Info(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel < Info {
		return
	}
	g.Logger.Sugar().Debugf(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
}

// Warn print warn messages
func (g Logger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel < Warn {
		return
	}
	g.Logger.Sugar().Warnf(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
}

// Error print error messages
func (g Logger) Error(ctx context.Context, msg string, data ...interface{}) {
	if g.LogLevel < Error {
		return
	}
	g.Logger.Sugar().Errorf(msg, append([]interface{}{utils.FileWithLineNum()}, data...))
}

// Trace print sql message
func (g Logger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	if g.LogLevel >= Silent {
		elapsed := time.Since(begin)
		switch {
		// gorm error log
		case err != nil && logger.LogLevel(g.LogLevel) >= logger.Error:
			sql, rows := fc()
			g.Logger.Error("gorm",
				zap.String("caller", utils.FileWithLineNum()),
				zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1e6),
				zap.String("sql", sql),
				zap.Int64("affect_rows", rows),
				zap.Error(err),
			)
		// gorm warning log
		case elapsed > g.SlowThreshold && g.SlowThreshold != 0 && logger.LogLevel(g.LogLevel) >= logger.Warn:
			sql, rows := fc()
			g.Logger.Warn("gorm",
				zap.String("caller", utils.FileWithLineNum()),
				zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1e6),
				zap.String("sql", sql),
				zap.Int64("affect_rows", rows),
			)
		// gorm info log, it can be considered all log or debug log and so on ...
		case logger.LogLevel(g.LogLevel) >= logger.Info:
			sql, rows := fc()
			g.Logger.Debug("gorm",
				zap.String("caller", utils.FileWithLineNum()),
				zap.Float64("elapsed_ms", float64(elapsed.Nanoseconds())/1e6),
				zap.String("sql", sql),
				zap.Int64("affect_rows", rows),
			)
		}
	}
}
