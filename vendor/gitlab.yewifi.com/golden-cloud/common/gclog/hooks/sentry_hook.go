package hooks

/*
 * @Date: 2020-10-29 14:58:35
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 15:27:16
 */

import (
	"errors"
	"fmt"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/sirupsen/logrus"
)

var _ logrus.Hook = (*SentryHook)(nil)

var defaultSentryHookLevels = []logrus.Level{
	logrus.PanicLevel,
	logrus.FatalLevel,
	logrus.ErrorLevel,
	logrus.WarnLevel,
}

type SentryHook struct {
	// DSN tells the SDK where to send the events to.
	sentryCli *sentry.Client
	levels    []logrus.Level
}

func NewSentryHook(dsn string, levels ...logrus.Level) (*SentryHook, error) {
	if levels == nil {
		levels = defaultSentryHookLevels
	}

	if dsn == "" {
		hook := &SentryHook{
			sentryCli: nil,
			levels:    levels,
		}
		return hook, nil
	}

	cli, err := sentry.NewClient(sentry.ClientOptions{Dsn: dsn})
	if err != nil {
		return nil, err
	}

	hook := &SentryHook{
		sentryCli: cli,
		levels:    levels,
	}
	return hook, nil
}

func MustNewSentryHook(dsn string, levels ...logrus.Level) *SentryHook {
	h, err := NewSentryHook(dsn, levels...)
	if err != nil {
		panic(err)
	}
	return h
}

func (h *SentryHook) Levels() []logrus.Level {
	return h.levels
}

func (h *SentryHook) Fire(entry *logrus.Entry) error {
	if h.sentryCli == nil {
		return nil
	}

	// 如果发生报错，需要立马flush上报到sentry平台，且上报的耗时最多100ms，不能影响正常业务
	defer h.sentryCli.Flush(time.Millisecond * 100)

	scope := sentry.NewScope()
	for key, value := range entry.Data {
		switch typedValue := value.(type) {
		case map[string]interface{}:
			for subKey, subValue := range typedValue {
				scope.SetTag(fmt.Sprintf("%v.%v", key, subKey), fmt.Sprintf("%+v", subValue))
			}
		default:
			scope.SetTag(key, fmt.Sprintf("%+v", value))
		}

	}

	h.sentryCli.CaptureException(errors.New(entry.Message), nil, scope)
	return nil
}
