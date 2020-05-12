package logger

import (
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

func newLogrusLogger(c Config) *logrus.Logger {
	level, err := logrus.ParseLevel(c.Level)
	if err != nil {
		level = logrus.WarnLevel
	}

	l := &logrus.Logger{
		Out:   os.Stdout,
		Hooks: make(logrus.LevelHooks),
		Level: level,
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

	return l
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
