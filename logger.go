package sanlog

import (
	"bytes"
	"fmt"
	"github.com/fatih/color"
	"os"
	"path/filepath"
	"runtime"
)

var (
	logQueue = make(chan *logValue, 10000)
)

func init() {
	go flushLog()
}

//Logger ...
type Logger struct {
	name       string
	writer     LogWriter
	level      LogLevel
	callerSkip int
}

type logValue struct {
	level  LogLevel
	value  []byte
	fileNo string
	writer LogWriter
}

func (l *Logger) SetCallerSkip(skip int) {
	l.callerSkip = skip
}

//SetLevel ...
func (l *Logger) SetLevel(level LogLevel) {
	l.level = level
}

//SetLogName ...
func (l *Logger) SetLogName(name string) {
	l.name = name
}

//SetWriter ...
func (l *Logger) SetWriter(w LogWriter) {
	l.writer = w
}

//SetFileRoller ...
func (l *Logger) SetFileRoller(logpath string, num int, sizeMB int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		panic(err)
	}
	w := NewRollFileWriter(logpath, l.name, num, sizeMB)
	l.writer = w
	return nil
}

//SetDayRoller ...
func (l *Logger) SetDayRoller(logpath string, num int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		return err
	}
	w := NewDateWriter(logpath, l.name, DAY, num)
	l.writer = w
	return nil
}

//SetHourRoller ...
func (l *Logger) SetHourRoller(logpath string, num int) error {
	if err := os.MkdirAll(logpath, 0755); err != nil {
		return err
	}
	w := NewDateWriter(logpath, l.name, HOUR, num)
	l.writer = w
	return nil
}

//Debug ...
func (l *Logger) Debug(v ...interface{}) {
	l.writef(DEBUG, "", v)
}

//Info ...
func (l *Logger) Info(v ...interface{}) {
	l.writef(INFO, "", v)
}

//Warn ...
func (l *Logger) Warn(v ...interface{}) {
	l.writef(WARN, "", v)
}

//Fatal ...
func (l *Logger) Fatal(v ...interface{}) {
	l.writef(FATAL, "", v)
}

//Error ...
func (l *Logger) Error(v ...interface{}) {
	l.writef(ERROR, "", v)
}

//Debugf ...
func (l *Logger) Debugf(format string, v ...interface{}) {
	l.writef(DEBUG, format, v)
}

//Infof ...
func (l *Logger) Infof(format string, v ...interface{}) {
	l.writef(INFO, format, v)
}

//Warnf ...
func (l *Logger) Warnf(format string, v ...interface{}) {
	l.writef(WARN, format, v)
}

//Errorf ...
func (l *Logger) Errorf(format string, v ...interface{}) {
	l.writef(ERROR, format, v)
}

//Fatalf ...
func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.writef(FATAL, format, v)
}

// Panic is equivalent to l.Print() followed by a call to panic().
func (l *Logger) Panic(v ...interface{}) {
	l.writef(PANIC, "", v)
	panic(v)
}

// Panicf is equivalent to l.Printf() followed by a call to panic().
func (l *Logger) Panicf(format string, v ...interface{}) {
	l.writef(PANIC, format, v)
	panic(fmt.Sprintf(format, v...))
}

func (l *Logger) writef(level LogLevel, format string, v []interface{}) {
	if level < l.level {
		return
	}

	if l.callerSkip == 0 {
		l.callerSkip = 2
	}

	buf := bytes.NewBuffer(nil)

	fmt.Fprintf(buf, "%s|", CurrDateTime)

	_, file, line, ok := runtime.Caller(l.callerSkip)
	if !ok {
		file = "???"
		line = 0
	} else {
		file = filepath.Base(file)
	}
	fmt.Fprintf(buf, "%s:%d|", file, line)

	buf.WriteString(level.String())
	buf.WriteByte('|')

	if format == "" {
		fmt.Fprint(buf, v...)
	} else {
		fmt.Fprintf(buf, format, v...)
	}

	buf.WriteByte('\n')

	msg := withColor(level, buf.String())

	logQueue <- &logValue{value: []byte(msg), writer: l.writer}
}

func flushLog() {
	for v := range logQueue {
		if v.writer != nil {
			v.writer.Write(v.value)
		}
	}
}

func withColor(level LogLevel, msg string) string {
	switch level {
	case DEBUG:
		return msg
	case INFO:
		return color.GreenString(msg)
	case WARN:
		return color.YellowString(msg)
	case ERROR:
		return color.RedString(msg)
	case FATAL:
		return color.RedString(msg)
	case PANIC:
		return color.RedString(msg)
	default:
		return msg
	}
	return msg
}
