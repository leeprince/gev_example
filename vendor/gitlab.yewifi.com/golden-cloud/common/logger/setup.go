package logger

import (
	"gitlab.yewifi.com/golden-cloud/common/logger/hooks"
)

type LoggerConf struct {
	LogLevel  uint32 `yaml:"log_level"`
	LogPath   string `yaml:"log_path"`
	SentryDsn string `yaml:"sentry_dsn"`
}

func SetUp(cnf *LoggerConf, appName string) {
	SetLogLevel(cnf.LogLevel)
	if cnf.SentryDsn != "" {
		hooks.SetEnableSentry()
		hooks.SetSentryDSN(cnf.SentryDsn)
	}
	SetAppName(appName)
	SetLogPath(cnf.LogPath)
	InitLog()
}
