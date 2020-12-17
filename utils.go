package sanlog

import (
	"os"
	"path/filepath"
)

var loggerMap = make(map[string]*Logger)

// GetLogger return an logger instance
func GetLogger(name string) *Logger {
	if lg, ok := loggerMap[name]; ok {
		return lg
	}
	lg := &Logger{
		name:   name,
		writer: &ConsoleWriter{},
		lFlags: Ldate | Lfile | Llevel | Lcolor,
	}
	loggerMap[name] = lg
	return lg
}

//NewRollFileWriter ...
func NewRollFileWriter(logpath, name string, num, sizeMB int) *RollFileWriter {
	w := &RollFileWriter{
		logpath: logpath,
		name:    name,
		num:     num,
		size:    int64(sizeMB) * 1024 * 1024,
	}
	fullPath := filepath.Join(logpath, name+".log")
	st, _ := os.Stat(fullPath)
	if st != nil {
		w.currSize = st.Size()
	}
	return w
}

//NewDateWriter ...
func NewDateWriter(logpath, name string, dateType DateType, num int) *DateWriter {
	w := &DateWriter{
		logpath:  logpath,
		name:     name,
		num:      num,
		dateType: dateType,
	}
	w.currDate = w.getCurrDate()
	return w
}
