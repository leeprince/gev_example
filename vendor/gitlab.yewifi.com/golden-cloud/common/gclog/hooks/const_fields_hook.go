package hooks

/*
 * @Date: 2020-10-29 15:03:21
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 15:05:52
 */

import (
	"fmt"

	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/maputil"
)

var _ logrus.Hook = (*ConstFieldsHook)(nil)

type ConstFieldsHook struct {
	dataKey     string
	constFields map[string]interface{}
	levels      []logrus.Level
}

func NewConstFieldsHook(dataKey string, levels ...logrus.Level) *ConstFieldsHook {
	if levels == nil {
		levels = logrus.AllLevels
	}
	return &ConstFieldsHook{
		dataKey:     dataKey,
		constFields: make(map[string]interface{}),
		levels:      levels,
	}
}

func (h *ConstFieldsHook) Levels() []logrus.Level {
	return h.levels
}

func (h *ConstFieldsHook) Add(key string, value interface{}) *ConstFieldsHook {
	h.constFields[key] = value
	return h
}

func (h *ConstFieldsHook) Fire(entry *logrus.Entry) error {
	obj, err := maputil.GetOrCreatePath(entry.Data, h.dataKey)
	if err != nil {
		fmt.Printf("[ERROR] gclog/hooks/const_fields_hook.go: GetOrCreatePath failed, dataKey = %v \n", h.dataKey)
		return err
	}
	for key, value := range h.constFields {
		obj[key] = value
	}
	return nil
}
