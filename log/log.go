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
	INFO(format string, v ...interface{})
	WARNING(format string, v ...interface{})
	ERROR(format string, v ...interface{})
	FATAL(format string, v ...interface{})
}

type logger struct{
	*log.Logger
	options *Options
}

var defaultLog = &logger {
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
	path string `default:"../log/gorpc"`   // 日志文件路径前缀，文件名为 gorpc.2019-09-26.log
	frame string `default:"../log/frame"`  // 框架日志打印路径，默认 ../log/frame.log
	level Level `default:"2"`          // 日志级别，默认为 debug
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
	defaultLog.Trace(format, v...)
}

func (log *logger) Trace(format string, v ...interface{}) {
	if log.options.level > TRACE {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[TRACE] ")
	buffer.WriteString(data)

}

func Debug(format string, v ...interface{}) {
	defaultLog.Debug(format, v...)
}

func (log *logger) Debug(format string, v ...interface{}) {
	if log.options.level > DEBUG {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[DEBUG] ")
	buffer.WriteString(data)
	Output(log, 2, buffer.String())
}

func Info(format string, v ...interface{}) {
	defaultLog.Info(format, v...)
}

func (log *logger) Info(format string, v ...interface{}) {
	if log.options.level > INFO {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[INFO] ")
	buffer.WriteString(data)
	Output(log, 2, buffer.String())
}

func Warning(format string, v ...interface{}) {
	defaultLog.Warning(format, v...)
}

func (log *logger) Warning(format string, v ...interface{}) {
	if log.options.level > WARNGING {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[WARNING] ")
	buffer.WriteString(data)
	Output(log, 2, buffer.String())
}

func Error(format string, v ...interface{}) {
	defaultLog.Error(format, v...)
}

func (log *logger) Error(format string, v ...interface{}) {
	if log.options.level > ERROR {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[ERROR] ")
	buffer.WriteString(data)
	Output(log, 2, buffer.String())
}

func Fatal(format string, v ...interface{}) {
	defaultLog.Fatal(format, v...)
}

func (log *logger) Fatal(format string, v ...interface{}) {
	if log.options.level > FATAL {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	var buffer bytes.Buffer
	buffer.WriteString("[FATAL] ")
	buffer.WriteString(data)
	Output(log, 2, buffer.String())
}

func Output(log *logger, calldepth int, s string) {
	if err := log.Output(2, s); err != nil {
		fmt.Println("log output fail...")
	}
}
