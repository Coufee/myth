package logc

import (
	"flag"
	"os"
	"runtime"

	"myth/go-essential/conf/env"
)

func init() {
	host, _ := os.Hostname()
	_Os = runtime.GOOS
	_Config = &Config{
		Family: env.AppID,
		Host:   host,
	}

	_Hander = newHandlers([]string{}, NewStdout())
	addFlag(flag.CommandLine)
}

func addFlag(args *flag.FlagSet) {
	args.BoolVar(&_StdOut, "log.stdout", _StdOut, "log enable stdout or not, or use LOG_STDOUT env variable.")
	args.StringVar(&_Dir, "log.dir", _Dir, "log file `path, or use LOG_DIR env variable.")
	args.StringVar(&_NetWork, "log.network", _NetWork, "log file `path, or use LOG_DIR env variable.")
	args.StringVar(&_Addr, "log.addr", _Addr, "log file `path, or use LOG_DIR env variable.")
	args.Var(&_Filter, "log.filter", "log field for sensitive message, or use LOG_FILTER env variable, format: field1,field2.")
}

// Init create logger with context.
func Init(conf *Config) {
	_InitOnce.Do(func() {
		var isNil bool
		if conf == nil {
			isNil = true
			conf = &Config{
				Stdout:      _StdOut,
				ServiceName: _ServiceName,
				Dir:         _Dir,
				Filter:      _Filter,
				NetWork:     _NetWork,
				Addr:        _Addr,
			}
		}
		if len(env.AppID) != 0 {
			conf.Family = env.AppID // for caster
		}
		conf.Host = env.Hostname
		if len(conf.Host) == 0 {
			host, _ := os.Hostname()
			conf.Host = host
		}
		var hs []Handler
		// when env is dev
		if conf.Stdout || (isNil && (env.DeployEnv == "" || env.DeployEnv == env.DeployEnvDev)) {
			hs = append(hs, NewStdout())
		}
		if conf.Dir != "" {
			hs = append(hs, NewFile(conf.Dir, conf.FileBufferSize, conf.RotateSize, conf.MaxLogFile))
		}
		if conf.ServiceName != "" {
			hs = append(hs, NewSysLog(conf.NetWork, conf.Addr, conf.ServiceName))
		}

		_Hander = newHandlers(conf.Filter, hs...)
		_Config = conf
	})
}

func logw(args []interface{}) []Field {
	if len(args)%2 != 0 {
		Warn("log: the variadic must be plural, the last one will ignored")
	}
	ds := make([]Field, 0, len(args)/2)
	for i := 0; i < len(args)-1; i = i + 2 {
		if key, ok := args[i].(string); ok {
			ds = append(ds, KV(key, args[i+1]))
		} else {
			Warn("log: key must be string, get %T, ignored", args[i])
		}
	}
	return ds
}

// SetFormat only effective on stdout and file handler
// %T time format at "15:04:05.999" on stdout handler, "15:04:05 MST" on file handler
// %t time format at "15:04:05" on stdout handler, "15:04" on file on file handler
// %D data format at "2006/01/02"
// %d data format at "01/02"
// %L log level e.g. INFO WARN ERROR
// %M log message and additional fields: key=value this is log message
// NOTE below pattern not support on file handler
// %f function name and line number e.g. model.Get:121
// %i instance id
// %e deploy env e.g. dev uat fat prod
// %z zone
// %S full file name and line number: /a/b/c/d.go:23
// %s final file name element and line number: d.go:23
func SetFormat(format string) {
	_Hander.SetFormat(format)
}

// Close close resource.
func Close() (err error) {
	err = _Hander.Close()
	_Hander = _defaultStdout
	return
}
