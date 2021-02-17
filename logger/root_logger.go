package logger

var rootLogger Logger

func SetupRootLogger(c Config) {
	rootLogger = newLogrusLogger(c)
}

func WithFields(fields Fields) Logger {
	if rootLogger == nil {
		panic("Logger not initialized")
	}

	return rootLogger.WithFields(fields)
}

func Debugf(format string, args ...interface{}) {
	WithFields(nil).Debugf(format, args...)
}

func Infof(format string, args ...interface{}) {
	WithFields(nil).Infof(format, args...)
}

func Warnf(format string, args ...interface{}) {
	WithFields(nil).Warnf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	WithFields(nil).Errorf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	WithFields(nil).Fatalf(format, args...)
}

func Debug(args ...interface{}) {
	WithFields(nil).Debug(args...)
}

func Info(args ...interface{}) {
	WithFields(nil).Info(args...)
}

func Warn(args ...interface{}) {
	WithFields(nil).Warn(args...)
}

func Error(args ...interface{}) {
	WithFields(nil).Error(args...)
}

func Fatal(args ...interface{}) {
	WithFields(nil).Fatal(args...)
}
