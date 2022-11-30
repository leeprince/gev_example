package hooks

import (
	"fmt"

	"github.com/sirupsen/logrus"
)

func NewFormatLogIdHook() *FormatLogIdHook {
	return &FormatLogIdHook{}
}

type FormatLogIdHook struct {
}

func (hook *FormatLogIdHook) Fire(entry *logrus.Entry) error {
	entry.Data["logId"] = "未获取到logId"
	if len(entry.Message) >= 16 {
		defer func() {
			if e := recover(); e != nil {
				entry.Message = "分割logId发生了错误:" + fmt.Sprintf("%s %s", e, entry.Message)
				entry.Data["logId"] = "获取logId错误"
			}
		}()
		logId := entry.Message[0:16]
		var e, d int
		for i := 0; i < len(logId); i++ {
			switch {
			case 96 < logId[i] && logId[i] < 123:
				e += 1
			case 47 < logId[i] && logId[i] < 58:
				d += 1
			}
		}
		// fmt.Println("entery message", entry.Message, " len:", len(entry.Message), " content:", entry.Message[16:28])
		if e+d == 16 {
			entry.Data["logId"] = logId
			entry.Message = entry.Message[16:len(entry.Message)]
		} else {
			entry.Data["logId"] = "nil"
		}
	}
	return nil
}

func (hook *FormatLogIdHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
