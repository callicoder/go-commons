package logger

type Format string

const (
	TextFormat Format = "text"
	JsonFormat Format = "json"
)

type Config struct {
	Level        string
	Format       string
	ReportCaller bool
}

type Logger interface {
	Debug(args ...interface{})
	Info(args ...interface{})
	Warn(args ...interface{})
	Error(args ...interface{})
	Fatal(args ...interface{})

	Debugf(format string, args ...interface{})
	Infof(format string, args ...interface{})
	Warnf(format string, args ...interface{})
	Errorf(format string, args ...interface{})
	Fatalf(format string, args ...interface{})

	WithFields(fields Fields) Logger
}

type Fields map[string]interface{}
