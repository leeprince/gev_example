package gclog

/*
 * @Date: 2020-10-29 15:52:25
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2021-02-19 11:28:13
 */

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	gormloggerlib "gorm.io/gorm/logger"
)

var _ gormloggerlib.Interface = (*GormLogger)(nil)

// Notes:
// 	gorm logger interface:
// 	type Interface interface {
// 		LogMode(LogLevel) Interface
// 		Info(context.Context, string, ...interface{})
// 		Warn(context.Context, string, ...interface{})
// 		Error(context.Context, string, ...interface{})
// 		Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error)
// 	}

type GormLogger struct {
	logrusLogger  *logrus.Logger
	slowThreshold time.Duration

	gormLoggerLevel gormloggerlib.LogLevel
}

func NewGormLogger(logrusLogger *logrus.Logger, slowThreshold time.Duration) *GormLogger {
	return &GormLogger{
		logrusLogger:  logrusLogger,
		slowThreshold: slowThreshold,
	}
}

func (l *GormLogger) LogMode(level gormloggerlib.LogLevel) gormloggerlib.Interface {
	// Notes:
	// gorm logger LogLevel =>
	// 	Silent: 1
	// 	Error: 2
	// 	Warn: 3
	// 	Info: 4
	l.gormLoggerLevel = level
	return l
}

func (l *GormLogger) Info(ctx context.Context, format string, args ...interface{}) {
	if l.gormLoggerLevel < gormloggerlib.Info {
		return
	}
	logID := GetLogIDFromCtx(ctx)
	l.logrusLogger.
		WithField("gorm.context", ctx).
		WithField("gorm.message", fmt.Sprintf(format, args...)).
		Info(logID, "")
}

func (l *GormLogger) Warn(ctx context.Context, format string, args ...interface{}) {
	if l.gormLoggerLevel < gormloggerlib.Warn {
		return
	}
	logID := GetLogIDFromCtx(ctx)
	l.logrusLogger.
		WithField("gorm.context", ctx).
		WithField("gorm.message", fmt.Sprintf(format, args...)).
		Warn(logID, "")
}

func (l *GormLogger) Error(ctx context.Context, format string, args ...interface{}) {
	if l.gormLoggerLevel < gormloggerlib.Error {
		return
	}
	logID := GetLogIDFromCtx(ctx)
	l.logrusLogger.
		WithField("gorm.context", ctx).
		WithField("gorm.message", fmt.Sprintf(format, args...)).
		Error(logID, "")
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	logID := GetLogIDFromCtx(ctx)

	if l.gormLoggerLevel > 0 {
		elapsed := time.Since(begin)
		switch {
		// NOTES(aiden.deng): 不打印 record not found 的错误
		case err != nil && err != gorm.ErrRecordNotFound && l.gormLoggerLevel >= gormloggerlib.Error:
			sql, rows := fc()
			t := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)
			l.logrusLogger.
				WithError(err).
				WithField("gorm.millisecond", elapsed.Milliseconds()).
				WithField("gorm.time", t).
				WithField("gorm.rows", rows).
				WithField("gorm.sql", sql).
				Error(logID, fmt.Sprintf("gorm logger: err = %+v", err))

		case elapsed > l.slowThreshold && l.slowThreshold != 0 && l.gormLoggerLevel >= gormloggerlib.Warn:
			sql, rows := fc()
			t := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)

			l.logrusLogger.
				WithField("gorm.millisecond", elapsed.Milliseconds()).
				WithField("gorm.time", t).
				WithField("gorm.rows", rows).
				WithField("gorm.sql", sql).
				Warn(logID, fmt.Sprintf("gorm logger: time = %s", t))

		case l.gormLoggerLevel >= gormloggerlib.Info:
			sql, rows := fc()
			t := fmt.Sprintf("%.3fms", float64(elapsed.Nanoseconds())/1e6)

			l.logrusLogger.
				WithField("gorm.millisecond", elapsed.Milliseconds()).
				WithField("gorm.time", t).
				WithField("gorm.rows", rows).
				WithField("gorm.sql", sql).
				Info(logID, fmt.Sprintf("gorm logger: sql = %s", sql))
		}
	}
}
