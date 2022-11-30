package gclog

/*
 * @Date: 2020-10-29 10:01:26
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-12-03 13:37:36
 */

import (
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
	fileutil "gitlab.yewifi.com/golden-cloud/common/file"
	"gitlab.yewifi.com/golden-cloud/common/gclog/formatters"
	"gopkg.in/natefinch/lumberjack.v2"
)

type OptionFunc func(*logrus.Logger) error

func WithDefaultJSONFormatter(jsoniterAPIs ...jsoniter.API) OptionFunc {

	return func(logger *logrus.Logger) error {
		var jsoniterAPI jsoniter.API
		if len(jsoniterAPIs) >= 1 {
			jsoniterAPI = jsoniterAPIs[0]
		}
		jsonFormatter := &formatters.JSONFormatter{
			TimestampFormat:  "2006-01-02 15:04:05.000",
			DisableTimestamp: false,
			DataKey:          "",
			FieldMap: formatters.FieldMap{
				logrus.FieldKeyTime:  "logTime",
				logrus.FieldKeyLevel: "level",
				logrus.FieldKeyMsg:   "message",
			},
			CallerPrettyfier: nil,
			PrettyPrint:      false,
			JSON:             jsoniterAPI,
		}

		logger.SetFormatter(jsonFormatter)
		return nil
	}
}

func WithOutputPath(dirPath, filename string) OptionFunc {
	return func(logger *logrus.Logger) error {
		var writer io.Writer
		var err error

		filePath := filepath.Join(dirPath, filename)
		if fileutil.PathExists(filePath) {
			writer, err = os.OpenFile(filePath, os.O_APPEND|os.O_WRONLY, os.ModeAppend)
		} else {
			writer, err = os.OpenFile(filePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
		}
		if err != nil {
			return err
		}

		if logger.Out == nil {
			logger.SetOutput(writer)
		} else {
			logger.SetOutput(io.MultiWriter(logger.Out, writer))
		}

		return nil
	}
}

func WithStdout(use ...bool) OptionFunc {
	return func(logger *logrus.Logger) error {
		useStdout := true

		if stdoutEnv := os.Getenv("STDOUT"); stdoutEnv != "" {
			if stdoutEnv == "false" || stdoutEnv == "0" || stdoutEnv == "off" {
				useStdout = false

			}

		} else if use != nil && len(use) > 0 && !use[0] {
			useStdout = false
		}

		if useStdout {
			if logger.Out == nil {
				logger.SetOutput(os.Stdout)
			} else {
				logger.SetOutput(io.MultiWriter(logger.Out, os.Stdout))
			}
		}
		return nil
	}
}

var (
	DefaultRotationFilePathChangeFunc func(string) (string, error) = defaultRotationFilePathChangeFunc
)

func defaultRotationFilePathChangeFunc(filePath string) (string, error) {
	// filePath 修改: /path/to/your.log => /path/to/your.2020-09-27.log
	ext := filepath.Ext(filePath)
	if ext == "" {
		filePath = filePath + ".%Y-%m-%d"
	} else {
		filePath = strings.TrimSuffix(filePath, ext) + ".%Y-%m-%d" + ext
	}
	return filePath, nil
}

func WithOutputPathAndRotation(dirPath, filename string, rotationFilePathChangeFunc func(string) (string, error), rotateOptions ...rotatelogs.Option) OptionFunc {

	return func(logger *logrus.Logger) error {
		var writer io.Writer
		var err error

		filePath := filepath.Join(dirPath, filename)
		if rotationFilePathChangeFunc == nil {
			rotationFilePathChangeFunc = DefaultRotationFilePathChangeFunc
		}
		filePath, err = rotationFilePathChangeFunc(filePath)
		if err != nil {
			return err
		}

		// 默认配置：Local时区，按天分割
		if rotateOptions == nil || len(rotateOptions) == 0 {
			rotateOptions = []rotatelogs.Option{
				rotatelogs.WithClock(rotatelogs.Local),
				rotatelogs.WithRotationTime(time.Hour * 24),
			}
		}

		writer, err = rotatelogs.New(
			filePath,
			rotateOptions...,
		)
		if err != nil {
			return err
		}

		if logger.Out == nil {
			logger.SetOutput(writer)
		} else {
			logger.SetOutput(io.MultiWriter(logger.Out, writer))
		}

		return nil
	}
}

func WithHook(h logrus.Hook) OptionFunc {
	return func(logger *logrus.Logger) error {
		// thread safe
		logger.AddHook(h)
		return nil
	}
}

func WithLogLevel(logLevel int) OptionFunc {
	return func(logger *logrus.Logger) error {
		logger.SetLevel(logrus.Level(logLevel))
		return nil
	}
}

// 日志切割
func WithRollingLog(filePath, fileName string) OptionFunc {
	return func(l *logrus.Logger) error {
		r := &lumberjack.Logger{
			Filename:   strings.TrimRight(filePath, "/") + "/" + fileName,
			MaxSize:    1024, //文件最大大小(MB)
			MaxBackups: 10,   //最多缓存文件数
			MaxAge:     30,   // 最多保留天数
			Compress:   false,
		}

		if l.Out == nil {
			l.SetOutput(r)
		} else {
			l.SetOutput(io.MultiWriter(l.Out, r))
		}

		return nil
	}

}
