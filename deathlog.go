package death

// Logger interface to log.
type Logger interface {
	Error(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
}

type defaultLogger struct{}

var logger = defaultLogger{}

// DefaultLogger returns a logger that does nothing
func DefaultLogger() Logger {
	return logger
}

func (d defaultLogger) Error(v ...interface{}) {}
func (d defaultLogger) Debug(v ...interface{}) {}
func (d defaultLogger) Info(v ...interface{})  {}
