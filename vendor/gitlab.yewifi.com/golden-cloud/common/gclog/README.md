<!--
 * @Date: 2020-11-02 16:11:24
 * @LastEditors: aiden.deng (Zhenpeng Deng)
 * @LastEditTime: 2020-11-07 17:46:52
-->


# gclog: Golden Cloud LOGger

## 1. 快速开始


把logger库改为gclog库即可使用兼容、性能更快、功能更丰富的的logger实例了：


```go
package main

import (
	"gitlab.yewifi.com/golden-cloud/common"
	"gitlab.yewifi.com/golden-cloud/common/gclog"
	"gitlab.yewifi.com/golden-cloud/common/logger"
)

func main() {
	loggerConf := &logger.LoggerConf{LogLevel: 10, LogPath: "./logs", SentryDsn: ""}

	// logger.SetUp(loggerConf, "my-app-name")
	// logger.Info(common.Uniqid(), "hello world")

	gclog.SetUp(loggerConf, "my-app-name")
	gclog.Info(common.Uniqid(), "hello world")
}
```

## 2. 介绍

之所以在已经有了common/logger库的前提下还重写一个logger库，是因为：
1. common/logger库代码结构较为混乱，功能揉杂，添加/修改功能时对后人的心理负担较大
2. 年久失修
3. 存在bug与性能问题（后面会说）
4. 已经使用了common/logger库的项目非常多，如果直接在common/logger库里重构，可能会造成预料之外的影响

综上，新建了一个common/gclog，特性：
1. 兼容性：这点必须排在第一，保持原有输出格式不变，同时也提供和common/logger一样的构造函数
2. 扩展性：尽可能利用logrus的特性，同时对修改关闭，对扩展开放
3. 可维护性：完善文档与注释，增加代码可读性，删除一些已经archived的第三方引用
4. 性能优化：更快的自定义json序列化方法、caller寻址方法等
5. bug修复：取logID时utf8隐藏bug、hook定位不明确等
6. 功能改造：一切均可自定义配置，包括序列化、immutable类字段、可变类字段、控制台输出等
7. 稳定性：该库已在风控项目生产环境中稳定运行了一段时间

具体地，做了以下事情：
1. 重写代码：保证不会再出现原库中的"logger.go文件中有一个全局log *logrus.Entry对象，log.go文件中有一个全局logger *logrus.Entry对象"的歧义现象。
2. 提供兼容构造方法`SetUp()`用于兼容老代码，提供`Default()`默认logger构造方法，以及提供给想要自定义扩展logger的底层构造方法`New()`。
3. 在原有JSONFormatter的基础上扩展, 结合`jsoniter`库，提供使用`jsoniter.ConfigCompatibleWithStandardLibrary`（默认，和`encoding/json`兼容）、`jsoniter.ConfigDefault`、`jsoniter.ConfigFastest`的不同初始化方式。
4. 优化caller寻址方法，不需要迭代50层函数栈去寻找caller。
5. logID解析方法更高效，且修复了原库中当message包含中文时直接对`string`取下标可能发生的bug，具体地是将message转成`rune`类型再解析logID。
6. 提供支持gorm2.0的logger，用于将SQL以log的方式打印出来，并检测函数执行时间，当超时时会以warning的log leve打log。配合上sentry可以方便的排查线上系统因SQL操作导致的性能瓶颈。
7. 优化sentry hook，去掉对maven库的引用，该库已被sentry官方摒弃。
8. 增加data hook，开发者可以通过logrus库的原生接口`logger.WithField("my-key", value)`将log数据保存在log json的data字段中。
9. 一切皆可配置，所有的logger默认行为，都是以**默认配置项**的方式去配置，而不是以硬编码的方式写死，最大化对修改关闭，对扩展开放。

## 3. 使用

### 3.1. 构造方法

gclog提供`SetUp()`方法来兼容common/logger的logger实例构造方法：
```go
func SetUp(cfg *commonlogger.LoggerConf, appName string)
```

