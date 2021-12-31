package cluster

import (
	"io"
	"regexp"

	"github.com/7vars/leikari"
)

var logRegex = regexp.MustCompile(`(?:\[(?P<logtype>[A-Z]+)\])\s(?P<msg>.+)`)

type logWrapper struct{
	log leikari.Logger
}

func newLogWrapper(log leikari.Logger) io.Writer {
	return &logWrapper{log}
}

func splitMessage(b []byte) (string, string) {
	match := logRegex.FindStringSubmatch(string(b))
	logtype := ""
	msg := ""
	for i, name := range logRegex.SubexpNames() {
		switch name {
		case "logtype":
			logtype = match[i]
		case "msg":
			msg = match[i]
		}
	}
	return logtype, msg
}

func (lw *logWrapper) Write(b []byte) (int, error) {
	ltype, msg := splitMessage(b)
	switch ltype {
	case "DEBUG":
		lw.log.Debug(msg)
	case "INFO":
		lw.log.Info(msg)
	case "ERR": 
		lw.log.Error(msg)
	default: 
		lw.log.Warn(msg)
	}
	return len(b), nil
}