# logger



## Usage

```go

package main

import (
	"io"

	"gitlab.yewifi.com/golden-cloud/common/logger"
)

func main() {
    
	logger.SetLogLevel(5) // logrus.DebugLevel
	logger.SetAppName("ticket-ocr/ocr-api")
	logger.SetLogPath("./")
	logger.AddStdoutWriter()
    
    logger.InitLog()

	logger.Info("hello world")
}

```