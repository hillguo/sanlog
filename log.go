package sanlog

const (
	calldepth = 4
)

func init (){
	l.SetCallerSkip(3)
	l.SetLevel(DEBUG)
}

var l = GetLogger("sanlog")

//SetLogger ...
func SetLogger(logger *Logger) {
	l = logger
}

//SetLevel ...
func SetLevel(level LogLevel) {
	l.level = level
}

//Debug ...
func Debug(v ...interface{}) {
	l.Debug(v...)
}

//Debugf ...
func Debugf(format string, v ...interface{}) {
	l.Debugf(format, v...)
}

//Info ...
func Info(v ...interface{}) {
	l.Info(v...)
}

//Infof ...
func Infof(format string, v ...interface{}) {
	l.Infof(format, v...)
}

//Warn ...
func Warn(v ...interface{}) {
	l.Warn(v...)
}

//Warnf ...
func Warnf(format string, v ...interface{}) {
	l.Warnf(format, v...)
}

//Error ...
func Error(v ...interface{}) {
	l.Error(v...)
}

//Errorf ...
func Errorf(format string, v ...interface{}) {
	l.Errorf(format, v...)
}

//Fatal ...
func Fatal(v ...interface{}) {
	l.Fatal(v...)
}

//Fatalf ...
func Fatalf(format string, v ...interface{}) {
	l.Fatalf(format, v...)
}

//Panic ...
func Panic(v ...interface{}) {
	l.Panic(v...)
}

//Panicf ...
func Panicf(format string, v ...interface{}) {
	l.Panicf(format, v...)
}
