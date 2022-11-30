package gclog

/*
 * @Date: 2020-10-29 15:51:38
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2021-01-07 11:22:13
 */

import (
	jsoniter "github.com/json-iterator/go"
	"gitlab.yewifi.com/golden-cloud/common/gclog/hooks"
	commonlogger "gitlab.yewifi.com/golden-cloud/common/logger"
)

func SetUp(cfg *commonlogger.LoggerConf, appName string) {
	optionFuncs := []OptionFunc{
		WithDefaultJSONFormatter(jsoniter.ConfigCompatibleWithStandardLibrary),
		WithLogLevel(int(cfg.LogLevel)),
		WithStdout(false),
		WithOutputPath(cfg.LogPath, "application.log"),
		WithHook(hooks.NewLogIDHook("context", "logId", 16, "未能获取logId")),
		WithHook(hooks.NewAutomaticCallerHook("context", "line", 7, 11, []string{"gclog/deprecated_api.go", "gclog/logger.go", "gclog/default_logger.go", "gclog/gorm_logger.go", "logrus"})),
		WithHook(hooks.NewConstFieldsHook("context").Add("app", appName).Add("hostname", mayHostname())),
		WithHook(hooks.NewDataHook("data")),
		WithHook(hooks.MustNewSentryHook(cfg.SentryDsn)),
	}

	logger, err := New(optionFuncs...)
	if err != nil {
		panic(err)
	}
	defaultLogger = logger
}
