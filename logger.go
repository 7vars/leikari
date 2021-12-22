package leikari

import (
	"fmt"
	"log"
	"os"
)

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
	logger log.Logger
}

func newLogger() Logger {
	return newSysLogger(*log.Default())
}

func newSysLogger(logger log.Logger) Logger {
	return &sysLogger{
		logger: logger,
	}
}

func (l *sysLogger) ForName(name string) Logger {
	return newSysLogger(*log.New(os.Stderr, name+" ", l.logger.Flags()))
}

func (l *sysLogger) Debug(v ...interface{}) {
	l.logger.Println(append([]interface{}{"[DEBUG]"}, v...)...)
}

func (l *sysLogger) Info(v ...interface{}) {
	l.logger.Println(append([]interface{}{"[INFO]"}, v...)...)
}

func (l *sysLogger) Warn(v ...interface{}) {
	l.logger.Println(append([]interface{}{"[WARN]"}, v...)...)
}

func (l *sysLogger) Error(v ...interface{}) {
	l.logger.Println(append([]interface{}{"[ERROR]"}, v...)...)
}

func (l *sysLogger) Fatal(v ...interface{}) {
	l.logger.Fatalln(append([]interface{}{"[FATAL]"}, v...)...)
}

func (l *sysLogger) Panic(v ...interface{}) {
	l.logger.Panicln(append([]interface{}{"[PANIC]"}, v...)...)
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