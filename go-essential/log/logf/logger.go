package log

import (
	"context"
	"github.com/sirupsen/logrus"
	"io"
	"sync"
	"time"
)

func New(stdLogger *logrus.Logger) *Logger {
	result := &Logger{
		StdLogger:   stdLogger,
		DebugLogger: logrus.New(),
	}

	result.DebugLogger.SetLevel(logrus.TraceLevel)
	return result
}

type Logger struct {
	StdLogger   *logrus.Logger
	DebugLogger *logrus.Logger

	entryPool sync.Pool
}

func (logger *Logger) newEntry() *Entry {
	return NewEntry(logger)
}

// Adds a field to the log entry, note that it doesn't log until you call
// Debug, Print, Info, Warn, Error, Fatal or Panic. It only creates a log entry.
// If you want multiple fields, use `WithFields`.
func (logger *Logger) WithField(key string, value interface{}) *Entry {
	return logger.newEntry().WithField(key, value)
}

// Adds a struct of fields to the log entry. All it does is call `WithField` for
// each `Field`.
func (logger *Logger) WithFields(fields Fields) *Entry {
	return logger.newEntry().WithFields(fields)
}

// Add an error as single field to the log entry.  All it does is call
// `WithError` for the given `error`.
func (logger *Logger) WithError(err error) *Entry {
	return logger.newEntry().WithError(err)
}

// Add a context to the log entry.
func (logger *Logger) WithContext(ctx context.Context) *Entry {
	return logger.newEntry().WithContext(ctx)
}

// Overrides the time of the log entry.
func (logger *Logger) WithTime(t time.Time) *Entry {
	return logger.newEntry().WithTime(t)
}

// 当无Context直接调用时，等于直接调用 StdLogger
func (logger *Logger) Log(level Level, args ...interface{}) {
	logrusLevel := logrus.Level(level)

	logger.StdLogger.Log(logrusLevel, args...)
}

func (logger *Logger) Logf(level Level, format string, args ...interface{}) {
	logrusLevel := logrus.Level(level)

	logger.StdLogger.Logf(logrusLevel, format, args...)
}

func (logger *Logger) Logln(level Level, args ...interface{}) {
	logrusLevel := logrus.Level(level)

	logger.StdLogger.Logln(logrusLevel, args...)
}

func (logger *Logger) Trace(args ...interface{}) {
	logger.Log(TraceLevel, args...)
}

func (logger *Logger) Debug(args ...interface{}) {
	logger.Log(DebugLevel, args...)
}

func (logger *Logger) Info(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *Logger) Print(args ...interface{}) {
	logger.Log(InfoLevel, args...)
}

func (logger *Logger) Warn(args ...interface{}) {
	logger.Log(WarnLevel, args...)
}

func (logger *Logger) Warning(args ...interface{}) {
	logger.Warn(args...)
}

func (logger *Logger) Error(args ...interface{}) {
	logger.Log(ErrorLevel, args...)
}

func (logger *Logger) Fatal(args ...interface{}) {
	logger.Log(FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panic(args ...interface{}) {
	logger.Log(PanicLevel, args...)
}

func (logger *Logger) Tracef(format string, args ...interface{}) {
	logger.Logf(TraceLevel, format, args...)
}

func (logger *Logger) Debugf(format string, args ...interface{}) {
	logger.Logf(DebugLevel, format, args...)
}

func (logger *Logger) Infof(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

func (logger *Logger) Printf(format string, args ...interface{}) {
	logger.Logf(InfoLevel, format, args...)
}

func (logger *Logger) Warnf(format string, args ...interface{}) {
	logger.Logf(WarnLevel, format, args...)
}

func (logger *Logger) Warningf(format string, args ...interface{}) {
	logger.Warnf(format, args...)
}

func (logger *Logger) Errorf(format string, args ...interface{}) {
	logger.Logf(ErrorLevel, format, args...)
}

func (logger *Logger) Fatalf(format string, args ...interface{}) {
	logger.Logf(FatalLevel, format, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicf(format string, args ...interface{}) {
	logger.Logf(PanicLevel, format, args...)
}

func (logger *Logger) Traceln(args ...interface{}) {
	logger.Logln(TraceLevel, args...)
}

func (logger *Logger) Debugln(args ...interface{}) {
	logger.Logln(DebugLevel, args...)
}

func (logger *Logger) Infoln(args ...interface{}) {
	logger.Logln(InfoLevel, args...)
}

func (logger *Logger) Println(args ...interface{}) {
	logger.Logln(InfoLevel, args...)
}

func (logger *Logger) Warnln(args ...interface{}) {
	logger.Logln(WarnLevel, args...)
}

func (logger *Logger) Warningln(args ...interface{}) {
	logger.Warnln(args...)
}

func (logger *Logger) Errorln(args ...interface{}) {
	logger.Logln(ErrorLevel, args...)
}

func (logger *Logger) Fatalln(args ...interface{}) {
	logger.Logln(FatalLevel, args...)
	logger.Exit(1)
}

func (logger *Logger) Panicln(args ...interface{}) {
	logger.Logln(PanicLevel, args...)
}

func (logger *Logger) Exit(code int) {
	logger.StdLogger.Exit(code)
	logger.DebugLogger.Exit(code)
}

//When file is opened with appending mode, it's safe to
//write concurrently to a file (within 4k message on Linux).
//In these cases user can choose to disable the lock.
func (logger *Logger) SetNoLock() {
	logger.StdLogger.SetNoLock()
	logger.DebugLogger.SetNoLock()
}

// SetLevel sets the logger level.
func (logger *Logger) SetLevel(level Level) {
	logger.StdLogger.SetLevel(logrus.Level(level))
}

// GetLevel returns the logger level.
func (logger *Logger) GetLevel() Level {
	return Level(logger.StdLogger.GetLevel())
}

// AddHook adds a hook to the logger hooks.
func (logger *Logger) AddHook(hook logrus.Hook) {
	logger.StdLogger.AddHook(hook)
	logger.DebugLogger.AddHook(hook)
}

// IsLevelEnabled checks if the log level of the logger is greater than the level param
func (logger *Logger) IsLevelEnabled(level Level) bool {
	return true
}

// SetFormatter sets the logger formatter.
func (logger *Logger) SetFormatter(formatter logrus.Formatter) {
	logger.StdLogger.SetFormatter(formatter)
	logger.DebugLogger.SetFormatter(formatter)
}

// SetOutput sets the logger output.
func (logger *Logger) SetOutput(output io.Writer) {
	logger.StdLogger.SetOutput(output)
	logger.DebugLogger.SetOutput(output)
}

func (logger *Logger) SetReportCaller(reportCaller bool) {
	logger.StdLogger.SetReportCaller(reportCaller)
	logger.DebugLogger.SetReportCaller(reportCaller)
}

// ReplaceHooks replaces the logger hooks and returns the old ones
func (logger *Logger) ReplaceHooks(hooks logrus.LevelHooks) logrus.LevelHooks {
	logger.StdLogger.ReplaceHooks(hooks)
	return logger.DebugLogger.ReplaceHooks(hooks)
}
