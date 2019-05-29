package sanlog

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"
)
const (
	DAY DateType = iota
	HOUR
)

type DateType uint8

type LogWriter interface {
	Write(v []byte)
	NeedPrefix() bool
}

type ConsoleWriter struct {

}

func (w *ConsoleWriter) Write(v []byte) {
	os.Stdout.Write(v)
}

func (w *ConsoleWriter) NeedPrefix() bool {
	return true
}

type RollFileWriter struct {
	logpath  string
	name     string
	num      int
	size     int64
	currSize int64
	currFile *os.File
	openTime int64
}


func (w *RollFileWriter) Write(v []byte) {
	if w.currFile == nil || w.openTime+10 < CurrUnixTime {
		fullPath := filepath.Join(w.logpath, w.name+".log")
		reOpenFile(fullPath, &w.currFile, &w.openTime)
	}
	if w.currFile == nil {
		return
	}
	n, _ := w.currFile.Write(v)
	w.currSize += int64(n)
	if w.currSize >= w.size {
		w.currSize = 0
		for i := w.num - 1; i >= 1; i-- {
			var n1, n2 string
			if i > 1 {
				n1 = strconv.Itoa(i - 1)
			}
			n2 = strconv.Itoa(i)
			p1 := filepath.Join(w.logpath, w.name+n1+".log")
			p2 := filepath.Join(w.logpath, w.name+n2+".log")
			if _, err := os.Stat(p1); !os.IsNotExist(err) {
				os.Rename(p1, p2)
			}
		}
		fullPath := filepath.Join(w.logpath, w.name+".log")
		reOpenFile(fullPath, &w.currFile, &w.openTime)
	}
}

func (w *RollFileWriter) NeedPrefix() bool {
	return true
}

func reOpenFile(path string, currFile **os.File, openTime *int64) {
	*openTime = CurrUnixTime
	if *currFile != nil {
		(*currFile).Close()
	}
	of, err := os.OpenFile(path, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	if err == nil {
		*currFile = of
	} else {
		fmt.Println("open log file error", err)
	}
}

type DateWriter struct {
	logpath  string
	name     string
	dateType DateType
	num      int
	currDate string
	currFile *os.File
	openTime int64
}

func (w *DateWriter) Write(v []byte) {
	if w.isExpired() {
		w.currDate = w.getCurrDate()
		w.cleanOldLogs()
		fullPath := filepath.Join(w.logpath, w.name+"_"+w.currDate+".log")
		reOpenFile(fullPath, &w.currFile, &w.openTime)
	}
	if w.currFile == nil || w.openTime+10 < CurrUnixTime {
		fullPath := filepath.Join(w.logpath, w.name+"_"+w.currDate+".log")
		reOpenFile(fullPath, &w.currFile, &w.openTime)
	}
	if w.currFile == nil {
		return
	}
	w.currFile.Write(v)
}

func (w *DateWriter) NeedPrefix() bool {
	return true
}

func (w *DateWriter) cleanOldLogs() {
	format := "20060102"
	duration := -time.Hour * 24
	if w.dateType == HOUR {
		format = "2006010215"
		duration = -time.Hour
	}

	t := time.Now()
	t = t.Add(duration * time.Duration(w.num))
	for i := 0; i < 30; i++ {
		t = t.Add(duration)
		k := t.Format(format)
		fullPath := filepath.Join(w.logpath, w.name+"_"+k+".log")
		if _, err := os.Stat(fullPath); !os.IsNotExist(err) {
			os.Remove(fullPath)
		}
	}
	return
}

func (w *DateWriter) getCurrDate() string {
	if w.dateType == HOUR {
		return CurrDateHour
	}
	return CurrDateDay // DAY
}

func (w *DateWriter) isExpired() bool {
	currDate := w.getCurrDate()
	return w.currDate != currDate
}