package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

type logrusLogEntry struct {
	entry *logrus.Entry
}

type logrusLogger struct {
	logger *logrus.Logger
}

func newLogrusLogger(c Config) Logger {
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		level = logrus.WarnLevel
	}

	l := &logrus.Logger{
		Out:          os.Stdout,
		Hooks:        make(logrus.LevelHooks),
		Level:        level,
		ReportCaller: c.ReportCaller,
	}

	format, err := parseFormat(c.Format)
	if err != nil {
		format = TextFormat
	}

	if format == TextFormat {
		l.Formatter = &logrus.TextFormatter{
			FullTimestamp:          true,
			DisableLevelTruncation: true,
		}
	} else {
		l.Formatter = &logrus.JSONFormatter{}
	}

	return &logrusLogger{l}
}

func parseFormat(format string) (Format, error) {
	switch strings.ToLower(format) {
	case "text":
		return TextFormat, nil
	case "json":
		return JsonFormat, nil
	}

	return "", fmt.Errorf("Not a valid logger format : %s", format)
}

func (l *logrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *logrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *logrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *logrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *logrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *logrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *logrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *logrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *logrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *logrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *logrusLogger) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		entry: l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *logrusLogEntry) Debug(args ...interface{}) {
	l.entry.Debug(args...)
}

func (l *logrusLogEntry) Info(args ...interface{}) {
	l.entry.Info(args...)
}

func (l *logrusLogEntry) Warn(args ...interface{}) {
	l.entry.Warn(args...)
}

func (l *logrusLogEntry) Error(args ...interface{}) {
	l.entry.Error(args...)
}

func (l *logrusLogEntry) Fatal(args ...interface{}) {
	l.entry.Fatal(args...)
}

func (l *logrusLogEntry) Debugf(format string, args ...interface{}) {
	l.entry.Debugf(format, args...)
}

func (l *logrusLogEntry) Infof(format string, args ...interface{}) {
	l.entry.Infof(format, args...)
}

func (l *logrusLogEntry) Warnf(format string, args ...interface{}) {
	l.entry.Warnf(format, args...)
}

func (l *logrusLogEntry) Errorf(format string, args ...interface{}) {
	l.entry.Errorf(format, args...)
}

func (l *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) Panicf(format string, args ...interface{}) {
	l.entry.Fatalf(format, args...)
}

func (l *logrusLogEntry) WithFields(fields Fields) Logger {
	return &logrusLogEntry{
		l.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for key, val := range fields {
		logrusFields[key] = val
	}
	return logrusFields
}
