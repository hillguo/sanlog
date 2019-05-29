package sanlog

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"runtime"
)

const (
	DEBUG LogLevel = iota
	INFO
	WARN
	ERROR
	OFF
)

var (
	logLevel LogLevel = DEBUG
	logQueue  = make(chan *logValue, 10000)
   	loggerMap = make(map[string]*Logger)
	writeDone = make(chan bool)
)


type LogLevel uint8

type Logger struct {
	name   string
	writer LogWriter
}

type logValue struct {
	level  LogLevel
	value  []byte
	fileNo string
	writer LogWriter
}

func init() {
	go flushLog(true)
}

func (lv *LogLevel) String() string {
	switch *lv {
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	default:
		return "UNKNOWN"
	}
}

func SetLevel(level LogLevel) {
	logLevel = level
}

func StringToLevel(level string) LogLevel {
	switch level {
	case "DEBUG":
		return DEBUG
	case "INFO":
		return INFO
	case "WARN":
		return WARN
	case "ERROR":
		return ERROR
	default:
		return DEBUG
	}
}

func (l *Logger) SetLogName(name string) {
	l.name = name
}

func (l *Logger) SetFileRoller(logpath string, num int, sizeMB int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		panic(err)
	}
	w := NewRollFileWriter(logpath, l.name, num, sizeMB)
	l.writer = w
	return nil
}

func (l *Logger) IsConsoleWriter() bool {
	if reflect.TypeOf(l.writer) == reflect.TypeOf(&ConsoleWriter{}) {
		return true
	}
	return false
}

func (l *Logger) SetWriter(w LogWriter) {
	l.writer = w
}

func (l *Logger) SetDayRoller(logpath string, num int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		return err
	}
	w := NewDateWriter(logpath, l.name, DAY, num)
	l.writer = w
	return nil
}

func (l *Logger) SetHourRoller(logpath string, num int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		return err
	}
	w := NewDateWriter(logpath, l.name, HOUR, num)
	l.writer = w
	return nil
}

func (l *Logger) SetConsole() {
	l.writer = &ConsoleWriter{}
}

func (l *Logger) Debug(v ...interface{}) {
	l.writef(DEBUG, "", v)
}

func (l *Logger) Info(v ...interface{}) {
	l.writef(INFO, "", v)
}

func (l *Logger) Warn(v ...interface{}) {
	l.writef(WARN, "", v)
}

func (l *Logger) Error(v ...interface{}) {
	l.writef(ERROR, "", v)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writef(DEBUG, format, v)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.writef(INFO, format, v)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.writef(WARN, format, v)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writef(ERROR, format, v)
}

func (l *Logger) writef(level LogLevel, format string, v []interface{}) {
	if level < logLevel {
		return
	}

	buf := bytes.NewBuffer(nil)
	if l.writer.NeedPrefix() {
		fmt.Fprintf(buf, "%s|", CurrDateTime)
		if logLevel == DEBUG {
			_, file, line, ok := runtime.Caller(2)
			if !ok {
				file = "???"
				line = 0
			} else {
				file = filepath.Base(file)
			}
			fmt.Fprintf(buf, "%s:%d|", file, line)
		}
	}
	buf.WriteString(level.String())
	buf.WriteByte('|')

	if format == "" {
		fmt.Fprint(buf, v...)
	} else {
		fmt.Fprintf(buf, format, v...)
	}
	if l.writer.NeedPrefix() {
		buf.WriteByte('\n')
	}
	logQueue <- &logValue{value: buf.Bytes(), writer: l.writer}
}

func FlushLogger() {
	flushLog(false)
}

func flushLog(sync bool) {
	if sync {
		for v := range logQueue {
			v.writer.Write(v.value)
		}
	} else {
		for {
			select {
			case v := <-logQueue:
				v.writer.Write(v.value)
				continue
			default:
				return
			}
		}
	}
}
