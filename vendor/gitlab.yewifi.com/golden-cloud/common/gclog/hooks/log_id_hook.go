package hooks

/*
 * @Date: 2020-10-29 14:50:20
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 14:56:55
 */

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/maputil"
)

var _ logrus.Hook = (*LogIDHook)(nil)

type LogIDHook struct {
	dataKey      string
	fieldName    string
	logIDLength  int
	defaultLogID string
	levels       []logrus.Level
}

func NewLogIDHook(dataKey string, fieldName string, logIDLength int, defaultLogID string, levels ...logrus.Level) *LogIDHook {
	if levels == nil {
		levels = logrus.AllLevels
	}
	return &LogIDHook{
		dataKey:      dataKey,
		fieldName:    fieldName,
		logIDLength:  logIDLength,
		defaultLogID: defaultLogID,
		levels:       levels,
	}
}

func (h *LogIDHook) Levels() []logrus.Level {
	return h.levels
}

func (h *LogIDHook) Fire(entry *logrus.Entry) error {
	obj, err := maputil.GetOrCreatePath(entry.Data, h.dataKey)
	if err != nil {
		fmt.Printf("[ERROR] gclog/hooks/log_id_hook.go: GetOrCreatePath failed, dataKey = %v \n", h.dataKey)
		return err
	}

	obj[h.fieldName] = h.defaultLogID

	messageRunes := []rune(entry.Message)
	if len(messageRunes) >= h.logIDLength {
		for i := 0; i < h.logIDLength; i++ {
			if !((97 <= messageRunes[i] && messageRunes[i] <= 122) || (48 <= messageRunes[i] && messageRunes[i] <= 57)) {
				return nil
			}
		}

		obj[h.fieldName] = string(messageRunes[:h.logIDLength])
		entry.Message = string(messageRunes[h.logIDLength:])
	}
	return nil
}
