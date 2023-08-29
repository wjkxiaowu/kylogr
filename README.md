# kylogr

kylogr 是用来方便做日志切割和回收的，主要基于file-rotatelogs对logrus做日志切分和回收。

## 用法

把包导入

```go
import(
    _ "github.com/wjkxiaowu/kylogr"
)
```



默认设置,

```golang
LOG_DIR            = "./logs"  //日志目录
LOG_NAME_PREFIX    = "log"     //日志文件名的前缀（最好以服务名称为前缀）
LOG_LEVEL          = "info"    //日志级别
LOG_NAME_SUFFIX    = "_%Y_%m_%d_%H_%M_%S.log" //日志文件名后缀和文件名字前缀合成完整文件名
LOG_ROTATION_COUNT = "7"   // 设置文件清理前最多保存的个数，与LOG_MAX_AGE_HOUR互斥
LOG_ROTATION_TIME  = "24" //小时 设置日志分割的时间,这里设置为一天分割一次, 设置日志切割时间间隔(1天）
LOG_MAX_AGE_HOUR   = "null" //小时 设置文件清理前的最长保存时间 与DEFAULT_LOG_ROTATION_COUNT互斥
LOG_FORMATTER      = "text" //日志输出格式
```

配置方法

```bash
export LOG_LEVEL=debug
```
