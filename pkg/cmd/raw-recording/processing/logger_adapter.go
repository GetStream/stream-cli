package processing

import (
	"fmt"
	"io"
)

// LogLevel represents the severity level of a log message
type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarn
	LogLevelError
)

type ProcessingLogger struct {
	writer io.Writer
	level  LogLevel
}

func NewProcessingLogger(writer io.Writer, level LogLevel) *ProcessingLogger {
	return &ProcessingLogger{
		writer: writer,
		level:  level,
	}
}

func (l *ProcessingLogger) log(level LogLevel, prefix, format string, args ...interface{}) {
	if level < l.level {
		return
	}
	msg := fmt.Sprintf(format, args...)
	_, _ = fmt.Fprintf(l.writer, "%s %s\n", prefix, msg)
}

func (l *ProcessingLogger) Debug(format string, args ...interface{}) {
	l.log(LogLevelDebug, "[DEBUG]", format, args...)
}

func (l *ProcessingLogger) Debugf(format string, args ...interface{}) {
	l.log(LogLevelDebug, "[DEBUG]", format, args...)
}

func (l *ProcessingLogger) Info(format string, args ...interface{}) {
	l.log(LogLevelInfo, "[INFO]", format, args...)
}

func (l *ProcessingLogger) Infof(format string, args ...interface{}) {
	l.log(LogLevelInfo, "[INFO]", format, args...)
}

func (l *ProcessingLogger) Warn(format string, args ...interface{}) {
	l.log(LogLevelWarn, "[WARN]", format, args...)
}

func (l *ProcessingLogger) Warnf(format string, args ...interface{}) {
	l.log(LogLevelWarn, "[WARN]", format, args...)
}

func (l *ProcessingLogger) Error(format string, args ...interface{}) {
	l.log(LogLevelError, "[ERROR]", format, args...)
}

func (l *ProcessingLogger) Errorf(format string, args ...interface{}) {
	l.log(LogLevelError, "[ERROR]", format, args...)
}
