package hooks

/*
 * @Date: 2020-09-28 14:42:27
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-10-29 15:09:19
 */

import (
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	"gitlab.yewifi.com/golden-cloud/common/gclog/internal/maputil"
)

var _ logrus.Hook = (*DataHook)(nil)

type DataHook struct {
	dataFieldName        string
	dataFieldTopItemName string
	levels               []logrus.Level
}

func NewDataHook(dataFieldName string, levels ...logrus.Level) *DataHook {
	if levels == nil {
		levels = logrus.AllLevels
	}
	return &DataHook{
		dataFieldName:        dataFieldName,
		dataFieldTopItemName: strings.Split(dataFieldName, ".")[0],
		levels:               levels,
	}
}

func (h *DataHook) Levels() []logrus.Level {
	return h.levels
}

func (h *DataHook) Fire(entry *logrus.Entry) error {

	for fieldName, fieldValue := range entry.Data {
		if fieldName == h.dataFieldTopItemName {
			continue
		}
		if fieldName == "context" || fieldName == "err" || fieldName == "error" {
			continue
		}

		dataObj, err := maputil.GetOrCreatePath(entry.Data, h.dataFieldName)
		if err != nil {
			fmt.Printf("[ERROR] gclog/hooks/data_hook.go: GetOrCreatePath failed, dataFieldName = %v \n", h.dataFieldName)
			return err
		}

		var obj map[string]interface{}
		var lastItemName string

		itemNames := strings.Split(fieldName, ".")
		nItems := len(itemNames)

		if nItems == 1 {
			lastItemName = itemNames[0]
			obj = dataObj
		} else {
			lastItemName = itemNames[nItems-1]
			fatherFieldName := strings.Join(itemNames[:nItems-1], ".")
			obj, err = maputil.GetOrCreatePath(dataObj, fatherFieldName)
		}
		if err != nil {
			return err
		}
		obj[lastItemName] = fieldValue
		delete(entry.Data, fieldName)
	}
	return nil
}
