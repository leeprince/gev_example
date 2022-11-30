package logger

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"sync"

	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/logger/hooks"
)

var onceInit sync.Once

var log *logrus.Entry

var logPath = ""

var logName = "application.log"

var appName = "app"

var customHooks CustomHooks

var withFields CustomFields

var writer io.Writer = nil

var addStdoutWriter bool

var logLevel logrus.Level = logrus.InfoLevel

var loggers sync.Map

func SetDefaultLogger(l *logrus.Entry) {
	log = l
	logger = l

	if l != nil && l.Logger != nil {
		formatter = l.Logger.Formatter
	}

}

func SetLogPath(p string) {
	logPath = p
}

func SetLogName(l string) {
	logName = l
}

func SetAppName(a string) {
	appName = a
}

func SetHooks(h CustomHooks) {
	customHooks = h
}

func SetWithFields(f CustomFields) {
	withFields = f
}

func SetWriter(w io.Writer) {
	writer = w
}

func SetLogLevel(l uint32) {
	logLevel = logrus.Level(l)
}

func AddStdoutWriter() {
	addStdoutWriter = true
}

// newDayLogger 按天分割日志
func newDayLogger(loggerName string, formatter logrus.Formatter) (*logrus.Entry, error) {

	log := logrus.New()
	log.Formatter = formatter //日志格式
	log.Level = logLevel

	if !addStdoutWriter { //默认不输出到 标准输出
		log.Out = ioutil.Discard
	}

	//添加默认的所有hook
	hks := defaultHooks()
	for _, hook := range hks {
		log.Hooks.Add(hook)
	}

	//添加按天切割日志的hook
	filepath := path.Join(logPath, loggerName+"_%Y%m%d.log")
	rHook, err := hooks.NewRotateHook(filepath, log.Formatter)
	if err != nil {
		return nil, err
	}

	log.Hooks.Add(rHook)

	withFields := defaultFields(appName)
	logger := log.WithFields(withFields)

	return logger, nil
}

// GetDayLogger (json日志格式)
func GetDayLogger(logName string) *logrus.Entry {
	key := logName + "_day"
	iLogger, ok := loggers.Load(key)

	if !ok { //初始化新logger(默认json格式)
		logger, err := newDayLogger(logName, formatter)
		if err != nil {
			logger = GetLogger()
		}

		loggers.Store(key, logger)
		return logger
	}

	//已有对应logName的logger存在
	logger, ok := iLogger.(*logrus.Entry)
	if !ok || logger == nil {
		return GetLogger()
	}

	return logger
}

// GetStatLogger (stat格式)
func GetStatLogger(logName string) *logrus.Entry {
	key := logName + "_stat"
	iLogger, ok := loggers.Load(key)

	if !ok { //初始化新logger(stat格式)
		logger, err := newDayLogger(logName, statFormat)
		if err != nil {
			logger = GetLogger()
		}

		loggers.Store(key, logger)
		return logger
	}

	//已有对应logName的logger存在
	logger, ok := iLogger.(*logrus.Entry)
	if !ok || logger == nil {
		return GetLogger()
	}

	return logger
}

func InitLog() {
	onceInit.Do(func() {
		if logPath == "" {
			logPath = "./"
		}
		logFile := logPath + "/" + logName
		fmt.Println("logFile --- ", logFile)
		var w io.Writer

		if writer == nil {
			var err error
			w, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			if err != nil && os.IsNotExist(err) {
				_, err = os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0777)
				if err != nil {
					panic("创建日志文件错误" + err.Error())
				}
				w, err = os.OpenFile(logFile, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
			}
			if err != nil {
				panic("打开日志文件错误" + err.Error())
			}
		} else {
			w = writer
		}

		if addStdoutWriter {
			w = io.MultiWriter(w, os.Stdout)
		}
		InitLogger(appName, w, withFields, customHooks, logLevel)
		log = GetLogger()
	})

}
