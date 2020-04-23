package log

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"time"
)

type Entry struct {
	stdEntry   *logrus.Entry
	debugEntry *logrus.Entry
}

func NewEntry(logger *Logger) *Entry {
	return &Entry{
		stdEntry:   logger.StdLogger.WithContext(nil),
		debugEntry: logger.DebugLogger.WithContext(nil),
	}
}

// Returns the string representation from the reader and ultimately the
// formatter.
func (entry *Entry) String() (string, error) {
	return entry.stdEntry.String()
}

// Add an error as single field (using the key defined in ErrorKey) to the Entry.
func (entry *Entry) WithError(err error) *Entry {
	entry.stdEntry = entry.stdEntry.WithError(err)
	entry.debugEntry = entry.debugEntry.WithError(err)
	return entry
}

// Add a context to the Entry.
func (entry *Entry) WithContext(ctx context.Context) *Entry {
	entry.stdEntry = entry.stdEntry.WithContext(ctx)
	entry.debugEntry = entry.debugEntry.WithContext(ctx)

	return entry
}

// Add a single field to the Entry.
func (entry *Entry) WithField(key string, value interface{}) *Entry {
	entry.stdEntry = entry.stdEntry.WithField(key, value)
	entry.debugEntry = entry.debugEntry.WithField(key, value)

	return entry
}

func (entry *Entry) WithFields(fields Fields) *Entry {
	entry.stdEntry = entry.stdEntry.WithFields(logrus.Fields(fields))
	entry.debugEntry = entry.debugEntry.WithFields(logrus.Fields(fields))

	return entry
}

// Overrides the time of the Entry.
func (entry *Entry) WithTime(t time.Time) *Entry {
	entry.stdEntry = entry.stdEntry.WithTime(t)
	entry.debugEntry = entry.debugEntry.WithTime(t)

	return entry
}

func (entry Entry) HasCaller() (has bool) {
	return entry.stdEntry.HasCaller()
}

func (entry *Entry) Log(level Level, args ...interface{}) {
	logrusLevel := logrus.Level(level)
	entry.stdEntry.Log(logrusLevel, args...)

	if entry.stdEntry.Logger.IsLevelEnabled(logrusLevel) == false && entry.IsDebugMode() {
		entry.debugEntry.Log(logrusLevel, args...)
	}
}

func (entry *Entry) Trace(args ...interface{}) {
	entry.Log(TraceLevel, args...)
}

func (entry *Entry) Debug(args ...interface{}) {
	entry.Log(DebugLevel, args...)
}

func (entry *Entry) Print(args ...interface{}) {
	entry.Info(args...)
}

func (entry *Entry) Info(args ...interface{}) {
	entry.Log(InfoLevel, args...)
}

func (entry *Entry) Warn(args ...interface{}) {
	entry.Log(WarnLevel, args...)
}

func (entry *Entry) Warning(args ...interface{}) {
	entry.Warn(args...)
}

func (entry *Entry) Error(args ...interface{}) {
	entry.Log(ErrorLevel, args...)
}

func (entry *Entry) Fatal(args ...interface{}) {
	entry.Log(FatalLevel, args...)
}

func (entry *Entry) Panic(args ...interface{}) {
	entry.Log(PanicLevel, args...)
	panic(fmt.Sprint(args...))
}

// Entry Printf family functions

func (entry *Entry) Logf(level Level, format string, args ...interface{}) {
	logrusLevel := logrus.Level(level)
	entry.stdEntry.Logf(logrusLevel, format, args...)

	if entry.stdEntry.Logger.IsLevelEnabled(logrusLevel) == false && entry.IsDebugMode() {
		entry.debugEntry.Logf(logrusLevel, format, args...)
	}
}

func (entry *Entry) Tracef(format string, args ...interface{}) {
	entry.Logf(TraceLevel, format, args...)
}

func (entry *Entry) Debugf(format string, args ...interface{}) {
	entry.Logf(DebugLevel, format, args...)
}

func (entry *Entry) Infof(format string, args ...interface{}) {
	entry.Logf(InfoLevel, format, args...)
}

func (entry *Entry) Printf(format string, args ...interface{}) {
	entry.Infof(format, args...)
}

func (entry *Entry) Warnf(format string, args ...interface{}) {
	entry.Logf(WarnLevel, format, args...)
}

func (entry *Entry) Warningf(format string, args ...interface{}) {
	entry.Warnf(format, args...)
}

func (entry *Entry) Errorf(format string, args ...interface{}) {
	entry.Logf(ErrorLevel, format, args...)
}

func (entry *Entry) Fatalf(format string, args ...interface{}) {
	entry.Logf(FatalLevel, format, args...)
}

func (entry *Entry) Panicf(format string, args ...interface{}) {
	entry.Logf(PanicLevel, format, args...)
}

// Entry Println family functions

func (entry *Entry) Logln(level Level, args ...interface{}) {
	logrusLevel := logrus.Level(level)
	entry.stdEntry.Logln(logrusLevel, args...)

	if entry.stdEntry.Logger.IsLevelEnabled(logrusLevel) == false && entry.IsDebugMode() {
		entry.debugEntry.Logln(logrusLevel, args...)
	}
}

func (entry *Entry) Traceln(args ...interface{}) {
	entry.Logln(TraceLevel, args...)
}

func (entry *Entry) Debugln(args ...interface{}) {
	entry.Logln(DebugLevel, args...)
}

func (entry *Entry) Infoln(args ...interface{}) {
	entry.Logln(InfoLevel, args...)
}

func (entry *Entry) Println(args ...interface{}) {
	entry.Infoln(args...)
}

func (entry *Entry) Warnln(args ...interface{}) {
	entry.Logln(WarnLevel, args...)
}

func (entry *Entry) Warningln(args ...interface{}) {
	entry.Warnln(args...)
}

func (entry *Entry) Errorln(args ...interface{}) {
	entry.Logln(ErrorLevel, args...)
}

func (entry *Entry) Fatalln(args ...interface{}) {
	entry.Logln(FatalLevel, args...)
}

func (entry *Entry) Panicln(args ...interface{}) {
	entry.Logln(PanicLevel, args...)
}

func (entry *Entry) IsDebugMode() bool {
	if entry.stdEntry.Context == nil {
		return false
	}

	value := entry.stdEntry.Context.Value("__logger_debug_mode__")
	if value == nil {
		return false
	}

	if debugMode, ok := value.(bool); !ok {
		return false
	} else {
		return debugMode
	}
}

func SetDebugModeContext(ctx context.Context, mode bool) context.Context {
	if ctx == nil {
		return context.WithValue(context.Background(), "__logger_debug_mode__", mode)
	}

	return context.WithValue(ctx, "__logger_debug_mode__", mode)
}
