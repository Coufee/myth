package logc

import (
	"context"
	"fmt"
)

// Info logs a message at the info log level.
func Debug(format string, args ...interface{}) {
	_Hander.Log(context.Background(), _LogLevelDebug, KVString(_log, fmt.Sprintf(format, args...)))
}

// Info logs a message at the info log level.
func Info(format string, args ...interface{}) {
	_Hander.Log(context.Background(), _LogLevelInfo, KVString(_log, fmt.Sprintf(format, args...)))
}

// Warn logs a message at the warning log level.
func Warn(format string, args ...interface{}) {
	_Hander.Log(context.Background(), _LogLevelWarn, KVString(_log, fmt.Sprintf(format, args...)))
}

// Error logs a message at the error log level.
func Error(format string, args ...interface{}) {
	_Hander.Log(context.Background(), _LogLevelError, KVString(_log, fmt.Sprintf(format, args...)))
}

// Error logs a message at the error log level.
func Emerg(format string, args ...interface{}) {
	_Hander.Log(context.Background(), _LogLevelEmerg, KVString(_log, fmt.Sprintf(format, args...)))
}

// Infoc logs a message at the info log level.
func Infoc(ctx context.Context, format string, args ...interface{}) {
	_Hander.Log(ctx, _LogLevelInfo, KVString(_log, fmt.Sprintf(format, args...)))
}

// Warnc logs a message at the warning log level.
func Warnc(ctx context.Context, format string, args ...interface{}) {
	_Hander.Log(ctx, _LogLevelWarn, KVString(_log, fmt.Sprintf(format, args...)))
}

// Errorc logs a message at the error log level.
func Errorc(ctx context.Context, format string, args ...interface{}) {
	_Hander.Log(ctx, _LogLevelError, KVString(_log, fmt.Sprintf(format, args...)))
}

// Infov logs a message at the info log level.
func Infov(ctx context.Context, args ...Field) {
	_Hander.Log(ctx, _LogLevelInfo, args...)
}

// Warnv logs a message at the warning log level.
func Warnv(ctx context.Context, args ...Field) {
	_Hander.Log(ctx, _LogLevelWarn, args...)
}

// Errorv logs a message at the error log level.
func Errorv(ctx context.Context, args ...Field) {
	_Hander.Log(ctx, _LogLevelError, args...)
}

// Infow logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Infow(ctx context.Context, args ...interface{}) {
	_Hander.Log(ctx, _LogLevelInfo, logw(args)...)
}

// Warnw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Warnw(ctx context.Context, args ...interface{}) {
	_Hander.Log(ctx, _LogLevelWarn, logw(args)...)
}

// Errorw logs a message with some additional context. The variadic key-value pairs are treated as they are in With.
func Errorw(ctx context.Context, args ...interface{}) {
	_Hander.Log(ctx, _LogLevelError, logw(args)...)
}
