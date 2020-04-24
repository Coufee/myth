package logc

import (
	"context"
	"time"

	pkgerr "github.com/pkg/errors"
)

type Handler interface {
	// Log handle log
	// variadic D is k-v struct represent log content
	Log(context.Context, LogLevel, ...Field)

	SetFormat(string)

	Close() error
}

func newHandlers(filters []string, handlers ...Handler) *Handlers {
	set := make(map[string]struct{})
	for _, k := range filters {
		set[k] = struct{}{}
	}
	return &Handlers{filters: set, handlers: handlers}
}

// Handlers a bundle for hander with filter function.
type Handlers struct {
	filters  map[string]struct{}
	handlers []Handler
}

func (hs Handlers) Log(ctx context.Context, lv LogLevel, args ...Field) {
	hasSource := false
	for i := range args {
		if _, ok := hs.filters[args[i].Key]; ok {
			args[i].Value = "***"
		}
		if args[i].Key == _source {
			hasSource = true
		}
	}
	if !hasSource {
		fn := getLineInfo(3)
		//errIncr(lv, fn)
		args = append(args, KVString(_source, fn))
	}
	args = append(args, KV(_time, time.Now()), KVInt64(_levelValue, int64(lv)), KVString(_level, lv.String()))
	for _, h := range hs.handlers {
		h.Log(ctx, lv, args...)
	}
}

func (hs Handlers) Close() (err error) {
	for _, h := range hs.handlers {
		if e := h.Close(); e != nil {
			err = pkgerr.WithStack(e)
		}
	}
	return
}

func (hs Handlers) SetFormat(format string) {
	for _, h := range hs.handlers {
		h.SetFormat(format)
	}
}
