package logger

type CompositeLogger struct {
	loggers []Logger
}

func NewCompositeLogger(loggers ...Logger) *CompositeLogger {
	return &CompositeLogger{
		loggers: loggers,
	}
}

func (cl *CompositeLogger) Info(msg string) {
	for _, logger := range cl.loggers {
		logger.Info(msg)
	}
}

func (cl *CompositeLogger) Error(msg string) {
	for _, logger := range cl.loggers {
		logger.Error(msg)
	}
}

func (cl *CompositeLogger) Warning(msg string) {
	for _, logger := range cl.loggers {
		logger.Warning(msg)
	}
}

func (cl *CompositeLogger) Panic(msg string) {
	for _, logger := range cl.loggers {
		logger.Panic(msg)
	}
}
