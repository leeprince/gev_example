package gclog

/*
 * @Date: 2020-10-28 17:41:29
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2022-01-10 11:12:08
 */

import (
	"os"

	jsoniter "github.com/json-iterator/go"
	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/gclog/hooks"
)

func New(optionFuncs ...OptionFunc) (*logrus.Logger, error) {
	logger := logrus.New()
	logger.Out = nil
	for _, optionFunc := range optionFuncs {
		err := optionFunc(logger)
		if err != nil {
			return nil, err
		}
	}

	// 如果logger的writer为nil，那么设为os.Stdout，防止logrus panic
	if logger.Out == nil {
		logger.Out = os.Stdout
	}

	return logger, nil
}

func MustNew(optionFuncs ...OptionFunc) *logrus.Logger {
	logger, err := New(optionFuncs...)
	if err != nil {
		panic(err)
	}
	return logger
}

func mayHostname() string {
	hostname, _ := os.Hostname()
	return hostname
}

func NewMinimum(appName string, logLevel int, sentryDsn string) *logrus.Logger {
	logger := MustNew(
		// 配置 formatter
		WithDefaultJSONFormatter(jsoniter.ConfigCompatibleWithStandardLibrary),

		// 配置日志level
		WithLogLevel(logLevel),

		// 开启控制台标准输出
		WithStdout(),

		// logID解析
		WithHook(hooks.NewLogIDHook("context", "logID", 16, "未能获取logID")),

		// caller寻址器
		WithHook(hooks.NewAutomaticCallerHook("context", "caller", 7, 11, []string{"gclog/logger.go", "gclog/default_logger.go", "gclog/gorm_logger.go", "logrus"})),

		// immutable字段
		WithHook(hooks.NewConstFieldsHook("context").Add("app", appName).Add("hostname", mayHostname())),

		// 调用logrus.Logger.WithField(key, value)时 字段存放的根路径
		WithHook(hooks.NewDataHook("data")),

		// sentry hook会根据字符串是否为空去判断是否开启sentry日志告警功能
		WithHook(hooks.MustNewSentryHook(sentryDsn)),
	)
	return logger
}

func Default(appName string, logDirPath string, logLevel int, sentryDsn string) *logrus.Logger {
	logger := MustNew(
		// 配置 formatter
		WithDefaultJSONFormatter(jsoniter.ConfigCompatibleWithStandardLibrary),

		// 配置日志level
		WithLogLevel(logLevel),

		// 默认不开启控制台标准输出
		WithStdout(false),

		// 默认不使用日志分割
		WithOutputPath(logDirPath, "application.log"),

		// logID解析
		WithHook(hooks.NewLogIDHook("context", "logID", 16, "未能获取logID")),

		// caller寻址器
		WithHook(hooks.NewAutomaticCallerHook("context", "caller", 7, 11, []string{"gclog/logger.go", "gclog/default_logger.go", "gclog/gorm_logger.go", "logrus"})),

		// immutable字段
		WithHook(hooks.NewConstFieldsHook("context").Add("app", appName).Add("hostname", mayHostname())),

		// 调用logrus.Logger.WithField(key, value)时 字段存放的根路径
		WithHook(hooks.NewDataHook("data")),

		// sentry hook会根据字符串是否为空去判断是否开启sentry日志告警功能
		WithHook(hooks.MustNewSentryHook(sentryDsn)),
	)
	return logger
}

// 使用日志滚动
func DefaultV2(appName string, logDirPath string, logLevel int, sentryDsn string) *logrus.Logger {
	logger := MustNew(
		// 配置 formatter
		WithDefaultJSONFormatter(jsoniter.ConfigCompatibleWithStandardLibrary),

		// 配置日志level
		WithLogLevel(logLevel),

		// 默认不开启控制台标准输出
		WithStdout(false),

		// 默认使用日志滚动
		WithRollingLog(logDirPath, "application.log"),

		// logID解析
		WithHook(hooks.NewLogIDHook("context", "logID", 16, "未能获取logID")),

		// caller寻址器
		WithHook(hooks.NewAutomaticCallerHook("context", "caller", 7, 11, []string{"gclog/logger.go", "gclog/default_logger.go", "gclog/gorm_logger.go", "logrus"})),

		// immutable字段
		WithHook(hooks.NewConstFieldsHook("context").Add("app", appName).Add("hostname", mayHostname())),

		// 调用logrus.Logger.WithField(key, value)时 字段存放的根路径
		WithHook(hooks.NewDataHook("data")),

		// sentry hook会根据字符串是否为空去判断是否开启sentry日志告警功能
		WithHook(hooks.MustNewSentryHook(sentryDsn)),
	)
	return logger
}
