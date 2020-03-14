package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

const (
	NULL = iota
	TRACE = 1
	DEBUG = 2
	INFO = 3
	WARNGING = 4
	ERROR = 5
	FATAL = 6
)

type Log interface {
	Trace(format string, v ...interface{})
	Debug(format string, v ...interface{})
	Info(format string, v ...interface{})
	Warning(format string, v ...interface{})
	Error(format string, v ...interface{})
	Fatal(format string, v ...interface{})
}

type logger struct{
	*log.Logger
	options *Options
}

var DefaultLog = &logger {
	Logger : log.New(os.Stdout, "", log.LstdFlags|log.Lshortfile),
	options : &Options {
		level : 2,
	},
}

type Level int

func (level Level) String() string {
	switch level {
		case TRACE :
			return "trace"
		case DEBUG :
			return "debug"
		case INFO :
			return "info"
		case WARNGING :
			return "warning"
		case ERROR :
			return "error"
		case FATAL :
			return "fatal"
		default :
			return "unkown"
	}
}

type Options struct {
	path string `default:"../log/gorpc"`   // log file path prefix, filename : gorpc.2019-09-26.log
	frame string `default:"../log/frame"`  // frame log print path, default : ../log/frame.log
	level Level `default:"2"`         	   // log level, default: debug
}

type Option func(*Options)

func WithPath(path string) Option {
	return func(o *Options) {
		o.path = path
	}
}

func WithFrame(frame string) Option {
	return func(o *Options) {
		o.frame = frame
	}
}

func WithLevel(level Level) Option {
	return func(o *Options) {
		o.level = level
	}
}

func Trace(format string, v ...interface{}) {
	DefaultLog.Trace(format, v...)
}

func (log *logger) Trace(format string, v ...interface{}) {
	if log.options.level > TRACE {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[TRACE] ")
	buffer.WriteString(data)
	log.Output(3, buffer.String())
}

func Debug(format string, v ...interface{}) {
	DefaultLog.Debug(format, v...)
}

func (log *logger) Debug(format string, v ...interface{}) {
	if log.options.level > DEBUG {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[DEBUG] ")
	buffer.WriteString(data)
	log.Output(3, buffer.String())
}

func Info(format string, v ...interface{}) {
	DefaultLog.Info(format, v...)
}

func (log *logger) Info(format string, v ...interface{}) {
	if log.options.level > INFO {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[INFO] ")
	buffer.WriteString(data)
	log.Output(3, buffer.String())
}

func Warning(format string, v ...interface{}) {
	DefaultLog.Warning(format, v...)
}

func (log *logger) Warning(format string, v ...interface{}) {
	if log.options.level > WARNGING {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[WARNING] ")
	buffer.WriteString(data)
	log.Output(3, buffer.String())
}

func Error(format string, v ...interface{}) {
	DefaultLog.Error(format, v...)
}

func (log *logger) Error(format string, v ...interface{}) {
	if log.options.level > ERROR {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[ERROR] ")
	buffer.WriteString(data)
	log.Output(3, buffer.String())
}

func Fatal(format string, v ...interface{}) {
	DefaultLog.Fatal(format, v...)
}

func (log *logger) Fatal(format string, v ...interface{}) {
	if log.options.level > FATAL {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[FATAL] ")
	buffer.WriteString(data)
	log.Output(3, buffer.String())
}

