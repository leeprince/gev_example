package hooks

import (
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/rifflock/lfshook"
	"github.com/sirupsen/logrus"
)

// RotateHook 日志分割
type RotateHook struct {
	logPath string
	format  logrus.Formatter
}

// NewRotateHook 获取按天切割日志的hook对象
func NewRotateHook(logPath string, format logrus.Formatter) (logrus.Hook, error) {
	//获取按天切割的writer对象
	writer, err := rotatelogs.New(logPath)
	if err != nil {
		return nil, err
	}

	//返回hook
	lfsHook := lfshook.NewHook(lfshook.WriterMap{
		logrus.DebugLevel: writer,
		logrus.InfoLevel:  writer,
		logrus.WarnLevel:  writer,
		logrus.ErrorLevel: writer,
		logrus.FatalLevel: writer,
		logrus.PanicLevel: writer,
		logrus.TraceLevel: writer,
	}, format)

	return lfsHook, nil
}
