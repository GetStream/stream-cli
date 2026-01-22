package processing

import (
	"github.com/GetStream/getstream-go/v3"
)

type ProcessingLogger struct {
	logger *getstream.DefaultLogger
}

func NewRawToolLogger(logger *getstream.DefaultLogger) *ProcessingLogger {
	return &ProcessingLogger{
		logger: logger,
	}
}

func (l *ProcessingLogger) Debug(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

func (l *ProcessingLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debug(format, args...)
}

func (l *ProcessingLogger) Info(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

func (l *ProcessingLogger) Infof(format string, args ...interface{}) {
	l.logger.Info(format, args...)
}

func (l *ProcessingLogger) Warn(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

func (l *ProcessingLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warn(format, args...)
}

func (l *ProcessingLogger) Error(format string, args ...interface{}) {
	l.logger.Error(format, args...)
}

func (l *ProcessingLogger) Errorf(format string, args ...interface{}) {
	l.logger.Error(format, args...)
}