在实际情况中，更推荐使用下列可配置性更好的构造方法：
```go
// 推荐：最底层的构造方法，从头开始根据给定的配置项去配置一个logger实例
func New(optionFuncs ...OptionFunc) (*logrus.Logger, error) 

// 推荐：最顶层的构造方法，提供一个开箱即用的logger实例
func Default(appName string, logDirPath string, logLevel int, sentryDsn string) *logrus.Logger 
```

### 3.2. 兼容性

gclog库与common一样，提供`SetUp(*logger.LoggerConf)`接口来初始化全局logger。


common/logger库：
```go
// File: t0_gclog/main.go
package main

import (
	"gitlab.yewifi.com/golden-cloud/common"
	"gitlab.yewifi.com/golden-cloud/common/gclog"
	"gitlab.yewifi.com/golden-cloud/common/logger"
)

func main() {
	loggerConf := &logger.LoggerConf{LogLevel: 10, LogPath: "./logs", SentryDsn: ""}
	logger.SetUp(loggerConf, "my-app-name")
	logger.Info(common.Uniqid(), "hello world")
}
```

输出（手动格式化后）：
```go
// 文件路径: ./logs/application.log 
{
  "context": {
    "app": "my-app-name",
    "line": "t0_gclog/main.go:16",
    "logId": "5f9fd84b00073af3"
  },
  "level": "info",
  "logTime": "2020-11-02 17:58:35",
  "message": "hello world"
}
```


gclog库：
```go
// File: t0_gclog/main.go
package main

import (
	"gitlab.yewifi.com/golden-cloud/common"
	"gitlab.yewifi.com/golden-cloud/common/gclog"
	"gitlab.yewifi.com/golden-cloud/common/logger"
)

func main() {
	loggerConf := &logger.LoggerConf{LogLevel: 10, LogPath: "./logs", SentryDsn: ""}
	gclog.SetUp(loggerConf, "my-app-name")
	gclog.Info(common.Uniqid(), "hello world")
}
```

输出（手动格式化后）：
```go
// 文件路径: ./logs/application.log 
{
  "context": {
    "app": "my-app-name",
    "line": "t0_gclog/main.go:17",
    "logId": "5f9fd85200094364"
  },
  "level": "info",
  "message": "hello world",
  "time": "2020-11-02 17:58:42"
}
```

可以看到，log的格式、日志文件名和构造函数`SetUp`，两者都是相同的。

### 3.2. New()




### 3.1. 默认logger

gclog库自带一个不需要初始化的全局logger，只输出到控制台：
```go
package main

import "gitlab.yewifi.com/golden-cloud/common/gclog"

func main() {
	gclog.Info("hello world")
}
```

输出：
```json
{"context":{"app":"gclog","caller":"t30_gclog/main.go:11","hostname":"red-coast-base"},"level":"info","message":"hello world","time":"2020-11-02 17:21:34"}
```

格式化看一下：
```json
{
  "context": {
    "app": "gclog", 
    "caller": "t30_gclog/main.go:11",
    "hostname": "red-coast-base"
  },
  "level": "info",
  "message": "hello world",
  "time": "2020-11-02 17:21:34"
}
```




### 默认log格式

```go
package main

import "gitlab.yewifi.com/golden-cloud/common/gclog"

func main() {
	gclog.Info("hello world")
}
```

输出：
```json
{"context":{"app":"gclog","caller":"t30_gclog/main.go:11","hostname":"red-coast-base"},"level":"info","message":"hello world","time":"2020-11-02 17:16:42"}
```

格式化看一下：
```json
{
  "context": {
    "app": "gclog",
    "caller": "t30_gclog/main.go:11",
    "hostname": "red-coast-base"
  },
  "level": "info",
  "message": "hello world",
  "time": "2020-11-02 17:17:11"
}
```

注意点：
- logId -> logID
- context新增hostname字段，用于k8s副本同时写log时区分是哪个pod写入的。



## 开发

提交MR并[@aiden.deng](https://gitlab.yewifi.com/aiden.deng)