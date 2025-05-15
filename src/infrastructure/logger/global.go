package logger

import "sync"

var (
	logger Logger
	once   sync.Once
)

type AsyncGlobalLogger struct {
	inner Logger
	queue chan func()
}

func Init(inner Logger) Logger {
	once.Do(func() {
		agl := &AsyncGlobalLogger{inner: inner, queue: make(chan func(), 1000)}
		go agl.startWorker()
		logger = agl
	})
	return logger
}

func Instance() Logger {
	if logger == nil {
		panic("Logger is not initialized. Call logger.Init() first.")
	}
	return logger
}

func (a *AsyncGlobalLogger) startWorker() {
	for fn := range a.queue {
		fn()
	}
}

func (a *AsyncGlobalLogger) Info(msg string) {
	a.queue <- func() {
		a.inner.Info(msg)
	}
}

func (a *AsyncGlobalLogger) Warning(msg string) {
	a.queue <- func() {
		a.inner.Warning(msg)
	}
}

func (a *AsyncGlobalLogger) Error(msg string) {
	a.queue <- func() {
		a.inner.Error(msg)
	}
}

func (a *AsyncGlobalLogger) Panic(msg string) {
	a.queue <- func() {
		a.inner.Panic(msg)
	}
}
