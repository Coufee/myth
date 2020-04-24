package logc

import (
	"fmt"
	"io"
	"log/syslog"
	"sync"
)

var (
	_InitOnce    *sync.Once = &sync.Once{}
	_Hander      Handler
	_Filter      logFilter
	_Config      *Config
	_StdOut      = false     //默认输出
	_ServiceName = "default" //默认服务名称
	_Dir         string
	_NetWork     string
	_Addr        string
	_Os          string
)

// Config log config.
type Config struct {
	Family         string
	Host           string
	Stdout         bool
	NetWork        string
	Addr           string
	ServiceName    string
	Dir            string
	FileBufferSize int64 // buffer size
	MaxLogFile     int   // MaxLogFile
	RotateSize     int64 // RotateSize
	// Filter tell log handler which field are sensitive message, use * instead.
	Filter []string
}

//// metricErrCount prometheus error counter.
//var (
//	metricErrCount = metric.NewBusinessMetricCount("log_error_total", "source")
//)

const (
	_timeFormat = "2006-01-02T15:04:05.999999"

	// log level defined in level.go.
	_levelValue = "level_value"
	//  log level name: INFO, WARN...
	_level = "level"
	// log time.
	_time = "time"
	// request path.
	// _title = "title"
	// log file.
	_source = "source"
	// common log filed.
	_log = "log"
	// app name.
	_appID = "app_id"
	// container ID.
	_instanceID = "instance_id"
	// uniq ID from trace.
	_tid = "traceid"
	// request time.
	// _ts = "ts"
	// requester.
	_caller = "caller"
	// container environment: prod, pre, uat, fat.
	_deplyEnv = "env"
	// container area.
	_zone = "zone"
	// mirror flag
	_mirror = "mirror"
	// color.
	_color = "color"
	// env_color
	_envColor = "env_color"
	// cluster.
	_cluster = "cluster"
)

const (
	osWindow = "windows"
	osMac    = "darwin"
	osLinux  = "linux"

	handlerTypeStdOut  = "StdOut"
	handlerTypeFileLog = "File"
	handlerTypeSysLog  = "SysLog"

	defaultSysPattern  = "[%D %T] [%L] [%S] %M"
	defaultStdPattern  = "%d %T %L %f - %M"
	defaultFilePattern = "[%D %T] [%L] [%S] %M"
)

type LogLevel int

const (
	_LogLevelDebug LogLevel = iota
	_LogLevelTrace
	_LogLevelInfo
	_LogLevelNotice
	_LogLevelWarn
	_LogLevelError
	_LogLevelFatal
	_LogLevelCrit
	_LogLevelAlert
	_LogLevelEmerg
)

var levelNames = [...]string{
	_LogLevelDebug:  "DEBUG",
	_LogLevelTrace:  "TRACE",
	_LogLevelInfo:   "INFO",
	_LogLevelNotice: "NOTICE",
	_LogLevelWarn:   "WARN",
	_LogLevelError:  "ERROR",
	_LogLevelFatal:  "FATAL",
	_LogLevelCrit:   "CRIT",
	_LogLevelAlert:  "ALERT",
	_LogLevelEmerg:  "EMERG",
}

var sysLevelPriority = [...]syslog.Priority{
	_LogLevelDebug:  syslog.LOG_DEBUG,
	_LogLevelInfo:   syslog.LOG_INFO,
	_LogLevelNotice: syslog.LOG_NOTICE,
	_LogLevelWarn:   syslog.LOG_WARNING,
	_LogLevelError:  syslog.LOG_ERR,
	_LogLevelCrit:   syslog.LOG_CRIT,
	_LogLevelAlert:  syslog.LOG_ALERT,
	_LogLevelEmerg:  syslog.LOG_EMERG,
}

func (l LogLevel) String() string {
	return levelNames[l]
}

type LogColor int

const (
	_Black LogColor = iota + 30
	_Red
	_Green
	_Yellow
	_Blue
	_Magenta
	_Cyan
	_White
)

var levelColor = [...]LogColor{
	_LogLevelDebug: _White,
	_LogLevelTrace: _Cyan,
	_LogLevelInfo:  _Magenta,
	_LogLevelWarn:  _Blue,
	_LogLevelError: _Yellow,
	_LogLevelFatal: _Green,
	_LogLevelEmerg: _Red,
}

func (c *LogColor) Add(s string) string {
	return fmt.Sprintf("\x1b[%dm%s\x1b[0m", uint8(*c), s)
}

// level idx
const (
	_FileInfoIdx = iota
	_FileWarnIdx
	_FileErrorIdx
	_FileTotalIdx

	_SysLogInfoIdx = iota
	_SysLogErrorIdx
	_SysLogEmergIdx
	_SysLogeTotalIdx
)

var _FileNames = map[int]string{
	_FileInfoIdx:  "info.log",
	_FileWarnIdx:  "warning.log",
	_FileErrorIdx: "error.log",
}

var _SysLogPriority = map[int]syslog.Priority{
	_SysLogInfoIdx:  syslog.LOG_INFO,
	_SysLogErrorIdx: syslog.LOG_ERR,
	_SysLogEmergIdx: syslog.LOG_EMERG,
}

// Render render log output
type Render interface {
	Render(io.Writer, LogLevel, map[string]interface{}) error
	RenderString(map[string]interface{}) string
}
