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

// general log interface for gorpc
type Log interface {
	Trace(v ...interface{})
	Debug(v ...interface{})
	Info(v ...interface{})
	Warning(v ...interface{})
	Error(v ...interface{})
	Fatal(v ...interface{})
	Tracef(format string, v ...interface{})
	Debugf(format string, v ...interface{})
	Infof(format string, v ...interface{})
	Warningf(format string, v ...interface{})
	Errorf(format string, v ...interface{})
	Fatalf(format string, v ...interface{})

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
	return "unknown"
}

type Options struct {
	path string `default:"../log/gorpc"`   // 日志文件路径前缀，文件名为 gorpc.4019-09-46.log
	frame string `default:"../log/frame"`  // 框架日志打印路径，默认 ../log/frame.log
	level Level `defalut:"debug"`          // 日志级别，默认为 debug
}

type Option func(*Options)

// set the log path
func WithPath(path string) Option {
	return func(o *Options) {
		o.path = path
	}
}

// set the frame log path
func WithFrame(frame string) Option {
	return func(o *Options) {
		o.frame = frame
	}
}

// set the log level
func WithLevel(level Level) Option {
	return func(o *Options) {
		o.level = level
	}
}

// Trace print trace log
func Trace(v ...interface{}) {
	DefaultLog.Trace(v...)
}

// Tracef print a formatted trace log
func Tracef(format string, v ...interface{}) {
	DefaultLog.Tracef(format, v...)
}

func (log *logger) Trace(v ...interface{}) {
	if log.options.level < TRACE {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4,"[TRACE] ", data)
}

func (log *logger) Tracef(format string, v ...interface{}) {
	if log.options.level < TRACE {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	Output(log, 4,"[TRACE] ", data)
}

// Debug print debug log
func Debug(v ...interface{}) {
	DefaultLog.Debug(v...)
}

// Debugf print a formatted debug log
func Debugf(format string, v ...interface{}) {
	DefaultLog.Debugf(format, v...)
}

func (log *logger) Debug(v ...interface{}) {
	if log.options.level < DEBUG {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4,"[DEBUG] ", data)
}

func (log *logger) Debugf(format string, v ...interface{}) {
	if log.options.level < DEBUG {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	Output(log, 4,"[DEBUG] ", data)
}

// Info print info log
func Info(v ...interface{}) {
	DefaultLog.Info(v...)
}

// Infof print a formatted info log
func Infof(format string, v ...interface{}) {
	DefaultLog.Infof(format, v...)
}

func (log *logger) Info(v ...interface{}) {
	if log.options.level < INFO {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4,"[INFO] ", data)
}

func (log *logger) Infof(format string, v ...interface{}) {
	if log.options.level < INFO {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	Output(log, 4,"[INFO] ", data)
}

// Warning print warning log
func Warning(v ...interface{}) {
	DefaultLog.Warning(v...)
}

// Warningf print a formatted warning log
func Warningf(format string, v ...interface{}) {
	DefaultLog.Warningf(format, v...)
}

func (log *logger) Warning(v ...interface{}) {
	if log.options.level < WARNGING {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4,"[WARNING] ", data)
}

func (log *logger) Warningf(format string, v ...interface{}) {
	if log.options.level < WARNGING {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	Output(log, 4,"[WARNING] ", data)
}

// Error print error log
func Error(v ...interface{}) {
	DefaultLog.Error(v...)
}

// Errorf print a formatted error log
func Errorf(format string, v ...interface{}) {
	DefaultLog.Errorf(format, v...)
}

func (log *logger) Error(v ...interface{}) {
	if log.options.level < ERROR {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4,"[ERROR] ", data)
}

func (log *logger) Errorf(format string, v ...interface{}) {
	if log.options.level < ERROR {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	Output(log, 4,"[ERROR] ", data)
}

// Fatal print fatal log
func Fatal(v ...interface{}) {
	DefaultLog.Fatal(v...)
}

// Fatalf print a formatted fatal log
func Fatalf(format string, v ...interface{}) {
	DefaultLog.Fatalf(format, v...)
}

func (log *logger) Fatal(v ...interface{}) {
	if log.options.level < FATAL {
		return
	}
	data := log.Prefix() + fmt.Sprint(v...)
	Output(log, 4,"[FATAL] ", data)
}

func (log *logger) Fatalf(format string, v ...interface{}) {
	if log.options.level < FATAL {
		return
	}
	data := log.Prefix() + fmt.Sprintf(format,v...)
	Output(log, 4,"[FATAL] ", data)
}

// call Output to write log
func Output(log *logger, calldepth int, prefix string, data string) {
	var buffer bytes.Buffer
	buffer.WriteString(prefix)
	buffer.WriteString(data)
	log.Output(calldepth, buffer.String())
}
