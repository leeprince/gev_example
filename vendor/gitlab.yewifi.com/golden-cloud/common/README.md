<!--
 * @Date: 2021-08-18 10:46:18
 * @LastEditors: aiden.deng(Zhenpeng Deng)
 * @LastEditTime: 2021-08-18 14:32:48
-->


# common: 高灯云 golang 公共库


## 模块

- **adaptime**: TODO
- **apiauthorization**: TODO
- **bigdata**: TODO
- **billingcli**: TODO
- **bizerr**: TODO
- **byteutil**: TODO
- **constval2**: TODO
- **crypto**: TODO
- **error**: TODO
- **file**: TODO
- **gcjsonpb**: TODO
- **gclog**: Golden Cloud LOGger, 日志库
- **graceful**: TODO
- **grpcmw**: TODO
- **hashid**: TODO
- **httpcli**: TODO
- **logger**: TODO
- **mathutil**: TODO
- **mysqlhelper2**: TODO
- **rabbit**: TODO
- **redishelper**: TODO
- **routine**: TODO
- **sign**: TODO
- **slices**: TODO
- **tracer**: jeager opentracing 调用链
- **txobjstorage**: TODO
- **util**: TODO
- **yamlutil**: TODO
- **common.go**: TODO
- **date_time.go**: TODO
- **err_filter.go**: TODO
- **error.go**: TODO
- **event_id.go**: 开放平台回调事件ID
- **marshal.go**: TODO
- **parse.go**: TODO
- **recovery.go**: TODO


### 模块:adaptime
### 模块:apiauthorization
### 模块:bigdata
### 模块:billingcli
### 模块:bizerr
### 模块:byteutil
### 模块:constval2
### 模块:crypto
### 模块:error
### 模块:file
### 模块:fuck.txt
### 模块:gcjsonpb
### 模块:gclog
### 模块:graceful
### 模块:grpcmw
### 模块:hashid
### 模块:httpcli
### 模块:logger
### 模块:mathutil
### 模块:mysqlhelper2
### 模块:rabbit
### 模块:redishelper
### 模块:routine
### 模块:sign
### 模块:slices
### 模块:tracer
### 模块:txobjstorage
### 模块:util
### 模块:yamlutil
### 模块:common.go
### 模块:date_time.go
### 模块:err_filter.go
### 模块:error.go

### 模块:event_id.go

用于生成开放平台回调事件ID `event_id`。`event_id` 默认长度固定为24位字符，生成时可指定前缀(比如logID)。

用法:

```go

import (
    "gitlab.yewifi.com/golden-cloud/common"
    "fmt"
)

func main() {
    eventID := common.EventID()
    fmt.Println(eventID)        
    // print: "52fdfc072182654f163f5f0f"，长度为24个字符

    logID := common.Uniqid()  // 611ca9680005ff48，16个字符
    eventID = common.EventID(logID)
    fmt.Println(logID)          
    // print: "611ca9680005ff489a621d72"，长度仍然为24个字符，前16个为log_id
}

```

效果:
```bash
go test -v -count=1 ./event_id_test.go
```
- 1协程，生成10000个 event id，无重复。
- 16协程，每个协程生成10w个 event id，并发执行，无重复。



性能:
```bash
go test -v -bench=. -benchmem event_id_test.go
```
```
goos: darwin
goarch: amd64
cpu: Intel(R) Core(TM) i9-9880H CPU @ 2.30GHz
BenchmarkEventID
BenchmarkEventID-16      6233388               193.8 ns/op            64 B/op          3 allocs/op
PASS
ok      command-line-arguments  1.924s
```


### 模块:marshal.go
### 模块:parse.go
### 模块:recovery.go