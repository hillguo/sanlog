package main

import "github.com/hillguo/sanlog"

func main(){
	l :=sanlog.GetLogger("test")
	l.Debug("test")
	l.Info("test")
	l.Error("test")
	for {

	}
}
