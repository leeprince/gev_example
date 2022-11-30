package hooks

/*
 * @Date: 2020-10-29 11:57:07
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 14:32:47
 */

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/maputil"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/runtimeutil"
)

var _ logrus.Hook = (*CallerHook)(nil)

type CallerHook struct {
	dataKey  string
	fieldKey string
	skip     int
	levels   []logrus.Level
}

func NewCallerHook(dataKey string, fieldKey string, skip int, levels ...logrus.Level) *CallerHook {
	if levels == nil {
		levels = logrus.AllLevels
	}
	return &CallerHook{
		dataKey:  dataKey,
		fieldKey: fieldKey,
		skip:     skip,
		levels:   levels,
	}
}

func (h *CallerHook) Levels() []logrus.Level {
	return h.levels
}

func (h *CallerHook) Fire(entry *logrus.Entry) error {
	obj, err := maputil.GetOrCreatePath(entry.Data, h.dataKey)
	if err != nil {
		fmt.Printf("[ERROR] gclog/caller_hook.go: GetOrCreatePath failed, dataKey = %v \n", h.dataKey)
		return err
	}

	callerName, lineNo := runtimeutil.GetCaller(h.skip)
	obj[h.fieldKey] = fmt.Sprintf("%s:%d", callerName, lineNo)
	return nil
}
