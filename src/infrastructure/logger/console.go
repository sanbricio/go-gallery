package logger

import (
	"fmt"
	loggerEntity "go-gallery/src/domain/entities/logger"
	"log"
	"time"

	"github.com/fatih/color"
)

type ConsoleLogger struct{}

func NewConsoleLogger() *ConsoleLogger {
	return &ConsoleLogger{}
}

func (l *ConsoleLogger) Info(msg string) {
	l.printLog(loggerEntity.INFO, msg)
}

func (l *ConsoleLogger) Error(msg string) {
	l.printLog(loggerEntity.ERROR, msg)
}

func (l *ConsoleLogger) Warning(msg string) {
	l.printLog(loggerEntity.WARNING, msg)
}

func (l *ConsoleLogger) Panic(msg string) {
	l.printLog(loggerEntity.PANIC, msg)
}

func (l *ConsoleLogger) printLog(level loggerEntity.LogLevel, msg string) {
	currentTime := time.Now().Format("02/01/06 15:04:05")
	var logColor func(format string, a ...any) string
	var prefix string

	switch level {
	case loggerEntity.INFO:
		prefix = loggerEntity.INFO.String()
		logColor = color.New(color.FgGreen).SprintfFunc()
	case loggerEntity.ERROR:
		prefix = loggerEntity.ERROR.String()
		logColor = color.New(color.FgRed).SprintfFunc()
	case loggerEntity.WARNING:
		prefix = loggerEntity.WARNING.String()
		logColor = color.New(color.FgYellow).SprintfFunc()
	case loggerEntity.PANIC:
		prefix = loggerEntity.PANIC.String()
		logColor = color.New(color.FgHiMagenta).SprintfFunc()
	}

	logMessage := fmt.Sprintf("%s: %s %s", prefix, currentTime, msg)
	log.Println(logColor("%s", logMessage))
}
