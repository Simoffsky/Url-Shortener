package log

import (
	"fmt"
	"strings"
	"time"
)

type Logger interface {
	Debug(string)
	Info(string)
	Warning(string)
	Error(string)
}

type LoggerLevel int

const (
	Debug LoggerLevel = iota
	Info
	Warning
	Error
)

func LevelFromString(level string) LoggerLevel {
	level = strings.ToUpper(level)
	switch level {
	case "DEBUG":
		return Debug
	case "INFO":
		return Info
	case "WARNING":
		return Warning
	case "ERROR":
		return Error
	default:
		fmt.Printf("[LOG]Unknown log level \"%s\", setting to DEBUG level\n", level)
		return Debug
	}
}

type defaultLogger struct {
	level      LoggerLevel
	prefix     string
	timeLayout string
}

func NewDefaultLogger(level LoggerLevel) *defaultLogger {
	return &defaultLogger{
		level:      level,
		prefix:     "",
		timeLayout: time.DateTime,
	}
}

func (l *defaultLogger) WithPrefix(prefix string) *defaultLogger {
	return &defaultLogger{
		level:      l.level,
		prefix:     l.prefix + prefix,
		timeLayout: time.DateTime,
	}
}

func (l *defaultLogger) WithTimePrefix(timeLayout string) *defaultLogger {
	return &defaultLogger{
		level:      l.level,
		prefix:     l.prefix,
		timeLayout: timeLayout,
	}
}

func (l *defaultLogger) logMessage(levelName, msg string) {
	fmt.Printf("%s [%s] %s: %s\n", l.prefix, time.Now().Format(l.timeLayout), levelName, msg)
}

func (l *defaultLogger) SetLevel(level LoggerLevel) {
	l.level = level
}

func (l *defaultLogger) Info(msg string) {
	if l.level <= Info {
		l.logMessage("INFO", msg)
	}
}

func (l *defaultLogger) Debug(msg string) {
	if l.level <= Debug {
		l.logMessage("DEBUG", msg)
	}
}

func (l *defaultLogger) Warning(msg string) {
	if l.level <= Warning {
		l.logMessage("WARNING", msg)

	}
}

func (l *defaultLogger) Error(msg string) {
	if l.level <= Error {
		l.logMessage("ERROR", msg)
	}
}
