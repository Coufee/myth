package logc

import (
	"context"
	"io"
	"log/syslog"
	"myth/go-essential/log/logc/internal/syswriter"
	"time"
)

type SysLogHandler struct {
	render      Render
	serviceName string
	sw          [_SysLogeTotalIdx]*syswriter.SysWriter
}

func NewSysLog(network, addr, serviceName string) *SysLogHandler {
	handler := &SysLogHandler{
		render:      newPatternRender(defaultSysPattern, handlerTypeSysLog),
		serviceName: serviceName,
	}

	writer, err := syslog.Dial(network, addr, syslog.LOG_LOCAL0, serviceName)
	if err != nil {
		panic(err)
	}

	for index, level := range _SysLogPriority {
		w, err := syswriter.New(writer, syswriter.Priority(level))
		if err != nil {
			panic(err)
		}

		handler.sw[index] = w
	}

	return handler
}

func (s *SysLogHandler) Log(ctx context.Context, lv LogLevel, args ...Field) {
	d := toMap(args...)
	// add extra fields
	addExtraField(ctx, d)
	d[_time] = time.Now().Format(_timeFormat)
	var w io.Writer
	switch lv {
	case _LogLevelEmerg:
		w = s.sw[_SysLogEmergIdx]
	case _LogLevelError:
		w = s.sw[_SysLogErrorIdx]
	default:
		w = s.sw[_SysLogInfoIdx]
	}

	s.render.Render(w, lv, d)
}

func (s *SysLogHandler) Close() error {
	for _, sw := range s.sw {
		// ignored error
		sw.Close()
	}

	return nil
}

func (s *SysLogHandler) SetFormat(format string) {
	s.render = newPatternRender(format, handlerTypeSysLog)
}
