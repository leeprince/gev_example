package hooks

/*
 * @Date: 2020-10-29 13:41:10
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 14:49:59
 */

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/maputil"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/runtimeutil"
)

var _ logrus.Hook = (*AutomaticCallerHook)(nil)

type AutomaticCallerHook struct {
	dataKey          string
	fieldKey         string
	minSkip          int
	maxSkip          int
	excludedPrefixes []string
	levels           []logrus.Level
}

func NewAutomaticCallerHook(dataKey string, fieldKey string, minSkip int, maxSkip int, excludedPrefixes []string, levels ...logrus.Level) *AutomaticCallerHook {
	if levels == nil {
		levels = logrus.AllLevels
	}
	return &AutomaticCallerHook{
		dataKey:          dataKey,
		fieldKey:         fieldKey,
		minSkip:          minSkip,
		maxSkip:          maxSkip,
		excludedPrefixes: excludedPrefixes,
		levels:           levels,
	}
}

func (h *AutomaticCallerHook) Levels() []logrus.Level {
	return h.levels
}

func (h *AutomaticCallerHook) Fire(entry *logrus.Entry) error {
	obj, err := maputil.GetOrCreatePath(entry.Data, h.dataKey)
	if err != nil {
		fmt.Printf("[ERROR] gclog/hooks/automatic_caller_hook.go: GetOrCreatePath failed, dataKey = %v \n", h.dataKey)
		return err
	}

	var callerName string
	var lineNo int
	for i := h.minSkip; i <= h.maxSkip; i++ {
		callerName, lineNo = runtimeutil.GetCaller(i)
		hasNoExcludedPrefixes := true
		for _, excludedPrefix := range h.excludedPrefixes {
			if strings.HasPrefix(callerName, excludedPrefix) {
				hasNoExcludedPrefixes = false
				break
			}
		}
		if hasNoExcludedPrefixes {
			break
		}
	}

	obj[h.fieldKey] = fmt.Sprintf("%s:%d", callerName, lineNo)
	return nil
}
