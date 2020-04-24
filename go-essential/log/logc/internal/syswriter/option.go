package syswriter

import (
	"log/syslog"
	"time"
)

var defaultOption = option{
	ServiceName: "default",
	ChanSize:    1024 * 8,
	Priority:    syslog.LOG_INFO,
	Batch:       false,
}

type option struct {
	Batch       bool
	ServiceName string
	ChanSize    int
	Priority    syslog.Priority

	// TODO export Option
	WriteTimeout time.Duration
}

// Option syswriter option
type Option func(opt *option)

// writer log level
func Priority(p syslog.Priority) Option {
	return func(opt *option) {
		opt.Priority = p
	}
}
