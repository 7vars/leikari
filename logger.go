package leikari

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type LogLevel int

const (
	LEVEL_DEBUG LogLevel = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_PANIC
)

func logLevel(level string) LogLevel {
	switch (strings.ToUpper(level)) {
	case "DEBUG":
		return LEVEL_DEBUG
	case "WARN": 
		return LEVEL_WARN
	case "ERROR":
		return LEVEL_ERROR
	case "FATAL": 
		return LEVEL_FATAL
	case "PANIC":
		return LEVEL_PANIC
	}
	return LEVEL_INFO
}

type Logger interface {
	ForName(string) Logger

	Debug(...interface{})
	Info(...interface{})
	Warn(...interface{})
	Error(...interface{})
	Fatal(...interface{})
	Panic(...interface{})

	Debugf(string, ...interface{})
	Infof(string, ...interface{})
	Warnf(string, ...interface{})
	Errorf(string, ...interface{})
	Fatalf(string, ...interface{})
	Panicf(string, ...interface{})
}

type emptyLogger struct {}

func Empty() Logger {
	return &emptyLogger{}
}

func (el emptyLogger) ForName(string) Logger { return el }
func (emptyLogger) Debug(...interface{}) {}
func (emptyLogger) Info(...interface{}) {}
func (emptyLogger) Warn(...interface{}) {}
func (emptyLogger) Error(...interface{}) {}
func (emptyLogger) Fatal(...interface{}) {}
func (emptyLogger) Panic(...interface{}) {}
func (emptyLogger) Debugf(string, ...interface{}) {}
func (emptyLogger) Infof(string, ...interface{}) {}
func (emptyLogger) Warnf(string, ...interface{}) {}
func (emptyLogger) Errorf(string, ...interface{}) {}
func (emptyLogger) Fatalf(string, ...interface{}) {}
func (emptyLogger) Panicf(string, ...interface{}) {}


type sysLogger struct{
	name string
	level LogLevel
	logger log.Logger
}

func newLogger(loglevel LogLevel) Logger {
	return newSysLogger("", loglevel, *log.Default())
}

func newSysLogger(name string, loglevel LogLevel, logger log.Logger) Logger {
	return &sysLogger{
		name: name,
		level: loglevel,
		logger: logger,
	}
}

func (l *sysLogger) ForName(name string) Logger {
	return newSysLogger(name, l.level, *log.New(os.Stderr, "", l.logger.Flags()))
}

func (l *sysLogger) appendPrefix(vals []interface{}, prefix ...interface{}) []interface{} {
	var result []interface{}
	result = append(result, prefix...)
	if l.name != "" {
		result = append(result, fmt.Sprintf("(%v)", l.name))
	}
	result = append(result, vals...)
	return result
}

func (l *sysLogger) Debug(v ...interface{}) {
	if l.level > LEVEL_DEBUG {
		return
	}
	l.logger.Println(l.appendPrefix(v, "[DEBUG]")...)
}

func (l *sysLogger) Info(v ...interface{}) {
	if l.level > LEVEL_INFO {
		return
	}
	l.logger.Println(l.appendPrefix(v, "[INFO] ")...)
}

func (l *sysLogger) Warn(v ...interface{}) {
	if l.level > LEVEL_WARN {
		return
	}
	l.logger.Println(l.appendPrefix(v, "[WARN] ")...)
}

func (l *sysLogger) Error(v ...interface{}) {
	if l.level > LEVEL_ERROR {
		return
	}
	l.logger.Println(l.appendPrefix(v, "[ERROR]")...)
}

func (l *sysLogger) Fatal(v ...interface{}) {
	if l.level > LEVEL_FATAL {
		return
	}
	l.logger.Println(l.appendPrefix(v, "[FATAL]")...)
}

func (l *sysLogger) Panic(v ...interface{}) {
	l.logger.Println(l.appendPrefix(v, "[PANIC]")...)
}

func (l *sysLogger) Debugf(format string, v ...interface{}) {
	l.Debug(fmt.Sprintf(format, v...))
}

func (l *sysLogger) Infof(format string, v ...interface{}) {
	l.Info(fmt.Sprintf(format, v...))
}

func (l *sysLogger) Warnf(format string, v ...interface{}) {
	l.Warn(fmt.Sprintf(format, v...))
}

func (l *sysLogger) Errorf(format string, v ...interface{}) {
	l.Error(fmt.Sprintf(format, v...))
}

func (l *sysLogger) Fatalf(format string, v ...interface{}) {
	l.Fatal(fmt.Sprintf(format, v...))
}

func (l *sysLogger) Panicf(format string, v ...interface{}) {
	l.Panic(fmt.Sprintf(format, v...))
}