package formatter

import (
	"github.com/sirupsen/logrus"
)

// StatFormatter  日志格式: 时间|日志id|msg
type StatFormatter struct {
}

// Format 转换格式
func (f StatFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	// logId
	logId, ok := entry.Data["logId"].(string)
	if !ok || logId == "" {
		logId = "未获取到logId"
	}

	// 调用函数
	line, ok := entry.Data["line"].(string)
	if !ok || line == "" {
		line = "未获取到line"
	}

	// 时间
	timeNow := entry.Time.Format("2006-01-02 15:04:05")

	// 时间|日志id|msg
	result := timeNow + "|" + logId + "|" + entry.Message + "\n"

	return []byte(result), nil
}
