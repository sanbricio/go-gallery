package logger

type Logger interface {
	Info(msg string)
	Error(msg string)
	Warning(msg string)
	Panic(msg string)
}
