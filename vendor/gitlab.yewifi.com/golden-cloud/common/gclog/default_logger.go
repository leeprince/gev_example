package gclog

/*
 * @Date: 2020-10-29 15:50:06
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2021-06-22 16:07:42
 */

import (
	"sync"

	"github.com/sirupsen/logrus"
	commonlogger "gitlab.yewifi.com/golden-cloud/common/logger"
)

var (
	defaultLogger *logrus.Logger
	once          sync.Once
)

func init() {
	once.Do(func() {
		defaultLogger = NewMinimum("gclog", int(logrus.DebugLevel), "")
	})
}

func SetDefaultLogger(logger *logrus.Logger) {
	defaultLogger = logger
	commonlogger.SetDefaultLogger(logrus.NewEntry(logger))
}

func GetDefaultLogger() *logrus.Logger {
	return defaultLogger
}

func WithField(key string, value interface{}) *logrus.Entry {
	return defaultLogger.WithField(key, value)
}

func WithFields(fields map[string]interface{}) *logrus.Entry {
	return defaultLogger.WithFields(fields)
}

func WithError(err error) *logrus.Entry {
	return defaultLogger.WithError(err)
}

func Trace(args ...interface{}) {
	defaultLogger.Trace(args...)
}

func Debug(args ...interface{}) {
	defaultLogger.Debug(args...)
}

func Info(args ...interface{}) {
	defaultLogger.Info(args...)
}

func Infof(format string, args ...interface{}) {
	defaultLogger.Infof(format, args...)
}

func Warn(args ...interface{}) {
	defaultLogger.Warn(args...)
}

func Error(args ...interface{}) {
	defaultLogger.Error(args...)
}

func Errorf(format string, args ...interface{}) {
	defaultLogger.Errorf(format, args...)
}

func Fatal(args ...interface{}) {
	defaultLogger.Fatal(args...)
}

func Panic(args ...interface{}) {
	defaultLogger.Panic(args...)
}

func PanicIfError(err error) {
	if err == nil {
		return
	}
	defaultLogger.Panic(err)
}
